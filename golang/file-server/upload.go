package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/pkg/errors"
)

type UploadInfo struct {
	Filename string
	Size     int64
	ChunkNum int
	ChunkMap map[int]bool
}

var (
	uploadTempDir = "./upload-temp"
	uploadRootDir = "./upload-files"
)

func UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		err := errors.New("invalid request method. only POST is allowed")
		renderError(w, err)
		return
	}

	// 解析上传参数
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		renderError(w, errors.Wrap(err, "parse multipart form failed"))
		return
	}
	filename := r.FormValue("filename")
	size, _ := strconv.ParseInt(r.FormValue("size"), 10, 64)
	chunkNum, _ := strconv.Atoi(r.FormValue("chunknum"))
	chunkIndex, _ := strconv.Atoi(r.FormValue("chunkindex"))

	log.Printf("filename: %s, size: %d, chunkNum: %d, chunkIndex: %d", filename, size, chunkNum, chunkIndex)
	if filename == "" || size == 0 || chunkNum == 0 || chunkIndex == 0 {
		renderError(w, errors.New("invalid request params"))
		return
	}

	// 如果是分片的第一块，则创建临时文件进行上传
	tmpFilename := filepath.Join(uploadTempDir, filename)
	if chunkIndex == 1 {
		os.MkdirAll(uploadTempDir, os.ModePerm)
		os.Remove(tmpFilename)
		tempFile, err := os.Create(tmpFilename)
		if err != nil {
			renderError(w, errors.Wrap(err, "create temp file failed"))
			return
		}
		tempFile.Close()
	}

	// 将分片数据写入临时文件
	tempFile, err := os.OpenFile(tmpFilename, os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		renderError(w, errors.Wrap(err, fmt.Sprintf("open tempfile(%s) failed", tmpFilename)))
		return
	}
	defer tempFile.Close()
	file, _, err := r.FormFile("file")
	if err != nil {
		renderError(w, errors.Wrap(err, "get file from request failed"))
		return
	}
	defer file.Close()
	_, err = io.Copy(tempFile, file)
	if err != nil {
		renderError(w, errors.Wrap(err, "write file to temp file failed"))
		return
	}

	// 更新已上传分片的标记
	chunkMap := make(map[int]bool)
	if chunkIndex != chunkNum {
		chunkMap[chunkIndex] = true
	} else {
		for i := 1; i <= chunkNum; i++ {
			chunkMap[i] = true
		}
	}

	// 保存已上传分片信息
	uploadInfo := UploadInfo{
		Filename: filename,
		Size:     size,
		ChunkNum: chunkNum,
		ChunkMap: chunkMap,
	}
	uploadInfoBytes, _ := json.Marshal(uploadInfo)
	tempFileInfo := filepath.Join(uploadTempDir, filename+".info")
	err = os.WriteFile(tempFileInfo, uploadInfoBytes, os.ModePerm)
	if err != nil {
		renderError(w, errors.Wrap(err, "write upload info to temp file failed"))
		return
	}

	// 如果所有分片都上传完成，则将临时文件重命名为正式文件
	if len(chunkMap) == chunkNum {
		if err := moveFile(filename); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		renderOK(w, http.StatusCreated, "upload success")
	} else {
		renderOK(w, http.StatusPartialContent, fmt.Sprintf("chunk %d uploaded success", chunkIndex))
	}
}

// moveFile 将上传完成的临时文件移动到 upload 目录中
func moveFile(filename string) error {
	if err := os.Rename(filepath.Join(uploadTempDir, filename), filepath.Join(uploadRootDir, filename)); err != nil {
		return err
	}
	if err := os.Remove(filepath.Join(uploadTempDir, filename+".info")); err != nil {
		return err
	}
	return nil
}

package main

import (
	"fmt"
	"github.com/pkg/errors"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var chunkSize int64 = 5 * 1024 * 1024 // 5MB

func DownloadFileHandler(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("filename")
	if filename == "" {
		http.Error(w, "File name is required", http.StatusBadRequest)
		return
	}

	downloadFilename := filepath.Join(uploadRootDir, filename)
	file, err := os.Open(downloadFilename)
	if err != nil {
		renderError(w, errors.Wrap(err, "open file failed"))
		return
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		renderError(w, errors.Wrap(err, "get file info failed"))
		return
	}

	// 获取文件大小，如果是分片下载，需要告诉客户端文件大小，否则客户端无法计算进度
	totalSize := fileInfo.Size()

	rangeHeader := r.Header.Get("Range")
	if rangeHeader == "" {
		// 如果请求头中不包含 Range，则说明是普通下载
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileInfo.Name()))
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", totalSize))
		io.Copy(w, file)
		return
	}

	// 如果请求头中包含 Range，则说明是分片下载
	// Range: bytes=0-1024
	rangeParts := strings.Split(rangeHeader, "=")[1]
	rangeStart, _ := strconv.ParseInt(strings.Split(rangeParts, "-")[0], 10, 64)
	rangeEnd := fileInfo.Size() - 1
	if rangeParts != "" {
		rangeEnd, _ = strconv.ParseInt(strings.Split(rangeParts, "-")[1], 10, 64)
		chunkSize = rangeEnd - rangeStart + 1
		// 告诉客户端本次请求的文件大小
		// Content-Range: bytes 0-1024/10240
		w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", rangeStart, rangeEnd, totalSize))
		w.WriteHeader(http.StatusPartialContent)
	}
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileInfo.Name()))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", chunkSize))
	w.Header().Set("Accept-Ranges", "bytes")

	buf := make([]byte, chunkSize)

	// 从 rangeStart 开始读取文件
	_, err = file.Seek(rangeStart, io.SeekStart)
	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			renderError(w, errors.Wrap(err, "read file failed"))
			return
		}
		if n == 0 {
			break
		}
		_, err = w.Write(buf[:n])
		if err != nil {
			log.Panicf("write file failed: %v\n", err)
		}
	}
}

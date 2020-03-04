// Package main ...
// Copyright 2019 The yeqown Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package main

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"docker.io/go-docker"
	"docker.io/go-docker/api/types"
)

type StreamMessage struct {
	Stream string `json:"stream,omitempty"`
}

type AuxMessage struct {
	Aux struct {
		ID string `json:"ID"`
	} `json:"aux"`
}

// https://godoc.org/docker.io/go-docker

func main() {
	cli, err := docker.NewEnvClient()
	if err != nil {
		panic(err)
	}

	dockerFile := "Dockerfile"
	buildCtx, err := tarBuildContext("./testdata")
	if err != nil {
		fmt.Printf("build error: %v", err)
		return
	}

	// https://godoc.org/docker.io/go-docker/api/types#ImageBuildOptions
	opts := types.ImageBuildOptions{
		Tags:           []string{"docker-go-testimage:latest"},
		SuppressOutput: false,
		NoCache:        true,
		CgroupParent:   "cgroup_parent",
		Context:        buildCtx,
		Dockerfile:     dockerFile,
		Remove:         true,
	}

	resp, err := cli.ImageBuild(context.Background(), buildCtx, opts)
	if err != nil {
		fmt.Println("*** error ***")
		panic(err)
	} else {
		fmt.Println("*** success ***")
		fmt.Println(resp)
	}
	defer resp.Body.Close()

	fmt.Println("*** ImageBuild output ***")
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	// var (
	// 	streamMessage StreamMessage
	// 	auxMessage    AuxMessage
	// )

	fmt.Println(buf.String())

	// for _, line := range strings.Split(buf.String(), "\n") {
	// 	if err = json.Unmarshal([]byte(line), &streamMessage); err == nil {
	// 		if streamMessage.Stream != "\n" {
	// 			fmt.Println(strings.TrimSpace(streamMessage.Stream))
	// 		}
	// 		// the AuxMessage struct is being caught above. Why?
	// 	} else if err = json.Unmarshal([]byte(line), &auxMessage); err == nil {
	// 		fmt.Printf(auxMessage.Aux.ID)
	// 	}
	// }
}

func tarBuildContext(dockerPath string) (buildCtx io.Reader, err error) {
	var (
		buf       *bytes.Buffer
		tarWriter *tar.Writer
	)

	buf = bytes.NewBuffer(nil)
	tarWriter = tar.NewWriter(buf)
	defer tarWriter.Close()

	err = filepath.Walk(dockerPath, func(path string, info os.FileInfo, err error) error {
		var (
			fd *os.File
		)

		// open file and read file stat
		if fd, err = os.Open(path); err != nil {
			return err
		}
		defer fd.Close()

		// skip folder
		if info.IsDir() {
			return nil
		}

		// write file header
		if rel, err := filepath.Rel(dockerPath, path); err != nil {
			return err
		} else {
			tarHeader := &tar.Header{
				Name:    rel,
				Size:    info.Size(),
				Mode:    int64(info.Mode()),
				ModTime: info.ModTime(),
			}
			fmt.Printf("Tarring => %s\n", rel)
			if err = tarWriter.WriteHeader(tarHeader); err != nil {
				fmt.Printf("Unable to write tar header: %s", err)
				return err
			}
		}

		// write file content
		if filebody, err := ioutil.ReadAll(fd); err != nil {
			return err
		} else {
			if _, err = tarWriter.Write(filebody); err != nil {
				fmt.Printf("[1m\033[31mUnable to write tar body: %s[0m", err)
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return bytes.NewReader(buf.Bytes()), nil
}

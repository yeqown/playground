package main

// FILE : signal/main.go
// USAGE: 测试 '信号' 和 '命名管道' 的配合
// SITUA: 通过信号触发进程的特定操作，同时向进程传输一些数据

import (
	"bufio"
	"io"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	ch := make(chan os.Signal, 1)
	timeout := time.NewTimer(time.Minute)

	signal.Notify(ch, syscall.SIGUSR1, syscall.SIGUSR2)

	println("pid: " + strconv.Itoa(os.Getpid()))

loop:
	for {
		select {
		case <-ch:
			println("got one signal")
			readFromNamedPipe()
		case <-timeout.C:
			println("timeout, quit")
			break loop
		}
	}

	println("main quit")
}

const (
	// namedPipeFile = "named.pipe"
	namedPipeFile = "named.pipe2"
)

func readFromNamedPipe() {
	println("readFromNamedPipe enter")

	fd, err := os.OpenFile(namedPipeFile, os.O_RDONLY, os.ModeNamedPipe)
	if err != nil {
		println("open pipe faild: " + err.Error())
		return
	}
	defer fd.Close()

	r := bufio.NewReader(fd)
	for {
		line, err := r.ReadBytes('\n')
		if err == io.EOF {

			break
		}

		println("line=" + string(line))
	}

	println("readFromNamedPipe quit")
}

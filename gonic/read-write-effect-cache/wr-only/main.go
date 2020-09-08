package main

import (
	"fmt"
	"io"
	"os"
	"time"
)

func getFilePattern() string {
	return fmt.Sprintf("./logs/rd-swr-%d.log", time.Now().Unix())
}

func prepare() io.Writer {
	fd, err := os.OpenFile(getFilePattern(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	return fd
}

func main() {
	w := prepare()
	line := []byte("this is a line, blabla blabla blabla blabla blabla blabla blabla blabla blabla blabla\n")

	for range time.NewTicker(10 * time.Millisecond).C {
		_, _ = w.Write(line)
	}
}

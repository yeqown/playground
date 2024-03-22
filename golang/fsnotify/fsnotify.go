package main

import (
    "log"
    "os"
    "os/signal"
    "syscall"

    "github.com/fsnotify/fsnotify"
)

func main() {
    watcher, err := fsnotify.NewWatcher()
    if err != nil {
        log.Fatal(err)
    }
    defer watcher.Close()

    done := make(chan bool)
    go func() {
        for {
            select {
            case event, ok := <-watcher.Events:
                if !ok {
                    return
                }
                log.Println("event:", event)
                if event.Op&fsnotify.Write == fsnotify.Write {
                    log.Println("modified file:", event.Name)
                }
            case err, ok := <-watcher.Errors:
                if !ok {
                    return
                }
                log.Println("error:", err)
            }
        }
    }()

    // 文件是 /tmp/testconfig 的软链接
    // ln -s /tmp/testconfig.yaml /tmp/testconfig.softlink
    err = watcher.Add("/tmp/testconfig.softlink")
    if err != nil {
        log.Fatal(err)
    }
    sigs := make(chan os.Signal, 1)
    signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
    <-sigs
    close(done)
}

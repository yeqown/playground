package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	mode = flag.String("mode", "exit", "choose mode to run and test")
	code = flag.Int("code", 0, "os.Exit code")
)

func main() {
	flag.Parse()
	fmt.Printf("mode=%s, ecode=%d\n", *mode, *code)

	// defer fmt.Println("main quit")

	switch *mode {
	case "panic":
		doPanic()
	case "exit":
		doOSExit(*code)
	}
}

//go:noinline
func doPanic() {
	// fmt.Println("panic calling")
	panic("doPanic")
}

//go:noinline
func doOSExit(code int) {
	// fmt.Println("os.Exit calling")
	os.Exit(code)
}

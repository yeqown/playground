package main

import (
	"C"
	"fmt"
)

// 通过export语句来导出指定函数

// SayHello go version of func `SayHello`
//export SayHello
func SayHello(s *C.char) {
	fmt.Println(C.GoString(s))
}

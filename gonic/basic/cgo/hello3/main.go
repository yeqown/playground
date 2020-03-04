package main

// #include<hello.h>
import "C"

// 通过声明一个hello.h来声明一份API，不用关系是用C或者C++还是Go来实现这个头文件。
// 最终调用还是通过“C”包来调用
func main() {
	C.SayHello(C.CString("hello cgo ver3"))
}

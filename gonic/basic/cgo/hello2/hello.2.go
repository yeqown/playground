package main

// void sayHello(const char *s);
import "C"

func main() {
	C.sayHello(C.CString("hello cgo"))
}

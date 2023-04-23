package main

import (
	"fmt"

	"github.com/agiledragon/gomonkey/v2"
)

func main() {
	// use gomonkey to mock add with add2
	patches := gomonkey.ApplyFunc(add, add2)
	_ = patches

	r := add(1, 2)
	fmt.Printf("add(1, 2) = %d\n", r)
}

func add(a, b int) int {
	fmt.Println("add is called")
	return a + b
}

func add2(a, b int) int {
	fmt.Println("add2 is called")
	return a - b
}

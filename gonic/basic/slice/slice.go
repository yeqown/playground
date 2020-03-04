package main

import (
	"fmt"
)

/*
 golang 函数调用使用参数传递（通过stack，也支持了多返回值这一特性）
 对于slice作为参数传递，可以使用s[pos] = value 这种方式修改，但是append操作无法修改(length == capcity)；
 但是如果传递的slice的capcity > length，那么append的值也会修改原slice（这是由于slice是基于数组的）
*/

func main() {
	s := make([]int, 10, 10)
	s2 := make([]int, 10, 20)

	// initialize two slice
	for i := 0; i < 10; i++ {
		s[i] = i + 1
		s2[i] = i + 1
	}

	// change slice (capcity == length)
	// s has no change, append would not change s
	change(s)

	// change slice (capcity > length)
	fmt.Printf("s2 call change before: s2[11] = %v\n", s2[:11]) // [1 2 3 4 5 6 7 8 9 10 0]
	change(s2)
	fmt.Printf("s2 call change after: s2[11] = %v\n", s2[:11]) // [-99 2 3 4 5 6 7 8 9 10 99]

	// copyS2 := make([]int, 11)
	// copy(copyS2, s2[:11])
	// copyS2[11] = 99, append changed copyS2
}

func change(s []int) {
	s[0] = -99
	s = append(s, 99)
}

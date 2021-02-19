package basic_test

import (
	"fmt"
	"testing"
)

/*
 golang 函数调用使用参数传递（通过stack，也支持了多返回值这一特性）
 对于slice作为参数传递，可以使用s[pos] = value 这种方式修改，但是append操作无法修改(length == capcity)；
 但是如果传递的slice的capcity > length，那么append的值也会修改原slice（这是由于slice是基于数组的）
*/

func Test_Slice_transfer(t *testing.T) {
	s := make([]int, 10, 10)
	s2 := make([]int, 10, 20)

	// initialize two slice
	for i := 0; i < 10; i++ {
		s[i] = i + 1
		s2[i] = i + 1
	}

	// change function
	change := func(s []int) {
		s[0] = -99
		s = append(s, 99)
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

func Test_Slice_grownup(t *testing.T) {
	const N = 1024
	var a [N]int

	println(len(a[:N-1:N]), cap(a[:N-1:N])) // [s:e:cap] len=1023 cap=1024
	println(len(a[:N:N]), cap(a[:N:N]))     // len=1024 cap=1024

	x := append(a[:N-1:N], 9, 9)
	y := append(a[:N:N], 9)

	println(cap(x), cap(y))
	// go < 1.16 在扩容时如果 old.len < 1024，则会扩容为 2cap （2048）, 否则扩容为 5/4 cap（1280）
	// cap(x) = 2048, old.len = 1023 <  1024 扩容 2 old.cap   => 2048
	// cap(y) = 1280, old.len = 1024 >= 1024 增加 1/4 old.cap => 1280

	// go1.16
	// 注意 1.16 之后 old.len < 1024 变更为 old.cap < 1024 判断
	// cap(x) = 1280, old.cap = 1024 >= 1024 增加 1/4 old.cap => 1280
	// cap(y) = 1280, old.cap = 1024 >= 1024 增加 1/4 old.cap => 1280
}

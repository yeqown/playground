package main

import (
	"fmt"

	"github.com/smallnest/weighted"
)

func main() {
	var sw *weighted.SW
	// sw := new(weighted.SW)
	data1 := map[int]int{
		1: 10, // left cap 10
		2: 20,
		3: 30,
	}
	sw = load(data1)
	testDividing(sw)

	data2 := map[int]int{
		1: 1,
		2: 11,
		3: 30,
	}
	sw = load(data2)
	testDividing(sw)

	data3 := map[int]int{
		1: -1,
		2: -2,
		3: 20,
	}
	sw = load(data3)
	testDividing(sw)
}

func load(data map[int]int) (sw *weighted.SW) {
	// clear pre data
	sw = new(weighted.SW)

	for k, v := range data {
		sw.Add(k, v)
	}

	return sw
}

func testDividing(sw *weighted.SW) {
	fmt.Println("........testDividing........")
	for i := 0; i < 10; i++ {
		fmt.Printf("sw get: %d\n", sw.Next().(int))
	}
}

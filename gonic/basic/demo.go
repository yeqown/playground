package main

var i = 0

//go:noinline
func f() int {
	i++
	return i
}

func main() {
	// i := 0
	// go:noinline
	// f := func() int {
	//	i++
	//	return i
	// }
	c := make(chan int, 1)
	c <- f()
	select {
	case c <- f():
	default:
		println(i)
	}
}

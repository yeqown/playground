package main

import (
	"time"
	"math"
)

func main() {

	ballast := make([]byte, 10<<30)
	_ = ballast
	<-time.After(time.Duration(math.MaxInt64))
}

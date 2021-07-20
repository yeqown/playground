package basic_test

import (
	"testing"
	"time"
)

func Test_Gorutine_lambda(t *testing.T) {
	var arr = []int{1, 2, 3, 4, 5}
	// var idx_o, v_o = 0, 0 FAILED
	for idx, v := range arr {
		// FAILED：idx_o, v_o = idx, v
		idx_o, v_o := idx, v
		go func() {
			println(idx_o, v_o)
		}()
	}

	time.Sleep(3 * time.Second)
}

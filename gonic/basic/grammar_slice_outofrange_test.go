package basic_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_struct_WhyNotOutOfRange(t *testing.T) {
	var x struct {
		s [][32]byte
	}

	// why not panic, but got 32
	assert.Equal(t, 32, len(x.s[99]))
}

func Test_panicOutOfRange(t *testing.T) {
	str := "wcd"

	// what's the difference between slice[index] and slice[:]:
	// one panic, another works normal.
	print(str[len(str)])
	print(str[len(str):])
}

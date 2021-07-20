package basic_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_struct_WhyNotOutOfRange(t *testing.T) {
	var x struct {
		s [][32]byte
	}

	assert.Equal(t, 99, len(x.s[99]))
}

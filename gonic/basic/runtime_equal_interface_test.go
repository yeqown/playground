package basic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type myErr struct{}

func (e myErr) Error() string {
	return "empty"
}

func getErr() *myErr {
	return nil
}

// Test_equal_interface all pass cases
func Test_equal_interface(t *testing.T) {
	var e error
	e = getErr()
	assert.Equal(t, false, e == nil)

	e2 := getErr()
	assert.Equal(t, true, e2 == nil)

	var e3 *myErr
	e3 = getErr()
	assert.Equal(t, true, e3 == nil)

	var e4 interface{}
	e4 = getErr()
	assert.Equal(t, false, e4 == nil)

	// NOTE: interface{} is dive
}

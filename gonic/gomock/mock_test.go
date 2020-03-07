package gomock_test

import (
	"errors"

	"github.com/golang/mock/gomock"
	mocktest "github.com/playground/gonic/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"testing"
)

func Test_AnimalQuack(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	duck := mocktest.NewMockIAnimal(ctrl)
	duck.
		EXPECT().
		Quack(gomock.Any()).
		AnyTimes().
		DoAndReturn(func(times int) error {
			if times >= 5 {
				return errors.New("too many time to quack")
			}
			return nil
		})

	educk := mocktest.NewMockIAnimal(ctrl)
	educk.
		EXPECT().
		Quack(gomock.Any()). // matcher 用法
		Return(nil).
		AnyTimes()

	assert.True(t, mocktest.AnimalQuack(duck, 5) != nil, nil)
	assert.True(t, mocktest.AnimalQuack(duck, 1) == nil, nil)
	assert.True(t, mocktest.AnimalQuack(educk, 5) == nil, nil)
	assert.True(t, mocktest.AnimalQuack(educk, 1) == nil, nil)
	assert.True(t, mocktest.AnimalQuack(educk, 150) == nil, nil)
}

func Test_testifyMock(t *testing.T) {
	_ = mock.Anything
}

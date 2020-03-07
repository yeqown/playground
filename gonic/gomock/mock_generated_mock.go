// Code generated by MockGen. DO NOT EDIT.
// Source: ./mock.go

// Package gomock is a generated GoMock package.
package gomock

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockIAnimal is a mock of IAnimal interface
type MockIAnimal struct {
	ctrl     *gomock.Controller
	recorder *MockIAnimalMockRecorder
}

// MockIAnimalMockRecorder is the mock recorder for MockIAnimal
type MockIAnimalMockRecorder struct {
	mock *MockIAnimal
}

// NewMockIAnimal creates a new mock instance
func NewMockIAnimal(ctrl *gomock.Controller) *MockIAnimal {
	mock := &MockIAnimal{ctrl: ctrl}
	mock.recorder = &MockIAnimalMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockIAnimal) EXPECT() *MockIAnimalMockRecorder {
	return m.recorder
}

// Quack mocks base method
func (m *MockIAnimal) Quack(times int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Quack", times)
	ret0, _ := ret[0].(error)
	return ret0
}

// Quack indicates an expected call of Quack
func (mr *MockIAnimalMockRecorder) Quack(times interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Quack", reflect.TypeOf((*MockIAnimal)(nil).Quack), times)
}

package tpl

import (
	"container/list"
	"fmt"
)

// stack ....
type stack struct {
	list *list.List
}

// newStack ...
func newStack() *stack {
	list := list.New()
	return &stack{list}
}

// Push method of stack
func (s *stack) Push(value interface{}) {
	s.list.PushBack(value)
}

// Pop method of stack
func (s *stack) Pop() interface{} {
	e := s.list.Back()
	if e != nil {
		s.list.Remove(e)
		return e.Value
	}
	return nil
}

// Peak method of stack
func (s *stack) Peak() interface{} {
	e := s.list.Back()
	if e != nil {
		return e.Value
	}

	return nil
}

// Len method of stack
func (s *stack) Len() int {
	return s.list.Len()
}

// Empty method of stack
func (s *stack) Empty() bool {
	return s.list.Len() == 0
}


func runeHelper(s *stack) (format string) {
	elem := s.list.Front()
	if elem == nil {
		return
	}
	for elem != nil {
		format += fmt.Sprintf("%s,", string(elem.Value.(rune)))
		elem = elem.Next()
	}
	return
}

func float64Helper(s *stack) (format string) {
	elem := s.list.Front()
	if elem == nil {
		return
	}
	for elem != nil {
		format += fmt.Sprintf("%.2f,", elem.Value.(float64))
		elem = elem.Next()
	}
	return
}

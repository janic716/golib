package set

import (
	"fmt"
	"strings"
	"sync"
)

const defaultSetSize = 16

type Set interface {
	Add(val interface{})
	Adds(slice interface{})
	Has(val interface{}) bool
	Clear()
	Remove(val interface{})
}

type SafeSet struct {
	m map[interface{}]bool
	sync.RWMutex
}

func NewSafeSet(eles ...interface{}) *SafeSet {
	len := len(eles) * 2
	if len == 0 {
		len = defaultSetSize
	}
	set := &SafeSet{m: make(map[interface{}]bool, len)}
	set.Adds(eles)
	return set
}

func (s *SafeSet) Add(ele interface{}) {
	s.Lock()
	s.m[ele] = true
	s.Unlock()
}

func (s *SafeSet) Adds(eles ...interface{}) {
	if len(eles) == 0 {
		return
	}
	s.Lock()
	for ele := range eles {
		s.m[ele] = true
	}
	s.Unlock()
}

func (s *SafeSet) Remove(val interface{}) {
	s.Lock()
	delete(s.m, val)
	s.Unlock()
}

func (s *SafeSet) Has(val interface{}) bool {
	s.RLock()
	_, found := s.m[val]
	s.RUnlock()
	return found
}

func (s *SafeSet) Clear() {
	s.Lock()
	s.m = make(map[interface{}]bool)
	s.Unlock()
}

func (s SafeSet) String() string {
	s.RLock()
	var list []string
	for k, _ := range s.m {
		list = append(list, fmt.Sprint(k))
	}
	s.RUnlock()
	return fmt.Sprintf("[%s]", strings.Join(list, " "))
}

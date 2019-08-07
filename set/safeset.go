package set

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

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

func NewHashSet(slice interface{}) *SafeSet {
	set := &SafeSet{m: make(map[interface{}]bool)}
	if slice != nil {
		tp := reflect.TypeOf(slice)
		if tp.Kind() != reflect.Slice {
			return set
		}
		value := reflect.ValueOf(slice)
		l := value.Len()
		for i := 0; i < l; i++ {
			v := value.Index(i)
			set.Add(v.Interface())
		}
	}
	return set
}

func (s *SafeSet) Add(val interface{}) {
	s.Lock()
	s.m[val] = true
	s.Unlock()
}

func (s *SafeSet) AddSlice(slice interface{}) {
	if slice != nil {
		tp := reflect.TypeOf(slice)
		kind := tp.Kind()
		if kind != reflect.Slice && kind != reflect.Array && kind != reflect.String {
			return
		}
		value := reflect.ValueOf(slice)
		l := value.Len()
		for i := 0; i < l; i++ {
			v := value.Index(i)
			s.Add(v.Interface())
		}
	}
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

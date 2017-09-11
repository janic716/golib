package set

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

type Set interface {
	Add(val interface{})
	AddSlice(slice interface{})
	Has(val interface{}) bool
	Clear()
	Remove(val interface{})
}

type HashSet struct {
	m map[interface{}]bool
	sync.RWMutex
}

func NewHashSet(slice interface{}) *HashSet {
	set := &HashSet{m: make(map[interface{}]bool)}
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

func (s *HashSet) Add(val interface{}) {
	s.Lock()
	defer s.Unlock()
	s.m[val] = true
}

func (s *HashSet) AddSlice(slice interface{}) {
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

func (s *HashSet) Remove(val interface{}) {
	s.Lock()
	defer s.Unlock()
	delete(s.m, val)
}

func (s *HashSet) Has(val interface{}) bool {
	s.RLock()
	s.RUnlock()
	_, found := s.m[val]
	return found
}

func (s *HashSet) Clear() {
	s.Lock()
	defer s.Unlock()
	s.m = make(map[interface{}]bool)
}

func (s HashSet) String() string {
	s.RLock()
	s.RUnlock()
	var list []string
	for k, _ := range s.m {
		list = append(list, fmt.Sprint(k))
	}
	return fmt.Sprintf("[%s]", strings.Join(list, " "))
}

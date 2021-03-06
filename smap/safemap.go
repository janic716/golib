package smap

import "sync"

type safeMap struct {
	sync.RWMutex
	data map[interface{}]interface{}
}

/*
// NewBeeMap return new safemap
func NewBeeMap() *BeeMap {
	return &BeeMap{
		lock: new(sync.RWMutex),
		bm:   make(map[interface{}]interface{}),
	}
}

// Get from maps return the k's value
func (m *BeeMap) Get(k interface{}) interface{} {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if val, ok := m.bm[k]; ok {
		return val
	}
	return nil
}

// Set Maps the given key and value. Returns false
// if the key is already in the map and changes nothing.
func (m *BeeMap) Set(k interface{}, v interface{}) bool {
	m.lock.Lock()
	defer m.lock.Unlock()
	if val, ok := m.bm[k]; !ok {
		m.bm[k] = v
	} else if val != v {
		m.bm[k] = v
	} else {
		return false
	}
	return true
}

// Check Returns true if k is exist in the map.
func (m *BeeMap) Check(k interface{}) bool {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if _, ok := m.bm[k]; !ok {
		return false
	}
	return true
}

// Delete the given key and value.
func (m *BeeMap) Delete(k interface{}) {
	m.lock.Lock()
	defer m.lock.Unlock()
	delete(m.bm, k)
}

// Items returns all items in safemap.
func (m *BeeMap) Items() map[interface{}]interface{} {
	m.lock.RLock()
	defer m.lock.RUnlock()
	r := make(map[interface{}]interface{})
	for k, v := range m.bm {
		r[k] = v
	}
	return r
}
 */
func NewSafeMap() *safeMap {
	return &safeMap{data: make(map[interface{}]interface{})}
}

func (s *safeMap) Get(key interface{}) interface{} {
	s.RLock()
	if val, found := s.data[key]; found {
		s.RUnlock()
		return val
	}
	s.RUnlock()
	return nil
}

func (s *safeMap) Set(key, val interface{}) {
	s.Lock()
	s.data[key] = val
	s.Unlock()
}

func (s *safeMap) Del(key interface{}) {
	s.Lock()
	delete(s.data, key)
	s.Unlock()
}

func (s *safeMap)ForeachKV(travelFunc func(k, v interface{}) bool) {
	s.RLock()
	for k, v := range s.data {
		if !travelFunc(k, v) {
			break
		}
	}
	s.RUnlock()
}
package server

import "sync"

type MapRW struct {
	mx sync.RWMutex
	m  map[string]*Hub
}

func (m *MapRW) Load(key string) (*Hub, bool) {
	m.mx.RLock()
	defer m.mx.RUnlock()
	val, ok := m.m[key]

	return val, ok
}

func (m *MapRW) Store(key string, value *Hub) {
	m.mx.Lock()
	defer m.mx.Unlock()
	m.m[key] = value
}

func (m *MapRW) LoadForce(key string) *Hub {
	m.mx.RLock()
	defer m.mx.RUnlock()

	val := m.m[key]

	return val
}

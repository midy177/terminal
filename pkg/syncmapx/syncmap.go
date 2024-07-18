// syncmap implements a generic synchronized map ontop of sync.Map.
package syncmapx

import (
	"golang.org/x/sync/syncmap"
)

// Map is a sync.Map wrapped in a generic manner.
type Map[K any, V any] struct {
	ds *syncmap.Map
}

// New returns a new Map that is generic.
func New[K any, V any]() *Map[K, V] {
	return &Map[K, V]{
		ds: &syncmap.Map{},
	}
}

// Load will load data from the Map.
func (m *Map[K, V]) Load(I K) (V, bool) {
	data, ok := m.ds.Load(I)
	if !ok {
		var res V
		return res, false
	}
	result, ok := data.(V)
	if !ok {
		var res V
		return res, false
	}
	return result, ok
}

// LoadOrStore will load or store data from the Map.
func (m *Map[K, V]) LoadOrStore(I K, J V) (V, bool) {
	data, ok := m.ds.LoadOrStore(I, J)
	if !ok {
		var res V
		return res, false
	}
	result, ok := data.(V)
	if !ok {
		var res V
		return res, false
	}
	return result, ok
}

// LoadAndDelete will load and delete data from the Map.
func (m *Map[K, V]) LoadAndDelete(I K) (V, bool) {
	data, ok := m.ds.LoadAndDelete(I)
	if !ok {
		var res V
		return res, false
	}
	result, ok := data.(V)
	if !ok {
		var res V
		return res, false
	}
	return result, ok
}

// Delete will delete data from the Map.
func (m *Map[K, V]) Delete(I K) {
	m.ds.Delete(I)
}

// Store will store data in the Map.
func (m *Map[K, V]) Store(I K, J V) {
	m.ds.Store(I, J)
}

// Range will iterate over the data in the Map and apply the func on it.
func (m *Map[K, V]) Range(L func(I K, J V) bool) {
	m.ds.Range(func(A, B any) bool {
		C, ok := A.(K)
		if !ok {
			return true
		}
		D, ok := B.(V)
		if !ok {
			return true
		}
		return L(C, D)
	})
}

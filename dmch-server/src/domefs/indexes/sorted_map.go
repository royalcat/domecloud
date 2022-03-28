package indexes

import (
	"sort"
	"sync"

	"golang.org/x/exp/constraints"
)

type SortedMap[T constraints.Ordered, G any] struct {
	mutex sync.RWMutex

	values map[T]G
	order  []T
}

func (m *SortedMap[T, G]) Get(key T) (G, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if v, ok := m.values[key]; ok {
		return v, true
	} else {
		return v, false
	}
}

// TODO переделать это нахуй на linked list
func (m *SortedMap[T, G]) Set(key T, value G) {
	m.mutex.Lock()

	m.values[key] = value
	index := sort.Search(len(m.order), func(i int) bool { return m.order[i] > key })
	m.order = append(m.order, key)
	copy(m.order[index+1:], m.order[index:])
	m.order[index] = key

	m.mutex.Unlock()
}

func (m *SortedMap[T, G]) All() []G {
	m.mutex.RLock()

	out := make([]G, len(m.order))
	for _, k := range m.order {
		out = append(out, m.values[k])
	}

	m.mutex.RUnlock()

	return out
}

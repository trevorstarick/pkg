package db

// src: https://www.reddit.com/r/golang/comments/twucb0/comment/j4x7xbx/?utm_source=share&utm_medium=web2x&context=3

import (
	"cmp"
	"fmt"
	"log/slog"
	"slices"
	"sync"
	"sync/atomic"

	bloom "github.com/bits-and-blooms/bloom/v3"
)

type Map[K comparable, V any] struct {
	m     sync.Map
	c     atomic.Int64
	bloom bloom.BloomFilter
}

func (m *Map[K, V]) Load(key K) (V, bool) {
	if !m.bloom.Test(fmt.Appendf([]byte{}, "%v", key)) {
		return *(new(V)), false
	}

	iface, ok := m.m.Load(key)
	if !ok {
		return *(new(V)), ok
	}

	//nolint:forcetypeassert // this is a Generic type, and we know it's a V
	return iface.(V), ok
}

func (m *Map[K, V]) LoadAndDelete(key K) (V, bool) {
	iface, loaded := m.m.LoadAndDelete(key)
	if loaded {
		m.c.Add(-1)

		//nolint:forcetypeassert // this is a Generic type, and we know it's a V
		return iface.(V), loaded
	}

	return (*new(V)), false
}

func (m *Map[K, V]) LoadOrStore(key K, value V) (V, bool) {
	iface, loaded := m.m.LoadOrStore(key, value)
	if !loaded {
		m.c.Add(1)
		m.bloom.Add(fmt.Appendf([]byte{}, "%v", key))
	}

	//nolint:forcetypeassert // this is a Generic type, and we know it's a V
	return iface.(V), loaded
}

func (m *Map[K, V]) store(key K, value V) {
	m.m.Store(key, value)
	m.bloom.Add(fmt.Appendf([]byte{}, "%v", key))
	m.c.Add(1)
}

func (m *Map[K, V]) Store(key K, value V) {
	m.store(key, value)
}

func (m *Map[K, V]) StoreAndReturn(key K, value V) V {
	m.store(key, value)

	return value
}

func (m *Map[K, V]) delete(key K) {
	m.m.Delete(key)
	m.c.Add(-1)
}

func (m *Map[K, V]) Delete(key K) {
	m.delete(key)
}

func (m *Map[K, V]) Range(f func(K, V) bool) {
	m.m.Range(func(key, value any) bool {
		//nolint:forcetypeassert // this is a Generic type, and we know it's a K,V
		return f(key.(K), value.(V))
	})
}

func (m *Map[K, V]) SortedRange(f func(K, V) bool) {
	keys := make([]K, 0)

	m.m.Range(func(key, _ any) bool {
		if k, ok := key.(K); ok {
			keys = append(keys, k)
		} else {
			slog.Warn("key is not of type K", "key", key)

			return true
		}

		return true
	})

	slices.SortStableFunc(keys, func(a, b K) int {
		if _, ok := any(a).(fmt.Stringer); ok {
			//nolint:forcetypeassert // we do the type assertion in the if statement, so we know it's a fmt.Stringer
			return cmp.Compare(any(a).(fmt.Stringer).String(), any(b).(fmt.Stringer).String())
		}

		switch any(a).(type) {
		case string:
			//nolint:forcetypeassert // we do the type assertion in the switch statement, so we know it's a string
			return cmp.Compare(any(a).(string), any(b).(string))
		case int:
			//nolint:forcetypeassert // we do the type assertion in the switch statement, so we know it's an int
			return any(a).(int) - any(b).(int)
		default:
			panic("unsupported key type")
		}
	})

	for _, key := range keys {
		value, _ := m.m.Load(key)
		//nolint:forcetypeassert // this is a Generic type, and we know it's a V
		if !f(key, value.(V)) {
			break
		}
	}
}

func (m *Map[K, V]) Contains(key K) bool {
	if !m.bloom.Test(fmt.Appendf([]byte{}, "%v", key)) {
		return false
	}

	_, ok := m.m.Load(key)

	return ok
}

func (m *Map[K, V]) Len() int {
	return int(m.c.Load())
}

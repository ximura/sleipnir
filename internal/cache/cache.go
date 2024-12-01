package cache

import "sync"

// Cache - generic memory store
// K - key is any comparable type
// V - value is any type
type Cache[K comparable, V any] struct {
	m     sync.Mutex
	store map[K]V
}

func New[K comparable, V any]() Cache[K, V] {
	return Cache[K, V]{
		store: make(map[K]V),
	}
}

func (c *Cache[K, V]) Set(key K, value V) {
	c.m.Lock()
	c.store[key] = value
	c.m.Unlock()
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.m.Lock()
	v, ok := c.store[key]
	c.m.Unlock()
	return v, ok
}

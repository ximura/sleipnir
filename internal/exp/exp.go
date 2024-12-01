package exp

import "sync"

// SimpleMutexCache - generic memory store, only supports write and reads
// use sync.Mutex for thread safety
// K - key is any comparable type
// V - value is any type
type SimpleMutexCache[K comparable, V any] struct {
	m     sync.Mutex
	store map[K]V
}

func NewSimpleMutexCache[K comparable, V any](size int) SimpleMutexCache[K, V] {
	return SimpleMutexCache[K, V]{
		store: make(map[K]V, size),
	}
}

func (c *SimpleMutexCache[K, V]) Set(key K, value V) {
	c.m.Lock()
	c.store[key] = value
	c.m.Unlock()
}

func (c *SimpleMutexCache[K, V]) Get(key K) (V, bool) {
	c.m.Lock()
	v, ok := c.store[key]
	c.m.Unlock()
	return v, ok
}

// SimpleRWMutexCache - generic memory store, only supports write and reads
// use sync.RWMutex for thread safety
// K - key is any comparable type
// V - value is any type
type SimpleRWMutexCache[K comparable, V any] struct {
	m     sync.RWMutex
	store map[K]V
}

func NewSimpleRWMutexCache[K comparable, V any](size int) SimpleRWMutexCache[K, V] {
	return SimpleRWMutexCache[K, V]{
		store: make(map[K]V, size),
	}
}

func (c *SimpleRWMutexCache[K, V]) Set(key K, value V) {
	c.m.Lock()
	c.store[key] = value
	c.m.Unlock()
}

func (c *SimpleRWMutexCache[K, V]) Get(key K) (V, bool) {
	c.m.RLock()
	v, ok := c.store[key]
	c.m.RUnlock()
	return v, ok
}

// SimpleRWMutexCache - generic memory store, only supports write and reads
// use sync.RWMutex for thread safety
type SimpleMapCache struct {
	store sync.Map
}

func NewSimpleMapCache() SimpleMapCache {
	return SimpleMapCache{
		store: sync.Map{},
	}
}

func (c *SimpleMapCache) Set(key, value any) {
	c.store.Store(key, value)
}

func (c *SimpleMapCache) Get(key any) (any, bool) {
	return c.store.Load(key)
}

package cache

// Cache - generic memory store
// K - key is any comparable type
// V - value is any type
type Cache[K comparable, V any] struct {
	store map[K]V
}

func New[K comparable, V any]() Cache[K, V] {
	return Cache[K, V]{
		store: make(map[K]V),
	}
}

func (c *Cache[K, V]) Set(key K, value V) {
	c.store[key] = value
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	v, ok := c.store[key]
	return v, ok
}

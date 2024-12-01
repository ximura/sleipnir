package cache

import "sync"

// Cache - generic memory store, only supports write and reads
// use sync.RWMutex for thread safety
// K - key is any comparable type
// V - value is any type
type Cache[K comparable, V any] struct {
	mu    sync.RWMutex
	store map[K]V
}

// Load returns the value stored in the cache for a key
// The ok result indicates whether value was found in the cache.
func (c *Cache[K, V]) Load(key K) (value V, ok bool) {
	return value, ok
}

// Store sets the value for a key.
func (c *Cache[K, V]) Store(key K, value V) {
	c.Swap(key, value)
}

// Clear deletes all the entries, resulting in an empty cache.
func (c *Cache[K, V]) Clear() {
}

// LoadOrStore returns the existing value for the key if present.
// Otherwise, it stores and returns the given value.
// The loaded result is true if the value was loaded, false if stored.
func (c *Cache[K, V]) LoadOrStore(key, value any) (actual V, loaded bool) {
	return actual, loaded
}

// LoadAndDelete deletes the value for a key, returning the previous value if any.
// The loaded result reports whether the key was present.
func (c *Cache[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	return value, loaded
}

// Delete deletes the value for a key.
func (c *Cache[K, V]) Delete(key K) {
	c.LoadAndDelete(key)
}

// Swap swaps the value for a key and returns the previous value if any.
// The loaded result reports whether the key was present.
func (c *Cache[K, V]) Swap(key K, value V) (previous V, loaded bool) {
	return value, loaded
}

// CompareAndSwap swaps the old and new values for key
// if the value stored in the map is equal to old.
// The old value must be of a comparable type.
func (c *Cache[K, V]) CompareAndSwap(key K, old, new V) (swapped bool) {
	return swapped
}

// CompareAndDelete deletes the entry for key if its value is equal to old.
// The old value must be of a comparable type.
//
// If there is no current value for key in the map, CompareAndDelete
// returns false (even if the old value is the nil interface value).
func (c *Cache[K, V]) CompareAndDelete(key K, old V) (deleted bool) {
	return deleted
}

// Range calls f sequentially for each key and value present in the map.
// If f returns false, range stops the iteration.
func (c *Cache[K, V]) Range(f func(key K, value V) bool) {
}

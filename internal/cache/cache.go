package cache

import (
	"sync"
)

func New[K comparable, V any]() Cache[K, V] {
	return Cache[K, V]{
		store: make(map[K]V),
	}
}

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
	c.mu.RLock()
	value, ok = c.store[key]
	c.mu.RUnlock()
	return value, ok
}

// Store sets the value for a key.
func (c *Cache[K, V]) Store(key K, value V) {
	c.mu.Lock()
	c.store[key] = value
	c.mu.Unlock()
}

// Clear deletes all the entries, resulting in an empty cache.
func (c *Cache[K, V]) Clear() {
	c.mu.Lock()
	// be aware that clear doesn't shrink map size it only deletes content
	clear(c.store)
	c.mu.Unlock()
}

// LoadOrStore returns the existing value for the key if present.
// Otherwise, it stores and returns the given value.
// The loaded result is true if the value was loaded, false if stored.
func (c *Cache[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	// try to check in read only mode before failing into write lock
	{
		c.mu.RLock()
		value, ok := c.store[key]
		if ok {
			c.mu.RUnlock()
			return value, true
		}
		c.mu.RUnlock()
	}

	// now check in write mode, and store value if not present
	{
		c.mu.Lock()
		defer c.mu.Unlock() // we would be in critical section until end of this function
		// we need to check again to prevent race conditions between read and write locks
		value, ok := c.store[key]
		if ok {
			return value, true
		}

		c.store[key] = value
	}
	return value, false
}

// LoadAndDelete deletes the value for a key, returning the previous value if any.
// The loaded result reports whether the key was present.
func (c *Cache[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	// try to check in read only mode before failing into write lock
	{
		c.mu.RLock()
		value, ok := c.store[key]
		if !ok {
			c.mu.RUnlock()
			return value, false
		}
		c.mu.RUnlock()
	}

	// now check in write mode, and delete value if stil present
	{
		c.mu.Lock()
		defer c.mu.Unlock() // we would be in critical section until end of this function
		value, ok := c.store[key]
		if !ok {
			// value not present nothing to delete
			return value, false
		}

		delete(c.store, key)
	}
	return value, true
}

// Delete deletes the value for a key.
func (c *Cache[K, V]) Delete(key K) {
	c.mu.Lock()
	delete(c.store, key)
	c.mu.Unlock()
}

// Swap swaps the value for a key and returns the previous value if any.
// The loaded result reports whether the key was present.
func (c *Cache[K, V]) Swap(key K, value V) (previous V, loaded bool) {
	c.mu.Lock()
	defer c.mu.Unlock() // we would be in critical section until end of this function
	previous, loaded = c.store[key]
	c.store[key] = value

	return previous, loaded
}

// CompareAndSwap swaps the old and new values for key
// if the value stored in the map is equal to old.
// The old value must be of a comparable type.
func (c *Cache[K, V]) CompareAndSwap(key K, old, new V, cmp func(old, value V) bool) (swapped bool) {
	// first check in read only mode
	{
		c.mu.RLock()
		value, ok := c.store[key]
		if !ok || !cmp(value, old) {
			c.mu.RUnlock()
			return false
		}
		c.mu.RUnlock()
	}

	// read check found value and it's comparable, we need to enter write mode and repeat checks
	{
		c.mu.Lock()
		defer c.mu.Unlock()
		value, ok := c.store[key]
		if !ok || !cmp(value, old) {
			// value not present nothing to swap
			// or cmp returned false
			return false
		}

		c.store[key] = new
	}

	return true
}

// CompareAndDelete deletes the entry for key if its value is equal to old.
// The old value must be of a comparable type.
//
// If there is no current value for key in the map, CompareAndDelete
// returns false.
func (c *Cache[K, V]) CompareAndDelete(key K, old V, cmp func(old, value V) bool) (deleted bool) {
	// first check in read only mode
	{
		c.mu.RLock()
		value, ok := c.store[key]
		if !ok || !cmp(value, old) {
			c.mu.RUnlock()
			return false
		}
		c.mu.RUnlock()
	}

	// read check found value and it's comparable, we need to enter write mode and repeat checks
	{
		c.mu.Lock()
		defer c.mu.Unlock()
		value, ok := c.store[key]
		if !ok || !cmp(value, old) {
			// value not present nothing to swap
			// or cmp returned false
			return false
		}

		delete(c.store, key)
	}
	return true
}

// Range calls f sequentially for each key and value present in the map.
// If f returns false, range stops the iteration.
func (c *Cache[K, V]) Range(f func(key K, value V) bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	for k, v := range c.store {
		if !f(k, v) {
			return
		}
	}
}

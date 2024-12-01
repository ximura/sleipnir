package cache_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/ximura/sleipnir/internal/cache"
	"gotest.tools/v3/assert"
)

func TestEmptySet(t *testing.T) {
	c := cache.New[string, int]()
	v, ok := c.Load("foo")
	assert.Equal(t, ok, false, "Got true for non existing item")
	assert.Equal(t, v, 0, "Got wrong value for non existing item")
}

func TestGetAfterSet(t *testing.T) {
	c := cache.New[string, int]()
	c.Store("foo", 1)
	v, ok := c.Load("foo")
	assert.Equal(t, ok, true, "Got false for existing item")
	assert.Equal(t, v, 1, "Got wrong value for existing item")
}

func TestShouldOverwriteExistingValue(t *testing.T) {
	c := cache.New[string, int]()
	c.Store("foo", 1)
	v, ok := c.Load("foo")
	assert.Equal(t, ok, true, "Got false for existing item")
	assert.Equal(t, v, 1, "Got wrong value for existing item")

	c.Store("foo", 2)
	v, ok = c.Load("foo")
	assert.Equal(t, ok, true, "Got false for existing item")
	assert.Equal(t, v, 2, "Got wrong value for existing item")
}

func TestParallelSet(t *testing.T) {
	n := 30
	c := cache.New[string, int]()

	var wg sync.WaitGroup
	wg.Add(n)
	for i := range n {
		go func(i int) {
			c.Store("foo", i)
			c.Store("bar", i)
			wg.Done()
		}(i)
	}

	wg.Wait()

	_, ok1 := c.Load("foo")
	_, ok2 := c.Load("bar")

	assert.Equal(t, ok1, true, "Got false for existing item")
	assert.Equal(t, ok2, true, "Got false for existing item")
}

// TestConcurrentClear tests concurrent behavior of Cache properties to ensure no data races.
// Checks for proper synchronization between Clear, Store, Load operations.
func TestConcurrentClear(t *testing.T) {
	c := cache.New[int, int]()

	wg := sync.WaitGroup{}
	wg.Add(30) // 10 goroutines for writing, 10 goroutines for reading, 10 goroutines for waiting

	// Writing data to the map concurrently
	for i := 0; i < 10; i++ {
		go func(k, v int) {
			defer wg.Done()
			c.Store(k, v)
		}(i, i*10)
	}

	// Reading data from the map concurrently
	for i := 0; i < 10; i++ {
		go func(k int) {
			defer wg.Done()
			if value, ok := c.Load(k); ok {
				t.Logf("Key: %v, Value: %v\n", k, value)
			} else {
				t.Logf("Key: %v not found\n", k)
			}
		}(i)
	}

	// Clearing data from the map concurrently
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			c.Clear()
		}()
	}

	wg.Wait()

	c.Clear()

	var err error
	c.Range(func(k, v int) bool {
		t.Errorf("after Clear, cache contains (%v, %v); expected to be empty", k, v)
		err = fmt.Errorf("cache not empty")

		return true
	})

	assert.NilError(t, err)
}

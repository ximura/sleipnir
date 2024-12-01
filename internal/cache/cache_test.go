package cache_test

import (
	"sync"
	"testing"

	"github.com/ximura/sleipnir/internal/cache"
	"gotest.tools/v3/assert"
)

func TestEmptySet(t *testing.T) {
	c := cache.New[string, int]()
	v, ok := c.Get("foo")
	assert.Equal(t, ok, false, "Got true for non existing item")
	assert.Equal(t, v, 0, "Got wrong value for non existing item")
}

func TestGetAfterSet(t *testing.T) {
	c := cache.New[string, int]()
	c.Set("foo", 1)
	v, ok := c.Get("foo")
	assert.Equal(t, ok, true, "Got false for existing item")
	assert.Equal(t, v, 1, "Got wrong value for existing item")
}

func TestShouldOverwriteExistingValue(t *testing.T) {
	c := cache.New[string, int]()
	c.Set("foo", 1)
	v, ok := c.Get("foo")
	assert.Equal(t, ok, true, "Got false for existing item")
	assert.Equal(t, v, 1, "Got wrong value for existing item")

	c.Set("foo", 2)
	v, ok = c.Get("foo")
	assert.Equal(t, ok, true, "Got false for existing item")
	assert.Equal(t, v, 2, "Got wrong value for existing item")
}

func TestParallelSet(t *testing.T) {
	n := 10
	c := cache.New[string, int]()

	var wg sync.WaitGroup
	wg.Add(n)
	for i := range n {
		go func(i int) {
			c.Set("foo", i)
			c.Set("bar", i)
			wg.Done()
		}(i)
	}

	wg.Wait()

	_, ok1 := c.Get("foo")
	_, ok2 := c.Get("bar")

	assert.Equal(t, ok1, true, "Got false for existing item")
	assert.Equal(t, ok2, true, "Got false for existing item")
}

package cache_test

import (
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

package cache

import "testing"

func TestRedisCache_Set(t *testing.T) {
	cache := NewCache()
	_ = cache.Set("hello", "world")
}

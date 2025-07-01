package cache

import (
	"fmt"
	"time"
)

var DefaultTtl = time.Minute

type Cache interface {
	Get(key string) (any, bool)
	SetTtl(key string, value any, ttl time.Duration)
	Del(key string)
}

// Do execute the function f and cache the result with the provided key and ttl, if the key exists in the cache, return the cached value
func Do[T any](cache Cache, key string, ttl time.Duration, f func() (T, error)) (T, error) {
	if cachedValue, ok := cache.Get(key); ok {
		if valueCast, ok := cachedValue.(T); ok {
			return valueCast, nil
		} else {
			cache.Del(key)
		}
	}
	value, err := f()
	if err != nil {
		return value, fmt.Errorf("DoWithCache: %w", err)
	}
	cache.SetTtl(key, value, ttl)
	return value, nil
}

// Get gets cached value by provided key from the provided cache,
// the bool return value indicates whether the key exists in the cache and the value is of the correct type
func Get[T any](cache Cache, key string) (T, bool) {
	var defaultValue T
	value, ok := cache.Get(key)
	if !ok {
		return defaultValue, false
	}
	t, ok := value.(T)
	if !ok {
		return defaultValue, false
	}
	return t, true
}

// SetTtl sets a key value pair to the provided cache with a ttl
func SetTtl(cache Cache, key string, value any, ttl time.Duration) {
	cache.SetTtl(key, value, ttl)
}

// Del deletes a key value pair from the provided cache
func Del(cache Cache, key string) {
	cache.Del(key)
}

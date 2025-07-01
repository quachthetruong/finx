package mock

import (
	"time"
)

// EmptyCache is an empty implementation of cache.Cache
type EmptyCache struct{}

func (m *EmptyCache) Get(_ string) (any, bool) {
	return nil, false
}

func (m *EmptyCache) SetTtl(_ string, _ any, _ time.Duration) {
}

func (m *EmptyCache) Del(_ string) {
}

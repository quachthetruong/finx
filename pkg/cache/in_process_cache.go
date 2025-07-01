package cache

import (
	"context"
	"time"

	"github.com/dgraph-io/ristretto"

	"financing-offer/pkg/shutdown"
)

type InProcessCache struct {
	store *ristretto.Cache
}

// Get gets cached value by provided key, the bool return value indicates whether the key exists in the cache
func (c *InProcessCache) Get(key string) (any, bool) {
	return c.store.Get(key)
}

// SetTtl sets a key value pair with a ttl
func (c *InProcessCache) SetTtl(key string, value any, ttl time.Duration) {
	c.store.SetWithTTL(key, value, 1, ttl)
}

// Del deletes a key value pair
func (c *InProcessCache) Del(key string) {
	c.store.Del(key)
}

func NewInProcessCache(task *shutdown.Tasks) (*InProcessCache, error) {
	r, err := ristretto.NewCache(
		&ristretto.Config{
			MaxCost:     1 << 30, // 1GB
			NumCounters: 1e7,     // 10M
			BufferItems: 64,
		},
	)
	if err != nil {
		return nil, err
	}
	task.AddShutdownTask(
		func(_ context.Context) error {
			r.Close()
			return nil
		},
	)
	return &InProcessCache{
		store: r,
	}, nil
}

package cache

import (
	"fmt"
	"time"

	"semrush/cache/internal/lfu"
	"semrush/cache/internal/lru"
	"semrush/cache/internal/timebased"
)

// CacheStrategy represents the cache eviction strategy.
type CacheStrategy int

const (
	// LRU strategy (Least Recently Used)
	LRU CacheStrategy = iota

	// LFU strategy (Least Frequently Used)
	LFU

	// TimeBased expiration strategy
	TimeBased
)

// Cache interface represents the methods a cache should implement.
type Cache interface {
	Set(key string, value any)
	Get(key string) (any, bool)
	Delete(key string)
	Clear()
	Len() int
}

// New creates a new cache instance based on the given strategy and capacity.
func New(strategy CacheStrategy, capacity int, expirationDuration time.Duration) (Cache, error) {
	switch strategy {
	case LRU:
		return lru.NewLRUCache(capacity), nil
	case LFU:
		return lfu.NewLFUCache(capacity), nil
	case TimeBased:
		if expirationDuration <= 0 {
			return nil, fmt.Errorf("expiration duration must be positive for TimeBased strategy")
		}
		return timebased.NewTimeBasedCache(capacity, expirationDuration), nil
	default:
		return nil, fmt.Errorf("unknown cache strategy: %d", strategy)
	}
}

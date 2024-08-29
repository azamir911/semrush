package timebased

import (
	"semrush/cache/internal/base"

	"sync"
	"time"
)

// item stores the value and its expiration time
type item struct {
	value      interface{}
	expiration time.Time
}

// TimeBasedCache represents a cache with time-based expiration of items.
type TimeBasedCache struct {
	*base.BaseCache
	expirationDuration time.Duration
	ttlMap             map[string]item
	stopChan           chan struct{}
	stopOnce           sync.Once // Ensures StopCleanup is called only once
}

// NewTimeBasedCache initializes a new TimeBasedCache with the given capacity and expiration duration.
func NewTimeBasedCache(capacity int, expirationDuration time.Duration) *TimeBasedCache {
	cache := &TimeBasedCache{
		BaseCache:          base.NewBaseCache(capacity),
		expirationDuration: expirationDuration,
		ttlMap:             make(map[string]item),
		stopChan:           make(chan struct{}),
	}

	go cache.cleanupExpiredItems()
	return cache
}

// Set adds a new item to the cache with an expiration time.
func (c *TimeBasedCache) Set(key string, value interface{}) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	if len(c.Data) >= c.Capacity {
		c.evict()
	}

	expiration := time.Now().Add(c.expirationDuration)
	c.Data[key] = value
	c.ttlMap[key] = item{value, expiration}
}

// Get retrieves an item from the cache if it hasn't expired.
func (c *TimeBasedCache) Get(key string) (interface{}, bool) {
	c.Mutex.RLock()
	defer c.Mutex.RUnlock()

	if item, ok := c.ttlMap[key]; ok {
		if time.Now().Before(item.expiration) {
			return item.value, true
		}
		// Item has expired
		c.Mutex.RUnlock()
		c.Mutex.Lock()
		defer c.Mutex.Unlock()
		delete(c.Data, key)
		delete(c.ttlMap, key)
		return nil, false
	}
	return nil, false
}

// evict removes the oldest item (based on expiration time) from the cache.
func (c *TimeBasedCache) evict() {
	var oldestKey string
	var oldestExpiration time.Time = time.Now().Add(c.expirationDuration)

	for key, item := range c.ttlMap {
		if item.expiration.Before(oldestExpiration) {
			oldestKey = key
			oldestExpiration = item.expiration
		}
	}

	if oldestKey != "" {
		delete(c.Data, oldestKey)
		delete(c.ttlMap, oldestKey)
	}
}

// cleanupExpiredItems periodically removes expired items from the cache.
func (c *TimeBasedCache) cleanupExpiredItems() {
	ticker := time.NewTicker(c.expirationDuration / 2)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.Mutex.Lock()
			now := time.Now()
			for key, item := range c.ttlMap {
				if now.After(item.expiration) {
					delete(c.Data, key)
					delete(c.ttlMap, key)
				}
			}
			c.Mutex.Unlock()
		case <-c.stopChan:
			return
		}
	}
}

// StopCleanup stops the cleanup goroutine when the cache is no longer in use.
func (c *TimeBasedCache) StopCleanup() {
	c.stopOnce.Do(func() {
		close(c.stopChan)
	})
}

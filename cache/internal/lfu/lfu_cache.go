package lfu

import (
	"math"

	"semrush/cache/internal/base"
)

// LFUCache represents a cache that evicts the least frequently used items.
type LFUCache struct {
	*base.BaseCache
	freq map[string]int
}

// NewLFUCache initializes a new LFUCache with the given capacity.
func NewLFUCache(capacity int) *LFUCache {
	return &LFUCache{
		BaseCache: base.NewBaseCache(capacity),
		freq:      make(map[string]int),
	}
}

// Set adds a new item to the cache, updating its frequency or evicting the least frequently used item if necessary.
func (c *LFUCache) Set(key string, value any) {
	c.BaseCache.Mutex.Lock()
	defer c.BaseCache.Mutex.Unlock()

	if _, exists := c.BaseCache.Data[key]; exists {
		// Update the existing item without increasing capacity
		c.freq[key]++
		c.BaseCache.Data[key] = value
		return
	}

	if len(c.BaseCache.Data) >= c.BaseCache.Capacity {
		c.evict()
	}

	c.BaseCache.Data[key] = value
	c.freq[key] = 1
}

// Get retrieves an item from the cache and increments its frequency.
func (c *LFUCache) Get(key string) (any, bool) {
	c.BaseCache.Mutex.RLock()
	defer c.BaseCache.Mutex.RUnlock()

	value, ok := c.BaseCache.Get(key)
	if ok {
		c.freq[key]++
	}
	return value, ok
}

// evict removes the least frequently used item from the cache.
func (c *LFUCache) evict() {
	var minFreq int = math.MaxInt
	var evictKey string

	for key, freq := range c.freq {
		if freq < minFreq {
			minFreq = freq
			evictKey = key
		}
	}

	if evictKey != "" {
		delete(c.BaseCache.Data, evictKey)
		delete(c.freq, evictKey)
	}
}

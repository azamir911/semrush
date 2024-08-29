package lru

import (
	"container/list"

	"semrush/cache/internal/base"
)

// LRUCache represents a cache that evicts the least recently used items.
type LRUCache struct {
	*base.BaseCache
	order *list.List
}

// NewLRUCache initializes a new LRUCache with the given capacity.
func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		BaseCache: base.NewBaseCache(capacity),
		order:     list.New(),
	}
}

func (c *LRUCache) Set(key string, value any) {
	c.BaseCache.Mutex.Lock()
	defer c.BaseCache.Mutex.Unlock()

	if _, exists := c.BaseCache.Data[key]; exists {
		// Update the existing item, but don't increase capacity
		c.updateOrder(key)
		c.BaseCache.Data[key] = value
		return
	}

	if len(c.BaseCache.Data) >= c.BaseCache.Capacity {
		c.evict()
	}
	c.BaseCache.Data[key] = value
	c.order.PushFront(key)
}

func (c *LRUCache) Get(key string) (any, bool) {
	c.BaseCache.Mutex.RLock()
	defer c.BaseCache.Mutex.RUnlock()

	value, ok := c.BaseCache.Get(key)
	if ok {
		c.updateOrder(key)
	}
	return value, ok
}

func (c *LRUCache) evict() {
	oldest := c.order.Back()
	if oldest != nil {
		c.order.Remove(oldest)
		delete(c.BaseCache.Data, oldest.Value.(string))
	}
}

func (c *LRUCache) updateOrder(key string) {
	for e := c.order.Front(); e != nil; e = e.Next() {
		if e.Value == key {
			c.order.Remove(e)
			break
		}
	}
	c.order.PushFront(key)
}

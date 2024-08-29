package base

import "sync"

// BaseCache provides common functionality for all cache implementations.
type BaseCache struct {
	Data     map[string]any
	Capacity int
	Mutex    sync.RWMutex
}

// NewBaseCache initializes a new BaseCache with a given capacity.
func NewBaseCache(capacity int) *BaseCache {
	return &BaseCache{
		Data:     make(map[string]any),
		Capacity: capacity,
	}
}

func (c *BaseCache) Set(key string, value any) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	// Eviction is no longer handled here, to be managed by specific cache types.
	c.Data[key] = value
}

func (c *BaseCache) Get(key string) (any, bool) {
	c.Mutex.RLock()
	defer c.Mutex.RUnlock()
	value, ok := c.Data[key]
	return value, ok
}

func (c *BaseCache) Delete(key string) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	delete(c.Data, key)
}

func (c *BaseCache) Clear() {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	c.Data = make(map[string]any)
}

func (c *BaseCache) Len() int {
	c.Mutex.RLock()
	defer c.Mutex.RUnlock()
	return len(c.Data)
}

package cache

import (
	"sync"
	"time"
)

type BasicCache struct {
	mu    sync.RWMutex
	items map[string]entry
}

func NewBasicCache() *BasicCache {
	return &BasicCache{
		items: make(map[string]entry),
	}
}

func (c *BasicCache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	e, ok := c.items[key]
	if !ok || time.Now().After(e.expiresAt) {
		return nil, false
	}
	return e.value, true
}

func (c *BasicCache) Set(key string, value []byte, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = entry{
		value:     value,
		expiresAt: time.Now().Add(ttl),
	}
}

func (c *BasicCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

func (c *BasicCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	clear(c.items)
}

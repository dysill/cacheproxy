package cache

import (
	"sync"
	"time"
)

type FIFOCache struct {
	mu       sync.Mutex
	items    map[string]*node
	capacity int
	list     *dll
}

func NewFIFOCache(capacity int) *FIFOCache {
	return &FIFOCache{
		capacity: capacity,
		items:    make(map[string]*node),
		list:     newDLL(),
	}
}

func (c *FIFOCache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	n, ok := c.items[key]
	if !ok {
		return nil, false
	}
	if time.Now().After(n.entry.expiresAt) {
		c.list.remove(n)
		delete(c.items, key)
		return nil, false
	}
	return n.entry.value, true
}

func (c *FIFOCache) Set(key string, value []byte, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	e := entry{
		value:     value,
		expiresAt: time.Now().Add(ttl),
	}

	if existing, ok := c.items[key]; ok {
		c.list.moveToFront(existing)
		existing.entry = e
		return
	}

	if len(c.items) >= c.capacity {
		oldest := c.list.removeBack()
		delete(c.items, oldest.key)
	}

	n := &node{
		key:   key,
		entry: e,
	}
	c.items[key] = n
	c.list.insertFront(n)
}

func (c *FIFOCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	n, ok := c.items[key]
	if !ok {
		return
	}
	delete(c.items, key)
	c.list.remove(n)
}

func (c *FIFOCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	clear(c.items)
	c.list = newDLL()
}

package cache

import (
	"sync"
	"time"
)

type LRUCache struct {
	mu       sync.Mutex
	items    map[string]*node
	capacity int
	list     *dll
}

func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		items:    make(map[string]*node),
		list:     newDLL(),
	}
}

func (c *LRUCache) Get(key string) ([]byte, bool) {
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

	c.list.moveToFront(n) // LRU
	return n.entry.value, true
}

func (c *LRUCache) Set(key string, value []byte, ttl time.Duration) {
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
		lru := c.list.removeBack()
		delete(c.items, lru.key)
	}

	n := &node{
		key:   key,
		entry: e,
	}
	c.items[key] = n
	c.list.insertFront(n)
}

func (c *LRUCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	n, ok := c.items[key]
	if !ok {
		return
	}
	delete(c.items, key)
	c.list.remove(n)
}

func (c *LRUCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	clear(c.items)
	c.list = newDLL()
}

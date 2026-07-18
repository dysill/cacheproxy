package cache

import (
	"fmt"
	"time"
)

type Cache interface {
	Get(key string) ([]byte, bool)
	Set(key string, value []byte, ttl time.Duration)
	Delete(key string)
	Clear()
}

func NewCache(policy string, capacity int) (Cache, error) {
	switch policy {
	case "basic":
		return NewBasicCache(), nil
	case "fifo":
		return NewFIFOCache(capacity), nil
	case "lru":
		return NewLRUCache(capacity), nil
	default:
		return nil, fmt.Errorf("invalid cache policy: %s", policy)
	}
}

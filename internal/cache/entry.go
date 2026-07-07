package cache

import "time"

type entry struct {
	value     []byte
	expiresAt time.Time
}

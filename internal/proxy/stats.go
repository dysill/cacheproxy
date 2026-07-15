package proxy

import (
	"encoding/json"
	"net/http"
	"sync/atomic"
)

type stats struct {
	hits   atomic.Int64
	misses atomic.Int64
}

func (h *ProxyHandler) StatsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int64{
		"hits":   h.stats.hits.Load(),
		"misses": h.stats.misses.Load(),
	})
}

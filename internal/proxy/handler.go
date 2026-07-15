package proxy

import (
	"io"
	"net/http"
	"time"

	"github.com/dysill/cacheproxy/internal/cache"
)

type ProxyHandler struct {
	cache       cache.Cache
	upstreamURL string
	client      *http.Client
	defaultTTL  time.Duration
	stats       stats
}

func NewProxyHandler(cache cache.Cache, upstreamURL string, ttl time.Duration) *ProxyHandler {
	return &ProxyHandler{
		cache:       cache,
		upstreamURL: upstreamURL,
		client:      &http.Client{},
		defaultTTL:  ttl,
	}
}

func (h *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path
	value, found := h.cache.Get(key)
	if found { // cache hit
		h.stats.hits.Add(1)
		w.Header().Set("Content-Type", "application/json") // Only handle JSON for now, could generalize later
		w.Write(value)
		return
	}
	h.stats.misses.Add(1)

	resp, err := h.client.Get(h.upstreamURL + key)
	if err != nil {
		http.Error(w, "failed to reach upstream", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "upstream returned an error", resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "failed to read upstream response", http.StatusInternalServerError)
		return
	}

	h.cache.Set(key, body, h.defaultTTL)
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

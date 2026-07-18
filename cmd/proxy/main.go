package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dysill/cacheproxy/internal/cache"
	"github.com/dysill/cacheproxy/internal/proxy"
)

func main() {
	port := flag.Int("port", 8081, "Port number for proxy")
	upstreamURL := flag.String("upstream", "http://localhost:8080", "Base URL of the upstream server")
	ttlSeconds := flag.Int("ttl", 60, "Default Cache TTL in seconds")
	policy := flag.String("policy", "lru", "Cache eviciton policy")
	capacity := flag.Int("capacity", 100, "Cache capacity")
	flag.Parse()
	defaultTTL := time.Duration(*ttlSeconds) * time.Second

	c, err := cache.NewCache(*policy, *capacity)
	if err != nil {
		log.Fatal(err)
	}
	h := proxy.NewProxyHandler(c, *upstreamURL, defaultTTL)
	mux := http.NewServeMux()
	mux.Handle("/", h)
	mux.HandleFunc("/stats", h.StatsHandler)

	addr := fmt.Sprintf(":%d", *port)
	log.Printf("proxy server listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

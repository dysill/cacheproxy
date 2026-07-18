# cacheproxy

An HTTP caching reverse proxy written in Go with choice of different cache evicition policies (LRU, FIFO) and TTL expiration. Built as a learning project mostly for Go standard library (http, sync packages etc).

## Structure

- `cmd/upstream` — dummy JSON server with `fast/{id}` and `slow/{id}` routes to test the proxy
- `cmd/proxy` — the proxy server itself
- `internal/cache` — the `Cache` interface and all of its implementations
- `internal/proxy` — the HTTP handler that ties caching to request forwarding

## Design notes

- FIFO and LRU both implemented with a doubly linked list + map
- FIFO `Set` puts the new key or moves the existing key to the front of the list to be evicted the latest; LRU does this but on both `Get` and `Set`
- TTL is checked lazily on read and does not revalidate
- Cache only stores HTTP body, not headers or status codes
- Only supports single upstream server specified as a flag to the proxy

## Running the program(s)

```bash
go run ./cmd/upstream -port 8080

go run ./cmd/proxy -port 8081 -policy lru -capacity 100 -upstream http://localhost:8080 -ttl 60
```

All given flag values are defaults. They can be omitted. Use -h flag for more info.

You can then experiment with a series of curl commands and see that the cache hit rate will be the expected one for the given eviction policy. Also note that cache hits on the slow endpoints will not have a delay. For example using LRU with capacity 2:

```bash
curl http://localhost:8081/slow/1 # miss
curl http://localhost:8081/slow/2 # miss
curl http://localhost:8081/slow/1 # hit, moves "1" to front
curl http://localhost:8081/slow/3 # miss, evicts "2"
curl http://localhost:8081/slow/1 # hit, moves "1" to front again
curl http://localhost:8081/slow/2 # miss, evicts "3"
curl http://localhost:8081/slow/stats # {"hits":2,"misses":4}
```

## Potential Future Work

- Generalize the proxy to support multiple arbitrary upstreams during runtime
- Read `Cache-Control`/`max-age` from upstream responses instead of some fixed TTL
- Support more eviction policies such as LFU
- Background TTL sweep instead of lazy expiration
- Implement solution to cache stampedes 
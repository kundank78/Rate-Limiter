Rate Limiter
==========

This package offers a Go implementation of Rate Limit Token Bucket and Sliding Window algorithms. 
It also supports a Redis-backed sliding window approach, utilizing the sorted sets data structure for distributed rate limiting.

### Usage

Initialize the rate limiter with max number of requests in a given time window in seconds. 
Sliding Window Distributed option uses redis sorted sets to maintain user request timestamps (Redis server config are currently hardcoded)
```go
limiter := pkg.Init(pkg.SlidingWindowDistributedAlgo, 10, 60)
router.GET("/albums", pkg.RateLimiter(limiter, getAlbums))
```

### Example

Execute `main.go` which runs a simple album server at `localhost:8000`. You can try different rate limiter options by modifying the 
rate limiter init() params.

### Next Steps
- Add support for user identifiers other than ip address 
- Add support for approximated sliding window algo 
- Support frameworks other than gin

package pkg

import (
	"sync"
	"time"
)

type TokenBucket struct {
	tokensMap  map[string]int
	maxRequest int
	windowTime int
	lock       sync.Mutex
}

func InitTokenBucket(maxRequest int, windowTime int) *TokenBucket {
	tokenBucket := &TokenBucket{
		tokensMap:  make(map[string]int),
		maxRequest: maxRequest,
		windowTime: windowTime,
	}
	go tokenBucket.tokenGeneration()

	return tokenBucket
}

func (tb *TokenBucket) isAllowed(userId string) bool {
	tb.lock.Lock()
	defer tb.lock.Unlock()
	if _, ok := tb.tokensMap[userId]; ok {
		tokens := tb.tokensMap[userId]
		if tokens > 0 {
			tb.tokensMap[userId] = tokens - 1
			return true
		}
		return false
	}
	tb.tokensMap[userId] = tb.maxRequest - 1
	return true
}

func (tb *TokenBucket) tokenGeneration() {
	for {
		select {
		case <-time.After(time.Duration(tb.windowTime) * time.Second):
			tb.lock.Lock()
			for k, _ := range tb.tokensMap {
				tb.tokensMap[k] = tb.maxRequest
			}
			tb.lock.Unlock()
		}
	}
}

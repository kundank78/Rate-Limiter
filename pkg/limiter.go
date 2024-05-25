package pkg

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

const (
	SlidingWindowAlgo            = "slidingWindow"
	SlidingWindowDistributedAlgo = "slidingWindowDistributed"
	TokenBucketAlgo              = "tokenBucket"
)

type Limiter interface {
	isAllowed(userId string) bool
}

func Init(rateLimiterType string, maxRequest int, timeWindow int) Limiter {
	switch rateLimiterType {
	case SlidingWindowAlgo:
		return InitSlidingWindow(maxRequest, timeWindow)
	case TokenBucketAlgo:
		return InitTokenBucket(maxRequest, timeWindow)
	case SlidingWindowDistributedAlgo:
		return InitSlidingWindowDistributed(maxRequest, timeWindow)
	default:
		return InitTokenBucket(maxRequest, timeWindow)
	}
}

func RateLimiter(limiter Limiter, next func(c *gin.Context)) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestRemoteAddr := c.Request.RemoteAddr
		requestRemoteAddr = strings.Split(requestRemoteAddr, ":")[0]
		if !limiter.isAllowed(requestRemoteAddr) {
			c.IndentedJSON(http.StatusTooManyRequests, gin.H{"message": "Request Limit Exceeded"})
			return
		} else {
			next(c)
		}
	}
}

func GetCurrentTime() int64 {
	return time.Now().UnixMilli()
}

package pkg

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

type SlidingWindowDistributed struct {
	redis      *redis.Client
	maxRequest int
	windowTime int
}

func InitSlidingWindowDistributed(maxRequest int, windowTime int) *SlidingWindowDistributed {
	swd := &SlidingWindowDistributed{
		redis:      InitRedis(),
		maxRequest: maxRequest,
		windowTime: windowTime,
	}
	return swd
}

func (swd *SlidingWindowDistributed) isAllowed(userId string) bool {
	currTime := GetCurrentTime()
	windowStart := currTime - int64(swd.windowTime*1000)

	ctx := context.Background()
	res := RedisClient.ZCount(ctx, userId, strconv.FormatInt(windowStart, 10), strconv.FormatInt(currTime, 10))
	val, err := res.Result()
	if err != nil {
		fmt.Printf("Error querying redis for user %v | error -> %v", userId, err)
		return true
	}

	if val < int64(swd.maxRequest) {
		res = RedisClient.ZAdd(ctx, userId, redis.Z{float64(currTime), currTime})
		if err != nil {
			fmt.Printf("Error adding new timestamp in redis for user %v | error -> %v", userId, err)
			return true
		}
		go func() {
			RedisClient.Expire(ctx, userId, time.Duration(swd.windowTime)*time.Second)
			res = RedisClient.ZRemRangeByScore(ctx, userId, "-inf", strconv.FormatInt(windowStart, 10))
			if err != nil {
				fmt.Printf("Error removing old timestamp in redis for user %v | error -> %v", userId, err)
			}
		}()

		return true
	}
	return false
}

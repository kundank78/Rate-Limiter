package pkg

import (
	"fmt"
	"sync"
	"time"
)

type SlidingWindow struct {
	tokensMap      map[string]map[int64]int
	maxRequest     int
	windowTime     int
	lock           sync.RWMutex
	cleanFrequency int
}

func InitSlidingWindow(maxRequest int, windowTime int) *SlidingWindow {
	sw := &SlidingWindow{
		tokensMap:      make(map[string]map[int64]int),
		maxRequest:     maxRequest,
		windowTime:     windowTime,
		lock:           sync.RWMutex{},
		cleanFrequency: 60,
	}
	go sw.cleanUserRequestMap()
	return sw
}

func (sw *SlidingWindow) isAllowed(userId string) bool {

	unixTime := GetCurrentTime()

	if sw.calculateTotalRequestCountInWindow(userId) < sw.maxRequest {
		sw.lock.Lock()
		defer sw.lock.Unlock()
		if _, ok := sw.tokensMap[userId]; ok {
			userTimestampCount := sw.tokensMap[userId]
			if _, ok := userTimestampCount[unixTime]; ok {
				count := userTimestampCount[unixTime]
				userTimestampCount[unixTime] = count + 1
			} else {
				userTimestampCount[unixTime] = 1
			}
		} else {
			userTimestampCount := make(map[int64]int)
			userTimestampCount[unixTime] = 1
			sw.tokensMap[userId] = userTimestampCount
		}
		return true
	}
	return false
}

func (sw *SlidingWindow) calculateTotalRequestCountInWindow(userId string) int {
	sw.lock.RLock()
	defer sw.lock.RUnlock()
	requestCount := 0
	if _, ok := sw.tokensMap[userId]; ok {
		userTimestampCount := sw.tokensMap[userId]
		currTime := GetCurrentTime()

		for timestamp, count := range userTimestampCount {
			if currTime-timestamp < int64(sw.windowTime) {
				requestCount = requestCount + count
			}
		}
	}
	return requestCount
}

func (sw *SlidingWindow) cleanUserRequestMap() {
	for {
		time.Sleep(time.Duration(sw.cleanFrequency) * time.Second)

		fmt.Println("Running Sliding Window Clean UP!!!")
		currTime := GetCurrentTime()

		deleteTimestampList := make(map[int64]bool)
		sw.lock.Lock()
		for userId, userTimestampCount := range sw.tokensMap {
			for timestamp, _ := range userTimestampCount {
				if currTime-timestamp > int64(sw.windowTime) {
					deleteTimestampList[timestamp] = true
				}
			}
			for timestamp, _ := range deleteTimestampList {
				fmt.Printf("Deleting timestamp %v for user -> %s \n", timestamp, userId)
				delete(userTimestampCount, timestamp)
			}
		}
		sw.lock.Unlock()
	}
}

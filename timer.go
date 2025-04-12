package main

import (
	"sync"
	"time"
)

type TimerID int

var (
	timerIDCounter TimerID = 0
	timeouts               = make(map[TimerID]*time.Timer)
	intervals              = make(map[TimerID]*time.Ticker)
	mu             sync.Mutex
)

func GsetTimeout(callback func(), delay time.Duration) TimerID {
	mu.Lock()
	defer mu.Unlock()

	timerIDCounter++
	id := timerIDCounter

	timeouts[id] = time.AfterFunc(delay, func() {
		callback()
		mu.Lock()
		delete(timeouts, id)
		mu.Unlock()
	})

	return id
}

func GclearTimeout(id TimerID) {
	mu.Lock()
	defer mu.Unlock()

	if timer, exists := timeouts[id]; exists {
		timer.Stop()
		delete(timeouts, id)
	}
}

func GsetInterval(callback func(), interval time.Duration) TimerID {
	mu.Lock()
	defer mu.Unlock()

	timerIDCounter++
	id := timerIDCounter

	intervals[id] = time.NewTicker(interval)

	go func() {
		for range intervals[id].C {
			callback()
		}
	}()

	return id
}

func GclearInterval(id TimerID) {
	mu.Lock()
	defer mu.Unlock()

	if ticker, exists := intervals[id]; exists {
		ticker.Stop()
		delete(intervals, id)
	}
}

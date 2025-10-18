package main

import (
	"sync"
	"time"
)

type RateLimiter struct {
	clients           map[string]*ClientInfo
	mutex             sync.RWMutex
	cleanup           *time.Ticker
	maxRequestsPerMin int
}

type ClientInfo struct {
	requests  int
	resetTime time.Time
	lastSeen  time.Time
}

func NewRateLimiter(maxRequestsPerMin int) *RateLimiter {
	rl := &RateLimiter{
		clients:           make(map[string]*ClientInfo),
		cleanup:           time.NewTicker(time.Minute),
		maxRequestsPerMin: maxRequestsPerMin,
	}

	go rl.cleanupExpiredClients()
	return rl
}

func (rl *RateLimiter) IsAllowed(clientIP string) bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	client, exists := rl.clients[clientIP]

	if !exists {
		rl.clients[clientIP] = &ClientInfo{
			requests:  1,
			resetTime: now.Add(time.Minute),
			lastSeen:  now,
		}
		return true
	}

	client.lastSeen = now

	if now.After(client.resetTime) {
		client.requests = 1
		client.resetTime = now.Add(time.Minute)
		return true
	}

	if client.requests < rl.maxRequestsPerMin {
		client.requests++
		return true
	}

	return false
}

func (rl *RateLimiter) cleanupExpiredClients() {
	for range rl.cleanup.C {
		rl.mutex.Lock()
		now := time.Now()
		for ip, client := range rl.clients {
			if now.Sub(client.lastSeen) > 10*time.Minute {
				delete(rl.clients, ip)
			}
		}
		rl.mutex.Unlock()
	}
}

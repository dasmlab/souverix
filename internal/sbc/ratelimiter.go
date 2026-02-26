package sbc

import (
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// RateLimiter implements per-IP rate limiting
type RateLimiter struct {
	limit   int
	window  time.Duration
	clients map[string]*clientLimit
	mu      sync.RWMutex
	log     *logrus.Logger
}

type clientLimit struct {
	count     int
	windowStart time.Time
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(limit int, window time.Duration, log *logrus.Logger) *RateLimiter {
	rl := &RateLimiter{
		limit:   limit,
		window:  window,
		clients: make(map[string]*clientLimit),
		log:     log,
	}

	// Cleanup old entries periodically
	go rl.cleanup()

	return rl
}

// Allow checks if a request from an IP should be allowed
func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	client, exists := rl.clients[ip]

	if !exists {
		rl.clients[ip] = &clientLimit{
			count:       1,
			windowStart: now,
		}
		return true
	}

	// Reset window if it has expired
	if now.Sub(client.windowStart) > rl.window {
		client.count = 1
		client.windowStart = now
		return true
	}

	// Check limit
	if client.count >= rl.limit {
		return false
	}

	client.count++
	return true
}

// cleanup removes old entries periodically
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for ip, client := range rl.clients {
			if now.Sub(client.windowStart) > 2*rl.window {
				delete(rl.clients, ip)
			}
		}
		rl.mu.Unlock()
	}
}


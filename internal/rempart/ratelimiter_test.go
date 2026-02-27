package rempart

import (
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

func TestNewRateLimiter(t *testing.T) {
	log := logrus.New()
	log.SetLevel(logrus.ErrorLevel)

	rl := NewRateLimiter(100, 60*time.Second, log)
	if rl == nil {
		t.Error("NewRateLimiter() returned nil")
	}
}

func TestRateLimiter_Allow(t *testing.T) {
	log := logrus.New()
	log.SetLevel(logrus.ErrorLevel)

	rl := NewRateLimiter(5, 1*time.Second, log)

	// First 5 should be allowed
	for i := 0; i < 5; i++ {
		if !rl.Allow("192.168.1.1") {
			t.Errorf("RateLimiter.Allow() should allow request %d", i+1)
		}
	}

	// 6th should be blocked
	if rl.Allow("192.168.1.1") {
		t.Error("RateLimiter.Allow() should block after limit")
	}
}

func TestRateLimiter_DifferentIPs(t *testing.T) {
	log := logrus.New()
	log.SetLevel(logrus.ErrorLevel)

	rl := NewRateLimiter(5, 1*time.Second, log)

	// Different IPs should have separate limits
	if !rl.Allow("192.168.1.1") {
		t.Error("RateLimiter.Allow() should allow for IP 1")
	}
	if !rl.Allow("192.168.1.2") {
		t.Error("RateLimiter.Allow() should allow for IP 2")
	}
}

func TestRateLimiter_WindowReset(t *testing.T) {
	log := logrus.New()
	log.SetLevel(logrus.ErrorLevel)

	rl := NewRateLimiter(2, 100*time.Millisecond, log)

	// Use up limit
	rl.Allow("192.168.1.1")
	rl.Allow("192.168.1.1")

	// Should be blocked
	if rl.Allow("192.168.1.1") {
		t.Error("RateLimiter.Allow() should block after limit")
	}

	// Wait for window to reset
	time.Sleep(150 * time.Millisecond)

	// Should be allowed again
	if !rl.Allow("192.168.1.1") {
		t.Error("RateLimiter.Allow() should allow after window reset")
	}
}

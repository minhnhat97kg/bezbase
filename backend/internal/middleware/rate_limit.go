package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

// RateLimitConfig holds the rate limiting configuration
type RateLimitConfig struct {
	MaxRequests int           // Maximum number of requests per window
	Window      time.Duration // Time window for rate limiting
	KeyFunc     func(c echo.Context) string // Function to extract key from context
}

// DefaultRateLimitConfig returns a default rate limit configuration
func DefaultRateLimitConfig() RateLimitConfig {
	return RateLimitConfig{
		MaxRequests: 100,
		Window:      time.Minute,
		KeyFunc: func(c echo.Context) string {
			return c.RealIP()
		},
	}
}

// RateLimiter holds the rate limiting state
type RateLimiter struct {
	config  RateLimitConfig
	clients map[string]*clientInfo
	mutex   sync.RWMutex
}

// clientInfo stores information about a client's requests
type clientInfo struct {
	requests  int
	resetTime time.Time
	mutex     sync.RWMutex
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(config RateLimitConfig) *RateLimiter {
	rl := &RateLimiter{
		config:  config,
		clients: make(map[string]*clientInfo),
	}
	
	// Start cleanup goroutine
	go rl.cleanup()
	
	return rl
}

// cleanup removes expired client entries
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(rl.config.Window)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			rl.mutex.Lock()
			now := time.Now()
			for key, client := range rl.clients {
				client.mutex.RLock()
				if now.After(client.resetTime) {
					delete(rl.clients, key)
				}
				client.mutex.RUnlock()
			}
			rl.mutex.Unlock()
		}
	}
}

// Allow checks if a request should be allowed
func (rl *RateLimiter) Allow(key string) (bool, int, time.Time) {
	rl.mutex.RLock()
	client, exists := rl.clients[key]
	rl.mutex.RUnlock()
	
	if !exists {
		rl.mutex.Lock()
		client = &clientInfo{
			requests:  0,
			resetTime: time.Now().Add(rl.config.Window),
		}
		rl.clients[key] = client
		rl.mutex.Unlock()
	}
	
	client.mutex.Lock()
	defer client.mutex.Unlock()
	
	now := time.Now()
	if now.After(client.resetTime) {
		client.requests = 0
		client.resetTime = now.Add(rl.config.Window)
	}
	
	if client.requests >= rl.config.MaxRequests {
		return false, client.requests, client.resetTime
	}
	
	client.requests++
	return true, client.requests, client.resetTime
}

// RateLimit returns a middleware that implements rate limiting
func RateLimit() echo.MiddlewareFunc {
	return RateLimitWithConfig(DefaultRateLimitConfig())
}

// RateLimitWithConfig returns a rate limiting middleware with custom configuration
func RateLimitWithConfig(config RateLimitConfig) echo.MiddlewareFunc {
	limiter := NewRateLimiter(config)
	
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			key := config.KeyFunc(c)
			allowed, requests, resetTime := limiter.Allow(key)
			
			// Set rate limit headers
			c.Response().Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", config.MaxRequests))
			c.Response().Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", config.MaxRequests-requests))
			c.Response().Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", resetTime.Unix()))
			
			if !allowed {
				return echo.NewHTTPError(http.StatusTooManyRequests, "Rate limit exceeded")
			}
			
			return next(c)
		}
	}
}

// AuthRateLimit provides rate limiting for authentication endpoints
func AuthRateLimit() echo.MiddlewareFunc {
	config := RateLimitConfig{
		MaxRequests: 5,
		Window:      time.Minute,
		KeyFunc: func(c echo.Context) string {
			return c.RealIP()
		},
	}
	return RateLimitWithConfig(config)
}

// APIRateLimit provides rate limiting for API endpoints
func APIRateLimit() echo.MiddlewareFunc {
	config := RateLimitConfig{
		MaxRequests: 100,
		Window:      time.Minute,
		KeyFunc: func(c echo.Context) string {
			return c.RealIP()
		},
	}
	return RateLimitWithConfig(config)
}
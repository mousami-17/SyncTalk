package middleware

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"realtime-chat/src/cache"
)

// Rate limit configuration
const (
	// Auth endpoints - stricter limits
	AuthRateLimit     = 50              // 5 attempts
	AuthRateWindow    = 15 * time.Minute // per 15 minutes
	
	// Message endpoints - moderate limits
	MessageRateLimit  = 100            // 100 messages
	MessageRateWindow = 1 * time.Minute // per minute
	
	// General API - lenient limits
	GeneralRateLimit  = 60             // 60 requests
	GeneralRateWindow = 1 * time.Minute // per minute
)

// RateLimitConfig holds rate limit settings
type RateLimitConfig struct {
	Limit  int
	Window time.Duration
	KeyPrefix string
}

// RateLimiter creates a rate limiting middleware
func RateLimiter(config RateLimitConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get client identifier (IP address)
		clientIP := c.IP()
		
		// Create Redis key
		key := fmt.Sprintf("ratelimit:%s:%s", config.KeyPrefix, clientIP)
		
		// Check rate limit
		allowed, remaining, resetTime, err := checkRateLimit(key, config.Limit, config.Window)
		if err != nil {
			log.Printf("Rate limit check error: %v", err)
			// On error, allow the request but log it
			return c.Next()
		}

		// Set rate limit headers
		c.Set("X-RateLimit-Limit", fmt.Sprintf("%d", config.Limit))
		c.Set("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
		c.Set("X-RateLimit-Reset", fmt.Sprintf("%d", resetTime.Unix()))

		if !allowed {
			retryAfter := time.Until(resetTime).Seconds()
			c.Set("Retry-After", fmt.Sprintf("%.0f", retryAfter))
			
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":   true,
				"message": "Too many requests. Please try again later.",
				"retryAfter": retryAfter,
			})
		}

		return c.Next()
	}
}

// checkRateLimit checks if request is within rate limit
func checkRateLimit(key string, limit int, window time.Duration) (allowed bool, remaining int, resetTime time.Time, err error) {
	ctx := context.Background()
	
	// Increment counter
	count, err := cache.RedisClient.Incr(ctx, key).Result()
	if err != nil {
		return false, 0, time.Time{}, err
	}

	// Set expiration on first request
	if count == 1 {
		cache.RedisClient.Expire(ctx, key, window)
	}

	// Get TTL for reset time
	ttl, err := cache.RedisClient.TTL(ctx, key).Result()
	if err != nil {
		return false, 0, time.Time{}, err
	}

	resetTime = time.Now().Add(ttl)
	remaining = limit - int(count)
	if remaining < 0 {
		remaining = 0
	}

	allowed = count <= int64(limit)
	return allowed, remaining, resetTime, nil
}

// AuthRateLimiter for authentication endpoints
func AuthRateLimiter() fiber.Handler {
	return RateLimiter(RateLimitConfig{
		Limit:     AuthRateLimit,
		Window:    AuthRateWindow,
		KeyPrefix: "auth",
	})
}

// MessageRateLimiter for message sending
func MessageRateLimiter() fiber.Handler {
	return RateLimiter(RateLimitConfig{
		Limit:     MessageRateLimit,
		Window:    MessageRateWindow,
		KeyPrefix: "message",
	})
}

// GeneralRateLimiter for general API endpoints
func GeneralRateLimiter() fiber.Handler {
	return RateLimiter(RateLimitConfig{
		Limit:     GeneralRateLimit,
		Window:    GeneralRateWindow,
		KeyPrefix: "general",
	})
}

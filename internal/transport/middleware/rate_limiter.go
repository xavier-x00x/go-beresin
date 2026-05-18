package middleware

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

// RateLimiter holds the Redis client and defines limit rules.
type RateLimiter struct {
	rdb *redis.Client
}

// NewRateLimiter creates a new RateLimiter instance.
func NewRateLimiter(rdb *redis.Client) *RateLimiter {
	return &RateLimiter{rdb: rdb}
}

// Lua script to atomically increment and retrieve TTL
const rateLimitLua = `
local key = KEYS[1]
local limit = tonumber(ARGV[1])
local duration = tonumber(ARGV[2])

local current = redis.call('get', key)
if current and tonumber(current) >= limit then
    local ttl = redis.call('ttl', key)
    if ttl < 0 then
        redis.call('expire', key, duration)
        ttl = duration
    end
    return {tonumber(current) + 1, ttl}
end

local newVal = redis.call('incr', key)
local ttl = redis.call('ttl', key)
if newVal == 1 or ttl < 0 then
    redis.call('expire', key, duration)
    ttl = duration
end
return {newVal, ttl}
`

// Limit applies a generic rate limiting rule.
func (rl *RateLimiter) Limit(group string, limit int, duration time.Duration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// If Redis client is nil, skip rate limiting (fail-open or log)
		if rl.rdb == nil {
			return c.Next()
		}

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		// Identify the client by IP address
		clientIP := c.IP()
		key := fmt.Sprintf("rate_limit:%s:%s", group, clientIP)

		durationSeconds := int(duration.Seconds())

		// Execute the Lua script atomically
		res, err := rl.rdb.Eval(ctx, rateLimitLua, []string{key}, limit, durationSeconds).Result()
		if err != nil {
			// Fail open and log error to avoid blocking valid traffic if Redis experiences issues
			log.Printf("[WARNING] [RateLimiter] Redis error: %v", err)
			return c.Next()
		}

		slice, ok := res.([]interface{})
		if !ok || len(slice) < 2 {
			return c.Next()
		}

		count := slice[0].(int64)
		ttl := slice[1].(int64)

		remaining := int64(limit) - count
		if remaining < 0 {
			remaining = 0
		}

		// Set rate limiting headers in response
		c.Set("X-RateLimit-Limit", strconv.Itoa(limit))
		c.Set("X-RateLimit-Remaining", strconv.FormatInt(remaining, 10))
		c.Set("X-RateLimit-Reset", strconv.FormatInt(ttl, 10))

		// If the limit has been exceeded, return 429 Too Many Requests
		if count > int64(limit) {
			return c.Status(fiber.StatusTooManyRequests).JSON(Response{
				Status:  "error",
				Message: "Too many requests. Please try again later.",
				Data: fiber.Map{
					"limit":             limit,
					"remaining":         0,
					"retry_after_sec":   ttl,
					"rate_limit_group":  group,
				},
			})
		}

		return c.Next()
	}
}

// GlobalLimit enforces general rate limiting (default: 60 requests per minute).
func (rl *RateLimiter) GlobalLimit() fiber.Handler {
	return rl.Limit("global", 60, 1*time.Minute)
}

// StrictLimit enforces very strict rate limiting for sensitive endpoints (default: 5 requests per minute).
func (rl *RateLimiter) StrictLimit() fiber.Handler {
	return rl.Limit("strict", 5, 1*time.Minute)
}

// IsIPBlocked checks if an IP address is currently blocked from logging in.
// Returns true and the remaining block duration if blocked.
func (rl *RateLimiter) IsIPBlocked(ctx context.Context, ip string) (bool, time.Duration, error) {
	if rl.rdb == nil {
		return false, 0, nil
	}
	key := fmt.Sprintf("auth_block:%s", ip)
	ttl, err := rl.rdb.TTL(ctx, key).Result()
	if err != nil {
		return false, 0, err
	}
	if ttl > 0 {
		return true, ttl, nil
	}
	return false, 0, nil
}

// RecordFailedLogin increments the failed login counter for an IP.
// If it reaches 5, blocks the IP for 15 minutes.
func (rl *RateLimiter) RecordFailedLogin(ctx context.Context, ip string) (int64, error) {
	if rl.rdb == nil {
		return 0, nil
	}
	counterKey := fmt.Sprintf("auth_failed_counter:%s", ip)
	blockKey := fmt.Sprintf("auth_block:%s", ip)

	count, err := rl.rdb.Incr(ctx, counterKey).Result()
	if err != nil {
		return 0, err
	}

	// Set initial TTL of 1 hour for the failed counter
	if count == 1 {
		rl.rdb.Expire(ctx, counterKey, 1*time.Hour)
	}

	if count >= 5 {
		// Block IP for 15 minutes
		err = rl.rdb.Set(ctx, blockKey, "blocked", 15*time.Minute).Err()
		if err != nil {
			return count, err
		}
		// Reset the failed counter
		rl.rdb.Del(ctx, counterKey)
	}

	return count, nil
}

// ResetFailedLogin clears the failed login counter for an IP upon successful login.
func (rl *RateLimiter) ResetFailedLogin(ctx context.Context, ip string) error {
	if rl.rdb == nil {
		return nil
	}
	counterKey := fmt.Sprintf("auth_failed_counter:%s", ip)
	return rl.rdb.Del(ctx, counterKey).Err()
}

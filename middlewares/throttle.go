package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	models "throttling-api/models/db"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type Throttle interface {
	Allow(ctx context.Context, clientID, endpoint string, limit int) (bool, error)
}

type SlidingLogThrottle struct {
	mux sync.Mutex

	windowDuration time.Duration
	requests       map[string][]time.Time //in memory storate
}

func NewSlidingLogThrottle(duration time.Duration) SlidingLogThrottle {
	return SlidingLogThrottle{
		windowDuration: duration,
		requests:       map[string][]time.Time{},
	}
}

func (t SlidingLogThrottle) Allow(ctx context.Context, clientID, endpoint string, limit int) (bool, error) {
	t.mux.Lock()
	defer t.mux.Unlock()

	now := time.Now()
	key := clientID + ":" + endpoint
	requests := t.requests[key]

	filteredRequests := []time.Time{}
	for _, req := range requests {
		if now.Sub(req) <= t.windowDuration {
			filteredRequests = append(filteredRequests, req)
		}
	}

	if len(filteredRequests) >= limit {
		return false, nil
	}

	t.requests[key] = append(filteredRequests, now)
	return true, nil
}

type FixedWindowThrottle struct {
	windowDuration time.Duration

	redisClient *redis.Client //redis storage
}

func NewFixedWindowThrottle(duration time.Duration, redisClient *redis.Client) FixedWindowThrottle {
	return FixedWindowThrottle{
		windowDuration: duration,
		redisClient:    redisClient,
	}
}

func (t *FixedWindowThrottle) getRedisKey(clientID, endpoint string) string {
	now := time.Now()
	windowStart := now.Truncate(t.windowDuration).Unix()
	return fmt.Sprintf("fixed_window:%s:%s::%d", clientID, endpoint, windowStart)
}

func (t FixedWindowThrottle) Allow(ctx context.Context, clientID, endpoint string, limit int) (bool, error) {
	redisKey := t.getRedisKey(clientID, endpoint)
	expireSeconds := int(t.windowDuration.Seconds())

	script := redis.NewScript(`
		local current = redis.call("INCR", KEYS[1])
		if current == 1 then
			redis.call("EXPIRE", KEYS[1], ARGV[1])
		end
		return current
	`)

	// Run atomically: INCR + EXPIRE if first increment
	res, err := script.Run(ctx, t.redisClient, []string{redisKey}, expireSeconds).Result()
	if err != nil {
		return false, err
	}

	count := res.(int64)
	return count <= int64(limit), nil
}

func ThrottleMiddleware(throttle Throttle, endpoint string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			abortWithError(c, http.StatusUnauthorized, "unauthorized")
			return
		}

		model := user.(models.User)
		allow, err := throttle.Allow(c.Request.Context(), model.ClientID, endpoint, model.RateLimit)
		if err != nil {
			abortWithError(c, http.StatusInternalServerError, "internal error")
			return
		}

		if !allow {
			abortWithError(c, http.StatusTooManyRequests, "rate limit exceeded")
			return
		}

		c.Next()
	}
}

package redis

import (
	"context"
	"os"
	"sync"

	"github.com/charmbracelet/log"
	"github.com/gomodule/redigo/redis"
)

var (
	redisPool *redis.Pool
	once      sync.Once
)

// GetRedisPool returns a shared Redis connection pool
func GetRedisPool() *redis.Pool {
	if redisPool == nil {
		once.Do(func() {
			redisURL := os.Getenv("REDIS_URL")
			if redisURL == "" {
				log.Fatal("REDIS_URL environment variable is required")
			}
			redisPool = &redis.Pool{
				MaxIdle:     10,
				MaxActive:   100,
				IdleTimeout: 240,
				Dial: func() (redis.Conn, error) {
					return redis.DialURL(redisURL)
				},
				DialContext: func(ctx context.Context) (redis.Conn, error) {
					return redis.DialURLContext(ctx, redisURL)
				},
			}
		})
	}
	return redisPool
}

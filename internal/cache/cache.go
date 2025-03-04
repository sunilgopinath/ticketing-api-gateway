package cache

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

// InitRedis initializes Redis client
func InitRedis() {
	redisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Redis server address
	})
	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
}

// GetCache retrieves cached response
func GetCache(ctx context.Context, key string) (string, error) {
	val, err := redisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil // Cache miss
	}
	return val, err
}

// SetCache stores API response in Redis
func SetCache(ctx context.Context, key, value string, ttl time.Duration) error {
	return redisClient.Set(ctx, key, value, ttl).Err()
}

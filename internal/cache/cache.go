package cache

import (
	"context"
	"log"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

var (
	cacheHitsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_hits_total",
			Help: "Total number of cache hits",
		},
		[]string{"endpoint", "instance"}, // Two labels
	)
	cacheMissesTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_misses_total",
			Help: "Total number of cache misses",
		},
		[]string{"endpoint", "instance"}, // Two labels
	)
)

func init() {
	prometheus.MustRegister(cacheHitsTotal)
	prometheus.MustRegister(cacheMissesTotal)
}

func InitRedis() {
	redisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
}

func GetCache(ctx context.Context, key string, endpoint, instance string) (string, error) {
	val, err := redisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		cacheMissesTotal.WithLabelValues(endpoint, instance).Inc()
		return "", nil
	}
	if err != nil {
		return "", err
	}
	cacheHitsTotal.WithLabelValues(endpoint, instance).Inc()
	return val, err
}

func SetCache(ctx context.Context, key, value string, ttl time.Duration, endpoint, instance string) error {
	return redisClient.Set(ctx, key, value, ttl).Err()
}

package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheRepository struct {
	client *redis.Client
}

func NewCacheRepository(addr string) *CacheRepository {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	return &CacheRepository{client: client}
}

func (r *CacheRepository) Set(key string, value string, expiration time.Duration) error {
	ctx := context.Background()
	return r.client.Set(ctx, key, value, expiration).Err()
}

func (r *CacheRepository) Get(key string) (string, error) {
	ctx := context.Background()
	return r.client.Get(ctx, key).Result()
}

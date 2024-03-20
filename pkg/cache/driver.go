package cache

import (
	"context"
	"fmt"
	"github.com/patrickmn/go-cache"
	"github.com/redis/go-redis/v9"
	"time"
)

type redisDriver struct {
	client *redis.Client
}

func (c redisDriver) Get(ctx context.Context, key string) (interface{}, error) {
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return val, nil
}

func (c redisDriver) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.client.Set(ctx, key, value, expiration).Err()
}

func (c redisDriver) Forget(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

func newRedisDriver(host string, port, db int) *redisDriver {
	return &redisDriver{
		client: redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("%s:%d", host, port),
			DB:   db,
		}),
	}
}

type inMemoryDriver struct {
	c *cache.Cache
}

func (i inMemoryDriver) Get(ctx context.Context, key string) (interface{}, error) {
	value, found := i.c.Get(key)
	if found {
		return value, nil
	}
	return nil, nil
}

func (i inMemoryDriver) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	i.c.Set(key, value, expiration)
	return nil
}

func (i inMemoryDriver) Forget(ctx context.Context, key string) error {
	i.c.Delete(key)
	return nil
}

func newInMemoryDriver(defaultExpiration, cleanupInterval time.Duration) *inMemoryDriver {
	c := cache.New(defaultExpiration, cleanupInterval)
	return &inMemoryDriver{
		c: c,
	}
}

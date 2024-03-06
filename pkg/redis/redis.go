package redis

import (
	"context"
	"github.com/kurneo/go-template/config"
	"github.com/redis/go-redis/v9"
	"strconv"
	"sync"
	"time"
)

var (
	redisOnce     sync.Once
	redisInstance Contact
)

const Nil = redis.Nil

type Contact interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	GetClient() *redis.Client
}

type Client struct {
	client *redis.Client
}

func (c Client) Get(ctx context.Context, key string) (string, error) {
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return val, Nil
		} else {
			return val, err
		}
	}
	return val, nil
}

func (c Client) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.client.Set(ctx, key, value, expiration).Err()
}

func (c Client) GetClient() *redis.Client {
	return c.client
}

func NewRedisClient(cfg config.Redis) (Contact, error) {
	var err error = nil
	db, err := strconv.Atoi(cfg.DB)
	if err != nil {
		return nil, err
	}

	redisOnce.Do(func() {
		redisInstance = &Client{
			client: redis.NewClient(&redis.Options{
				Addr: cfg.Addr,
				DB:   db,
			}),
		}
	})

	return redisInstance, err
}

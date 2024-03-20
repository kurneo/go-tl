package cache

import (
	"context"
	"errors"
	"sync"
	"time"
)

type Contact interface {
	Get(ctx context.Context, key string) (interface{}, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Forget(ctx context.Context, key string) error
}

type Config struct {
	Driver string
	Redis  struct {
		Host string
		Port int
		DB   int
	}
	InMemory struct {
		DefaultExpiration      time.Duration
		DefaultCleanUpInterval time.Duration
	}
}

var (
	cacheOnce     sync.Once
	cacheInstance Contact
)

const (
	DriverRedis    = "redis"
	DriverInMemory = "in-memory"
)

func New(c Config) (Contact, error) {
	var err error

	cacheOnce.Do(func() {
		switch c.Driver {
		case DriverRedis:
			cacheInstance = newRedisDriver(c.Redis.Host, c.Redis.Port, c.Redis.DB)
		case DriverInMemory:
			cacheInstance = newInMemoryDriver(c.InMemory.DefaultExpiration, c.InMemory.DefaultCleanUpInterval)
		default:
			err = errors.New("cache driver is invalid")
		}
	})

	if err != nil {
		return nil, err
	}

	return cacheInstance, err
}

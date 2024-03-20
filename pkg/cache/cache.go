package cache

import (
	"context"
	"errors"
	"github.com/spf13/viper"
	"sync"
	"time"
)

type Contact interface {
	Get(ctx context.Context, key string) (interface{}, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Forget(ctx context.Context, key string) error
}

var (
	cacheOnce     sync.Once
	cacheInstance Contact
)

const (
	DriverRedis    = "redis"
	DriverInMemory = "in-memory"
)

func New() (Contact, error) {
	var err error

	driver := viper.GetString("CACHE_DRIVER")

	cacheOnce.Do(func() {
		switch driver {
		case DriverRedis:
			cacheInstance = newRedisDriver(
				viper.GetString("CACHE_REDIS_HOST"),
				viper.GetInt("CACHE_REDIS_PORT"),
				viper.GetInt("CACHE_REDIS_DB"),
			)
		case DriverInMemory:
			cacheInstance = newInMemoryDriver(
				viper.GetDuration("CACHE_IN_MEMORY_DEFAULT_EXPIRATION"),
				viper.GetDuration("CACHE_IN_MEMORY_CLEANUP_INTERVAL"))
		default:
			err = errors.New("cache driver is invalid")
		}
	})

	if err != nil {
		return nil, err
	}

	return cacheInstance, err
}

package store

import (
	"context"
	"fmt"
	"time"

	"github.com/nortoo/usms/internal/pkg/etc"
	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client

func GetRedisClient() *redis.Client {
	return rdb
}

// InitRedis creates a new Redis client.
func InitRedis(config *etc.Store) error {
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port),
		Password: config.Redis.Password, // no password set
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	return err
}

package store

import (
	"context"
	"fmt"
	"time"

	"github.com/nortoo/usms/internal/pkg/etc"
	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client

type RedisCli struct {
	rdb *redis.Client
}

// NewRedisCli creates a new Redis client.
func NewRedisCli(config *etc.Store) (*RedisCli, error) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port),
		Password: config.Redis.Password, // no password set
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return &RedisCli{rdb: rdb}, nil
}

func (c *RedisCli) GetRDB() *redis.Client {
	return c.rdb
}

func (c *RedisCli) ListKeys(ctx context.Context, key string) (keys []string, err error) {
	var cursor uint64

	// Loop until the cursor returns to 0
	for {
		results, nextCursor, err := c.GetRDB().Scan(ctx, cursor, key, 0).Result()
		if err != nil {
			return nil, err
		}

		// Append the retrieved keys to our list
		keys = append(keys, results...)

		// Update the cursor for the next iteration
		cursor = nextCursor

		// If the cursor is 0, it means we have iterated through all keys
		if cursor == 0 {
			return keys, nil
		}
	}
}

package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache[T any] interface{}

type Cache[T any] struct {
	rdb        redis.UniversalClient
	transcoder Transcoder[T]
}

func NewRedisCache[T any](client redis.UniversalClient) *Cache[T] {
	return &Cache[T]{rdb: client, transcoder: NewPipelineTranscoder[T]()}
}

func NewRedisCacheWithTranscoder[T any](client redis.UniversalClient, transcoder Transcoder[T]) *Cache[T] {
	return &Cache[T]{rdb: client, transcoder: transcoder}
}

func (c *Cache[T]) Set(ctx context.Context, value T, key string, ttl time.Duration) error {
	var err error
	var str string

	str, err = c.transcoder.Encode(value)
	if err = c.rdb.Set(ctx, key, str, ttl).Err(); err != nil {
		return err
	}

	return nil
}

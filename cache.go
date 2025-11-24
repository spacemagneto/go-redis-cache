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
	if err != nil {
		return err
	}

	if err = c.rdb.Set(ctx, key, str, ttl).Err(); err != nil {
		return err
	}

	return nil
}

func (c *Cache[T]) Get(ctx context.Context, key string) (T, error) {
	var res T
	var err error
	var result string

	result, err = c.rdb.Get(ctx, key).Result()
	if err != nil {
		return res, err
	}

	res, err = c.transcoder.Decode(result)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (c *Cache[T]) GetWithTTL(ctx context.Context, key string) (T, time.Duration, error) {
	var res T
	var err error
	var ttl time.Duration
	var result string

	result, err = c.rdb.Get(ctx, key).Result()
	if err != nil {
		return res, 0, err
	}

	ttl, err = c.rdb.TTL(ctx, key).Result()
	if err != nil {
		return res, 0, err
	}

	res, err = c.transcoder.Decode(result)
	if err != nil {
		return res, 0, err
	}

	return res, ttl, nil
}

func (c *Cache[T]) Exists(ctx context.Context, key string) (bool, error) {
	var err error
	var exist int64
	exist, err = c.rdb.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}

	return exist == 1, err
}

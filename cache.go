package cache

import "github.com/redis/go-redis/v9"

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

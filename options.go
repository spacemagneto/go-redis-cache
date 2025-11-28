package cache

import (
	"time"

	"github.com/redis/go-redis/v9"
)

const defaultTTl = 2 * time.Hour

type options[T any] func(c *Cache[T])

func WithClient[T any](rdb redis.UniversalClient) options[T] {
	return func(c *Cache[T]) {
		c.rdb = rdb
	}
}

func WithTTL[T any](ttl time.Duration) options[T] {
	return func(c *Cache[T]) {
		if ttl == 0 {
			ttl = defaultTTl
		}

		c.ttl = ttl
	}
}

func WithTranscoder[T any](t Transcoder[T]) options[T] {
	return func(c *Cache[T]) {
		if t == nil {
			t = &defaultTranscoder[T]{}
		}

		c.transcoder = t
	}
}

package cache

import (
	"time"

	"github.com/redis/go-redis/v9"
)

// default TTL used when no explicit TTL is configured.
const defaultTTl = 2 * time.Hour

// options defines the functional options pattern for configuring Cache[T].
// Each option is a function that mutates the cache instance during construction.
type options[T any] func(c *Cache[T])

// WithClient sets the Redis client to be used by the cache.
// This option is mandatory - cache creation fails without a valid client.
func WithClient[T any](rdb redis.UniversalClient) options[T] {
	return func(c *Cache[T]) {
		c.rdb = rdb
	}
}

// WithTTL configures the default expiration time for cached values.
// If the provided duration is zero or negative, the global default (2 hours) is used.
// This value is only used when Set is called without an explicit TTL.
func WithTTL[T any](ttl time.Duration) options[T] {
	return func(c *Cache[T]) {
		if ttl <= 0 {
			ttl = defaultTTl
		}

		c.ttl = ttl
	}
}

// WithTranscoder allows injection of a custom transcoder for type T.
// If nil is passed, the default high-performance transcoder is used automatically.
// This enables custom serialization strategies (e.g. protobuf, msgpack, etc.).
func WithTranscoder[T any](t Transcoder[T]) options[T] {
	return func(c *Cache[T]) {
		if t == nil {
			t = &defaultTranscoder[T]{}
		}

		c.transcoder = t
	}
}

package cache

import (
	"context"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

// cache provides a generic Redis-backed caching layer for values of any type T.
// It combines a Redis client, a configurable transcoder for serialization, and a default TTL.
// The struct is intentionally lightweight and immutable after construction — all fields
// are set during creation and never changed afterward.
type cache[T any] struct {
	rdb        redis.UniversalClient
	transcoder Transcoder[T]
	ttl        time.Duration
}

// NewRedisCache constructs a fully configured cache[T] using the provided options.
// It applies all options in order, validates required fields, and sets sensible defaults.
// Returns an error only if no Redis client was provided.
func NewRedisCache[T any](opts ...options[T]) (*cache[T], error) {
	cache := &cache[T]{}

	for _, opt := range opts {
		opt(cache)
	}

	if cache.rdb == nil {
		return nil, ErrEmptyRedisClient
	}

	if cache.ttl == 0 {
		cache.ttl = defaultTTl
	}

	if cache.transcoder == nil {
		cache.transcoder = &defaultTranscoder[T]{}
	}

	return cache, nil
}

// Set stores the given value in Redis under the specified key.
// If ttl is zero, the cache's configured default TTL is used.
// The operation performs full validation and serialization before interacting with Redis.
// Returns nil on success or an appropriate error if any step fails.
func (c *cache[T]) Set(ctx context.Context, value T, key string, ttl time.Duration) error {
	if ttl == 0 {
		ttl = c.ttl
	}

	if IsNil[T](value) {
		return ErrValueIsEmpty
	}

	if strings.TrimSpace(key) == "" {
		return ErrKeyIsEmpty
	}

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

// Get retrieves a value from the cache using the provided key.
// The method first validates the key for meaningful content before interacting with Redis.
// It fetches the stored encoded data, handles any retrieval errors, and then delegates the decoding of the representation back into the original type using the transcoder.
// Method returns the reconstructed value and any error generated during the validation, retrieval, or decoding process, ensuring clear communication of the outcome to the caller.
func (c *cache[T]) Get(ctx context.Context, key string) (T, error) {
	var res T

	if strings.TrimSpace(key) == "" {
		return res, ErrKeyIsEmpty
	}

	var err error
	var result string

	result, err = c.rdb.Get(ctx, key).Result()
	if err != nil {
		return res, err
	}

	return c.transcoder.Decode(result)
}

// GetWithTTL method retrieves both the stored value and the remaining lifetime associated with the key in the cache.
// The method ensures that the key contains meaningful content before interacting with redis, then fetches the encoded value, checks for retrieval errors, retrieves the ttl, and decodes the stored representation into its original form.
// This method provides callers with both the reconstructed value and the duration until expiration, allowing them to understand not only the data but also its temporal validity within the cache.
// Errors from key validation, redis retrieval, ttl lookup, or decoding are all surfaced clearly to the caller to maintain predictable behavior.
func (c *cache[T]) GetWithTTL(ctx context.Context, key string) (T, time.Duration, error) {
	var res T

	if strings.TrimSpace(key) == "" {
		return res, 0, ErrKeyIsEmpty
	}

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

// Exists method checks whether a value associated with the provided key is present in the cache.
// The method validates that the key contains meaningful content before querying redis to avoid invalid lookups.
// It queries redis for the existence of the key and interprets the returned count to determine presence.
// Method returns a boolean indicating existence along with any error encountered during validation or redis interaction.
func (c *cache[T]) Exists(ctx context.Context, key string) (bool, error) {
	if strings.TrimSpace(key) == "" {
		return false, ErrKeyIsEmpty
	}

	var err error
	var exist int64

	exist, err = c.rdb.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}

	return exist == 1, err
}

// Delete method removes a value associated with the provided key from the cache.
// The method validates that the key contains meaningful content before interacting with redis to avoid invalid delete operations.
// It issues a delete command to redis and propagates any error that occurs during the removal process.
// This behavior ensures that cache state changes are performed safely and that failures are clearly reported to the caller.
func (c *cache[T]) Delete(ctx context.Context, key string) error {
	if strings.TrimSpace(key) == "" {
		return ErrKeyIsEmpty
	}

	if err := c.rdb.Del(ctx, key).Err(); err != nil {
		return err
	}

	return nil
}

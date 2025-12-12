package cache

import "errors"

var (
	// ErrEmptyRedisClient is returned when attempting to create a cache without providing a Redis client.
	// The Redis client is mandatory for all cache operations — construction fails if it is missing.
	ErrEmptyRedisClient = errors.New("redis client is empty")
	// ErrValueIsEmpty is returned by Set when the value being cached is nil.
	// Nil values cannot be meaningfully stored and would lead to ambiguous cache state.
	ErrValueIsEmpty = errors.New("value is empty")
	// ErrKeyIsEmpty is returned when an operation is attempted with an empty or whitespace-only key.
	// Cache keys must be non-empty and contain at least one non-whitespace character to be valid.
	ErrKeyIsEmpty = errors.New("cache key is empty or contains whitespace")
)

package cache

import "errors"

var (
	ErrEmptyRedisClient = errors.New("redis client is empty")
	ErrValueIsEmpty     = errors.New("value is empty")
	ErrKeyIsEmpty       = errors.New("cache key is empty or contains whitespace")
)

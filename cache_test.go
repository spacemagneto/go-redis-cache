package cache

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCache(t *testing.T) {
	t.Parallel()

	// Retrieve the singleton or global instance of TestContextInstance.
	// This instance represents the shared context object for testing purposes.
	testContext := TestContextInstance
	// Assert that the TestContextInstance is not nil.
	// This ensures that the instance has been correctly initialized and is available for use.
	// A nil value would indicate a failure in the setup or initialization process, which could
	// lead to runtime errors in operations relying on this context.
	assert.NotNil(t, testContext, "Expected TestContextInstance to be initialized and non-nil, but it was nil")

	// Create a new context for managing timeouts and cancellations in the Redis operations.
	// The `context.Background()` function returns an empty context that can be used as the root context.
	// It is typically used as the starting point for creating derived contexts or for situations
	// where no specific context with a deadline or cancellation is needed. This context is passed
	// to Redis operations to ensure that they can be properly managed or canceled if needed.
	ctx := testContext.context

	cache := NewRedisCache[user](testContext.GetRedis())

	t.Run("BaseRedisCache", func(t *testing.T) {
		baseCache := NewRedisCache[user](testContext.GetRedis())
		assert.NotNil(t, baseCache)
		assert.NotEmpty(t, baseCache.rdb)
		assert.NotEmpty(t, baseCache.transcoder)
	})

	t.Run("InitRedisCacheWithTranscoder", func(t *testing.T) {
		transcoder := NewPipelineTranscoder[user]()
		cacheWithTranscoder := NewRedisCacheWithTranscoder[user](testContext.GetRedis(), transcoder)
		assert.NotNil(t, cacheWithTranscoder)
		assert.NotEmpty(t, cacheWithTranscoder.rdb)
		assert.NotEmpty(t, cacheWithTranscoder.transcoder)
		assert.Equal(t, transcoder, cacheWithTranscoder.transcoder)
	})

	t.Run("SuccessSetInCache", func(t *testing.T) {
		payload := user{ID: 1, Name: "Name", Email: "email@gmail.com", Age: 11}
		key := fmt.Sprintf("cache_key_user_%d", payload.ID)
		err := cache.Set(ctx, payload, key, 10*time.Minute)
		assert.NoError(t, err)

		exist, err := cache.Exists(ctx, key)
		assert.NoError(t, err)
		assert.True(t, exist)
	})

	t.Run("SuccessSetAndGetData", func(t *testing.T) {
		expectPayload := user{ID: 2, Name: "Name2", Email: "email2@gmail.com", Age: 22}

		key := fmt.Sprintf("cache_key_user_%d", expectPayload.ID)
		err := cache.Set(ctx, expectPayload, key, 10*time.Minute)
		assert.NoError(t, err)

		payload, err := cache.Get(ctx, key)
		assert.NoError(t, err)
		assert.NotEmpty(t, payload)
		assert.Equal(t, expectPayload, payload)
	})

	t.Run("SuccessDelete", func(t *testing.T) {
		payload := user{ID: 3, Name: "Name3", Email: "email3@gmail.com", Age: 33}

		key := fmt.Sprintf("cache_key_user_%d", payload.ID)
		err := cache.Set(ctx, payload, key, 10*time.Minute)
		assert.NoError(t, err)

		deleteErr := cache.Delete(ctx, key)
		assert.NoError(t, deleteErr)

		exist, err := cache.Exists(ctx, key)
		assert.NoError(t, err)
		assert.False(t, exist)
	})
}

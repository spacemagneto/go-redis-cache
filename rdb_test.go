package cache

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type user struct {
	ID    int    `json:"id"`
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
	Age   int    `json:"age,omitempty"`
}

func TestCache(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	redisAddress := os.Getenv("REDIS_ADDRESS")

	rdb := redis.NewUniversalClient(&redis.UniversalOptions{Addrs: []string{redisAddress}, PoolSize: 10})

	pingErr := rdb.Ping(ctx).Err()
	assert.NoError(t, pingErr)

	cache, err := NewRedisCache[*user](WithClient[*user](rdb))
	assert.NoError(t, err)

	t.Run("InitRedisCacheWithTranscoder", func(t *testing.T) {
		transcoder := &defaultTranscoder[user]{}
		cacheWithTranscoder, err := NewRedisCache[user](WithClient[user](rdb), WithTranscoder[user](transcoder))
		assert.NoError(t, err)

		assert.NotNil(t, cacheWithTranscoder)
		assert.NotEmpty(t, cacheWithTranscoder.rdb)
		assert.NotNil(t, cacheWithTranscoder.transcoder)
		assert.Equal(t, transcoder, cacheWithTranscoder.transcoder)
		assert.Equal(t, defaultTTl, cacheWithTranscoder.ttl)
	})

	t.Run("NewRedisCacheWithoutOptions", func(t *testing.T) {
		_, err = NewRedisCache[user]()
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrEmptyRedisClient)
	})

	t.Run("NewRedisCacheWithTTLSetCustomTTL", func(t *testing.T) {
		c, err := NewRedisCache[user](
			WithClient[user](rdb),
			WithTranscoder[user](defaultTranscoder[user]{}),
			WithTTL[user](15*time.Minute),
		)

		assert.NoError(t, err)
		assert.Equal(t, 15*time.Minute, c.ttl)
	})

	t.Run("NewRedisCacheWithZeroTTL", func(t *testing.T) {
		c, err := NewRedisCache[user](
			WithClient[user](rdb),
			WithTranscoder[user](defaultTranscoder[user]{}),
			WithTTL[user](0),
		)

		assert.NoError(t, err)
		assert.Equal(t, defaultTTl, c.ttl)
	})

	t.Run("NewRedisCacheWithNegativeTTL", func(t *testing.T) {
		c, err := NewRedisCache[user](
			WithClient[user](rdb),
			WithTranscoder[user](defaultTranscoder[user]{}),
			WithTTL[user](-11),
		)

		assert.NoError(t, err)
		assert.Equal(t, defaultTTl, c.ttl)
	})

	t.Run("NewRedisCacheWithNilTranscoder", func(t *testing.T) {
		c, err := NewRedisCache[user](
			WithClient[user](rdb),
			WithTranscoder[user](nil),
			WithTTL[user](0),
		)

		assert.NoError(t, err)
		assert.Equal(t, defaultTTl, c.ttl)
		assert.Equal(t, &defaultTranscoder[user]{}, c.transcoder)
	})

	t.Run("SuccessSetInCache", func(t *testing.T) {
		payload := &user{ID: rand.Int(), Name: "Name", Email: "email@gmail.com", Age: rand.Int()}
		key := fmt.Sprintf("cache_key_user_%d", payload.ID)
		err = cache.Set(ctx, payload, key, 10*time.Minute)
		assert.NoError(t, err)

		exist, err := cache.Exists(ctx, key)
		assert.NoError(t, err)
		assert.True(t, exist)
	})

	t.Run("SetWithZeroTTL", func(t *testing.T) {
		payload := &user{ID: rand.Int(), Name: "Name", Email: "email@gmail.com", Age: rand.Int()}
		key := fmt.Sprintf("cache_key_user_%d", payload.ID)
		err = cache.Set(ctx, payload, key, 0)
		assert.NoError(t, err)

		exist, err := cache.Exists(ctx, key)
		assert.NoError(t, err)
		assert.True(t, exist)
	})

	t.Run("SetWithDoneContext", func(t *testing.T) {
		payload := &user{ID: rand.Int(), Name: "Name", Email: "email@gmail.com", Age: rand.Int()}
		key := fmt.Sprintf("cache_key_user_%d", payload.ID)

		doneCtx, cancel := context.WithCancel(ctx)
		cancel()

		err = cache.Set(doneCtx, payload, key, 10*time.Minute)
		assert.Error(t, err)
	})

	t.Run("SetWithNilValue", func(t *testing.T) {
		err = cache.Set(ctx, nil, "key", 10*time.Minute)
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrValueIsEmpty)
	})

	t.Run("SetWithEmptyKey", func(t *testing.T) {
		payload := &user{ID: rand.Int(), Name: "Name", Email: "email@gmail.com", Age: rand.Int()}
		err = cache.Set(ctx, payload, "      ", 10*time.Minute)
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrKeyIsEmpty)
	})

	t.Run("SuccessSetAndGetData", func(t *testing.T) {
		expectPayload := &user{ID: rand.Int(), Name: "Name2", Email: "email2@gmail.com", Age: rand.Int()}

		key := fmt.Sprintf("cache_key_user_%d", expectPayload.ID)
		err = cache.Set(ctx, expectPayload, key, 10*time.Minute)
		assert.NoError(t, err)

		payload, err := cache.Get(ctx, key)
		assert.NoError(t, err)
		assert.NotEmpty(t, payload)
		assert.Equal(t, expectPayload, payload)
	})

	t.Run("GetWithEmptyKey", func(t *testing.T) {
		_, err = cache.Get(ctx, "")
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrKeyIsEmpty)
	})

	t.Run("GetWithTTL", func(t *testing.T) {
		expectPayload := &user{ID: rand.Int(), Name: "Name", Email: "email@gmail.com", Age: rand.Int()}

		key := fmt.Sprintf("cache_key_user_%d", expectPayload.ID)
		err = cache.Set(ctx, expectPayload, key, 10*time.Minute)
		assert.NoError(t, err)

		payload, ttl, err := cache.GetWithTTL(ctx, key)
		assert.NoError(t, err)
		assert.NotEmpty(t, payload)
		assert.NotEqual(t, ttl, 0)
		assert.Equal(t, expectPayload, payload)
	})

	t.Run("GetWithDoneContext", func(t *testing.T) {
		doneCtx, cancel := context.WithCancel(ctx)
		cancel()

		_, err = cache.Get(doneCtx, "key")
		assert.Error(t, err)
	})

	t.Run("GetWithTTLDoneContext", func(t *testing.T) {
		doneCtx, cancel := context.WithCancel(ctx)
		cancel()

		_, _, err = cache.GetWithTTL(doneCtx, "key")
		assert.Error(t, err)
	})

	t.Run("GetWithTTLWithEmptyKey", func(t *testing.T) {
		_, _, err = cache.GetWithTTL(ctx, "")
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrKeyIsEmpty)
	})

	t.Run("ExistWithEmptyKey", func(t *testing.T) {
		_, err = cache.Exists(ctx, "")
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrKeyIsEmpty)
	})

	t.Run("ExistWithDoneContext", func(t *testing.T) {
		doneCtx, cancel := context.WithCancel(ctx)
		cancel()

		_, err = cache.Exists(doneCtx, "key")
		assert.Error(t, err)
	})

	t.Run("SuccessDelete", func(t *testing.T) {
		payload := &user{ID: rand.Int(), Name: "Name", Email: "email@gmail.com", Age: rand.Int()}

		key := fmt.Sprintf("cache_key_user_%d", payload.ID)
		err = cache.Set(ctx, payload, key, 10*time.Minute)
		assert.NoError(t, err)

		deleteErr := cache.Delete(ctx, key)
		assert.NoError(t, deleteErr)

		exist, err := cache.Exists(ctx, key)
		assert.NoError(t, err)
		assert.False(t, exist)
	})

	t.Run("DeleteWithEmptyKey", func(t *testing.T) {
		err = cache.Delete(ctx, "")
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrKeyIsEmpty)
	})

	t.Run("DeleteWithRandomKey", func(t *testing.T) {
		err = cache.Delete(ctx, fmt.Sprintf("rand_key_1"))
		assert.NoError(t, err)
	})

	t.Run("DeleteWithDoneContext", func(t *testing.T) {
		doneCtx, cancel := context.WithCancel(ctx)
		cancel()

		err = cache.Delete(doneCtx, "key")
		assert.Error(t, err)
	})
}

func TestCacheFailedEncode(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	redisAddress := os.Getenv("REDIS_ADDRESS")

	rdb := redis.NewUniversalClient(&redis.UniversalOptions{Addrs: []string{redisAddress}})

	pingErr := rdb.Ping(ctx).Err()
	assert.NoError(t, pingErr)

	mockTranscoder := NewMockTranscoder[*user](t)

	cache, err := NewRedisCache[*user](WithClient[*user](rdb), WithTranscoder[*user](mockTranscoder))
	assert.NoError(t, err)
	assert.NotNil(t, cache)

	mockTranscoder.EXPECT().Encode(mock.Anything).Return("", errors.New("failed encode data"))

	payload := &user{ID: 12414214, Name: "Name", Email: "email@gmail.com", Age: 876541}
	key := fmt.Sprintf("cache_key_user_%d", payload.ID)
	err = cache.Set(ctx, payload, key, 10*time.Minute)
	assert.Error(t, err)
}

func TestCacheFailedDecode(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	redisAddress := os.Getenv("REDIS_ADDRESS")

	rdb := redis.NewUniversalClient(&redis.UniversalOptions{Addrs: []string{redisAddress}})

	pingErr := rdb.Ping(ctx).Err()
	assert.NoError(t, pingErr)

	mockTranscoder := NewMockTranscoder[*user](t)

	cache, err := NewRedisCache[*user](WithClient[*user](rdb), WithTranscoder[*user](mockTranscoder))
	assert.NoError(t, err)
	assert.NotNil(t, cache)

	t.Run("GetDataWithFailedEncode", func(t *testing.T) {
		var zero *user

		expectEncode := "AAAAAAAAAAAAAAAAAAAAAA"

		mockTranscoder.EXPECT().Encode(mock.Anything).Return(expectEncode, nil)
		mockTranscoder.EXPECT().Decode(expectEncode).Return(zero, errors.New("failed decode"))

		expectPayload := &user{ID: 2, Name: "Name2", Email: "email2@gmail.com", Age: 22}

		key := fmt.Sprintf("cache_key_user_%d", expectPayload.ID)
		err = cache.Set(ctx, expectPayload, key, 10*time.Minute)
		assert.NoError(t, err)

		_, err = cache.Get(ctx, key)
		assert.Error(t, err)
	})

	t.Run("GetTTLWithFailedEncode", func(t *testing.T) {
		var zero *user

		expectEncode := "AAAAAAAAAAAAAAAAAAAAAA"

		mockTranscoder.EXPECT().Encode(mock.Anything).Return(expectEncode, nil)
		mockTranscoder.EXPECT().Decode(expectEncode).Return(zero, errors.New("failed decode"))

		expectPayload := &user{ID: rand.Int(), Name: "Name2", Email: "email2@gmail.com", Age: rand.Int()}

		key := fmt.Sprintf("cache_key_user_%d", expectPayload.ID)
		err = cache.Set(ctx, expectPayload, key, 10*time.Minute)
		assert.NoError(t, err)

		_, _, err = cache.GetWithTTL(ctx, key)
		assert.Error(t, err)
	})
}

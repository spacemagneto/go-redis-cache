# Redis Cache Wrapper
This library provides a generic, type-safe Redis-backed cache wrapper built on top of go-redis. It is designed to offer a minimal, predictable API with strong validation, configurable defaults, and pluggable serialization while remaining lightweight and easy to integrate into existing codebases.

The cache uses Go generics to allow storing and retrieving values of any type while maintaining compile-time type safety. All interactions with Redis are explicit, validated, and error-aware, ensuring that failures are never hidden from the caller.

## Installation

The library can be installed using the standard Go module tooling. It requires Go with generics support and a compatible version of go-redis.
```go
go get github.com/spacemagneto/go-redis-cache
```
After installation, the package can be imported directly into your application code and used alongside an existing go-redis client instance.

## Cache Behavior

### Set
The Set operation stores a value in Redis under the specified key. If no TTL is provided, the cache default TTL is applied automatically.
Before interacting with Redis, the method validates that the value is not empty and that the key contains meaningful content. The value is then encoded using the configured transcoder. Only after successful validation and encoding does the cache perform the Redis write.
Any error during validation, encoding, or Redis interaction is returned immediately.

### Get
The Get operation retrieves a value from Redis using the provided key. The key is validated before any Redis call is made.
The stored value is fetched in its encoded form and then decoded back into the original type using the configured transcoder. Errors from Redis or decoding are propagated directly to the caller.

### GetWithTTL
The GetWithTTL operation retrieves both the cached value and its remaining time to live. This method performs the same validation and decoding as Get, while also querying Redis for the key expiration duration.
This allows callers to reason about both the cached data and its temporal validity.

### Exists
The Exists operation checks whether a key is present in Redis. The key is validated before querying Redis.
The method returns a boolean indicating existence along with any error returned by Redis. This provides a lightweight way to check cache presence without retrieving or decoding the stored value.

### Delete
The Delete operation removes a key from Redis. The key is validated before the delete command is issued.
Any error returned by Redis during deletion is propagated to the caller, ensuring that cache state changes are always observable.

# Usage Examples

The cache is created once during application initialization and reused across the application lifetime. A Redis client must be provided, while other options are optional.

```go
ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
defer cancel()

rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
cache, err := NewRedisCache[string](
    WithClient[string](rdb),
    WithTTL[string](time.Hour),
)
if err != nil {
    log.Fatal(err)
}
```

Storing a value using the default TTL is done by passing zero as the ttl argument.

```go
err = cache.Set(ctx, "example-value", "example-key", 0)
if err != nil {
    log.Fatal(err)
}
```

Storing a value with an explicit TTL overrides the cache default for that operation only.

```go
err = cache.Set(ctx, "example-value", "example-key", 30*time.Minute)
if err != nil {
    log.Fatal(err)
}
```

Retrieving a value from the cache returns the decoded value or an error if retrieval or decoding fails.

```go
value, err := cache.Get(ctx, "example-key")
if err != nil {
    log.Fatal(err)
}
```

Retrieving a value together with its remaining TTL provides visibility into the expiration state of the entry.

```go
value, ttl, err := cache.GetWithTTL(ctx, "example-key")
if err != nil {
    log.Fatal(err)
}
```

Checking for key existence avoids decoding and is useful for conditional logic.

```go
exists, err := cache.Exists(ctx, "example-key")
    if err != nil {
    log.Fatal(err)
}
```

Deleting a cached value removes it from Redis immediately.


-------------------------------------------------
# License

This package is licensed under the Apache License, Version 2.0. See the LICENSE file for details.
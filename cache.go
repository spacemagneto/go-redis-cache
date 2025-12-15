package cache

import (
	"context"
	"time"
)

// Cache interface defines a generic contract for storing, retrieving, and managing values in a cache system.
// It provides methods for writing values, reading values, checking for existence, retrieving expiration information, and removing entries.
// The interface abstracts the underlying storage mechanism, allowing different cache implementations to be used interchangeably.
// This design enables consistent cache behavior while supporting multiple storage backends.
type Cache[T any] interface {
	// Set method stores a value in the cache using the provided key and lifetime duration.
	// Method validates the key and value, applies a default lifetime when no duration is provided, and persists the value in the underlying storage.
	// Any error produced during validation, encoding, or storage is returned to the caller for handling.
	// This method ensures that cache entries are written in a controlled and predictable manner.
	Set(ctx context.Context, value T, key string, ttl time.Duration) error

	// Get method retrieves a value from the cache using the provided key.
	// Method validates the key before attempting retrieval and returns the stored value after decoding it into its original form.
	// Any error produced during validation, retrieval, or decoding is returned to the caller.
	// This method provides a safe mechanism for accessing cached data.
	Get(ctx context.Context, key string) (T, error)

	// GetWithTTL method retrieves a value from the cache along with its remaining lifetime.
	// Method validates the key, fetches the stored data, retrieves the associated expiration duration, and decodes the value.
	// Both the decoded value and the ttl are returned together with any error encountered during the process.
	// This method allows callers to evaluate both the data and its expiration status.
	GetWithTTL(ctx context.Context, key string) (T, time.Duration, error)

	// Exists method checks whether a value associated with the provided key is present in the cache.
	// Method validates the key before querying the underlying storage system.
	// It returns a boolean indicating presence along with any error encountered during the existence check.
	// This method provides a lightweight way to determine cache membership.
	Exists(ctx context.Context, key string) (bool, error)

	// The Delete method removes a value associated with the provided key from the cache.
	// Method validates the key before issuing a delete command to the underlying storage system.
	// Any error encountered during deletion is returned to the caller.
	// This method allows callers to explicitly remove cache entries.
	Delete(ctx context.Context, key string) error
}

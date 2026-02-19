package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// Cache is a simple abstraction over Valkey/Redis used by repositories and
// services. Values are stored as JSON.
type Cache interface {
	// Get unmarshals the cached value at key into dest. If the key does not
	// exist, an error is returned.
	Get(ctx context.Context, key string, dest interface{}) error

	// Set marshals value to JSON and stores it at key with the given TTL.
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error

	// Delete removes the value at key (no-op if it does not exist).
	Delete(ctx context.Context, key string) error

	// Exists returns true if the key exists.
	Exists(ctx context.Context, key string) (bool, error)

	// Ping checks connectivity to the cache backend.
	Ping(ctx context.Context) error

	// GetClient returns the underlying Redis client for advanced operations (e.g., PubSub).
	GetClient() *redis.Client
}


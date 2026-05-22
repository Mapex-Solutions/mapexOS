package constants

import "time"

// Test Configuration
const (
	DefaultTimeout    = 30 * time.Second
	ShortTimeout      = 5 * time.Second
	AsyncWaitTime     = 100 * time.Millisecond // Wait for async operations
	CacheInvalidation = 100 * time.Millisecond // Wait for cache invalidation
)

// Test JWT Token (for development/testing only)
const (
	TestJWTSecret = "test-secret-key-for-e2e-tests"
)

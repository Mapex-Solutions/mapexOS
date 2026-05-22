package redis

// CacheKeyBuilderAdapter is the Redis-specific implementation of
// CacheKeyBuilderPort for the assettemplates module. It centralizes key
// formatting so the application layer never touches the raw Redis key
// format.
type CacheKeyBuilderAdapter struct{}

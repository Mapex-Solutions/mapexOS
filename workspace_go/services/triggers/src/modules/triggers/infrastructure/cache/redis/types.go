package redis

// KeyBuilder is the Redis adapter implementing ports.TriggerCacheKeyBuilderPort.
// It produces Redis cache keys using the prefixes defined in constants.go.
//
// This adapter is stateless — a single shared instance is safe for concurrent
// use and is registered once in the DI container (see triggers/module.go).
type KeyBuilder struct{}

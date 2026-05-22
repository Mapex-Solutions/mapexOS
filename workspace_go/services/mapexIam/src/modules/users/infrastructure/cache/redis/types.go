package redis

// CounterCacheAdapter is the infrastructure adapter that implements
// users/application/ports.CounterCachePort.
//
// Architecture Pattern: Hexagonal Architecture (Ports & Adapters)
//   - Port: ports.CounterCachePort
//   - Adapter: CounterCacheAdapter
//
// The adapter encapsulates the Redis key layout for the users counter cache,
// keeping infrastructure-specific knowledge (key prefix) out of the application
// layer.
type CounterCacheAdapter struct{}

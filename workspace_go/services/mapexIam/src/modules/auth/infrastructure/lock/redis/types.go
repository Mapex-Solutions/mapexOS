package redis

import (
	redisLock "github.com/Mapex-Solutions/mapexGoKit/infrastructure/redisLock"
)

// LockManagerAdapter is the infrastructure adapter that implements
// auth/application/ports.LockManagerPort by delegating to the
// redisLock.LockManager driver from mapexGoKit.
//
// Architecture Pattern: Hexagonal Architecture (Ports & Adapters)
//   - Port: ports.LockManagerPort
//   - Adapter: LockManagerAdapter (this struct)
//
// The adapter keeps the concrete *redisLock.LockManager driver out of the
// application layer's DI struct.
type LockManagerAdapter struct {
	lm *redisLock.LockManager
}

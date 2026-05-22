package redis

import (
	"context"
	"fmt"
	"time"

	"mapexIam/src/modules/auth/application/ports"

	redsync "github.com/go-redsync/redsync/v4"
	redisLock "github.com/Mapex-Solutions/mapexGoKit/infrastructure/redisLock"
)

// NewLockManagerAdapter builds a LockManagerAdapter around the given
// redisLock.LockManager driver. Returns the adapter as the application port
// interface so callers do not depend on the concrete driver type.
//
// Parameters:
//   - lm: The underlying redisLock.LockManager driver (provided by bootstrap).
//
// Returns:
//   - ports.LockManagerPort: The adapter as an interface.
func NewLockManagerAdapter(lm *redisLock.LockManager) ports.LockManagerPort {
	return &LockManagerAdapter{lm: lm}
}

// SetLock acquires a distributed lock with the given TTL. The returned value
// is an opaque handle (*redsync.Mutex underneath) and MUST be passed back to
// SetUnlock to release the lock.
func (a *LockManagerAdapter) SetLock(ctx context.Context, key string, ttl time.Duration) (any, error) {
	return a.lm.SetLock(ctx, key, ttl)
}

// SetUnlock releases the distributed lock referenced by the opaque mutex
// handle previously returned by SetLock.
func (a *LockManagerAdapter) SetUnlock(ctx context.Context, mutex any) error {
	m, ok := mutex.(*redsync.Mutex)
	if !ok {
		return fmt.Errorf("[INFRA:REDIS] LockManagerAdapter.SetUnlock: invalid mutex handle type %T", mutex)
	}
	return a.lm.SetUnlock(ctx, m)
}

// Compile-time check to ensure LockManagerAdapter implements LockManagerPort.
var _ ports.LockManagerPort = (*LockManagerAdapter)(nil)

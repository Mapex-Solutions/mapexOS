package ports

import (
	"context"
	"time"
)

// LockManagerPort is the driven port for distributed locking used by the
// authorization / coverage cache repositories.
//
// Architecture Pattern: Hexagonal Architecture (Ports & Adapters)
//   - Port: LockManagerPort (this interface)
//   - Adapter: infrastructure/lock/redis.LockManagerAdapter (wraps
//     github.com/Mapex-Solutions/mapexGoKit/infrastructure/redisLock.LockManager)
//
// The mutex handle returned by SetLock is treated as an opaque token (any)
// so the application/infrastructure boundary does not leak the concrete
// redsync.Mutex driver type. The same token MUST be passed back into
// SetUnlock by the caller.
//
// Only the methods actually consumed by the cache repositories are exposed.
type LockManagerPort interface {
	// SetLock attempts to acquire a distributed lock with the given TTL.
	// Returns an opaque mutex handle that MUST be passed to SetUnlock to
	// release the lock.
	SetLock(ctx context.Context, key string, ttl time.Duration) (any, error)

	// SetUnlock releases a distributed lock previously acquired via SetLock.
	// The mutex argument MUST be the exact token returned by SetLock.
	SetUnlock(ctx context.Context, mutex any) error
}

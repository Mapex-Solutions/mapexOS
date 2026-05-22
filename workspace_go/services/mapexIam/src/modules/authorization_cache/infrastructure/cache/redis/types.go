package redis

import (
	common "github.com/Mapex-Solutions/mapexGoKit/infrastructure/common/ports"
	"mapexIam/src/modules/authorization_cache/domain/repositories"
)

// AuthCacheRepository is the Redis adapter implementing the domain
// AuthCacheRepository port. It centralizes auth/coverage/role cache
// invalidation so memberships, roles, groups, and organizations all share
// the same versioning + delete semantics on Redis DB 5 (SharedCache).
//
// Cache key layout:
//   - User auth: auth:org:{orgId}:user:{userId}:ver (version pointer + versioned payload, 30d TTL)
//   - Coverage:  user:{userId}:orgs (simple DEL)
//   - Role:      role:{roleId} (simple DEL)
type AuthCacheRepository struct {
	cache common.SharedCache
}

// Compile-time port check.
var _ repositories.AuthCacheRepository = (*AuthCacheRepository)(nil)

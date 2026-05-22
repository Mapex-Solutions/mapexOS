package redis

import (
	"mapexIam/src/modules/auth/application/di"

	common "github.com/Mapex-Solutions/mapexGoKit/infrastructure/common/ports"
)

// AuthorizationCacheRepository implements the domain AuthorizationCacheRepository interface.
// This is the Adapter (infrastructure implementation) following Hexagonal Architecture.
//
// Architecture Pattern: Hexagonal Architecture (Ports & Adapters)
//   - Port: repositories.AuthorizationCacheRepository (domain interface)
//   - Adapter: AuthorizationCacheRepository (this implementation)
//
// This repository handles building authorization cache (user permissions per org).
// It implements the recursive inheritance algorithm and uses Redis distributed locks
// to prevent cache stampede.
//
// Responsibilities:
//   - Orchestrates other services to gather permission data
//   - Implements recursive inheritance algorithm
//   - Manages versioned cache storage
//   - Prevents cache stampede using distributed locks
//
// Cache Strategy:
//   - Version Pointer: auth:org:{orgId}:user:{userId}:ver → integer (no TTL, persists)
//   - Versioned Cache: auth:org:{orgId}:user:{userId}:v{N} → permissions array (TTL: 30 days)
type AuthorizationCacheRepository struct {
	deps di.AuthorizationCacheRepositoryDI
}

// CoverageCacheRepository implements the domain CoverageCacheRepository interface.
// This is the Adapter (infrastructure implementation) following Hexagonal Architecture.
//
// Architecture Pattern: Hexagonal Architecture (Ports & Adapters)
//   - Port: repositories.CoverageCacheRepository (domain interface)
//   - Adapter: CoverageCacheRepository (this implementation)
//
// This repository handles building coverage cache with hierarchy expansion.
// It uses cursor pagination to fetch ALL user memberships without hard limits.
//
// Responsibilities:
//   - Fetches ALL user memberships using cursor pagination
//   - Expands hierarchy for recursive memberships (scope="recursive")
//   - Builds flat list of accessible orgIds for queries
//   - Fetches organization metadata
//   - Manages cache storage with TTL
//   - Prevents cache stampede using distributed locks
//
// Cache Strategy:
//   - Key format: coverage:user:{userId}
//   - TTL: 30 days
//   - Contains: UserAccess (accessibleOrgIds + detailed organizations)
type CoverageCacheRepository struct {
	deps di.CoverageCacheRepositoryDI
}

// SessionRepository implements the domain SessionRepository interface using Redis cache.
// It stores user session refresh tokens with TTL for authentication token rotation.
//
// Architecture Pattern: Hexagonal Architecture (Ports & Adapters)
//   - This is the Adapter (infrastructure implementation)
//   - Implements: repositories.SessionRepository (domain port)
//
// Cache Strategy:
//   - Key pattern: USER_SESSION:REFRESH:{userId}:{sessionId}
//   - Storage: Redis (via common.Cache interface)
//   - TTL: Configurable per session (typically 7 days)
//
// Security:
//   - Sessions are isolated by userId + sessionId combination
//   - Deletion prevents token reuse (logout/revocation)
type SessionRepository struct {
	cache common.Cache
}

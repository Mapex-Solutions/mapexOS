package repositories

import (
	"context"
)

// CoverageCacheRepository defines the contract for building coverage cache.
// This follows Hexagonal Architecture - Port pattern.
//
// Architecture Pattern: Repository Pattern
//   - Port: CoverageCacheRepository (domain interface)
//   - Adapter: Implementation in infrastructure/cache
//
// This repository is responsible for:
//   - Building coverage cache (list of organizations user has access to WITH hierarchy expansion)
//   - Fetching ALL user memberships using cursor pagination
//   - Expanding hierarchy for recursive memberships (scope="recursive")
//   - Extracting flat list of accessible orgIds for queries
//   - Fetching organization metadata
//   - Managing cache storage with TTL
//
// Cache Strategy:
//   - Key format: coverage:user:{userId}
//   - TTL: 30 days
//   - Contains: UserAccess (accessibleOrgIds + detailed organizations)
//   - Uses distributed locks to prevent cache stampede
//
// Usage:
//   - Called by middleware on cache miss
//   - Never called eagerly (lazy build strategy)
type CoverageCacheRepository interface {
	// BuildCache builds the coverage cache for a user (list of accessible organizations).
	// Uses cursor pagination to fetch ALL user memberships without hard limits.
	// Uses Redis distributed lock to prevent cache stampede.
	//
	// Algorithm:
	//  1. Acquire Redis lock (prevents multiple builds)
	//  2. Use GetAllUserMemberships() to fetch ALL memberships (cursor pagination)
	//  3. For each membership with scope="recursive", expand hierarchy (get all children)
	//  4. Build flat list of accessible orgIds (deduplicated)
	//  5. Fetch organization metadata (id, name, type, pathKey, scope, membershipId, roleIds)
	//  6. Save to cache with 30-day TTL
	//  7. Release lock
	//
	// If lock cannot be acquired:
	//  - Polls every 500ms for 10s waiting for cache
	//  - Returns cached UserAccess when ready
	//  - Timeout error if cache not ready after 10s
	//
	// Parameters:
	//   - ctx: Request-scoped context for cancellation and timeouts
	//   - userId: Unique identifier of the user
	//
	// Returns:
	//   - *UserAccess: Complete access information (accessibleOrgIds + organizations)
	//   - error: If build fails or timeout occurs
	BuildCache(ctx context.Context, userId string) (*UserAccess, error)

	// GetCachedAccess retrieves the cached coverage for a user.
	// Returns cache miss error if not found.
	//
	// Parameters:
	//   - ctx: Request-scoped context
	//   - userId: Unique identifier of the user
	//
	// Returns:
	//   - *UserAccess: Complete access information
	//   - error: If cache not found or retrieval fails
	GetCachedAccess(ctx context.Context, userId string) (*UserAccess, error)

	// InvalidateCache removes the coverage cache for a user.
	// Called when user memberships change (create/update/delete).
	//
	// Parameters:
	//   - ctx: Request-scoped context
	//   - userId: Unique identifier of the user
	//
	// Returns:
	//   - error: If invalidation fails
	InvalidateCache(ctx context.Context, userId string) error
}

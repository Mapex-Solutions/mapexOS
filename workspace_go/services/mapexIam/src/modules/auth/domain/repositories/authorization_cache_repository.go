package repositories

import (
	"context"
)

// AuthorizationCacheRepository defines the contract for building authorization cache.
// This follows Hexagonal Architecture - Port pattern.
//
// Architecture Pattern: Repository Pattern
//   - Port: AuthorizationCacheRepository (domain interface)
//   - Adapter: Implementation in infrastructure/cache
//
// This repository is responsible for:
//   - Building authorization cache (computing permissions from memberships, roles, orgs)
//   - Using recursive inheritance algorithm
//   - Managing versioned cache storage
//   - Coordinating with other services to gather permission data
//
// Cache Strategy:
//   - Version Pointer: auth:org:{orgId}:user:{userId}:ver → integer (no TTL, persists)
//   - Versioned Cache: auth:org:{orgId}:user:{userId}:v{N} → permissions array (TTL: 30 days)
//   - Uses distributed locks to prevent cache stampede
//
// Usage:
//   - Called by internal API endpoints on cache miss
//   - Never called eagerly (lazy build strategy)
type AuthorizationCacheRepository interface {
	// BuildCache builds the authorization cache for a user in a specific organization.
	// Implements recursive inheritance algorithm with Redis distributed lock to prevent stampede.
	//
	// Algorithm:
	//  1. Acquire Redis lock (prevents multiple builds)
	//  2. Get organization (pathKey, rolePolicy)
	//  3. Get local memberships (direct in org)
	//  4. Get recursive memberships (inherited from parents) if rolePolicy="merge"
	//  5. Resolve all roleIds → permissions
	//  6. Remove duplicates
	//  7. Get next version number (round-robin 1-100)
	//  8. Save to cache with version
	//  9. Release lock
	//
	// If lock cannot be acquired:
	//  - Polls every 500ms for 10s waiting for cache
	//  - Returns cached permissions when ready
	//  - Timeout error if cache not ready after 10s
	//
	// Parameters:
	//   - ctx: Request-scoped context for cancellation and timeouts
	//   - userId: Unique identifier of the user
	//   - orgId: Unique identifier of the organization
	//
	// Returns:
	//   - []string: Array of permission strings
	//   - int: Cache version number
	//   - error: If build fails or timeout occurs
	BuildCache(ctx context.Context, userId, orgId string) ([]string, int, error)

	// GetOrBuildCache attempts to read permissions from cache (O(1) hot path).
	// Falls back to BuildCache on cache miss (cold path).
	//
	// This method is designed for the public endpoint where:
	//   - Most requests hit cache (O(1) Redis GET)
	//   - Cold start triggers full build (rare)
	//
	// Parameters:
	//   - ctx: Request-scoped context for cancellation and timeouts
	//   - userId: Unique identifier of the user (from JWT)
	//   - orgId: Unique identifier of the organization (from X-Org-Context header)
	//
	// Returns:
	//   - []string: Array of permission strings
	//   - int: Cache version number
	//   - error: If both cache read and build fail
	GetOrBuildCache(ctx context.Context, userId, orgId string) ([]string, int, error)
}

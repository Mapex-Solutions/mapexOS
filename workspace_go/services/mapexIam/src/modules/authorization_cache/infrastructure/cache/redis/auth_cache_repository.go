package redis

import (
	"context"
	"fmt"
	"strconv"

	common "github.com/Mapex-Solutions/mapexGoKit/infrastructure/common/ports"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
	"mapexIam/src/modules/authorization_cache/domain/repositories"
)

// New creates a new AuthCacheRepository instance.
//
// Parameters:
//   - cache: SharedCache implementation (typically Redis DB 5 for cross-service auth data)
//
// Returns:
//   - repositories.AuthCacheRepository: The repository implementation as interface
func New(cache common.SharedCache) repositories.AuthCacheRepository {
	return &AuthCacheRepository{cache: cache}
}

// InvalidateUserAuth invalidates the authorization cache for a specific user in a specific organization
// using the versioning strategy.
//
// Cache Strategy (Versioning):
//   - Version Pointer (persistent): auth:org:{orgId}:user:{userId}:ver → integer (1-100)
//   - Versioned Cache (30 days TTL): auth:org:{orgId}:user:{userId}:v{version} → permissions data
//   - On invalidation: Bump version (old cache expires naturally via TTL)
//   - Round robin prevents collisions (1-100)
//
// Benefits:
//   - O(1) invalidation per user (just increment version)
//   - No need to delete thousands of cache keys
//   - Old cache expires gracefully via TTL
//   - Memory efficient (version keys are 1-3 bytes)
//
// Implementation:
//   - Get current version from cache (default 0 if not exists)
//   - Calculate new version ((currentVer % 100) + 1)
//   - Set new version in cache (no TTL - persists forever)
//   - Old versioned cache expires naturally via its TTL
//
// Parameters:
//   - ctx: Request-scoped context for cancellation and timeouts
//   - userId: Unique identifier of the user
//   - orgId: Unique identifier of the organization
//
// Returns:
//   - error: Error if invalidation fails
func (r *AuthCacheRepository) InvalidateUserAuth(ctx context.Context, userId string, orgId string) error {
	// Build version pointer key: auth:org:{orgId}:user:{userId}:ver
	verKey := fmt.Sprintf("auth:org:%s:user:%s:ver", orgId, userId)

	// Get current version (default to 0 if not exists)
	var currentVerStr string
	currentVer := 0
	err := r.cache.Get(ctx, verKey, &currentVerStr)
	if err == nil && currentVerStr != "" {
		currentVer, _ = strconv.Atoi(currentVerStr)
	}

	// Calculate new version (round robin 1-100)
	newVer := (currentVer % 100) + 1

	// Set NEW version (no TTL on version pointer - persists forever)
	// Old versioned cache will expire naturally via its TTL (30 minutes)
	if err := r.cache.Set(ctx, verKey, fmt.Sprintf("%d", newVer)); err != nil {
		logger.Error(err, fmt.Sprintf("[REPO:AuthCache] Failed to set new auth version for user=%s org=%s", userId, orgId))
		return err
	}

	logger.Info(fmt.Sprintf("[REPO:AuthCache] Invalidated auth cache for user=%s org=%s (version: %d -> %d)", userId, orgId, currentVer, newVer))
	return nil
}

// InvalidateCoverage removes the coverage cache for a specific user.
// Coverage refers to the list of organizations that a user has access to through their memberships.
//
// Cache Strategy:
//   - Coverage cache key pattern: user:{userId}:orgs
//   - TTL: 30 minutes
//   - Contains: List of organizations with metadata (id, name, type, pathKey)
//   - Compression: GZIP (500 bytes → 150 bytes)
//
// This forces a rebuild of the coverage cache on the next request to /users/me/coverage.
//
// Implementation:
//   - Builds cache key: user:{userId}:orgs
//   - Executes DEL command on Redis
//   - Logs operation for audit trail
//
// Parameters:
//   - ctx: Request-scoped context for cancellation and timeouts
//   - userId: Unique identifier of the user
//
// Returns:
//   - error: Error if invalidation fails
func (r *AuthCacheRepository) InvalidateCoverage(ctx context.Context, userId string) error {
	// Build cache key following the pattern: user:{userId}:orgs
	cacheKey := fmt.Sprintf("user:%s:orgs", userId)

	// Delete from cache (simple DEL - no versioning needed)
	if err := r.cache.Del(ctx, cacheKey); err != nil {
		logger.Error(err, fmt.Sprintf("[REPO:AuthCache] Failed to invalidate coverage cache for userId=%s", userId))
		return err
	}

	logger.Info(fmt.Sprintf("[REPO:AuthCache] Invalidated coverage cache for userId=%s", userId))
	return nil
}

// InvalidateRole removes a role from the cache.
// When a role's permissions change, this ensures that all users with this role
// will get fresh permission data on the next authorization check.
//
// Cache Strategy:
//   - Roles are cached separately with key pattern: role:{roleId}
//   - TTL: 1 hour (roles change rarely)
//   - When role is updated/deleted: DEL role:{roleId}
//   - This is O(1) invalidation - no need to touch user authorization caches
//   - Multiple users benefit from the same role cache (shared cache)
//
// Implementation:
//   - Builds cache key: role:{roleId}
//   - Executes DEL command on Redis
//   - Logs operation for audit trail
//
// Parameters:
//   - ctx: Request-scoped context for cancellation and timeouts
//   - roleId: Unique identifier of the role to invalidate
//
// Returns:
//   - error: Error if invalidation fails
func (r *AuthCacheRepository) InvalidateRole(ctx context.Context, roleId string) error {
	// Build cache key following the pattern: role:{roleId}
	cacheKey := fmt.Sprintf("role:%s", roleId)

	// Delete from cache
	if err := r.cache.Del(ctx, cacheKey); err != nil {
		logger.Error(err, fmt.Sprintf("[REPO:AuthCache] Failed to invalidate role cache for roleId=%s", roleId))
		return err
	}

	logger.Info(fmt.Sprintf("[REPO:AuthCache] Invalidated role cache for roleId=%s", roleId))
	return nil
}

// Compile-time check to ensure AuthCacheRepository implements the domain interface
var _ repositories.AuthCacheRepository = (*AuthCacheRepository)(nil)

package redis

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"mapexIam/src/modules/auth/domain/repositories"
	"mapexIam/src/modules/auth/application/di"

	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// NewAuthorizationCacheRepository creates a new AuthorizationCacheRepository instance.
// Uses DI struct pattern for clean dependency management.
//
// Parameters:
//   - deps: Dependency injection container with all required dependencies
//
// Returns:
//   - repositories.AuthorizationCacheRepository: The repository implementation as interface
func NewAuthorizationCacheRepository(deps di.AuthorizationCacheRepositoryDI) repositories.AuthorizationCacheRepository {
	return &AuthorizationCacheRepository{
		deps: deps,
	}
}

// BuildCache builds the authorization cache for a user in a specific organization.
// Implements recursive inheritance algorithm with Redis distributed lock to prevent stampede.
//
// Algorithm:
//  1. Try acquire Redis lock (TTL: 10s)
//  2. If lock acquired:
//     a. Get organization (pathKey, rolePolicy)
//     b. Get local memberships (direct in org)
//     c. Get recursive memberships (inherited from parents) if rolePolicy="merge"
//     d. Resolve all roleIds → permissions
//     e. Remove duplicates
//     f. Get next version number (round-robin 1-100)
//     g. Save to cache with version
//     h. Release lock
//  3. If lock failed:
//     a. Poll every 500ms for 10s waiting for cache
//     b. Return cached permissions when ready
//     c. Timeout error if cache not ready after 10s
func (r *AuthorizationCacheRepository) BuildCache(ctx context.Context, userId, orgId string) ([]string, int, error) {
	normalizedOrgId := normalizeOrgId(orgId)
	lockKey := fmt.Sprintf("build:lock:auth:org:%s:user:%s", normalizedOrgId, userId)
	lockTTL := 10 * time.Second

	logger.Info(fmt.Sprintf("[REPO:AuthCache] Attempting to build cache for user=%s org=%s", userId, orgId))

	// Try to acquire lock
	mutex, err := r.deps.LockMgr.SetLock(ctx, lockKey, lockTTL)

	if err == nil {
		// Lock acquired → This instance builds cache
		defer r.deps.LockMgr.SetUnlock(ctx, mutex)

		logger.Info(fmt.Sprintf("[REPO:AuthCache] Lock acquired, building cache for user=%s org=%s", userId, orgId))

		permissions, version, buildErr := r.buildPermissions(ctx, userId, orgId)
		if buildErr != nil {
			logger.Error(buildErr, fmt.Sprintf("[REPO:AuthCache] Failed to build cache for user=%s org=%s", userId, orgId))
			return nil, 0, buildErr
		}

		logger.Info(fmt.Sprintf("[REPO:AuthCache] Built cache for user=%s org=%s version=%d (%d permissions)",
			userId, orgId, version, len(permissions)))
		return permissions, version, nil
	}

	// Lock not acquired → Wait for cache to be ready
	logger.Info(fmt.Sprintf("[REPO:AuthCache] Lock not acquired, waiting for cache to be ready user=%s org=%s", userId, orgId))

	maxWait := 10 * time.Second
	pollInterval := 500 * time.Millisecond
	startTime := time.Now()

	for time.Since(startTime) < maxWait {
		// Check if cache is ready
		permissions, version, cacheErr := r.getCachedPermissions(ctx, userId, orgId)
		if cacheErr == nil && len(permissions) > 0 {
			logger.Info(fmt.Sprintf("[REPO:AuthCache] Cache ready for user=%s org=%s version=%d", userId, orgId, version))
			return permissions, version, nil
		}

		// Wait before retry
		time.Sleep(pollInterval)
	}

	logger.Error(nil, fmt.Sprintf("[REPO:AuthCache] Timeout waiting for cache user=%s org=%s", userId, orgId))
	return nil, 0, errors.New("timeout waiting for cache build")
}

// GetOrBuildCache attempts to read permissions from cache (O(1) hot path).
// Falls back to BuildCache on cache miss (cold path).
//
// This method is designed for the public GET /auth/me/permissions endpoint
// where most requests hit cache. Only on miss does it trigger a full build.
func (r *AuthorizationCacheRepository) GetOrBuildCache(ctx context.Context, userId, orgId string) ([]string, int, error) {
	// 1. Try cache first (O(1) — hot path)
	permissions, version, err := r.getCachedPermissions(ctx, userId, orgId)
	if err == nil && len(permissions) > 0 {
		logger.Info(fmt.Sprintf("[REPO:AuthCache] Cache hit for user=%s org=%s version=%d", userId, orgId, version))
		return permissions, version, nil
	}

	// 2. Cache miss — full build (cold path)
	logger.Info(fmt.Sprintf("[REPO:AuthCache] Cache miss for user=%s org=%s, triggering build", userId, orgId))
	return r.BuildCache(ctx, userId, orgId)
}

// buildPermissions implements the recursive inheritance algorithm to build permissions
func (r *AuthorizationCacheRepository) buildPermissions(ctx context.Context, userId, orgId string) ([]string, int, error) {
	allRoleIds := make(map[string]bool) // Use map to avoid duplicates

	// Get ALL user memberships (direct + via groups) using GetAllUserMemberships
	// This is the single source of truth for user access, including group memberships
	allMemberships, err := r.deps.MembershipService.GetAllUserMemberships(ctx, userId)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get all memberships: %w", err)
	}

	logger.Info(fmt.Sprintf("[REPO:AuthCache] Found %d total memberships (direct + groups) for user=%s", len(allMemberships), userId))

	// When orgId is empty, build ROOT permissions (mapex.*)
	if orgId == "" {
		logger.Info(fmt.Sprintf("[REPO:AuthCache] Building ROOT permissions for user=%s", userId))

		// Extract all roleIds from all memberships (process in memory)
		for _, membership := range allMemberships {
			if membership.Enabled {
				for _, roleId := range membership.RoleIds {
					allRoleIds[roleId.Hex()] = true
				}
			}
		}
	} else {
		// Standard flow: build permissions for specific org

		// 1. Get organization (pathKey, rolePolicy)
		org, err := r.deps.OrgService.GetOrganizationById(ctx, &orgId)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to get organization: %w", err)
		}

		// 2. Filter memberships for this specific org (local memberships)
		localCount := 0
		for _, membership := range allMemberships {
			if membership.Enabled && membership.OrgID != nil && membership.OrgID.Hex() == orgId {
				for _, roleId := range membership.RoleIds {
					allRoleIds[roleId.Hex()] = true
				}
				localCount++
			}
		}

		logger.Info(fmt.Sprintf("[REPO:AuthCache] Found %d local memberships for user=%s org=%s", localCount, userId, orgId))

		// 3. Get recursive memberships (inherited from parents) if rolePolicy="merge"
		// Default rolePolicy is "merge" for backward compatibility
		rolePolicy := "merge"
		if org.AccessPolicy.RolePolicy != "" {
			rolePolicy = org.AccessPolicy.RolePolicy
		}

		if rolePolicy == "merge" {
			// Filter by pathKey prefix in memory for recursive memberships
			matchedCount := 0
			for _, membership := range allMemberships {
				if membership.Enabled && membership.Scope == "recursive" {
					if org.PathKey != nil && *org.PathKey != "" && membership.OrgPathKey != "" && isPrefix(membership.OrgPathKey, *org.PathKey) {
						for _, roleId := range membership.RoleIds {
							allRoleIds[roleId.Hex()] = true
						}
						matchedCount++
					}
				}
			}
			logger.Info(fmt.Sprintf("[REPO:AuthCache] Found %d recursive memberships (matched by pathKey) for user=%s org=%s", matchedCount, userId, orgId))
		}
	}

	// 4. Resolve all roleIds → permissions
	allPermissions := []string{}
	for roleId := range allRoleIds {
		rolePermissions, err := r.getRolePermissions(ctx, roleId)
		if err != nil {
			logger.Error(err, fmt.Sprintf("[REPO:AuthCache] Failed to get role permissions for roleId=%s", roleId))
			continue
		}
		allPermissions = append(allPermissions, rolePermissions...)
	}

	// 5. Remove duplicates
	uniquePermissions := removeDuplicates(allPermissions)

	// 6. Get next version number
	version, err := r.getNextVersion(ctx, userId, orgId)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get next version: %w", err)
	}

	// 7. Save to cache (versioned)
	normalizedOrgId := normalizeOrgId(orgId)
	cacheKey := fmt.Sprintf("auth:org:%s:user:%s:v%d", normalizedOrgId, userId, version)
	ttl := 30 * 24 * time.Hour // 30 days

	if err := r.deps.Cache.SetEx(ctx, cacheKey, uniquePermissions, ttl); err != nil {
		return nil, 0, fmt.Errorf("failed to cache permissions: %w", err)
	}

	return uniquePermissions, version, nil
}

// getCachedPermissions retrieves permissions from cache using current version
func (r *AuthorizationCacheRepository) getCachedPermissions(ctx context.Context, userId, orgId string) ([]string, int, error) {
	normalizedOrgId := normalizeOrgId(orgId)
	verKey := fmt.Sprintf("auth:org:%s:user:%s:ver", normalizedOrgId, userId)

	var versionStr string
	if err := r.deps.Cache.Get(ctx, verKey, &versionStr); err != nil {
		return nil, 0, fmt.Errorf("version not found: %w", err)
	}

	version, err := strconv.Atoi(versionStr)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid version format: %w", err)
	}

	cacheKey := fmt.Sprintf("auth:org:%s:user:%s:v%d", normalizedOrgId, userId, version)
	var permissions []string
	if err := r.deps.Cache.Get(ctx, cacheKey, &permissions); err != nil {
		return nil, 0, fmt.Errorf("cache not found: %w", err)
	}

	return permissions, version, nil
}

// getNextVersion gets the next version number and cleans up old versions
// Strategy:
//   - Current version (nextVer): No TTL, always available
//   - Previous version (currentVer): TTL 10s - active clients can finish
//   - Older versions (penultimate and before): DELETE immediately
func (r *AuthorizationCacheRepository) getNextVersion(ctx context.Context, userId, orgId string) (int, error) {
	normalizedOrgId := normalizeOrgId(orgId)
	verKey := fmt.Sprintf("auth:org:%s:user:%s:ver", normalizedOrgId, userId)

	var currentVerStr string
	currentVer := 0
	err := r.deps.Cache.Get(ctx, verKey, &currentVerStr)
	if err == nil && currentVerStr != "" {
		currentVer, _ = strconv.Atoi(currentVerStr)
	}

	nextVer := (currentVer % 100) + 1

	// Set next version as current
	if err := r.deps.Cache.Set(ctx, verKey, fmt.Sprintf("%d", nextVer)); err != nil {
		return 0, err
	}

	// Clean up old versions asynchronously to avoid blocking the response
	go r.cleanupOldVersions(context.Background(), userId, normalizedOrgId, currentVer, nextVer)

	return nextVer, nil
}

// cleanupOldVersions removes old cache versions to prevent memory bloat
// - Previous version (currentVer): Set TTL 10s for active clients
// - Penultimate and older: Delete immediately
func (r *AuthorizationCacheRepository) cleanupOldVersions(ctx context.Context, userId, normalizedOrgId string, currentVer, nextVer int) {
	// Set TTL 10s on previous version (allow active clients to finish)
	// We need to re-fetch the value and set it with TTL since there's no EXPIRE command in the interface
	if currentVer > 0 {
		prevKey := fmt.Sprintf("auth:org:%s:user:%s:v%d", normalizedOrgId, userId, currentVer)

		// Get current value
		var prevPermissions []string
		if err := r.deps.Cache.Get(ctx, prevKey, &prevPermissions); err == nil {
			// Re-set with 10s TTL
			if err := r.deps.Cache.SetEx(ctx, prevKey, prevPermissions, 10*time.Second); err != nil {
				logger.Error(err, fmt.Sprintf("[REPO:AuthCache] Failed to set TTL on previous version user=%s org=%s v%d", userId, normalizedOrgId, currentVer))
			}
		}
	}

	// Delete penultimate and older versions (v0 to currentVer-1, excluding currentVer)
	// We keep: nextVer (current, no TTL) and currentVer (previous, TTL 10s)
	for v := 1; v <= 100; v++ {
		// Skip current and next version
		if v == currentVer || v == nextVer {
			continue
		}

		oldKey := fmt.Sprintf("auth:org:%s:user:%s:v%d", normalizedOrgId, userId, v)
		if err := r.deps.Cache.Del(ctx, oldKey); err != nil {
			// Log but don't fail - cleanup is best-effort
			logger.Error(err, fmt.Sprintf("[REPO:AuthCache] Failed to delete old version user=%s org=%s v%d", userId, normalizedOrgId, v))
		}
	}

	logger.Info(fmt.Sprintf("[REPO:AuthCache] Cleaned old versions for user=%s org=%s (kept v%d with TTL 10s, v%d current)", userId, normalizedOrgId, currentVer, nextVer))
}

// getRolePermissions uses RoleService to fetch permissions for a given role
func (r *AuthorizationCacheRepository) getRolePermissions(ctx context.Context, roleId string) ([]string, error) {
	roleResponse, err := r.deps.RoleService.GetRoleById(ctx, &roleId)
	if err != nil {
		return nil, fmt.Errorf("failed to find role: %w", err)
	}

	if roleResponse.Permissions == nil {
		return []string{}, nil
	}

	return *roleResponse.Permissions, nil
}

// Helper functions

// normalizeOrgId converts empty orgId to "global" for better cache key identification.
// This prevents cache keys like "auth:org::user:" and makes them "auth:org:global:user:"
// which is more readable and indicates these are global/hierarchical permissions (mapex.*, admin_vendor.*, etc.)
//
// Parameters:
//   - orgId: Organization ID (can be empty)
//
// Returns:
//   - string: "global" if empty, otherwise returns orgId unchanged
func normalizeOrgId(orgId string) string {
	if orgId == "" {
		return "global"
	}
	return orgId
}

// isPrefix checks if prefix is a prefix of path for hierarchical pathKey validation.
// Used to determine if a membership's pathKey is an ancestor of the current organization.
//
// Parameters:
//   - prefix: The prefix to check (e.g., "0001/0002/")
//   - path: The full path to validate (e.g., "0001/0002/0003/")
//
// Returns:
//   - bool: true if prefix is a prefix of path, false otherwise
func isPrefix(prefix, path string) bool {
	if len(prefix) > len(path) {
		return false
	}
	return path[:len(prefix)] == prefix
}

// removeDuplicates removes duplicate strings from a slice while preserving order.
// Used to deduplicate permissions from multiple roles.
//
// Parameters:
//   - slice: Slice of strings that may contain duplicates
//
// Returns:
//   - []string: New slice with duplicates removed
func removeDuplicates(slice []string) []string {
	keys := make(map[string]bool)
	result := []string{}
	for _, entry := range slice {
		if _, exists := keys[entry]; !exists {
			keys[entry] = true
			result = append(result, entry)
		}
	}
	return result
}

// Compile-time check to ensure AuthorizationCacheRepository implements the domain interface
var _ repositories.AuthorizationCacheRepository = (*AuthorizationCacheRepository)(nil)

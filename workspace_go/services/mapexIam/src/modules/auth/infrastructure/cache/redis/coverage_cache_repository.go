package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"mapexIam/src/modules/auth/domain/repositories"
	"mapexIam/src/modules/auth/application/di"
	membershipPorts "mapexIam/src/modules/memberships/application/ports"
	orgPorts "mapexIam/src/modules/organizations/application/ports"

	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// NewCoverageCacheRepository creates a new CoverageCacheRepository instance.
// Uses DI struct pattern for clean dependency management.
//
// Parameters:
//   - deps: Dependency injection container with all required dependencies
//
// Returns:
//   - repositories.CoverageCacheRepository: The repository implementation as interface
func NewCoverageCacheRepository(deps di.CoverageCacheRepositoryDI) repositories.CoverageCacheRepository {
	return &CoverageCacheRepository{
		deps: deps,
	}
}

// BuildCache builds the coverage cache for a user (list of accessible organizations).
// Uses cursor pagination to fetch ALL user memberships without hard limits.
// Uses Redis distributed lock to prevent cache stampede.
//
// Algorithm:
//  1. Try acquire Redis lock (TTL: 10s)
//  2. If lock acquired:
//     a. Use GetAllUserMemberships() to fetch ALL memberships (cursor pagination)
//     b. For each membership with scope="recursive", expand hierarchy (get all children)
//     c. Build flat list of accessible orgIds (deduplicated)
//     d. Build detailed organization coverage list
//     e. Save to cache with 30-day TTL
//     f. Release lock
//  3. If lock failed:
//     a. Poll every 500ms for 10s waiting for cache
//     b. Return cached UserAccess when ready
//     c. Timeout error if cache not ready after 10s
func (r *CoverageCacheRepository) BuildCache(ctx context.Context, userId string) (*repositories.UserAccess, error) {
	lockKey := fmt.Sprintf("build:lock:coverage:user:%s", userId)
	lockTTL := 10 * time.Second

	logger.Info(fmt.Sprintf("[REPO:CoverageCache] Attempting to build cache for user=%s", userId))

	// Try to acquire lock
	mutex, err := r.deps.LockMgr.SetLock(ctx, lockKey, lockTTL)

	if err == nil {
		// Lock acquired → This instance builds cache
		defer r.deps.LockMgr.SetUnlock(ctx, mutex)

		logger.Info(fmt.Sprintf("[REPO:CoverageCache] Lock acquired, building cache for user=%s", userId))

		userAccess, buildErr := r.buildCoverage(ctx, userId)
		if buildErr != nil {
			logger.Error(buildErr, fmt.Sprintf("[REPO:CoverageCache] Failed to build cache for user=%s", userId))
			return nil, buildErr
		}

		logger.Info(fmt.Sprintf("[REPO:CoverageCache] Built cache for user=%s (%d accessible orgs)", userId, len(userAccess.AccessibleOrgIds)))
		return userAccess, nil
	}

	// Lock not acquired → Wait for cache to be ready
	logger.Info(fmt.Sprintf("[REPO:CoverageCache] Lock not acquired, waiting for cache to be ready user=%s", userId))

	maxWait := 10 * time.Second
	pollInterval := 500 * time.Millisecond
	startTime := time.Now()

	for time.Since(startTime) < maxWait {
		// Check if cache is ready
		userAccess, cacheErr := r.GetCachedAccess(ctx, userId)
		if cacheErr == nil && len(userAccess.AccessibleOrgIds) >= 0 {
			logger.Info(fmt.Sprintf("[REPO:CoverageCache] Cache ready for user=%s", userId))
			return userAccess, nil
		}

		// Wait before retry
		time.Sleep(pollInterval)
	}

	logger.Error(nil, fmt.Sprintf("[REPO:CoverageCache] Timeout waiting for cache user=%s", userId))
	return nil, errors.New("timeout waiting for cache build")
}

// GetCachedAccess retrieves the cached coverage for a user.
// Returns cache miss error if not found.
func (r *CoverageCacheRepository) GetCachedAccess(ctx context.Context, userId string) (*repositories.UserAccess, error) {
	cacheKey := fmt.Sprintf("coverage:user:%s", userId)

	var userAccess repositories.UserAccess
	if err := r.deps.Cache.Get(ctx, cacheKey, &userAccess); err != nil {
		return nil, fmt.Errorf("cache not found: %w", err)
	}

	return &userAccess, nil
}

// InvalidateCache removes the coverage cache for a user.
// Called when user memberships change (create/update/delete).
func (r *CoverageCacheRepository) InvalidateCache(ctx context.Context, userId string) error {
	cacheKey := fmt.Sprintf("coverage:user:%s", userId)

	logger.Info(fmt.Sprintf("[REPO:CoverageCache] Invalidating cache for user=%s", userId))

	if err := r.deps.Cache.Del(ctx, cacheKey); err != nil {
		logger.Error(err, fmt.Sprintf("[REPO:CoverageCache] Failed to invalidate cache for user=%s", userId))
		return err
	}

	logger.Info(fmt.Sprintf("[REPO:CoverageCache] Invalidated cache for user=%s", userId))
	return nil
}

// buildCoverage implements the coverage building logic with hierarchy expansion
func (r *CoverageCacheRepository) buildCoverage(ctx context.Context, userId string) (*repositories.UserAccess, error) {
	// 1. Use GetAllUserMemberships to fetch ALL memberships (cursor pagination, no hard limit)
	memberships, err := r.deps.MembershipService.GetAllUserMemberships(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to get user memberships: %w", err)
	}

	logger.Info(fmt.Sprintf("[REPO:CoverageCache] Found %d memberships for user=%s", len(memberships), userId))

	// 2. If no memberships found, return empty access
	if len(memberships) == 0 {
		logger.Info(fmt.Sprintf("[REPO:CoverageCache] User %s has no memberships", userId))

		emptyAccess := &repositories.UserAccess{
			UserID:           userId,
			AccessibleOrgIds: []string{},
			Organizations:    []repositories.OrganizationCoverage{},
			LastUpdated:      time.Now(),
			Version:          1,
		}

		// Cache empty result
		if err := r.saveToCache(ctx, userId, emptyAccess); err != nil {
			return nil, err
		}

		return emptyAccess, nil
	}

	// 3. Expand hierarchy and build coverage
	accessibleOrgIds, organizations, err := r.expandHierarchyAndBuildCoverage(ctx, memberships)
	if err != nil {
		return nil, fmt.Errorf("failed to expand hierarchy: %w", err)
	}

	// 4. Build UserAccess
	userAccess := &repositories.UserAccess{
		UserID:           userId,
		AccessibleOrgIds: accessibleOrgIds,
		Organizations:    organizations,
		LastUpdated:      time.Now(),
		Version:          1,
	}

	// 5. Save to cache
	if err := r.saveToCache(ctx, userId, userAccess); err != nil {
		return nil, err
	}

	return userAccess, nil
}

// expandHierarchyAndBuildCoverage expands hierarchy for recursive memberships and builds coverage list
func (r *CoverageCacheRepository) expandHierarchyAndBuildCoverage(
	ctx context.Context,
	memberships []*membershipPorts.Membership,
) ([]string, []repositories.OrganizationCoverage, error) {

	// Maps for deduplication
	orgIdsMap := make(map[string]bool)
	coverageMap := make(map[string]*repositories.OrganizationCoverage)

	for _, membership := range memberships {
		if membership == nil {
			continue
		}

		orgID := membership.OrgID.Hex()
		scope := membership.Scope
		membershipID := membership.ID.Hex()

		// Get roleIds
		roleIDs := make([]string, 0, len(membership.RoleIds))
		for _, roleID := range membership.RoleIds {
			roleIDs = append(roleIDs, roleID.Hex())
		}

		// Get org details
		org, err := r.deps.OrgService.GetOrganizationById(ctx, &orgID)
		if err != nil {
			logger.Error(err, fmt.Sprintf("[REPO:CoverageCache] Failed to get organization orgId=%s", orgID))
			continue
		}

		// Determine rolePolicy (default "merge" for backward compatibility)
		rolePolicy := "merge"
		if org.AccessPolicy.RolePolicy != "" {
			rolePolicy = org.AccessPolicy.RolePolicy
		}

		// Add direct organization
		orgIdsMap[orgID] = true
		coverageMap[orgID] = &repositories.OrganizationCoverage{
			ID:           org.ID.Hex(),
			Name:         *org.Name,
			Type:         *org.Type,
			PathKey:      getStringValue(org.PathKey),
			Scope:        scope,
			MembershipID: membershipID,
			RoleIDs:      roleIDs,
			RolePolicy:   rolePolicy,
		}

		// If scope is "recursive", expand hierarchy (get all children)
		if scope == "recursive" && org.PathKey != nil {
			children, err := r.getChildOrganizations(ctx, *org.PathKey)
			if err != nil {
				logger.Error(err, fmt.Sprintf("[REPO:CoverageCache] Failed to get children for pathKey=%s", *org.PathKey))
				continue
			}

			// Add children to accessible list (filter by rolePolicy)
			addedCount := 0
			for _, child := range children {
				// Determine child's rolePolicy (default "merge")
				childRolePolicy := "merge"
				if child.AccessPolicy.RolePolicy != "" {
					childRolePolicy = child.AccessPolicy.RolePolicy
				}

				// Only add child if it accepts inherited permissions (rolePolicy="merge")
				if childRolePolicy != "merge" {
					logger.Info(fmt.Sprintf("[REPO:CoverageCache] Skipping child orgId=%s (rolePolicy=%s, does not accept inheritance)", child.ID.Hex(), childRolePolicy))
					continue
				}

				childID := child.ID.Hex()
				orgIdsMap[childID] = true
				addedCount++

				// Only add to coverageMap if not already there (direct membership takes precedence)
				if _, exists := coverageMap[childID]; !exists {
					coverageMap[childID] = &repositories.OrganizationCoverage{
						ID:           childID,
						Name:         child.Name,
						Type:         child.Type,
						PathKey:      child.PathKey,
						Scope:        "inherited", // Mark as inherited from parent
						MembershipID: membershipID,
						RoleIDs:      roleIDs,
						RolePolicy:   childRolePolicy,
					}
				}
			}

			logger.Info(fmt.Sprintf("[REPO:CoverageCache] Added %d/%d children (filtered by rolePolicy=merge) for pathKey=%s", addedCount, len(children), *org.PathKey))
		}
	}

	// Convert maps to slices
	accessibleOrgIds := make([]string, 0, len(orgIdsMap))
	for id := range orgIdsMap {
		accessibleOrgIds = append(accessibleOrgIds, id)
	}

	organizations := make([]repositories.OrganizationCoverage, 0, len(coverageMap))
	for _, coverage := range coverageMap {
		organizations = append(organizations, *coverage)
	}

	return accessibleOrgIds, organizations, nil
}

// getChildOrganizations retrieves all child organizations for a given pathKey using prefix match.
// Returns organization entities with accessPolicy included (for rolePolicy filtering).
func (r *CoverageCacheRepository) getChildOrganizations(ctx context.Context, parentPathKey string) ([]orgPorts.Organization, error) {
	// Use OrganizationService to get child organizations
	// This follows Hexagonal Architecture: Infrastructure uses Service Port (cross-domain)
	// Note: GetChildOrganizationsByPathKey now includes accessPolicy in projection
	children, err := r.deps.OrgService.GetChildOrganizationsByPathKey(ctx, parentPathKey)
	if err != nil {
		logger.Error(err, fmt.Sprintf("[REPO:CoverageCache] Failed to get children for pathKey=%s", parentPathKey))
		return nil, err
	}

	return children, nil
}

// saveToCache saves UserAccess to Redis cache
func (r *CoverageCacheRepository) saveToCache(ctx context.Context, userId string, userAccess *repositories.UserAccess) error {
	cacheKey := fmt.Sprintf("coverage:user:%s", userId)
	ttl := 30 * 24 * time.Hour // 30 days

	if err := r.deps.Cache.SetEx(ctx, cacheKey, userAccess, ttl); err != nil {
		return fmt.Errorf("failed to cache coverage: %w", err)
	}

	return nil
}

// getStringValue safely extracts string value from pointer.
// Returns empty string if pointer is nil.
//
// Parameters:
//   - ptr: Pointer to string (can be nil)
//
// Returns:
//   - string: Value if pointer is not nil, empty string otherwise
func getStringValue(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}

// Compile-time check to ensure CoverageCacheRepository implements the domain interface
var _ repositories.CoverageCacheRepository = (*CoverageCacheRepository)(nil)

package repositories

import (
	"context"
)

// AuthCacheRepository defines the contract for managing authorization cache invalidation.
// This follows Hexagonal Architecture - Port pattern.
//
// This is a SHARED repository used by multiple modules:
//   - MembershipService (for invalidating user auth when memberships change)
//   - RoleService (for invalidating auth when roles change)
//   - GroupService (for invalidating auth when users are added to groups)
//   - OrganizationService (for invalidating coverage when orgs change)
//
// Architecture Note:
//   - Normally, each module manages its own cache
//   - This is an EXCEPTION because authorization cache is shared across multiple domains
//   - Centralizes cache invalidation to avoid code duplication
//   - Follows the same pattern as auth/SessionRepository
//
// Cache Strategy:
//   - User Auth: Versioning strategy (auth:org:{orgId}:user:{userId}:ver)
//   - Coverage: Simple deletion (user:{userId}:orgs)
//   - Role: Simple deletion (role:{roleId})
type AuthCacheRepository interface {
	// InvalidateUserAuth invalidates the authorization cache for a user in an organization.
	// This should be called when:
	//   - A membership is created/updated/deleted
	//   - A user is added/removed from a group
	//   - A role's permissions are updated
	//
	// Uses versioning strategy:
	//   - Bumps version pointer (1-100 round robin)
	//   - Old cached data expires naturally via TTL
	//   - O(1) invalidation per user
	//
	// Parameters:
	//   - ctx: Request-scoped context for cancellation and timeouts
	//   - userId: Unique identifier of the user
	//   - orgId: Unique identifier of the organization
	//
	// Returns:
	//   - error: Error if invalidation fails
	InvalidateUserAuth(ctx context.Context, userId string, orgId string) error

	// InvalidateCoverage invalidates the coverage cache for a user.
	// Coverage refers to the list of organizations that a user has access to.
	//
	// This should be called when:
	//   - A membership is created (user gains access to new org)
	//   - A membership is deleted (user loses access to org)
	//   - A user is added/removed from a group
	//   - An organization is updated (metadata changed)
	//
	// Cache key: user:{userId}:orgs
	//
	// Parameters:
	//   - ctx: Request-scoped context for cancellation and timeouts
	//   - userId: Unique identifier of the user
	//
	// Returns:
	//   - error: Error if invalidation fails
	InvalidateCoverage(ctx context.Context, userId string) error

	// InvalidateRole removes a role from the cache.
	// When a role's permissions change, this ensures users with this role
	// will get fresh permission data on next authorization check.
	//
	// This should be called when:
	//   - A role's permissions are updated
	//   - A role is deleted
	//
	// Cache key: role:{roleId}
	//
	// Parameters:
	//   - ctx: Request-scoped context for cancellation and timeouts
	//   - roleId: Unique identifier of the role
	//
	// Returns:
	//   - error: Error if invalidation fails
	InvalidateRole(ctx context.Context, roleId string) error
}

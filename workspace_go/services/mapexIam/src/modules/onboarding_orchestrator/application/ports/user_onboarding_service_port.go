package ports

import (
	"context"
	"mapexIam/src/modules/onboarding_orchestrator/application/dtos"

	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
)

// UserOnboardingServicePort defines the contract for user onboarding operations.
// This interface follows Hexagonal Architecture (Ports & Adapters) pattern,
// enabling dependency inversion and decoupling from concrete implementations.
//
// Architecture Pattern: Hexagonal Architecture (Ports & Adapters)
//   - Port: UserOnboardingServicePort (this interface)
//   - Adapter: UserOnboardingService (implements the contract)
//
// The port allows for:
//   - Easy mocking in tests
//   - Swappable implementations
//   - Clear separation between domain logic and infrastructure concerns
//
// Used by:
//   - HTTP Handlers (for HTTP interface layer)
type UserOnboardingServicePort interface {
	// CreateUserWithMemberships creates a user and associates them with multiple organizations/roles atomically.
	// This method orchestrates the complete user onboarding flow:
	//  1. Gets OrgID from RequestContext and Scope from org.AccessPolicy.DefaultScope
	//  2. Creates the user
	//  3. Creates membership for the user (direct roles or group)
	//  4. Returns combined response with user and memberships
	//
	// Parameters:
	//   - ctx: The context for controlling cancellation and timeouts
	//   - requestContext: Contains OrgContext (selected orgId) for determining membership org and scope
	//   - dto: A pointer to CreateUserWithMembershipsDto containing user data and membership data
	//
	// Returns:
	//   - A pointer to UserOnboardingResponse with the created user and memberships
	//   - An error if any step fails
	CreateUserWithMemberships(ctx context.Context, requestContext *reqCtx.RequestContext, dto *dtos.CreateUserWithMembershipsDto) (*dtos.UserOnboardingResponse, error)

	// UpdateUserWithAccess updates user data and replaces their access configuration atomically.
	// This method orchestrates the complete user update flow:
	//  1. Gets OrgID from RequestContext and Scope from org.AccessPolicy.DefaultScope
	//  2. Updates user data (only provided fields)
	//  3. Removes existing memberships/groups for user in current org
	//  4. Creates new membership (direct roles) OR adds to group
	//  5. Returns combined response with updated user and new memberships
	//
	// All operations run in a MongoDB transaction for ACID compliance.
	//
	// Parameters:
	//   - ctx: The context for controlling cancellation and timeouts
	//   - requestContext: Contains OrgContext (selected orgId) for determining membership org and scope
	//   - userID: The ID of the user to update
	//   - dto: A pointer to UpdateUserWithAccessDto containing user data updates and access config
	//
	// Returns:
	//   - A pointer to UserOnboardingResponse with the updated user and new memberships
	//   - An error if any step fails (transaction will be rolled back)
	UpdateUserWithAccess(ctx context.Context, requestContext *reqCtx.RequestContext, userID string, dto *dtos.UpdateUserWithAccessDto) (*dtos.UserOnboardingResponse, error)
}

package ports

import (
	"context"

	"mapexIam/src/modules/auth/application/dtos"
)

// AuthServicePort defines the inbound port (Hexagonal Architecture) for authentication operations.
// This port allows external layers (handlers, other services) to depend on authentication operations
// without coupling to the concrete AuthService implementation.
//
// Architecture Pattern: Hexagonal Architecture (Ports & Adapters)
//   - Port: This interface (defines the contract)
//   - Adapter: AuthService (implements the contract)
//
// Used by:
//   - HTTP Handlers (for API endpoints)
//   - Other application services requiring authentication operations
type AuthServicePort interface {
	// Login authenticates a user by email and password.
	// Issues both an access token (short-lived) and a refresh token (longer-lived).
	//
	// Parameters:
	//   - ctx: Request-scoped context
	//   - dto: Login credentials (email and password)
	//
	// Returns:
	//   - interface{}: Map containing "access_token", "refresh_token", and "user" data
	//   - error: Error if authentication fails
	Login(ctx context.Context, dto *dtos.LoginDTO) (interface{}, error)

	// RefreshToken validates a refresh token and issues new tokens.
	// Implements token rotation for security.
	//
	// Parameters:
	//   - ctx: Request-scoped context
	//   - refreshToken: The refresh token to validate
	//
	// Returns:
	//   - map[string]interface{}: Map containing "access_token", "refresh_token", and "user" data
	//   - error: Error if token is invalid or expired
	RefreshToken(ctx context.Context, refreshToken string) (map[string]interface{}, error)

	// Logout invalidates a refresh token, ending the user's session.
	//
	// Parameters:
	//   - ctx: Request-scoped context
	//   - token: The access token to extract session info from
	//
	// Returns:
	//   - error: Error if logout fails
	Logout(ctx context.Context, token string) error

	// GetMyCoverage returns the list of organizations accessible by a user.
	// Used by UI to populate organization selectors and navigation tree.
	// Tries cache first, builds coverage if cache miss.
	//
	// Parameters:
	//   - ctx: Request-scoped context
	//   - userId: The user ID to fetch coverage for
	//
	// Returns:
	//   - map[string]interface{}: Map containing "organizations" and "lastUpdated"
	//   - error: Error if coverage fetch fails
	GetMyCoverage(ctx context.Context, userId string) (map[string]interface{}, error)

	// GetMyPermissions returns the resolved permissions for a user in a specific organization.
	// Reads from Redis cache (O(1)), triggers build on cache miss.
	// Used by the UI to control visibility of UI elements and route access.
	//
	// Parameters:
	//   - ctx: Request-scoped context
	//   - userId: The user ID from JWT token
	//   - orgId: The organization ID from X-Org-Context header
	//
	// Returns:
	//   - map[string]interface{}: Map containing "permissions" ([]string) and "version" (int)
	//   - error: Error if permission fetch fails
	GetMyPermissions(ctx context.Context, userId, orgId string) (map[string]interface{}, error)
}

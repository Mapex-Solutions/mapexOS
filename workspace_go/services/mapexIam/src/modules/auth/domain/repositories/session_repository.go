package repositories

import (
	"context"
	"time"
)

// SessionRepository defines the contract for managing user authentication sessions.
// Sessions are stored with refresh tokens for authentication token rotation and revocation.
//
// Architecture Pattern: Hexagonal Architecture (Ports & Adapters)
//   - Port: This interface (outbound port - defined by domain)
//   - Adapter: Infrastructure implementation (e.g., Redis cache)
//
// Domain Concept:
//   - A session represents an authenticated user's refresh token lifecycle
//   - Sessions are identified by userId + sessionId (jti from JWT)
//   - Sessions expire after a configurable TTL
//
// Used by:
//   - AuthService (for login, token refresh, and logout operations)
type SessionRepository interface {
	// StoreRefreshToken persists a refresh token for a user session.
	// This should be called when:
	//   - A user logs in (new session)
	//   - A refresh token is rotated (token refresh)
	//
	// Parameters:
	//   - ctx: Request-scoped context for cancellation and timeouts
	//   - userId: Unique identifier of the user
	//   - sessionId: Session identifier (typically JWT jti claim)
	//   - refreshToken: The refresh token to store
	//   - ttl: Time-to-live for the session
	//
	// Returns:
	//   - error: Error if storage fails
	StoreRefreshToken(
		ctx context.Context,
		userId string,
		sessionId string,
		refreshToken string,
		ttl time.Duration,
	) error

	// GetRefreshToken retrieves a stored refresh token for validation.
	// This should be called when:
	//   - Validating a refresh token request
	//
	// Parameters:
	//   - ctx: Request-scoped context for cancellation and timeouts
	//   - userId: Unique identifier of the user
	//   - sessionId: Session identifier (typically JWT jti claim)
	//
	// Returns:
	//   - string: The stored refresh token
	//   - error: Error if token not found or retrieval fails
	GetRefreshToken(
		ctx context.Context,
		userId string,
		sessionId string,
	) (string, error)

	// InvalidateRefreshToken removes a refresh token, ending the session.
	// This should be called when:
	//   - A user logs out
	//   - A session needs to be terminated (security breach, forced logout)
	//
	// Parameters:
	//   - ctx: Request-scoped context for cancellation and timeouts
	//   - userId: Unique identifier of the user
	//   - sessionId: Session identifier (typically JWT jti claim)
	//
	// Returns:
	//   - error: Error if invalidation fails
	//
	// Security Note:
	//   - Once invalidated, the refresh token cannot be used to issue new access tokens
	//   - Access tokens already issued remain valid until they expire
	InvalidateRefreshToken(
		ctx context.Context,
		userId string,
		sessionId string,
	) error
}

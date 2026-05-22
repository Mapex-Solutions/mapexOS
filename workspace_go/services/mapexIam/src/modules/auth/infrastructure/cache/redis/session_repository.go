package redis

import (
	"context"
	"fmt"
	"time"

	common "github.com/Mapex-Solutions/mapexGoKit/infrastructure/common/ports"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
	"mapexIam/src/modules/auth/domain/repositories"
)

// New creates a new SessionRepository instance.
//
// Parameters:
//   - cache: Cache implementation (typically Redis)
//
// Returns:
//   - repositories.SessionRepository: The repository implementation
func New(cache common.Cache) repositories.SessionRepository {
	return &SessionRepository{
		cache: cache,
	}
}

// StoreRefreshToken persists a refresh token in cache with TTL.
//
// Implementation:
//   - Builds cache key from userId and sessionId
//   - Stores token with configurable TTL
//   - Logs operation for audit trail
//
// Cache Key: USER_SESSION:REFRESH:{userId}:{sessionId}
func (r *SessionRepository) StoreRefreshToken(
	ctx context.Context,
	userId string,
	sessionId string,
	refreshToken string,
	ttl time.Duration,
) error {
	cacheKey := fmt.Sprintf("USER_SESSION:REFRESH:%s:%s", userId, sessionId)

	if err := r.cache.SetEx(ctx, cacheKey, refreshToken, ttl); err != nil {
		logger.Error(err, fmt.Sprintf("[REPO:Session] Failed to store refresh token for userId=%s sessionId=%s", userId, sessionId))
		return err
	}

	logger.Info(fmt.Sprintf("[REPO:Session] Stored refresh token for userId=%s sessionId=%s", userId, sessionId))
	return nil
}

// GetRefreshToken retrieves a stored refresh token.
//
// Implementation:
//   - Builds cache key from userId and sessionId
//   - Retrieves token from cache
//   - Returns error if not found or cache failure
//
// Cache Key: USER_SESSION:REFRESH:{userId}:{sessionId}
func (r *SessionRepository) GetRefreshToken(
	ctx context.Context,
	userId string,
	sessionId string,
) (string, error) {
	cacheKey := fmt.Sprintf("USER_SESSION:REFRESH:%s:%s", userId, sessionId)
	var cachedToken string

	if err := r.cache.Get(ctx, cacheKey, &cachedToken); err != nil {
		logger.Error(err, fmt.Sprintf("[REPO:Session] Failed to get refresh token for userId=%s sessionId=%s", userId, sessionId))
		return "", err
	}

	return cachedToken, nil
}

// InvalidateRefreshToken removes a refresh token from cache.
//
// Implementation:
//   - Builds cache key from userId and sessionId
//   - Deletes the key from cache
//   - Logs operation for audit trail
//
// Security:
//   - Once deleted, the refresh token cannot be used
//   - Access tokens remain valid until they expire
//
// Cache Key: USER_SESSION:REFRESH:{userId}:{sessionId}
func (r *SessionRepository) InvalidateRefreshToken(
	ctx context.Context,
	userId string,
	sessionId string,
) error {
	cacheKey := fmt.Sprintf("USER_SESSION:REFRESH:%s:%s", userId, sessionId)

	if err := r.cache.Del(ctx, cacheKey); err != nil {
		logger.Error(err, fmt.Sprintf("[REPO:Session] Failed to invalidate refresh token for userId=%s sessionId=%s", userId, sessionId))
		return err
	}

	logger.Info(fmt.Sprintf("[REPO:Session] Invalidated refresh token for userId=%s sessionId=%s", userId, sessionId))
	return nil
}

// Compile-time check to ensure SessionRepository implements the domain interface
var _ repositories.SessionRepository = (*SessionRepository)(nil)

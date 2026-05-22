package services

import (
	"context"

	userDtos "mapexIam/src/modules/users/application/dtos"

	customErrors "github.com/Mapex-Solutions/mapexGoKit/microservices/http/customErrors"
	httpStatus "github.com/Mapex-Solutions/mapexGoKit/microservices/http/status"
	utilsJWT "github.com/Mapex-Solutions/mapexGoKit/utils/jwt"
)

// validateRefreshClaims parses the refresh JWT and extracts (userId,
// sessionId). Any parse / signature failure surfaces as 401.
func (s *AuthService) validateRefreshClaims(refreshToken string) (string, string, error) {
	claims, err := utilsJWT.ParseRefreshToken(s.di.AuthConfig.Secret, refreshToken)
	if err != nil {
		return "", "", &customErrors.ServerCustomError{
			Code:   httpStatus.UNAUTHORIZED,
			Errors: []string{"Invalid or expired refresh token."},
		}
	}
	return claims.Subject, claims.ID, nil
}

// assertCachedRefreshMatches confirms the presented token is the one
// currently registered for the (userId, sessionId) pair. Cache miss or
// mismatch yields 401 so a stolen token from a previous rotation cannot
// be replayed.
func (s *AuthService) assertCachedRefreshMatches(ctx context.Context, userId, sessionId, refreshToken string) error {
	cachedToken, err := s.di.SessionRepo.GetRefreshToken(ctx, userId, sessionId)
	if err != nil || cachedToken != refreshToken {
		return &customErrors.ServerCustomError{
			Code:   httpStatus.UNAUTHORIZED,
			Errors: []string{"Refresh token not found or mismatch."},
		}
	}
	return nil
}

// loadEnabledUserById reloads the user and confirms the account is
// still enabled. Soft-deleted or disabled accounts cannot rotate tokens
// even when their refresh token is still in cache.
func (s *AuthService) loadEnabledUserById(ctx context.Context, userId string) (*userDtos.UserResponse, error) {
	user, err := s.di.UserService.GetUserById(ctx, &userId)
	if user == nil || err != nil || !*user.Enabled {
		return nil, &customErrors.ServerCustomError{
			Code:   httpStatus.UNAUTHORIZED,
			Errors: []string{"User not found or disabled."},
		}
	}
	return user, nil
}

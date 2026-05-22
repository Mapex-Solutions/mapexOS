package services

import (
	"context"

	"mapexIam/src/modules/auth/application/constants"
	"mapexIam/src/modules/auth/application/di"
	"mapexIam/src/modules/auth/application/dtos"
	"mapexIam/src/modules/auth/application/ports"
	userDtos "mapexIam/src/modules/users/application/dtos"
	userPorts "mapexIam/src/modules/users/application/ports"

	customErrors "github.com/Mapex-Solutions/mapexGoKit/microservices/http/customErrors"
	httpStatus "github.com/Mapex-Solutions/mapexGoKit/microservices/http/status"
	utilsJWT "github.com/Mapex-Solutions/mapexGoKit/utils/jwt"
	mapper "github.com/Mapex-Solutions/mapexGoKit/utils/mapper"
	random "github.com/Mapex-Solutions/mapexGoKit/utils/random"
)

// Compile-time check to ensure AuthService implements AuthServicePort
var _ ports.AuthServicePort = (*AuthService)(nil)

// New creates a new AuthService.
func New(deps di.AuthServiceDI) ports.AuthServicePort {
	return &AuthService{di: deps}
}

// Login orchestrates the credential-based authentication flow:
// load user by email -> reject when missing or disabled -> verify password
// -> mint a session id and sign access + refresh JWTs -> persist the
// refresh token in cache for revocation/rotation -> return the tokens +
// the user response DTO.
func (s *AuthService) Login(ctx context.Context, dto *dtos.LoginDTO) (interface{}, error) {
	user, err := s.loadUserForLogin(ctx, dto.Email)
	if err != nil {
		return nil, err
	}
	if err := s.verifyLoginCredentials(user, dto.Password); err != nil {
		return nil, err
	}

	sessionId, _ := random.GenerateSessionID(4)
	accessToken, refreshToken := s.signSessionTokens(user.ID.Hex(), user.Email, sessionId)

	s.di.SessionRepo.StoreRefreshToken(ctx, user.ID.Hex(), sessionId, refreshToken, constants.RefreshTokenTTL)
	userResponse, _ := mapper.EntityToDto[userPorts.User, userDtos.UserResponse](user)

	return map[string]interface{}{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user":          userResponse,
	}, nil
}

// RefreshToken orchestrates token rotation: parse the incoming refresh
// token -> match it against the cached value -> reload the user and
// confirm enablement -> sign fresh access + refresh JWTs and overwrite
// the cached entry. Any divergence yields 401.
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (map[string]interface{}, error) {
	userId, sessionId, err := s.validateRefreshClaims(refreshToken)
	if err != nil {
		return nil, err
	}
	if err := s.assertCachedRefreshMatches(ctx, userId, sessionId, refreshToken); err != nil {
		return nil, err
	}

	user, err := s.loadEnabledUserById(ctx, userId)
	if err != nil {
		return nil, err
	}

	emailStr := ""
	if user.Email != nil {
		emailStr = *user.Email
	}
	accessToken, newRefreshToken := s.signSessionTokens(user.ID.Hex(), emailStr, sessionId)
	s.di.SessionRepo.StoreRefreshToken(ctx, user.ID.Hex(), sessionId, newRefreshToken, constants.RefreshTokenTTL)

	return map[string]interface{}{
		"access_token":  accessToken,
		"refresh_token": newRefreshToken,
		"user":          user,
	}, nil
}

// Logout invalidates the user's refresh token by removing it from cache,
// effectively ending the authenticated session. Token parse failure yields
// 401; cache-delete failure yields 500.
func (s *AuthService) Logout(ctx context.Context, token string) error {
	claims, err := utilsJWT.ParseJWT(s.di.AuthConfig.Secret, token)
	if err != nil {
		return &customErrors.ServerCustomError{
			Code:   httpStatus.UNAUTHORIZED,
			Errors: []string{"Invalid or expired refresh token."},
		}
	}
	if delErr := s.di.SessionRepo.InvalidateRefreshToken(ctx, claims.UserID, claims.ID); delErr != nil {
		return &customErrors.ServerCustomError{
			Code:   httpStatus.INTERNAL_SERVER_ERROR,
			Errors: []string{"Failed to log out user session."},
		}
	}
	return nil
}

// GetMyCoverage returns the list of organizations accessible by the
// caller. Cache-aside via CoverageCacheRepo: try the cache first; on
// miss build the coverage from the source of truth and let the
// repository repopulate the cache.
func (s *AuthService) GetMyCoverage(ctx context.Context, userId string) (map[string]interface{}, error) {
	userAccess, err := s.di.CoverageCacheRepo.GetCachedAccess(ctx, userId)
	if err != nil {
		userAccess, err = s.di.CoverageCacheRepo.BuildCache(ctx, userId)
		if err != nil {
			return nil, &customErrors.ServerCustomError{
				Code:   httpStatus.INTERNAL_SERVER_ERROR,
				Errors: []string{"Failed to fetch user coverage"},
			}
		}
	}
	return map[string]interface{}{
		"organizations": userAccess.Organizations,
		"lastUpdated":   userAccess.LastUpdated,
	}, nil
}

// GetMyPermissions returns the caller's resolved permissions in one
// org. AuthCacheRepo.GetOrBuildCache encapsulates the cache-first /
// build-on-miss path so this method stays a thin envelope around the
// 500-on-error path.
func (s *AuthService) GetMyPermissions(ctx context.Context, userId, orgId string) (map[string]interface{}, error) {
	permissions, version, err := s.di.AuthCacheRepo.GetOrBuildCache(ctx, userId, orgId)
	if err != nil {
		return nil, &customErrors.ServerCustomError{
			Code:   httpStatus.INTERNAL_SERVER_ERROR,
			Errors: []string{"Failed to fetch user permissions"},
		}
	}
	return map[string]interface{}{
		"permissions": permissions,
		"version":     version,
	}, nil
}

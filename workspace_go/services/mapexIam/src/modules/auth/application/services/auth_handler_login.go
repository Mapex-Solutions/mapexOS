package services

import (
	"context"

	"mapexIam/src/modules/auth/application/constants"
	userPorts "mapexIam/src/modules/users/application/ports"

	customErrors "github.com/Mapex-Solutions/mapexGoKit/microservices/http/customErrors"
	httpStatus "github.com/Mapex-Solutions/mapexGoKit/microservices/http/status"
	utilsPassword "github.com/Mapex-Solutions/mapexGoKit/utils/bcrypt/password"
	utilsJWT "github.com/Mapex-Solutions/mapexGoKit/utils/jwt"
)

// loadUserForLogin fetches a user by email and returns it. Missing user
// surfaces as 401 to avoid leaking which addresses are registered.
func (s *AuthService) loadUserForLogin(ctx context.Context, email string) (*userPorts.User, error) {
	user, err := s.di.UserService.GetUserByEmail(ctx, &email)
	if user == nil || err != nil {
		return nil, &customErrors.ServerCustomError{
			Code:   httpStatus.UNAUTHORIZED,
			Errors: []string{"Access is denied due to invalid credentials."},
		}
	}
	return user, nil
}

// verifyLoginCredentials enforces enablement and password match. Wrong
// password yields 401; disabled account yields 403 so administrators can
// distinguish "wrong password" from "blocked user" in audit logs.
func (s *AuthService) verifyLoginCredentials(user *userPorts.User, plainPassword string) error {
	if !user.Enabled {
		return &customErrors.ServerCustomError{
			Code:   httpStatus.FORBIDDEN,
			Errors: []string{"Access is denied, user blocked."},
		}
	}
	if !utilsPassword.CheckPassword(*user.Password, plainPassword) {
		return &customErrors.ServerCustomError{
			Code:   httpStatus.UNAUTHORIZED,
			Errors: []string{"Access is denied due to invalid credentials."},
		}
	}
	return nil
}

// signSessionTokens mints both an access JWT (short-lived) and a refresh
// JWT (longer-lived) for one session. The TTLs live in application/constants
// so login and refresh share the same values.
func (s *AuthService) signSessionTokens(userIdHex, email, sessionId string) (string, string) {
	secretKey := s.di.AuthConfig.Secret
	accessToken, _ := utilsJWT.SignJWT(secretKey, userIdHex, sessionId, email, constants.AccessTokenTTL)
	refreshToken, _ := utilsJWT.SignRefreshToken(secretKey, userIdHex, sessionId, constants.RefreshTokenTTL)
	return accessToken, refreshToken
}

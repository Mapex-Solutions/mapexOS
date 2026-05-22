package handlers

import (
	"github.com/gofiber/fiber/v2"

	"mapexIam/src/modules/auth/application/dtos"
	"mapexIam/src/modules/auth/application/ports"

	middlewaresAuth "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/auth"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/response"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/status"
)

// Login authenticates a user with email and password credentials.
// Returns JWT access token, refresh token, and user information upon successful authentication.
//
// POST /auth/login
// Body: { "email": "user@example.com", "password": "secret" }
//
// Returns:
//   - 200: { "accessToken": "...", "refreshToken": "...", "user": {...} }
//   - 400: Invalid credentials or validation error
//   - 401: Authentication failed
func Login(service ports.AuthServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// retrieve the timeout‐aware Context you set in ContextInjector
		ctx := c.UserContext()

		bodyData, _ := requestValidation.GetDTO[*dtos.LoginDTO](c, "bodyDTO")
		retData, err := service.Login(ctx, bodyData)

		if err == nil {
			return response.Success(c, retData)
		}
		return err
	}
}

// Logout invalidates the user's session by removing the access token from cache.
// Requires valid JWT token in Authorization header.
//
// POST /auth/logout
// Headers: Authorization: Bearer <token>
//
// Returns:
//   - 200: { "success": true }
//   - 401: Missing or invalid token
//   - 500: Failed to logout
func Logout(service ports.AuthServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// Retrieve the timeout‐aware Context you set in ContextInjector
		ctx := c.UserContext()

		// Get token from request
		token, ok := c.Locals("token").(string)
		if !ok || token == "" {
			return response.InternalServerError(c, "invalid or missing token", nil)
		}

		err := service.Logout(ctx, token)

		if err == nil {
			return response.Success(c, map[string]bool{"success": true})
		}
		return err
	}
}

// RefreshToken generates a new access token using a valid refresh token.
// The refresh token is extracted from the request by RefreshTokenExtractor middleware.
//
// POST /auth/refresh
// Headers: Authorization: Bearer <access_token>
// Body: { "refreshToken": "..." }
//
// Returns:
//   - 200: { "accessToken": "...", "refreshToken": "..." }
//   - 401: Invalid or expired refresh token
//   - 500: Failed to refresh token
func RefreshToken(service ports.AuthServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// retrieve the timeout‐aware Context you set in ContextInjector
		ctx := c.UserContext()

		// Get refresh token from request
		refreshToken, ok := c.Locals("refreshToken").(string)
		if !ok || refreshToken == "" {
			return response.InternalServerError(c, "invalid or missing refresh token", nil)
		}

		retData, err := service.RefreshToken(ctx, refreshToken)

		if err == nil {
			return response.Success(c, retData)
		}
		return err
	}
}

// GetMyCoverage returns the list of organizations accessible by the authenticated user.
// Used by the UI to populate organization selectors and navigation tree.
// Tries cache first, builds coverage if cache miss.
//
// GET /auth/users/me/coverage
// Headers: Authorization: Bearer <token>
//
// Returns:
//   - 200: { "organizations": [...], "lastUpdated": "..." }
//   - 401: Missing or invalid userId in token
//   - 500: Failed to fetch coverage
func GetMyCoverage(service ports.AuthServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// retrieve the timeout‐aware Context you set in ContextInjector
		ctx := c.UserContext()

		// Get userId from JWT token
		userId, ok := middlewaresAuth.GetUserIdFromToken(c)
		if !ok {
			return response.Custom(c, status.UNAUTHORIZED, []string{"missing or invalid userId in token"})
		}

		// Call service to get coverage
		result, err := service.GetMyCoverage(ctx, userId)
		if err != nil {
			return err
		}

		return response.Success(c, result)
	}
}

// GetMyPermissions returns the resolved permissions for the authenticated user
// in the current organization context. Reads from Redis cache (O(1)), builds on miss.
//
// GET /auth/me/permissions
// Headers: Authorization: Bearer <token>
// Headers: X-Org-Context: <orgId>
//
// Returns:
//   - 200: { "permissions": ["assets.list", "assets.create", ...], "version": 42 }
//   - 401: Missing or invalid userId in token
//   - 500: Failed to fetch permissions
func GetMyPermissions(service ports.AuthServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// retrieve the timeout-aware Context you set in ContextInjector
		ctx := c.UserContext()

		// Get userId from JWT token
		userId, ok := middlewaresAuth.GetUserIdFromToken(c)
		if !ok {
			return response.Custom(c, status.UNAUTHORIZED, []string{"missing or invalid userId in token"})
		}

		// Get orgId from X-Org-Context header
		orgId := c.Get("X-Org-Context")

		// Call service to get permissions
		result, err := service.GetMyPermissions(ctx, userId, orgId)
		if err != nil {
			return err
		}

		return response.Success(c, result)
	}
}

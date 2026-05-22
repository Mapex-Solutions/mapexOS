package handlers

import (
	"github.com/gofiber/fiber/v2"

	"mapexIam/src/modules/auth/application/dtos"
	"mapexIam/src/modules/auth/domain/repositories"

	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/response"
)

// BuildAuthorizationCache builds the authorization cache for a user in a specific organization.
// This is an internal endpoint that should be protected with API Key middleware.
//
// POST /internal/auth/build-authorization
// Body: { "userId": "...", "orgId": "..." }
//
// Returns:
//   - 200: { "permissions": [...], "version": 42 }
//   - 400: Invalid request
//   - 500: Build failed
func BuildAuthorizationCache(authCacheRepo repositories.AuthorizationCacheRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		bodyData, _ := requestValidation.GetDTO[*dtos.BuildAuthorizationCacheRequest](c, "bodyDTO")

		// Build cache using AuthorizationCacheRepository
		permissions, version, err := authCacheRepo.BuildCache(ctx, bodyData.UserID, bodyData.OrgID)
		if err != nil {
			return response.InternalServerError(c, "Failed to build authorization cache", err)
		}

		return response.Success(c, map[string]interface{}{
			"permissions": permissions,
			"version":     version,
		})
	}
}

// BuildCoverageCache builds the coverage cache for a user (list of accessible organizations).
// This is an internal endpoint that should be protected with API Key middleware.
//
// POST /internal/auth/build-coverage
// Body: { "userId": "..." }
//
// Returns:
//   - 200: { "organizations": [...] }
//   - 400: Invalid request
//   - 500: Build failed
func BuildCoverageCache(coverageCacheRepo repositories.CoverageCacheRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		bodyData, _ := requestValidation.GetDTO[*dtos.BuildCoverageCacheRequest](c, "bodyDTO")

		// Build cache using CoverageCacheRepository
		organizations, err := coverageCacheRepo.BuildCache(ctx, bodyData.UserID)
		if err != nil {
			return response.InternalServerError(c, "Failed to build coverage cache", err)
		}

		return response.Success(c, map[string]interface{}{
			"organizations": organizations,
		})
	}
}

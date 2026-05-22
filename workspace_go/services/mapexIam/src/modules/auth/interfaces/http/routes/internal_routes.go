package routes

import (
	"github.com/gofiber/fiber/v2"

	"mapexIam/src/modules/auth/application/dtos"
	"mapexIam/src/modules/auth/domain/repositories"
	"mapexIam/src/modules/auth/interfaces/http/handlers"

	validation "github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
)

// RegisterInternalRoutes registers internal routes for cache building operations.
// These routes are protected with API Key middleware (configured in module.go).
// Should only be called by internal services.
//
// Endpoints:
//   - POST /internal/auth/build-authorization - Build authorization cache for a user in an organization
//   - POST /internal/auth/build-coverage - Build coverage cache for a user
//
// All routes require X-API-Key header with valid API key.
func RegisterInternalRoutes(
	group fiber.Router,
	authCacheRepo repositories.AuthorizationCacheRepository,
	coverageCacheRepo repositories.CoverageCacheRepository,
) {
	// Build authorization cache (permissions for user in org)
	buildAuthDto := validation.NewValidation(&dtos.BuildAuthorizationCacheRequest{}, nil, nil)
	group.Post("/build-authorization", validation.ValidationMiddleware(buildAuthDto), handlers.BuildAuthorizationCache(authCacheRepo))

	// Build coverage cache (list of accessible organizations)
	buildCoverageDto := validation.NewValidation(&dtos.BuildCoverageCacheRequest{}, nil, nil)
	group.Post("/build-coverage", validation.ValidationMiddleware(buildCoverageDto), handlers.BuildCoverageCache(coverageCacheRepo))
}

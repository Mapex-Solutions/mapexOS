package routes

import (
	"mapexIam/src/modules/organizations/application/ports"

	"github.com/gofiber/fiber/v2"
)

// RegisterInternalRoutes registers internal API routes for organizations.
// These routes are used for inter-service communication (e.g., cache population).
// Authentication: X-API-Key header (no JWT).
//
// Internal routes are designed for:
//   - Low-latency service-to-service communication
//   - Bypassing permission checks (API Key handles authorization)
//
// Security: These routes should ONLY be accessible via internal network/service mesh.
func RegisterInternalRoutes(app fiber.Router, service ports.OrganizationServicePort) {
	_ = app.Group("/api/internal/v1/organizations")
	_ = service
	// Retention policies endpoint removed - now served by Events service at /api/v1/retention
}

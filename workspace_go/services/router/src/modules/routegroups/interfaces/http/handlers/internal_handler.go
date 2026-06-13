package handlers

import (
	"strings"

	"router/src/modules/routegroups/application/dtos"
	"router/src/modules/routegroups/application/ports"

	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/response"
	web "github.com/Mapex-Solutions/mapexGoKit/microservices/http/web"
)

// GetRouteGroupsByIds handles internal API requests to fetch multiple route groups by IDs.
// This endpoint is designed for MS-to-MS communication via API Key authentication.
//
// Query Parameters:
//   - ids (required): Comma-separated list of route group IDs (e.g., "id1,id2,id3")
//   - projection (optional): Comma-separated fields to return (e.g., "name,enabled,version")
//
// Authentication: API Key (internal only)
//
// Example:
//
//	GET /internal/route_groups?ids=507f1f77bcf86cd799439011,507f1f77bcf86cd799439012&projection=name,enabled
//
// Response:
//
//	[
//	  { "id": "507f1f77bcf86cd799439011", "name": "Group 1", "enabled": true },
//	  { "id": "507f1f77bcf86cd799439012", "name": "Group 2", "enabled": true }
//	]
//
// Parameters:
//   - service: RouteGroupServicePort interface for business logic
//
// Returns:
//   - web.Handler: Configured handler function
func GetRouteGroupsByIds(service ports.RouteGroupServicePort) web.Handler {
	return func(c *web.Ctx) error {
		ctx := c.UserContext()

		// Get validated query DTO
		queryData, _ := requestValidation.GetDTO[*dtos.RouteGroupInternalIdsQuery](c, "queryDTO")

		// Split comma-separated IDs into array
		idsStr := strings.TrimSpace(queryData.Ids)
		if idsStr == "" {
			return response.InternalServerError(c, "ids parameter is required", nil)
		}

		ids := strings.Split(idsStr, ",")

		// Trim whitespace from each ID
		for i := range ids {
			ids[i] = strings.TrimSpace(ids[i])
		}

		// Call service to fetch route groups by IDs
		routeGroups, err := service.GetRouteGroupsByIds(ctx, ids)

		if err != nil {
			return response.InternalServerError(c, err.Error(), nil)
		}

		// TODO: Apply projection if provided
		// For now, return full objects (projection can be implemented later)

		return response.Success(c, routeGroups)
	}
}

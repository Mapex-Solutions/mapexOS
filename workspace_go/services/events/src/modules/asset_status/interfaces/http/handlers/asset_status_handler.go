package handlers

import (
	"github.com/gofiber/fiber/v2"

	"events/src/modules/asset_status/application/dtos"
	"events/src/modules/asset_status/application/ports"

	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/response"
)

// GetConnectivityHistory returns a Fiber handler serving
// GET /api/v1/events/assets/:assetUUID/connectivity_history (asset-scoped).
//
// The :assetUUID path param is injected into the query's AssetUUID filter
// before delegating to the service, so the service layer always sees the
// same query shape regardless of which route was hit.
//
// Returns:
//   - 200 OK with cursor-paginated results { items, nextCursor, prevCursor, hasNext, hasPrevious }
//   - 400 Bad Request if validation fails
//   - 403 Forbidden if request context missing
//   - 500 Internal Server Error on service failure
func GetConnectivityHistory(service ports.AssetStatusServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		requestContext, ok := c.Locals("requestContext").(*reqCtx.RequestContext)
		if !ok {
			return fiber.NewError(fiber.StatusForbidden, "Request context not found")
		}

		queryData, _ := requestValidation.GetDTO[*dtos.AssetConnectivityHistoryQuery](c, "queryDTO")

		if assetUUID := c.Params("assetUUID"); assetUUID != "" {
			queryData.AssetUUID = &assetUUID
		}

		result, err := service.ListAssetConnectivityHistory(ctx, requestContext, queryData)
		if err != nil {
			return err
		}
		return response.Success(c, result)
	}
}

// ListConnectivityHistory returns a Fiber handler serving
// GET /api/v1/events/connectivity_history (org-wide listing).
//
// No path param. The handler delegates to the same service method as the
// asset-scoped variant — the query's AssetUUID field is optional here.
//
// Returns: same shape as GetConnectivityHistory.
func ListConnectivityHistory(service ports.AssetStatusServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		requestContext, ok := c.Locals("requestContext").(*reqCtx.RequestContext)
		if !ok {
			return fiber.NewError(fiber.StatusForbidden, "Request context not found")
		}

		queryData, _ := requestValidation.GetDTO[*dtos.AssetConnectivityHistoryQuery](c, "queryDTO")

		result, err := service.ListAssetConnectivityHistory(ctx, requestContext, queryData)
		if err != nil {
			return err
		}
		return response.Success(c, result)
	}
}

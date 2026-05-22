package handlers

import (
	"github.com/gofiber/fiber/v2"

	"assets/src/modules/assets/application/dtos"
	"assets/src/modules/assets/application/ports"

	contractsCommon "github.com/Mapex-Solutions/MapexOS/contracts/common"
	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/response"
)

// CreateAsset returns a Fiber handler that creates a new asset.
//
// Following Hexagonal Architecture, this handler accepts the service port interface
// and delegates business logic to the service layer.
//
// It uses RequestContext (injected by coverage middleware) which contains:
//   - OrgContext: The selected organization ID from X-Org-Context header
//   - OrgContextData: Organization data including PathKey for hierarchical filtering
//
// The handler passes the full RequestContext to the service layer, which extracts
// the needed fields (orgId, pathKey) for multi-tenant support.
//
// It expects a validated DTO of type dtos.AssetCreateDTO to be stored
// in the Fiber context under the key "bodyDTO" (usually populated by
// requestValidation middleware).
//
// Parameters:
//   - service: The AssetServicePort interface for asset business operations
//
// Returns:
//   - A Fiber handler function that processes the asset creation request
func CreateAsset(service ports.AssetServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		// Get RequestContext from coverage middleware
		requestContext, ok := c.Locals("requestContext").(*reqCtx.RequestContext)
		if !ok {
			return response.InternalServerError(c, "requestContext not found in request context", nil)
		}

		bodyData, _ := requestValidation.GetDTO[*dtos.AssetCreateDTO](c, "bodyDTO")

		// Pass requestContext to service (contains OrgContext and OrgContextData)
		retData, err := service.CreateAsset(ctx, requestContext, bodyData)

		if err != nil {
			return err
		}
		return response.Created(c, retData)
	}
}

// GetAssetById returns a Fiber handler that retrieves an asset by its ID.
//
// Following Hexagonal Architecture, this handler accepts the service port interface
// and delegates business logic to the service layer.
//
// It expects a validated DTO of type dtos.AssetIdDto to be stored
// in the Fiber context under the key "paramsDTO" (usually populated by
// requestValidation middleware).
//
// Parameters:
//   - service: The AssetServicePort interface for asset business operations
//
// Returns:
//   - A Fiber handler function that processes the asset retrieval request
func GetAssetById(service ports.AssetServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		asset, _ := requestValidation.GetDTO[*dtos.AssetIdDto](c, "paramsDTO")
		retData, err := service.GetAssetById(ctx, &asset.AssetId)

		if err != nil {
			return err
		}
		return response.Success(c, retData)
	}
}

// UpdateAssetById returns a Fiber handler that updates an asset by its ID.
//
// Following Hexagonal Architecture, this handler accepts the service port interface
// and delegates business logic to the service layer.
//
// It expects two validated DTOs:
//   - dtos.AssetIdDto stored in the Fiber context under the key "paramsDTO"
//   - dtos.AssetUpdateDTO stored in the Fiber context under the key "bodyDTO"
//
// (Both are usually populated by requestValidation middleware)
//
// Parameters:
//   - service: The AssetServicePort interface for asset business operations
//
// Returns:
//   - A Fiber handler function that processes the asset update request
func UpdateAssetById(service ports.AssetServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Retrieve the timeout-aware Context you set in ContextInjector
		ctx := c.UserContext()

		asset, _ := requestValidation.GetDTO[*dtos.AssetIdDto](c, "paramsDTO")
		bodyData, _ := requestValidation.GetDTO[*dtos.AssetUpdateDTO](c, "bodyDTO")
		retData, err := service.UpdateAssetById(ctx, &asset.AssetId, bodyData)

		if err != nil {
			return err
		}
		return response.Success(c, retData)
	}
}

// GenerateMqttPassword returns a Fiber handler that issues a fresh
// random alphanumeric MQTT password for the operator to drop into the
// asset form. Stateless — the endpoint reads no asset and writes no
// state. The operator remains free to type a custom password instead;
// the platform validates only at create / change-password time.
func GenerateMqttPassword(service ports.AssetServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()
		retData, err := service.GenerateMqttPassword(ctx)
		if err != nil {
			return err
		}
		return response.Success(c, retData)
	}
}

// DeleteAssetById returns a Fiber handler that deletes an asset by its ID.
//
// Following Hexagonal Architecture, this handler accepts the service port interface
// and delegates business logic to the service layer.
//
// It expects a validated DTO of type dtos.AssetIdDto to be stored
// in the Fiber context under the key "paramsDTO" (usually populated by
// requestValidation middleware).
//
// Parameters:
//   - service: The AssetServicePort interface for asset business operations
//
// Returns:
//   - A Fiber handler function that processes the asset deletion request
func DeleteAssetById(service ports.AssetServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		asset, _ := requestValidation.GetDTO[*dtos.AssetIdDto](c, "paramsDTO")
		retData, err := service.DeleteAssetById(ctx, &asset.AssetId)

		if err != nil {
			return err
		}
		return response.Success(c, retData)
	}
}

// GetAssets returns a Fiber handler that retrieves a paginated and filtered list of assets.
// Uses RequestContext from coverage middleware for context-aware org filtering with hierarchical support.
//
// The handler extracts RequestContext from c.Locals("requestContext") which was set by InjectRequestContext middleware.
// This provides access to:
//   - ScopedOrgIds: All accessible organization IDs
//   - OrgContext: Optional org filter from X-Org-Context header
//   - OrgContextData: Detailed org data with PathKey
//   - CoverageOrgs: Full coverage data with hierarchical information
//
// Query supports hierarchical filtering via includeChildren parameter:
//   - OrgContext + includeChildren=true: Returns org and all descendants (PathKey range)
//   - OrgContext + includeChildren=false: Returns specific org only
//   - No OrgContext: Returns all accessible orgs
//
// It expects a validated DTO of type dtos.AssetQueryDTO in the "queryDTO" context key
// (populated by requestValidation middleware) containing optional filters such as:
//   - name, status, assetTemplateId, category, assetType (filters)
//   - page, perPage (pagination)
//   - projection (field selection)
//   - includeChildren (hierarchical query flag)
//
// Returns:
//   - 200 OK with paginated asset list
//   - 400 Bad Request if query validation fails
//   - 500 Internal Server Error on service failure or requestContext not found
func GetAssets(service ports.AssetServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Retrieve the timeout-aware Context you set in ContextInjector
		ctx := c.UserContext()

		// Get RequestContext from coverage middleware
		requestContext, ok := c.Locals("requestContext").(*reqCtx.RequestContext)
		if !ok {
			return response.InternalServerError(c, "requestContext not found in request context", nil)
		}

		queryData, _ := requestValidation.GetDTO[*dtos.AssetQueryDTO](c, "queryDTO")
		retData, err := service.GetAssets(ctx, requestContext, queryData)

		if err != nil {
			return err
		}
		return response.Success(c, retData)
	}
}

// GetAssetCount returns a Fiber handler that returns the total count of assets.
// Uses cached count with 6h TTL, invalidated on create/delete.
//
// Parameters:
//   - service: The AssetServicePort interface for asset business operations
//
// Returns:
//   - A Fiber handler function that processes the asset count request
func GetAssetCount(service ports.AssetServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		requestContext, ok := c.Locals("requestContext").(*reqCtx.RequestContext)
		if !ok {
			return response.InternalServerError(c, "requestContext not found in request context", nil)
		}

		count, err := service.CountAssets(ctx, requestContext)
		if err != nil {
			return err
		}

		return response.Success(c, contractsCommon.CounterResponse{Count: count})
	}
}

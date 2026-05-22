package handlers

import (
	"github.com/gofiber/fiber/v2"

	"assets/src/modules/assettemplates/application/dtos"
	"assets/src/modules/assettemplates/application/ports"

	contractsCommon "github.com/Mapex-Solutions/MapexOS/contracts/common"
	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/response"
)

// CreateAssetTemplate returns a Fiber handler that creates a new assettemplate.
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
// It expects a validated DTO of type dtos.AssetTemplateCreateDTO to be stored
// in the Fiber context under the key "bodyDTO" (usually populated by
// requestValidation middleware).
//
// Parameters:
//   - service: The AssetTemplateServicePort interface for asset template business operations
//
// Returns:
//   - A Fiber handler function that processes the asset template creation request
func CreateAssetTemplate(service ports.AssetTemplateServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		// Get RequestContext from coverage middleware
		requestContext, ok := c.Locals("requestContext").(*reqCtx.RequestContext)
		if !ok {
			return response.InternalServerError(c, "requestContext not found in request context", nil)
		}

		bodyData, _ := requestValidation.GetDTO[*dtos.AssetTemplateCreateDTO](c, "bodyDTO")

		// Pass requestContext to service (contains OrgContext and OrgContextData)
		retData, err := service.CreateAssetTemplate(ctx, requestContext, bodyData)

		if err != nil {
			return err
		}
		return response.Created(c, retData)
	}
}

// GetAssetTemplateById returns a Fiber handler that retrieves an asset template by its ID.
//
// Following Hexagonal Architecture, this handler accepts the service port interface
// and delegates business logic to the service layer.
//
// It expects a validated DTO of type dtos.AssetTemplateIdDto to be stored
// in the Fiber context under the key "paramsDTO" (usually populated by
// requestValidation middleware).
//
// Parameters:
//   - service: The AssetTemplateServicePort interface for asset template business operations
//
// Returns:
//   - A Fiber handler function that processes the asset template retrieval request
func GetAssetTemplateById(service ports.AssetTemplateServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// retrieve the timeout‐aware Context you set in ContextInjector
		ctx := c.UserContext()

		assettemplate, _ := requestValidation.GetDTO[*dtos.AssetTemplateIdDto](c, "paramsDTO")
		retData, err := service.GetAssetTemplateById(ctx, &assettemplate.AssetTemplateId)

		if err != nil {
			return err
		}
		return response.Success(c, retData)
	}
}

// UpdateAssetTemplateById returns a Fiber handler that updates an asset template by its ID.
//
// Following Hexagonal Architecture, this handler accepts the service port interface
// and delegates business logic to the service layer.
//
// It expects two validated DTOs:
//   - dtos.AssetTemplateIdDto stored in the Fiber context under the key "paramsDTO"
//   - dtos.AssetTemplateUpdateDTO stored in the Fiber context under the key "bodyDTO"
//
// (Both are usually populated by requestValidation middleware)
//
// Parameters:
//   - service: The AssetTemplateServicePort interface for asset template business operations
//
// Returns:
//   - A Fiber handler function that processes the asset template update request
func UpdateAssetTemplateById(service ports.AssetTemplateServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// retrieve the timeout‐aware Context you set in ContextInjector
		ctx := c.UserContext()

		assettemplate, _ := requestValidation.GetDTO[*dtos.AssetTemplateIdDto](c, "paramsDTO")
		bodyData, _ := requestValidation.GetDTO[*dtos.AssetTemplateUpdateDTO](c, "bodyDTO")
		retData, err := service.UpdateAssetTemplateById(ctx, &assettemplate.AssetTemplateId, bodyData)

		if err != nil {
			return err
		}
		return response.Success(c, retData)
	}
}

// DeleteAssetTemplateById returns a Fiber handler that deletes an asset template by its ID.
//
// Following Hexagonal Architecture, this handler accepts the service port interface
// and delegates business logic to the service layer.
//
// It expects a validated DTO of type dtos.AssetTemplateIdDto to be stored
// in the Fiber context under the key "paramsDTO" (usually populated by
// requestValidation middleware).
//
// Parameters:
//   - service: The AssetTemplateServicePort interface for asset template business operations
//
// Returns:
//   - A Fiber handler function that processes the asset template deletion request
func DeleteAssetTemplateById(service ports.AssetTemplateServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// retrieve the timeout‐aware Context you set in ContextInjector
		ctx := c.UserContext()

		assettemplate, _ := requestValidation.GetDTO[*dtos.AssetTemplateIdDto](c, "paramsDTO")
		retData, err := service.DeleteAssetTemplateById(ctx, &assettemplate.AssetTemplateId)

		if err != nil {
			return err
		}
		return response.Success(c, retData)
	}
}

// GetAssetTemplates returns a Fiber handler that retrieves a paginated and filtered list of asset templates.
// Uses scopedOrgIds from coverage middleware for multi-tenant access control.
//
// The handler extracts scopedOrgIds from c.Locals() which was set by the coverage middleware.
// This ensures users can only query asset templates within their accessible organizations.
//
// Parameters:
//   - service: The AssetTemplateServicePort interface for asset template business operations
//
// Returns:
//   - A Fiber handler function that processes the asset template listing request
func GetAssetTemplates(service ports.AssetTemplateServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Retrieve the timeout-aware Context you set in ContextInjector
		ctx := c.UserContext()

		// Get RequestContext from coverage middleware
		requestContext, ok := c.Locals("requestContext").(*reqCtx.RequestContext)
		if !ok {
			return response.InternalServerError(c, "requestContext not found in request context", nil)
		}

		queryData, _ := requestValidation.GetDTO[*dtos.AssetTemplateQueryDto](c, "queryDTO")
		retData, err := service.GetAssetTemplates(ctx, requestContext, queryData)

		if err != nil {
			return err
		}
		return response.Success(c, retData)
	}
}

// GetAvailableFields returns a Fiber handler that retrieves only the available fields
// of an asset template by its ID.
//
// This endpoint is optimized for performance using Redis cache with 24-hour TTL.
// It returns a lightweight response containing only the availableFields array,
// which is used by the Rules module for autocomplete in Event conditions.
//
// Following Hexagonal Architecture, this handler accepts the service port interface
// and delegates business logic to the service layer.
//
// Response format:
//
//	{
//	  "availableFields": ["eventType", "data.temperature", "data.location.lat", ...]
//	}
//
// Parameters:
//   - service: The AssetTemplateServicePort interface for asset template business operations
//
// Returns:
//   - A Fiber handler function that processes the available fields retrieval request
func GetAvailableFields(service ports.AssetTemplateServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		// Get RequestContext from coverage middleware (for future multi-tenant validation)
		requestContext, ok := c.Locals("requestContext").(*reqCtx.RequestContext)
		if !ok {
			return response.InternalServerError(c, "requestContext not found in request context", nil)
		}

		// Get and validate assetTemplateId from params
		assetTemplateIdDto, _ := requestValidation.GetDTO[*dtos.AssetTemplateIdDto](c, "paramsDTO")

		// Call service to get available fields (with caching)
		result, err := service.GetAvailableFields(ctx, &assetTemplateIdDto.AssetTemplateId, requestContext)

		if err != nil {
			return err
		}

		return response.Success(c, result)
	}
}

// GetAssetTemplateCount returns a Fiber handler that returns the total count of asset templates.
// Uses cached count with 6h TTL, invalidated on create/delete.
//
// Parameters:
//   - service: The AssetTemplateServicePort interface for asset template business operations
//
// Returns:
//   - A Fiber handler function that processes the asset template count request
func GetAssetTemplateCount(service ports.AssetTemplateServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		requestContext, ok := c.Locals("requestContext").(*reqCtx.RequestContext)
		if !ok {
			return response.InternalServerError(c, "requestContext not found in request context", nil)
		}

		count, err := service.CountAssetTemplates(ctx, requestContext)
		if err != nil {
			return err
		}

		return response.Success(c, contractsCommon.CounterResponse{Count: count})
	}
}

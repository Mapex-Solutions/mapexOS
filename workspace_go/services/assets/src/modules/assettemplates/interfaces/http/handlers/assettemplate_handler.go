package handlers

import (
	web "github.com/Mapex-Solutions/mapexGoKit/microservices/http/web"

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
func CreateAssetTemplate(service ports.AssetTemplateServicePort) web.Handler {
	return func(c *web.Ctx) error {
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
func GetAssetTemplateById(service ports.AssetTemplateServicePort) web.Handler {
	return func(c *web.Ctx) error {

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
func UpdateAssetTemplateById(service ports.AssetTemplateServicePort) web.Handler {
	return func(c *web.Ctx) error {

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
func DeleteAssetTemplateById(service ports.AssetTemplateServicePort) web.Handler {
	return func(c *web.Ctx) error {

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
func GetAssetTemplates(service ports.AssetTemplateServicePort) web.Handler {
	return func(c *web.Ctx) error {
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
func GetAvailableFields(service ports.AssetTemplateServicePort) web.Handler {
	return func(c *web.Ctx) error {
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
func GetAssetTemplateCount(service ports.AssetTemplateServicePort) web.Handler {
	return func(c *web.Ctx) error {
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

// GetFieldVocabulary returns a Fiber handler that serves the curated,
// multi-tenant field vocabulary to the asset-template authoring UI.
//
// The fields are grouped by category in a fixed display order, and each
// field's hint plus the group label are resolved to the language requested
// via the optional ?lang= query parameter (en-US fallback).
//
// Response data shape:
//
//	{
//	  "groups": [
//	    { "category": "climate", "label": "Climate", "fields": [
//	        { "value": "temperature", "hint": "Ambient temperature", "type": "number", "unit": "°C" }
//	    ] }
//	  ]
//	}
//
// Parameters:
//   - service: The AssetTemplateServicePort interface for asset template business operations
//
// Returns:
//   - A Fiber handler function that processes the field vocabulary request
func GetFieldVocabulary(service ports.AssetTemplateServicePort) web.Handler {
	return func(c *web.Ctx) error {
		ctx := c.UserContext()

		requestContext, ok := c.Locals("requestContext").(*reqCtx.RequestContext)
		if !ok {
			return response.InternalServerError(c, "requestContext not found in request context", nil)
		}

		queryData, _ := requestValidation.GetDTO[*dtos.FieldVocabularyQuery](c, "queryDTO")
		lang := ""
		if queryData != nil {
			lang = queryData.Lang
		}

		result, err := service.GetFieldVocabulary(ctx, requestContext, lang)
		if err != nil {
			return err
		}
		return response.Success(c, result)
	}
}

// CreateMigrationPlan creates a scheduled template migration plan.
// Uses RequestContext (coverage middleware) to scope the plan to the caller's org.
func CreateMigrationPlan(service ports.AssetTemplateServicePort) web.Handler {
	return func(c *web.Ctx) error {
		ctx := c.UserContext()

		requestContext, ok := c.Locals("requestContext").(*reqCtx.RequestContext)
		if !ok {
			return response.InternalServerError(c, "requestContext not found in request context", nil)
		}

		bodyData, _ := requestValidation.GetDTO[*dtos.MigrationPlanCreateRequest](c, "bodyDTO")
		retData, err := service.CreateMigrationPlan(ctx, requestContext, bodyData)
		if err != nil {
			return err
		}
		return response.Created(c, retData)
	}
}

// ListMigrationPlans returns a paginated, filtered list of migration plans.
// Uses RequestContext (coverage middleware) for automatic org filtering.
func ListMigrationPlans(service ports.AssetTemplateServicePort) web.Handler {
	return func(c *web.Ctx) error {
		ctx := c.UserContext()

		requestContext, ok := c.Locals("requestContext").(*reqCtx.RequestContext)
		if !ok {
			return response.InternalServerError(c, "requestContext not found in request context", nil)
		}

		queryData, _ := requestValidation.GetDTO[*dtos.MigrationPlanQueryDTO](c, "queryDTO")
		retData, err := service.GetMigrationPlans(ctx, requestContext, queryData)
		if err != nil {
			return err
		}
		return response.Success(c, retData)
	}
}

// GetMigrationPlan returns a single migration plan by its id.
func GetMigrationPlan(service ports.AssetTemplateServicePort) web.Handler {
	return func(c *web.Ctx) error {
		ctx := c.UserContext()

		params, _ := requestValidation.GetDTO[*dtos.MigrationPlanIdDto](c, "paramsDTO")
		retData, err := service.GetMigrationPlanById(ctx, &params.Id)
		if err != nil {
			return err
		}
		return response.Success(c, retData)
	}
}

// UpdateMigrationPlan partially updates an editable migration plan.
func UpdateMigrationPlan(service ports.AssetTemplateServicePort) web.Handler {
	return func(c *web.Ctx) error {
		ctx := c.UserContext()

		params, _ := requestValidation.GetDTO[*dtos.MigrationPlanIdDto](c, "paramsDTO")
		bodyData, _ := requestValidation.GetDTO[*dtos.MigrationPlanUpdateRequest](c, "bodyDTO")
		retData, err := service.UpdateMigrationPlanById(ctx, &params.Id, bodyData)
		if err != nil {
			return err
		}
		return response.Success(c, retData)
	}
}

// CancelMigrationPlan cancels a pending or scheduled migration plan.
func CancelMigrationPlan(service ports.AssetTemplateServicePort) web.Handler {
	return func(c *web.Ctx) error {
		ctx := c.UserContext()

		params, _ := requestValidation.GetDTO[*dtos.MigrationPlanIdDto](c, "paramsDTO")
		if err := service.CancelMigrationPlanById(ctx, &params.Id); err != nil {
			return err
		}
		return response.Success(c, map[string]bool{"success": true})
	}
}

// ListMigrationExecutions returns the per-asset executions of a migration plan.
func ListMigrationExecutions(service ports.AssetTemplateServicePort) web.Handler {
	return func(c *web.Ctx) error {
		ctx := c.UserContext()

		params, _ := requestValidation.GetDTO[*dtos.MigrationPlanIdDto](c, "paramsDTO")
		queryData, _ := requestValidation.GetDTO[*dtos.MigrationExecutionQueryDTO](c, "queryDTO")
		retData, err := service.GetMigrationExecutions(ctx, &params.Id, queryData)
		if err != nil {
			return err
		}
		return response.Success(c, retData)
	}
}

// InstallFromMarketplace returns a handler that installs a marketplace asset
// template into the caller's organization (org-scoped, never system).
func InstallFromMarketplace(service ports.AssetTemplateServicePort) web.Handler {
	return func(c *web.Ctx) error {
		requestContext, ok := c.Locals("requestContext").(*reqCtx.RequestContext)
		if !ok {
			return response.InternalServerError(c, "requestContext not found in request context", nil)
		}
		params, _ := requestValidation.GetDTO[*dtos.InstallParams](c, "paramsDTO")
		body, _ := requestValidation.GetDTO[*dtos.InstallBody](c, "bodyDTO")
		result, err := service.InstallFromMarketplace(c.UserContext(), requestContext, params.MarketplaceVendor, params.MarketplaceSlug, body.ShareWithChildren)
		if err != nil {
			return err
		}
		return response.Success(c, result)
	}
}

// UninstallFromMarketplace returns a handler that removes the caller org's
// installation of a marketplace asset template. The shared content stays.
func UninstallFromMarketplace(service ports.AssetTemplateServicePort) web.Handler {
	return func(c *web.Ctx) error {
		requestContext, ok := c.Locals("requestContext").(*reqCtx.RequestContext)
		if !ok {
			return response.InternalServerError(c, "requestContext not found in request context", nil)
		}
		params, _ := requestValidation.GetDTO[*dtos.InstallParams](c, "paramsDTO")
		if err := service.UninstallFromMarketplace(c.UserContext(), requestContext, params.MarketplaceVendor, params.MarketplaceSlug); err != nil {
			return err
		}
		return response.Success(c, map[string]bool{"success": true})
	}
}

// CheckInstalled returns a handler reporting which of the given marketplace GUIDs
// the caller org has installed, in one call — driving the listing's Install vs
// Uninstall toggle.
func CheckInstalled(service ports.AssetTemplateServicePort) web.Handler {
	return func(c *web.Ctx) error {
		requestContext, ok := c.Locals("requestContext").(*reqCtx.RequestContext)
		if !ok {
			return response.InternalServerError(c, "requestContext not found in request context", nil)
		}
		body, _ := requestValidation.GetDTO[*dtos.MarketplaceInstalledCheckRequest](c, "bodyDTO")
		installed, err := service.InstalledGuids(c.UserContext(), requestContext, body.MarketplaceGuids)
		if err != nil {
			return err
		}
		return response.Success(c, dtos.MarketplaceInstalledCheckResponse{Installed: installed})
	}
}

package ports

import (
	ctx "context"

	"assets/src/modules/assettemplates/application/dtos"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

// AssetTemplateServicePort defines the contract for asset template business operations.
//
// This port interface enables Hexagonal Architecture by decoupling the business
// logic from its implementation. Handlers and routes depend on this interface
// rather than the concrete service implementation.
//
// Methods:
//   - CreateAssetTemplate: Creates a new asset template
//   - GetAssetTemplateById: Retrieves an asset template by its ID
//   - UpdateAssetTemplateById: Updates an existing asset template
//   - DeleteAssetTemplateById: Deletes an asset template
type AssetTemplateServicePort interface {
	// CreateAssetTemplate creates a new asset template from the provided DTO.
	// Uses RequestContext to populate orgId and pathKey for multi-tenant support.
	//
	// Parameters:
	//   - ctx: Context for controlling cancellation and timeouts
	//   - requestContext: Contains OrgContext (selected orgId) and OrgContextData (pathKey)
	//   - dto: Asset template creation data
	//
	// Returns:
	//   - AssetTemplateResponse: The created asset template
	//   - error: If creation fails
	CreateAssetTemplate(ctx ctx.Context, requestContext *reqCtx.RequestContext, dto *dtos.AssetTemplateCreateDTO) (*dtos.AssetTemplateResponse, error)

	// GetAssetTemplateById retrieves an asset template by its unique identifier.
	//
	// Parameters:
	//   - ctx: Context for controlling cancellation and timeouts
	//   - assetTemplateId: The unique identifier of the asset template
	//
	// Returns:
	//   - AssetTemplateResponse: The asset template if found
	//   - error: If not found or retrieval fails
	GetAssetTemplateById(ctx ctx.Context, assetTemplateId *string) (*dtos.AssetTemplateResponse, error)

	// UpdateAssetTemplateById updates an existing asset template.
	//
	// Parameters:
	//   - ctx: Context for controlling cancellation and timeouts
	//   - assetTemplateId: The unique identifier of the asset template
	//   - dto: Update data (only non-nil fields will be updated)
	//
	// Returns:
	//   - AssetTemplateResponse: The updated asset template
	//   - error: If update fails or asset template not found
	UpdateAssetTemplateById(ctx ctx.Context, assetTemplateId *string, dto *dtos.AssetTemplateUpdateDTO) (*dtos.AssetTemplateResponse, error)

	// DeleteAssetTemplateById removes an asset template by its unique identifier.
	//
	// Parameters:
	//   - ctx: Context for controlling cancellation and timeouts
	//   - assetTemplateId: The unique identifier of the asset template to delete
	//
	// Returns:
	//   - map[string]bool: Success indicator ({"success": true})
	//   - error: If deletion fails or asset template not found
	DeleteAssetTemplateById(ctx ctx.Context, assetTemplateId *string) (map[string]bool, error)

	// GetAssetTemplates retrieves a paginated and filtered list of asset templates.
	// Uses RequestContext from coverage middleware for context-aware org filtering with hierarchical support.
	//
	// Parameters:
	//   - ctx: Context for controlling cancellation and timeouts
	//   - requestContext: Contains ScopedOrgIds, OrgContext, OrgContextData with PathKey
	//   - query: Filters, pagination, projection, and includeChildren flag
	//
	// Returns:
	//   - PaginatedResult: Matching asset templates and pagination metadata
	//   - error: If query fails
	//
	// Security: Uses orgfilter.BuildOrgFilter() for automatic org filtering based on RequestContext
	GetAssetTemplates(ctx ctx.Context, requestContext *reqCtx.RequestContext, query *dtos.AssetTemplateQueryDto) (*model.PaginatedResult[dtos.AssetTemplateResponse], error)

	// CountAssetTemplates returns the total count of asset templates for the given org context.
	// Uses Redis cache with 6h TTL, invalidated on create/delete.
	//
	// Parameters:
	//   - ctx: Context for controlling cancellation and timeouts
	//   - requestContext: Contains org access data from coverage middleware
	//
	// Returns:
	//   - int64: Total count of matching asset templates
	//   - error: If query fails
	CountAssetTemplates(ctx ctx.Context, requestContext *reqCtx.RequestContext) (int64, error)

	// UpdateManufacturerName updates the manufacturerName field in all asset templates
	// that reference the given manufacturer list ID.
	// Called by NATS consumer when a manufacturer list name changes.
	//
	// Parameters:
	//   - ctx: Context for controlling cancellation and timeouts
	//   - manufacturerId: The list ID (MongoDB ObjectId as string)
	//   - newName: The new manufacturer name
	//
	// Returns:
	//   - error: If update fails
	UpdateManufacturerName(ctx ctx.Context, manufacturerId string, newName string) error

	// UpdateModelName updates the modelName field in all asset templates
	// that reference the given model list ID.
	// Called by NATS consumer when a model list name changes.
	//
	// Parameters:
	//   - ctx: Context for controlling cancellation and timeouts
	//   - modelId: The list ID (MongoDB ObjectId as string)
	//   - newName: The new model name
	//
	// Returns:
	//   - error: If update fails
	UpdateModelName(ctx ctx.Context, modelId string, newName string) error

	// UpdateCategoryName updates the categoryName field in all asset templates
	// that reference the given category list ID.
	// Called by NATS consumer when a category list name changes.
	//
	// Parameters:
	//   - ctx: Context for controlling cancellation and timeouts
	//   - categoryId: The list ID (MongoDB ObjectId as string)
	//   - newName: The new category name
	//
	// Returns:
	//   - error: If update fails
	UpdateCategoryName(ctx ctx.Context, categoryId string, newName string) error

	// GetAvailableFields retrieves only the available fields of an asset template.
	// Uses Redis cache with 24-hour TTL for performance optimization.
	// Cache is invalidated when template is created or updated.
	//
	// Parameters:
	//   - ctx: Context for controlling cancellation and timeouts
	//   - assetTemplateId: The unique identifier of the asset template
	//   - requestContext: Contains OrgContext for multi-tenant validation
	//
	// Returns:
	//   - map with availableFields, availableFieldsUpdatedAt, and cached flag
	//   - error: If retrieval fails or template not found
	GetAvailableFields(ctx ctx.Context, assetTemplateId *string, requestContext *reqCtx.RequestContext) (map[string]interface{}, error)

	// HandleListNameUpdated processes a NATS list-name-updated message, unmarshals
	// the ListNameUpdatedEvent payload, and dispatches to the appropriate denormalization
	// handler (manufacturer/model/category). Acks/Nacks/Rejects the message internally.
	//
	// This port method exists so the NATS consumer file can remain pure
	// wiring (NewConsumer + lambda delegation) — the JSON unmarshal +
	// switch dispatcher lives in the application service, not the consumer
	// file.
	//
	// Parameters:
	//   - msg: The NATS message carrying a ListNameUpdatedEvent JSON payload
	HandleListNameUpdated(msg *natsModel.Message)

	// ProcessL2WriteRetry is invoked by the template L2 sync fallback
	// consumer when a previous synchronous L2 write failed. Re-fetches
	// the template from Mongo (source of truth) and re-runs
	// syncTemplateL2. On success emits the existing FANOUT invalidation.
	// The consumer NAKs on error so NATS retries with backoff.
	ProcessL2WriteRetry(ctx ctx.Context, templateId string) error

	// GetTemplateByIdForCacheFallback retrieves a template by ID and repopulates L2 cache.
	//
	// This method is used by internal API for TieredCache fallback.
	// When L2 (MinIO) cache miss occurs, consuming services call this endpoint
	// to fetch the template from MongoDB and repopulate the cache.
	//
	// Flow:
	//  1. Fetch template from MongoDB by ID
	//  2. Write scripts to MinIO (L2) to repopulate cache
	//  3. Return the full template response
	//
	// Parameters:
	//   - ctx: Context for controlling cancellation and timeouts
	//   - templateId: The template ID (MongoDB ObjectId as string)
	//
	// Returns:
	//   - AssetTemplateResponse: The full template data
	//   - error: If template not found
	GetTemplateByIdForCacheFallback(ctx ctx.Context, templateId string) (*dtos.AssetTemplateResponse, error)
}

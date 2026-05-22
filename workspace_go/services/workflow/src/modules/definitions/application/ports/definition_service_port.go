package ports

import (
	ctx "context"

	"workflow/src/modules/definitions/application/dtos"

	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

// DefinitionServicePort defines the contract for workflow definition business operations.
type DefinitionServicePort interface {
	// CreateDefinition creates a new workflow definition.
	CreateDefinition(ctx ctx.Context, requestContext *reqCtx.RequestContext, dto *dtos.DefinitionCreateDTO) (*dtos.DefinitionResponse, error)

	// GetDefinitionById retrieves a workflow definition by its unique identifier.
	GetDefinitionById(ctx ctx.Context, definitionId *string) (*dtos.DefinitionResponse, error)

	// UpdateDefinitionById updates an existing workflow definition.
	UpdateDefinitionById(ctx ctx.Context, definitionId *string, dto *dtos.DefinitionUpdateDTO) (*dtos.DefinitionResponse, error)

	// DeleteDefinitionById removes a workflow definition.
	DeleteDefinitionById(ctx ctx.Context, definitionId *string) (map[string]bool, error)

	// GetDefinitions retrieves a paginated and filtered list of workflow definitions.
	GetDefinitions(ctx ctx.Context, requestContext *reqCtx.RequestContext, query *dtos.DefinitionQueryDTO) (*model.PaginatedResult[dtos.DefinitionResponse], error)

	// CountDefinitions returns the total count of workflow definitions for the given org context.
	CountDefinitions(ctx ctx.Context, requestContext *reqCtx.RequestContext) (int64, error)

	// GetNodeScript retrieves script source for a code node (internal endpoint for TieredCache fallback).
	// Fetches from MongoDB, repopulates L2 (MinIO), and returns the script string.
	GetNodeScript(ctx ctx.Context, definitionId string, nodeId string) (string, error)
}

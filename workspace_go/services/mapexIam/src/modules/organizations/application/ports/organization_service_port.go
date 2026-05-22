package ports

import (
	ctx "context"
	"mapexIam/src/modules/organizations/application/dtos"
	"mapexIam/src/modules/organizations/domain/entities"

	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

// OrganizationServicePort defines the contract for organization business logic operations.
// This interface follows Hexagonal Architecture (Ports & Adapters) pattern,
// enabling dependency inversion and decoupling from concrete implementations.
//
// The port allows for:
//   - Easy mocking in tests
//   - Swappable implementations
//   - Clear separation between domain logic and infrastructure concerns
type OrganizationServicePort interface {
	// CreateOrganization creates a new organization with proper hierarchical PathKey and Code.
	// Handles parent-child relationships and generates unique codes based on parent's childCount.
	CreateOrganization(context ctx.Context, dto *dtos.CreateOrganizationDto) (*dtos.OrganizationResponse, error)

	// GetOrganizationById retrieves an organization entity by its unique identifier.
	GetOrganizationById(context ctx.Context, organizationId *string) (*dtos.OrganizationResponse, error)

	// UpdateOrganizationById updates an existing organization entity.
	// If AccessPolicy changes, invalidates authorization cache for all users in this organization.
	UpdateOrganizationById(context ctx.Context, organizationId *string, dto *dtos.UpdateOrganizationDto) (*dtos.OrganizationResponse, error)

	// DeleteOrganizationById removes an organization entity by its unique ID.
	DeleteOrganizationById(context ctx.Context, organizationId *string) (map[string]bool, error)

	// GetOrganizations retrieves a paginated and filtered list of Organization DTOs.
	// Uses RequestContext from coverage middleware for context-aware org filtering with hierarchical support.
	//
	// Parameters:
	//   - context: Request-scoped context
	//   - reqContext: RequestContext with ScopedOrgIds, OrgContext, and PathKey data
	//   - query: Query parameters including filters, pagination, and includeChildren flag
	//
	// Returns filtered results based on:
	//   - OrgContext + includeChildren: PathKey range query (org and descendants)
	//   - OrgContext alone: Direct orgId query (specific org only)
	//   - No context: $in query (all accessible orgs)
	GetOrganizations(context ctx.Context, reqContext *reqCtx.RequestContext, query *dtos.OrganizationQueryDto) (*model.PaginatedResult[dtos.OrganizationResponse], error)

	// GetOrganizationsTree retrieves organizations in a hierarchical tree structure
	// with cursor-based pagination for UI navigation components.
	GetOrganizationsTree(context ctx.Context, orgId *string, query *dtos.TreeQueryDto) (*dtos.TreeResponseDto, error)

	// GetChildOrganizationsByPathKey retrieves all child organizations for a given parent pathKey.
	// Uses pathKey prefix match to find all descendants in the hierarchy.
	//
	// Example: parentPathKey = "mapex" -> returns all orgs with pathKey starting with "mapex."
	//
	// Parameters:
	//   - context: Request-scoped context
	//   - parentPathKey: PathKey of the parent organization
	//
	// Returns:
	//   - []entities.Organization: List of child organizations
	//   - error: If query fails
	GetChildOrganizationsByPathKey(context ctx.Context, parentPathKey string) ([]entities.Organization, error)

}

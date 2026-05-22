package handlers

import (
	"github.com/gofiber/fiber/v2"

	"mapexIam/src/modules/organizations/application/dtos"
	"mapexIam/src/modules/organizations/application/ports"

	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/response"
)

// CreateOrganization returns a Fiber handler that creates a new organization.
//
// It expects a validated DTO of type dtos.CreateOrganizationDto in the "bodyDTO" context key
// (populated by requestValidation middleware).
//
// Organizations support hierarchical multi-tenant architecture with PathKey, CustomerId, Scope, and OrgId.
//
// Returns:
//   - 201 Created with organization data
//   - 400 Bad Request if validation fails
//   - 500 Internal Server Error on service failure
func CreateOrganization(service ports.OrganizationServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// retrieve the timeout‐aware Context you set in ContextInjector
		ctx := c.UserContext()

		bodyData, _ := requestValidation.GetDTO[*dtos.CreateOrganizationDto](c, "bodyDTO")

		// Transform DTO (validation, normalization)
		if err := bodyData.Transform(); err != nil {
			return response.BadRequest(c, []string{err.Error()})
		}

		retData, err := service.CreateOrganization(ctx, bodyData)

		if err != nil {
			return err
		} else {
			return response.Created(c, retData)
		}
	}
}

// GetOrganizationById returns a Fiber handler that retrieves an organization by its unique identifier.
//
// It expects a validated DTO of type dtos.OrganizationIdDto in the "paramsDTO" context key
// (populated by requestValidation middleware).
//
// Returns:
//   - 200 OK with organization data if found
//   - 404 Not Found if organization doesn't exist
//   - 500 Internal Server Error on service failure
func GetOrganizationById(service ports.OrganizationServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// retrieve the timeout‐aware Context you set in ContextInjector
		ctx := c.UserContext()

		params, _ := requestValidation.GetDTO[*dtos.OrganizationIdDto](c, "paramsDTO")
		retData, err := service.GetOrganizationById(ctx, &params.OrganizationId)

		if err != nil {
			return err
		} else {
			return response.Success(c, retData)
		}
	}
}

// UpdateOrganizationById returns a Fiber handler that updates an existing organization's information.
//
// It expects validated DTOs:
//   - dtos.OrganizationIdDto in "paramsDTO" for the organization identifier
//   - dtos.UpdateOrganizationDto in "bodyDTO" for the fields to update
//
// Both DTOs are populated by requestValidation middleware.
//
// Returns:
//   - 200 OK with updated organization data
//   - 404 Not Found if organization doesn't exist
//   - 400 Bad Request if validation fails
//   - 500 Internal Server Error on service failure
func UpdateOrganizationById(service ports.OrganizationServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// retrieve the timeout‐aware Context you set in ContextInjector
		ctx := c.UserContext()

		params, _ := requestValidation.GetDTO[*dtos.OrganizationIdDto](c, "paramsDTO")
		bodyData, _ := requestValidation.GetDTO[*dtos.UpdateOrganizationDto](c, "bodyDTO")

		// Transform DTO (validation, normalization)
		if err := bodyData.Transform(); err != nil {
			return response.BadRequest(c, []string{err.Error()})
		}

		retData, err := service.UpdateOrganizationById(ctx, &params.OrganizationId, bodyData)

		if err != nil {
			return err
		} else {
			return response.Created(c, retData)
		}
	}
}

// DeleteOrganizationById returns a Fiber handler that permanently deletes an organization.
//
// It expects a validated DTO of type dtos.OrganizationIdDto in the "paramsDTO" context key
// (populated by requestValidation middleware).
//
// Returns:
//   - 200 OK with success flag if deletion succeeds
//   - 404 Not Found if organization doesn't exist
//   - 500 Internal Server Error on service failure
func DeleteOrganizationById(service ports.OrganizationServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// retrieve the timeout‐aware Context you set in ContextInjector
		ctx := c.UserContext()

		params, _ := requestValidation.GetDTO[*dtos.OrganizationIdDto](c, "paramsDTO")
		retData, err := service.DeleteOrganizationById(ctx, &params.OrganizationId)

		if err != nil {
			return err
		} else {
			return response.Success(c, retData)
		}
	}
}

// GetOrganizations returns a Fiber handler that retrieves a paginated and filtered list of organizations.
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
// It expects a validated DTO of type dtos.OrganizationQueryDto in the "queryDTO" context key
// (populated by requestValidation middleware) containing optional filters such as:
//   - type, name, enabled (filters)
//   - page, perPage (pagination)
//   - projection (field selection)
//   - includeChildren (hierarchical query flag)
//
// Returns:
//   - 200 OK with paginated organization list
//   - 400 Bad Request if query validation fails
//   - 500 Internal Server Error on service failure or requestContext not found
func GetOrganizations(service ports.OrganizationServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// retrieve the timeout‐aware Context you set in ContextInjector
		ctx := c.UserContext()

		// Get RequestContext from coverage middleware
		requestContext, ok := c.Locals("requestContext").(*reqCtx.RequestContext)
		if !ok {
			return response.InternalServerError(c, "requestContext not found in request context", nil)
		}

		queryData, _ := requestValidation.GetDTO[*dtos.OrganizationQueryDto](c, "queryDTO")
		retData, err := service.GetOrganizations(ctx, requestContext, queryData)

		if err != nil {
			return err
		} else {
			return response.Success(c, retData)
		}
	}
}

// GetOrganizationsTree returns a Fiber handler that retrieves organizations in a tree structure
// with cursor-based pagination for hierarchical navigation in UI components.
// Uses X-Org-Context header to determine the root organization for the tree.
func GetOrganizationsTree(service ports.OrganizationServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {

		ctx := c.UserContext()

		queryData, _ := requestValidation.GetDTO[*dtos.TreeQueryDto](c, "queryDTO")
		orgId := c.Get("X-Org-Context") // Get org context from new header

		retData, err := service.GetOrganizationsTree(ctx, &orgId, queryData)

		if err != nil {
			return err
		} else {
			return response.Success(c, retData)
		}
	}
}

package handlers

import (
	"github.com/gofiber/fiber/v2"

	"triggers/src/modules/triggers/application/dtos"
	"triggers/src/modules/triggers/application/ports"

	contractsCommon "github.com/Mapex-Solutions/MapexOS/contracts/common"
	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/response"
)

// CreateTrigger returns a Fiber handler that creates a new trigger.
//
// Following Hexagonal Architecture, this handler accepts the service port interface
// and delegates business logic to the service layer.
//
// It uses RequestContext (injected by coverage middleware) which contains:
//   - OrgContext: The selected organization ID from X-Org-Context header
//   - OrgContextData: Organization data including PathKey for hierarchical filtering
//   - UserContext: The user ID from JWT token
//
// The handler passes the full RequestContext to the service layer, which extracts
// the needed fields (orgId, pathKey, userId) for multi-tenant support.
//
// IMPORTANT: orgId and pathKey are populated by the service from RequestContext,
// NOT from the DTO! This prevents client manipulation.
//
// Parameters:
//   - service: The TriggerServicePort interface for trigger business operations
//
// Returns:
//   - A Fiber handler function that processes the trigger creation request
func CreateTrigger(service ports.TriggerServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		// Get RequestContext from coverage middleware
		requestContext, ok := c.Locals("requestContext").(*reqCtx.RequestContext)
		if !ok {
			return response.InternalServerError(c, "requestContext not found in request context", nil)
		}

		// Get validated DTO from request validation middleware
		bodyData, _ := requestValidation.GetDTO[*dtos.CreateTriggerDto](c, "bodyDTO")

		// Pass requestContext to service (contains OrgContext, OrgContextData, UserContext)
		retData, err := service.CreateTrigger(ctx, requestContext, bodyData)

		if err != nil {
			return err
		}
		return response.Created(c, retData)
	}
}

// GetTriggerById returns a Fiber handler that retrieves a trigger by its ID.
//
// Following Hexagonal Architecture, this handler accepts the service port interface
// and delegates business logic to the service layer.
//
// Parameters:
//   - service: The TriggerServicePort interface for trigger business operations
//
// Returns:
//   - A Fiber handler function that processes the trigger retrieval request
func GetTriggerById(service ports.TriggerServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		// Get trigger ID from URL params
		triggerId := c.Params("id")

		retData, err := service.GetTriggerById(ctx, &triggerId)

		if err != nil {
			return err
		}
		return response.Success(c, retData)
	}
}

// UpdateTriggerById returns a Fiber handler that updates a trigger by its ID.
//
// Following Hexagonal Architecture, this handler accepts the service port interface
// and delegates business logic to the service layer.
//
// Parameters:
//   - service: The TriggerServicePort interface for trigger business operations
//
// Returns:
//   - A Fiber handler function that processes the trigger update request
func UpdateTriggerById(service ports.TriggerServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		// Get RequestContext for updatedBy field
		requestContext, ok := c.Locals("requestContext").(*reqCtx.RequestContext)
		if !ok {
			return response.InternalServerError(c, "requestContext not found in request context", nil)
		}

		// Get trigger ID from URL params
		triggerId := c.Params("id")

		// Get validated DTO from request validation middleware
		bodyData, _ := requestValidation.GetDTO[*dtos.UpdateTriggerDto](c, "bodyDTO")

		retData, err := service.UpdateTriggerById(ctx, requestContext, &triggerId, bodyData)

		if err != nil {
			return err
		}
		return response.Success(c, retData)
	}
}

// DeleteTriggerById returns a Fiber handler that deletes a trigger by its ID.
//
// Following Hexagonal Architecture, this handler accepts the service port interface
// and delegates business logic to the service layer.
//
// Parameters:
//   - service: The TriggerServicePort interface for trigger business operations
//
// Returns:
//   - A Fiber handler function that processes the trigger deletion request
func DeleteTriggerById(service ports.TriggerServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		// Get trigger ID from URL params
		triggerId := c.Params("id")

		retData, err := service.DeleteTriggerById(ctx, &triggerId)

		if err != nil {
			return err
		}
		return response.Success(c, retData)
	}
}

// GetTriggers returns a Fiber handler that retrieves a paginated list of triggers.
//
// Following Hexagonal Architecture, this handler accepts the service port interface
// and delegates business logic to the service layer.
//
// It uses RequestContext from coverage middleware for org filtering with hierarchical support.
//
// Parameters:
//   - service: The TriggerServicePort interface for trigger business operations
//
// Returns:
//   - A Fiber handler function that processes the trigger list request
func GetTriggers(service ports.TriggerServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		// Get RequestContext from coverage middleware
		requestContext, ok := c.Locals("requestContext").(*reqCtx.RequestContext)
		if !ok {
			return response.InternalServerError(c, "requestContext not found in request context", nil)
		}

		// Get query parameters from request validation middleware
		queryData, _ := requestValidation.GetDTO[*dtos.TriggerQueryDto](c, "queryDTO")

		retData, err := service.GetTriggers(ctx, requestContext, queryData)

		if err != nil {
			return err
		}
		return response.Success(c, retData)
	}
}

// GetTriggerCount returns a Fiber handler that returns the total count of triggers.
// Uses cached count with 6h TTL, invalidated on create/delete.
//
// Parameters:
//   - service: The TriggerServicePort interface for trigger business operations
//
// Returns:
//   - A Fiber handler function that processes the trigger count request
func GetTriggerCount(service ports.TriggerServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		requestContext, ok := c.Locals("requestContext").(*reqCtx.RequestContext)
		if !ok {
			return response.InternalServerError(c, "requestContext not found in request context", nil)
		}

		count, err := service.CountTriggers(ctx, requestContext)
		if err != nil {
			return err
		}

		return response.Success(c, contractsCommon.CounterResponse{Count: count})
	}
}

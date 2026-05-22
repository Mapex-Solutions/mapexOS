package handlers

import (
	"github.com/gofiber/fiber/v2"

	"events/src/modules/retention/application/dtos"
	"events/src/modules/retention/application/ports"

	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/response"
)

// GetRetentionPolicies returns a Fiber handler that retrieves a paginated and filtered list of retention policies.
// Uses RequestContext from coverage middleware for context-aware org filtering.
//
// Parameters:
//   - service: The RetentionServicePort interface for retention policy business operations
//
// Returns:
//   - A Fiber handler function that processes the retention policy list request
func GetRetentionPolicies(service ports.RetentionServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		// Get RequestContext from coverage middleware
		requestContext, ok := c.Locals("requestContext").(*reqCtx.RequestContext)
		if !ok {
			return response.InternalServerError(c, "requestContext not found in request context", nil)
		}

		queryData, dtoErr := requestValidation.GetDTO[*dtos.RetentionPolicyQueryDTO](c, "queryDTO")
		if dtoErr != nil {
			return response.BadRequest(c, []string{dtoErr.Error()})
		}
		retData, err := service.GetRetentionPolicies(ctx, requestContext, queryData)

		if err != nil {
			return err
		}
		return response.Success(c, retData)
	}
}

// GetRetentionPolicyById returns a Fiber handler that retrieves a retention policy by its ID.
//
// Parameters:
//   - service: The RetentionServicePort interface for retention policy business operations
//
// Returns:
//   - A Fiber handler function that processes the retention policy retrieval request
func GetRetentionPolicyById(service ports.RetentionServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		params, dtoErr := requestValidation.GetDTO[*dtos.RetentionPolicyParamsDTO](c, "paramsDTO")
		if dtoErr != nil {
			return response.BadRequest(c, []string{dtoErr.Error()})
		}
		retData, err := service.GetRetentionPolicyById(ctx, &params.RetentionPolicyId)

		if err != nil {
			return err
		}
		return response.Success(c, retData)
	}
}

// UpsertRetentionPolicy returns a Fiber handler that creates or updates a retention policy.
// Uses PUT with upsert semantics on orgId + type.
//
// Parameters:
//   - service: The RetentionServicePort interface for retention policy business operations
//
// Returns:
//   - A Fiber handler function that processes the retention policy upsert request
func UpsertRetentionPolicy(service ports.RetentionServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		// Get RequestContext from coverage middleware
		requestContext, ok := c.Locals("requestContext").(*reqCtx.RequestContext)
		if !ok {
			return response.InternalServerError(c, "requestContext not found in request context", nil)
		}

		bodyData, dtoErr := requestValidation.GetDTO[*dtos.RetentionPolicyUpsertDTO](c, "bodyDTO")
		if dtoErr != nil {
			return response.BadRequest(c, []string{dtoErr.Error()})
		}

		retData, err := service.UpsertRetentionPolicy(ctx, requestContext, bodyData)
		if err != nil {
			return err
		}
		return response.Success(c, retData)
	}
}

// DeleteRetentionPolicyById returns a Fiber handler that deletes a retention policy by its ID.
//
// Parameters:
//   - service: The RetentionServicePort interface for retention policy business operations
//
// Returns:
//   - A Fiber handler function that processes the retention policy deletion request
func DeleteRetentionPolicyById(service ports.RetentionServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		params, dtoErr := requestValidation.GetDTO[*dtos.RetentionPolicyParamsDTO](c, "paramsDTO")
		if dtoErr != nil {
			return response.BadRequest(c, []string{dtoErr.Error()})
		}
		retData, err := service.DeleteRetentionPolicyById(ctx, &params.RetentionPolicyId)

		if err != nil {
			return err
		}
		return response.Success(c, retData)
	}
}

package handlers

import (
	"github.com/gofiber/fiber/v2"

	"events/src/modules/events/application/dtos"
	"events/src/modules/events/application/ports"

	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/response"
)

// GetEventsRaw returns a Fiber handler that retrieves raw events using cursor-based pagination.
// Uses RequestContext from coverage middleware for context-aware org filtering.
//
// The handler extracts RequestContext from c.Locals() which was set by the InjectRequestContext middleware.
// This ensures users can only query events within their accessible organizations with zero extra queries.
//
// Cursor pagination is used instead of offset pagination because:
//   - Raw events can have millions of records
//   - COUNT queries are expensive on large ClickHouse tables
//   - Cursor pagination uses timestamp index for efficient seeks
//
// Query parameters:
//   - cursor: timestamp to start from (RFC3339 format, optional)
//   - direction: "next" (older items) or "prev" (newer items), default: "next"
//   - limit: max items to return (default: 20, max: 100)
//   - sortAsc: false = DESC (newest first), true = ASC (oldest first)
//   - threadId: filter by thread ID / data source ID (optional)
//   - source: filter by source (http_gateway, mqtt_gateway, etc.)
//   - startTime: filter events after this timestamp (RFC3339 format)
//   - endTime: filter events before this timestamp (RFC3339 format)
//   - includeChildren: include child orgs hierarchically (default: false)
//
// Returns:
//   - 200 OK with cursor-paginated events data (items, nextCursor, prevCursor, hasNext, hasPrevious)
//   - 400 Bad Request if validation fails
//   - 403 Forbidden if user doesn't have access to org
//   - 500 Internal Server Error on service failure
func GetEventsRaw(service ports.EventServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// retrieve the timeout-aware Context you set in ContextInjector
		ctx := c.UserContext()

		// Get RequestContext from coverage middleware
		requestContext, ok := c.Locals("requestContext").(*reqCtx.RequestContext)
		if !ok {
			return fiber.NewError(fiber.StatusForbidden, "Request context not found")
		}

		// Get validated query DTO
		queryData, _ := requestValidation.GetDTO[*dtos.EventsRawQueryDto](c, "queryDTO")

		// Call service with RequestContext for automatic org filtering
		result, err := service.GetEventsRaw(ctx, requestContext, queryData)

		if err != nil {
			return err
		} else {
			return response.Success(c, result)
		}
	}
}

// GetEventsJsExec returns a Fiber handler that retrieves JS Executor events using cursor-based pagination.
// Uses RequestContext from coverage middleware for context-aware org filtering.
//
// The handler extracts RequestContext from c.Locals() which was set by the InjectRequestContext middleware.
// This ensures users can only query events within their accessible organizations with zero extra queries.
//
// Cursor pagination is used instead of offset pagination because:
//   - JS exec events can have millions of records
//   - COUNT queries are expensive on large ClickHouse tables
//   - Cursor pagination uses timestamp index for efficient seeks
//
// Query parameters:
//   - cursor: timestamp to start from (RFC3339 format, optional)
//   - direction: "next" (older items) or "prev" (newer items), default: "next"
//   - limit: max items to return (default: 20, max: 100)
//   - sortAsc: false = DESC (newest first), true = ASC (oldest first)
//   - assetUuid: filter by asset UUID (optional)
//   - assetId: filter by asset MongoDB ID (optional)
//   - success: filter by execution success status (optional)
//   - startTime: filter events after this timestamp (RFC3339 format)
//   - endTime: filter events before this timestamp (RFC3339 format)
//   - includeChildren: include child orgs hierarchically (default: false)
//
// Returns:
//   - 200 OK with cursor-paginated events data (items, nextCursor, prevCursor, hasNext, hasPrevious)
//   - 400 Bad Request if validation fails
//   - 403 Forbidden if user doesn't have access to org
//   - 500 Internal Server Error on service failure
func GetEventsJsExec(service ports.EventServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// retrieve the timeout-aware Context you set in ContextInjector
		ctx := c.UserContext()

		// Get RequestContext from coverage middleware
		requestContext, ok := c.Locals("requestContext").(*reqCtx.RequestContext)
		if !ok {
			return fiber.NewError(fiber.StatusForbidden, "Request context not found")
		}

		// Get validated query DTO
		queryData, _ := requestValidation.GetDTO[*dtos.EventsJsExecQueryDto](c, "queryDTO")

		// Call service with RequestContext for automatic org filtering
		result, err := service.GetEventsJsExec(ctx, requestContext, queryData)

		if err != nil {
			return err
		} else {
			return response.Success(c, result)
		}
	}
}

// GetEventsRouter returns a Fiber handler that retrieves router events using cursor-based pagination.
// Uses RequestContext from coverage middleware for context-aware org filtering.
//
// The handler extracts RequestContext from c.Locals() which was set by the InjectRequestContext middleware.
// This ensures users can only query events within their accessible organizations with zero extra queries.
//
// Cursor pagination is used instead of offset pagination because:
//   - Router events can have millions of records
//   - COUNT queries are expensive on large ClickHouse tables
//   - Cursor pagination uses timestamp index for efficient seeks
//
// Query parameters:
//   - cursor: timestamp to start from (RFC3339 format, optional)
//   - direction: "next" (older items) or "prev" (newer items), default: "next"
//   - limit: max items to return (default: 20, max: 100)
//   - sortAsc: false = DESC (newest first), true = ASC (oldest first)
//   - threadId: filter by thread ID (optional)
//   - assetId: filter by asset ID (optional)
//   - routerId: filter by router/RouteGroup ID (optional)
//   - success: filter by success status (optional)
//   - publishedCount: filter by published count (optional)
//   - startTime: filter events after this timestamp (RFC3339 format)
//   - endTime: filter events before this timestamp (RFC3339 format)
//   - includeChildren: include child orgs hierarchically (default: false)
//
// Returns:
//   - 200 OK with cursor-paginated events data (items, nextCursor, prevCursor, hasNext, hasPrevious)
//   - 400 Bad Request if validation fails
//   - 403 Forbidden if user doesn't have access to org
//   - 500 Internal Server Error on service failure
func GetEventsRouter(service ports.EventServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// retrieve the timeout-aware Context you set in ContextInjector
		ctx := c.UserContext()

		// Get RequestContext from coverage middleware
		requestContext, ok := c.Locals("requestContext").(*reqCtx.RequestContext)
		if !ok {
			return fiber.NewError(fiber.StatusForbidden, "Request context not found")
		}

		// Get validated query DTO
		queryData, _ := requestValidation.GetDTO[*dtos.EventsRouterQueryDto](c, "queryDTO")

		// Call service with RequestContext for automatic org filtering
		result, err := service.GetEventsRouter(ctx, requestContext, queryData)

		if err != nil {
			return err
		} else {
			return response.Success(c, result)
		}
	}
}

// GetEventsBusinessRule returns a Fiber handler that retrieves business rule events using cursor-based pagination.
// Uses RequestContext from coverage middleware for context-aware org filtering.
//
// The handler extracts RequestContext from c.Locals() which was set by the InjectRequestContext middleware.
// This ensures users can only query events within their accessible organizations with zero extra queries.
//
// Cursor pagination is used instead of offset pagination because:
//   - Business rule events can have millions of records
//   - COUNT queries are expensive on large ClickHouse tables
//   - Cursor pagination uses timestamp index for efficient seeks
//
// Query parameters:
//   - cursor: timestamp to start from (RFC3339 format, optional)
//   - direction: "next" (older items) or "prev" (newer items), default: "next"
//   - limit: max items to return (default: 20, max: 100)
//   - sortAsc: false = DESC (newest first), true = ASC (oldest first)
//   - threadId: filter by thread ID (optional)
//   - ruleId: filter by rule template ID (optional)
//   - businessRuleId: filter by business rule ID (optional)
//   - matched: filter by matched status (optional)
//   - startTime: filter events after this timestamp (RFC3339 format)
//   - endTime: filter events before this timestamp (RFC3339 format)
//   - includeChildren: include child orgs hierarchically (default: false)
//
// Returns:
//   - 200 OK with cursor-paginated events data (items, nextCursor, prevCursor, hasNext, hasPrevious)
//   - 400 Bad Request if validation fails
//   - 403 Forbidden if user doesn't have access to org
//   - 500 Internal Server Error on service failure
func GetEventsBusinessRule(service ports.EventServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// retrieve the timeout-aware Context you set in ContextInjector
		ctx := c.UserContext()

		// Get RequestContext from coverage middleware
		requestContext, ok := c.Locals("requestContext").(*reqCtx.RequestContext)
		if !ok {
			return fiber.NewError(fiber.StatusForbidden, "Request context not found")
		}

		// Get validated query DTO
		queryData, _ := requestValidation.GetDTO[*dtos.EventsBusinessRuleQueryDto](c, "queryDTO")

		// Call service with RequestContext for automatic org filtering
		result, err := service.GetEventsBusinessRule(ctx, requestContext, queryData)

		if err != nil {
			return err
		} else {
			return response.Success(c, result)
		}
	}
}

// GetEventsTrigger returns a Fiber handler that retrieves trigger events using cursor-based pagination.
// Uses RequestContext from coverage middleware for context-aware org filtering.
//
// The handler extracts RequestContext from c.Locals() which was set by the InjectRequestContext middleware.
// This ensures users can only query events within their accessible organizations with zero extra queries.
//
// Cursor pagination is used instead of offset pagination because:
//   - Trigger events can have millions of records
//   - COUNT queries are expensive on large ClickHouse tables
//   - Cursor pagination uses timestamp index for efficient seeks
//
// Query parameters:
//   - cursor: timestamp to start from (RFC3339 format, optional)
//   - direction: "next" (older items) or "prev" (newer items), default: "next"
//   - limit: max items to return (default: 20, max: 100)
//   - sortAsc: false = DESC (newest first), true = ASC (oldest first)
//   - triggerId: filter by trigger ID (optional)
//   - triggerType: filter by trigger type (http, mqtt, email, etc.)
//   - category: filter by category (technical, communication)
//   - source: filter by source (router)
//   - success: filter by success status (optional)
//   - startTime: filter events after this timestamp (RFC3339 format)
//   - endTime: filter events before this timestamp (RFC3339 format)
//   - includeChildren: include child orgs hierarchically (default: false)
//
// Returns:
//   - 200 OK with cursor-paginated events data (items, nextCursor, prevCursor, hasNext, hasPrevious)
//   - 400 Bad Request if validation fails
//   - 403 Forbidden if user doesn't have access to org
//   - 500 Internal Server Error on service failure
func GetEventsTrigger(service ports.EventServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// retrieve the timeout-aware Context you set in ContextInjector
		ctx := c.UserContext()

		// Get RequestContext from coverage middleware
		requestContext, ok := c.Locals("requestContext").(*reqCtx.RequestContext)
		if !ok {
			return fiber.NewError(fiber.StatusForbidden, "Request context not found")
		}

		// Get validated query DTO
		queryData, _ := requestValidation.GetDTO[*dtos.EventsTriggerQueryDto](c, "queryDTO")

		// Call service with RequestContext for automatic org filtering
		result, err := service.GetEventsTrigger(ctx, requestContext, queryData)

		if err != nil {
			return err
		} else {
			return response.Success(c, result)
		}
	}
}

// GetEventsWorkflow returns a Fiber handler that retrieves workflow execution events using cursor-based pagination.
// Uses RequestContext from coverage middleware for context-aware org filtering.
//
// Query parameters:
//   - cursor: timestamp to start from (RFC3339 format, optional)
//   - direction: "next" (older items) or "prev" (newer items), default: "next"
//   - limit: max items to return (default: 20, max: 100)
//   - sortAsc: false = DESC (newest first), true = ASC (oldest first)
//   - eventTrackerId: filter by event tracker ID (optional)
//   - workflowUUID: filter by workflow UUID (optional)
//   - instanceId: filter by instance ID (optional)
//   - definitionId: filter by definition ID (optional)
//   - status: filter by terminal status (completed, failed, cancelled) (optional)
//   - success: filter by success status (optional)
//   - startTime: filter events after this timestamp (RFC3339 format)
//   - endTime: filter events before this timestamp (RFC3339 format)
//   - includeChildren: include child orgs hierarchically (default: false)
//
// Returns:
//   - 200 OK with cursor-paginated events data
//   - 400 Bad Request if validation fails
//   - 403 Forbidden if user doesn't have access to org
//   - 500 Internal Server Error on service failure
func GetEventsWorkflow(service ports.EventServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		requestContext, ok := c.Locals("requestContext").(*reqCtx.RequestContext)
		if !ok {
			return fiber.NewError(fiber.StatusForbidden, "Request context not found")
		}

		queryData, _ := requestValidation.GetDTO[*dtos.EventsWorkflowQueryDto](c, "queryDTO")

		result, err := service.GetEventsWorkflow(ctx, requestContext, queryData)
		if err != nil {
			return err
		}
		return response.Success(c, result)
	}
}

// GetWorkflowEventByExecutionId returns a Fiber handler that retrieves a single workflow event by executionId.
func GetWorkflowEventByExecutionId(service ports.EventServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		requestContext, ok := c.Locals("requestContext").(*reqCtx.RequestContext)
		if !ok {
			return fiber.NewError(fiber.StatusForbidden, "Request context not found")
		}

		paramsData, _ := requestValidation.GetDTO[*dtos.EventsWorkflowExecutionIdParamDto](c, "paramsDTO")
		result, err := service.GetWorkflowEventByExecutionId(ctx, requestContext, paramsData.ExecutionId)
		if err != nil {
			return err
		}
		return response.Success(c, result)
	}
}

// GetEventStoreDetail returns a Fiber handler that retrieves a single event by eventTrackerId.
// Resolves EVA fieldIds to field names using the appropriate template (based on source).
// Returns the event with an advancedSearch map containing resolved field names.
func GetEventStoreDetail(service ports.EventServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		eventTrackerId := c.Params("eventTrackerId")
		if eventTrackerId == "" {
			return fiber.NewError(fiber.StatusBadRequest, "eventTrackerId is required")
		}

		result, err := service.GetEventStoreDetail(ctx, eventTrackerId)
		if err != nil {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}

		return response.Success(c, result)
	}
}

// GetEventsStore returns a Fiber handler that retrieves processed events using cursor-based pagination.
// Uses POST with JSON body to support EVA dynamic field filters (EvaFilters array).
// Uses RequestContext from coverage middleware for context-aware org filtering.
//
// The handler extracts RequestContext from c.Locals() which was set by the InjectRequestContext middleware.
// This ensures users can only query events within their accessible organizations with zero extra queries.
//
// Cursor pagination is used instead of offset pagination because:
//   - Processed events can have millions of records
//   - COUNT queries are expensive on large ClickHouse tables
//   - Cursor pagination uses created timestamp index for efficient seeks
//
// Request body (JSON):
//   - cursor: timestamp to start from (RFC3339 format, optional)
//   - direction: "next" (older items) or "prev" (newer items), default: "next"
//   - limit: max items to return (default: 20, max: 100)
//   - sortAsc: false = DESC (newest first), true = ASC (oldest first)
//   - threadId: filter by thread ID for distributed tracing (optional)
//   - assetId: filter by asset ID (optional)
//   - templateId: filter by asset template ID (optional)
//   - eventType: filter by event type (telemetry, alarm, command) (optional)
//   - source: filter by source service (optional)
//   - startTime: filter events after this timestamp (RFC3339 format)
//   - endTime: filter events before this timestamp (RFC3339 format)
//   - includeChildren: include child orgs hierarchically (default: false)
//   - evaFilters: array of EVA dynamic field filters with operators (optional)
//
// Returns:
//   - 200 OK with cursor-paginated events data (items, nextCursor, prevCursor, hasNext, hasPrevious)
//   - 400 Bad Request if validation fails
//   - 403 Forbidden if user doesn't have access to org
//   - 500 Internal Server Error on service failure
func GetEventsStore(service ports.EventServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// retrieve the timeout-aware Context you set in ContextInjector
		ctx := c.UserContext()

		// Get RequestContext from coverage middleware
		requestContext, ok := c.Locals("requestContext").(*reqCtx.RequestContext)
		if !ok {
			return fiber.NewError(fiber.StatusForbidden, "Request context not found")
		}

		// Get validated body DTO (POST with EvaFilters support)
		queryData, _ := requestValidation.GetDTO[*dtos.EventsStoreQueryDto](c, "bodyDTO")

		// Call service with RequestContext for automatic org filtering
		result, err := service.GetEventsStore(ctx, requestContext, queryData)

		if err != nil {
			return err
		} else {
			return response.Success(c, result)
		}
	}
}

// GetEventsDLQCounts returns a Fiber handler that retrieves DLQ entry counts grouped by service type.
func GetEventsDLQCounts(service ports.EventServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		requestContext, ok := c.Locals("requestContext").(*reqCtx.RequestContext)
		if !ok {
			return fiber.NewError(fiber.StatusForbidden, "Request context not found")
		}

		queryData, _ := requestValidation.GetDTO[*dtos.EventsDLQCountsQueryDto](c, "queryDTO")

		result, err := service.GetEventsDLQCounts(ctx, requestContext, queryData)
		if err != nil {
			return err
		}
		return response.Success(c, result)
	}
}

// GetEventsDLQ returns a Fiber handler that retrieves DLQ events using cursor-based pagination.
// Uses RequestContext from coverage middleware for context-aware org filtering.
func GetEventsDLQ(service ports.EventServicePort) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		requestContext, ok := c.Locals("requestContext").(*reqCtx.RequestContext)
		if !ok {
			return fiber.NewError(fiber.StatusForbidden, "Request context not found")
		}

		queryData, _ := requestValidation.GetDTO[*dtos.EventsDLQQueryDto](c, "queryDTO")

		result, err := service.GetEventsDLQ(ctx, requestContext, queryData)
		if err != nil {
			return err
		}
		return response.Success(c, result)
	}
}

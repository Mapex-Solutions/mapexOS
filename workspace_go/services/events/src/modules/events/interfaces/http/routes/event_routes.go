package routes

import (
	"github.com/gofiber/fiber/v2"

	"events/src/modules/events/application/dtos"
	"events/src/modules/events/application/ports"
	"events/src/modules/events/interfaces/http/handlers"

	coverageMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/coverage"
	permissionMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/permission"
	validation "github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
	perms "github.com/Mapex-Solutions/MapexOS/permissions/events"
)

func RegisterRoutes(group fiber.Router, service ports.EventServicePort) {

	/**
	 * Events Raw Routes - ClickHouse raw events storage
	 */

	// Get raw events with filters, pagination, and projection
	// Uses InjectRequestContext middleware for context-aware org filtering
	eventsRawQueryDto := validation.NewValidation(nil, &dtos.EventsRawQueryDto{}, nil)
	group.Get("/raw",
		validation.ValidationMiddleware(eventsRawQueryDto),
		permissionMw.RequirePermission(perms.EventsRawList),
		coverageMw.InjectRequestContext(),
		handlers.GetEventsRaw(service),
	)

	/**
	 * Events JS Executor Routes - ClickHouse JS executor debug events
	 */

	// Get JS executor events with filters, cursor pagination
	// Uses InjectRequestContext middleware for context-aware org filtering
	eventsJsExecQueryDto := validation.NewValidation(nil, &dtos.EventsJsExecQueryDto{}, nil)
	group.Get("/jsexec",
		validation.ValidationMiddleware(eventsJsExecQueryDto),
		permissionMw.RequirePermission(perms.EventsJsExecutorList),
		coverageMw.InjectRequestContext(),
		handlers.GetEventsJsExec(service),
	)

	/**
	 * Events Router Routes - ClickHouse router execution history events
	 */

	// Get router events with filters, cursor pagination
	// Uses InjectRequestContext middleware for context-aware org filtering
	eventsRouterQueryDto := validation.NewValidation(nil, &dtos.EventsRouterQueryDto{}, nil)
	group.Get("/router",
		validation.ValidationMiddleware(eventsRouterQueryDto),
		permissionMw.RequirePermission(perms.EventsRouterList),
		coverageMw.InjectRequestContext(),
		handlers.GetEventsRouter(service),
	)

	/**
	 * Events Business Rule Routes - ClickHouse business rule execution history events
	 */

	// Get business rule events with filters, cursor pagination
	// Uses InjectRequestContext middleware for context-aware org filtering
	eventsBusinessRuleQueryDto := validation.NewValidation(nil, &dtos.EventsBusinessRuleQueryDto{}, nil)
	group.Get("/businessrule",
		validation.ValidationMiddleware(eventsBusinessRuleQueryDto),
		permissionMw.RequirePermission(perms.EventsBusinessRuleList),
		coverageMw.InjectRequestContext(),
		handlers.GetEventsBusinessRule(service),
	)

	/**
	 * Events Trigger Routes - ClickHouse trigger execution history events
	 */

	// Get trigger events with filters, cursor pagination
	// Uses InjectRequestContext middleware for context-aware org filtering
	eventsTriggerQueryDto := validation.NewValidation(nil, &dtos.EventsTriggerQueryDto{}, nil)
	group.Get("/trigger",
		validation.ValidationMiddleware(eventsTriggerQueryDto),
		permissionMw.RequirePermission(perms.EventsTriggerList),
		coverageMw.InjectRequestContext(),
		handlers.GetEventsTrigger(service),
	)

	/**
	 * Events Workflow Routes - ClickHouse workflow execution history events
	 */

	// Get workflow events with filters, cursor pagination
	// Uses InjectRequestContext middleware for context-aware org filtering
	eventsWorkflowQueryDto := validation.NewValidation(nil, &dtos.EventsWorkflowQueryDto{}, nil)
	group.Get("/workflow",
		validation.ValidationMiddleware(eventsWorkflowQueryDto),
		permissionMw.RequirePermission(perms.EventsWorkflowList),
		coverageMw.InjectRequestContext(),
		handlers.GetEventsWorkflow(service),
	)

	// Get single workflow event by executionId (MongoDB _id hex)
	eventsWorkflowExecIdParam := validation.NewValidation(nil, nil, &dtos.EventsWorkflowExecutionIdParamDto{})
	group.Get("/workflow/execution/:executionId",
		validation.ValidationMiddleware(eventsWorkflowExecIdParam),
		permissionMw.RequirePermission(perms.EventsWorkflowList),
		coverageMw.InjectRequestContext(),
		handlers.GetWorkflowEventByExecutionId(service),
	)

	// Events DLQ Routes - Dead Letter Queue events from all services
	eventsDLQCountsQueryDto := validation.NewValidation(nil, &dtos.EventsDLQCountsQueryDto{}, nil)
	group.Get("/dlq/counts",
		validation.ValidationMiddleware(eventsDLQCountsQueryDto),
		permissionMw.RequirePermission(perms.EventsDLQList),
		coverageMw.InjectRequestContext(),
		handlers.GetEventsDLQCounts(service),
	)

	eventsDLQQueryDto := validation.NewValidation(nil, &dtos.EventsDLQQueryDto{}, nil)
	group.Get("/dlq",
		validation.ValidationMiddleware(eventsDLQQueryDto),
		permissionMw.RequirePermission(perms.EventsDLQList),
		coverageMw.InjectRequestContext(),
		handlers.GetEventsDLQ(service),
	)

	/**
	 * Events Store Routes - ClickHouse processed events with EVA fields
	 */

	// Query processed events with optional EVA dynamic field filters
	// Uses POST to support EvaFilters array in request body
	eventsStoreQueryDto := validation.NewValidation(&dtos.EventsStoreQueryDto{}, nil, nil)
	group.Post("/store/query",
		validation.ValidationMiddleware(eventsStoreQueryDto),
		permissionMw.RequirePermission(perms.EventsProcessedList),
		coverageMw.InjectRequestContext(),
		handlers.GetEventsStore(service),
	)

	// Get single event detail with resolved EVA field names (advancedSearch)
	// Resolves fieldIds to field names based on source: "asset"→AssetTemplate, "rule"→BusinessRule
	group.Get("/store/:eventTrackerId",
		permissionMw.RequirePermission(perms.EventsProcessedRead),
		handlers.GetEventStoreDetail(service),
	)
}

package routes

import (
	"events/src/modules/events/application/dtos"
	"events/src/modules/events/application/ports"
	"events/src/modules/events/interfaces/http/handlers"

	perms "github.com/Mapex-Solutions/MapexOS/permissions/events"
	coverageMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/coverage"
	permissionMw "github.com/Mapex-Solutions/mapexGoKit/microservices/http/middlewares/permission"
	validation "github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/swagger"
	web "github.com/Mapex-Solutions/mapexGoKit/microservices/http/web"
)

// RegisterRoutes registers the events read-side HTTP routes (ClickHouse history
// queries). Base path: /api/v1/events.
//
// Routes are registered through the swagger wrapper: each NewValidation declares
// the input contract once (used to both validate and document), the module tag is
// declared once on Wrap, and each route's summary, description, and response type
// are attached via the fluent builder.
func RegisterRoutes(group web.Router, service ports.EventServicePort) {

	r := swagger.Wrap(group).Tag("Events")

	// Raw events (ClickHouse raw-events storage), cursor-paginated.
	eventsRawQueryDto := validation.NewValidation(nil, &dtos.EventsRawQueryDto{}, nil)
	r.Get("/raw", eventsRawQueryDto, swagger.Expose,
		permissionMw.RequirePermission(perms.EventsRawList),
		coverageMw.InjectRequestContext(),
		handlers.GetEventsRaw(service),
	).
		Summary("List raw events").
		Description("Returns cursor-paginated raw events from ClickHouse, filtered and scoped to the caller's organization.").
		Returns(&dtos.EventsRawCursorResultDto{})

	// JS-executor debug events, cursor-paginated.
	eventsJsExecQueryDto := validation.NewValidation(nil, &dtos.EventsJsExecQueryDto{}, nil)
	r.Get("/jsexec", eventsJsExecQueryDto, swagger.Expose,
		permissionMw.RequirePermission(perms.EventsJsExecutorList),
		coverageMw.InjectRequestContext(),
		handlers.GetEventsJsExec(service),
	).
		Summary("List JS-executor events").
		Description("Returns cursor-paginated JS-executor debug events (decode/validate/transform script runs) from ClickHouse.").
		Returns(&dtos.EventsJsExecCursorResultDto{})

	// Router execution history events, cursor-paginated.
	eventsRouterQueryDto := validation.NewValidation(nil, &dtos.EventsRouterQueryDto{}, nil)
	r.Get("/router", eventsRouterQueryDto, swagger.Expose,
		permissionMw.RequirePermission(perms.EventsRouterList),
		coverageMw.InjectRequestContext(),
		handlers.GetEventsRouter(service),
	).
		Summary("List router events").
		Description("Returns cursor-paginated router execution-history events (RouteGroup match + fan-out outcomes) from ClickHouse.").
		Returns(&dtos.EventsRouterCursorResultDto{})

	// Business-rule execution history events, cursor-paginated.
	eventsBusinessRuleQueryDto := validation.NewValidation(nil, &dtos.EventsBusinessRuleQueryDto{}, nil)
	r.Get("/businessrule", eventsBusinessRuleQueryDto, swagger.Expose,
		permissionMw.RequirePermission(perms.EventsBusinessRuleList),
		coverageMw.InjectRequestContext(),
		handlers.GetEventsBusinessRule(service),
	).
		Summary("List business-rule events").
		Description("Returns cursor-paginated business-rule execution-history events from ClickHouse.").
		Returns(&dtos.EventsBusinessRuleCursorResultDto{})

	// Trigger execution history events, cursor-paginated.
	eventsTriggerQueryDto := validation.NewValidation(nil, &dtos.EventsTriggerQueryDto{}, nil)
	r.Get("/trigger", eventsTriggerQueryDto, swagger.Expose,
		permissionMw.RequirePermission(perms.EventsTriggerList),
		coverageMw.InjectRequestContext(),
		handlers.GetEventsTrigger(service),
	).
		Summary("List trigger events").
		Description("Returns cursor-paginated trigger execution-history events from ClickHouse.").
		Returns(&dtos.EventsTriggerCursorResultDto{})

	// Workflow execution history events, cursor-paginated.
	eventsWorkflowQueryDto := validation.NewValidation(nil, &dtos.EventsWorkflowQueryDto{}, nil)
	r.Get("/workflow", eventsWorkflowQueryDto, swagger.Expose,
		permissionMw.RequirePermission(perms.EventsWorkflowList),
		coverageMw.InjectRequestContext(),
		handlers.GetEventsWorkflow(service),
	).
		Summary("List workflow events").
		Description("Returns cursor-paginated workflow execution-history events from ClickHouse.").
		Returns(&dtos.EventsWorkflowCursorResultDto{})

	// Single workflow event by executionId (MongoDB _id hex).
	eventsWorkflowExecIdParam := validation.NewValidation(nil, nil, &dtos.EventsWorkflowExecutionIdParamDto{})
	r.Get("/workflow/execution/:executionId", eventsWorkflowExecIdParam, swagger.Expose,
		permissionMw.RequirePermission(perms.EventsWorkflowList),
		coverageMw.InjectRequestContext(),
		handlers.GetWorkflowEventByExecutionId(service),
	).
		Summary("Get workflow event by execution ID").
		Description("Returns a single workflow execution-history event by its executionId (MongoDB ObjectId hex).").
		Returns(&dtos.EventsWorkflowResponseDto{})

	// DLQ entry counts grouped by service type.
	eventsDLQCountsQueryDto := validation.NewValidation(nil, &dtos.EventsDLQCountsQueryDto{}, nil)
	r.Get("/dlq/counts", eventsDLQCountsQueryDto, swagger.Expose,
		permissionMw.RequirePermission(perms.EventsDLQList),
		coverageMw.InjectRequestContext(),
		handlers.GetEventsDLQCounts(service),
	).
		Summary("Count DLQ entries").
		Description("Returns dead-letter-queue entry counts grouped by service type, for the caller's organization.").
		Returns(&dtos.EventsDLQCountsResultDto{})

	// DLQ events from all services, cursor-paginated.
	eventsDLQQueryDto := validation.NewValidation(nil, &dtos.EventsDLQQueryDto{}, nil)
	r.Get("/dlq", eventsDLQQueryDto, swagger.Expose,
		permissionMw.RequirePermission(perms.EventsDLQList),
		coverageMw.InjectRequestContext(),
		handlers.GetEventsDLQ(service),
	).
		Summary("List DLQ events").
		Description("Returns cursor-paginated dead-letter-queue events from all services, scoped to the caller's organization.").
		Returns(&dtos.EventsDLQCursorResultDto{})

	// Processed events with EVA dynamic-field filters. POST to carry the
	// EvaFilters array in the request body.
	eventsStoreQueryDto := validation.NewValidation(&dtos.EventsStoreQueryDto{}, nil, nil)
	r.Post("/store/query", eventsStoreQueryDto, swagger.Expose,
		permissionMw.RequirePermission(perms.EventsProcessedList),
		coverageMw.InjectRequestContext(),
		handlers.GetEventsStore(service),
	).
		Summary("Query processed events").
		Description("Returns cursor-paginated processed events (with resolved EVA fields) from ClickHouse, supporting an EvaFilters array of dynamic-field conditions in the request body.").
		Returns(&dtos.EventsStoreCursorResultDto{})

	// Single processed-event detail with resolved EVA field names.
	eventsStoreDetailDto := validation.NewValidation(nil, nil, nil)
	r.Get("/store/:eventTrackerId", eventsStoreDetailDto, swagger.Expose,
		permissionMw.RequirePermission(perms.EventsProcessedRead),
		handlers.GetEventStoreDetail(service),
	).
		Summary("Get processed-event detail").
		Description("Returns a single processed event by eventTrackerId, with EVA fieldIds resolved to human-readable field names based on their source (asset template or business rule).").
		Returns(&dtos.EventsStoreDetailResponseDto{})
}

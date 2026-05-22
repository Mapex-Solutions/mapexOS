package ports

import (
	ctx "context"

	"events/src/modules/events/application/dtos"

	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
)

// EventServicePort defines the contract for event storage operations.
//
// This port interface enables Hexagonal Architecture by decoupling the business
// logic from its implementation. NATS consumers depend on this interface
// rather than the concrete service implementation.
//
// The Events service responsibility is to persist events to ClickHouse for
// analytics, reporting, and historical queries.
type EventServicePort interface {
	// ProcessEvent handles the processing and storage of processed events.
	// This method is called by the NATS consumer for each message received.
	//
	// Flow:
	//   1. Parse the incoming event data (JSON)
	//   2. Extract required fields (assetId, orgId, pathKey, etc.)
	//   3. Store the event in ClickHouse via EventRepository
	//
	// Parameters:
	//   - data: The message payload containing event data
	//   - index: Position in the batch (for logging/debugging)
	//   - headers: Message headers containing metadata
	//
	// Returns:
	//   - nil: Event stored successfully (will be ACKed by NATS)
	//   - error: Storage failed (will be NAKed and redelivered)
	ProcessEvent(data []byte, index int, headers map[string][]string) error

	// ProcessRawEventBatch processes a complete batch of raw events at once.
	// This is the PRIMARY method for raw events - receives all messages from NATS
	// batch and performs a single bulk insert to ClickHouse.
	//
	// Flow:
	//   1. Receive all messages from NATS batch
	//   2. Parse and validate each message (Reject invalid ones to DLQ)
	//   3. Map to RawEvent entities
	//   4. Bulk insert to ClickHouse (SaveRawEventBatch)
	//   5. Handle all messages: Ack (success), Nack (DB error), Reject (invalid)
	//
	// The service handles ALL message lifecycle decisions:
	//   - msg.Reject(reason): Invalid JSON/validation errors → DLQ immediately
	//   - msg.Nack(err): Bulk insert failed → retry with backoff
	//   - msg.Ack(): Successfully processed → removed from queue
	//
	// Parameters:
	//   - messages: Slice of Message pointers with retry-aware Ack/Nack/Reject methods
	//
	// Returns:
	//   - Always nil (service handles all messages internally)
	ProcessRawEventBatch(messages []*natsModel.Message) error

	// GetEventsRaw retrieves raw events using cursor-based pagination.
	// This method is optimized for large datasets as it avoids COUNT queries.
	// Uses RequestContext for context-aware organization filtering.
	//
	// Parameters:
	//   - c: Request-scoped context for timeout and cancellation
	//   - requestContext: RequestContext from InjectRequestContext middleware
	//   - query: Query DTO containing filters and cursor pagination options
	//
	// Returns:
	//   - EventsRawCursorResultDto containing events and cursor metadata
	//   - error: Error if query fails
	GetEventsRaw(c ctx.Context, requestContext *reqCtx.RequestContext, query *dtos.EventsRawQueryDto) (*dtos.EventsRawCursorResultDto, error)

	/**
	 * JS Executor Events
	 */

	// ProcessJsExecEventBatch processes a complete batch of JS Executor events at once.
	// This is the PRIMARY method for JS exec events - receives all messages from NATS
	// batch and performs a single bulk insert to ClickHouse.
	//
	// Flow:
	//   1. Receive all messages from NATS batch
	//   2. Parse and validate each message (Reject invalid ones to DLQ)
	//   3. Map to JsExecEvent entities using mapper.DtoToEntity
	//   4. Bulk insert to ClickHouse (SaveJsExecEventBatch)
	//   5. Handle all messages: Ack (success), Nack (DB error), Reject (invalid)
	//
	// The service handles ALL message lifecycle decisions:
	//   - msg.Reject(reason): Invalid JSON/validation errors → DLQ immediately
	//   - msg.Nack(err): Bulk insert failed → retry with backoff
	//   - msg.Ack(): Successfully processed → removed from queue
	//
	// Parameters:
	//   - messages: Slice of Message pointers with retry-aware Ack/Nack/Reject methods
	//
	// Returns:
	//   - Always nil (service handles all messages internally)
	ProcessJsExecEventBatch(messages []*natsModel.Message) error

	// GetEventsJsExec retrieves JS Executor events using cursor-based pagination.
	// This method is optimized for large datasets as it avoids COUNT queries.
	// Uses RequestContext for context-aware organization filtering.
	//
	// Parameters:
	//   - c: Request-scoped context for timeout and cancellation
	//   - requestContext: RequestContext from InjectRequestContext middleware
	//   - query: Query DTO containing filters and cursor pagination options
	//
	// Returns:
	//   - EventsJsExecCursorResultDto containing events and cursor metadata
	//   - error: Error if query fails
	GetEventsJsExec(c ctx.Context, requestContext *reqCtx.RequestContext, query *dtos.EventsJsExecQueryDto) (*dtos.EventsJsExecCursorResultDto, error)

	/**
	 * Dead Letter Queue (DLQ) Events
	 */

	// ProcessDLQEventBatch processes a complete batch of DLQ events at once.
	// This method receives messages from the MAPEXOS-DLQ stream and stores them
	// in ClickHouse for debugging and analysis.
	//
	// Note: This consumer does NOT use retry/DLQ policy itself (no infinite loop).
	// If storage fails, messages are simply ACKed to avoid redelivery loops.
	//
	// Flow:
	//   1. Receive all messages from NATS batch
	//   2. Parse and map to DLQEvent entities
	//   3. Bulk insert to ClickHouse
	//   4. ACK all messages (even on failure to avoid loops)
	//
	// Parameters:
	//   - messages: Slice of Message pointers from NATS
	//
	// Returns:
	//   - Always nil (service handles all messages internally)
	ProcessDLQEventBatch(messages []*natsModel.Message) error

	// GetEventsDLQ retrieves DLQ events using cursor-based pagination.
	// This method is optimized for large datasets as it avoids COUNT queries.
	// Uses RequestContext for context-aware organization filtering.
	//
	// Parameters:
	//   - c: Request-scoped context for timeout and cancellation
	//   - requestContext: RequestContext from InjectRequestContext middleware
	//   - query: Query DTO containing filters and cursor pagination options
	//
	// Returns:
	//   - EventsDLQCursorResultDto containing events and cursor metadata
	//   - error: Error if query fails
	GetEventsDLQ(c ctx.Context, requestContext *reqCtx.RequestContext, query *dtos.EventsDLQQueryDto) (*dtos.EventsDLQCursorResultDto, error)

	// GetEventsDLQCounts retrieves DLQ entry counts grouped by service type.
	GetEventsDLQCounts(c ctx.Context, requestContext *reqCtx.RequestContext, query *dtos.EventsDLQCountsQueryDto) (*dtos.EventsDLQCountsResultDto, error)

	/**
	 * Router Events
	 */

	// ProcessRouterEventBatch processes a complete batch of router events at once.
	// This is the PRIMARY method for router events - receives all messages from NATS
	// batch and performs a single bulk insert to ClickHouse.
	//
	// Flow:
	//   1. Receive all messages from NATS batch
	//   2. Parse and validate each message (Reject invalid ones to DLQ)
	//   3. Map to RouterEvent entities
	//   4. Bulk insert to ClickHouse (SaveRouterEventBatch)
	//   5. Handle all messages: Ack (success), Nack (DB error), Reject (invalid)
	//
	// Parameters:
	//   - messages: Slice of Message pointers with retry-aware Ack/Nack/Reject methods
	//
	// Returns:
	//   - Always nil (service handles all messages internally)
	ProcessRouterEventBatch(messages []*natsModel.Message) error

	// GetEventsRouter retrieves router events using cursor-based pagination.
	// This method is optimized for large datasets as it avoids COUNT queries.
	// Uses RequestContext for context-aware organization filtering.
	//
	// Parameters:
	//   - c: Request-scoped context for timeout and cancellation
	//   - requestContext: RequestContext from InjectRequestContext middleware
	//   - query: Query DTO containing filters and cursor pagination options
	//
	// Returns:
	//   - EventsRouterCursorResultDto containing events and cursor metadata
	//   - error: Error if query fails
	GetEventsRouter(c ctx.Context, requestContext *reqCtx.RequestContext, query *dtos.EventsRouterQueryDto) (*dtos.EventsRouterCursorResultDto, error)

	/**
	 * Business Rule Events
	 */

	// ProcessBusinessRuleEventBatch processes a complete batch of business rule events at once.
	// This is the PRIMARY method for business rule events - receives all messages from NATS
	// batch and performs a single bulk insert to ClickHouse.
	//
	// Flow:
	//   1. Receive all messages from NATS batch
	//   2. Parse and validate each message (Reject invalid ones to DLQ)
	//   3. Map to BusinessRuleEvent entities
	//   4. Bulk insert to ClickHouse (SaveBusinessRuleEventBatch)
	//   5. Handle all messages: Ack (success), Nack (DB error), Reject (invalid)
	//
	// Parameters:
	//   - messages: Slice of Message pointers with retry-aware Ack/Nack/Reject methods
	//
	// Returns:
	//   - Always nil (service handles all messages internally)
	ProcessBusinessRuleEventBatch(messages []*natsModel.Message) error

	// GetEventsBusinessRule retrieves business rule events using cursor-based pagination.
	// This method is optimized for large datasets as it avoids COUNT queries.
	// Uses RequestContext for context-aware organization filtering.
	//
	// Parameters:
	//   - c: Request-scoped context for timeout and cancellation
	//   - requestContext: RequestContext from InjectRequestContext middleware
	//   - query: Query DTO containing filters and cursor pagination options
	//
	// Returns:
	//   - EventsBusinessRuleCursorResultDto containing events and cursor metadata
	//   - error: Error if query fails
	GetEventsBusinessRule(c ctx.Context, requestContext *reqCtx.RequestContext, query *dtos.EventsBusinessRuleQueryDto) (*dtos.EventsBusinessRuleCursorResultDto, error)

	/**
	 * Trigger Events
	 */

	// ProcessTriggerEventBatch processes a complete batch of trigger events at once.
	// This is the PRIMARY method for trigger events - receives all messages from NATS
	// batch and performs a single bulk insert to ClickHouse.
	//
	// Flow:
	//   1. Receive all messages from NATS batch
	//   2. Parse and validate each message (Reject invalid ones to DLQ)
	//   3. Map to TriggerEvent entities
	//   4. Bulk insert to ClickHouse (SaveTriggerEventBatch)
	//   5. Handle all messages: Ack (success), Nack (DB error), Reject (invalid)
	//
	// Parameters:
	//   - messages: Slice of Message pointers with retry-aware Ack/Nack/Reject methods
	//
	// Returns:
	//   - Always nil (service handles all messages internally)
	ProcessTriggerEventBatch(messages []*natsModel.Message) error

	// GetEventsTrigger retrieves trigger events using cursor-based pagination.
	// This method is optimized for large datasets as it avoids COUNT queries.
	// Uses RequestContext for context-aware organization filtering.
	//
	// Parameters:
	//   - c: Request-scoped context for timeout and cancellation
	//   - requestContext: RequestContext from InjectRequestContext middleware
	//   - query: Query DTO containing filters and cursor pagination options
	//
	// Returns:
	//   - EventsTriggerCursorResultDto containing events and cursor metadata
	//   - error: Error if query fails
	GetEventsTrigger(c ctx.Context, requestContext *reqCtx.RequestContext, query *dtos.EventsTriggerQueryDto) (*dtos.EventsTriggerCursorResultDto, error)

	/**
	 * Workflow Events
	 */

	// ProcessWorkflowEventBatch processes a complete batch of workflow execution events at once.
	// Flow:
	//   1. Receive all messages from NATS batch
	//   2. Parse and validate each message (Reject invalid ones to DLQ)
	//   3. Map to WorkflowEvent entities
	//   4. Bulk insert to ClickHouse (SaveWorkflowEventBatch)
	//   5. Handle all messages: Ack (success), Nack (DB error), Reject (invalid)
	//
	// Parameters:
	//   - messages: Slice of Message pointers with retry-aware Ack/Nack/Reject methods
	//
	// Returns:
	//   - Always nil (service handles all messages internally)
	ProcessWorkflowEventBatch(messages []*natsModel.Message) error

	// GetEventsWorkflow retrieves workflow events using cursor-based pagination.
	// Uses RequestContext for context-aware organization filtering.
	//
	// Parameters:
	//   - c: Request-scoped context for timeout and cancellation
	//   - requestContext: RequestContext from InjectRequestContext middleware
	//   - query: Query DTO containing filters and cursor pagination options
	//
	// Returns:
	//   - EventsWorkflowCursorResultDto containing events and cursor metadata
	//   - error: Error if query fails
	GetEventsWorkflow(c ctx.Context, requestContext *reqCtx.RequestContext, query *dtos.EventsWorkflowQueryDto) (*dtos.EventsWorkflowCursorResultDto, error)

	// GetWorkflowEventByExecutionId retrieves a single workflow event by executionId (MongoDB _id hex).
	GetWorkflowEventByExecutionId(c ctx.Context, requestContext *reqCtx.RequestContext, executionId string) (*dtos.EventsWorkflowResponseDto, error)

	/**
	 * Event Store (Processed Events with EVA)
	 */

	// ProcessEventStoreBatch processes a batch of processed events with EVA field mapping.
	// This is the PRIMARY method for storing processed events - receives all messages from NATS
	// batch, resolves EVA field mappings via template cache, and performs bulk insert to ClickHouse.
	//
	// EVA (Entity-Value-Attribute) Mapping Flow:
	//   1. Parse event with assetTemplateId
	//   2. Fetch template from cache (L0/L1/Fallback HTTP)
	//   3. Extract DynamicFields with fieldId mappings
	//   4. Map event data fields to eva_number, eva_string, eva_bool, eva_date MAPs
	//   5. Bulk insert to ClickHouse events table
	//
	// The service handles ALL message lifecycle decisions:
	//   - msg.Reject(reason): Invalid JSON/validation errors → DLQ immediately
	//   - msg.Nack(err): Bulk insert failed → retry with backoff
	//   - msg.Ack(): Successfully processed → removed from queue
	//
	// Parameters:
	//   - messages: Slice of Message pointers with retry-aware Ack/Nack/Reject methods
	//
	// Returns:
	//   - Always nil (service handles all messages internally)
	ProcessEventStoreBatch(messages []*natsModel.Message) error

	// GetEventsStore retrieves processed events using cursor-based pagination.
	// List view: no EVA fields, just core data + payload.
	//
	// Parameters:
	//   - c: Request-scoped context for timeout and cancellation
	//   - requestContext: RequestContext from InjectRequestContext middleware
	//   - query: Query DTO containing filters and cursor pagination options
	//
	// Returns:
	//   - EventsStoreCursorResultDto containing events and cursor metadata
	//   - error: Error if query fails
	GetEventsStore(c ctx.Context, requestContext *reqCtx.RequestContext, query *dtos.EventsStoreQueryDto) (*dtos.EventsStoreCursorResultDto, error)

	// GetEventStoreDetail retrieves a single event by eventTrackerId and resolves EVA fields.
	// Checks source to determine where to fetch DynamicFields:
	//   - "asset" → AssetTemplate via TieredCache
	//   - "rule"  → BusinessRule (future)
	//
	// Returns the event with advancedSearch map containing resolved field names.
	GetEventStoreDetail(c ctx.Context, eventTrackerId string) (*dtos.EventsStoreDetailResponseDto, error)

	// HandleTemplateInvalidate processes a single FANOUT message from
	// mapexos.fanout.template.invalidate. The kit's SubscribeFanout callback
	// auto-Acks; bad payloads are logged and dropped (no Reject/Nack).
	HandleTemplateInvalidate(msg *natsModel.Message)
}

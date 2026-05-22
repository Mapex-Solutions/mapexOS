import { z, StringAndNotBeEmptyOrOptional, IsBoolean, IsStringDateFormat, NumberIntAndPositive, IsString } from '@mapexos/validations';

/**
 * Events Raw Query schema - Used for cursor-based pagination of raw events
 *
 * Based on Go contract: EventsRawQuery (contracts/services/events/events/dtos.go)
 *
 * Cursor pagination fields (from CursorQueryDTO):
 *   - cursor: timestamp to start from (RFC3339 format)
 *   - direction: "next" (older items) or "prev" (newer items)
 *   - limit: max items to return (default: 20, max: 100)
 *   - sortAsc: false = DESC (newest first), true = ASC (oldest first)
 *   - includeChildren: include child orgs hierarchically (default: false)
 *
 * Module-specific filters:
 *   - threadId: filter by thread ID (data source identifier)
 *   - source: filter by source (http_gateway, mqtt_gateway, etc.)
 *   - startTime: filter events after this timestamp (ISO 8601)
 *   - endTime: filter events before this timestamp (ISO 8601)
 *
 * Why cursor pagination?
 *   - Raw events can have millions of records
 *   - Offset pagination requires COUNT which is expensive on large tables
 *   - Cursor pagination uses timestamp index directly (fast seeks)
 */
export const ZodEventsRawQuerySchema = z.object({
	// CursorQueryDTO fields
	cursor: IsStringDateFormat.optional(),
	direction: z.enum(['next', 'prev']).optional(),
	limit: NumberIntAndPositive.max(100).optional(),
	sortAsc: IsBoolean.optional(),
	includeChildren: IsBoolean.optional(),

	// Module-specific filters
	eventTrackerId: IsString.optional(),
	threadId: IsString.optional(),
	source: IsString.optional(),
	success: IsBoolean.optional(),
	startTime: IsStringDateFormat.optional(),
	endTime: IsStringDateFormat.optional(),
});

/**
 * Events Raw Response schema - Single raw event item
 *
 * Based on Go contract: EventsRawResponse (contracts/services/events/events/dtos.go)
 */
export const ZodEventsRawResponseSchema = z.object({
	created: IsStringDateFormat,
	eventTrackerId: IsString.optional(),
	threadId: IsString,
	orgId: IsString,
	source: IsString,
	name: IsString.optional(),
	description: IsString.optional(),
	event: z.record(IsString, z.any()),
	metadata: z.record(IsString, z.any()).optional(),
	success: IsBoolean,
	error: IsString.optional(),
	retentionDays: z.number().int().min(0).optional(),
});

/**
 * Events Raw Cursor Result schema - Cursor-paginated response
 *
 * Based on Go contract: EventsRawCursorResult (contracts/services/events/events/dtos.go)
 */
export const ZodEventsRawCursorResultSchema = z.object({
	/** List of raw events */
	items: z.array(ZodEventsRawResponseSchema),

	/** Timestamp to use for fetching the next page (older items in DESC order) */
	nextCursor: IsStringDateFormat.optional().nullable(),

	/** Timestamp to use for fetching the previous page (newer items in DESC order) */
	prevCursor: IsStringDateFormat.optional().nullable(),

	/** Indicates if there are more items after the current page */
	hasNext: IsBoolean,

	/** Indicates if there are more items before the current page */
	hasPrevious: IsBoolean,
});

// ============================================================================
// JS Executor Events Schemas
// ============================================================================

/**
 * Execution time filter operator enum
 * Valid values: lt (<), lte (<=), gt (>), gte (>=), between
 */
export const ExecTimeOperatorEnum = z.enum(['lt', 'lte', 'gt', 'gte', 'between']);

/**
 * Events JS Executor Query schema - Used for cursor-based pagination of JS executor events
 *
 * Based on Go contract: EventsJsExecQuery (contracts/services/events/events/events_jsexecutor.dto.go)
 *
 * Cursor pagination fields (from CursorQueryDTO):
 *   - cursor: timestamp to start from (RFC3339 format)
 *   - direction: "next" (older items) or "prev" (newer items)
 *   - limit: max items to return (default: 20, max: 100)
 *   - sortAsc: false = DESC (newest first), true = ASC (oldest first)
 *   - includeChildren: include child orgs hierarchically (default: false)
 *
 * Module-specific filters:
 *   - threadId: filter by thread ID (asset UUID)
 *   - success: filter by execution success status
 *   - startTime: filter events after this timestamp (ISO 8601)
 *   - endTime: filter events before this timestamp (ISO 8601)
 *   - execTimeOp: operator for execution time filter (lt, lte, gt, gte, between)
 *   - execTimeValue: execution time value in milliseconds
 *   - execTimeValueEnd: end value for "between" operator (milliseconds)
 */
export const ZodEventsJsExecQuerySchema = z.object({
	// CursorQueryDTO fields
	cursor: IsStringDateFormat.optional(),
	direction: z.enum(['next', 'prev']).optional(),
	limit: NumberIntAndPositive.max(100).optional(),
	sortAsc: IsBoolean.optional(),
	includeChildren: IsBoolean.optional(),

	// Module-specific filters
	eventTrackerId: IsString.optional(),
	threadId: IsString.optional(),
	success: IsBoolean.optional(),
	startTime: IsStringDateFormat.optional(),
	endTime: IsStringDateFormat.optional(),

	// Execution time filter (in milliseconds)
	execTimeOp: ExecTimeOperatorEnum.optional(),
	execTimeValue: NumberIntAndPositive.optional(),
	execTimeValueEnd: NumberIntAndPositive.optional(),
});

/**
 * Events JS Executor Response schema - Single JS executor event item
 *
 * Based on Go contract: EventsJsExecResponse (contracts/services/events/events/events_jsexecutor.dto.go)
 *
 * Field names standardized with events_raw pattern:
 *   - threadId: asset UUID (or 'n/a' if not available)
 *   - name: data source name
 *   - description: data source description
 *   - event: standardized payload (JSON string)
 */
export const ZodEventsJsExecResponseSchema = z.object({
	created: IsStringDateFormat,
	eventTrackerId: IsString.optional(),
	threadId: IsString,
	orgId: IsString,
	pathKey: IsString.optional(),
	name: IsString.optional(),
	description: IsString.optional(),
	event: IsString.optional(),
	success: IsBoolean,
	failedAt: IsString.optional(),
	totalExecutionTime: z.number().int().min(0),
	error: IsString.optional(),
	retentionDays: z.number().int().min(0).optional(),
});

/**
 * Events JS Executor Cursor Result schema - Cursor-paginated response
 *
 * Based on Go contract: EventsJsExecCursorResult (contracts/services/events/events/events_jsexecutor.dto.go)
 */
export const ZodEventsJsExecCursorResultSchema = z.object({
	/** List of JS executor events */
	items: z.array(ZodEventsJsExecResponseSchema),

	/** Timestamp to use for fetching the next page (older items in DESC order) */
	nextCursor: IsStringDateFormat.optional().nullable(),

	/** Timestamp to use for fetching the previous page (newer items in DESC order) */
	prevCursor: IsStringDateFormat.optional().nullable(),

	/** Indicates if there are more items after the current page */
	hasNext: IsBoolean,

	/** Indicates if there are more items before the current page */
	hasPrevious: IsBoolean,
});

// ============================================================================
// Router Events Schemas
// ============================================================================

/**
 * Events Router Query schema - Used for cursor-based pagination of router events
 *
 * Based on Go contract: EventsRouterQuery (contracts/services/events/events/events_router.dto.go)
 *
 * Cursor pagination fields (from CursorQueryDTO):
 *   - cursor: timestamp to start from (RFC3339 format)
 *   - direction: "next" (older items) or "prev" (newer items)
 *   - limit: max items to return (default: 20, max: 100)
 *   - sortAsc: false = DESC (newest first), true = ASC (oldest first)
 *   - includeChildren: include child orgs hierarchically (default: false)
 *
 * Module-specific filters:
 *   - threadId: filter by thread ID (asset ID)
 *   - assetId: filter by asset ID
 *   - routerId: filter by router/RouteGroup ID
 *   - success: filter by success status
 *   - publishedCount: filter by published count
 *   - startTime: filter events after this timestamp (ISO 8601)
 *   - endTime: filter events before this timestamp (ISO 8601)
 */
export const ZodEventsRouterQuerySchema = z.object({
	// CursorQueryDTO fields
	cursor: IsStringDateFormat.optional(),
	direction: z.enum(['next', 'prev']).optional(),
	limit: NumberIntAndPositive.max(100).optional(),
	sortAsc: IsBoolean.optional(),
	includeChildren: IsBoolean.optional(),

	// Module-specific filters
	eventTrackerId: IsString.optional(),
	threadId: IsString.optional(),
	assetId: IsString.optional(),
	routerId: IsString.optional(),
	success: IsBoolean.optional(),
	publishedCount: NumberIntAndPositive.optional(),
	startTime: IsStringDateFormat.optional(),
	endTime: IsStringDateFormat.optional(),
});

/**
 * Events Router Response schema - Single router event item
 *
 * Based on Go contract: EventsRouterResponse (contracts/services/events/events/events_router.dto.go)
 *
 * Fields:
 *   - threadId: asset ID used for grouping events
 *   - assetId: asset that triggered the routing
 *   - routerId: RouteGroup ID
 *   - name: RouteGroup name
 *   - description: RouteGroup description
 *   - totalRouters: number of routers in the RouteGroup
 *   - matchedCount: number of routers that matched conditions
 *   - publishedCount: number of routers that published events
 *   - event: routers array as JSON string
 *   - success: true if at least one router published
 */
export const ZodEventsRouterResponseSchema = z.object({
	created: IsStringDateFormat,
	eventTrackerId: IsString.optional(),
	threadId: IsString,
	orgId: IsString,
	pathKey: IsString.optional(),
	assetId: IsString,
	routerId: IsString,
	name: IsString.optional(),
	description: IsString.optional(),
	totalRouters: z.number().int().min(0),
	matchedCount: z.number().int().min(0),
	publishedCount: z.number().int().min(0),
	event: IsString,
	success: IsBoolean,
	error: IsString.optional(),
	retentionDays: z.number().int().min(0).optional(),
});

/**
 * Events Router Cursor Result schema - Cursor-paginated response
 *
 * Based on Go contract: EventsRouterCursorResult (contracts/services/events/events/events_router.dto.go)
 */
export const ZodEventsRouterCursorResultSchema = z.object({
	/** List of router events */
	items: z.array(ZodEventsRouterResponseSchema),

	/** Timestamp to use for fetching the next page (older items in DESC order) */
	nextCursor: IsStringDateFormat.optional().nullable(),

	/** Timestamp to use for fetching the previous page (newer items in DESC order) */
	prevCursor: IsStringDateFormat.optional().nullable(),

	/** Indicates if there are more items after the current page */
	hasNext: IsBoolean,

	/** Indicates if there are more items before the current page */
	hasPrevious: IsBoolean,
});

// ============================================================================
// Business Rule Events Schemas
// ============================================================================

/**
 * Events Business Rule Query schema - Used for cursor-based pagination of business rule events
 *
 * Based on Go contract: EventsBusinessRuleQuery (contracts/services/events/events/events_businessrule.dto.go)
 *
 * Cursor pagination fields (from CursorQueryDTO):
 *   - cursor: timestamp to start from (RFC3339 format)
 *   - direction: "next" (older items) or "prev" (newer items)
 *   - limit: max items to return (default: 20, max: 100)
 *   - sortAsc: false = DESC (newest first), true = ASC (oldest first)
 *   - includeChildren: include child orgs hierarchically (default: false)
 *
 * Module-specific filters:
 *   - threadId: filter by thread ID
 *   - ruleId: filter by rule template ID
 *   - businessRuleId: filter by business rule ID
 *   - matched: filter by matched status
 *   - startTime: filter events after this timestamp (ISO 8601)
 *   - endTime: filter events before this timestamp (ISO 8601)
 */
export const ZodEventsBusinessRuleQuerySchema = z.object({
	// CursorQueryDTO fields
	cursor: IsStringDateFormat.optional(),
	direction: z.enum(['next', 'prev']).optional(),
	limit: NumberIntAndPositive.max(100).optional(),
	sortAsc: IsBoolean.optional(),
	includeChildren: IsBoolean.optional(),

	// Module-specific filters
	eventTrackerId: IsString.optional(),
	threadId: IsString.optional(),
	ruleId: IsString.optional(),
	businessRuleId: IsString.optional(),
	matched: IsBoolean.optional(),
	startTime: IsStringDateFormat.optional(),
	endTime: IsStringDateFormat.optional(),
});

/**
 * Events Business Rule Response schema - Single business rule event item
 *
 * Based on Go contract: EventsBusinessRuleResponse (contracts/services/events/events/events_businessrule.dto.go)
 *
 * Fields:
 *   - threadId: external correlation identifier
 *   - ruleId: rule template ID
 *   - businessRuleId: business rule instance ID
 *   - businessRuleName: business rule name for display
 *   - businessRuleDescription: business rule description
 *   - matched: whether the rule matched
 *   - durationMs: execution duration in milliseconds
 *   - conditionsEvaluated/Matched: evaluation metrics
 *   - evaluationTree: hierarchical evaluation tree (JSON string)
 *   - conditionLogs: condition logs array (JSON string)
 *   - actionsToDispatch: actions to dispatch (JSON string)
 */
export const ZodEventsBusinessRuleResponseSchema = z.object({
	created: IsStringDateFormat,
	eventTrackerId: IsString.optional(),
	threadId: IsString,
	orgId: IsString,
	pathKey: IsString.optional(),
	ruleId: IsString,
	businessRuleId: IsString,
	businessRuleName: IsString,
	businessRuleDescription: IsString.optional(),
	matched: IsBoolean,
	durationMs: z.number().int().min(0),
	conditionsEvaluated: z.number().int().min(0),
	conditionsMatched: z.number().int().min(0),
	groupsEvaluated: z.number().int().min(0),
	maxDepthReached: z.number().int().min(0),
	finalState: IsString.optional(),
	stateChanges: IsString.optional(),
	evaluationTree: IsString.optional(),
	conditionLogs: IsString.optional(),
	actionsToDispatch: IsString.optional(),
	retentionDays: z.number().int().min(0).optional(),
});

/**
 * Events Business Rule Cursor Result schema - Cursor-paginated response
 *
 * Based on Go contract: EventsBusinessRuleCursorResult (contracts/services/events/events/events_businessrule.dto.go)
 */
export const ZodEventsBusinessRuleCursorResultSchema = z.object({
	/** List of business rule events */
	items: z.array(ZodEventsBusinessRuleResponseSchema),

	/** Timestamp to use for fetching the next page (older items in DESC order) */
	nextCursor: IsStringDateFormat.optional().nullable(),

	/** Timestamp to use for fetching the previous page (newer items in DESC order) */
	prevCursor: IsStringDateFormat.optional().nullable(),

	/** Indicates if there are more items after the current page */
	hasNext: IsBoolean,

	/** Indicates if there are more items before the current page */
	hasPrevious: IsBoolean,
});

// ============================================================================
// Trigger Events Schemas
// ============================================================================

/**
 * Events Trigger Query schema - Used for cursor-based pagination of trigger events
 *
 * Based on Go contract: EventsTriggerQuery (contracts/services/events/events/events_trigger.dto.go)
 *
 * Cursor pagination fields (from CursorQueryDTO):
 *   - cursor: timestamp to start from (RFC3339 format)
 *   - direction: "next" (older items) or "prev" (newer items)
 *   - limit: max items to return (default: 20, max: 100)
 *   - sortAsc: false = DESC (newest first), true = ASC (oldest first)
 *   - includeChildren: include child orgs hierarchically (default: false)
 *
 * Module-specific filters:
 *   - triggerId: filter by trigger ID
 *   - triggerType: filter by trigger type (http, mqtt, email, etc.)
 *   - category: filter by category (technical, communication)
 *   - success: filter by success status
 *   - startTime: filter events after this timestamp (ISO 8601)
 *   - endTime: filter events before this timestamp (ISO 8601)
 */
export const ZodEventsTriggerQuerySchema = z.object({
	// CursorQueryDTO fields
	cursor: IsStringDateFormat.optional(),
	direction: z.enum(['next', 'prev']).optional(),
	limit: NumberIntAndPositive.max(100).optional(),
	sortAsc: IsBoolean.optional(),
	includeChildren: IsBoolean.optional(),

	// Module-specific filters
	eventTrackerId: IsString.optional(),
	triggerId: IsString.optional(),
	triggerType: IsString.optional(),
	category: IsString.optional(),
	source: IsString.optional(),
	success: IsBoolean.optional(),
	startTime: IsStringDateFormat.optional(),
	endTime: IsStringDateFormat.optional(),
});

/**
 * Events Trigger Response schema - Single trigger event item
 *
 * Based on Go contract: EventsTriggerResponse (contracts/services/events/events/events_trigger.dto.go)
 *
 * Fields:
 *   - triggerId: trigger identifier
 *   - triggerName: trigger name for display
 *   - triggerType: trigger type (http, mqtt, rabbitmq, nats, websocket, email, teams, slack)
 *   - category: trigger category (technical, communication)
 *   - success: whether the trigger executed successfully
 *   - durationMs: execution duration in milliseconds
 *   - requestData: resolved config sent to trigger (JSON string)
 *   - responseData: response from trigger (JSON string)
 */
export const ZodEventsTriggerResponseSchema = z.object({
	created: IsStringDateFormat,
	eventTrackerId: IsString.optional(),
	orgId: IsString,
	pathKey: IsString.optional(),
	triggerId: IsString,
	triggerName: IsString,
	triggerType: IsString,
	category: IsString,
	source: IsString,
	success: IsBoolean,
	durationMs: z.number().int().min(0),
	error: IsString.optional(),
	requestData: IsString.optional(),
	responseData: IsString.optional(),
	retentionDays: z.number().int().min(0).optional(),
});

/**
 * Events Trigger Cursor Result schema - Cursor-paginated response
 *
 * Based on Go contract: EventsTriggerCursorResult (contracts/services/events/events/events_trigger.dto.go)
 */
export const ZodEventsTriggerCursorResultSchema = z.object({
	/** List of trigger events */
	items: z.array(ZodEventsTriggerResponseSchema),

	/** Timestamp to use for fetching the next page (older items in DESC order) */
	nextCursor: IsStringDateFormat.optional().nullable(),

	/** Timestamp to use for fetching the previous page (newer items in DESC order) */
	prevCursor: IsStringDateFormat.optional().nullable(),

	/** Indicates if there are more items after the current page */
	hasNext: IsBoolean,

	/** Indicates if there are more items before the current page */
	hasPrevious: IsBoolean,
});

// ============================================================================
// Event Store Schemas (Processed Events with EVA fields)
// ============================================================================

/**
 * Events Store Query schema - Used for cursor-based pagination of processed events
 *
 * Based on Go contract: EventsStoreQuery (contracts/services/events/events/events_store.dto.go)
 *
 * Cursor pagination fields (from CursorQueryDTO):
 *   - cursor: timestamp to start from (RFC3339 format)
 *   - direction: "next" (older items) or "prev" (newer items)
 *   - limit: max items to return (default: 20, max: 100)
 *   - sortAsc: false = DESC (newest first), true = ASC (oldest first)
 *   - includeChildren: include child orgs hierarchically (default: false)
 *
 * Module-specific filters:
 *   - threadId: filter by thread ID for distributed tracing
 *   - assetId: filter by asset ID
 *   - assetTemplateId: filter by asset template ID
 *   - eventType: filter by event type (telemetry, alarm, command)
 *   - source: filter by source service
 *   - startTime: filter events after this timestamp (ISO 8601)
 *   - endTime: filter events before this timestamp (ISO 8601)
 */
/**
 * EVA Filter schema - Single filter condition on an EVA MAP column field
 *
 * Based on Go contract: EvaFilter (contracts/services/events/events/events_store.dto.go)
 *
 * The frontend resolves the fieldId from the AssetTemplate.DynamicFields array
 * and sends it along with the bucket type, operator, and value(s).
 */
export const ZodEvaFilterSchema = z.object({
	/** uint16 key in the EVA MAP column (from DynamicField.fieldId property, auto-increment starting at 1) */
	fieldId: z.number().int().min(1).max(65535),

	/** EVA type bucket: "number", "string", "bool", "date" */
	bucket: z.enum(['number', 'string', 'bool', 'date']),

	/** Comparison operator */
	operator: z.enum(['eq', 'neq', 'gt', 'gte', 'lt', 'lte', 'between', 'like']),

	/** Primary comparison value (sent as string, cast per bucket type on backend) */
	value: IsString,

	/** Range end value, used only with "between" operator */
	endValue: IsString.optional(),
});

export const ZodEventsStoreQuerySchema = z.object({
	// CursorQueryDTO fields
	cursor: IsStringDateFormat.optional(),
	direction: z.enum(['next', 'prev']).optional(),
	limit: NumberIntAndPositive.max(100).optional(),
	sortAsc: IsBoolean.optional(),
	includeChildren: IsBoolean.optional(),

	// Module-specific filters
	threadId: IsString.optional(),
	assetId: IsString.optional(),
	assetTemplateId: IsString.optional(),
	eventType: IsString.optional(),
	source: IsString.optional(),
	startTime: IsStringDateFormat.optional(),
	endTime: IsStringDateFormat.optional(),

	// EVA dynamic field filters (sent via POST body as JSON array)
	evaFilters: z.array(ZodEvaFilterSchema).optional(),
});

/**
 * Events Store Response schema - Single processed event with EVA fields
 *
 * Based on Go contract: EventsStoreResponse (contracts/services/events/events/events_store.dto.go)
 *
 * EVA (Entity-Value-Attribute) fields:
 *   - evaNumber: Map<uint16, float64> for numeric values
 *   - evaString: Map<uint16, string> for text values
 *   - evaBool: Map<uint16, bool> for boolean values
 *   - evaDate: Map<uint16, time.Time> for date values
 *
 * Frontend can resolve field names using AssetTemplate.DynamicFields where
 * each field has a fieldId (uint16) that maps to these EVA maps.
 */
export const ZodEventsStoreResponseSchema = z.object({
	/** Event creation timestamp (not 'timestamp' - uses 'created' for consistency) */
	created: IsStringDateFormat,

	/** Thread ID for distributed tracing */
	threadId: IsString.optional(),

	/** Asset ID that generated this event */
	assetId: IsString,

	/** Asset Template ID used to process this event */
	assetTemplateId: IsString.optional(),

	/** Human-readable asset name (denormalized at write-time by Router) */
	assetName: IsString.optional(),

	/** Asset description (denormalized at write-time by Router) */
	assetDescription: IsString.optional(),

	/** Human-readable template name (denormalized at write-time by Router) */
	templateName: IsString.optional(),

	/** Template description (denormalized at write-time by Router) */
	templateDescription: IsString.optional(),

	/** Organization ID */
	orgId: IsString,

	/** Path key for hierarchical filtering */
	pathKey: IsString.optional(),

	/** Event type (telemetry, alarm, command) — not currently populated */
	eventType: IsString.optional(),

	/** Source service that created this event (asset, rule) */
	source: IsString.optional(),

	/** Original payload as JSON string */
	payload: IsString,

	/** Event metadata as JSON string */
	metadata: IsString.optional(),

	/** EVA Number fields - Map of fieldId to numeric value */
	evaNumber: z.record(z.string(), z.number()).optional(),

	/** EVA String fields - Map of fieldId to text value */
	evaString: z.record(z.string(), IsString).optional(),

	/** EVA Boolean fields - Map of fieldId to boolean value */
	evaBool: z.record(z.string(), IsBoolean).optional(),

	/** EVA Date fields - Map of fieldId to date value */
	evaDate: z.record(z.string(), IsStringDateFormat).optional(),

	/** Retention period in days */
	retentionDays: z.number().int().min(0).optional(),
});

/**
 * Events Store Cursor Result schema - Cursor-paginated response for processed events
 *
 * Based on Go contract: EventsStoreCursorResult (contracts/services/events/events/events_store.dto.go)
 */
export const ZodEventsStoreCursorResultSchema = z.object({
	/** List of processed events */
	items: z.array(ZodEventsStoreResponseSchema),

	/** Timestamp to use for fetching the next page (older items in DESC order) */
	nextCursor: IsStringDateFormat.optional().nullable(),

	/** Timestamp to use for fetching the previous page (newer items in DESC order) */
	prevCursor: IsStringDateFormat.optional().nullable(),

	/** Indicates if there are more items after the current page */
	hasNext: IsBoolean,

	/** Indicates if there are more items before the current page */
	hasPrevious: IsBoolean,
});

// ============================================================================
// Workflow Event Schemas (Workflow execution history from ClickHouse)
// ============================================================================

/**
 * Events Workflow Query schema - Used for cursor-based pagination of workflow execution events
 *
 * Based on Go contract: EventsWorkflowQuery (contracts/services/events/events/events_workflow.dto.go)
 *
 * Module-specific filters:
 *   - eventTrackerId: filter by event tracker ID for end-to-end tracing
 *   - workflowUUID: filter by workflow execution UUID
 *   - instanceId: filter by workflow instance ID
 *   - definitionId: filter by workflow definition ID
 *   - status: filter by terminal status (completed, failed, cancelled)
 *   - success: filter by success status
 *   - startTime, endTime: time range filters (ISO 8601)
 */
export const ZodEventsWorkflowQuerySchema = z.object({
	// CursorQueryDTO fields
	cursor: IsStringDateFormat.optional(),
	direction: z.enum(['next', 'prev']).optional(),
	limit: NumberIntAndPositive.max(100).optional(),
	sortAsc: IsBoolean.optional(),
	includeChildren: IsBoolean.optional(),

	// Module-specific filters
	eventTrackerId: IsString.optional(),
	workflowUUID: IsString.optional(),
	instanceId: IsString.optional(),
	definitionId: IsString.optional(),
	status: IsString.optional(),
	success: IsBoolean.optional(),
	startTime: IsStringDateFormat.optional(),
	endTime: IsStringDateFormat.optional(),
});

/**
 * Events Workflow Execution ID Param schema - URL param for getting a single workflow event
 */
export const ZodEventsWorkflowExecutionIdParamSchema = z.object({
	executionId: IsString,
});

/**
 * Events Workflow Response schema - Single workflow execution event item
 *
 * Based on Go contract: EventsWorkflowResponse (contracts/services/events/events/events_workflow.dto.go)
 */
export const ZodEventsWorkflowResponseSchema = z.object({
	created: IsStringDateFormat,
	finished: IsStringDateFormat,
	executionId: IsString.optional(),
	eventTrackerId: IsString.optional(),
	orgId: IsString,
	pathKey: IsString.optional(),
	workflowUUID: IsString.optional(),
	instanceId: IsString,
	definitionId: IsString,
	workflowName: IsString,
	status: IsString,
	success: IsBoolean,
	durationMs: z.number().int().min(0),
	errorMessage: IsString.optional(),
	executionPath: IsString.optional(),
	nodeOutputs: IsString.optional(),
	errorInfo: IsString.optional(),
	eventPayload: IsString.optional(),
	triggerSource: IsString.optional(),
	parentExecutionId: IsString.optional(),
	depth: z.number().int().min(0),
	retentionDays: z.number().int().min(0).optional(),
});

/**
 * Events Workflow Cursor Result schema - Cursor-paginated response
 *
 * Based on Go contract: EventsWorkflowCursorResult (contracts/services/events/events/events_workflow.dto.go)
 */
export const ZodEventsWorkflowCursorResultSchema = z.object({
	/** List of workflow execution events */
	items: z.array(ZodEventsWorkflowResponseSchema),

	/** Timestamp to use for fetching the next page (older items in DESC order) */
	nextCursor: IsStringDateFormat.optional().nullable(),

	/** Timestamp to use for fetching the previous page (newer items in DESC order) */
	prevCursor: IsStringDateFormat.optional().nullable(),

	/** Indicates if there are more items after the current page */
	hasNext: IsBoolean,

	/** Indicates if there are more items before the current page */
	hasPrevious: IsBoolean,
});

/**
 * Events DLQ Query schema - Dead Letter Queue events from all services
 */
export const ZodEventsDLQQuerySchema = z.object({
	cursor: IsStringDateFormat.optional(),
	direction: z.enum(['next', 'prev']).optional(),
	limit: NumberIntAndPositive.max(100).optional(),
	sortAsc: IsBoolean.optional(),
	includeChildren: IsBoolean.optional(),
	eventTrackerId: IsString.optional(),
	serviceName: IsString.optional(),
	serviceType: IsString.optional(),
	eventType: IsString.optional(),
	lastError: IsString.optional(),
	startTime: IsStringDateFormat.optional(),
	endTime: IsStringDateFormat.optional(),
});

export const ZodEventsDLQResponseSchema = z.object({
	created: IsString,
	eventTrackerId: IsString.optional(),
	id: IsString,
	orgId: IsString,
	pathKey: IsString.optional(),
	serviceName: IsString,
	serviceType: IsString,
	eventType: IsString,
	originalSubject: IsString,
	originalStream: IsString,
	originalData: IsString,
	originalHeaders: IsString.optional(),
	lastError: IsString,
	errorCount: z.number(),
	firstDelivery: IsString,
	lastDelivery: IsString,
	totalDeliveries: z.number(),
	consumerName: IsString,
	retentionDays: z.number(),
});

export const ZodEventsDLQCursorResultSchema = z.object({
	items: z.array(ZodEventsDLQResponseSchema),
	nextCursor: IsStringDateFormat.optional().nullable(),
	prevCursor: IsStringDateFormat.optional().nullable(),
	hasNext: IsBoolean,
	hasPrevious: IsBoolean,
});

/**
 * Events DLQ Counts schemas
 */
export const ZodEventsDLQCountsQuerySchema = z.object({
	includeChildren: IsBoolean.optional(),
	startTime: IsStringDateFormat.optional(),
	endTime: IsStringDateFormat.optional(),
});

export const ZodEventsDLQServiceCountSchema = z.object({
	serviceType: IsString,
	count: z.number(),
});

export const ZodEventsDLQCountsResultSchema = z.object({
	counts: z.array(ZodEventsDLQServiceCountSchema),
	total: z.number(),
});

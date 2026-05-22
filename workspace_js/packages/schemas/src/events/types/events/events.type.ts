import { z } from 'zod';
import {
	ZodEventsRawQuerySchema,
	ZodEventsRawResponseSchema,
	ZodEventsRawCursorResultSchema,
	ZodEventsJsExecQuerySchema,
	ZodEventsJsExecResponseSchema,
	ZodEventsJsExecCursorResultSchema,
	ZodEventsRouterQuerySchema,
	ZodEventsRouterResponseSchema,
	ZodEventsRouterCursorResultSchema,
	ZodEventsBusinessRuleQuerySchema,
	ZodEventsBusinessRuleResponseSchema,
	ZodEventsBusinessRuleCursorResultSchema,
	ZodEventsTriggerQuerySchema,
	ZodEventsTriggerResponseSchema,
	ZodEventsTriggerCursorResultSchema,
	ZodEvaFilterSchema,
	ZodEventsStoreQuerySchema,
	ZodEventsStoreResponseSchema,
	ZodEventsStoreCursorResultSchema,
	ZodEventsWorkflowQuerySchema,
	ZodEventsWorkflowExecutionIdParamSchema,
	ZodEventsWorkflowResponseSchema,
	ZodEventsWorkflowCursorResultSchema,
	ZodEventsDLQQuerySchema,
	ZodEventsDLQResponseSchema,
	ZodEventsDLQCursorResultSchema,
	ZodEventsDLQCountsQuerySchema,
	ZodEventsDLQServiceCountSchema,
	ZodEventsDLQCountsResultSchema,
} from '@/events/schemas/events/events.schema';

// ============================================================================
// Events Raw Types
// ============================================================================

/**
 * Events Raw Query type - Used for cursor-based pagination of raw events
 */
export type EventsRawQuery = z.infer<typeof ZodEventsRawQuerySchema>;

/**
 * Events Raw Response type - Single raw event item
 */
export type EventsRawResponse = z.infer<typeof ZodEventsRawResponseSchema>;

/**
 * Events Raw Cursor Result type - Cursor-paginated response from API
 */
export type EventsRawCursorResult = z.infer<typeof ZodEventsRawCursorResultSchema>;

// ============================================================================
// Events JS Executor Types
// ============================================================================

/**
 * Events JS Executor Query type - Used for cursor-based pagination of JS executor events
 */
export type EventsJsExecQuery = z.infer<typeof ZodEventsJsExecQuerySchema>;

/**
 * Events JS Executor Response type - Single JS executor event item
 */
export type EventsJsExecResponse = z.infer<typeof ZodEventsJsExecResponseSchema>;

/**
 * Events JS Executor Cursor Result type - Cursor-paginated response from API
 */
export type EventsJsExecCursorResult = z.infer<typeof ZodEventsJsExecCursorResultSchema>;

// ============================================================================
// Events Router Types
// ============================================================================

/**
 * Events Router Query type - Used for cursor-based pagination of router events
 */
export type EventsRouterQuery = z.infer<typeof ZodEventsRouterQuerySchema>;

/**
 * Events Router Response type - Single router event item
 */
export type EventsRouterResponse = z.infer<typeof ZodEventsRouterResponseSchema>;

/**
 * Events Router Cursor Result type - Cursor-paginated response from API
 */
export type EventsRouterCursorResult = z.infer<typeof ZodEventsRouterCursorResultSchema>;

// ============================================================================
// Events Business Rule Types
// ============================================================================

/**
 * Events Business Rule Query type - Used for cursor-based pagination of business rule events
 */
export type EventsBusinessRuleQuery = z.infer<typeof ZodEventsBusinessRuleQuerySchema>;

/**
 * Events Business Rule Response type - Single business rule event item
 */
export type EventsBusinessRuleResponse = z.infer<typeof ZodEventsBusinessRuleResponseSchema>;

/**
 * Events Business Rule Cursor Result type - Cursor-paginated response from API
 */
export type EventsBusinessRuleCursorResult = z.infer<typeof ZodEventsBusinessRuleCursorResultSchema>;

// ============================================================================
// Events Trigger Types
// ============================================================================

/**
 * Events Trigger Query type - Used for cursor-based pagination of trigger events
 */
export type EventsTriggerQuery = z.infer<typeof ZodEventsTriggerQuerySchema>;

/**
 * Events Trigger Response type - Single trigger event item
 */
export type EventsTriggerResponse = z.infer<typeof ZodEventsTriggerResponseSchema>;

/**
 * Events Trigger Cursor Result type - Cursor-paginated response from API
 */
export type EventsTriggerCursorResult = z.infer<typeof ZodEventsTriggerCursorResultSchema>;

// ============================================================================
// Events Store Types (Processed Events with EVA fields)
// ============================================================================

/**
 * EVA Filter type - Single filter condition on an EVA dynamic field
 */
export type EvaFilter = z.infer<typeof ZodEvaFilterSchema>;

/**
 * Events Store Query type - Used for cursor-based pagination of processed events
 * Supports EVA dynamic field filters via evaFilters array (POST body)
 */
export type EventsStoreQuery = z.infer<typeof ZodEventsStoreQuerySchema>;

/**
 * Events Store Response type - Single processed event with EVA fields
 */
export type EventsStoreResponse = z.infer<typeof ZodEventsStoreResponseSchema>;

/**
 * Events Store Cursor Result type - Cursor-paginated response from API
 */
export type EventsStoreCursorResult = z.infer<typeof ZodEventsStoreCursorResultSchema>;

// ============================================================================
// Events Workflow Types
// ============================================================================

/**
 * Events Workflow Query type - Used for cursor-based pagination of workflow execution events
 */
export type EventsWorkflowQuery = z.infer<typeof ZodEventsWorkflowQuerySchema>;

/**
 * Events Workflow Execution ID Param type - URL param for getting a single workflow event
 */
export type EventsWorkflowExecutionIdParam = z.infer<typeof ZodEventsWorkflowExecutionIdParamSchema>;

/**
 * Events Workflow Response type - Single workflow execution event item
 */
export type EventsWorkflowResponse = z.infer<typeof ZodEventsWorkflowResponseSchema>;

/**
 * Events Workflow Cursor Result type - Cursor-paginated response from API
 */
export type EventsWorkflowCursorResult = z.infer<typeof ZodEventsWorkflowCursorResultSchema>;

// DLQ Types
export type EventsDLQQuery = z.infer<typeof ZodEventsDLQQuerySchema>;
export type EventsDLQResponse = z.infer<typeof ZodEventsDLQResponseSchema>;
export type EventsDLQCursorResult = z.infer<typeof ZodEventsDLQCursorResultSchema>;
export type EventsDLQCountsQuery = z.infer<typeof ZodEventsDLQCountsQuerySchema>;
export type EventsDLQServiceCount = z.infer<typeof ZodEventsDLQServiceCountSchema>;
export type EventsDLQCountsResult = z.infer<typeof ZodEventsDLQCountsResultSchema>;

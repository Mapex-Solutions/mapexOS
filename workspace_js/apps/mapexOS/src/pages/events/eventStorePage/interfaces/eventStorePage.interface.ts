/**
 * EventStorePage Interfaces
 */

/**
 * Fixed filters state for event store page
 * These are the static filters available for all events
 *
 * Based on EventsStoreQuery from backend DTO
 */
export interface EventStorePageFilters {
	/** Filter by thread ID for distributed tracing */
	threadId: string | undefined;
	/** Filter by asset ID */
	assetId: string | undefined;
	/** Filter by asset template ID */
	assetTemplateId: string | undefined;
	/** Filter by event type (telemetry, alarm, command) */
	eventType: string | undefined;
	/** Filter by source service */
	source: string | undefined;
	/** Filter events after this timestamp (ISO 8601) */
	startTime: string | undefined;
	/** Filter events before this timestamp (ISO 8601) */
	endTime: string | undefined;
}

/**
 * Cursor state for event store page
 * Used for cursor-based pagination
 */
export interface EventStorePageCursor {
	/** Current cursor timestamp (for fetching next/prev page) */
	current: string | undefined;
	/** Next cursor timestamp (for fetching older items) */
	next: string | undefined;
	/** Previous cursor timestamp (for fetching newer items) */
	prev: string | undefined;
	/** Whether there are more items after current page */
	hasNext: boolean;
	/** Whether there are more items before current page */
	hasPrevious: boolean;
}

/**
 * Event store item response
 * Represents a single processed event from the store with EVA fields
 *
 * Based on EventsStoreResponse from backend DTO
 */
export interface EventStoreItem {
	/** Event creation timestamp */
	created: string;
	/** Thread ID for distributed tracing */
	threadId?: string;
	/** Asset ID that generated this event */
	assetId: string;
	/** Asset template ID used to process this event */
	assetTemplateId?: string;
	/** Organization ID */
	orgId: string;
	/** Path key for hierarchical filtering */
	pathKey?: string;
	/** Event type (telemetry, alarm, command) — not currently populated */
	eventType?: string;
	/** Source service that created this event */
	source?: string;
	/** Original payload as JSON string */
	payload: string;
	/** Event metadata as JSON string */
	metadata?: string;
	/** EVA Number fields - Map of fieldId to numeric value */
	evaNumber?: Record<string, number>;
	/** EVA String fields - Map of fieldId to text value */
	evaString?: Record<string, string>;
	/** EVA Boolean fields - Map of fieldId to boolean value */
	evaBool?: Record<string, boolean>;
	/** EVA Date fields - Map of fieldId to date value */
	evaDate?: Record<string, string>;
	/** Retention period in days */
	retentionDays?: number;

	// Human-readable names (denormalized at write-time by Router)
	/** Asset display name */
	assetName?: string;
	/** Asset description */
	assetDescription?: string;
	/** Template display name */
	templateName?: string;
	/** Template description */
	templateDescription?: string;
}

package clickhouseRepo

import (
	"context"
	"fmt"
	"time"

	"events/src/modules/events/domain/entities"

	chModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/clickhouse/model"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

/**
 * Raw Events Repository Methods
 */

// SaveRawEvent stores a single raw event in the events_raw ClickHouse table.
//
// Uses chModel.Table for automatic JSON marshaling of map fields.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - event: The raw event to store
//
// Returns:
//   - error: If the insert fails
func (r *EventRepositoryClickHouse) SaveRawEvent(ctx context.Context, event *entities.RawEvent) error {
	if r.rawEventTable == nil {
		return fmt.Errorf("events_raw table model not initialized")
	}

	if err := r.rawEventTable.Insert(ctx, event); err != nil {
		logger.Error(err, "[REPO:Event] Failed to insert raw event")
		return fmt.Errorf("failed to insert raw event: %w", err)
	}

	logger.Info(fmt.Sprintf("[REPO:Event] Raw event saved: threadId=%s, source=%s", event.ThreadId, event.Source))
	return nil
}

// SaveRawEventBatch stores multiple raw events in the events_raw ClickHouse table efficiently.
//
// Uses chModel.Table for automatic JSON marshaling of map fields.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - events: Slice of raw events to store
//
// Returns:
//   - error: If the batch insert fails
func (r *EventRepositoryClickHouse) SaveRawEventBatch(ctx context.Context, events []*entities.RawEvent) error {
	if len(events) == 0 {
		return nil
	}

	if r.rawEventTable == nil {
		return fmt.Errorf("events_raw table model not initialized")
	}

	if err := r.rawEventTable.InsertBatch(ctx, events); err != nil {
		logger.Error(err, "[REPO:Event] Failed to save raw event batch")
		return fmt.Errorf("failed to save raw event batch: %w", err)
	}

	logger.Info(fmt.Sprintf("[REPO:Event] Raw batch saved: %d events", len(events)))
	return nil
}

// QueryEventsRaw retrieves a paginated and filtered list of raw events from ClickHouse.
//
// Uses chModel.Table with advanced filters for complex queries including
// time range filtering via BETWEEN operator.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - orgFilter: Organization filter (orgId, pathKey patterns)
//   - threadId, source: Module-specific filters (optional)
//   - startTime, endTime: Time range filters (optional)
//   - pagination: Pagination options (page, perPage)
//   - sort: Sort order (e.g., "created:desc")
//
// Returns:
//   - Paginated result containing RawEvent entities
//   - error: If query fails
func (r *EventRepositoryClickHouse) QueryEventsRaw(
	ctx context.Context,
	orgFilter chModel.Map,
	threadId *string,
	source *string,
	startTime *time.Time,
	endTime *time.Time,
	pagination *chModel.PaginationOpts,
	sort string,
) (*chModel.PaginatedResult[entities.RawEvent], error) {

	if r.rawEventTable == nil {
		return nil, fmt.Errorf("events_raw table model not initialized")
	}

	// Build filters using chModel.Filter for advanced operators
	filters := []chModel.Filter{}

	// Apply org filter
	if orgId, ok := orgFilter["orgId"]; ok && orgId != nil {
		filters = append(filters, chModel.Filter{
			Field:    "org_id",
			Operator: chModel.OpEqual,
			Value:    orgId,
		})
	}

	// Apply pathKey filter with LIKE pattern support
	if pathKey, ok := orgFilter["pathKey"]; ok && pathKey != nil {
		if pathKeyStr, isString := pathKey.(string); isString {
			if len(pathKeyStr) > 0 && pathKeyStr[len(pathKeyStr)-1] == '%' {
				filters = append(filters, chModel.Filter{
					Field:    "path_key",
					Operator: chModel.OpLike,
					Value:    pathKeyStr,
				})
			} else {
				filters = append(filters, chModel.Filter{
					Field:    "path_key",
					Operator: chModel.OpEqual,
					Value:    pathKeyStr,
				})
			}
		} else if pathKeyMap, isMap := pathKey.(map[string]interface{}); isMap {
			if regex, hasRegex := pathKeyMap["$regex"]; hasRegex {
				regexStr := fmt.Sprintf("%v", regex)
				// Convert MongoDB regex to LIKE pattern
				if len(regexStr) > 0 && regexStr[0] == '^' {
					regexStr = regexStr[1:]
				}
				filters = append(filters, chModel.Filter{
					Field:    "path_key",
					Operator: chModel.OpLike,
					Value:    regexStr + "%",
				})
			}
		}
	}

	// Apply module-specific filters
	if threadId != nil && *threadId != "" {
		filters = append(filters, chModel.Filter{
			Field:    "thread_id",
			Operator: chModel.OpEqual,
			Value:    *threadId,
		})
	}

	if source != nil && *source != "" {
		filters = append(filters, chModel.Filter{
			Field:    "source",
			Operator: chModel.OpEqual,
			Value:    *source,
		})
	}

	// Apply time range filters using BETWEEN or individual operators
	if startTime != nil && endTime != nil {
		filters = append(filters, chModel.Filter{
			Field:    "created",
			Operator: chModel.OpBetween,
			Value:    *startTime,
			EndValue: *endTime,
		})
	} else {
		if startTime != nil {
			filters = append(filters, chModel.Filter{
				Field:    "created",
				Operator: chModel.OpGreaterEqual,
				Value:    *startTime,
			})
		}
		if endTime != nil {
			filters = append(filters, chModel.Filter{
				Field:    "created",
				Operator: chModel.OpLessEqual,
				Value:    *endTime,
			})
		}
	}

	// Use FindWithFilters for advanced filtering
	return r.rawEventTable.FindWithFilters(ctx, filters, pagination, sort)
}

// QueryEventsRawCursor retrieves raw events using cursor-based pagination.
//
// This method is optimized for large datasets as it avoids COUNT queries.
// Uses created as the cursor for efficient seeks.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - orgFilter: Organization filter (orgId, pathKey patterns)
//   - threadId, source: Module-specific filters (optional)
//   - success: Filter by auth success status (optional)
//   - startTime, endTime: Time range filters (optional)
//   - cursorOpts: Cursor pagination options (cursor, direction, limit, sortAsc)
//
// Returns:
//   - TimeCursorResult containing RawEvent entities and cursor metadata
//   - error: If query fails
func (r *EventRepositoryClickHouse) QueryEventsRawCursor(
	ctx context.Context,
	orgFilter chModel.Map,
	eventTrackerId *string,
	threadId *string,
	source *string,
	success *bool,
	startTime *time.Time,
	endTime *time.Time,
	cursorOpts *chModel.TimeCursorOpts,
) (*chModel.TimeCursorResult[entities.RawEvent], error) {

	// Log received parameters for debugging
	logger.Info(fmt.Sprintf("[REPO:Event] QueryEventsRawCursor called with: orgFilter=%+v, eventTrackerId=%v, threadId=%v, source=%v, success=%v, startTime=%v, endTime=%v, cursorOpts=%+v",
		orgFilter,
		eventTrackerId,
		threadId,
		source,
		success,
		startTime,
		endTime,
		cursorOpts,
	))

	if r.rawEventTable == nil {
		return nil, fmt.Errorf("events_raw table model not initialized")
	}

	// Build filters using chModel.Filter for advanced operators
	filters := []chModel.Filter{}

	// Apply org filter
	if orgId, ok := orgFilter["orgId"]; ok && orgId != nil {
		filters = append(filters, chModel.Filter{
			Field:    "org_id",
			Operator: chModel.OpEqual,
			Value:    orgId,
		})
	}

	// Apply pathKey filter with LIKE pattern support
	if pathKey, ok := orgFilter["pathKey"]; ok && pathKey != nil {
		if pathKeyStr, isString := pathKey.(string); isString {
			if len(pathKeyStr) > 0 && pathKeyStr[len(pathKeyStr)-1] == '%' {
				filters = append(filters, chModel.Filter{
					Field:    "path_key",
					Operator: chModel.OpLike,
					Value:    pathKeyStr,
				})
			} else {
				filters = append(filters, chModel.Filter{
					Field:    "path_key",
					Operator: chModel.OpEqual,
					Value:    pathKeyStr,
				})
			}
		} else if pathKeyMap, isMap := pathKey.(map[string]interface{}); isMap {
			if regex, hasRegex := pathKeyMap["$regex"]; hasRegex {
				regexStr := fmt.Sprintf("%v", regex)
				// Convert MongoDB regex to LIKE pattern
				if len(regexStr) > 0 && regexStr[0] == '^' {
					regexStr = regexStr[1:]
				}
				filters = append(filters, chModel.Filter{
					Field:    "path_key",
					Operator: chModel.OpLike,
					Value:    regexStr + "%",
				})
			}
		}
	}

	// Apply eventTrackerId filter for end-to-end tracing
	if eventTrackerId != nil && *eventTrackerId != "" {
		filters = append(filters, chModel.Filter{
			Field:    "event_tracker_id",
			Operator: chModel.OpEqual,
			Value:    *eventTrackerId,
		})
	}

	// Apply module-specific filters
	if threadId != nil && *threadId != "" {
		filters = append(filters, chModel.Filter{
			Field:    "thread_id",
			Operator: chModel.OpEqual,
			Value:    *threadId,
		})
	}

	if source != nil && *source != "" {
		filters = append(filters, chModel.Filter{
			Field:    "source",
			Operator: chModel.OpEqual,
			Value:    *source,
		})
	}

	// Apply success filter
	if success != nil {
		// Convert bool to UInt8 (0 or 1) for ClickHouse
		var successValue uint8 = 0
		if *success {
			successValue = 1
		}
		filters = append(filters, chModel.Filter{
			Field:    "success",
			Operator: chModel.OpEqual,
			Value:    successValue,
		})
	}

	// Apply time range filters (bounds for the data, not cursor)
	// Note: ClickHouse DateTime is stored without timezone, so we use UTC directly
	// The client should send timestamps in UTC for consistent filtering
	if startTime != nil {
		// Use UTC to match ClickHouse DateTime storage
		startUTC := startTime.UTC()
		filters = append(filters, chModel.Filter{
			Field:    "created",
			Operator: chModel.OpGreaterEqual,
			Value:    startUTC,
		})
	}
	if endTime != nil {
		// Use UTC to match ClickHouse DateTime storage
		endUTC := endTime.UTC()
		filters = append(filters, chModel.Filter{
			Field:    "created",
			Operator: chModel.OpLessEqual,
			Value:    endUTC,
		})
	}

	// Use FindByCursor for cursor-based pagination
	return r.rawEventTable.FindByCursor(ctx, filters, cursorOpts)
}

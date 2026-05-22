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
 * JS Executor Events Repository Methods
 */

// SaveJsExecEventBatch stores multiple JS Executor events in ClickHouse efficiently.
//
// Uses chModel.Table for automatic field mapping and bulk insert.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - events: Slice of JS Executor events to store
//
// Returns:
//   - error: If the batch insert fails
func (r *EventRepositoryClickHouse) SaveJsExecEventBatch(ctx context.Context, events []*entities.JsExecEvent) error {
	if len(events) == 0 {
		return nil
	}

	if r.jsExecEventTable == nil {
		return fmt.Errorf("events_jsexecutor table model not initialized")
	}

	if err := r.jsExecEventTable.InsertBatch(ctx, events); err != nil {
		logger.Error(err, "[REPO:Event] Failed to save jsexec event batch")
		return fmt.Errorf("failed to save jsexec event batch: %w", err)
	}

	logger.Info(fmt.Sprintf("[REPO:Event] JsExec batch saved: %d events", len(events)))
	return nil
}

// QueryEventsJsExecCursor retrieves JS Executor events using cursor-based pagination.
//
// This method is optimized for large datasets as it avoids COUNT queries.
// Uses created as the cursor for efficient seeks.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - orgFilter: Organization filter (orgId, pathKey patterns)
//   - threadId: Filter by thread ID (asset UUID, optional)
//   - success: Filter by execution success status (optional)
//   - startTime, endTime: Time range filters (optional)
//   - execTimeOp: Operator for execution time filter (lt, lte, gt, gte, between)
//   - execTimeValue: Execution time value in milliseconds (optional)
//   - execTimeValueEnd: End value for "between" operator (optional)
//   - cursorOpts: Cursor pagination options (cursor, direction, limit, sortAsc)
//
// Returns:
//   - TimeCursorResult containing JsExecEvent entities and cursor metadata
//   - error: If query fails
func (r *EventRepositoryClickHouse) QueryEventsJsExecCursor(
	ctx context.Context,
	orgFilter chModel.Map,
	eventTrackerId *string,
	threadId *string,
	success *bool,
	startTime *time.Time,
	endTime *time.Time,
	execTimeOp *string,
	execTimeValue *uint32,
	execTimeValueEnd *uint32,
	cursorOpts *chModel.TimeCursorOpts,
) (*chModel.TimeCursorResult[entities.JsExecEvent], error) {

	logger.Info(fmt.Sprintf("[REPO:Event] QueryEventsJsExecCursor called with: orgFilter=%+v, eventTrackerId=%v, threadId=%v, success=%v, startTime=%v, endTime=%v, execTimeOp=%v, execTimeValue=%v, execTimeValueEnd=%v, cursorOpts=%+v",
		orgFilter,
		eventTrackerId,
		threadId,
		success,
		startTime,
		endTime,
		execTimeOp,
		execTimeValue,
		execTimeValueEnd,
		cursorOpts,
	))

	if r.jsExecEventTable == nil {
		return nil, fmt.Errorf("events_jsexecutor table model not initialized")
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

	// Apply time range filters
	if startTime != nil {
		startUTC := startTime.UTC()
		filters = append(filters, chModel.Filter{
			Field:    "created",
			Operator: chModel.OpGreaterEqual,
			Value:    startUTC,
		})
	}
	if endTime != nil {
		endUTC := endTime.UTC()
		filters = append(filters, chModel.Filter{
			Field:    "created",
			Operator: chModel.OpLessEqual,
			Value:    endUTC,
		})
	}

	// Apply execution time filter
	if execTimeOp != nil && execTimeValue != nil {
		switch *execTimeOp {
		case "lt":
			filters = append(filters, chModel.Filter{
				Field:    "total_execution_time",
				Operator: chModel.OpLess,
				Value:    *execTimeValue,
			})
		case "lte":
			filters = append(filters, chModel.Filter{
				Field:    "total_execution_time",
				Operator: chModel.OpLessEqual,
				Value:    *execTimeValue,
			})
		case "gt":
			filters = append(filters, chModel.Filter{
				Field:    "total_execution_time",
				Operator: chModel.OpGreater,
				Value:    *execTimeValue,
			})
		case "gte":
			filters = append(filters, chModel.Filter{
				Field:    "total_execution_time",
				Operator: chModel.OpGreaterEqual,
				Value:    *execTimeValue,
			})
		case "between":
			// For "between" operator, we need both start and end values
			filters = append(filters, chModel.Filter{
				Field:    "total_execution_time",
				Operator: chModel.OpGreaterEqual,
				Value:    *execTimeValue,
			})
			if execTimeValueEnd != nil {
				filters = append(filters, chModel.Filter{
					Field:    "total_execution_time",
					Operator: chModel.OpLessEqual,
					Value:    *execTimeValueEnd,
				})
			}
		}
	}

	// Use FindByCursor for cursor-based pagination
	return r.jsExecEventTable.FindByCursor(ctx, filters, cursorOpts)
}

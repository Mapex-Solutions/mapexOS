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
 * Trigger Events Repository Methods
 */

// SaveTriggerEventBatch stores multiple trigger events in ClickHouse efficiently.
//
// Uses chModel.Table for automatic field mapping and bulk insert.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - events: Slice of trigger events to store
//
// Returns:
//   - error: If the batch insert fails
func (r *EventRepositoryClickHouse) SaveTriggerEventBatch(ctx context.Context, events []*entities.TriggerEvent) error {
	if len(events) == 0 {
		return nil
	}

	if r.triggerEventTable == nil {
		return fmt.Errorf("events_trigger table model not initialized")
	}

	if err := r.triggerEventTable.InsertBatch(ctx, events); err != nil {
		logger.Error(err, "[REPO:Event] Failed to save trigger event batch")
		return fmt.Errorf("failed to save trigger event batch: %w", err)
	}

	logger.Info(fmt.Sprintf("[REPO:Event] Trigger batch saved: %d events", len(events)))
	return nil
}

// QueryEventsTriggerCursor retrieves trigger events using cursor-based pagination.
//
// This method is optimized for large datasets as it avoids COUNT queries.
// Uses created as the cursor for efficient seeks.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - orgFilter: Organization filter (orgId, pathKey patterns)
//   - triggerId: Filter by trigger ID (optional)
//   - triggerType: Filter by trigger type (http, mqtt, email, etc.) (optional)
//   - category: Filter by category (technical, communication) (optional)
//   - source: Filter by source (router) (optional)
//   - success: Filter by success status (optional)
//   - startTime, endTime: Time range filters (optional)
//   - cursorOpts: Cursor pagination options (cursor, direction, limit, sortAsc)
//
// Returns:
//   - TimeCursorResult containing TriggerEvent entities and cursor metadata
//   - error: If query fails
func (r *EventRepositoryClickHouse) QueryEventsTriggerCursor(
	ctx context.Context,
	orgFilter chModel.Map,
	eventTrackerId *string,
	triggerId *string,
	triggerType *string,
	category *string,
	source *string,
	success *bool,
	startTime *time.Time,
	endTime *time.Time,
	cursorOpts *chModel.TimeCursorOpts,
) (*chModel.TimeCursorResult[entities.TriggerEvent], error) {

	logger.Info(fmt.Sprintf("[REPO:Event] QueryEventsTriggerCursor called with: orgFilter=%+v, eventTrackerId=%v, triggerId=%v, triggerType=%v, category=%v, source=%v, success=%v, startTime=%v, endTime=%v, cursorOpts=%+v",
		orgFilter,
		eventTrackerId,
		triggerId,
		triggerType,
		category,
		source,
		success,
		startTime,
		endTime,
		cursorOpts,
	))

	if r.triggerEventTable == nil {
		return nil, fmt.Errorf("events_trigger table model not initialized")
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

	// Apply trigger specific filters
	if triggerId != nil && *triggerId != "" {
		filters = append(filters, chModel.Filter{
			Field:    "trigger_id",
			Operator: chModel.OpEqual,
			Value:    *triggerId,
		})
	}

	if triggerType != nil && *triggerType != "" {
		filters = append(filters, chModel.Filter{
			Field:    "trigger_type",
			Operator: chModel.OpEqual,
			Value:    *triggerType,
		})
	}

	if category != nil && *category != "" {
		filters = append(filters, chModel.Filter{
			Field:    "category",
			Operator: chModel.OpEqual,
			Value:    *category,
		})
	}

	if source != nil && *source != "" {
		filters = append(filters, chModel.Filter{
			Field:    "source",
			Operator: chModel.OpEqual,
			Value:    *source,
		})
	}

	if success != nil {
		successVal := uint8(0)
		if *success {
			successVal = 1
		}
		filters = append(filters, chModel.Filter{
			Field:    "success",
			Operator: chModel.OpEqual,
			Value:    successVal,
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

	// Use FindByCursor for cursor-based pagination
	return r.triggerEventTable.FindByCursor(ctx, filters, cursorOpts)
}

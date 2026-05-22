package clickhouseRepo

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"events/src/modules/events/domain/entities"

	dtos "github.com/Mapex-Solutions/MapexOS/contracts/services/events/events"
	chModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/clickhouse/model"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

/**
 * Event Store Repository Methods (Processed Events with EVA)
 */

// castEvaValue converts a string value to the appropriate Go type for ClickHouse based on the EVA bucket.
//
// Parameters:
//   - bucket: EVA type ("number", "string", "bool", "date")
//   - val: String value to cast
//
// Returns:
//   - interface{}: Typed value for ClickHouse query parameter
//   - error: If parsing fails
func castEvaValue(bucket string, val string) (interface{}, error) {
	switch bucket {
	case "number":
		return strconv.ParseFloat(val, 64)
	case "string":
		return val, nil
	case "bool":
		lower := strings.ToLower(val)
		if lower == "true" || lower == "1" {
			return uint8(1), nil
		}
		return uint8(0), nil
	case "date":
		return time.Parse(time.RFC3339, val)
	default:
		return val, nil
	}
}

// SaveEventStoreBatch stores multiple processed events with EVA fields in ClickHouse.
//
// This method is the PRIMARY batch insert for processed events from the router service.
// Events are stored with EVA (Entity-Value-Attribute) MAP fields for efficient querying.
//
// EVA Fields:
//   - eva_number: Map(UInt16, Float64) for numeric values
//   - eva_string: Map(UInt16, String) for text values
//   - eva_bool: Map(UInt16, UInt8) for boolean values
//   - eva_date: Map(UInt16, DateTime64) for datetime values
//
// Uses chModel.Table for automatic field mapping and bulk insert.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - events: Slice of processed events with EVA fields to store
//
// Returns:
//   - error: If the batch insert fails
func (r *EventRepositoryClickHouse) SaveEventStoreBatch(ctx context.Context, events []*entities.Event) error {
	if len(events) == 0 {
		return nil
	}

	if r.eventTable == nil {
		return fmt.Errorf("events table model not initialized")
	}

	if err := r.eventTable.InsertBatch(ctx, events); err != nil {
		logger.Error(err, "[REPO:Event] Failed to save event store batch")
		return fmt.Errorf("failed to save event store batch: %w", err)
	}

	logger.Info(fmt.Sprintf("[REPO:Event] Event store batch saved: %d events", len(events)))
	return nil
}

// QueryEventsStoreCursor retrieves processed events using cursor-based pagination.
//
// This method is optimized for large datasets as it avoids COUNT queries.
// Uses created as the cursor for efficient seeks.
// Supports both indexed column filters and EVA MAP column filters.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - orgFilter: Organization filter (orgId, pathKey patterns)
//   - eventTrackerId: Filter by event tracker ID (optional)
//   - threadId: Filter by thread ID for distributed tracing (optional)
//   - assetId: Filter by asset ID (optional)
//   - templateId: Filter by asset template ID (optional)
//   - eventType: Filter by event type (telemetry, alarm, command) (optional)
//   - source: Filter by source service (optional)
//   - startTime, endTime: Time range filters (optional)
//   - evaFilters: EVA dynamic field filters with operators (optional)
//   - cursorOpts: Cursor pagination options (cursor, direction, limit, sortAsc)
//
// Returns:
//   - TimeCursorResult containing Event entities and cursor metadata
//   - error: If query fails
func (r *EventRepositoryClickHouse) QueryEventsStoreCursor(
	ctx context.Context,
	orgFilter chModel.Map,
	eventTrackerId *string,
	threadId *string,
	assetId *string,
	templateId *string,
	eventType *string,
	source *string,
	startTime *time.Time,
	endTime *time.Time,
	evaFilters []dtos.EvaFilter,
	cursorOpts *chModel.TimeCursorOpts,
) (*chModel.TimeCursorResult[entities.Event], error) {

	logger.Info(fmt.Sprintf("[REPO:Event] QueryEventsStoreCursor called with: orgFilter=%+v, eventTrackerId=%v, threadId=%v, assetId=%v, templateId=%v, eventType=%v, source=%v, startTime=%v, endTime=%v, evaFilters=%d, cursorOpts=%+v",
		orgFilter,
		eventTrackerId,
		threadId,
		assetId,
		templateId,
		eventType,
		source,
		startTime,
		endTime,
		len(evaFilters),
		cursorOpts,
	))

	if r.eventTable == nil {
		return nil, fmt.Errorf("events table model not initialized")
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

	// Apply event store specific filters
	if eventTrackerId != nil && *eventTrackerId != "" {
		filters = append(filters, chModel.Filter{
			Field:    "event_tracker_id",
			Operator: chModel.OpEqual,
			Value:    *eventTrackerId,
		})
	}

	if threadId != nil && *threadId != "" {
		filters = append(filters, chModel.Filter{
			Field:    "thread_id",
			Operator: chModel.OpEqual,
			Value:    *threadId,
		})
	}

	if assetId != nil && *assetId != "" {
		filters = append(filters, chModel.Filter{
			Field:    "asset_id",
			Operator: chModel.OpEqual,
			Value:    *assetId,
		})
	}

	if templateId != nil && *templateId != "" {
		filters = append(filters, chModel.Filter{
			Field:    "asset_template_id",
			Operator: chModel.OpEqual,
			Value:    *templateId,
		})
	}

	if eventType != nil && *eventType != "" {
		filters = append(filters, chModel.Filter{
			Field:    "event_type",
			Operator: chModel.OpEqual,
			Value:    *eventType,
		})
	}

	if source != nil && *source != "" {
		filters = append(filters, chModel.Filter{
			Field:    "source",
			Operator: chModel.OpEqual,
			Value:    *source,
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

	// Apply EVA MAP column filters
	// Each filter accesses a MAP column by fieldId: eva_number[42] >= 10.5
	for _, ef := range evaFilters {
		col, ok := evaBucketColumn[ef.Bucket]
		if !ok {
			logger.Info(fmt.Sprintf("[REPO:Event] Skipping unknown EVA bucket: %s", ef.Bucket))
			continue
		}

		op, ok := evaOperatorMap[ef.Operator]
		if !ok {
			logger.Info(fmt.Sprintf("[REPO:Event] Skipping unknown EVA operator: %s", ef.Operator))
			continue
		}

		// Use mapContains to skip rows where the MAP key doesn't exist.
		// Without this, ClickHouse returns the type's default value (0, "", etc.)
		// which can produce false positives (e.g. 0 < 49 = true).
		field := fmt.Sprintf("mapContains(%s, %d) AND %s[%d]", col, ef.FieldId, col, ef.FieldId)

		value, err := castEvaValue(ef.Bucket, ef.Value)
		if err != nil {
			logger.Info(fmt.Sprintf("[REPO:Event] Skipping EVA filter with invalid value: field=%d, bucket=%s, value=%s, err=%v", ef.FieldId, ef.Bucket, ef.Value, err))
			continue
		}

		if op == chModel.OpBetween {
			endValue, err := castEvaValue(ef.Bucket, ef.EndValue)
			if err != nil {
				logger.Info(fmt.Sprintf("[REPO:Event] Skipping EVA BETWEEN filter with invalid endValue: field=%d, endValue=%s, err=%v", ef.FieldId, ef.EndValue, err))
				continue
			}
			filters = append(filters, chModel.Filter{
				Field:    field,
				Operator: op,
				Value:    value,
				EndValue: endValue,
			})
		} else {
			filters = append(filters, chModel.Filter{
				Field:    field,
				Operator: op,
				Value:    value,
			})
		}
	}

	// Use FindByCursor for cursor-based pagination
	return r.eventTable.FindByCursor(ctx, filters, cursorOpts)
}

// GetEventStoreByTrackerId retrieves a single event by event_tracker_id.
// Uses FindByCursor with limit 1 for efficient single-row lookup.
func (r *EventRepositoryClickHouse) GetEventStoreByTrackerId(ctx context.Context, eventTrackerId string) (*entities.Event, error) {
	if r.eventTable == nil {
		return nil, fmt.Errorf("events table model not initialized")
	}

	filters := []chModel.Filter{
		{
			Field:    "event_tracker_id",
			Operator: chModel.OpEqual,
			Value:    eventTrackerId,
		},
	}

	cursorOpts := &chModel.TimeCursorOpts{
		Limit:   1,
		SortAsc: false,
	}

	result, err := r.eventTable.FindByCursor(ctx, filters, cursorOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to query event by trackerId: %w", err)
	}

	if len(result.Items) == 0 {
		return nil, fmt.Errorf("event not found: %s", eventTrackerId)
	}

	return &result.Items[0], nil
}

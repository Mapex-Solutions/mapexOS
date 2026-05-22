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
 * Dead Letter Queue (DLQ) Events Repository Methods
 */

// SaveDLQEvent stores a single DLQ event in ClickHouse.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - event: DLQ event to store
//
// Returns:
//   - error: If the insert fails
func (r *EventRepositoryClickHouse) SaveDLQEvent(ctx context.Context, event *entities.DLQEvent) error {
	if event == nil {
		return nil
	}

	if r.dlqEventTable == nil {
		return fmt.Errorf("events_dlq table model not initialized")
	}

	events := []*entities.DLQEvent{event}
	if err := r.dlqEventTable.InsertBatch(ctx, events); err != nil {
		logger.Error(err, "[REPO:Event] Failed to save DLQ event")
		return fmt.Errorf("failed to save DLQ event: %w", err)
	}

	logger.Info(fmt.Sprintf("[REPO:Event] DLQ event saved: id=%s, service=%s", event.ID, event.ServiceName))
	return nil
}

// SaveDLQEventBatch stores multiple DLQ events in ClickHouse efficiently.
//
// Uses chModel.Table for automatic field mapping and bulk insert.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - events: Slice of DLQ events to store
//
// Returns:
//   - error: If the batch insert fails
func (r *EventRepositoryClickHouse) SaveDLQEventBatch(ctx context.Context, events []*entities.DLQEvent) error {
	if len(events) == 0 {
		return nil
	}

	if r.dlqEventTable == nil {
		return fmt.Errorf("events_dlq table model not initialized")
	}

	if err := r.dlqEventTable.InsertBatch(ctx, events); err != nil {
		logger.Error(err, "[REPO:Event] Failed to save DLQ event batch")
		return fmt.Errorf("failed to save DLQ event batch: %w", err)
	}

	logger.Info(fmt.Sprintf("[REPO:Event] DLQ batch saved: %d events", len(events)))
	return nil
}

// QueryEventsDLQCursor retrieves DLQ events using cursor-based pagination.
//
// This method is optimized for large datasets as it avoids COUNT queries.
// Uses created as the cursor for efficient seeks.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - orgFilter: Organization filter (orgId, pathKey patterns)
//   - serviceName: Filter by service name (optional)
//   - serviceType: Filter by service type (optional)
//   - eventType: Filter by event type (optional)
//   - startTime, endTime: Time range filters (optional)
//   - cursorOpts: Cursor pagination options (cursor, direction, limit, sortAsc)
//
// Returns:
//   - TimeCursorResult containing DLQEvent entities and cursor metadata
//   - error: If query fails
func (r *EventRepositoryClickHouse) QueryEventsDLQCursor(
	ctx context.Context,
	orgFilter chModel.Map,
	eventTrackerId *string,
	serviceName *string,
	serviceType *string,
	eventType *string,
	lastError *string,
	startTime *time.Time,
	endTime *time.Time,
	cursorOpts *chModel.TimeCursorOpts,
) (*chModel.TimeCursorResult[entities.DLQEvent], error) {

	logger.Debug(fmt.Sprintf("[REPO:Event] QueryEventsDLQCursor: orgFilter=%+v, serviceName=%v, serviceType=%v, eventType=%v, lastError=%v",
		orgFilter, serviceName, serviceType, eventType, lastError,
	))

	if r.dlqEventTable == nil {
		return nil, fmt.Errorf("events_dlq table model not initialized")
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

	// Apply service-specific filters
	if serviceName != nil && *serviceName != "" {
		filters = append(filters, chModel.Filter{
			Field:    "service_name",
			Operator: chModel.OpEqual,
			Value:    *serviceName,
		})
	}

	if serviceType != nil && *serviceType != "" {
		filters = append(filters, chModel.Filter{
			Field:    "service_type",
			Operator: chModel.OpEqual,
			Value:    *serviceType,
		})
	}

	if eventType != nil && *eventType != "" {
		filters = append(filters, chModel.Filter{
			Field:    "event_type",
			Operator: chModel.OpEqual,
			Value:    *eventType,
		})
	}

	if lastError != nil && *lastError != "" {
		filters = append(filters, chModel.Filter{
			Field:    "last_error",
			Operator: chModel.OpLike,
			Value:    "%" + *lastError + "%",
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
	return r.dlqEventTable.FindByCursor(ctx, filters, cursorOpts)
}

// CountByServiceType counts DLQ entries grouped by service_type.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - orgFilter: Organization filter (orgId, pathKey patterns)
//   - startTime, endTime: Time range filters (optional)
//
// Returns:
//   - Slice of DLQServiceCount with service_type and count
//   - error: If query fails
func (r *EventRepositoryClickHouse) CountByServiceType(
	ctx context.Context,
	orgFilter chModel.Map,
	startTime *time.Time,
	endTime *time.Time,
) ([]entities.DLQServiceCount, error) {

	query := "SELECT service_type, count() as count FROM events_dlq WHERE 1=1"
	args := []interface{}{}

	// Apply org filter
	if orgId, ok := orgFilter["orgId"]; ok && orgId != nil {
		query += " AND org_id = ?"
		args = append(args, orgId)
	}
	if pathKey, ok := orgFilter["pathKey"]; ok && pathKey != nil {
		if pathKeyStr, isString := pathKey.(string); isString {
			if len(pathKeyStr) > 0 && pathKeyStr[len(pathKeyStr)-1] == '%' {
				query += " AND path_key LIKE ?"
				args = append(args, pathKeyStr)
			} else {
				query += " AND path_key = ?"
				args = append(args, pathKeyStr)
			}
		}
	}

	// Apply time range filters
	if startTime != nil {
		query += " AND created >= ?"
		args = append(args, startTime.UTC())
	}
	if endTime != nil {
		query += " AND created <= ?"
		args = append(args, endTime.UTC())
	}

	query += " GROUP BY service_type ORDER BY count DESC"

	rows, err := r.conn.Query(ctx, query, args...)
	if err != nil {
		logger.Error(err, "[REPO:Event] Failed to count DLQ by service type")
		return nil, fmt.Errorf("failed to count DLQ by service type: %w", err)
	}
	defer rows.Close()

	var results []entities.DLQServiceCount
	for rows.Next() {
		var item entities.DLQServiceCount
		if err := rows.Scan(&item.ServiceType, &item.Count); err != nil {
			logger.Error(err, "[REPO:Event] Failed to scan DLQ service count")
			continue
		}
		results = append(results, item)
	}

	return results, nil
}

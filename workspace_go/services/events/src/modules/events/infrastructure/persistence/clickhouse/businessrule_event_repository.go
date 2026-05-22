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
 * Business Rule Events Repository Methods
 */

// SaveBusinessRuleEventBatch stores multiple business rule events in ClickHouse efficiently.
//
// Uses chModel.Table for automatic field mapping and bulk insert.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - events: Slice of business rule events to store
//
// Returns:
//   - error: If the batch insert fails
func (r *EventRepositoryClickHouse) SaveBusinessRuleEventBatch(ctx context.Context, events []*entities.BusinessRuleEvent) error {
	if len(events) == 0 {
		return nil
	}

	if r.businessRuleEventTable == nil {
		return fmt.Errorf("events_businessrule table model not initialized")
	}

	if err := r.businessRuleEventTable.InsertBatch(ctx, events); err != nil {
		logger.Error(err, "[REPO:Event] Failed to save business rule event batch")
		return fmt.Errorf("failed to save business rule event batch: %w", err)
	}

	logger.Info(fmt.Sprintf("[REPO:Event] Business rule batch saved: %d events", len(events)))
	return nil
}

// QueryEventsBusinessRuleCursor retrieves business rule events using cursor-based pagination.
//
// This method is optimized for large datasets as it avoids COUNT queries.
// Uses created as the cursor for efficient seeks.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - orgFilter: Organization filter (orgId, pathKey patterns)
//   - threadId: Filter by thread ID (optional)
//   - ruleId: Filter by rule template ID (optional)
//   - businessRuleId: Filter by business rule ID (optional)
//   - matched: Filter by matched status (optional)
//   - startTime, endTime: Time range filters (optional)
//   - cursorOpts: Cursor pagination options (cursor, direction, limit, sortAsc)
//
// Returns:
//   - TimeCursorResult containing BusinessRuleEvent entities and cursor metadata
//   - error: If query fails
func (r *EventRepositoryClickHouse) QueryEventsBusinessRuleCursor(
	ctx context.Context,
	orgFilter chModel.Map,
	eventTrackerId *string,
	threadId *string,
	ruleId *string,
	businessRuleId *string,
	matched *bool,
	startTime *time.Time,
	endTime *time.Time,
	cursorOpts *chModel.TimeCursorOpts,
) (*chModel.TimeCursorResult[entities.BusinessRuleEvent], error) {

	logger.Info(fmt.Sprintf("[REPO:Event] QueryEventsBusinessRuleCursor called with: orgFilter=%+v, eventTrackerId=%v, threadId=%v, ruleId=%v, businessRuleId=%v, matched=%v, startTime=%v, endTime=%v, cursorOpts=%+v",
		orgFilter,
		eventTrackerId,
		threadId,
		ruleId,
		businessRuleId,
		matched,
		startTime,
		endTime,
		cursorOpts,
	))

	if r.businessRuleEventTable == nil {
		return nil, fmt.Errorf("events_businessrule table model not initialized")
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

	// Apply business rule specific filters
	if threadId != nil && *threadId != "" {
		filters = append(filters, chModel.Filter{
			Field:    "thread_id",
			Operator: chModel.OpEqual,
			Value:    *threadId,
		})
	}

	if ruleId != nil && *ruleId != "" {
		filters = append(filters, chModel.Filter{
			Field:    "rule_id",
			Operator: chModel.OpEqual,
			Value:    *ruleId,
		})
	}

	if businessRuleId != nil && *businessRuleId != "" {
		filters = append(filters, chModel.Filter{
			Field:    "business_rule_id",
			Operator: chModel.OpEqual,
			Value:    *businessRuleId,
		})
	}

	if matched != nil {
		matchedVal := uint8(0)
		if *matched {
			matchedVal = 1
		}
		filters = append(filters, chModel.Filter{
			Field:    "matched",
			Operator: chModel.OpEqual,
			Value:    matchedVal,
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
	return r.businessRuleEventTable.FindByCursor(ctx, filters, cursorOpts)
}

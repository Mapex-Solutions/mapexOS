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
 * Workflow Events Repository Methods
 */

// SaveWorkflowEventBatch stores multiple workflow events in ClickHouse efficiently.
//
// Uses chModel.Table for automatic field mapping and bulk insert.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - events: Slice of workflow events to store
//
// Returns:
//   - error: If the batch insert fails
func (r *EventRepositoryClickHouse) SaveWorkflowEventBatch(ctx context.Context, events []*entities.WorkflowEvent) error {
	if len(events) == 0 {
		return nil
	}

	if r.workflowEventTable == nil {
		return fmt.Errorf("events_workflow table model not initialized")
	}

	if err := r.workflowEventTable.InsertBatch(ctx, events); err != nil {
		logger.Error(err, "[REPO:Event] Failed to save workflow event batch")
		return fmt.Errorf("failed to save workflow event batch: %w", err)
	}

	logger.Info(fmt.Sprintf("[REPO:Event] Workflow batch saved: %d events", len(events)))
	return nil
}

// QueryEventsWorkflowCursor retrieves workflow events using cursor-based pagination.
//
// This method is optimized for large datasets as it avoids COUNT queries.
// Uses created as the cursor for efficient seeks.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - orgFilter: Organization filter (orgId, pathKey patterns)
//   - eventTrackerId: Filter by event tracker ID (optional)
//   - workflowUUID: Filter by workflow UUID (optional)
//   - instanceId: Filter by instance ID (optional)
//   - definitionId: Filter by definition ID (optional)
//   - status: Filter by terminal status (completed, failed, cancelled) (optional)
//   - success: Filter by success status (optional)
//   - startTime, endTime: Time range filters (optional)
//   - cursorOpts: Cursor pagination options (cursor, direction, limit, sortAsc)
//
// Returns:
//   - TimeCursorResult containing WorkflowEvent entities and cursor metadata
//   - error: If query fails
func (r *EventRepositoryClickHouse) QueryEventsWorkflowCursor(
	ctx context.Context,
	orgFilter chModel.Map,
	eventTrackerId *string,
	workflowUUID *string,
	instanceId *string,
	definitionId *string,
	status *string,
	success *bool,
	startTime *time.Time,
	endTime *time.Time,
	cursorOpts *chModel.TimeCursorOpts,
) (*chModel.TimeCursorResult[entities.WorkflowEvent], error) {

	if r.workflowEventTable == nil {
		return nil, fmt.Errorf("events_workflow table model not initialized")
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

	// Apply workflow specific filters
	if workflowUUID != nil && *workflowUUID != "" {
		filters = append(filters, chModel.Filter{
			Field:    "workflow_uuid",
			Operator: chModel.OpEqual,
			Value:    *workflowUUID,
		})
	}

	if instanceId != nil && *instanceId != "" {
		filters = append(filters, chModel.Filter{
			Field:    "instance_id",
			Operator: chModel.OpEqual,
			Value:    *instanceId,
		})
	}

	if definitionId != nil && *definitionId != "" {
		filters = append(filters, chModel.Filter{
			Field:    "definition_id",
			Operator: chModel.OpEqual,
			Value:    *definitionId,
		})
	}

	if status != nil && *status != "" {
		filters = append(filters, chModel.Filter{
			Field:    "status",
			Operator: chModel.OpEqual,
			Value:    *status,
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
	return r.workflowEventTable.FindByCursor(ctx, filters, cursorOpts)
}

// FindWorkflowEventByExecutionId retrieves a single workflow event by org_id + execution_id.
func (r *EventRepositoryClickHouse) FindWorkflowEventByExecutionId(ctx context.Context, orgId string, executionId string) (*entities.WorkflowEvent, error) {
	if r.workflowEventTable == nil {
		return nil, fmt.Errorf("events_workflow table model not initialized")
	}

	filters := []chModel.Filter{
		{Field: "org_id", Operator: chModel.OpEqual, Value: orgId},
		{Field: "execution_id", Operator: chModel.OpEqual, Value: executionId},
	}

	pagination := &chModel.PaginationOpts{Page: 1, PerPage: 1}
	result, err := r.workflowEventTable.FindWithFilters(ctx, filters, pagination, "-created")
	if err != nil {
		return nil, fmt.Errorf("[REPO:Event] FindWorkflowEventByExecutionId failed: %w", err)
	}

	if len(result.Items) == 0 {
		return nil, nil
	}

	return &result.Items[0], nil
}

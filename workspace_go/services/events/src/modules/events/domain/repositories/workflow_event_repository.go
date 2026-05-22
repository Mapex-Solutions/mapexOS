package repositories

import (
	"context"
	"time"

	"events/src/modules/events/domain/entities"

	chModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/clickhouse/model"
)

// WorkflowEventRepository defines storage operations for Workflow execution events.
type WorkflowEventRepository interface {
	SaveWorkflowEventBatch(ctx context.Context, events []*entities.WorkflowEvent) error

	QueryEventsWorkflowCursor(
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
	) (*chModel.TimeCursorResult[entities.WorkflowEvent], error)

	// FindWorkflowEventByExecutionId retrieves a single workflow event by org_id + execution_id.
	FindWorkflowEventByExecutionId(ctx context.Context, orgId string, executionId string) (*entities.WorkflowEvent, error)
}

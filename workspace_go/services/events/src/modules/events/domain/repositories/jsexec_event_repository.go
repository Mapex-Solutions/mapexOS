package repositories

import (
	"context"
	"time"

	"events/src/modules/events/domain/entities"

	chModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/clickhouse/model"
)

// JsExecEventRepository defines storage operations for JS Executor events.
type JsExecEventRepository interface {
	SaveJsExecEventBatch(ctx context.Context, events []*entities.JsExecEvent) error

	QueryEventsJsExecCursor(
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
	) (*chModel.TimeCursorResult[entities.JsExecEvent], error)
}

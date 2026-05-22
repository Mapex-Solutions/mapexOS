package repositories

import (
	"context"
	"time"

	"events/src/modules/events/domain/entities"

	chModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/clickhouse/model"
)

// TriggerEventRepository defines storage operations for Trigger events.
type TriggerEventRepository interface {
	SaveTriggerEventBatch(ctx context.Context, events []*entities.TriggerEvent) error

	QueryEventsTriggerCursor(
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
	) (*chModel.TimeCursorResult[entities.TriggerEvent], error)
}

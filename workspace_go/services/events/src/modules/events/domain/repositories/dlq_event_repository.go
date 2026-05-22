package repositories

import (
	"context"
	"time"

	"events/src/modules/events/domain/entities"

	chModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/clickhouse/model"
)

// DLQEventRepository defines storage operations for Dead Letter Queue events.
type DLQEventRepository interface {
	SaveDLQEvent(ctx context.Context, event *entities.DLQEvent) error
	SaveDLQEventBatch(ctx context.Context, events []*entities.DLQEvent) error

	QueryEventsDLQCursor(
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
	) (*chModel.TimeCursorResult[entities.DLQEvent], error)

	CountByServiceType(
		ctx context.Context,
		orgFilter chModel.Map,
		startTime *time.Time,
		endTime *time.Time,
	) ([]entities.DLQServiceCount, error)
}

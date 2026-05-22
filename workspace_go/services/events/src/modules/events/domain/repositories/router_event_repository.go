package repositories

import (
	"context"
	"time"

	"events/src/modules/events/domain/entities"

	chModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/clickhouse/model"
)

// RouterEventRepository defines storage operations for Router events.
type RouterEventRepository interface {
	SaveRouterEventBatch(ctx context.Context, events []*entities.RouterEvent) error

	QueryEventsRouterCursor(
		ctx context.Context,
		orgFilter chModel.Map,
		eventTrackerId *string,
		threadId *string,
		assetId *string,
		routerId *string,
		success *bool,
		publishedCount *int,
		startTime *time.Time,
		endTime *time.Time,
		cursorOpts *chModel.TimeCursorOpts,
	) (*chModel.TimeCursorResult[entities.RouterEvent], error)
}

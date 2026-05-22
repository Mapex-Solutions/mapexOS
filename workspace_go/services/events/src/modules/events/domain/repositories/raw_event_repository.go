package repositories

import (
	"context"
	"time"

	"events/src/modules/events/domain/entities"

	chModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/clickhouse/model"
)

// RawEventRepository defines storage operations for raw (unprocessed) events.
type RawEventRepository interface {
	SaveRawEvent(ctx context.Context, event *entities.RawEvent) error
	SaveRawEventBatch(ctx context.Context, events []*entities.RawEvent) error

	// DEPRECATED: Use QueryEventsRawCursor for large datasets.
	QueryEventsRaw(
		ctx context.Context,
		orgFilter chModel.Map,
		threadId *string,
		source *string,
		startTime *time.Time,
		endTime *time.Time,
		pagination *chModel.PaginationOpts,
		sort string,
	) (*chModel.PaginatedResult[entities.RawEvent], error)

	QueryEventsRawCursor(
		ctx context.Context,
		orgFilter chModel.Map,
		eventTrackerId *string,
		threadId *string,
		source *string,
		success *bool,
		startTime *time.Time,
		endTime *time.Time,
		cursorOpts *chModel.TimeCursorOpts,
	) (*chModel.TimeCursorResult[entities.RawEvent], error)
}

package repositories

import (
	"context"
	"time"

	"events/src/modules/events/domain/entities"

	dtos "github.com/Mapex-Solutions/MapexOS/contracts/services/events/events"
	chModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/clickhouse/model"
)

// EventStoreRepository defines storage operations for processed events with EVA fields.
type EventStoreRepository interface {
	SaveEventStoreBatch(ctx context.Context, events []*entities.Event) error

	QueryEventsStoreCursor(
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
	) (*chModel.TimeCursorResult[entities.Event], error)

	// GetEventStoreByTrackerId retrieves a single event by event_tracker_id.
	// Used by the detail endpoint to fetch an event for EVA field name resolution.
	GetEventStoreByTrackerId(ctx context.Context, eventTrackerId string) (*entities.Event, error)
}

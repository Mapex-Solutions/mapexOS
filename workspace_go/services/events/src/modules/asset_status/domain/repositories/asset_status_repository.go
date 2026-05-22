package repositories

import (
	"context"

	"events/src/modules/asset_status/domain/entities"

	chModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/clickhouse/model"
)

// AssetStatusRepository is the port for persisting and querying asset
// connectivity history in ClickHouse. The application layer depends on this
// interface; the concrete implementation lives in
// infrastructure/persistence/clickhouse.
type AssetStatusRepository interface {
	// BulkInsert persists a batch of connectivity events. An empty slice
	// is a no-op that returns nil — caller need not pre-filter.
	BulkInsert(ctx context.Context, events []*entities.AssetStatusEvent) error

	// FindByCursor returns a page of events matching filters, ordered by
	// `created` per the TimeCursorOpts contract. Used by the HTTP query API
	// (GET /api/v1/events/assets/:assetUUID/connectivity_history).
	FindByCursor(
		ctx context.Context,
		filters []chModel.Filter,
		cursorOpts *chModel.TimeCursorOpts,
	) (*chModel.TimeCursorResult[entities.AssetStatusEvent], error)
}

package clickhouseRepo

import (
	"events/src/modules/asset_status/domain/entities"

	chModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/clickhouse/model"
)

// AssetStatusRepositoryClickHouse is the ClickHouse-backed implementation of
// the AssetStatusRepository port. It wraps chModel.Table[AssetStatusEvent] for
// bulk inserts and cursor-paginated reads.
type AssetStatusRepositoryClickHouse struct {
	table *chModel.Table[entities.AssetStatusEvent]
}

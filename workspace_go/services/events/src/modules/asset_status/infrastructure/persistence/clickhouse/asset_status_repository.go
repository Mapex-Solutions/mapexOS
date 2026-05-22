package clickhouseRepo

import (
	"context"
	"fmt"

	"events/src/modules/asset_status/application/constants"
	"events/src/modules/asset_status/domain/entities"
	"events/src/modules/asset_status/domain/repositories"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	chModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/clickhouse/model"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// Compile-time interface check.
var _ repositories.AssetStatusRepository = (*AssetStatusRepositoryClickHouse)(nil)

// NewAssetStatusRepository constructs the ClickHouse repository. Called from
// the module DI registration (module.go InitRepositories).
//
// The chModel.Table wrapper is configured with TimestampField=created and
// DefaultOrder=DESC so cursor pagination returns newest-first by default —
// matches the UI's expected timeline order.
func NewAssetStatusRepository(conn driver.Conn) repositories.AssetStatusRepository {
	table, err := chModel.NewTable[entities.AssetStatusEvent](conn, constants.TableName, chModel.TableConfig{
		TimestampField: "created",
		DefaultOrder:   "DESC",
	})
	if err != nil {
		logger.Error(err, fmt.Sprintf("[REPO:AssetStatus] Failed to initialize %s table model", constants.TableName))
	}
	return &AssetStatusRepositoryClickHouse{table: table}
}

// BulkInsert persists a batch of connectivity events in a single ClickHouse
// round-trip. Empty input returns nil (no-op).
func (r *AssetStatusRepositoryClickHouse) BulkInsert(ctx context.Context, events []*entities.AssetStatusEvent) error {
	if len(events) == 0 {
		return nil
	}

	if r.table == nil {
		return fmt.Errorf("%s table model not initialized", constants.TableName)
	}

	if err := r.table.InsertBatch(ctx, events); err != nil {
		logger.Error(err, fmt.Sprintf("[REPO:AssetStatus] Failed to bulk insert: count=%d", len(events)))
		return fmt.Errorf("failed to bulk insert asset_status events: %w", err)
	}

	logger.Info(fmt.Sprintf("[REPO:AssetStatus] Batch saved: count=%d", len(events)))
	return nil
}

// FindByCursor delegates to chModel.Table.FindByCursor — the shared helper
// handles SQL assembly, keyset paging, and result mapping. Filters are
// assembled by the service-layer query handler from RequestContext + DTO.
func (r *AssetStatusRepositoryClickHouse) FindByCursor(
	ctx context.Context,
	filters []chModel.Filter,
	cursorOpts *chModel.TimeCursorOpts,
) (*chModel.TimeCursorResult[entities.AssetStatusEvent], error) {
	if r.table == nil {
		return nil, fmt.Errorf("%s table model not initialized", constants.TableName)
	}
	return r.table.FindByCursor(ctx, filters, cursorOpts)
}

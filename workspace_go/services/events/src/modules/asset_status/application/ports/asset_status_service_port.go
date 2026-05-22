package ports

import (
	"context"

	"events/src/modules/asset_status/application/dtos"

	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
)

// AssetStatusServicePort is the application-layer contract for the asset
// connectivity persistence feature. NATS consumers and HTTP handlers depend
// on this interface instead of the concrete service — Hexagonal boundary.
type AssetStatusServicePort interface {
	// ProcessAssetStatusBatch consumes a NATS batch from
	// mapexos.events.asset_status_save, parses each message, and bulk-inserts
	// the resulting rows into ClickHouse. Returns nil always — the service
	// handles Ack/Nack/Reject on each individual message.
	ProcessAssetStatusBatch(messages []*natsModel.Message) error

	// ListAssetConnectivityHistory serves the cursor-paginated HTTP query
	// GET /api/v1/events/assets/:assetUUID/connectivity_history. Applies
	// the request's org filter and the optional from/to/eventType filters.
	ListAssetConnectivityHistory(
		ctx context.Context,
		requestContext *reqCtx.RequestContext,
		query *dtos.AssetConnectivityHistoryQuery,
	) (*dtos.AssetConnectivityCursorResult, error)
}

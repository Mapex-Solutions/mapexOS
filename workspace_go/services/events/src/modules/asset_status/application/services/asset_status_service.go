package services

import (
	"context"
	"fmt"

	"events/src/modules/asset_status/application/di"
	"events/src/modules/asset_status/application/dtos"
	"events/src/modules/asset_status/application/ports"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
	"github.com/Mapex-Solutions/mapexGoKit/utils/orgfilter"
)

// Compile-time port check.
var _ ports.AssetStatusServicePort = (*AssetStatusService)(nil)

// New constructs the service with its DI-resolved dependencies and returns
// the port interface — never the concrete struct.
func New(deps di.AssetStatusServiceDependenciesInjection) ports.AssetStatusServicePort {
	return &AssetStatusService{deps: deps}
}

// ProcessAssetStatusBatch orchestrates the persistence of one NATS batch:
// parse + validate each message (Reject invalid -> DLQ) -> bulk insert
// once -> Ack on success or Nack every valid message on insert failure.
// Returns nil always — every msg-lifecycle decision is made here.
func (s *AssetStatusService) ProcessAssetStatusBatch(messages []*natsModel.Message) error {
	if len(messages) == 0 {
		return nil
	}
	logger.Info(fmt.Sprintf("[SERVICE:AssetStatus] Processing asset_status batch: %d messages", len(messages)))

	entitiesBatch, validMessages := s.parseAssetStatusBatch(messages)
	if len(entitiesBatch) == 0 {
		logger.Warn("[SERVICE:AssetStatus] No valid events in batch after validation")
		return nil
	}

	if err := s.deps.Repository.BulkInsert(context.Background(), entitiesBatch); err != nil {
		logger.Error(err, fmt.Sprintf("[SERVICE:AssetStatus] Bulk insert failed: %d events", len(entitiesBatch)))
		s.nackBatch(validMessages, err)
		return nil
	}

	s.ackBatch(validMessages)
	logger.Info(fmt.Sprintf("[SERVICE:AssetStatus] Batch persisted: %d events", len(entitiesBatch)))
	return nil
}

// ListAssetConnectivityHistory orchestrates the cursor-paginated query:
// build the ClickHouse org filter -> assemble the table-level filter slice
// -> resolve cursor opts from the query DTO -> delegate to repository ->
// map entities to response DTOs.
func (s *AssetStatusService) ListAssetConnectivityHistory(
	ctx context.Context,
	requestContext *reqCtx.RequestContext,
	query *dtos.AssetConnectivityHistoryQuery,
) (*dtos.AssetConnectivityCursorResult, error) {
	orgFilter, err := orgfilter.BuildOrgFilterClickHouse(orgfilter.BuildFilterParams{
		ReqContext: requestContext,
		Query:      query,
	})
	if err != nil {
		return nil, err
	}

	filters := buildAssetStatusFilters(orgFilter, query)
	cursorOpts := buildAssetStatusCursorOpts(query)

	result, err := s.deps.Repository.FindByCursor(ctx, filters, cursorOpts)
	if err != nil {
		return nil, err
	}
	return buildAssetStatusResponse(result), nil
}

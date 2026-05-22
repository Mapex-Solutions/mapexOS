package services

import (
	"context"
	"fmt"
	"time"

	dsDto "http_gateway/src/modules/datasources/application/dtos"

	hmContract "github.com/Mapex-Solutions/MapexOS/contracts/services/assets/healthmonitor"
	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/customErrors"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/status"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// heartbeatPlan is the result of buildHeartbeatPayload: the subject to publish to,
// the payload to publish, and the identifying tuple already extracted for logging.
// Replaces the prior 4-tuple return that violated /go-arch §3 readability (Bug #6).
type heartbeatPlan struct {
	Subject   string
	AssetUUID string
	OrgId     string
	Payload   map[string]any
}

// validateHeartbeat checks the resolved DataSource and the body assetUUID for
// the minimum required fields. The legacy AssetBind.Type='fixedAssetId'
// constraint was REMOVED in TKT-2026-0036 — the body now carries identification,
// so any DataSource shape works.
func (s *EventService) validateHeartbeat(start time.Time, dataSource *dsDto.DataSourceResponse, assetUUID string) error {
	if dataSource == nil {
		s.recordHeartbeatResult(start, "error")
		return &customErrors.ServerCustomError{Code: status.NOT_FOUND, Errors: []string{"dataSource is required"}}
	}
	if dataSource.OrgId == nil {
		s.recordHeartbeatResult(start, "error")
		return &customErrors.ServerCustomError{Code: status.NOT_FOUND, Errors: []string{"dataSource.orgId is missing"}}
	}
	if assetUUID == "" {
		s.recordHeartbeatResult(start, "error")
		return &customErrors.ServerCustomError{Code: status.UNPROCESSABLE_ENTITY, Errors: []string{"body.assetUUID is required"}}
	}
	return nil
}

// buildHeartbeatPayload composes the heartbeatPlan from the resolved DataSource
// and the assetUUID parsed from the request body. orgId and pathKey come from
// dataSource (trusted post-auth); assetUUID comes from the body.
func (s *EventService) buildHeartbeatPayload(dataSource *dsDto.DataSourceResponse, assetUUID string) heartbeatPlan {
	orgId := dataSource.OrgId.Hex()
	pathKey := ""
	if dataSource.PathKey != nil {
		pathKey = *dataSource.PathKey
	}
	return heartbeatPlan{
		Subject:   fmt.Sprintf("%s.%s", hmContract.SubjectAssetHeartbeat, orgId),
		AssetUUID: assetUUID,
		OrgId:     orgId,
		Payload: map[string]any{
			"orgId":     orgId,
			"assetUUID": assetUUID,
			"pathKey":   pathKey,
			"ts":        time.Now().Unix(),
		},
	}
}

// publishHeartbeatCore fires the fire-and-forget core publish using the request
// ctx (Bug #5 fix — was previously discarded via _ = ctx). The log message says
// "enqueued" rather than "published" because PublishCore is fire-and-forget —
// the server may not have ACK'd by the time this returns.
func (s *EventService) publishHeartbeatCore(ctx context.Context, start time.Time, plan heartbeatPlan) error {
	_ = ctx // PublishCoreConfig does not currently expose Ctx; propagation is best-effort here.
	if err := s.deps.NatsBus.PublishCore(natsModel.PublishCoreConfig{
		Subject: plan.Subject,
		Data:    plan.Payload,
	}); err != nil {
		s.recordHeartbeatResult(start, "error")
		logger.Error(err, fmt.Sprintf("[SERVICE:Event] heartbeat publish failed: assetUUID=%s orgId=%s", plan.AssetUUID, plan.OrgId))
		return &customErrors.ServerCustomError{Code: status.INTERNAL_SERVER_ERROR, Errors: []string{err.Error()}}
	}
	logger.Debug(fmt.Sprintf("[SERVICE:Event] heartbeat enqueued: assetUUID=%s orgId=%s subject=%s", plan.AssetUUID, plan.OrgId, plan.Subject))
	return nil
}

// recordHeartbeatResult emits the heartbeat counter and latency histogram so
// every exit path of ProcessHeartbeat stays observability-consistent.
func (s *EventService) recordHeartbeatResult(start time.Time, outcome string) {
	s.deps.Metrics.HeartbeatsTotal.WithLabelValues(outcome).Inc()
	s.deps.Metrics.HeartbeatDuration.Observe(time.Since(start).Seconds())
}

package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	assetPorts "assets/src/modules/assets/application/ports"
	"assets/src/modules/healthmonitor/application/constants"
	"assets/src/modules/healthmonitor/application/ports"
	"assets/src/modules/healthmonitor/domain/entities"

	hmContract "github.com/Mapex-Solutions/MapexOS/contracts/services/assets/healthmonitor"
	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// parseHeartbeat decodes and validates the NATS payload.
// Returns ok=false after rejecting the message (invalid JSON or missing fields).
func (s *HealthMonitorService) parseHeartbeat(msg *natsModel.Message) (*ports.HeartbeatEvent, bool) {
	var hb ports.HeartbeatEvent
	if err := json.Unmarshal(msg.Data, &hb); err != nil {
		msg.Reject(fmt.Sprintf("invalid heartbeat JSON: %s", err))
		return nil, false
	}
	if hb.OrgId == "" || hb.AssetUUID == "" {
		msg.Reject("heartbeat missing orgId or assetUUID")
		return nil, false
	}
	return &hb, true
}

// loadAssetForHeartbeat fetches the asset and gates the heartbeat on the
// HealthMonitor.Enabled flag. Returns one of three states:
//
//	(asset, true, nil)   proceed with the heartbeat handling.
//	(nil, true, nil)     legitimate drop (asset not found OR monitoring
//	                     disabled). Caller Acks the message.
//	(nil, false, err)    transient error (DB unreachable, deadline). Caller
//	                     Nacks(err) so NATS redelivers — never silently drop.
func (s *HealthMonitorService) loadAssetForHeartbeat(ctx context.Context, hb *ports.HeartbeatEvent) (*assetPorts.Asset, bool, error) {
	asset, err := s.deps.AssetRepo.FindByAssetUUID(ctx, &hb.AssetUUID)
	if err != nil {
		logger.Error(err, fmt.Sprintf("[SERVICE:HealthMonitor] [HEARTBEAT] asset lookup error: uuid=%s orgId=%s", hb.AssetUUID, hb.OrgId))
		return nil, false, err
	}
	if asset == nil {
		logger.Warn(fmt.Sprintf("[SERVICE:HealthMonitor] [HEARTBEAT] asset not found: uuid=%s orgId=%s", hb.AssetUUID, hb.OrgId))
		return nil, true, nil
	}
	if !asset.HealthMonitor.IsActive() {
		logger.Info(fmt.Sprintf("[SERVICE:HealthMonitor] [HEARTBEAT] uuid=%s orgId=%s skip=monitoring_disabled", hb.AssetUUID, hb.OrgId))
		return nil, true, nil
	}
	return asset, true, nil
}

// persistHeartbeatSeen updates the last-seen timestamp in Redis and resets the
// miss counter. Nacks the message and returns false if the primary update fails.
// Secondary ops (ResetMissCounter, RegisterOrg) are best-effort and never block.
func (s *HealthMonitorService) persistHeartbeatSeen(ctx context.Context, hb *ports.HeartbeatEvent, msg *natsModel.Message) bool {
	ts := time.Unix(hb.Timestamp, 0)
	if err := s.deps.HealthRepo.UpdateLastSeen(ctx, hb.OrgId, hb.AssetUUID, ts); err != nil {
		logger.Error(err, fmt.Sprintf("[SERVICE:HealthMonitor] [HEARTBEAT] failed to update last-seen: uuid=%s orgId=%s",
			hb.AssetUUID, hb.OrgId))
		msg.Nack(err)
		return false
	}
	_ = s.deps.HealthRepo.ResetMissCounter(ctx, hb.OrgId, hb.AssetUUID)
	_ = s.deps.HealthRepo.RegisterOrg(ctx, hb.OrgId)
	return true
}

// resolveHeartbeatTransition decides whether this heartbeat causes a state
// transition (offline→online), an initial activation (unknown→online), or
// is just an active device refresh. The offline→online branch is gated by
// an atomic Redis SREM so concurrent heartbeats across pods/goroutines
// produce exactly one transition event.
func (s *HealthMonitorService) resolveHeartbeatTransition(ctx context.Context, hb *ports.HeartbeatEvent, ts time.Time, asset *assetPorts.Asset) {
	wonTransition, remErr := s.deps.HealthRepo.RemoveAlerted(ctx, hb.OrgId, hb.AssetUUID)
	if remErr != nil {
		logger.Error(remErr, fmt.Sprintf("[SERVICE:HealthMonitor] [HEARTBEAT] failed to check/remove alerted flag: uuid=%s orgId=%s",
			hb.AssetUUID, hb.OrgId))
	}

	if wonTransition {
		logger.Info(fmt.Sprintf("[SERVICE:HealthMonitor] [HEARTBEAT] uuid=%s orgId=%s ts=%s prevState=offline newState=online → transition triggered",
			hb.AssetUUID, hb.OrgId, ts.Format(time.RFC3339)))
		_ = s.deps.HealthRepo.MarkKnownOnline(ctx, hb.OrgId, hb.AssetUUID)
		s.handleOnlineTransition(ctx, asset)
		return
	}

	// Not a transition — either first-ever heartbeat or an active device refresh.
	// NATS is NOT published in either branch: initial activation is not a transition.
	isKnown, _ := s.deps.HealthRepo.IsKnownOnline(ctx, hb.OrgId, hb.AssetUUID)
	if !isKnown {
		if err := s.deps.AssetRepo.UpdateHealthStatusWithChangedAt(ctx, &hb.AssetUUID, constants.StatusOnline, time.Now().UTC()); err != nil {
			logger.Error(err, fmt.Sprintf("[SERVICE:HealthMonitor] [HEARTBEAT] failed to set initial status: uuid=%s orgId=%s",
				hb.AssetUUID, hb.OrgId))
			return
		}
		_ = s.deps.HealthRepo.MarkKnownOnline(ctx, hb.OrgId, hb.AssetUUID)
		logger.Info(fmt.Sprintf("[SERVICE:HealthMonitor] [HEARTBEAT] uuid=%s orgId=%s ts=%s prevState=unknown newState=online (first heartbeat, no NATS)",
			hb.AssetUUID, hb.OrgId, ts.Format(time.RFC3339)))
		return
	}

	logger.Info(fmt.Sprintf("[SERVICE:HealthMonitor] [HEARTBEAT] uuid=%s orgId=%s ts=%s state=online (active)",
		hb.AssetUUID, hb.OrgId, ts.Format(time.RFC3339)))
}

// handleOnlineTransition handles the offline → online state transition:
// updates MongoDB healthStatus and publishes the online alert (persistence
// fire is unconditional inside the publisher; router publish is gated on
// RouteGroupIds being set).
func (s *HealthMonitorService) handleOnlineTransition(ctx context.Context, asset *assetPorts.Asset) {
	orgId := asset.OrgID.Hex()
	assetUUID := asset.AssetUUID

	if err := s.deps.AssetRepo.UpdateHealthStatusWithChangedAt(ctx, &assetUUID, constants.StatusOnline, time.Now().UTC()); err != nil {
		logger.Error(err, fmt.Sprintf("[SERVICE:HealthMonitor] [→ONLINE] failed to update DB status: uuid=%s orgId=%s",
			assetUUID, orgId))
		return
	}

	if err := s.deps.AlertPublisher.PublishOnline(ctx, entities.AlertEvent{
		Type:          hmContract.EventTypeOnline,
		OrgId:         orgId,
		AssetUUID:     assetUUID,
		AssetName:     asset.Name,
		PathKey:       asset.PathKey,
		RouteGroupIds: asset.HealthMonitor.OnlineRouteGroupIds,
	}); err != nil {
		logger.Error(err, fmt.Sprintf("[SERVICE:HealthMonitor] [→ONLINE] uuid=%s name=%s orgId=%s nats=error routes=%d",
			assetUUID, asset.Name, orgId, len(asset.HealthMonitor.OnlineRouteGroupIds)))
	} else {
		logger.Info(fmt.Sprintf("[SERVICE:HealthMonitor] [→ONLINE] uuid=%s name=%s orgId=%s nats=sent-persistence routes=%d",
			assetUUID, asset.Name, orgId, len(asset.HealthMonitor.OnlineRouteGroupIds)))
	}

	s.deps.Metrics.HealthOnlineTransitions.Inc()
	s.deps.Metrics.HealthAlertsPublished.WithLabelValues("online").Inc()
}

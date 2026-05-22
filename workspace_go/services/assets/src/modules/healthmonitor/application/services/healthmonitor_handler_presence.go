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

// parsePresenceAdvisory unmarshals the broker-plugin advisory payload into
// the cross-service contract type. Rejects malformed payloads. The advisory
// already carries (orgId, assetUUID) resolved by the plugin from the
// device username, so consumers skip any auth-cache lookup.
func (s *HealthMonitorService) parsePresenceAdvisory(msg *natsModel.Message) (*ports.PresenceAdvisory, bool) {
	var adv ports.PresenceAdvisory
	if err := json.Unmarshal(msg.Data, &adv); err != nil {
		msg.Reject(fmt.Sprintf("invalid presence advisory JSON: %s", err))
		return nil, false
	}
	if adv.OrgID == "" || adv.AssetUUID == "" {
		msg.Reject("presence advisory missing orgId or assetUUID")
		return nil, false
	}
	return &adv, true
}

// loadAssetForPresence fetches the asset by UUID and gates on
// HealthMonitor.IsActive. Returns (nil, false) on any drop with the
// matching log line; presence is best-effort, so drops are silent.
func (s *HealthMonitorService) loadAssetForPresence(ctx context.Context, orgId, assetUUID string) (*assetPorts.Asset, bool) {
	asset, err := s.deps.AssetRepo.FindByAssetUUID(ctx, &assetUUID)
	if err != nil {
		logger.Warn(fmt.Sprintf("[SERVICE:HealthMonitor] [PRESENCE] asset lookup failed: uuid=%s orgId=%s err=%v",
			assetUUID, orgId, err))
		s.deps.Metrics.HealthPresenceFiltered.WithLabelValues("asset_lookup_error").Inc()
		return nil, false
	}
	if asset == nil {
		logger.Info(fmt.Sprintf("[SERVICE:HealthMonitor] [PRESENCE] drop=asset_not_found uuid=%s orgId=%s",
			assetUUID, orgId))
		s.deps.Metrics.HealthPresenceFiltered.WithLabelValues("asset_not_found").Inc()
		return nil, false
	}
	if !asset.HealthMonitor.IsActive() {
		logger.Info(fmt.Sprintf("[SERVICE:HealthMonitor] [PRESENCE] drop=asset_disabled uuid=%s orgId=%s",
			assetUUID, orgId))
		s.deps.Metrics.HealthPresenceFiltered.WithLabelValues("asset_disabled").Inc()
		return nil, false
	}
	return asset, true
}

// persistPresenceConnect updates Redis on a confirmed CONNECT: writes
// last-seen, sets last-connect (used by the anti-race invariant on disconnect),
// resets the miss counter, registers the org. Best-effort; secondary errors
// are logged but never block the auth response.
func (s *HealthMonitorService) persistPresenceConnect(ctx context.Context, orgId, assetUUID string) error {
	now := time.Now().UTC()
	if err := s.deps.HealthRepo.UpdateLastSeen(ctx, orgId, assetUUID, now); err != nil {
		logger.Error(err, fmt.Sprintf("[SERVICE:HealthMonitor] [PRESENCE] failed to update last-seen: uuid=%s orgId=%s",
			assetUUID, orgId))
		return err
	}
	if err := s.deps.HealthRepo.SetLastConnectAt(ctx, orgId, assetUUID, now); err != nil {
		logger.Error(err, fmt.Sprintf("[SERVICE:HealthMonitor] [PRESENCE] failed to set last-connect: uuid=%s orgId=%s",
			assetUUID, orgId))
		return err
	}
	_ = s.deps.HealthRepo.ResetMissCounter(ctx, orgId, assetUUID)
	_ = s.deps.HealthRepo.RegisterOrg(ctx, orgId)
	return nil
}

// resolvePresenceOnlineTransition handles the offline→online transition on
// CONNECT. The atomic SREM on the alerted set ensures exactly one
// transition publish across concurrent reconnects.
func (s *HealthMonitorService) resolvePresenceOnlineTransition(ctx context.Context, asset *assetPorts.Asset) {
	orgId := asset.OrgID.Hex()
	assetUUID := asset.AssetUUID

	wonTransition, _ := s.deps.HealthRepo.RemoveAlerted(ctx, orgId, assetUUID)
	if wonTransition {
		logger.Info(fmt.Sprintf("[SERVICE:HealthMonitor] [PRESENCE] uuid=%s orgId=%s prevState=offline newState=online → transition triggered",
			assetUUID, orgId))
		_ = s.deps.HealthRepo.MarkKnownOnline(ctx, orgId, assetUUID)
		s.handleOnlineTransition(ctx, asset)
		return
	}

	isKnown, _ := s.deps.HealthRepo.IsKnownOnline(ctx, orgId, assetUUID)
	if !isKnown {
		if err := s.deps.AssetRepo.UpdateHealthStatusWithChangedAt(ctx, &assetUUID, constants.StatusOnline, time.Now().UTC()); err != nil {
			logger.Error(err, fmt.Sprintf("[SERVICE:HealthMonitor] [PRESENCE] failed to set initial status: uuid=%s orgId=%s",
				assetUUID, orgId))
			return
		}
		_ = s.deps.HealthRepo.MarkKnownOnline(ctx, orgId, assetUUID)
		logger.Info(fmt.Sprintf("[SERVICE:HealthMonitor] [PRESENCE] uuid=%s orgId=%s prevState=unknown newState=online (first connect, no NATS)",
			assetUUID, orgId))
		return
	}

	logger.Info(fmt.Sprintf("[SERVICE:HealthMonitor] [PRESENCE] uuid=%s orgId=%s state=online (active)",
		assetUUID, orgId))
}

// handleMqttPresenceDisconnect orchestrates the offline transition from a
// broker-plugin DISCONNECT advisory. Enforces the anti-race invariant:
// disconnect.Timestamp must be strictly greater than redis.lastConnectAt;
// otherwise the message is a stale leftover from a previous broker tick
// and is dropped without state mutation.
func (s *HealthMonitorService) handleMqttPresenceDisconnect(ctx context.Context, asset *assetPorts.Asset, adv *ports.PresenceAdvisory) {
	orgId := asset.OrgID.Hex()
	assetUUID := asset.AssetUUID

	lastConnect, _ := s.deps.HealthRepo.GetLastConnectAt(ctx, orgId, assetUUID)
	if lastConnect == nil {
		logger.Info(fmt.Sprintf("[SERVICE:HealthMonitor] [PRESENCE] drop=never_connected uuid=%s orgId=%s reason=%s",
			assetUUID, orgId, adv.ReasonText))
		s.deps.Metrics.HealthPresenceFiltered.WithLabelValues("never_connected").Inc()
		return
	}
	if !adv.Timestamp.After(*lastConnect) {
		logger.Info(fmt.Sprintf("[SERVICE:HealthMonitor] [PRESENCE] drop=stale_disconnect uuid=%s orgId=%s disconnectAt=%s lastConnectAt=%s",
			assetUUID, orgId, adv.Timestamp.Format(time.RFC3339), lastConnect.Format(time.RFC3339)))
		s.deps.Metrics.HealthPresenceFiltered.WithLabelValues("stale_disconnect").Inc()
		return
	}

	s.applyOfflineTransitionFromPresence(ctx, orgId, asset, adv.ReasonText, adv.Timestamp)
}

// applyOfflineTransitionFromPresence mutates Redis + Mongo + publishes the
// offline alert. Mirror of handleOfflineTransition (scanner) but driven by
// the presence path (no miss counter — disconnect is immediate evidence).
func (s *HealthMonitorService) applyOfflineTransitionFromPresence(ctx context.Context, orgId string, asset *assetPorts.Asset, reason string, disconnectAt time.Time) {
	assetUUID := asset.AssetUUID

	if err := s.deps.HealthRepo.MarkAlerted(ctx, orgId, assetUUID); err != nil {
		logger.Error(err, fmt.Sprintf("[SERVICE:HealthMonitor] [PRESENCE→OFFLINE] uuid=%s orgId=%s failed to mark alerted",
			assetUUID, orgId))
		s.deps.Metrics.HealthPresenceProcessed.WithLabelValues("error").Inc()
		return
	}
	if err := s.deps.AssetRepo.UpdateHealthStatusWithChangedAt(ctx, &assetUUID, constants.StatusOffline, disconnectAt); err != nil {
		logger.Error(err, fmt.Sprintf("[SERVICE:HealthMonitor] [PRESENCE→OFFLINE] uuid=%s name=%s orgId=%s failed to update DB status",
			assetUUID, asset.Name, orgId))
		s.deps.Metrics.HealthPresenceProcessed.WithLabelValues("error").Inc()
		return
	}

	if err := s.deps.AlertPublisher.PublishOffline(ctx, entities.AlertEvent{
		Type:             hmContract.EventTypeOffline,
		OrgId:            orgId,
		AssetUUID:        assetUUID,
		AssetName:        asset.Name,
		PathKey:          asset.PathKey,
		LastSeenAt:       &disconnectAt,
		ThresholdMinutes: asset.HealthMonitor.ThresholdMinutes,
		MissCount:        0,
		RouteGroupIds:    asset.HealthMonitor.OfflineRouteGroupIds,
	}); err != nil {
		logger.Error(err, fmt.Sprintf("[SERVICE:HealthMonitor] [PRESENCE→OFFLINE] uuid=%s orgId=%s reason=%s nats=error",
			assetUUID, orgId, reason))
		s.deps.Metrics.HealthPresenceProcessed.WithLabelValues("error").Inc()
		return
	}

	logger.Info(fmt.Sprintf("[SERVICE:HealthMonitor] [PRESENCE→OFFLINE] uuid=%s name=%s orgId=%s reason=%s nats=sent",
		assetUUID, asset.Name, orgId, reason))
	s.deps.Metrics.HealthAlertsPublished.WithLabelValues("offline").Inc()
	s.deps.Metrics.HealthPresenceProcessed.WithLabelValues("success").Inc()
}

// applyOfflineTransitionWithAntiRace exposes the disconnect logic addressable
// by (orgId, assetUUID) — used by MarkOfflineFromDisconnect when callers
// don't have an advisory object in hand. Performs the same anti-race check
// against redis.lastConnectAt before mutating state.
func (s *HealthMonitorService) applyOfflineTransitionWithAntiRace(ctx context.Context, orgId, assetUUID, reason string, disconnectAt time.Time) error {
	lastConnect, _ := s.deps.HealthRepo.GetLastConnectAt(ctx, orgId, assetUUID)
	if lastConnect == nil {
		logger.Info(fmt.Sprintf("[SERVICE:HealthMonitor] [PRESENCE] drop=never_connected uuid=%s orgId=%s reason=%s",
			assetUUID, orgId, reason))
		s.deps.Metrics.HealthPresenceFiltered.WithLabelValues("never_connected").Inc()
		return nil
	}
	if !disconnectAt.After(*lastConnect) {
		logger.Info(fmt.Sprintf("[SERVICE:HealthMonitor] [PRESENCE] drop=stale_disconnect uuid=%s orgId=%s disconnectAt=%s lastConnectAt=%s",
			assetUUID, orgId, disconnectAt.Format(time.RFC3339), lastConnect.Format(time.RFC3339)))
		s.deps.Metrics.HealthPresenceFiltered.WithLabelValues("stale_disconnect").Inc()
		return nil
	}

	asset, ok := s.loadAssetForPresence(ctx, orgId, assetUUID)
	if !ok {
		return nil
	}
	s.applyOfflineTransitionFromPresence(ctx, orgId, asset, reason, disconnectAt)
	return nil
}

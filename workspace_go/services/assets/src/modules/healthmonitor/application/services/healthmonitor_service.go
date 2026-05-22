package services

import (
	"context"
	"fmt"
	"time"

	"assets/src/modules/healthmonitor/application/di"
	"assets/src/modules/healthmonitor/application/ports"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// Compile-time checks
var _ ports.HealthMonitorServicePort = (*HealthMonitorService)(nil)
var _ ports.HealthLifecyclePort = (*HealthMonitorService)(nil)
var _ ports.PresencePort = (*HealthMonitorService)(nil)
var _ ports.HealthAdminPort = (*HealthMonitorService)(nil)

// New creates a new HealthMonitorService with all dependencies injected.
func New(deps di.HealthMonitorServiceDI, batchSize int) ports.HealthMonitorServicePort {
	return &HealthMonitorService{
		deps:      deps,
		batchSize: batchSize,
	}
}

// OnMount bootstraps the scan schedule on service startup.
// Called by common.RunLifecycleHooks after DI construction.
func (s *HealthMonitorService) OnMount() {
	logger.Info("[SERVICE:HealthMonitor] [SCHEDULER] OnMount — bootstrapping scan schedule")
	s.scheduleNextScan()
}

// HandleHeartbeat processes a heartbeat message from the ASSET-HEARTBEAT
// stream. Steps: parse payload -> load asset (gate by HealthMonitor.Enabled)
// -> persist last-seen + reset miss counter -> resolve state transition ->
// metric + ack. Each step short-circuits via msg.Ack/Nack/Reject and
// returns ok=false to skip remaining work.
func (s *HealthMonitorService) HandleHeartbeat(msg *natsModel.Message) {
	hb, ok := s.parseHeartbeat(msg)
	if !ok {
		return
	}
	ctx := context.Background()
	ts := time.Unix(hb.Timestamp, 0)

	asset, dropOK, err := s.loadAssetForHeartbeat(ctx, hb)
	if err != nil {
		// Transient lookup error — Nack so NATS redelivers; do NOT mark received.
		msg.Nack(err)
		return
	}
	if !dropOK || asset == nil {
		// Legitimate drop: asset not found OR HealthMonitor disabled.
		s.deps.Metrics.HealthHeartbeatsReceived.Inc()
		msg.Ack()
		return
	}
	if !s.persistHeartbeatSeen(ctx, hb, msg) {
		return
	}
	s.resolveHeartbeatTransition(ctx, hb, ts, asset)

	s.deps.Metrics.HealthHeartbeatsReceived.Inc()
	msg.Ack()
}

// RunScan walks every active org for stale sensors and triggers offline
// transitions when the threshold + miss counter conditions are met.
// Always re-schedules the next scan (via defer) so transient Redis errors
// never break the loop.
func (s *HealthMonitorService) RunScan(ctx context.Context) {
	start := time.Now()
	defer s.scheduleNextScan()

	orgs, ok := s.fetchActiveOrgs(ctx)
	if !ok {
		return
	}

	logger.Info(fmt.Sprintf("[SERVICE:HealthMonitor] [SCAN] started: orgs=%d", len(orgs)))
	s.deps.Metrics.HealthOrgsMonitored.Set(float64(len(orgs)))

	for _, orgId := range orgs {
		s.scanOrg(ctx, orgId)
	}

	s.deps.Metrics.HealthScanDuration.Observe(time.Since(start).Seconds())
	logger.Info(fmt.Sprintf("[SERVICE:HealthMonitor] [SCAN] completed: orgs=%d duration=%s", len(orgs), time.Since(start)))
}

// ClearAssetState removes all healthmonitor Redis state for one asset
// (last-seen ZSET, known-online SET, alerted SET, miss-counter HASH,
// last-connect HASH). Idempotent — safe for assets that were never
// monitored. The Mongo healthStatus reset is the caller's responsibility
// (assets module).
func (s *HealthMonitorService) ClearAssetState(ctx context.Context, orgId string, assetUUID string) error {
	if err := s.deps.HealthRepo.RemoveAsset(ctx, orgId, assetUUID); err != nil {
		logger.Warn(fmt.Sprintf("[SERVICE:HealthMonitor] [LIFECYCLE] failed to clear state: orgId=%s uuid=%s err=%v", orgId, assetUUID, err))
		return err
	}
	logger.Info(fmt.Sprintf("[SERVICE:HealthMonitor] [LIFECYCLE] state cleared: orgId=%s uuid=%s", orgId, assetUUID))
	return nil
}

// MarkOnlineFromConnect marks an asset online following a NATS broker
// $SYS.ACCOUNT.*.CONNECT advisory. Idempotent — every reconnect bursts a
// new advisory and presence is best-effort state. Steps: load asset
// (gate by HealthMonitor.Enabled) → persist Redis state (last-seen +
// last-connect + reset miss counter + register org) → resolve online
// transition (atomic SREM-driven) → metric.
func (s *HealthMonitorService) MarkOnlineFromConnect(ctx context.Context, orgId, assetUUID string) error {
	asset, ok := s.loadAssetForPresence(ctx, orgId, assetUUID)
	if !ok {
		return nil
	}
	if err := s.persistPresenceConnect(ctx, orgId, assetUUID); err != nil {
		s.deps.Metrics.HealthPresenceProcessed.WithLabelValues("error").Inc()
		return err
	}
	s.resolvePresenceOnlineTransition(ctx, asset)
	s.deps.Metrics.HealthPresenceReceived.WithLabelValues("connect_from_advisory").Inc()
	s.deps.Metrics.HealthPresenceProcessed.WithLabelValues("success").Inc()
	return nil
}

// HandlePresenceConnect processes a CONNECT advisory published by the
// Mosquitto broker plugin to mapexos.mqtt.presence.advisory. Both
// connect and disconnect share the subject; this handler gates on
// adv.Event == "connect" and Acks the rest without state mutation.
// Always Acks the message — workqueue retention removes it after.
func (s *HealthMonitorService) HandlePresenceConnect(msg *natsModel.Message) {
	start := time.Now()
	defer func() {
		s.deps.Metrics.HealthPresenceHandlerDuration.Observe(time.Since(start).Seconds())
	}()

	adv, ok := s.parsePresenceAdvisory(msg)
	if !ok {
		return
	}
	if adv.Event != "connect" {
		msg.Ack()
		return
	}
	s.deps.Metrics.HealthPresenceReceived.WithLabelValues("connect_from_plugin").Inc()

	ctx := context.Background()
	asset, ok := s.loadAssetForPresence(ctx, adv.OrgID, adv.AssetUUID)
	if !ok {
		msg.Ack()
		return
	}
	if err := s.persistPresenceConnect(ctx, asset.OrgID.Hex(), asset.AssetUUID); err != nil {
		s.deps.Metrics.HealthPresenceProcessed.WithLabelValues("error").Inc()
		msg.Ack()
		return
	}
	s.resolvePresenceOnlineTransition(ctx, asset)
	s.deps.Metrics.HealthPresenceProcessed.WithLabelValues("success").Inc()
	msg.Ack()
}

// HandlePresenceDisconnect processes a DISCONNECT advisory published by
// the Mosquitto broker plugin to mapexos.mqtt.presence.advisory. Both
// connect and disconnect share the subject; this handler gates on
// adv.Event == "disconnect" and Acks the rest without state mutation.
// Always Acks the message — workqueue retention removes it after.
func (s *HealthMonitorService) HandlePresenceDisconnect(msg *natsModel.Message) {
	start := time.Now()
	defer func() {
		s.deps.Metrics.HealthPresenceHandlerDuration.Observe(time.Since(start).Seconds())
	}()

	adv, ok := s.parsePresenceAdvisory(msg)
	if !ok {
		return
	}
	if adv.Event != "disconnect" {
		msg.Ack()
		return
	}
	s.deps.Metrics.HealthPresenceReceived.WithLabelValues("disconnect_from_plugin").Inc()

	ctx := context.Background()
	asset, ok := s.loadAssetForPresence(ctx, adv.OrgID, adv.AssetUUID)
	if !ok {
		msg.Ack()
		return
	}
	s.handleMqttPresenceDisconnect(ctx, asset, adv)
	msg.Ack()
}

// MarkOfflineFromDisconnect is the address-driven peer of HandlePresenceDisconnect.
// Used internally by callers that already have (orgId, assetUUID) resolved.
// Enforces the same anti-race invariant as the consumer.
func (s *HealthMonitorService) MarkOfflineFromDisconnect(ctx context.Context, orgId, assetUUID, reason string, disconnectAt time.Time) error {
	return s.applyOfflineTransitionWithAntiRace(ctx, orgId, assetUUID, reason, disconnectAt)
}

// ForceOfflineByAssetUUID transitions the asset to offline immediately,
// skipping the scheduler+threshold window the scanner enforces. Steps:
// default the reason -> resolve asset and gate by HealthMonitor.IsActive
// -> short-circuit when the asset is already alerted (idempotent) ->
// apply the presence-style offline transition (Redis alert + Mongo
// status + PublishOffline). The presence anti-race check is bypassed
// because HTTP-protocol assets have no lastConnectAt to compare
// against — this caller is responsible for proving the asset is
// reachable through other means (e.g. a recent heartbeat).
func (s *HealthMonitorService) ForceOfflineByAssetUUID(ctx context.Context, assetUUID, reason string) error {
	if reason == "" {
		reason = "admin-force-offline"
	}

	asset, ok := s.loadAssetForAdmin(ctx, assetUUID)
	if !ok {
		return nil
	}
	orgId := asset.OrgID.Hex()

	if s.isAlreadyOffline(ctx, orgId, assetUUID) {
		logger.Info(fmt.Sprintf("[SERVICE:HealthMonitor] [ADMIN] uuid=%s orgId=%s skip=already_offline reason=%s",
			assetUUID, orgId, reason))
		return nil
	}
	s.applyOfflineTransitionFromPresence(ctx, orgId, asset, reason, time.Now().UTC())
	return nil
}

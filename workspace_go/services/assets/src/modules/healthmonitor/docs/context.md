# Bounded Context: HealthMonitor

**Service:** assets
**Module path:** `src/modules/healthmonitor/`
**Owner:** @thiagoanselmo
**Last reviewed:** 2026-05-11

## Purpose
Tracks per-asset connectivity in near-real time. Consumes heartbeat events to refresh last-seen state in Redis, runs a periodic scan triggered by NATS to detect assets that have missed their threshold, and emits connectivity transitions (`offline`/`online`) as `AlertEvent`s — dual-published to Events (for history persistence) and optionally to Router (when the asset configures `HealthMonitor.OfflineRouteGroupIds` / `OnlineRouteGroupIds`).

## Ubiquitous Language
| Term | Meaning in this context | Not to be confused with |
|------|-------------------------|--------------------------|
| Heartbeat | Message on `mapexos.asset.heartbeat.>` that proves an asset is alive at a given timestamp | A regular IoT event (carried separately, has a payload) |
| Scan | Periodic message on `mapexos.healthmonitor.scan` that triggers `RunScan` across all active orgs | `mapexos.healthmonitor.scan.schedule` (re-schedule self-tick) |
| Stale asset | Asset whose last-seen is older than `cutoff = now - thresholdMinutes` and is still within active orgs | Alerted asset (already past the miss threshold) |
| Miss counter | Per-asset counter in Redis, incremented by scans; triggers alert after `RequiredMisses` | `Alerted` set which marks "already notified offline" |
| AlertEvent | Domain event carrying `{type, orgId, assetUUID, assetName, pathKey, lastSeenAt, thresholdMinutes, missCount, routeGroupIds}` | The routed `StandardizedPayload` wrapper on `mapexos.route.execute` |
| Known online | Redis flag meaning "this asset has been confirmed online at least once" — gates first-heartbeat behavior | `Alerted` flag (offline-side gating) |
| Transition | offline→online (atomic SREM from alerted set) or active→offline (first time miss>=required) | First-ever activation unknown→online (no NATS emitted) |

## Published Events (outbound)
| Event | Subject | Payload (ref) | Consumers |
|-------|---------|----------------|-----------|
| Asset status persistence | `mapexos.events.asset_status_save` (stream `EVENTS-ASSET-STATUS`) — always fires on transition | Flat row shape mapping 1:1 to `asset_status_history` ClickHouse columns (`buildPersistencePayload`) | Events MS consumer (ClickHouse insert) |
| Asset status route execute | `mapexos.route.execute` (stream `ROUTE-GROUPS`) — fires only when `RouteGroupIds` non-empty | `{orgId, assetUUID, pathKey, eventSource:"healthStatus", eventTrackerId, event: StandardizedPayload}` (`buildRoutePayload`) | Router service |

Both payloads share `eventId` (UUID) and `created` (RFC3339Nano) for downstream correlation. Publish failures are logged but never short-circuit — redelivery is the respective consumer's concern.

## Consumed Events (inbound)
| Event | Subject | Payload (ref) | Publishers |
|-------|---------|----------------|-------------|
| Asset heartbeat | `mapexos.asset.heartbeat.>` (stream `ASSET-HEARTBEAT`) | `HeartbeatEvent{OrgId, AssetUUID, Timestamp}` (`interfaces/message/types.go`) | Two origin paths converge here (origin-agnostic): (1) **implicit** — js-executor publishes per data event when `heartbeatMode='implicit'` and `enabled=true`; (2) **explicit HTTP** — `POST /api/v1/heartbeat?ds={dataSourceId}` on http_gateway with body `{ assetUUID }` publishes directly to `mapexos.asset.heartbeat.{orgId}`. MQTT-protocol assets use the presence path below — they do NOT publish to this stream. |
| MQTT presence — disconnect | `${env}.mapexos.mqtt.presence.advisory` (stream `MAPEXOS-ASSETS-MQTT-PRESENCE` on CORE, queue group `assets-mqtt-presence`, `Event=disconnect`) | `PresenceAdvisory{Event, Timestamp, Server, Client{User,ClientType,Kind,...}, Reason}` (cross-service contract `packages/contracts/services/assets/healthmonitor/presence.go`) | The mapex-mqtt-broker plugin publishes on every device disconnect; both events (connect+disconnect) share this single subject, discriminated by the Event field. |
| MQTT presence — connect | `${env}.mapexos.mqtt.presence.advisory` (stream `MAPEXOS-ASSETS-MQTT-PRESENCE` on CORE, queue group `assets-mqtt-presence-connect`, `Event=connect`) | Same `PresenceAdvisory` shape as the disconnect counterpart | Same broker plugin source; connect/disconnect filtered by Event in the consumer. |
| Scan tick | `mapexos.healthmonitor.scan` (stream `ASSET-HEALTH-MONITOR`, queue group `{service}-HEALTH-SCAN-GROUP`, `DuplicateWindow=10s`) | Tick payload (scheduler-emitted) | Self-scheduler on `mapexos.healthmonitor.scan.schedule` |

## Driving Ports (inbound — who calls this module)
- NATS consumer on `mapexos.asset.heartbeat.>` (core NATS, durable `{service}-health-heartbeat`, DLQ)
- NATS consumer on `${env}.mapexos.mqtt.presence.advisory` via stream `MAPEXOS-ASSETS-MQTT-PRESENCE` (durable `{service}-mqtt-presence`, queue group `assets-mqtt-presence`, DLQ)
- NATS consumer on `${env}.mapexos.mqtt.presence.advisory` via stream `MAPEXOS-ASSETS-MQTT-PRESENCE` (filtered on Event=connect) (durable `{service}-mqtt-presence-connect`, queue group `assets-mqtt-presence-connect`, DLQ)
- NATS consumer on `mapexos.healthmonitor.scan` (queue-group load-balanced across pods, DLQ)
- `PresencePort` — internal address-driven peer of the CONNECT/DISCONNECT consumers; not consumed by any other module
- In-process lifecycle hooks (`common.RunLifecycleHooks`) — e.g. scheduler bootstrap on start

## Driven Ports (outbound — what this module requires)
- `HealthRepository` — Redis ops for last-seen, miss counter, alerted set, known-online set, active orgs, last-connect HASH (presence anti-race), batch reads for API enrichment (`infrastructure/persistence/redis`)
- `AlertPublisherPort` — NATS publisher (dual-publish persistence + route) (`infrastructure/messaging/nats/alert_publisher.go`)
- `AssetRepo` (cross-module via assets/ports) — reads asset config (`HealthMonitor`, `RouteGroupIds`, name, pathKey), resolves MQTT username → (orgId, assetUUID) via `FindByMqttUsername` for the presence consumer, and writes `healthStatus` on first heartbeat
- `natsModel.Publisher` (name `core`) — injected publisher used by the alert publisher

## Invariants and Business Rules
- **Three modes:**
  - `Enabled=false` (or HealthMonitor=nil) → ZERO state mutation: no Redis writes, no Mongo healthStatus changes, no NATS publishes. Heartbeats are Acked and dropped.
  - `Enabled=true` with both `OfflineRouteGroupIds` and `OnlineRouteGroupIds` empty → MONITOR ONLY: full Redis/Mongo state tracking + persistence subject `mapexos.events.asset_status_save` (events service → ClickHouse `asset_status_history`). NO publish to `mapexos.route.execute`.
  - `Enabled=true` with at least one route group non-empty → MONITOR + ROUTE: same as monitor-only PLUS publish to `mapexos.route.execute` with the configured route groups for the matching transition direction.
- Offline→online transition MUST be gated by atomic `RemoveAlerted` (Redis SREM) — exactly-once across pods/goroutines per reconnection
- First-ever heartbeat (unknown→online) sets `healthStatus=online` on the asset and marks `known`, but does NOT emit any NATS event (not a transition)
- Every write to `healthStatus` MUST also write `healthStatusChangedAt` in the same Mongo `$set` — all three transition sites (first heartbeat, offline→online, online→offline) go through `AssetRepo.UpdateHealthStatusWithChangedAt(..., time.Now().UTC())` so the flip timestamp the API surfaces cannot diverge from the status it's describing
- Scan cutoff = `now - asset.HealthMonitor.ThresholdMinutes`; offline fires only after `missCount >= RequiredMisses`
- `AlertEvent` routing is conditional on `len(RouteGroupIds) > 0` — persistence publish is unconditional
- `eventSource="healthStatus"` on the route payload tells Router to read `HealthMonitor.{Offline,Online}RouteGroupIds` instead of the default `Asset.RouteGroupIds`
- Scan consumer's `DuplicateWindow=10s` MUST stay below the scan interval (60s) to avoid the stream-level 15m dedupe silently swallowing re-schedules
- Only `trigger` and `workflow` router kinds may appear in health route groups (invariant enforced upstream by the Assets module at Create/Update)
- `handleHeartbeat` is intentionally **origin-agnostic**: payload `{orgId, assetUUID, ts}` produces identical Redis side-effects regardless of source (js-executor implicit publish vs HTTP gateway). Heartbeat origin selection is per-asset via `HealthMonitorConfig.heartbeatMode` (`'implicit'` default; `'explicit'` delegates the publish to the device for HTTP, or to the MQTT presence path for MQTT).
- **MQTT presence**: MQTT-protocol assets in `heartbeatMode='explicit'` signal liveness via the mapex-mqtt-broker plugin — the broker publishes `${env}.mapexos.mqtt.presence.advisory` with `Event=connect` on CONNECT and `Event=disconnect` on DISCONNECT. No firmware change required beyond keeping the MQTT connection open.
- **Anti-race on disconnect**: every CONNECT writes `lastConnectAt` (Redis HASH `hm:lc:{orgId}` field `assetUUID`); the disconnect handler skips messages where `disconnect.Timestamp <= redis.lastConnectAt`. Protects multi-broker-replica reconnect scenarios where a stale DISCONNECT arrives after a fresh CONNECT on a different broker replica.
- **3-layer MQTT-client filter** runs before any I/O on `HandlePresenceDisconnect`: `client_type=='mqtt'`, `kind=='Client'`, `user!=''`. Internal NATS service connections (assets MS, http_gateway, etc.) are dropped in microseconds with `HealthPresenceFiltered{reason="non_mqtt_client"}`.
- **Heartbeat lookup error semantics**: `loadAssetForHeartbeat` distinguishes transient errors (Nack → NATS redelivers) from legitimate drops (Ack → asset not found / monitoring disabled). Caller in `HandleHeartbeat` reflects this 3-state contract.

## Known Cross-Context Interactions
- Router service: source of heartbeat events; downstream consumer of `mapexos.route.execute` for offline/online routing
- Events service: consumer of `mapexos.events.asset_status_save` — persists the transition history in ClickHouse
- Assets module (same service): provides `HealthRepository` read access used for API enrichment (`healthStatus`, `healthStatusChangedAt`, `lastSeenAt`) and writes the atomic `(healthStatus, healthStatusChangedAt)` pair on every transition via `UpdateHealthStatusWithChangedAt`; provides `HealthMonitor` config and `RouteGroupIds` on the asset
- Assets module (same service): calls `HealthLifecyclePort.ClearAssetState` on Asset Update (when `HealthMonitor.Enabled` transitions true→false) and on Asset Delete to purge all Redis health state for the asset
- Redis (shared infra): hot-path state (last-seen, miss, alerted, known, active-orgs sets)

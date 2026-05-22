package ports

import (
	"context"
	"time"
)

// PresencePort exposes liveness write operations driven by NATS broker
// advisories: $SYS.ACCOUNT.*.CONNECT (online) and $SYS.ACCOUNT.*.DISCONNECT
// (offline). Idempotent and best-effort — callers MUST NOT block on
// these operations.
//
// Implemented by HealthMonitorService (healthmonitor module). Consumed
// by the healthmonitor presence consumer family — the CONNECT consumer
// invokes MarkOnlineFromConnect and the DISCONNECT consumer invokes
// MarkOfflineFromDisconnect. Other modules MUST NOT depend on this
// port directly; the broker advisory chain is the only blessed
// presence trigger.
type PresencePort interface {
	// MarkOnlineFromConnect marks an asset online following a NATS
	// $SYS.ACCOUNT.*.CONNECT advisory. Implementations MUST be
	// idempotent — every reconnect bursts a fresh CONNECT and presence
	// is best-effort state.
	MarkOnlineFromConnect(ctx context.Context, orgId string, assetUUID string) error

	// MarkOfflineFromDisconnect marks an asset offline following a NATS
	// $SYS.ACCOUNT.*.DISCONNECT advisory. Enforces the anti-race
	// invariant (disconnectAt > redis.lastConnectAt) so a stale
	// disconnect after a fresh reconnect on a different broker replica is
	// dropped.
	MarkOfflineFromDisconnect(ctx context.Context, orgId string, assetUUID string, reason string, disconnectAt time.Time) error
}

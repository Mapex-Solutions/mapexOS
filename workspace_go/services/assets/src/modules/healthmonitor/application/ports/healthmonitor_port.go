package ports

import (
	"context"
	"time"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"

	"assets/src/modules/healthmonitor/domain/entities"
)

// HealthMonitorServicePort defines the application service interface.
type HealthMonitorServicePort interface {
	HandleHeartbeat(msg *natsModel.Message)
	HandlePresenceConnect(msg *natsModel.Message)
	HandlePresenceDisconnect(msg *natsModel.Message)
	RunScan(ctx context.Context)
}

// HealthLifecyclePort exposes asset lifecycle hooks for cross-module callers
// (assets CRUD handlers). Idempotent — safe to call for assets that were
// never monitored and safe to call multiple times for the same asset.
//
// The current implementation clears the four Redis state stores managed by
// the healthmonitor module (last-seen ZSET, known-online SET, alerted SET,
// miss-counter HASH). It does NOT touch Mongo `healthStatus` — the assets
// module is responsible for resetting that field after this port returns
// (so the Mongo write is observable in the same caller's transaction
// boundary).
type HealthLifecyclePort interface {
	ClearAssetState(ctx context.Context, orgId string, assetUUID string) error
}

// HealthAdminPort exposes administrative state transitions that bypass the
// per-asset scan + threshold cycle. The only operation today is
// ForceOfflineByAssetUUID, used by e2e journeys to assert offline-action
// route group wiring without waiting the configured scan interval. Implementations
// must be idempotent — a second call once the asset is already offline is a no-op.
type HealthAdminPort interface {
	ForceOfflineByAssetUUID(ctx context.Context, assetUUID string, reason string) error
}

// HealthRepository defines the Redis operations for health state tracking.
type HealthRepository interface {
	// Hot path — called per heartbeat
	UpdateLastSeen(ctx context.Context, orgId string, assetUUID string, ts time.Time) error
	ResetMissCounter(ctx context.Context, orgId string, assetUUID string) error
	IsAlerted(ctx context.Context, orgId string, assetUUID string) (bool, error)
	// RemoveAlerted atomically removes the asset from the alerted (offline) set.
	// Returns (true, nil) only for the ONE caller that actually removed the entry —
	// concurrent heartbeats get (false, nil). Use this bool to gate the
	// offline→online transition so it fires exactly once per reconnection.
	RemoveAlerted(ctx context.Context, orgId string, assetUUID string) (bool, error)
	RegisterOrg(ctx context.Context, orgId string) error
	IsKnownOnline(ctx context.Context, orgId string, assetUUID string) (bool, error)
	MarkKnownOnline(ctx context.Context, orgId string, assetUUID string) error

	// Scanner — called per scan cycle
	FindStale(ctx context.Context, orgId string, cutoff time.Time, offset int64, limit int64) ([]string, error)
	IncrementMiss(ctx context.Context, orgId string, assetUUID string) (int64, error)
	MarkAlerted(ctx context.Context, orgId string, assetUUID string) error
	GetActiveOrgs(ctx context.Context) ([]string, error)

	// API enrichment — called per GetById / List
	GetLastSeen(ctx context.Context, orgId string, assetUUID string) (*time.Time, error)
	GetLastSeenBatch(ctx context.Context, orgId string, assetUUIDs []string) (map[string]*time.Time, error)
	IsAlertedBatch(ctx context.Context, orgId string, assetUUIDs []string) (map[string]bool, error)

	// Cleanup
	RemoveAsset(ctx context.Context, orgId string, assetUUID string) error

	// Presence — called from the CONNECT advisory consumer and the disconnect anti-race check.
	// SetLastConnectAt records the timestamp of the most recent successful
	// MQTT CONNECT for the asset. GetLastConnectAt returns nil when the
	// asset has never connected (presence consumer treats nil as drop reason).
	SetLastConnectAt(ctx context.Context, orgId string, assetUUID string, ts time.Time) error
	GetLastConnectAt(ctx context.Context, orgId string, assetUUID string) (*time.Time, error)
}

// AlertPublisherPort defines the interface for publishing health alert events to Router.
type AlertPublisherPort interface {
	PublishOffline(ctx context.Context, event entities.AlertEvent) error
	PublishOnline(ctx context.Context, event entities.AlertEvent) error
}

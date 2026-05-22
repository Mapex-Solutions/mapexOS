package ports

import (
	"time"

	"assets/src/modules/assets/domain/entities"
	"assets/src/modules/assets/domain/repositories"
	healthPorts "assets/src/modules/healthmonitor/application/ports"
)

// AssetCertificateInput is the cross-module payload used by sibling
// modules (e.g., mqttcerts) when they emit a freshly-issued device
// cert and need the assets module to mirror the metadata onto the
// asset entity and run the standard L2 / FANOUT side effects. PEM
// bytes are intentionally absent — they are returned to the operator
// once at issue time and never persisted server-side.
type AssetCertificateInput struct {
	Serial      string
	Fingerprint string
	SubjectCN   string
	IssuedAt    time.Time
	ExpiresAt   time.Time
}

// Asset is the public type alias for the domain entity.
// Cross-module consumers import this instead of domain/entities directly.
type Asset = entities.Asset

// HealthMonitorConfig is the public type alias for the domain entity.
type HealthMonitorConfig = entities.HealthMonitorConfig

// AssetWithTemplate is the public type alias for the entity-with-template
// projection used by paginated reads (and tests that need to build the
// list-path fixture).
type AssetWithTemplate = entities.AssetWithTemplate

// AssetRepository is the public type alias for the domain repository interface.
// Cross-module consumers inject this instead of domain/repositories directly.
type AssetRepository = repositories.AssetRepository

// HealthLifecyclePort is the public type alias for the healthmonitor lifecycle port.
// Cross-module consumers (assets) inject this to clear Redis health state when
// HealthMonitor.Enabled toggles from true to false, or when an asset is deleted.
type HealthLifecyclePort = healthPorts.HealthLifecyclePort

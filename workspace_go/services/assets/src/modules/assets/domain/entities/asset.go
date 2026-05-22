package entities

import (
	"time"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

type NoneConfig struct{}

// MqttConfig is the persistent MQTT-protocol identity for an asset.
// ClientId + Username are the public identity the device presents on
// CONNECT; AuthType is the platform-declared credential mode the
// broker enforces (password XOR cert — mutual exclusion).
//
// AuthType=password: the asset MUST have a PasswordHash (the bcrypt
// hash, returned to the broker via the read-model). Plaintext
// password is shown to the operator exactly once at create or
// change-password time and never retrievable afterwards.
//
// AuthType=cert: PasswordHash is empty by construction; the broker
// authenticates by cert serial against the active AssetCertificate
// on the parent Asset.
type MqttConfig struct {
	ClientId     string         `bson:"clientId"`
	Username     string         `bson:"username"`
	AuthType     string         `bson:"authType"`
	PasswordHash string         `bson:"passwordHash,omitempty"`
	CertTTL      *CertTTLConfig `bson:"certTTL,omitempty"`
}

// CertTTLConfig is the operator-declared validity window applied
// when the mqttcerts module signs this asset's device cert. Only
// meaningful when AuthType=cert. Bounds + unit→day mapping live in
// the mqttcerts application constants so the rule applies the same
// whether the value flows in through the API or through a script.
type CertTTLConfig struct {
	Value int    `bson:"value"`
	Unit  string `bson:"unit"`
}

type ProtocolType struct {
	Type string      `bson:"type"`
	Http *NoneConfig `bson:"http,omitempty"`
	Mqtt *MqttConfig `bson:"mqtt,omitempty"`
}

// HealthMonitorConfig configures sensor inactivity monitoring for this asset.
//
// HeartbeatMode chooses how the platform learns the device is alive:
//   - "" or "implicit" (default): js-executor emits a heartbeat for every
//     data event the device sends.
//   - "explicit": js-executor SKIPS implicit publishes; the path is chosen
//     by the asset's protocol — MQTT-protocol assets use NATS broker
//     presence ($SYS.ACCOUNT.*.CONNECT + $SYS.ACCOUNT.*.DISCONNECT
//     advisories; no device-side topic to publish); HTTP-protocol assets
//     POST to /api/v1/heartbeat?ds={dataSourceId} with body { assetUUID }.
type HealthMonitorConfig struct {
	Enabled              bool     `bson:"enabled"`
	ThresholdMinutes     int      `bson:"thresholdMinutes"`
	RequiredMisses       int      `bson:"requiredMisses"`
	HeartbeatMode        string   `bson:"heartbeatMode,omitempty"`
	OfflineRouteGroupIds []string `bson:"offlineRouteGroupIds"`
	OnlineRouteGroupIds  []string `bson:"onlineRouteGroupIds"`
}

// IsActive reports whether health monitoring is enabled for the asset.
// Pointer receiver makes the call nil-safe — callers can write
// `asset.HealthMonitor.IsActive()` even when HealthMonitor is nil.
func (h *HealthMonitorConfig) IsActive() bool {
	return h != nil && h.Enabled
}

// ResolvedMode returns "explicit" when the operator explicitly chose
// explicit mode; otherwise "implicit" (covers nil receiver, empty string,
// and any unrecognized value — back-compat default).
func (h *HealthMonitorConfig) ResolvedMode() string {
	if h == nil {
		return "implicit"
	}
	if h.HeartbeatMode == "explicit" {
		return "explicit"
	}
	return "implicit"
}

type Asset struct {
	ID           model.ObjectId `bson:"_id,omitempty"`
	Name         string         `bson:"name"`
	Enabled      bool           `bson:"enabled"`
	DebugEnabled bool           `bson:"debugEnabled"`
	Description  *string        `bson:"description,omitempty"`

	AssetUUID       string         `bson:"assetUUID"`
	AssetTemplateID model.ObjectId `bson:"assetTemplateId"`

	// Multi-tenant fields
	OrgID      model.ObjectId  `bson:"orgId"`
	PathKey    string          `bson:"pathKey"`              // Hierarchical path for range queries
	CustomerID *model.ObjectId `bson:"customerId,omitempty"` // Tenant anchor (denormalized)

	RouteGroupIds []string `bson:"routeGroupIds"`

	HealthMonitor         *HealthMonitorConfig `bson:"healthMonitor,omitempty"`
	HealthStatus          string               `bson:"healthStatus"`
	HealthStatusChangedAt *time.Time           `bson:"healthStatusChangedAt,omitempty"`

	Protocol  ProtocolType `bson:"protocol"`
	Latitude  *float64     `bson:"latitude,omitempty"`
	Longitude *float64     `bson:"longitude,omitempty"`

	// CurrentCert is the asset's currently-active mTLS device cert
	// metadata (one cert per asset; nil = no cert active). PEM bytes
	// are NEVER persisted — issued certs are returned to the operator
	// once at issue time and discarded server-side. See the mqttcerts
	// bounded context for the lifecycle.
	CurrentCert *AssetCertificate `bson:"currentCert,omitempty"`

	Created time.Time `bson:"created"`
	Updated time.Time `bson:"updated"`
}

// AssetCertificate is the embedded subdoc carrying the active cert's
// metadata. NEVER carries the PEM — PEM is returned once at issue
// time and discarded. Domain entity: bson-only tags, no json tags,
// no cross-service contract imports.
type AssetCertificate struct {
	Serial      string    `bson:"serial"`
	Fingerprint string    `bson:"fingerprint"`
	SubjectCN   string    `bson:"subjectCN"`
	IssuedAt    time.Time `bson:"issuedAt"`
	ExpiresAt   time.Time `bson:"expiresAt"`
}

// Keep these as requested
func (u *Asset) GetCreated() time.Time { return u.Created }
func (u *Asset) GetUpdated() time.Time { return u.Updated }

// PATCH/UPDATE payload (every field optional)
type AssetUpdate struct {
	ID           *model.ObjectId `bson:"_id,omitempty"`
	Name         *string         `bson:"name,omitempty"`
	Enabled      *bool           `bson:"enabled,omitempty"`
	DebugEnabled *bool           `bson:"debugEnabled,omitempty"`
	Description  *string         `bson:"description,omitempty"`

	AssetUUID       *string         `bson:"assetUUID,omitempty"`
	AssetTemplateID *model.ObjectId `bson:"assetTemplateId,omitempty"`

	// Multi-tenant fields (usually not updated, but available)
	OrgID      *model.ObjectId `bson:"orgId,omitempty"`
	PathKey    *string         `bson:"pathKey,omitempty"`
	CustomerID *model.ObjectId `bson:"customerId,omitempty"`

	RouteGroupIds *[]string `bson:"routeGroupIds,omitempty"`

	HealthMonitor         *HealthMonitorConfig `bson:"healthMonitor,omitempty"`
	HealthStatus          *string              `bson:"healthStatus,omitempty"`
	HealthStatusChangedAt *time.Time           `bson:"healthStatusChangedAt,omitempty"`

	Protocol  *ProtocolType `bson:"protocol,omitempty"`
	Latitude  *float64      `bson:"latitude,omitempty"`
	Longitude *float64      `bson:"longitude,omitempty"`

	Created *time.Time `bson:"created"`
	Updated *time.Time `bson:"updated"`
}

func (u *AssetUpdate) GetCreated() *time.Time { return u.Created }
func (u *AssetUpdate) GetUpdated() *time.Time { return u.Updated }

// AssetWithTemplate represents an Asset with template classification data joined via $lookup.
// This entity is used specifically for aggregation queries that need template information.
// It embeds the base Asset entity and adds template classification fields.
type AssetWithTemplate struct {
	ID           model.ObjectId `bson:"_id,omitempty"`
	Name         string         `bson:"name"`
	Enabled      bool           `bson:"enabled"`
	DebugEnabled bool           `bson:"debugEnabled"`
	Description  *string        `bson:"description,omitempty"`

	AssetUUID       string         `bson:"assetUUID"`
	AssetTemplateID model.ObjectId `bson:"assetTemplateId"`

	// Template classification data (populated from $lookup aggregation)
	CategoryId       *model.ObjectId `bson:"categoryId,omitempty"`
	CategoryName     *string         `bson:"categoryName,omitempty"`
	ManufacturerId   *model.ObjectId `bson:"manufacturerId,omitempty"`
	ManufacturerName *string         `bson:"manufacturerName,omitempty"`
	ModelId          *model.ObjectId `bson:"modelId,omitempty"`
	ModelName        *string         `bson:"modelName,omitempty"`
	Version          *string         `bson:"version,omitempty"`

	// Multi-tenant fields
	OrgID      model.ObjectId  `bson:"orgId"`
	PathKey    string          `bson:"pathKey"`
	CustomerID *model.ObjectId `bson:"customerId,omitempty"`

	RouteGroupIds []string `bson:"routeGroupIds"`

	// Health monitoring — HealthMonitor config + HealthStatusChangedAt flow through
	// the $lookup aggregation so the List path can surface them on AssetResponse
	// without an extra query. HealthMonitor MUST be projected so enrichHealthStatusBatch
	// can run per-asset (guarded by HealthMonitor.Enabled).
	HealthMonitor         *HealthMonitorConfig `bson:"healthMonitor,omitempty"`
	HealthStatus          string               `bson:"healthStatus,omitempty"`
	HealthStatusChangedAt *time.Time           `bson:"healthStatusChangedAt,omitempty"`

	Protocol  ProtocolType `bson:"protocol"`
	Latitude  *float64     `bson:"latitude,omitempty"`
	Longitude *float64     `bson:"longitude,omitempty"`

	// CurrentCert is projected from the same Mongo field used by the
	// canonical Asset entity. The list aggregation MUST $project it
	// (see templateProjectStage) — otherwise the response DTO comes
	// back without currentCert and the UI's "no certificate" warning
	// fires on every cert-mode row even after a successful issue.
	CurrentCert *AssetCertificate `bson:"currentCert,omitempty"`

	Created time.Time `bson:"created"`
	Updated time.Time `bson:"updated"`
}

func (u *AssetWithTemplate) GetCreated() time.Time { return u.Created }
func (u *AssetWithTemplate) GetUpdated() time.Time { return u.Updated }

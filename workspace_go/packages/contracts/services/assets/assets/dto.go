package assets

import (
	"fmt"
	"time"

	"github.com/Mapex-Solutions/MapexOS/contracts/common"
	"github.com/Mapex-Solutions/MapexOS/contracts/common/query"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

/**
 * SHARED READ MODEL (CQRS Pattern)
 *
 * AssetReadModel represents the shared read model for cross-service queries.
 * This is a CQRS Read Model - a denormalized projection optimized for reads.
 *
 * OWNERSHIP: Asset Service (write only)
 * CONSUMERS: MQTT Gateway, Router, JS-Executor, Triggers, Asset Service (read only)
 *
 * Key Format: read:asset:{assetUUID}
 * TTL: 7 days
 * Redis DB: 5 (SharedCache)
 */
type AssetReadModel struct {
	// Identity
	ID   string `json:"id"`
	UUID string `json:"uuid"`

	// Multi-tenant
	OrgId   string `json:"orgId"`
	PathKey string `json:"pathKey"`

	// Status
	Enabled      bool `json:"enabled"`
	DebugEnabled bool `json:"debugEnabled"`

	// Asset data
	Name               string   `json:"name"`
	Description        string   `json:"description,omitempty"`
	AssetTemplateID    string   `json:"assetTemplateId,omitempty"`
	AssetTemplateOrgID string   `json:"assetTemplateOrgId,omitempty"`
	RouteGroupIds      []string `json:"routeGroupIds,omitempty"`

	// Health monitoring
	HealthMonitor *HealthMonitorConfig `json:"healthMonitor,omitempty"`
	HealthStatus  string               `json:"healthStatus,omitempty"`

	// Protocol (for MQTT/HTTP gateway)
	Protocol *ProtocolType `json:"protocol,omitempty"`

	// CurrentCert is the asset's currently-active MQTT device cert
	// metadata. Empty / nil when the asset has no active cert.
	// Consumed by the broker plugin to validate cert-mode CONNECTs
	// (serial-equality compare with the client cert presented at TLS
	// handshake).
	CurrentCert *AssetCertificate `json:"currentCert,omitempty"`

	// Location
	Latitude  *float64 `json:"latitude,omitempty"`
	Longitude *float64 `json:"longitude,omitempty"`

	// Timestamps
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}

// AssetCertificate mirrors the embedded subdoc on the Asset entity.
// PEM bytes are NEVER persisted or transmitted — only the metadata
// fields needed for the broker plugin to validate a cert-mode CONNECT
// (serial match) and for ops dashboards (expiry).
type AssetCertificate struct {
	Serial      string    `json:"serial"`
	Fingerprint string    `json:"fingerprint"`
	SubjectCN   string    `json:"subjectCN"`
	IssuedAt    time.Time `json:"issuedAt"`
	ExpiresAt   time.Time `json:"expiresAt"`
}

type NoneConfig struct{}

// MqttConfig carries the device-MQTT identity for an asset. ClientId +
// Username form the public identity the device presents on CONNECT;
// AuthType is the platform-declared credential mode the broker
// enforces (mutual exclusion — password XOR cert).
//
// Mode semantics:
//   - AuthType=password: device CONNECTs with username + password on
//     the plaintext listener; the broker bcrypt-compares against
//     PasswordHash. The asset MUST have a PasswordHash (set at create
//     or via subsequent rotation).
//   - AuthType=cert: device CONNECTs with an mTLS client cert on port
//     8883; the broker validates the cert serial against CurrentCert.
//     Username on the CONNECT may be empty — the broker uses the
//     cert's Subject CN ("{orgId}:{assetUUID}") as the lookup key.
//
// Field asymmetry between request and read-model contexts:
//   - Password is the plaintext credential, present ONLY on create/update
//     request bodies (and only when AuthType=password) — the platform
//     bcrypt-hashes it before persistence and the response path always
//     omits it (one-shot delivery, never retrievable).
//   - PasswordHash is the bcrypt hash, present ONLY on the read-model
//     path (AssetReadModel emitted to L2 MinIO / served by the internal
//     read-model endpoint). The broker plugin consumes it for local
//     bcrypt-compare on password-mode CONNECTs. Request bodies never
//     carry it — validation rejects the field there.
type MqttConfig struct {
	ClientId     string `json:"clientId" validate:"required,min=3"`
	Username     string `json:"username" validate:"required,min=3"`
	AuthType     string `json:"authType" validate:"required,oneof=password cert"`
	Password     string `json:"password,omitempty" validate:"omitempty,min=8"`
	PasswordHash string `json:"passwordHash,omitempty"`

	// CertTTL is the operator-declared validity window for the asset's
	// MQTT device cert. Only meaningful when AuthType=cert. The
	// mqttcerts module reads this when signing a new cert (IssueCert)
	// so each asset can declare its own rotation cadence — a short TTL
	// suits dev / test (e.g. 1 day) while industrial deployments stick
	// to the year-scale default. Absent = platform default (1 year).
	CertTTL *CertTTLConfig `json:"certTTL,omitempty" validate:"omitempty"`
}

// CertTTLConfig declares the operator's preferred validity window
// for the device cert. The mqttcerts signer converts (Value, Unit)
// into a day count and uses it as the cert's NotAfter offset.
//
// Bounds (enforced by the signer, not by struct tags so the same
// constants can also gate the UI's input validation):
//   - Min total: 1 day
//   - Max total: 3650 days (10 years)
//
// Unit values map to days as: day=1, week=7, month=30, year=365.
// Month and year approximations are deliberate — calendar-accurate
// math would require the issuance date plus locale-aware rules and
// the cert's resolution is already day-scale.
type CertTTLConfig struct {
	Value int    `json:"value" validate:"required,min=1"`
	Unit  string `json:"unit"  validate:"required,oneof=day week month year"`
}

type ProtocolType struct {
	Type string      `json:"type" validate:"required,oneof=http mqtt lorawan"`
	Http *NoneConfig `json:"http,omitempty"`
	Mqtt *MqttConfig `json:"mqtt,omitempty"`
}

// HealthMonitorConfig represents health monitoring configuration for an asset.
//
// HeartbeatMode chooses how the platform learns the device is alive:
//   - "implicit" (default, omitted = same): js-executor emits a heartbeat for
//     every data event the device sends.
//   - "explicit": js-executor SKIPS implicit publishes; the path is chosen by
//     the asset's protocol — MQTT-protocol assets use NATS broker presence
//     ($SYS.ACCOUNT.*.CONNECT + $SYS.ACCOUNT.*.DISCONNECT advisories;
//     no device-side topic to publish); HTTP-protocol assets POST to
//     /api/v1/heartbeat?ds={dataSourceId} with body { assetUUID }.
type HealthMonitorConfig struct {
	Enabled              *bool    `json:"enabled,omitempty"`
	ThresholdMinutes     *int     `json:"thresholdMinutes,omitempty" validate:"omitempty,min=10"`
	RequiredMisses       *int     `json:"requiredMisses,omitempty" validate:"omitempty,min=1"`
	HeartbeatMode        *string  `json:"heartbeatMode,omitempty" validate:"omitempty,oneof=implicit explicit"`
	OfflineRouteGroupIds []string `json:"offlineRouteGroupIds,omitempty" validate:"omitempty,max=3,dive,mongoid"`
	OnlineRouteGroupIds  []string `json:"onlineRouteGroupIds,omitempty" validate:"omitempty,max=3,dive,mongoid"`
}

type AssetId struct {
	AssetId string `params:"assetId" validate:"required,mongoid"`
}

type AssetUUID struct {
	AssetUUID string `params:"assetUUID" validate:"required,min=5"`
}

type AssetCreate struct {
	Name         string  `json:"name" validate:"required,min=1"`
	Enabled      bool    `json:"enabled" validate:"required"`
	DebugEnabled bool    `json:"debugEnabled"`
	Description  *string `json:"description,omitempty" validate:"omitempty,max=500"`

	AssetUUID       string `json:"assetUUID" validate:"required,min=5"`
	AssetTemplateID string `json:"assetTemplateId" validate:"required,mongoid"`

	// Multi-tenant fields (populated automatically by coverage middleware)
	OrgID         *model.ObjectId `json:"orgId,omitempty" validate:"omitempty"`
	PathKey       *string         `json:"pathKey,omitempty" validate:"omitempty"`
	RouteGroupIds []string        `json:"routeGroupIds" validate:"required,min=1,max=3,dive,mongoid"`

	HealthMonitor *HealthMonitorConfig `json:"healthMonitor,omitempty"`

	Protocol  ProtocolType `json:"protocol"`
	Latitude  *float64     `json:"latitude,omitempty"`
	Longitude *float64     `json:"longitude,omitempty"`

	Created *common.NullTime `json:"created,omitempty"`
	Updated *common.NullTime `json:"updated,omitempty"`
}

type AssetUpdate struct {
	Name         *string `json:"name,omitempty" validate:"omitempty,min=1"`
	Enabled      *bool   `json:"enabled,omitempty" validate:"omitempty"`
	DebugEnabled *bool   `json:"debugEnabled,omitempty" validate:"omitempty"`
	Description  *string `json:"description,omitempty" validate:"omitempty,max=500"`

	AssetUUID       *string `json:"assetUUID,omitempty" validate:"omitempty"`
	AssetTemplateID *string `json:"assetTemplateId,omitempty" validate:"omitempty,mongoid"`

	OrgId         *model.ObjectId `json:"orgId,omitempty" validate:"omitempty"`
	RouteGroupIds *[]string       `json:"routeGroupIds,omitempty" validate:"omitempty,min=1,max=3,dive,mongoid"`

	HealthMonitor *HealthMonitorConfig `json:"healthMonitor,omitempty"`

	Protocol  *ProtocolType `json:"protocol,omitempty"`
	Latitude  *float64      `json:"latitude,omitempty"`
	Longitude *float64      `json:"longitude,omitempty"`

	Created *common.NullTime `json:"created,omitempty"`
	Updated *common.NullTime `json:"updated,omitempty"`
}

// AssetQuery represents query parameters for listing assets.
// Embeds BaseQueryDTO for standard pagination, sorting, and hierarchy support.
//
// Standard fields (from BaseQueryDTO):
//   - Projection: comma-separated fields to return
//   - Page: page number (default: 1)
//   - PerPage: items per page (default: 20)
//   - Sort: sort order (default: "created:desc")
//   - IncludeChildren: include child orgs hierarchically (default: false)
//
// Module-specific filters:
//   - Name: filter by asset name (partial match)
//   - Enabled: filter by enabled status (true/false)
//   - AssetUUID: filter by device UUID
//   - AssetTemplateID: filter by template
//   - Category: filter by category
//   - AssetType: filter by type
//
// Organization filtering is handled automatically via RequestContext:
//   - No manual orgId/pathKey/customerId needed
//   - Context-aware filtering via X-Org-Context header
//   - Hierarchical queries via includeChildren parameter
type AssetQuery struct {
	query.BaseQueryDTO

	// Module-specific filters
	Name            *string `query:"name" validate:"omitempty,max=100"`
	Enabled         *bool   `query:"enabled" validate:"omitempty"`
	AssetUUID       *string `query:"assetUUID" validate:"omitempty"`
	AssetTemplateID *string `query:"assetTemplateId" validate:"omitempty,mongoid"`
	Category        *string `query:"category" validate:"omitempty,mongoid"`
	AssetType       *string `query:"assetType" validate:"omitempty,mongoid"`

	// Classification filters (for template-based filtering)
	CategoryId     *string `query:"categoryId" validate:"omitempty,mongoid"`
	ManufacturerId *string `query:"manufacturerId" validate:"omitempty,mongoid"`
	ModelId        *string `query:"modelId" validate:"omitempty,mongoid"`

	// Health monitoring filter
	HealthStatus *string `query:"healthStatus" validate:"omitempty,oneof=online offline unknown"`
}

type AssetResponse struct {
	ID           *common.ObjectID `json:"id,omitempty"`
	Name         *string          `json:"name,omitempty"`
	Enabled      *bool            `json:"enabled,omitempty"`
	DebugEnabled *bool            `json:"debugEnabled,omitempty"`
	Description  *string          `json:"description,omitempty"`

	AssetUUID       *string `json:"assetUUID,omitempty"`
	AssetTemplateID *string `json:"assetTemplateId,omitempty"`

	// Template classification data (populated from AssetTemplate lookup)
	AssetTemplateName *string `json:"assetTemplateName,omitempty"`
	AssetIdPath       *string `json:"assetIdPath,omitempty"` // Path to extract asset UUID from payload (from template)
	CategoryId        *string `json:"categoryId,omitempty"`
	CategoryName      *string `json:"categoryName,omitempty"`
	ManufacturerId    *string `json:"manufacturerId,omitempty"`
	ManufacturerName  *string `json:"manufacturerName,omitempty"`
	ModelId           *string `json:"modelId,omitempty"`
	ModelName         *string `json:"modelName,omitempty"`
	Version           *string `json:"version,omitempty"`

	OrgId         *common.ObjectID `json:"orgId,omitempty"`
	PathKey       *string          `json:"pathKey,omitempty"`
	CustomerID    *common.ObjectID `json:"customerId,omitempty"`
	RouteGroupIds   *[]string  `json:"routeGroupIds,omitempty"`
	RouteGroupNames *[]string  `json:"routeGroupNames,omitempty"` // Populated from Router service lookup

	HealthMonitor         *HealthMonitorConfig `json:"healthMonitor,omitempty"`
	HealthStatus          *string              `json:"healthStatus,omitempty"`
	HealthStatusChangedAt *time.Time           `json:"healthStatusChangedAt,omitempty"` // Mongo-persisted flip time; null if never transitioned
	LastSeenAt            *time.Time           `json:"lastSeenAt,omitempty"`            // Enriched from Redis (real-time)

	Protocol  *ProtocolType    `json:"protocol,omitempty"`
	Latitude  *float64         `json:"latitude,omitempty"`
	Longitude *float64         `json:"longitude,omitempty"`

	// CurrentCert is the asset's currently-active MQTT device cert
	// metadata. nil when the asset has no active cert (password-mode
	// asset, or cert-mode asset that hasn't issued its first cert yet).
	// The UI uses presence/absence here to derive the auth-mode radio
	// (cert when set; password when nil) and to gate "Generate cert" vs
	// "Rotate cert" affordances in the asset details drawer.
	CurrentCert *AssetCertificate `json:"currentCert,omitempty"`

	Created *common.NullTime `json:"created,omitempty"`
	Updated *common.NullTime `json:"updated,omitempty"`
}

func (d *AssetResponse) SetCreated(t *common.NullTime) { d.Created = t }
func (d *AssetResponse) SetUpdated(t *common.NullTime) { d.Updated = t }

// GenerateMqttPasswordResponse is returned by GET
// /api/v1/assets/_generate_mqtt_password. Stateless — does not touch
// any asset; just returns a strong random alphanumeric string the UI
// can drop into the password field. The operator remains free to type
// a custom password instead.
type GenerateMqttPasswordResponse struct {
	Password string `json:"password"`
}

/** TRANSFORMATIONS **/

func (a *AssetCreate) Transform() error {
	// Validate RouteGroupIds uniqueness (no duplicates)
	seen := make(map[string]bool)
	for _, id := range a.RouteGroupIds {
		if seen[id] {
			return fmt.Errorf("duplicate routeGroupId found: %s", id)
		}
		seen[id] = true
	}

	return nil
}

func (a *AssetUpdate) Transform() error {
	// Validate RouteGroupIds uniqueness (no duplicates) if provided
	if a.RouteGroupIds != nil {
		seen := make(map[string]bool)
		for _, id := range *a.RouteGroupIds {
			if seen[id] {
				return fmt.Errorf("duplicate routeGroupId found: %s", id)
			}
			seen[id] = true
		}
	}

	return nil
}

func (p *ProtocolType) Transform() error {
	if p.Type == "mqtt" && p.Mqtt == nil {
		return fmt.Errorf("field 'protocol.mqtt' must be provided when type is 'mqtt'")
	}

	return nil
}

// Note: Internal DTOs (AssetInternalId, AssetInternalUpdate, AssetUUIDParam,
// AssetScriptsResponse, AssetRouteGroupsResponse) were removed as consuming
// services now fetch asset data via TieredCache (L2 = MinIO) instead of
// internal API endpoints. Cache invalidation is handled via NATS FANOUT.

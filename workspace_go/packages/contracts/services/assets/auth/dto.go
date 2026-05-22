// Package auth holds the cross-service auth projection contract used by
// the broker plugin to decide MQTT CONNECTs. The projection is written
// by the assets service after every CRUD that affects auth fields and
// is read from MinIO bucket `mapex-asset-auth` (key {assetUUID}.json)
// or fetched via the assets service internal endpoint
// GET /internal/asset-auth/:assetUUID as the L3 fallback.
//
// The projection is intentionally slimmer than AssetReadModel (defined
// in services/assets/assets/dto.go): it carries only the fields the
// broker plugin needs to authenticate a CONNECT. Keeping the auth
// payload narrow lets the broker cache it densely on its L1 Pebble
// store and pull it cheaply from L2 on cold lookups.
//
// Reciprocity: mirrored by workspace_js/packages/schemas/src/assets/schemas/auth.
//
// Contracts stay leaf-level — no imports from services/.
package auth

// AuthProjection is the slim auth-only payload stored at
// mapex-asset-auth/{assetUUID}.json. The broker plugin consumes it on
// every CONNECT lookup; the assets service writes it as a side effect
// of every CRUD that touches auth fields (password rotation, cert
// issue/revoke, asset enable/disable).
//
// Type is intentionally open so future auth surfaces (http_api,
// lorawan_key, etc.) can reuse the same projection shape and the same
// bucket without breaking the broker. Today only "mqtt" is valid; the
// validator enforces that until additional surfaces are added.
type AuthProjection struct {
	AssetUUID         string `json:"assetUUID"         validate:"required,min=1"`
	OrgId             string `json:"orgId"             validate:"required"`
	Enabled           bool   `json:"enabled"`
	Type              string `json:"type"              validate:"required,oneof=mqtt"`
	AuthType          string `json:"authType"          validate:"required,oneof=password cert"`
	PasswordHash      string `json:"passwordHash,omitempty"`
	CurrentCertSerial string `json:"currentCertSerial,omitempty"`
}

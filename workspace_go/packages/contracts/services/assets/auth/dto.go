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
// Type discriminates the auth surface: "mqtt" sets the Mqtt block, "lorawan"
// sets the Lorawan block. Both are nested and symmetric — the broker reads
// Mqtt, the LNS reads Lorawan. The shape is open so future surfaces (http_api,
// etc.) add their own block without breaking existing readers.
type AuthProjection struct {
	AssetUUID string `json:"assetUUID" validate:"required,min=1"`
	OrgId     string `json:"orgId"     validate:"required"`
	Enabled   bool   `json:"enabled"`
	Type      string `json:"type"      validate:"required,oneof=mqtt lorawan"`

	// Mqtt is set when Type == "mqtt". It carries the broker's CONNECT-decision
	// fields: the auth mode plus the password hash (password mode) or the
	// current cert serial (cert mode).
	Mqtt *MqttAuth `json:"mqtt,omitempty"`

	// Lorawan is set when Type == "lorawan". It carries the device identity +
	// profile and the envelope-encrypted key material the LNS decrypts at
	// hydrate.
	Lorawan *LorawanAuth `json:"lorawan,omitempty"`
}

// MqttAuth is the slim MQTT auth block in the projection. The broker plugin
// decides every CONNECT from these fields alone: AuthType selects the mode,
// PasswordHash backs password mode, CurrentCertSerial backs cert mode.
type MqttAuth struct {
	AuthType          string `json:"authType,omitempty" validate:"omitempty,oneof=password cert"`
	PasswordHash      string `json:"passwordHash,omitempty"`
	CurrentCertSerial string `json:"currentCertSerial,omitempty"`
}

// LorawanAuth is the slim LoRaWAN auth block in the projection. Identity +
// profile are clear; the secret key material lives encrypted in Keys.
type LorawanAuth struct {
	DevEUI     string        `json:"devEui"`
	JoinEUI    string        `json:"joinEui,omitempty"`
	Region     string        `json:"region"`
	Class      string        `json:"class"`
	MacVersion string        `json:"macVersion"`
	PhyVersion string        `json:"phyVersion"`
	Activation string        `json:"activation"`
	Keys       EncryptedKeys `json:"keys"`
}

// EncryptedKeys holds the four envelope fields produced by the goKit envelope
// primitive. The plaintext sealed inside is a LorawanKeyMaterial JSON.
type EncryptedKeys struct {
	EncryptedDEK []byte `json:"encryptedDek"`
	DekNonce     []byte `json:"dekNonce"`
	EncryptedKey []byte `json:"encryptedKey"`
	KeyNonce     []byte `json:"keyNonce"`
}

// LorawanKeyMaterial is the plaintext shape sealed inside EncryptedKeys. The
// assets MS marshals it before encryption; the LNS unmarshals it after
// decryption. Both sides MUST agree on these field names.
type LorawanKeyMaterial struct {
	AppKey  string `json:"appKey,omitempty"`
	NwkKey  string `json:"nwkKey,omitempty"`
	DevAddr string `json:"devAddr,omitempty"`
	NwkSKey string `json:"nwkSKey,omitempty"`
	AppSKey string `json:"appSKey,omitempty"`
}

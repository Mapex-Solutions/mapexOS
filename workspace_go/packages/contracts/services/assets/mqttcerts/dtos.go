// Package mqttcerts carries the cross-service contracts for the
// Assets MS device-cert lifecycle endpoints. Wire shapes only.
package mqttcerts

import "time"

// IssueCertRequest is POST /api/v1/mqtt_certs.
// Force=true accepts replacing an asset's existing currentCert (prior
// cert moves to mqttRevokedCertificates with reason=replaced).
type IssueCertRequest struct {
	AssetUUID string `json:"assetUUID" validate:"required"`
	Force     bool   `json:"force"`
}

// IssueCertResponse is the one-shot return: cert + key + ca-chain.
// Server NEVER persists key bytes; the operator downloads once via
// the frontend zip helper.
type IssueCertResponse struct {
	Serial      string    `json:"serial"`
	Fingerprint string    `json:"fingerprint"`
	SubjectCN   string    `json:"subjectCN"`
	IssuedAt    time.Time `json:"issuedAt"`
	ExpiresAt   time.Time `json:"expiresAt"`
	CertPEM     []byte    `json:"certPEM"`
	KeyPEM      []byte    `json:"keyPEM"`
	CAChainPEM  []byte    `json:"caChainPEM"`
}

// RevokedCertResponse is one row of the revoked list endpoint.
type RevokedCertResponse struct {
	Serial      string    `json:"serial"`
	Fingerprint string    `json:"fingerprint"`
	AssetUUID   string    `json:"assetUUID"`
	OrgID       string    `json:"orgId"`
	SubjectCN   string    `json:"subjectCN"`
	IssuedAt    time.Time `json:"issuedAt"`
	RevokedAt   time.Time `json:"revokedAt"`
	Reason      string    `json:"reason"`
}

// ListRevokedQuery binds the query string for GET /api/v1/mqtt_certs.
type ListRevokedQuery struct {
	AssetUUID string `query:"assetUUID" validate:"required"`
}

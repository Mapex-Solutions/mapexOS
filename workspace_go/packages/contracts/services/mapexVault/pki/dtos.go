// Package pki carries the cross-service contracts for mapexVault's PKI
// bounded context. Wire shapes only — no business logic.
//
// Reciprocity per /go-arch §4: every shape here has a matching Zod
// schema under workspace_js/packages/schemas/src/services/mapexVault/pki/.
package pki

import "time"

// IntermediateCABundleResponse is the wire shape for
// GET /internal/pki/intermediate_ca_bundle. Returned ONLY to the Assets
// MS at boot. PrivateKeyPEM is plaintext bytes after envelope decryption;
// the caller is responsible for not logging this field.
type IntermediateCABundleResponse struct {
	CertPEM       []byte    `json:"certPEM"`
	PrivateKeyPEM []byte    `json:"privateKeyPEM"`
	SubjectCN     string    `json:"subjectCN"`
	NotBefore     time.Time `json:"notBefore"`
	NotAfter      time.Time `json:"notAfter"`
	Fingerprint   string    `json:"fingerprint"`
}

// SignServerRequest is the wire shape for POST /internal/pki/sign_server.
// Used by the broker provisioning script (and by tests). All fields are
// required.
type SignServerRequest struct {
	CN      string   `json:"cn" validate:"required"`
	SANs    []string `json:"sans" validate:"required,min=1,dive,required"`
	TTLDays int      `json:"ttlDays" validate:"required,min=1"`
}

// SignServerResponse carries the server cert produced by mapexVault.
// Caller persists certPEM to disk; serial is hex-encoded for human
// readability + log correlation.
type SignServerResponse struct {
	CertPEM   []byte    `json:"certPEM"`
	SerialHex string    `json:"serialHex"`
	NotAfter  time.Time `json:"notAfter"`
}

// CAChainResponse is the wire shape for GET /internal/pki/ca_chain.
// Returns root + intermediate concatenated as PEM (the "trust anchor
// bundle"). Public material — safe to expose anywhere.
type CAChainResponse struct {
	Chain []byte `json:"chain"`
}

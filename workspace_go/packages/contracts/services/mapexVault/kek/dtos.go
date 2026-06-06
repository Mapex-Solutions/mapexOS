// Package kek carries the cross-service contract for mapexVault's key
// encryption key (KEK) bounded context. Wire shapes only, no business logic.
//
// The KEK endpoint is internal Go-to-Go (Vault to assets MS / LNS), so there is
// no TS counterpart.
package kek

// KEKResponse is the wire shape for GET /internal/kek/:context. Returned ONLY to
// authorized internal services (assets MS, LNS) that envelope-encrypt device
// secrets locally. Kek is the 64-hex AES-256 key after envelope decryption; the
// caller must never log or persist it.
type KEKResponse struct {
	Context string `json:"context"`
	Kek     string `json:"kek"`
}

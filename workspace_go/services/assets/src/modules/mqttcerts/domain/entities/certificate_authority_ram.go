package entities

import "time"

// CertificateAuthorityRAM is the in-memory representation of the
// intermediate CA bundle fetched from mapexVault at boot. Held in
// atomic.Pointer for hot-swap-safe reads. NEVER persisted, NEVER
// serialized over the wire — no bson tags, no json tags.
type CertificateAuthorityRAM struct {
	CertPEM       []byte
	PrivateKeyPEM []byte
	SubjectCN     string
	NotBefore     time.Time
	NotAfter      time.Time
	Fingerprint   string
}

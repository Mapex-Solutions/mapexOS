package mapexvault_client

import (
	"time"

	"github.com/Mapex-Solutions/mapexGoKit/infrastructure/httpclient"
)

// MapexVaultClient is the concrete HTTP adapter for the
// MapexVaultClientPort. Wraps the platform's gokit httpclient so
// transport, X-API-Key injection, and timeout policy stay aligned
// with every other internal-API caller (router, mapexos, etc).
type MapexVaultClient struct {
	client *httpclient.HTTPClient
}

// intermediateCABundleWire decodes the on-wire JSON; field names match
// the Go contract IntermediateCABundleResponse json tags exactly.
type intermediateCABundleWire struct {
	CertPEM       []byte    `json:"certPEM"`
	PrivateKeyPEM []byte    `json:"privateKeyPEM"`
	SubjectCN     string    `json:"subjectCN"`
	NotBefore     time.Time `json:"notBefore"`
	NotAfter      time.Time `json:"notAfter"`
	Fingerprint   string    `json:"fingerprint"`
}

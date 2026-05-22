package ports

import (
	"context"

	"assets/src/modules/mqttcerts/domain/entities"
)

// MapexVaultClientPort wraps the HTTP call to
// `GET /internal/pki/intermediate_ca_bundle` on mapexVault.
type MapexVaultClientPort interface {
	FetchIntermediateCABundle(ctx context.Context) (*entities.CertificateAuthorityRAM, error)
}

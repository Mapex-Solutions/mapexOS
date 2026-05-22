package ports

import (
	"context"

	dtos "mapexVault/src/modules/pki/application/dtos"
)

// PkiServicePort is the entry point for the PKI bounded context. CA
// documents are seeded into Mongo by the mongodb-init container at
// deploy time; the service itself does not bootstrap.
type PkiServicePort interface {
	GetIntermediateCABundle(ctx context.Context) (*dtos.IntermediateCABundleResponse, error)
	GetCAChain(ctx context.Context) (*dtos.CAChainResponse, error)
	SignServerCert(ctx context.Context, req *dtos.SignServerRequest) (*dtos.SignServerResponse, error)
}

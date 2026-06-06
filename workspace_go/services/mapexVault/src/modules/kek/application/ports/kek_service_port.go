package ports

import (
	"context"

	dtos "mapexVault/src/modules/kek/application/dtos"
)

// KekServicePort is the entry point for the KEK bounded context. KEK documents
// are seeded into Mongo by the mongodb-init container at deploy time; the
// service itself does not bootstrap.
type KekServicePort interface {
	GetKEK(ctx context.Context, keyContext string) (*dtos.KEKResponse, error)
}

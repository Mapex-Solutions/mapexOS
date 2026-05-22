package repositories

import (
	"context"

	"mapexVault/src/modules/pki/domain/constants"
	"mapexVault/src/modules/pki/domain/entities"
)

// CARepository persists CertificateAuthority records to MongoDB.
// Application layer depends on this interface; infrastructure provides
// the concrete adapter.
type CARepository interface {
	FindByKind(ctx context.Context, kind constants.CAKind) (*entities.CertificateAuthority, error)
	Create(ctx context.Context, ca *entities.CertificateAuthority) error
	CountByKind(ctx context.Context, kind constants.CAKind) (int64, error)
}

package repositories

import (
	"context"

	"mapexVault/src/modules/kek/domain/entities"
)

// EncryptionKeyRepository reads platform KEK records from MongoDB. The
// application layer depends on this interface; infrastructure provides the
// concrete adapter. Records are seeded by the mongodb-init container.
type EncryptionKeyRepository interface {
	FindByContext(ctx context.Context, keyContext string) (*entities.EncryptionKey, error)
}

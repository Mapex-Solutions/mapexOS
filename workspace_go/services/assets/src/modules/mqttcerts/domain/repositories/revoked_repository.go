package repositories

import (
	"context"

	"assets/src/modules/mqttcerts/domain/entities"
)

// RevokedRepository persists RevokedCertificate rows.
type RevokedRepository interface {
	Create(ctx context.Context, r *entities.RevokedCertificate) error
	FindByAssetUUID(ctx context.Context, assetUUID string) ([]*entities.RevokedCertificate, error)
	FindBySerial(ctx context.Context, serial string) (*entities.RevokedCertificate, error)
	DeleteByAssetUUID(ctx context.Context, assetUUID string) (int64, error)
	CountByAssetUUID(ctx context.Context, assetUUID string) (int64, error)
}

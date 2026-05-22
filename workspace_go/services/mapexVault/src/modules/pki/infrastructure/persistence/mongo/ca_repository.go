package mongo

import (
	"context"

	"mapexVault/src/modules/pki/domain/constants"
	"mapexVault/src/modules/pki/domain/entities"
	"mapexVault/src/modules/pki/domain/repositories"

	manager "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/manager"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

type caRepository struct {
	model *model.Model[entities.CertificateAuthority]
}

// NewCARepository constructs the Mongo adapter for CertificateAuthority
// records and ensures indexes idempotently:
//   - kind UNIQUE (one record per kind in the hierarchy)
//   - isSystem + kind compound (defense in depth)
func NewCARepository(m *manager.MongoManager) repositories.CARepository {
	mdl := model.New[entities.CertificateAuthority](m.GetDatabase(), CollectionName, model.Config{
		Indexes: []model.IndexDefinition{
			{
				Name:   "kind_unique",
				Keys:   map[string]int{"kind": 1},
				Unique: true,
			},
			{
				Name: "isSystem_kind",
				Keys: map[string]int{"isSystem": 1, "kind": 1},
			},
		},
	})
	logger.Info("[REPO:CertificateAuthority] indexes ensured collection=" + CollectionName)
	return &caRepository{model: mdl}
}

func (r *caRepository) FindByKind(ctx context.Context, kind constants.CAKind) (*entities.CertificateAuthority, error) {
	filter := model.Map{"kind": string(kind)}
	return r.model.FindOne(ctx, &filter)
}

func (r *caRepository) Create(ctx context.Context, ca *entities.CertificateAuthority) error {
	_, err := r.model.CreateOne(ctx, ca)
	return err
}

func (r *caRepository) CountByKind(ctx context.Context, kind constants.CAKind) (int64, error) {
	filter := model.Map{"kind": string(kind)}
	return r.model.DIRECT().CountDocuments(ctx, filter)
}

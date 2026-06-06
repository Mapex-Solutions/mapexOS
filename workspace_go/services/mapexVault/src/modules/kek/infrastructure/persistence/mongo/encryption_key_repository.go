package mongo

import (
	"context"

	"mapexVault/src/modules/kek/domain/entities"
	"mapexVault/src/modules/kek/domain/repositories"

	manager "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/manager"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

type encryptionKeyRepository struct {
	model *model.Model[entities.EncryptionKey]
}

// NewEncryptionKeyRepository constructs the Mongo adapter for EncryptionKey
// records and ensures indexes idempotently:
//   - context UNIQUE (one KEK per context)
//   - isSystem + context compound (defense in depth)
func NewEncryptionKeyRepository(m *manager.MongoManager) repositories.EncryptionKeyRepository {
	mdl := model.New[entities.EncryptionKey](m.GetDatabase(), CollectionName, model.Config{
		Indexes: []model.IndexDefinition{
			{
				Name:   "context_unique",
				Keys:   map[string]int{"context": 1},
				Unique: true,
			},
			{
				Name: "isSystem_context",
				Keys: map[string]int{"isSystem": 1, "context": 1},
			},
		},
	})
	logger.Info("[REPO:EncryptionKey] indexes ensured collection=" + CollectionName)
	return &encryptionKeyRepository{model: mdl}
}

func (r *encryptionKeyRepository) FindByContext(ctx context.Context, keyContext string) (*entities.EncryptionKey, error) {
	filter := model.Map{"context": keyContext}
	return r.model.FindOne(ctx, &filter)
}

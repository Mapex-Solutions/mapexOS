package mongo

import (
	"context"
	"fmt"

	"assets/src/modules/mqttcerts/domain/entities"
	"assets/src/modules/mqttcerts/domain/repositories"

	manager "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/manager"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// NewRevokedRepository builds the Mongo adapter and ensures indexes:
//   - serial UNIQUE
//   - assetUUID + revokedAt desc
//   - orgId + revokedAt desc
//   - revokedAt TTL 30 days (auto-delete via Mongo)
func NewRevokedRepository(m *manager.MongoManager) repositories.RevokedRepository {
	ttl := RevokedTTLSeconds
	mdl := model.New[entities.RevokedCertificate](m.GetDatabase(), CollectionName, model.Config{
		Indexes: []model.IndexDefinition{
			{Name: "serial_unique", Keys: map[string]int{"serial": 1}, Unique: true},
			{Name: "assetUUID_revokedAt", Keys: map[string]int{"assetUUID": 1, "revokedAt": -1}},
			{Name: "orgId_revokedAt", Keys: map[string]int{"orgId": 1, "revokedAt": -1}},
			{Name: "revokedAt_ttl", Keys: map[string]int{"revokedAt": 1}, ExpireAfterSeconds: &ttl},
		},
	})
	logger.Info(fmt.Sprintf("[REPO:RevokedCertificate] indexes ensured collection=%s ttl=%ds", CollectionName, ttl))
	return &revokedRepository{model: mdl}
}

func (r *revokedRepository) Create(ctx context.Context, e *entities.RevokedCertificate) error {
	_, err := r.model.CreateOne(ctx, e)
	return err
}

func (r *revokedRepository) FindByAssetUUID(ctx context.Context, assetUUID string) ([]*entities.RevokedCertificate, error) {
	filter := model.Map{"assetUUID": assetUUID}
	cur, err := r.model.DIRECT().Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var rows []*entities.RevokedCertificate
	for cur.Next(ctx) {
		var e entities.RevokedCertificate
		if err := cur.Decode(&e); err != nil {
			return nil, err
		}
		rows = append(rows, &e)
	}
	return rows, nil
}

func (r *revokedRepository) FindBySerial(ctx context.Context, serial string) (*entities.RevokedCertificate, error) {
	filter := model.Map{"serial": serial}
	return r.model.FindOne(ctx, &filter)
}

func (r *revokedRepository) DeleteByAssetUUID(ctx context.Context, assetUUID string) (int64, error) {
	res, err := r.model.DIRECT().DeleteMany(ctx, model.Map{"assetUUID": assetUUID})
	if err != nil {
		return 0, err
	}
	return res.DeletedCount, nil
}

func (r *revokedRepository) CountByAssetUUID(ctx context.Context, assetUUID string) (int64, error) {
	return r.model.DIRECT().CountDocuments(ctx, model.Map{"assetUUID": assetUUID})
}

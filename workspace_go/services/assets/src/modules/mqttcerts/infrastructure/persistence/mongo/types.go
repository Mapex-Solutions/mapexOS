package mongo

import (
	"assets/src/modules/mqttcerts/domain/entities"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

// revokedRepository is the Mongo adapter for the RevokedRepository port.
// Backed by the shared `model.Model[T]` wrapper from mapexGoKit; the
// constructor ensures indexes (serial UNIQUE, assetUUID+revokedAt,
// orgId+revokedAt) and TTL 30d on revokedAt.
type revokedRepository struct {
	model *model.Model[entities.RevokedCertificate]
}

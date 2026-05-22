package collection

import (
	"http_gateway/src/modules/datasources/domain/entities"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

// repository is the MongoDB-backed adapter implementing repositories.DataSourceRepository.
// It wraps the generic *model.Model[entities.DataSource] used by all repository methods
// declared in dataSource_repository.go.
type repository struct {
	model *model.Model[entities.DataSource]
}

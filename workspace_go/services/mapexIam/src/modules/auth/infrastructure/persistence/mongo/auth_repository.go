package collection

import (
	"context"

	"mapexIam/src/modules/auth/domain/entities"
	"mapexIam/src/modules/auth/domain/repositories"

	manager "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/manager"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

// New creates and returns a generic repository for the XXX entity.
// It accepts a *MongoManager to obtain the database connection,
//
// Then calls model.New to initialize a Model[XXX],
// targeting the configured collection name with default settings.
func New(m *manager.MongoManager) repositories.AuthRepository {
	mdl := model.New[entities.Auth](m.GetDatabase(), collectionName, model.Config{})
	return &repository{model: mdl}
}

//
// START REPOSITORY METHODS
//

// Create inserts a new Auth entity into the repository.
// It accepts a context for cancellation and timeouts, and a pointer to a Auth entity to be persisted.
// Internally, it calls the underlying model’s CreateOne method to store the document in MongoDB.
// It returns the created Auth (populated with any database-assigned fields) and an error if something goes wrong.
func (r *repository) Login(ctx context.Context, u *entities.Auth) (*entities.Auth, error) {
	auth, err := r.model.CreateOne(ctx, u)
	return auth, err
}

// Compile-time check to ensure repository implements AuthRepository
var _ repositories.AuthRepository = (*repository)(nil)

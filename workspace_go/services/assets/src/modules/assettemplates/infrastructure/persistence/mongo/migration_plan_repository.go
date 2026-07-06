package collection

import (
	"context"

	"assets/src/modules/assettemplates/domain/entities"
	"assets/src/modules/assettemplates/domain/repositories"
	"assets/src/modules/assettemplates/infrastructure/persistence/mongo/constants"

	manager "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/manager"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

// migrationPlanRepository is the MongoDB-backed adapter implementing
// repositories.MigrationPlanRepository.
type migrationPlanRepository struct {
	model *model.Model[entities.MigrationPlan]
}

// NewMigrationPlanRepository creates a Mongo-backed MigrationPlanRepository.
func NewMigrationPlanRepository(m *manager.MongoManager) repositories.MigrationPlanRepository {
	mdl := model.New[entities.MigrationPlan](m.GetDatabase(), constants.MigrationPlansCollection, model.Config{
		Indexes: constants.MigrationPlansIndexes,
	})
	return &migrationPlanRepository{model: mdl}
}

func (r *migrationPlanRepository) Create(ctx context.Context, p *entities.MigrationPlan) (*entities.MigrationPlan, error) {
	return r.model.CreateOne(ctx, p)
}

func (r *migrationPlanRepository) FindById(ctx context.Context, id *string) (*entities.MigrationPlan, error) {
	retData, _ := r.model.FindByID(ctx, *id)
	return retData, nil
}

func (r *migrationPlanRepository) FindByIdAndUpdate(ctx context.Context, id *string, payload map[string]any) (*entities.MigrationPlan, error) {
	_id, _ := model.ToObjectID(*id)
	query := model.Map{"_id": _id}
	update := model.Map{"$set": payload}

	returnDoc := model.ReturnDoc(1)
	options := model.CommonOpts{ReturnDocument: &returnDoc}

	retData, _ := r.model.FindOneAndUpdate(ctx, &query, &update, &options)
	return retData, nil
}

func (r *migrationPlanRepository) FindWithFilters(
	ctx context.Context,
	filters model.Map,
	pagination *model.PaginationOpts,
	projection model.Map,
) (*model.PaginatedResult[entities.MigrationPlan], error) {
	opts := &model.CommonOpts{}
	if len(projection) > 0 {
		opts.Projection = projection
	}
	return r.model.FindByOffset(ctx, filters, pagination, opts)
}

func (r *migrationPlanRepository) IncrementCounter(ctx context.Context, planID, field string, delta int) error {
	_id, _ := model.ToObjectID(planID)
	_, err := r.model.FindAndUpdateMany(ctx, model.Map{"_id": _id}, model.Map{"$inc": model.Map{field: delta}})
	return err
}

// Compile-time check that migrationPlanRepository implements the port.
var _ repositories.MigrationPlanRepository = (*migrationPlanRepository)(nil)

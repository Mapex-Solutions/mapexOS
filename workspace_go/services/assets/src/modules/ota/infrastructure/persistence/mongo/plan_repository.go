package collection

import (
	"context"

	"assets/src/modules/ota/domain/entities"
	"assets/src/modules/ota/domain/repositories"

	manager "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/manager"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

// collectionPlans is the MongoDB collection for OTA plans.
const collectionPlans = "ota_plans"

// planRepository is the MongoDB-backed adapter implementing
// repositories.OTAPlanRepository.
type planRepository struct {
	model *model.Model[entities.OTAPlan]
}

// NewPlanRepository creates a Mongo-backed OTAPlanRepository.
func NewPlanRepository(m *manager.MongoManager) repositories.OTAPlanRepository {
	mdl := model.New[entities.OTAPlan](m.GetDatabase(), collectionPlans, model.Config{
		Indexes: []model.IndexDefinition{
			{Name: "idx_org", Keys: map[string]int{"orgId": 1}},
			{Name: "idx_status", Keys: map[string]int{"status": 1}},
			{Name: "idx_source_template", Keys: map[string]int{"sourceTemplateId": 1}},
			{Name: "idx_firmware", Keys: map[string]int{"firmwareId": 1}},
		},
	})
	return &planRepository{model: mdl}
}

func (r *planRepository) Create(ctx context.Context, p *entities.OTAPlan) (*entities.OTAPlan, error) {
	return r.model.CreateOne(ctx, p)
}

func (r *planRepository) FindById(ctx context.Context, id *string) (*entities.OTAPlan, error) {
	retData, _ := r.model.FindByID(ctx, *id)
	return retData, nil
}

func (r *planRepository) FindByIdAndUpdate(ctx context.Context, id *string, payload map[string]any) (*entities.OTAPlan, error) {
	_id, _ := model.ToObjectID(*id)
	query := model.Map{"_id": _id}
	update := model.Map{"$set": payload}

	returnDoc := model.ReturnDoc(1)
	options := model.CommonOpts{ReturnDocument: &returnDoc}

	retData, _ := r.model.FindOneAndUpdate(ctx, &query, &update, &options)
	return retData, nil
}

func (r *planRepository) DeleteById(ctx context.Context, id *string) error {
	_id, _ := model.ToObjectID(*id)
	query := model.Map{"_id": _id}
	return r.model.DeleteOne(ctx, &query)
}

func (r *planRepository) FindWithFilters(
	ctx context.Context,
	filters model.Map,
	pagination *model.PaginationOpts,
	projection model.Map,
) (*model.PaginatedResult[entities.OTAPlan], error) {
	opts := &model.CommonOpts{}
	if len(projection) > 0 {
		opts.Projection = projection
	}
	return r.model.FindByOffset(ctx, filters, pagination, opts)
}

func (r *planRepository) IncrementCounter(ctx context.Context, planID, field string, delta int) error {
	_id, _ := model.ToObjectID(planID)
	_, err := r.model.FindAndUpdateMany(ctx, model.Map{"_id": _id}, model.Map{"$inc": model.Map{field: delta}})
	return err
}

// Compile-time check that planRepository implements the port.
var _ repositories.OTAPlanRepository = (*planRepository)(nil)

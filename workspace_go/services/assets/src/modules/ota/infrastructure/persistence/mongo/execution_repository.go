package collection

import (
	"context"

	"assets/src/modules/ota/domain/entities"
	"assets/src/modules/ota/domain/repositories"

	manager "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/manager"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

// collectionExecutions is the MongoDB collection for per-asset OTA executions.
const collectionExecutions = "ota_executions"

// executionRepository is the MongoDB-backed adapter implementing
// repositories.OTAExecutionRepository.
type executionRepository struct {
	model *model.Model[entities.OTAExecution]
}

// NewExecutionRepository creates a Mongo-backed OTAExecutionRepository. The
// collection is indexed by planId and assetId (plus a compound planId+state for
// fast per-state counts during reconciliation).
func NewExecutionRepository(m *manager.MongoManager) repositories.OTAExecutionRepository {
	mdl := model.New[entities.OTAExecution](m.GetDatabase(), collectionExecutions, model.Config{
		Indexes: []model.IndexDefinition{
			{Name: "idx_plan", Keys: map[string]int{"planId": 1}},
			{Name: "idx_asset", Keys: map[string]int{"assetId": 1}},
			{Name: "idx_plan_state", Keys: map[string]int{"planId": 1, "state": 1}},
		},
	})
	return &executionRepository{model: mdl}
}

func (r *executionRepository) BulkInsert(ctx context.Context, execs []*entities.OTAExecution) (int64, error) {
	if len(execs) == 0 {
		return 0, nil
	}
	items := make([]entities.OTAExecution, 0, len(execs))
	for _, e := range execs {
		items = append(items, *e)
	}
	created, err := r.model.CreateMany(ctx, items)
	if err != nil {
		return 0, err
	}
	return int64(len(created)), nil
}

func (r *executionRepository) FindById(ctx context.Context, id *string) (*entities.OTAExecution, error) {
	retData, _ := r.model.FindByID(ctx, *id)
	return retData, nil
}

func (r *executionRepository) FindByIdAndUpdate(ctx context.Context, id *string, payload map[string]any) (*entities.OTAExecution, error) {
	_id, _ := model.ToObjectID(*id)
	query := model.Map{"_id": _id}
	update := model.Map{"$set": payload}

	returnDoc := model.ReturnDoc(1)
	options := model.CommonOpts{ReturnDocument: &returnDoc}

	retData, _ := r.model.FindOneAndUpdate(ctx, &query, &update, &options)
	return retData, nil
}

func (r *executionRepository) FindByPlan(
	ctx context.Context,
	planID *string,
	filters model.Map,
	pagination *model.PaginationOpts,
) (*model.PaginatedResult[entities.OTAExecution], error) {
	_id, _ := model.ToObjectID(*planID)
	if filters == nil {
		filters = model.Map{}
	}
	filters["planId"] = _id
	return r.model.FindByOffset(ctx, filters, pagination, &model.CommonOpts{})
}

func (r *executionRepository) FindWithFilters(
	ctx context.Context,
	filters model.Map,
	pagination *model.PaginationOpts,
	projection model.Map,
) (*model.PaginatedResult[entities.OTAExecution], error) {
	opts := &model.CommonOpts{}
	if len(projection) > 0 {
		opts.Projection = projection
	}
	return r.model.FindByOffset(ctx, filters, pagination, opts)
}

func (r *executionRepository) CountByState(ctx context.Context, planID *string, state entities.ExecutionState) (int64, error) {
	_id, _ := model.ToObjectID(*planID)
	return r.model.DIRECT().CountDocuments(ctx, model.Map{"planId": _id, "state": string(state)})
}

func (r *executionRepository) UpdateMany(ctx context.Context, filter model.Map, update model.Map) (int64, error) {
	result, err := r.model.FindAndUpdateMany(ctx, filter, update)
	if err != nil {
		return 0, err
	}
	return result.MatchedCount, nil
}

// Compile-time check that executionRepository implements the port.
var _ repositories.OTAExecutionRepository = (*executionRepository)(nil)

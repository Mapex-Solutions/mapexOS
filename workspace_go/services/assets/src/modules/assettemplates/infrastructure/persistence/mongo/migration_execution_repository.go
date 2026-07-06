package collection

import (
	"context"
	"time"

	"assets/src/modules/assettemplates/domain/entities"
	"assets/src/modules/assettemplates/domain/repositories"
	"assets/src/modules/assettemplates/infrastructure/persistence/mongo/constants"

	manager "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/manager"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

// migrationExecutionRepository is the MongoDB-backed adapter implementing
// repositories.MigrationExecutionRepository.
type migrationExecutionRepository struct {
	model *model.Model[entities.MigrationExecution]
}

// NewMigrationExecutionRepository creates a Mongo-backed MigrationExecutionRepository.
// The collection is indexed by planId and status (plus a compound planId+status
// for fast per-status counts during reconciliation).
func NewMigrationExecutionRepository(m *manager.MongoManager) repositories.MigrationExecutionRepository {
	mdl := model.New[entities.MigrationExecution](m.GetDatabase(), constants.MigrationExecutionsCollection, model.Config{
		Indexes: constants.MigrationExecutionsIndexes,
	})
	return &migrationExecutionRepository{model: mdl}
}

func (r *migrationExecutionRepository) BulkInsert(ctx context.Context, execs []*entities.MigrationExecution) (int, error) {
	if len(execs) == 0 {
		return 0, nil
	}
	items := make([]entities.MigrationExecution, 0, len(execs))
	for _, e := range execs {
		items = append(items, *e)
	}
	created, err := r.model.CreateMany(ctx, items)
	if err != nil {
		return 0, err
	}
	return len(created), nil
}

func (r *migrationExecutionRepository) FindByPlan(
	ctx context.Context,
	planID *string,
	filters model.Map,
	pagination *model.PaginationOpts,
) (*model.PaginatedResult[entities.MigrationExecution], error) {
	_id, _ := model.ToObjectID(*planID)
	if filters == nil {
		filters = model.Map{}
	}
	filters["planId"] = _id
	return r.model.FindByOffset(ctx, filters, pagination, &model.CommonOpts{})
}

func (r *migrationExecutionRepository) FindByIdAndUpdate(ctx context.Context, id *string, payload map[string]any) (*entities.MigrationExecution, error) {
	_id, _ := model.ToObjectID(*id)
	query := model.Map{"_id": _id}
	update := model.Map{"$set": payload}

	returnDoc := model.ReturnDoc(1)
	options := model.CommonOpts{ReturnDocument: &returnDoc}

	retData, _ := r.model.FindOneAndUpdate(ctx, &query, &update, &options)
	return retData, nil
}

func (r *migrationExecutionRepository) CountByPlanAndStatus(ctx context.Context, planID string, status entities.ExecStatus) (int64, error) {
	_id, _ := model.ToObjectID(planID)
	return r.model.DIRECT().CountDocuments(ctx, model.Map{"planId": _id, "status": string(status)})
}

func (r *migrationExecutionRepository) UpdateByPlanAndAsset(ctx context.Context, planID, assetID model.ObjectId, fields map[string]any) error {
	filter := model.Map{"planId": planID, "assetId": assetID}
	_, err := r.model.DIRECT().UpdateOne(ctx, filter, model.Map{"$set": fields})
	return err
}

func (r *migrationExecutionRepository) CancelPending(ctx context.Context, planID model.ObjectId) (int64, error) {
	filter := model.Map{"planId": planID, "status": string(entities.ExecPending)}
	update := model.Map{"$set": model.Map{"status": string(entities.ExecCancelled), "updated": time.Now()}}
	result, err := r.model.FindAndUpdateMany(ctx, filter, update)
	if err != nil {
		return 0, err
	}
	return result.MatchedCount, nil
}

// Compile-time check that migrationExecutionRepository implements the port.
var _ repositories.MigrationExecutionRepository = (*migrationExecutionRepository)(nil)

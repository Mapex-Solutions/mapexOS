package collection

import (
	"context"

	"assets/src/modules/ota/domain/entities"
	"assets/src/modules/ota/domain/repositories"

	manager "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/manager"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

// collectionFirmware is the MongoDB collection for firmware artifacts.
const collectionFirmware = "ota_firmware"

// firmwareRepository is the MongoDB-backed adapter implementing
// repositories.FirmwareRepository.
type firmwareRepository struct {
	model *model.Model[entities.Firmware]
}

// NewFirmwareRepository creates a Mongo-backed FirmwareRepository.
func NewFirmwareRepository(m *manager.MongoManager) repositories.FirmwareRepository {
	mdl := model.New[entities.Firmware](m.GetDatabase(), collectionFirmware, model.Config{
		Indexes: []model.IndexDefinition{
			{Name: "idx_org", Keys: map[string]int{"orgId": 1}},
			{Name: "idx_target_template", Keys: map[string]int{"targetTemplateId": 1}},
			{Name: "idx_status", Keys: map[string]int{"status": 1}},
		},
	})
	return &firmwareRepository{model: mdl}
}

func (r *firmwareRepository) Create(ctx context.Context, f *entities.Firmware) (*entities.Firmware, error) {
	return r.model.CreateOne(ctx, f)
}

func (r *firmwareRepository) FindById(ctx context.Context, id *string) (*entities.Firmware, error) {
	retData, _ := r.model.FindByID(ctx, *id)
	return retData, nil
}

func (r *firmwareRepository) FindByIdAndUpdate(ctx context.Context, id *string, payload map[string]any) (*entities.Firmware, error) {
	_id, _ := model.ToObjectID(*id)
	query := model.Map{"_id": _id}
	update := model.Map{"$set": payload}

	returnDoc := model.ReturnDoc(1)
	options := model.CommonOpts{ReturnDocument: &returnDoc}

	retData, _ := r.model.FindOneAndUpdate(ctx, &query, &update, &options)
	return retData, nil
}

func (r *firmwareRepository) DeleteById(ctx context.Context, id *string) error {
	_id, _ := model.ToObjectID(*id)
	query := model.Map{"_id": _id}
	return r.model.DeleteOne(ctx, &query)
}

func (r *firmwareRepository) FindWithFilters(
	ctx context.Context,
	filters model.Map,
	pagination *model.PaginationOpts,
	projection model.Map,
) (*model.PaginatedResult[entities.Firmware], error) {
	opts := &model.CommonOpts{}
	if len(projection) > 0 {
		opts.Projection = projection
	}
	return r.model.FindByOffset(ctx, filters, pagination, opts)
}

// Compile-time check that firmwareRepository implements the port.
var _ repositories.FirmwareRepository = (*firmwareRepository)(nil)

package services

import (
	ctx "context"
	"fmt"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/customErrors"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/status"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// parseClassificationId converts the raw classification list ID to ObjectId.
// Returns BAD_REQUEST when the ID is malformed so callers can fail fast.
func (s *AssetTemplateService) parseClassificationId(rawId string, label string) (model.ObjectId, error) {
	objectId, err := model.ToObjectID(rawId)
	if err != nil {
		return model.ObjectId{}, &customErrors.ServerCustomError{
			Code:   status.BAD_REQUEST,
			Errors: []string{fmt.Sprintf("Invalid %s ID", label)},
		}
	}
	return objectId, nil
}

// denormalizeTemplateField runs the cross-template UpdateMany that pushes the
// new classification name into every template referencing the given list id.
func (s *AssetTemplateService) denormalizeTemplateField(c ctx.Context, idField string, idValue model.ObjectId, nameField string, newName string) (int64, error) {
	filter := model.Map{idField: idValue}
	update := model.Map{"$set": model.Map{nameField: newName}}
	return s.deps.AssetTemplateRepo.UpdateMany(c, filter, update)
}

// logClassificationDenorm emits the post-denormalization audit log so the
// number of touched templates stays observable.
func (s *AssetTemplateService) logClassificationDenorm(label string, matched int64) {
	logger.Info(fmt.Sprintf("[SERVICE:AssetTemplate] Updated %s name for %d templates", label, matched))
}

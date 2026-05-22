package services

import (
	"assets/src/modules/assettemplates/application/dtos"
	"assets/src/modules/assettemplates/domain/entities"
	ctx "context"
	"time"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/customErrors"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/status"
)

// recordTemplateOp emits the count + duration metrics for one CRUD attempt
// so every exit path of every public method stays observability-consistent.
func (s *AssetTemplateService) recordTemplateOp(op, status string, start time.Time) {
	s.deps.Metrics.TemplateOperations.WithLabelValues(op, status).Inc()
	s.deps.Metrics.TemplateOperationDuration.WithLabelValues(op).Observe(time.Since(start).Seconds())
}

// fetchTemplateById is an internal method that fetches a template entity by ID.
// Used by GetAssetTemplateById, GetTemplateByIdForCacheFallback, and DeleteAssetTemplateById
// to avoid duplication of fetch + validation logic.
func (s *AssetTemplateService) fetchTemplateById(ctx ctx.Context, templateId *string) (*entities.Assettemplate, error) {
	template, err := s.deps.AssetTemplateRepo.FindById(ctx, templateId)
	if err != nil {
		return nil, err
	}
	if template == nil {
		return nil, &customErrors.ServerCustomError{Code: status.NOT_FOUND, Errors: []string{"Asset Template not found"}}
	}
	return template, nil
}

// populateOrgContext sets orgId and pathKey from RequestContext into the DTO.
// This helper is used for both Template and Local resource creation.
//
// Parameters:
//   - requestContext: Contains OrgContext (selected orgId) and OrgContextData (pathKey)
//   - dto: The DTO to populate
func (s *AssetTemplateService) populateOrgContext(requestContext *reqCtx.RequestContext, dto *dtos.AssetTemplateCreateDTO) {
	if requestContext.OrgContext != nil && *requestContext.OrgContext != "" {
		orgObjectId, err := model.ToObjectID(*requestContext.OrgContext)
		if err == nil {
			dto.OrgID = &orgObjectId
		}
	}

	if requestContext.OrgContextData != nil && requestContext.OrgContextData.PathKey != "" {
		pathKey := requestContext.OrgContextData.PathKey
		dto.PathKey = &pathKey
	}
}

// processDynamicFieldsUpdate handles the FieldId assignment for UPDATE operations.
// - Preserves existing FieldId for fields with same name
// - Assigns new FieldId for new fields
// - Marks removed fields as deprecated (Status=0) to preserve historical data
//
// Parameters:
//   - existingTemplate: The current template from database
//   - incomingFields: The fields from the update DTO
//
// Returns:
//   - DynamicFieldsResult with processed fields and updated NextFieldId
func (s *AssetTemplateService) processDynamicFieldsUpdate(
	existingTemplate *entities.Assettemplate,
	incomingFields []dtos.DynamicField,
) DynamicFieldsResult {
	// Build map of existing fields: fieldName -> DynamicField
	existingFieldsMap := make(map[string]entities.DynamicField)
	for _, field := range existingTemplate.DynamicFields {
		existingFieldsMap[field.Field] = field
	}

	// Track which existing fields are still present
	presentFields := make(map[string]bool)

	// Start with existing NextFieldId or 1 if not set
	nextFieldId := existingTemplate.NextFieldId
	if nextFieldId == 0 {
		nextFieldId = 1
	}

	// Process incoming fields
	processedFields := make([]entities.DynamicField, 0, len(incomingFields))
	for _, incoming := range incomingFields {
		if existing, found := existingFieldsMap[incoming.Field]; found {
			// Existing field: preserve FieldId, update other properties
			processedFields = append(processedFields, entities.DynamicField{
				FieldId:       existing.FieldId,
				Field:         incoming.Field,
				Value:         incoming.Value,
				Type:          incoming.Type,
				Status:        1, // Keep active
				LatitudePath:  incoming.LatitudePath,
				LongitudePath: incoming.LongitudePath,
			})
			presentFields[incoming.Field] = true
		} else {
			// New field: assign new FieldId
			processedFields = append(processedFields, entities.DynamicField{
				FieldId:       nextFieldId,
				Field:         incoming.Field,
				Value:         incoming.Value,
				Type:          incoming.Type,
				Status:        1, // Active
				LatitudePath:  incoming.LatitudePath,
				LongitudePath: incoming.LongitudePath,
			})
			nextFieldId++
		}
	}

	// Mark removed fields as deprecated (Status=0)
	// This preserves historical query capability
	for _, existing := range existingTemplate.DynamicFields {
		if !presentFields[existing.Field] && existing.Status == 1 {
			processedFields = append(processedFields, entities.DynamicField{
				FieldId:       existing.FieldId,
				Field:         existing.Field,
				Value:         existing.Value,
				Type:          existing.Type,
				Status:        0, // Deprecated
				LatitudePath:  existing.LatitudePath,
				LongitudePath: existing.LongitudePath,
			})
		}
	}

	return DynamicFieldsResult{
		Fields:      processedFields,
		NextFieldId: nextFieldId,
	}
}

// convertIdFieldsInMap converts string ID fields to ObjectId in the update map.
// Helper to avoid duplication in update methods.
func (s *AssetTemplateService) convertIdFieldsInMap(fieldsToUpdate map[string]interface{}) {
	if categoryId, ok := fieldsToUpdate["categoryId"].(string); ok && categoryId != "" {
		id, _ := model.ToObjectID(categoryId)
		fieldsToUpdate["categoryId"] = id
	}
	if manufacturerId, ok := fieldsToUpdate["manufacturerId"].(string); ok && manufacturerId != "" {
		id, _ := model.ToObjectID(manufacturerId)
		fieldsToUpdate["manufacturerId"] = id
	}
	if modelId, ok := fieldsToUpdate["modelId"].(string); ok && modelId != "" {
		id, _ := model.ToObjectID(modelId)
		fieldsToUpdate["modelId"] = id
	}
}

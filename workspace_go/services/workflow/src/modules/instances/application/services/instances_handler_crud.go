package services

import (
	"fmt"

	"workflow/src/modules/instances/application/dtos"
	"workflow/src/modules/instances/domain/entities"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	"github.com/Mapex-Solutions/mapexGoKit/utils/mapper"
)

// buildInstanceEntity translates the create DTO into a fresh WorkflowInstance
// entity. Defaults applied here (IsSystem=false, Enabled=true) cannot be set
// by the caller — the API surface is intentionally narrowed.
func (s *InstancesService) buildInstanceEntity(dto *dtos.InstanceCreateDTO) *entities.WorkflowInstance {
	return &entities.WorkflowInstance{
		DefinitionID:      dto.DefinitionID,
		DefinitionVersion: dto.DefinitionVersion,
		DefinitionName:    dto.DefinitionName,
		Name:              dto.Name,
		Description:       dto.Description,
		ExternalInputs:    dto.ExternalInputs,
		IsSystem:          false,
		IsTemplate:        dto.IsTemplate,
		UniqueExecution:   dto.UniqueExecution,
		WorkflowUUID:      dto.WorkflowUUID,
		Enabled:           true,
	}
}

// applyOrgContextToInstance copies the resolved org id and pathKey from the
// request-context (set by the coverage middleware) onto the entity so the
// stored document carries the same multi-tenant labels as future queries.
func (s *InstancesService) applyOrgContextToInstance(entity *entities.WorkflowInstance, requestContext *reqCtx.RequestContext) {
	if requestContext.OrgContext != nil && *requestContext.OrgContext != "" {
		if orgObjectId, err := model.ToObjectID(*requestContext.OrgContext); err == nil {
			entity.OrgID = &orgObjectId
		}
	}
	if requestContext.OrgContextData != nil && requestContext.OrgContextData.PathKey != "" {
		entity.PathKey = requestContext.OrgContextData.PathKey
	}
}

// buildInstanceUpdatePayload collects the not-nil fields of the update DTO
// into a map suitable for FindByIdAndUpdate. Nil-valued DTO fields are
// preserved on the stored document — partial updates are by-design.
func (s *InstancesService) buildInstanceUpdatePayload(dto *dtos.InstanceUpdateDTO) map[string]any {
	payload := make(map[string]any)
	if dto.Name != nil {
		payload["name"] = *dto.Name
	}
	if dto.Description != nil {
		payload["description"] = *dto.Description
	}
	if dto.ExternalInputs != nil {
		payload["externalInputs"] = dto.ExternalInputs
	}
	if dto.Enabled != nil {
		payload["enabled"] = *dto.Enabled
	}
	return payload
}

// toInstanceResponse converts a persisted entity to the wire DTO. The shared
// mapper utility handles bson→json field-name translation; failures here mean
// a struct-tag mismatch and surface as 500s for fast diagnosis.
func (s *InstancesService) toInstanceResponse(entity *entities.WorkflowInstance) (*dtos.InstanceResponse, error) {
	resp, err := mapper.EntityToDto[entities.WorkflowInstance, dtos.InstanceResponse](entity)
	if err != nil {
		return nil, fmt.Errorf("failed to map entity to response: %w", err)
	}
	return resp, nil
}

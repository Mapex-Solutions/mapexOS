package services

import (
	"mapexIam/src/modules/lists/application/dtos"

	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	customErrors "github.com/Mapex-Solutions/mapexGoKit/microservices/http/customErrors"
	httpStatus "github.com/Mapex-Solutions/mapexGoKit/microservices/http/status"
	"github.com/Mapex-Solutions/mapexGoKit/utils/orgfilter"
)

// applyListScope decides between system / template-ancestor / org-local
// scope and applies the corresponding multi-tenant fields to the create
// DTO. Returns 4xx when the caller cannot create the requested scope.
func (s *ListService) applyListScope(rc *reqCtx.RequestContext, dto *dtos.ListCreateDTO) error {
	if dto.IsSystem {
		dto.OrgID = nil
		dto.PathKey = nil
		return nil
	}
	if dto.IsTemplate {
		if err := orgfilter.ValidateTemplateCreation(rc.OrgContextData.PathKey); err != nil {
			return &customErrors.ServerCustomError{Code: httpStatus.FORBIDDEN, Errors: []string{err.Error()}}
		}
		populateListOrgContext(rc, dto)
		return nil
	}
	if err := orgfilter.ValidateOrgContextForNonSystem(rc); err != nil {
		return &customErrors.ServerCustomError{Code: httpStatus.BAD_REQUEST, Errors: []string{err.Error()}}
	}
	populateListOrgContext(rc, dto)
	return nil
}

// populateListOrgContext copies orgId + pathKey from request context onto
// the DTO so the entity inherits multi-tenant scoping.
func populateListOrgContext(rc *reqCtx.RequestContext, dto *dtos.ListCreateDTO) {
	if rc.OrgContext != nil && *rc.OrgContext != "" {
		if orgObjectId, err := model.ToObjectID(*rc.OrgContext); err == nil {
			dto.OrgID = &orgObjectId
		}
	}
	if rc.OrgContextData != nil && rc.OrgContextData.PathKey != "" {
		pathKey := rc.OrgContextData.PathKey
		dto.PathKey = &pathKey
	}
}

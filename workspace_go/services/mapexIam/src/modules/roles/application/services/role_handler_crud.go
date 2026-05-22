package services

import (
	"fmt"

	events "mapexIam/src/modules/cache_invalidation/application/events"
	"mapexIam/src/modules/roles/application/dtos"

	cacheInvalidation "github.com/Mapex-Solutions/MapexOS/contracts/services/mapexIam/cache_invalidation"
	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	customErrors "github.com/Mapex-Solutions/mapexGoKit/microservices/http/customErrors"
	httpStatus "github.com/Mapex-Solutions/mapexGoKit/microservices/http/status"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
	"github.com/Mapex-Solutions/mapexGoKit/utils/orgfilter"
)

// applyRoleScope decides between system / template-ancestor / org-local
// scope and applies the corresponding multi-tenant fields to the create
// DTO. Returns 4xx when the caller cannot create the requested scope.
func (s *RoleService) applyRoleScope(rc *reqCtx.RequestContext, dto *dtos.CreateRoleDto) error {
	if dto.IsSystem {
		dto.OrgID = nil
		dto.PathKey = nil
		return nil
	}
	if dto.IsTemplate {
		if err := orgfilter.ValidateTemplateCreation(rc.OrgContextData.PathKey); err != nil {
			return &customErrors.ServerCustomError{Code: httpStatus.FORBIDDEN, Errors: []string{err.Error()}}
		}
		populateRoleOrgContext(rc, dto)
		return nil
	}
	if err := orgfilter.ValidateOrgContextForNonSystem(rc); err != nil {
		return &customErrors.ServerCustomError{Code: httpStatus.BAD_REQUEST, Errors: []string{err.Error()}}
	}
	populateRoleOrgContext(rc, dto)
	return nil
}

// populateRoleOrgContext copies orgId + pathKey from request context onto
// the DTO so the entity inherits multi-tenant scoping.
func populateRoleOrgContext(rc *reqCtx.RequestContext, dto *dtos.CreateRoleDto) {
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

// publishRolePermissionsChanged emits the NATS cache-invalidation event
// when a role update changed the permission list. Best-effort; failures
// are logged but never block the caller.
func (s *RoleService) publishRolePermissionsChanged(roleId string, oldPerms, newPerms []string) {
	logger.Info(fmt.Sprintf("[SERVICE:Role] Permissions changed for role=%s - publishing cache invalidation event", roleId))
	event := events.NewRolePermissionsChangedEvent(roleId, oldPerms, newPerms, "")
	subject := fmt.Sprintf(cacheInvalidation.RolePermissionsChangedSubjectFormat, roleId)
	if err := s.deps.NatsBus.Publish(natsModel.PublishConfig{Subject: subject, Data: event}); err != nil {
		logger.Error(err, fmt.Sprintf("[SERVICE:Role] Failed to publish RolePermissionsChangedEvent for role=%s", roleId))
		return
	}
	logger.Info(fmt.Sprintf("[SERVICE:Role] Published RolePermissionsChangedEvent for role=%s to subject=%s", roleId, subject))
}

// publishRoleDeleted emits the NATS cache-invalidation event after a
// role is deleted so consumers drop every cached principal that
// referenced this role.
func (s *RoleService) publishRoleDeleted(roleId string) {
	logger.Info(fmt.Sprintf("[SERVICE:Role] Role deleted role=%s - publishing cache invalidation event", roleId))
	event := events.NewRoleDeletedEvent(roleId, "")
	subject := fmt.Sprintf(cacheInvalidation.RoleDeletedSubjectFormat, roleId)
	if err := s.deps.NatsBus.Publish(natsModel.PublishConfig{Subject: subject, Data: event}); err != nil {
		logger.Error(err, fmt.Sprintf("[SERVICE:Role] Failed to publish RoleDeletedEvent for role=%s", roleId))
		return
	}
	logger.Info(fmt.Sprintf("[SERVICE:Role] Published RoleDeletedEvent for role=%s to subject=%s", roleId, subject))
}

package services

import (
	ctx "context"
	"fmt"

	"mapexIam/src/modules/memberships/application/dtos"
	"mapexIam/src/modules/memberships/domain/entities"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	customErrors "github.com/Mapex-Solutions/mapexGoKit/microservices/http/customErrors"
	httpStatus "github.com/Mapex-Solutions/mapexGoKit/microservices/http/status"
)

// bindMembershipTenantFields resolves the target organization and copies
// its multi-tenant fields (OrgID, OrgPathKey, CustomerID) onto the
// freshly mapped entity. Missing org -> 404; missing/invalid org id ->
// 400 so admins can distinguish real input errors from data drift.
func (s *MembershipService) bindMembershipTenantFields(c ctx.Context, entity *entities.Membership, dto *dtos.CreateMembershipDto) error {
	org, err := s.deps.OrgService.GetOrganizationById(c, &dto.OrgID)
	if err != nil || org == nil {
		return &customErrors.ServerCustomError{
			Code:   httpStatus.NOT_FOUND,
			Errors: []string{"Organization not found"},
		}
	}
	if org.ID == nil {
		return &customErrors.ServerCustomError{
			Code:   httpStatus.BAD_REQUEST,
			Errors: []string{"Organization ID is missing"},
		}
	}
	orgObjectID, err := model.ToObjectID(*org.ID)
	if err != nil {
		return &customErrors.ServerCustomError{
			Code:   httpStatus.BAD_REQUEST,
			Errors: []string{"Invalid organization ID format"},
		}
	}

	orgPathKey := ""
	if org.PathKey != nil {
		orgPathKey = *org.PathKey
	}

	entity.OrgID = &orgObjectID
	entity.OrgPathKey = orgPathKey
	entity.CustomerID = org.CustomerID
	return nil
}

// convertMembershipObjectIds translates the assignee + role ids on the
// create DTO from string form to ObjectID. Any malformed id fails fast
// with 400 so the caller sees the offending value in the error message.
func convertMembershipObjectIds(entity *entities.Membership, dto *dtos.CreateMembershipDto) error {
	assigneeObjectID, err := model.ToObjectID(dto.AssigneeID)
	if err != nil {
		return &customErrors.ServerCustomError{
			Code:   httpStatus.BAD_REQUEST,
			Errors: []string{"Invalid assignee ID format"},
		}
	}
	entity.AssigneeID = assigneeObjectID

	roleObjectIds, err := convertStringsToObjectIds(dto.RoleIds, "Invalid role ID format")
	if err != nil {
		return err
	}
	entity.RoleIds = roleObjectIds
	return nil
}

// convertStringsToObjectIds turns a []string into []ObjectID, returning
// a 400 ServerCustomError naming the first malformed id. Used by both
// create and update paths.
func convertStringsToObjectIds(ids []string, errorPrefix string) ([]model.ObjectId, error) {
	out := make([]model.ObjectId, 0, len(ids))
	for _, idStr := range ids {
		objectID, err := model.ToObjectID(idStr)
		if err != nil {
			return nil, &customErrors.ServerCustomError{
				Code:   httpStatus.BAD_REQUEST,
				Errors: []string{fmt.Sprintf("%s: %s", errorPrefix, idStr)},
			}
		}
		out = append(out, objectID)
	}
	return out, nil
}

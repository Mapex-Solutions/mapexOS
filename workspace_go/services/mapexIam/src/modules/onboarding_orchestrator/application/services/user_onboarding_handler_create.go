package services

import (
	ctx "context"
	"fmt"

	groupDtos "mapexIam/src/modules/groups/application/dtos"
	membershipDtos "mapexIam/src/modules/memberships/application/dtos"
	"mapexIam/src/modules/onboarding_orchestrator/application/dtos"

	membershipContractDtos "github.com/Mapex-Solutions/MapexOS/contracts/services/mapexIam/memberships"
	userDtos "github.com/Mapex-Solutions/MapexOS/contracts/services/mapexIam/users"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	customErrors "github.com/Mapex-Solutions/mapexGoKit/microservices/http/customErrors"
	httpStatus "github.com/Mapex-Solutions/mapexGoKit/microservices/http/status"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// resolveOnboardingOrg loads the org context's default scope for onboarding.
// Returns (orgID, defaultScope, orgName) and a 4xx ServerCustomError when
// the org id from the request context cannot be resolved.
func (s *UserOnboardingService) resolveOnboardingOrg(c ctx.Context, requestContext *reqCtx.RequestContext) (string, string, string, error) {
	orgID := *requestContext.OrgContext
	logger.Debug(fmt.Sprintf("[SERVICE:Onboarding] Using organization from context: %s", orgID))

	org, err := s.deps.OrgService.GetOrganizationById(c, &orgID)
	if err != nil || org == nil {
		if err == nil {
			err = fmt.Errorf("organization not found: %s", orgID)
		}
		logger.Error(err, "[SERVICE:Onboarding] Organization from context not found")
		return "", "", "", &customErrors.ServerCustomError{
			Code:   httpStatus.BAD_REQUEST,
			Errors: []string{fmt.Sprintf("Organization not found: %s", orgID)},
		}
	}

	defaultScope := "local"
	if org.AccessPolicy != nil && org.AccessPolicy.DefaultScope != "" {
		defaultScope = org.AccessPolicy.DefaultScope
	}

	orgName := ""
	if org.Name != nil {
		orgName = *org.Name
	}
	logger.Debug(fmt.Sprintf("[SERVICE:Onboarding] Organization: %s, DefaultScope: %s", orgName, defaultScope))
	return orgID, defaultScope, orgName, nil
}

// createOnboardingUser persists the new user via UserService inside the
// caller's transaction context.
func (s *UserOnboardingService) createOnboardingUser(txCtx ctx.Context, dto *dtos.CreateUserWithMembershipsDto) (*userDtos.UserResponse, error) {
	logger.Info(fmt.Sprintf("[SERVICE:Onboarding] [TXN] Creating user: %s %s (%s)", dto.FirstName, dto.LastName, dto.Email))
	userCreateDto := &userDtos.UserCreate{
		Email:                   dto.Email,
		Password:                dto.Password,
		ChangePasswordNextLogin: dto.ChangePasswordNextLogin,
		FirstName:               dto.FirstName,
		LastName:                dto.LastName,
		Phone:                   dto.Phone,
		JobTitle:                dto.JobTitle,
		Enabled:                 dto.Enabled,
		Avatar:                  dto.Avatar,
	}

	createdUser, err := s.deps.UserService.CreateUser(txCtx, userCreateDto)
	if err != nil {
		logger.Error(err, "[SERVICE:Onboarding] [TXN] Failed to create user - transaction will rollback")
		return nil, err
	}
	logger.Info(fmt.Sprintf("[SERVICE:Onboarding] [TXN] User created: %s", *createdUser.ID))
	return createdUser, nil
}

// attachInitialMembership creates the direct user-org membership when the
// onboarding payload carries one. Returns the resulting membership list
// (empty when no direct membership was requested).
func (s *UserOnboardingService) attachInitialMembership(
	txCtx ctx.Context,
	dto *dtos.CreateUserWithMembershipsDto,
	createdUser *userDtos.UserResponse,
	orgID, defaultScope string,
) ([]*membershipContractDtos.MembershipResponse, error) {
	createdMemberships := make([]*membershipContractDtos.MembershipResponse, 0)
	if len(dto.Memberships) == 0 {
		return createdMemberships, nil
	}

	membershipData := dto.Memberships[0]
	logger.Debug(fmt.Sprintf("[SERVICE:Onboarding] [TXN] Creating direct membership: OrgID=%s, Roles=%v, Scope=%s",
		orgID, membershipData.Roles, defaultScope))

	membershipCreateDto := &membershipDtos.CreateMembershipDto{
		AssigneeType: "user",
		AssigneeID:   createdUser.ID.Hex(),
		OrgID:        orgID,
		RoleIds:      membershipData.Roles,
		Scope:        defaultScope,
		Enabled:      true,
	}

	createdMembership, err := s.deps.MembershipService.CreateMembership(txCtx, membershipCreateDto)
	if err != nil {
		logger.Error(err, "[SERVICE:Onboarding] [TXN] Failed to create membership - transaction will rollback")
		return nil, &customErrors.ServerCustomError{
			Code:   httpStatus.INTERNAL_SERVER_ERROR,
			Errors: []string{fmt.Sprintf("Failed to create membership for organization %s", orgID)},
		}
	}
	createdMemberships = append(createdMemberships, createdMembership)
	logger.Debug(fmt.Sprintf("[SERVICE:Onboarding] [TXN] Direct membership created: %s", *createdMembership.ID))
	return createdMemberships, nil
}

// attachInitialGroup honors the Groups[0] entry of the onboarding payload:
// when ExistingGroup is set, it just adds the user to that group; when
// NewGroup is set, it creates the group first (with its own roles) and
// then adds the user.
func (s *UserOnboardingService) attachInitialGroup(
	txCtx ctx.Context,
	requestContext *reqCtx.RequestContext,
	dto *dtos.CreateUserWithMembershipsDto,
	createdUser *userDtos.UserResponse,
	orgID string,
) error {
	if len(dto.Groups) == 0 {
		return nil
	}
	groupAccessData := dto.Groups[0]

	targetGroupID := ""
	if groupAccessData.ExistingGroup != nil {
		targetGroupID = groupAccessData.ExistingGroup.GroupID
		logger.Debug(fmt.Sprintf("[SERVICE:Onboarding] [TXN] Using existing group: GroupID=%s", targetGroupID))
	}
	if groupAccessData.NewGroup != nil {
		newGroupData := groupAccessData.NewGroup
		logger.Debug(fmt.Sprintf("[SERVICE:Onboarding] [TXN] Creating new group: Name=%s, Roles=%v", newGroupData.Name, newGroupData.RoleIds))

		orgObjectId, err := model.ToObjectID(orgID)
		if err != nil {
			logger.Error(err, "[SERVICE:Onboarding] [TXN] Invalid orgID format - transaction will rollback")
			return &customErrors.ServerCustomError{
				Code:   httpStatus.BAD_REQUEST,
				Errors: []string{fmt.Sprintf("Invalid organization ID format: %s", orgID)},
			}
		}

		groupCreateDto := &groupDtos.CreateGroupDto{
			Name:        newGroupData.Name,
			Description: newGroupData.Description,
			Enabled:     true,
			RoleIds:     newGroupData.RoleIds,
			OrgID:       &orgObjectId,
		}

		createdGroup, err := s.deps.GroupService.CreateGroup(txCtx, requestContext, groupCreateDto)
		if err != nil {
			logger.Error(err, "[SERVICE:Onboarding] [TXN] Failed to create new group - transaction will rollback")
			return &customErrors.ServerCustomError{
				Code:   httpStatus.INTERNAL_SERVER_ERROR,
				Errors: []string{fmt.Sprintf("Failed to create group '%s'", newGroupData.Name)},
			}
		}
		targetGroupID = createdGroup.ID.Hex()
		logger.Debug(fmt.Sprintf("[SERVICE:Onboarding] [TXN] New group created: %s", targetGroupID))
	}

	logger.Debug(fmt.Sprintf("[SERVICE:Onboarding] [TXN] Adding user to group: GroupID=%s", targetGroupID))
	if err := s.deps.GroupService.AddMemberToGroup(txCtx, targetGroupID, createdUser.ID.Hex()); err != nil {
		logger.Error(err, fmt.Sprintf("[SERVICE:Onboarding] [TXN] Failed to add user to group %s - transaction will rollback", targetGroupID))
		return &customErrors.ServerCustomError{
			Code:   httpStatus.INTERNAL_SERVER_ERROR,
			Errors: []string{fmt.Sprintf("Failed to add user to group %s", targetGroupID)},
		}
	}
	logger.Debug("[SERVICE:Onboarding] [TXN] User added to group successfully")
	return nil
}

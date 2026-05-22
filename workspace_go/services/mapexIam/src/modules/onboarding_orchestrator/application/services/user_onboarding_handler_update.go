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
	"github.com/Mapex-Solutions/mapexGoKit/utils/typeconv"
)

// updateOnboardingUser persists the user-field patch via UserService
// inside the caller's transaction context.
func (s *UserOnboardingService) updateOnboardingUser(txCtx ctx.Context, userID string, dto *dtos.UpdateUserWithAccessDto) (*userDtos.UserResponse, error) {
	logger.Info(fmt.Sprintf("[SERVICE:Onboarding] [TXN] Updating user: %s", userID))
	userUpdateDto := &userDtos.UserUpdate{
		FirstName:               dto.FirstName,
		LastName:                dto.LastName,
		Phone:                   dto.Phone,
		JobTitle:                dto.JobTitle,
		Enabled:                 dto.Enabled,
		Avatar:                  dto.Avatar,
		Password:                dto.Password,
		ChangePasswordNextLogin: dto.ChangePasswordNextLogin,
	}
	updatedUser, err := s.deps.UserService.UpdateUserById(txCtx, &userID, userUpdateDto)
	if err != nil {
		logger.Error(err, "[SERVICE:Onboarding] [TXN] Failed to update user - transaction will rollback")
		return nil, err
	}
	logger.Info(fmt.Sprintf("[SERVICE:Onboarding] [TXN] User updated: %s", *updatedUser.ID))
	return updatedUser, nil
}

// replaceDirectMemberships applies the REPLACE strategy on direct
// memberships: delete every existing direct user-org membership for the
// given userId+orgId, then recreate the desired list. Returns the new
// membership responses.
func (s *UserOnboardingService) replaceDirectMemberships(
	txCtx ctx.Context,
	userID, orgID, defaultScope string,
	desired []dtos.MembershipData,
) ([]*membershipContractDtos.MembershipResponse, error) {
	logger.Debug(fmt.Sprintf("[SERVICE:Onboarding] [TXN] Processing direct memberships update (%d items)", len(desired)))

	membershipQuery := &membershipDtos.MembershipQueryDto{
		AssigneeID:   &userID,
		AssigneeType: typeconv.PtrString("user"),
		OrgID:        &orgID,
	}

	existingMemberships, err := s.deps.MembershipService.GetAllMemberships(txCtx, membershipQuery)
	if err != nil {
		logger.Warn(fmt.Sprintf("[SERVICE:Onboarding] [TXN] Failed to query existing memberships: %v", err))
	}

	for _, membership := range existingMemberships {
		membershipID := membership.ID.Hex()
		logger.Debug(fmt.Sprintf("[SERVICE:Onboarding] [TXN] Deleting existing membership: %s", membershipID))
		if _, delErr := s.deps.MembershipService.DeleteMembershipById(txCtx, &membershipID); delErr != nil {
			logger.Error(delErr, fmt.Sprintf("[SERVICE:Onboarding] [TXN] Failed to delete membership %s", membershipID))
			return nil, &customErrors.ServerCustomError{
				Code:   httpStatus.INTERNAL_SERVER_ERROR,
				Errors: []string{fmt.Sprintf("Failed to remove old membership: %s", membershipID)},
			}
		}
	}
	logger.Debug(fmt.Sprintf("[SERVICE:Onboarding] [TXN] Removed %d existing direct memberships", len(existingMemberships)))

	createdMemberships := make([]*membershipContractDtos.MembershipResponse, 0, len(desired))
	for i, membershipData := range desired {
		scope := defaultScope
		if membershipData.Scope != nil && *membershipData.Scope != "" {
			scope = *membershipData.Scope
		}
		logger.Debug(fmt.Sprintf("[SERVICE:Onboarding] [TXN] Creating direct membership %d: OrgID=%s, Roles=%v, Scope=%s",
			i+1, orgID, membershipData.Roles, scope))

		membershipCreateDto := &membershipDtos.CreateMembershipDto{
			AssigneeType: "user",
			AssigneeID:   userID,
			OrgID:        orgID,
			RoleIds:      membershipData.Roles,
			Scope:        scope,
			Enabled:      true,
		}

		createdMembership, err := s.deps.MembershipService.CreateMembership(txCtx, membershipCreateDto)
		if err != nil {
			logger.Error(err, "[SERVICE:Onboarding] [TXN] Failed to create membership")
			return nil, &customErrors.ServerCustomError{
				Code:   httpStatus.INTERNAL_SERVER_ERROR,
				Errors: []string{fmt.Sprintf("Failed to create membership for organization %s", orgID)},
			}
		}
		createdMemberships = append(createdMemberships, createdMembership)
		logger.Debug(fmt.Sprintf("[SERVICE:Onboarding] [TXN] Direct membership created: %s", *createdMembership.ID))
	}
	return createdMemberships, nil
}

// diffApplyGroupMemberships applies the DIFF strategy on group access:
// resolve desired group ids (creating new groups when needed), compute
// add/remove against the user's current groups in this org, then run
// the GroupService add/remove calls. nil pointer = leave unchanged;
// [] = remove from all groups.
func (s *UserOnboardingService) diffApplyGroupMemberships(
	txCtx ctx.Context,
	requestContext *reqCtx.RequestContext,
	userID, orgID string,
	groups []dtos.GroupAccessData,
) error {
	logger.Debug(fmt.Sprintf("[SERVICE:Onboarding] [TXN] Processing group memberships update (%d desired groups)", len(groups)))

	currentGroupIds, err := s.deps.GroupService.GetUserGroupsInOrg(txCtx, userID, orgID)
	if err != nil {
		logger.Warn(fmt.Sprintf("[SERVICE:Onboarding] [TXN] Failed to query user groups: %v", err))
		currentGroupIds = []string{}
	}
	logger.Debug(fmt.Sprintf("[SERVICE:Onboarding] [TXN] User currently in %d groups: %v", len(currentGroupIds), currentGroupIds))

	desiredGroupIds := make(map[string]bool)
	for _, groupAccessData := range groups {
		targetGroupID := ""
		if groupAccessData.ExistingGroup != nil {
			targetGroupID = groupAccessData.ExistingGroup.GroupID
			logger.Debug(fmt.Sprintf("[SERVICE:Onboarding] [TXN] Desired existing group: %s", targetGroupID))
		}
		if groupAccessData.NewGroup != nil {
			newGroupData := groupAccessData.NewGroup
			logger.Debug(fmt.Sprintf("[SERVICE:Onboarding] [TXN] Creating new group: Name=%s, Roles=%v", newGroupData.Name, newGroupData.RoleIds))

			orgObjectId, oErr := model.ToObjectID(orgID)
			if oErr != nil {
				logger.Error(oErr, "[SERVICE:Onboarding] [TXN] Invalid orgID format")
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

			createdGroup, cErr := s.deps.GroupService.CreateGroup(txCtx, requestContext, groupCreateDto)
			if cErr != nil {
				logger.Error(cErr, "[SERVICE:Onboarding] [TXN] Failed to create new group")
				return &customErrors.ServerCustomError{
					Code:   httpStatus.INTERNAL_SERVER_ERROR,
					Errors: []string{fmt.Sprintf("Failed to create group '%s'", newGroupData.Name)},
				}
			}
			targetGroupID = createdGroup.ID.Hex()
			logger.Debug(fmt.Sprintf("[SERVICE:Onboarding] [TXN] New group created: %s", targetGroupID))
		}
		if targetGroupID != "" {
			desiredGroupIds[targetGroupID] = true
		}
	}

	currentGroupSet := make(map[string]bool)
	for _, gid := range currentGroupIds {
		currentGroupSet[gid] = true
	}

	groupsToAdd := []string{}
	groupsToRemove := []string{}
	for gid := range desiredGroupIds {
		if !currentGroupSet[gid] {
			groupsToAdd = append(groupsToAdd, gid)
		}
	}
	for _, gid := range currentGroupIds {
		if !desiredGroupIds[gid] {
			groupsToRemove = append(groupsToRemove, gid)
		}
	}
	logger.Debug(fmt.Sprintf("[SERVICE:Onboarding] [TXN] Group diff: ADD %v, REMOVE %v", groupsToAdd, groupsToRemove))

	for _, groupID := range groupsToRemove {
		logger.Debug(fmt.Sprintf("[SERVICE:Onboarding] [TXN] Removing user from group: %s", groupID))
		if rErr := s.deps.GroupService.RemoveMemberFromGroup(txCtx, groupID, userID); rErr != nil {
			logger.Error(rErr, fmt.Sprintf("[SERVICE:Onboarding] [TXN] Failed to remove user from group %s", groupID))
			return &customErrors.ServerCustomError{
				Code:   httpStatus.INTERNAL_SERVER_ERROR,
				Errors: []string{fmt.Sprintf("Failed to remove user from group: %s", groupID)},
			}
		}
	}
	for _, groupID := range groupsToAdd {
		logger.Debug(fmt.Sprintf("[SERVICE:Onboarding] [TXN] Adding user to group: %s", groupID))
		if aErr := s.deps.GroupService.AddMemberToGroup(txCtx, groupID, userID); aErr != nil {
			logger.Error(aErr, fmt.Sprintf("[SERVICE:Onboarding] [TXN] Failed to add user to group %s", groupID))
			return &customErrors.ServerCustomError{
				Code:   httpStatus.INTERNAL_SERVER_ERROR,
				Errors: []string{fmt.Sprintf("Failed to add user to group: %s", groupID)},
			}
		}
	}
	logger.Debug(fmt.Sprintf("[SERVICE:Onboarding] [TXN] Group membership updated: added %d, removed %d", len(groupsToAdd), len(groupsToRemove)))
	return nil
}

package services

import (
	ctx "context"
	"fmt"

	"mapexIam/src/modules/onboarding_orchestrator/application/dtos"

	membershipContractDtos "github.com/Mapex-Solutions/MapexOS/contracts/services/mapexIam/memberships"
	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// runCreateOnboardingTransaction encapsulates the ACID transaction that
// creates the user, attaches the initial direct membership when requested,
// and assigns the initial group. Any failure rolls everything back.
func (s *UserOnboardingService) runCreateOnboardingTransaction(c ctx.Context, requestContext *reqCtx.RequestContext, dto *dtos.CreateUserWithMembershipsDto, orgID string, defaultScope string) (*dtos.UserOnboardingResponse, error) {
	logger.Info("[SERVICE:Onboarding] Starting ACID transaction for user creation + memberships")
	result, err := s.deps.MongoManager.RunTransaction(c, func(txCtx ctx.Context) (interface{}, error) {
		createdUser, err := s.createOnboardingUser(txCtx, dto)
		if err != nil {
			return nil, err
		}
		createdMemberships, err := s.attachInitialMembership(txCtx, dto, createdUser, orgID, defaultScope)
		if err != nil {
			return nil, err
		}
		if err := s.attachInitialGroup(txCtx, requestContext, dto, createdUser, orgID); err != nil {
			return nil, err
		}
		return &dtos.UserOnboardingResponse{
			User:        createdUser,
			Memberships: createdMemberships,
		}, nil
	})
	if err != nil {
		logger.Error(err, "[SERVICE:Onboarding] Transaction failed - all changes rolled back")
		return nil, err
	}
	return result.(*dtos.UserOnboardingResponse), nil
}

// runUpdateAccessTransaction encapsulates the ACID transaction that updates
// user fields, applies REPLACE on direct memberships when the patch carries
// a Memberships pointer, and applies DIFF on groups when the Groups pointer
// is set. nil pointer = leave as-is; [] = remove all; [items] = desired.
func (s *UserOnboardingService) runUpdateAccessTransaction(c ctx.Context, requestContext *reqCtx.RequestContext, userID string, dto *dtos.UpdateUserWithAccessDto, orgID string, defaultScope string) (*dtos.UserOnboardingResponse, error) {
	logger.Info("[SERVICE:Onboarding] Starting ACID transaction for user update + access management")
	result, err := s.deps.MongoManager.RunTransaction(c, func(txCtx ctx.Context) (interface{}, error) {
		updatedUser, err := s.updateOnboardingUser(txCtx, userID, dto)
		if err != nil {
			return nil, err
		}
		createdMemberships := make([]*membershipContractDtos.MembershipResponse, 0)
		if dto.Memberships != nil {
			memberships, mErr := s.replaceDirectMemberships(txCtx, userID, orgID, defaultScope, *dto.Memberships)
			if mErr != nil {
				return nil, mErr
			}
			createdMemberships = memberships
		}
		if dto.Groups != nil {
			if gErr := s.diffApplyGroupMemberships(txCtx, requestContext, userID, orgID, *dto.Groups); gErr != nil {
				return nil, gErr
			}
		}
		return &dtos.UserOnboardingResponse{
			User:        updatedUser,
			Memberships: createdMemberships,
		}, nil
	})
	if err != nil {
		logger.Error(err, "[SERVICE:Onboarding] Transaction failed - all changes rolled back")
		return nil, err
	}
	return result.(*dtos.UserOnboardingResponse), nil
}

// logOnboardingCreated emits the success log for the create flow with the
// access-type discriminator (group vs direct).
func (s *UserOnboardingService) logOnboardingCreated(response *dtos.UserOnboardingResponse, dto *dtos.CreateUserWithMembershipsDto, orgName string) {
	accessType := "group"
	if len(dto.Memberships) > 0 {
		accessType = "direct"
	}
	logger.Info(fmt.Sprintf("[SERVICE:Onboarding] Transaction committed successfully - User %s created with %s access in org %s",
		*response.User.ID, accessType, orgName))
}

package services

import (
	ctx "context"
	"fmt"

	"mapexIam/src/modules/onboarding_orchestrator/application/di"
	"mapexIam/src/modules/onboarding_orchestrator/application/dtos"
	"mapexIam/src/modules/onboarding_orchestrator/application/ports"

	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// Compile-time check to ensure UserOnboardingService implements UserOnboardingServicePort
var _ ports.UserOnboardingServicePort = (*UserOnboardingService)(nil)

// New creates a new UserOnboardingService with injected dependencies.
func New(deps di.UserOnboardingServiceDependenciesInjection) ports.UserOnboardingServicePort {
	return &UserOnboardingService{deps: deps}
}

// CreateUserWithMemberships orchestrates onboarding within a single ACID
// transaction: resolve org context -> run transactional create -> log result.
// Steps inside the transaction (create user, attach direct membership, attach
// group) live in runCreateOnboardingTransaction so this orchestration stays
// readable.
func (s *UserOnboardingService) CreateUserWithMemberships(c ctx.Context, requestContext *reqCtx.RequestContext, dto *dtos.CreateUserWithMembershipsDto) (*dtos.UserOnboardingResponse, error) {
	logger.Info(fmt.Sprintf("[SERVICE:Onboarding] Starting user onboarding for email: %s", dto.Email))
	orgID, defaultScope, orgName, err := s.resolveOnboardingOrg(c, requestContext)
	if err != nil {
		return nil, err
	}
	response, err := s.runCreateOnboardingTransaction(c, requestContext, dto, orgID, defaultScope)
	if err != nil {
		return nil, err
	}
	s.logOnboardingCreated(response, dto, orgName)
	return response, nil
}

// UpdateUserWithAccess orchestrates a declarative user update inside a single
// ACID transaction. Per-field semantics live in
// runUpdateAccessTransaction (REPLACE direct memberships when the patch
// carries the pointer, DIFF group memberships when the Groups pointer is set,
// nil = leave as-is).
func (s *UserOnboardingService) UpdateUserWithAccess(c ctx.Context, requestContext *reqCtx.RequestContext, userID string, dto *dtos.UpdateUserWithAccessDto) (*dtos.UserOnboardingResponse, error) {
	logger.Info(fmt.Sprintf("[SERVICE:Onboarding] Starting user update with access for userId: %s", userID))
	orgID, defaultScope, orgName, err := s.resolveOnboardingOrg(c, requestContext)
	if err != nil {
		return nil, err
	}
	response, err := s.runUpdateAccessTransaction(c, requestContext, userID, dto, orgID, defaultScope)
	if err != nil {
		return nil, err
	}
	logger.Info(fmt.Sprintf("[SERVICE:Onboarding] Transaction committed successfully - User %s updated in org %s",
		*response.User.ID, orgName))
	return response, nil
}

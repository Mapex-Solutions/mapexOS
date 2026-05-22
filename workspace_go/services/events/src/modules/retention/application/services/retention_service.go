package services

import (
	ctx "context"
	"encoding/json"
	"fmt"

	"events/src/modules/retention/application/constants"
	"events/src/modules/retention/application/di"
	"events/src/modules/retention/application/dtos"
	"events/src/modules/retention/application/ports"
	"events/src/modules/retention/domain/entities"

	retentionContracts "github.com/Mapex-Solutions/MapexOS/contracts/services/events/retention"
	orgDtos "github.com/Mapex-Solutions/MapexOS/contracts/services/mapexIam/organizations"
	common "github.com/Mapex-Solutions/mapexGoKit/infrastructure/common/ports"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	customErrors "github.com/Mapex-Solutions/mapexGoKit/microservices/http/customErrors"
	httpStatus "github.com/Mapex-Solutions/mapexGoKit/microservices/http/status"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
	"github.com/Mapex-Solutions/mapexGoKit/utils/orgfilter"
)

// Compile-time check to ensure RetentionService implements RetentionServicePort interface.
var _ ports.RetentionServicePort = (*RetentionService)(nil)

// New creates and returns a new instance of RetentionService.
func New(deps di.RetentionServiceDependenciesInjection) ports.RetentionServicePort {
	return &RetentionService{
		deps: deps,
	}
}

// GetRetentionPolicies orchestrates the paginated list:
// build the org filter -> apply optional type/$in filter -> resolve
// pagination + projection -> delegate to repository -> map entities to
// response DTOs.
func (s *RetentionService) GetRetentionPolicies(
	c ctx.Context,
	requestContext *reqCtx.RequestContext,
	query *dtos.RetentionPolicyQueryDTO,
) (*model.PaginatedResult[dtos.RetentionPolicyResponse], error) {
	filters := s.buildPolicyListFilters(requestContext, query)
	pagination := &model.PaginationOpts{Page: int64(query.GetPage()), PerPage: int64(query.GetPerPage())}
	projection := orgfilter.BuildProjection(query.Projection)

	result, err := s.deps.RetentionRepo.FindWithFilters(c, filters, pagination, projection)
	if err != nil {
		return nil, err
	}

	dtoItems := make([]dtos.RetentionPolicyResponse, len(result.Items))
	for i, entity := range result.Items {
		dtoItems[i] = *s.entityToResponse(&entity)
	}
	return &model.PaginatedResult[dtos.RetentionPolicyResponse]{
		Items:      dtoItems,
		Pagination: result.Pagination,
	}, nil
}

// GetRetentionPolicyById fetches a single policy by id. Returns 404 when
// the id is unknown.
func (s *RetentionService) GetRetentionPolicyById(c ctx.Context, policyId *string) (*dtos.RetentionPolicyResponse, error) {
	policy, err := s.deps.RetentionRepo.FindById(c, policyId)
	if err != nil || policy == nil || policy.ID.IsZero() {
		return nil, &customErrors.ServerCustomError{Code: httpStatus.NOT_FOUND, Errors: []string{"retention policy not found"}}
	}
	return s.entityToResponse(policy), nil
}

// UpsertRetentionPolicy orchestrates create-or-update by orgId+type:
// validate org context + retention limits -> build the entity from the
// patch -> upsert -> drop the cache key for the pair -> apply the
// ClickHouse TTL when the policy governs a TTL-bound table.
func (s *RetentionService) UpsertRetentionPolicy(
	c ctx.Context,
	requestContext *reqCtx.RequestContext,
	dto *dtos.RetentionPolicyUpsertDTO,
) (*dtos.RetentionPolicyResponse, error) {
	if err := orgfilter.ValidateOrgContextForNonSystem(requestContext); err != nil {
		return nil, &customErrors.ServerCustomError{Code: httpStatus.BAD_REQUEST, Errors: []string{err.Error()}}
	}
	if err := retentionContracts.ValidateRetentionPolicy(dto.Type, dto.RetentionDays); err != nil {
		return nil, &customErrors.ServerCustomError{Code: httpStatus.BAD_REQUEST, Errors: []string{err.Error()}}
	}

	orgObjectId, pathKey := extractOrgScope(requestContext)
	policyEntity := buildPolicyEntity(dto, pathKey)

	upserted, err := s.deps.RetentionRepo.Upsert(c, &orgObjectId, dto.Type, policyEntity)
	if err != nil {
		return nil, err
	}

	s.invalidatePolicyCache(c, orgObjectId.Hex(), dto.Type)
	s.applyTTLOnUpsert(c, dto.Type, dto.RetentionDays)
	return s.entityToResponse(upserted), nil
}

// DeleteRetentionPolicyById orchestrates deletion: load the policy (404
// on miss) -> delete in repository -> drop the cache key for its
// orgId+type pair.
func (s *RetentionService) DeleteRetentionPolicyById(c ctx.Context, policyId *string) (map[string]bool, error) {
	policy, err := s.deps.RetentionRepo.FindById(c, policyId)
	if err != nil || policy == nil {
		return nil, &customErrors.ServerCustomError{Code: httpStatus.NOT_FOUND, Errors: []string{"retention policy not found"}}
	}

	if err := s.deps.RetentionRepo.DeleteById(c, policyId); err != nil {
		return nil, err
	}

	if policy.OrgId != nil {
		s.invalidatePolicyCache(c, policy.OrgId.Hex(), policy.Type)
	}
	return map[string]bool{"success": true}, nil
}

// GetRetentionDays returns the cached retention days for one org+table
// pair, falling back to the contract default when no policy exists.
// Cache-aside: the GetOrSetEx callback hits Mongo on miss and rewrites
// the cache with the standard TTL.
func (s *RetentionService) GetRetentionDays(c ctx.Context, orgId string, tableName string) (uint16, error) {
	cacheKey := buildRetentionCacheKey(orgId, tableName)

	var policy entities.RetentionPolicy
	_, err := s.deps.CacheRepo.GetOrSetEx(common.GetOrSetParams{
		Ctx:      c,
		CacheKey: cacheKey,
		CacheTTL: int(constants.RetentionCacheTTL.Seconds()),
		Dest:     &policy,
		Callback: func() (interface{}, error) {
			return s.deps.RetentionRepo.FindByOrgIdAndType(c, &orgId, tableName)
		},
	})

	if err != nil || policy.ID.IsZero() {
		return defaultRetentionFor(orgId, tableName), nil
	}
	return policy.RetentionDays, nil
}

// CreateDefaultPolicies seeds the per-table defaults for a freshly
// created organization. Per-table failures are logged but do not abort
// the seeding loop — a partial result is still better than none.
func (s *RetentionService) CreateDefaultPolicies(c ctx.Context, orgId string, pathKey string) error {
	orgObjectId, err := model.ToObjectID(orgId)
	if err != nil {
		return fmt.Errorf("invalid orgId: %s", orgId)
	}

	for retentionType, limits := range retentionContracts.RetentionPoliciesLimits {
		policyEntity := &entities.RetentionPolicy{
			Name:          limits.Name,
			Type:          retentionType,
			RetentionDays: limits.DefaultDays,
			PathKey:       pathKey,
			Enabled:       true,
		}
		if _, upErr := s.deps.RetentionRepo.Upsert(c, &orgObjectId, retentionType, policyEntity); upErr != nil {
			logger.Error(upErr, fmt.Sprintf("[SERVICE:Retention] Failed to create default policy for org=%s type=%s", orgId, retentionType))
			continue
		}
	}
	logger.Info(fmt.Sprintf("[SERVICE:Retention] Created %d default policies for org=%s", len(retentionContracts.RetentionPoliciesLimits), orgId))
	return nil
}

// SeedPlatformPolicies upserts platform-scoped (no orgId) policies and
// applies the matching ClickHouse TTL. Idempotent.
func (s *RetentionService) SeedPlatformPolicies(c ctx.Context) error {
	policyEntity := &entities.RetentionPolicy{
		Name:          constants.AssetStatusHistoryPolicyName,
		Type:          constants.TableAssetStatusHistory,
		RetentionDays: constants.AssetStatusHistoryDefaultDays,
		Enabled:       true,
	}
	if _, err := s.deps.RetentionRepo.Upsert(c, nil, constants.TableAssetStatusHistory, policyEntity); err != nil {
		logger.Error(err, fmt.Sprintf("[SERVICE:Retention] Failed to seed platform policy for %s", constants.TableAssetStatusHistory))
		return err
	}
	if ttlErr := s.ApplyAssetStatusHistoryTTL(c, constants.AssetStatusHistoryDefaultDays); ttlErr != nil {
		logger.Warn(fmt.Sprintf("[SERVICE:Retention] Seed succeeded but failed to apply ClickHouse TTL: %v", ttlErr))
	}
	logger.Info(fmt.Sprintf("[SERVICE:Retention] Seeded platform policy: %s (%d days)", constants.TableAssetStatusHistory, constants.AssetStatusHistoryDefaultDays))
	return nil
}

// ApplyAssetStatusHistoryTTL issues an ALTER TABLE asset_status_history
// MODIFY TTL statement against ClickHouse. Validates the requested days
// against the contract bounds before issuing the DDL.
func (s *RetentionService) ApplyAssetStatusHistoryTTL(c ctx.Context, days uint16) error {
	if s.deps.ClickHouseConn == nil {
		return fmt.Errorf("ClickHouseConn not injected — cannot apply TTL")
	}
	if days < constants.AssetStatusHistoryMinDays || days > constants.AssetStatusHistoryMaxDays {
		return fmt.Errorf("TTL days %d out of range [%d, %d]", days, constants.AssetStatusHistoryMinDays, constants.AssetStatusHistoryMaxDays)
	}

	stmt := fmt.Sprintf("ALTER TABLE %s MODIFY TTL created + toIntervalDay(%d)", constants.TableAssetStatusHistory, days)
	if err := s.deps.ClickHouseConn.Exec(c, stmt); err != nil {
		logger.Error(err, fmt.Sprintf("[SERVICE:Retention] TTL apply failed: table=%s days=%d", constants.TableAssetStatusHistory, days))
		return err
	}
	logger.Info(fmt.Sprintf("[SERVICE:Retention] TTL applied: table=%s days=%d", constants.TableAssetStatusHistory, days))
	return nil
}

// HandleOrgCreatedEvent decodes the NATS payload and triggers default
// policy creation for the new organization. Returns nil on success or
// unrecoverable parse failure (caller acks); a non-nil error nacks the
// message for retry.
func (s *RetentionService) HandleOrgCreatedEvent(msg *natsModel.Message) error {
	var event orgDtos.OrganizationCreatedEvent
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		logger.Error(err, "[SERVICE:Retention] Failed to unmarshal organization.created event")
		return nil
	}

	if event.OrgId == "" || event.PathKey == "" {
		logger.Warn("[SERVICE:Retention] Skipping organization.created event with empty orgId or pathKey")
		return nil
	}

	if err := s.CreateDefaultPolicies(ctx.Background(), event.OrgId, event.PathKey); err != nil {
		logger.Error(err, fmt.Sprintf("[SERVICE:Retention] Failed to create default policies for org=%s", event.OrgId))
		return err
	}
	logger.Info(fmt.Sprintf("[SERVICE:Retention] Created default policies for org=%s", event.OrgId))
	return nil
}


package services

import (
	ctx "context"
	"fmt"
	"strings"

	"events/src/modules/retention/application/constants"
	"events/src/modules/retention/application/dtos"
	"events/src/modules/retention/domain/entities"

	retentionContracts "github.com/Mapex-Solutions/MapexOS/contracts/services/events/retention"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
	"github.com/Mapex-Solutions/mapexGoKit/utils/orgfilter"
)

// buildRetentionCacheKey creates a cache key for a retention policy.
// Format: RETENTION_POLICY:{orgId}:{retentionType}.
func buildRetentionCacheKey(orgId string, retentionType string) string {
	return fmt.Sprintf("%s:%s:%s", constants.CacheKeyPrefix, orgId, retentionType)
}

// entityToResponse converts a RetentionPolicy entity to a RetentionPolicyResponse DTO.
func (s *RetentionService) entityToResponse(entity *entities.RetentionPolicy) *dtos.RetentionPolicyResponse {
	idStr := entity.ID.Hex()
	nameStr := entity.Name
	typeStr := entity.Type
	retentionDays := entity.RetentionDays
	pathKeyStr := entity.PathKey
	enabled := entity.Enabled
	created := entity.Created.Format("2006-01-02T15:04:05Z")
	updated := entity.Updated.Format("2006-01-02T15:04:05Z")

	resp := &dtos.RetentionPolicyResponse{
		ID:            &idStr,
		Name:          &nameStr,
		Type:          &typeStr,
		RetentionDays: &retentionDays,
		PathKey:       &pathKeyStr,
		Enabled:       &enabled,
		Created:       &created,
		Updated:       &updated,
	}

	if entity.OrgId != nil {
		orgIdStr := entity.OrgId.Hex()
		resp.OrgId = &orgIdStr
	}

	return resp
}

// buildPolicyListFilters wraps the list-mode org filter with the optional
// type / $in filter from the query DTO.
func (s *RetentionService) buildPolicyListFilters(rc *reqCtx.RequestContext, query *dtos.RetentionPolicyQueryDTO) model.Map {
	filters, _ := orgfilter.BuildOrgFilter(orgfilter.BuildFilterParams{ReqContext: rc, Query: query})
	if query.Type != nil && *query.Type != "" {
		types := strings.Split(*query.Type, ",")
		if len(types) == 1 {
			filters["type"] = types[0]
		} else {
			filters["type"] = model.Map{"$in": types}
		}
	}
	return filters
}

// invalidatePolicyCache drops the cache key for one org+table pair so the
// next read repopulates from Mongo.
func (s *RetentionService) invalidatePolicyCache(c ctx.Context, orgIdHex, tableName string) {
	cacheKey := buildRetentionCacheKey(orgIdHex, tableName)
	s.deps.CacheRepo.Del(c, cacheKey)
}

// applyTTLOnUpsert applies the ClickHouse TTL only for tables that currently
// participate in TTL bindings (today, only asset_status_history). Failures
// are logged — operators retry via a fresh upsert.
func (s *RetentionService) applyTTLOnUpsert(c ctx.Context, tableName string, days uint16) {
	if tableName != constants.TableAssetStatusHistory {
		return
	}
	if err := s.ApplyAssetStatusHistoryTTL(c, days); err != nil {
		logger.Warn(fmt.Sprintf("[SERVICE:Retention] Upsert committed but TTL apply failed: %v", err))
	}
}

// extractOrgScope pulls orgObjectId + pathKey from the request context for
// upsert. Empty values are normalized to the zero ObjectId / empty string
// (used by the seed path).
func extractOrgScope(rc *reqCtx.RequestContext) (model.ObjectId, string) {
	var orgObjectId model.ObjectId
	pathKey := ""
	if rc.OrgContext != nil && *rc.OrgContext != "" {
		orgObjectId, _ = model.ToObjectID(*rc.OrgContext)
	}
	if rc.OrgContextData != nil && rc.OrgContextData.PathKey != "" {
		pathKey = rc.OrgContextData.PathKey
	}
	return orgObjectId, pathKey
}

// buildPolicyEntity assembles the entity passed to Upsert. Defaults Enabled
// to true when the patch did not specify it.
func buildPolicyEntity(dto *dtos.RetentionPolicyUpsertDTO, pathKey string) *entities.RetentionPolicy {
	enabled := true
	if dto.Enabled != nil {
		enabled = *dto.Enabled
	}
	return &entities.RetentionPolicy{
		Name:          dto.Name,
		Type:          dto.Type,
		RetentionDays: dto.RetentionDays,
		PathKey:       pathKey,
		Enabled:       enabled,
	}
}

// defaultRetentionFor returns the contract-defined default for a table when
// no policy exists for the org, falling back further to the global constant
// when even the contract has no entry.
func defaultRetentionFor(orgId, tableName string) uint16 {
	if limits, exists := retentionContracts.RetentionPoliciesLimits[tableName]; exists {
		return limits.DefaultDays
	}
	logger.Warn(fmt.Sprintf("[SERVICE:Retention] No policy found for org %s table %s, using fallback default (%d days)", orgId, tableName, constants.DefaultRetentionDays))
	return constants.DefaultRetentionDays
}

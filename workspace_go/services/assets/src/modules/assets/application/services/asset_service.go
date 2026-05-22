package services

import (
	ctx "context"
	"time"

	"assets/src/modules/assets/application/di"
	"assets/src/modules/assets/application/dtos"
	"assets/src/modules/assets/application/ports"

	assetsContract "github.com/Mapex-Solutions/MapexOS/contracts/services/assets/assets"
	assetsAuthContract "github.com/Mapex-Solutions/MapexOS/contracts/services/assets/auth"
	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	customErrors "github.com/Mapex-Solutions/mapexGoKit/microservices/http/customErrors"
	httpStatus "github.com/Mapex-Solutions/mapexGoKit/microservices/http/status"
)

// Compile-time check to ensure AssetService implements AssetServicePort interface.
var _ ports.AssetServicePort = (*AssetService)(nil)

// New creates and returns a new instance of AssetService.
//
// This constructor follows Hexagonal Architecture by:
//   - Accepting dependencies through a DI struct (single parameter pattern)
//   - Returning the service port interface (not concrete type)
//   - Enabling loose coupling and testability
//
// Parameters:
//   - deps: Aggregated dependencies (repositories, NATS bus) injected by dig
//
// Returns:
//   - AssetServicePort: The service port interface implementation
func New(deps di.AssetServiceDependenciesInjection) ports.AssetServicePort {
	return &AssetService{
		deps: deps,
	}
}

// CreateAsset orchestrates asset creation:
// bind org context -> validate health-monitor invariants -> build entity
// from DTO -> bcrypt the MQTT password (mqtt-protocol only) -> persist
// once -> fan out side effects (auth cache, MinIO read model, counter
// cache) -> build the response DTO. All metric emission is centralized
// in recordAssetOp so every exit path stays observability-consistent.
func (s *AssetService) CreateAsset(c ctx.Context, requestContext *reqCtx.RequestContext, dto *dtos.AssetCreateDTO) (*dtos.AssetResponse, error) {
	start := time.Now()

	s.bindOrgContextOnCreate(requestContext, dto)
	if err := validateHealthMonitorConfig(c, s.deps.RouteGroupPort, dto.HealthMonitor); err != nil {
		s.recordAssetOp("create", "error", start)
		return nil, err
	}
	entity := s.buildEntityFromDto(dto)
	if err := s.hashMqttPasswordIfNeeded(entity, dto); err != nil {
		s.recordAssetOp("create", "error", start)
		return nil, err
	}
	persisted, err := s.persistNewAsset(c, entity)
	if err != nil {
		s.recordAssetOp("create", "error", start)
		return nil, err
	}
	s.fanoutCreateSideEffects(c, requestContext, persisted)
	s.recordAssetOp("create", "success", start)
	return s.buildCreateResponse(persisted), nil
}

// GenerateMqttPassword returns a fresh random alphanumeric password
// for the operator to drop into the asset's MQTT config before submit.
// Stateless — does not touch any asset record.
func (s *AssetService) GenerateMqttPassword(_ ctx.Context) (*dtos.GenerateMqttPasswordResponseDTO, error) {
	pwd, err := generateAlphanumericPassword(24)
	if err != nil {
		return nil, err
	}
	return &dtos.GenerateMqttPasswordResponseDTO{Password: pwd}, nil
}

// GetAssetById orchestrates a single-asset read:
// fetch from Mongo -> shape the response DTO -> enrich with template
// classification + route-group names + Mongo health flip -> overlay
// real-time Redis health state. Returns 404 when the asset does not exist.
func (s *AssetService) GetAssetById(c ctx.Context, assetId *string) (*dtos.AssetResponse, error) {
	start := time.Now()

	entity, err := s.deps.AssetRepo.FindById(c, assetId)
	if err != nil || entity == nil || entity.ID.IsZero() {
		s.recordAssetOp("read", "error", start)
		return nil, &customErrors.ServerCustomError{Code: httpStatus.NOT_FOUND, Errors: []string{"Asset not found"}}
	}

	resp := s.buildGetByIdResponse(entity)
	s.enrichWithTemplateClassification(c, entity, resp)
	s.enrichWithRouteGroupNames(c, entity, resp)
	s.enrichHealthStatus(c, resp, entity.OrgID.Hex())

	s.recordAssetOp("read", "success", start)
	return resp, nil
}

// GetByUUID returns the asset entity for the given device UUID. Thin
// delegation to the repository's UUID-indexed query; the entity-level
// shape exposes the JWT metadata (Protocol.Mqtt.TokenJti) that the
// devices module needs for jti-match validation on refresh. Returns
// NOT_FOUND when the UUID is unknown.
func (s *AssetService) GetByUUID(c ctx.Context, assetUUID string) (*ports.Asset, error) {
	asset, err := s.deps.AssetRepo.FindByAssetUUID(c, &assetUUID)
	if err != nil || asset == nil || asset.ID.IsZero() {
		return nil, &customErrors.ServerCustomError{
			Code:   httpStatus.NOT_FOUND,
			Errors: []string{"Asset not found for UUID: " + assetUUID},
		}
	}
	return asset, nil
}

// UpdateAssetById orchestrates a partial update:
// load the prior entity -> validate health-monitor invariants -> apply the
// patch in Mongo -> fan out side effects (auth cache, MinIO write, FANOUT
// invalidation, health-state cleanup on enable->disable) -> build the
// response DTO. Returns 404 when the target id is unknown.
func (s *AssetService) UpdateAssetById(c ctx.Context, assetId *string, dto *dtos.AssetUpdateDTO) (*dtos.AssetResponse, error) {
	start := time.Now()

	before, err := s.deps.AssetRepo.FindById(c, assetId)
	if err != nil || before == nil || before.ID.IsZero() {
		s.recordAssetOp("update", "error", start)
		return nil, &customErrors.ServerCustomError{Code: httpStatus.NOT_FOUND, Errors: []string{"Asset not found"}}
	}
	if err := validateHealthMonitorConfig(c, s.deps.RouteGroupPort, dto.HealthMonitor); err != nil {
		s.recordAssetOp("update", "error", start)
		return nil, err
	}
	after, err := s.applyAssetPatch(c, assetId, dto)
	if err != nil {
		s.recordAssetOp("update", "error", start)
		return nil, err
	}
	s.fanoutUpdateSideEffects(c, before, after)

	s.recordAssetOp("update", "success", start)
	return s.buildUpdateResponse(after), nil
}

// DeleteAssetById orchestrates asset deletion in cache-first order:
// load the asset -> clear all caches and Redis health state -> publish
// FANOUT invalidation -> finally delete from Mongo. The cache-first order
// guarantees that any retry after a partial failure repopulates from
// Mongo (still present) instead of leaving orphan cache entries.
func (s *AssetService) DeleteAssetById(c ctx.Context, assetId *string) (map[string]bool, error) {
	start := time.Now()

	if asset, _ := s.deps.AssetRepo.FindById(c, assetId); asset != nil {
		s.tearDownAssetCaches(c, asset)
	}
	if err := s.deps.AssetRepo.DeleteById(c, assetId); err != nil {
		s.recordAssetOp("delete", "error", start)
		return nil, err
	}

	s.recordAssetOp("delete", "success", start)
	return map[string]bool{"success": true}, nil
}

// GetAssets orchestrates the paginated list:
// build the org filter from RequestContext -> apply asset+template
// filters -> run the optimized aggregation -> map entities to DTOs ->
// enrich health status in batch when a single org is in scope.
func (s *AssetService) GetAssets(c ctx.Context, requestContext *reqCtx.RequestContext, query *dtos.AssetQueryDTO) (*model.PaginatedResult[dtos.AssetResponse], error) {
	start := time.Now()

	assetFilters, templateFilters, err := s.buildListFilters(requestContext, query)
	if err != nil {
		s.recordAssetOp("list", "error", start)
		return nil, err
	}

	result, err := s.deps.AssetRepo.FindWithFiltersAndTemplate(c, assetFilters, templateFilters,
		&model.PaginationOpts{Page: int64(query.GetPage()), PerPage: int64(query.GetPerPage())},
		model.Map{"created": -1},
	)
	if err != nil {
		s.recordAssetOp("list", "error", start)
		return nil, err
	}

	dtoItems := s.mapListEntitiesToDtos(result.Items)
	if requestContext.OrgContext != nil && *requestContext.OrgContext != "" {
		s.enrichHealthStatusBatch(c, dtoItems, *requestContext.OrgContext)
	}

	s.recordAssetOp("list", "success", start)
	s.deps.Metrics.AssetListResultsCount.Observe(float64(len(dtoItems)))
	return &model.PaginatedResult[dtos.AssetResponse]{Items: dtoItems, Pagination: result.Pagination}, nil
}

// GetAssetByMqttUsername fetches an asset by its MQTT username. Used by
// the auth module on cache miss; returns 404 when the username is unknown.
func (s *AssetService) GetAssetByMqttUsername(c ctx.Context, username string) (*dtos.AssetResponse, error) {
	asset, err := s.deps.AssetRepo.FindByMqttUsername(c, username)
	if err != nil || asset == nil || asset.ID.IsZero() {
		return nil, &customErrors.ServerCustomError{
			Code:   httpStatus.NOT_FOUND,
			Errors: []string{"Asset not found for username: " + username},
		}
	}
	return s.buildSimpleResponse(asset), nil
}


// GetAssetReadModelByUUID returns the denormalized read model used by
// TieredCache fallback. Repopulates MinIO (L2) before returning so the
// next request from any consumer hits the cache.
func (s *AssetService) GetAssetReadModelByUUID(c ctx.Context, assetUUID string) (*assetsContract.AssetReadModel, error) {
	asset, err := s.deps.AssetRepo.FindByAssetUUID(c, &assetUUID)
	if err != nil || asset == nil || asset.ID.IsZero() {
		return nil, &customErrors.ServerCustomError{
			Code:   httpStatus.NOT_FOUND,
			Errors: []string{"Asset not found for UUID: " + assetUUID},
		}
	}
	templateOrgId := s.writeAssetMetadata(c, asset)
	return s.buildReadModel(asset, templateOrgId), nil
}

// GetAuthProjectionByUUID returns the slim auth projection consumed
// by the broker plugin. Spawns an async goroutine to write the
// projection back to MinIO (mapex-asset-auth/{assetUUID}.json) — the
// HTTP response does NOT wait for the warm-up. Uses a fresh
// background context so the goroutine survives the request lifecycle.
func (s *AssetService) GetAuthProjectionByUUID(c ctx.Context, assetUUID string) (*assetsAuthContract.AuthProjection, error) {
	asset, err := s.deps.AssetRepo.FindByAssetUUID(c, &assetUUID)
	if err != nil || asset == nil || asset.ID.IsZero() {
		return nil, &customErrors.ServerCustomError{
			Code:   httpStatus.NOT_FOUND,
			Errors: []string{"Asset not found for UUID: " + assetUUID},
		}
	}
	projection := buildAuthProjection(asset)
	go func() {
		_ = s.deps.AssetStoragePort.WriteAssetAuth(ctx.Background(), projection)
	}()
	return &projection, nil
}

// ProcessL2WriteRetry is the public entry point for the L2 sync
// fallback consumer (asset_l2sync). On receipt of a retry hint the
// service re-fetches the current asset state from Mongo (NEVER
// trusting the event payload — a stale event can't overwrite newer
// data) and re-runs syncAssetL2. On success the existing FANOUT
// invalidation is emitted so caches downstream refresh. Returns an
// error if Mongo lookup fails — the consumer NAKs and NATS retries
// with backoff.
func (s *AssetService) ProcessL2WriteRetry(c ctx.Context, assetId string) error {
	asset, err := s.fetchAssetByID(c, assetId)
	if err != nil {
		return err
	}
	if asset == nil {
		// Asset deleted between failure and retry — ACK and drop.
		return nil
	}
	s.syncAssetL2(c, asset)
	s.publishAssetInvalidate(c, asset)
	return nil
}

// SetCurrentCert orchestrates the cross-module cert-issued reflection:
// load the asset by UUID (=before) -> apply the cert subdoc via
// FindByIdAndUpdate -> reuse fanoutUpdateSideEffects so the L2 +
// FANOUT cache layers see the new serial. Errors surface NOT_FOUND
// when the asset has been deleted between issue and reflection, and
// upstream errors otherwise.
func (s *AssetService) SetCurrentCert(c ctx.Context, assetUUID string, cert ports.AssetCertificateInput) error {
	before, err := s.loadAssetForCertSync(c, assetUUID)
	if err != nil {
		return err
	}
	after, err := s.applyCurrentCertPatch(c, before, cert)
	if err != nil {
		return err
	}
	s.fanoutUpdateSideEffects(c, before, after)
	return nil
}

// ClearCurrentCertBySerial orchestrates the cross-module revoke
// reflection: locate the asset by `currentCert.serial` -> apply a
// $set/$unset that drops the subdoc -> run fanoutUpdateSideEffects.
// Returns the cleared asset's UUID so the audit-row writer in
// mqttcerts can stamp the asset link.
func (s *AssetService) ClearCurrentCertBySerial(c ctx.Context, serial string) (string, error) {
	before, err := s.findAssetByCertSerial(c, serial)
	if err != nil {
		return "", err
	}
	after, err := s.applyClearCurrentCert(c, before)
	if err != nil {
		return "", err
	}
	s.fanoutUpdateSideEffects(c, before, after)
	return before.AssetUUID, nil
}

// CountAssets returns the total count of assets for the caller's org
// context with cache-aside semantics: try the Redis counter first, fall
// back to Mongo CountDocuments and re-cache on miss.
func (s *AssetService) CountAssets(c ctx.Context, requestContext *reqCtx.RequestContext) (int64, error) {
	start := time.Now()

	orgId := ""
	if requestContext.OrgContext != nil {
		orgId = *requestContext.OrgContext
	}
	cacheKey := s.deps.CacheKeyBuilder.BuildCounterCacheKey(orgId)

	if count, ok := s.tryCachedCount(c, cacheKey); ok {
		s.recordAssetOp("count", "success", start)
		return count, nil
	}

	count, err := s.countFromRepo(c, requestContext)
	if err != nil {
		s.recordAssetOp("count", "error", start)
		return 0, err
	}
	s.cacheCount(c, cacheKey, count)
	s.recordAssetOp("count", "success", start)
	return count, nil
}

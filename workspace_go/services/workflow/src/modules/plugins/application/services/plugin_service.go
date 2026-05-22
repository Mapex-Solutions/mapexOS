package services

import (
	"context"
	"time"

	"workflow/src/modules/plugins/application/di"
	"workflow/src/modules/plugins/application/dtos"
	"workflow/src/modules/plugins/application/ports"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/customErrors"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/status"
)

// Compile-time check
var _ ports.PluginServicePort = (*PluginService)(nil)

// New creates and returns a new PluginService.
func New(deps di.PluginServiceDependenciesInjection) ports.PluginServicePort {
	return &PluginService{deps: deps}
}

// CreatePlugin persists a new plugin manifest, stamps it with the org
// context, then warms the local cache and broadcasts a FANOUT invalidation
// so peer pods drop any stale view of the prior pluginId.
func (s *PluginService) CreatePlugin(ctx context.Context, requestContext *reqCtx.RequestContext, entity *dtos.PluginManifestResponse) (*dtos.PluginManifestResponse, error) {
	start := time.Now()
	s.applyOrgContextToPlugin(entity, requestContext)
	s.applyPluginCreateDefaults(entity)

	created, err := s.deps.PluginRepo.Create(ctx, entity)
	if err != nil {
		s.trackPluginMetrics("create", "error", start)
		return nil, err
	}
	s.invalidateAndFanoutPlugin(ctx, created.PluginID, "create")
	s.trackPluginMetrics("create", "success", start)
	return created, nil
}

// UpdatePluginById applies a partial update against the existing manifest
// and broadcasts a FANOUT invalidation. The pre-fetch is required because
// the cache key is the manifest's `pluginId` (not Mongo `_id`), and we need
// it to invalidate the right entry even if the DTO does not carry it.
func (s *PluginService) UpdatePluginById(ctx context.Context, pluginId *string, dto *dtos.PluginManifestUpdate) (*dtos.PluginManifestResponse, error) {
	start := time.Now()
	existing, err := s.fetchPluginForUpdate(ctx, pluginId, start)
	if err != nil {
		return nil, err
	}
	fieldsToUpdate := s.buildPluginUpdatePayload(dto)
	updated, err := s.persistPluginUpdate(ctx, pluginId, fieldsToUpdate, start)
	if err != nil {
		return nil, err
	}
	s.invalidateAndFanoutPlugin(ctx, existing.PluginID, "update")
	s.trackPluginMetrics("update", "success", start)
	return updated, nil
}

// DeletePluginById removes the manifest, then invalidates the local cache
// and publishes a FANOUT message. The pre-delete fetch is best-effort: a
// fetch failure logs a warning but does not block the delete itself.
func (s *PluginService) DeletePluginById(ctx context.Context, pluginId *string) (map[string]bool, error) {
	start := time.Now()
	existing := s.fetchPluginForDelete(ctx, pluginId)

	if err := s.deps.PluginRepo.DeleteById(ctx, pluginId); err != nil {
		s.trackPluginMetrics("delete", "error", start)
		return nil, err
	}
	s.invalidateAfterPluginDelete(ctx, existing)
	s.trackPluginMetrics("delete", "success", start)
	return map[string]bool{"success": true}, nil
}

// GetPluginById fetches a manifest by its Mongo ObjectId. A missing or
// zero-id row surfaces as the canonical 404 envelope.
func (s *PluginService) GetPluginById(ctx context.Context, pluginId *string) (*dtos.PluginManifestResponse, error) {
	start := time.Now()
	plugin, err := s.deps.PluginRepo.FindById(ctx, pluginId)
	if err != nil || plugin == nil || plugin.ID.IsZero() {
		s.trackPluginMetrics("read", "error", start)
		return nil, notFoundPluginManifest()
	}
	s.trackPluginMetrics("read", "success", start)
	return plugin, nil
}

// GetPluginByPluginId reads through the TieredCache (L0→L1→Mongo) keyed by
// the manifest's external pluginId string. Loader errors and a nil result
// both surface as 404 — callers cannot distinguish "not loaded" from
// "deleted" and should not.
func (s *PluginService) GetPluginByPluginId(ctx context.Context, pluginId string) (*dtos.PluginManifestResponse, error) {
	start := time.Now()
	plugin, err := s.deps.PluginLoader.GetManifest(ctx, pluginId)
	if err != nil {
		s.trackPluginMetrics("read", "error", start)
		return nil, err
	}
	if plugin == nil {
		s.trackPluginMetrics("read", "error", start)
		return nil, notFoundPluginManifest()
	}
	s.trackPluginMetrics("read", "success", start)
	return plugin, nil
}

// GetPlugins returns a paginated list across the multi-tenant visibility
// rules: org-local manifests, plus template manifests from ancestor orgs
// when the query opts in. Filters and pagination assembly live in helpers
// so the orchestration stays terse.
func (s *PluginService) GetPlugins(ctx context.Context, requestContext *reqCtx.RequestContext, query *dtos.PluginQueryDTO) (*model.PaginatedResult[dtos.PluginManifestResponse], error) {
	start := time.Now()
	filters := s.buildPluginListFilters(requestContext, query)
	pagination := s.buildPluginPagination(query)

	result, err := s.deps.PluginRepo.FindWithFilters(ctx, filters, pagination, nil)
	if err != nil {
		s.trackPluginMetrics("list", "error", start)
		return nil, err
	}
	s.trackPluginMetrics("list", "success", start)
	s.deps.Metrics.PluginListResultsCount.Observe(float64(len(result.Items)))
	return result, nil
}

// GetEnabledPlugins returns the editor's boot-time manifest list. Loader
// errors are observed via the list_enabled metric pair so the dashboard
// surfaces saturation issues independently from the per-plugin reads.
func (s *PluginService) GetEnabledPlugins(ctx context.Context) ([]dtos.PluginManifestResponse, error) {
	start := time.Now()
	plugins, err := s.deps.PluginLoader.GetAllEnabled(ctx)
	if err != nil {
		s.trackPluginMetrics("list_enabled", "error", start)
		return nil, err
	}
	s.trackPluginMetrics("list_enabled", "success", start)
	return plugins, nil
}

// HandleFanoutEvent processes a FANOUT cache-invalidation message: parse
// the payload (drop on bad shape) and invalidate the matching manifest
// across L0+L1 of the local TieredCache. The MinIO L2 source of truth is
// not touched.
func (s *PluginService) HandleFanoutEvent(msg *natsModel.Message) {
	payload, ok := s.parseFanoutPayload(msg)
	if !ok {
		return
	}
	ctx := context.Background()
	s.invalidatePluginCache(ctx, payload)
}

// trackPluginMetrics is the shared metric-recording shim used by every
// plugin operation: increments the per-(op, outcome) counter and observes
// the duration histogram. Centralising it keeps every public method readable.
func (s *PluginService) trackPluginMetrics(op, outcome string, start time.Time) {
	s.deps.Metrics.PluginOperations.WithLabelValues(op, outcome).Inc()
	s.deps.Metrics.PluginOperationDuration.WithLabelValues(op).Observe(time.Since(start).Seconds())
}

// notFoundPluginManifest wraps the canonical 404 envelope used by every
// plugin endpoint, keeping the message consistent across CRUD operations.
func notFoundPluginManifest() error {
	return &customErrors.ServerCustomError{Code: status.NOT_FOUND, Errors: []string{"Plugin manifest not found"}}
}

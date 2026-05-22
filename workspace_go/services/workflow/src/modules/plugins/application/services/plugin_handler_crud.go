package services

import (
	"context"
	"fmt"
	"time"

	"workflow/src/modules/plugins/application/dtos"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// applyOrgContextToPlugin copies the resolved org id and pathKey from the
// request context onto the manifest entity so the stored document carries
// the same multi-tenant labels as future queries. Pre-existing values on
// the entity are kept when the request context is empty.
func (s *PluginService) applyOrgContextToPlugin(entity *dtos.PluginManifestResponse, requestContext *reqCtx.RequestContext) {
	if requestContext.OrgContext != nil && *requestContext.OrgContext != "" {
		if orgObjectId, err := model.ToObjectID(*requestContext.OrgContext); err == nil {
			entity.OrgId = &orgObjectId
		}
	}
	if requestContext.OrgContextData != nil && requestContext.OrgContextData.PathKey != "" {
		entity.PathKey = requestContext.OrgContextData.PathKey
	}
}

// applyPluginCreateDefaults stamps the create/updated timestamps with a
// single `now` so the row's "first seen" time is the same as its initial
// version time — downstream tooling assumes both are the create instant.
func (s *PluginService) applyPluginCreateDefaults(entity *dtos.PluginManifestResponse) {
	now := time.Now()
	entity.Created = now
	entity.Updated = now
}

// fetchPluginForUpdate loads the pre-update manifest. A missing row surfaces
// as 404 via the same envelope the post-save fallback uses, so the client
// never has to distinguish a vanished row from a vanished update.
func (s *PluginService) fetchPluginForUpdate(ctx context.Context, pluginId *string, start time.Time) (*dtos.PluginManifestResponse, error) {
	existing, err := s.deps.PluginRepo.FindById(ctx, pluginId)
	if err != nil || existing == nil || existing.ID.IsZero() {
		s.trackPluginMetrics("update", "error", start)
		return nil, notFoundPluginManifest()
	}
	return existing, nil
}

// buildPluginUpdatePayload converts the partial DTO to a Mongo $set map and
// stamps the updated timestamp. Field-by-field copy lives in the package-
// level buildUpdateMap helper; this wrapper just adds the timestamp.
func (s *PluginService) buildPluginUpdatePayload(dto *dtos.PluginManifestUpdate) map[string]any {
	fieldsToUpdate := buildUpdateMap(dto)
	fieldsToUpdate["updated"] = time.Now()
	return fieldsToUpdate
}

// persistPluginUpdate runs the update and validates the post-save row. A
// repo error or a nil/zero-id result both surface as the canonical 404 — a
// row may have been deleted between the pre-fetch and the save (race) and
// the same envelope serves both cases.
func (s *PluginService) persistPluginUpdate(ctx context.Context, pluginId *string, fieldsToUpdate map[string]any, start time.Time) (*dtos.PluginManifestResponse, error) {
	updated, err := s.deps.PluginRepo.FindByIdAndUpdate(ctx, pluginId, fieldsToUpdate)
	if err != nil {
		s.trackPluginMetrics("update", "error", start)
		return nil, err
	}
	if updated == nil || updated.ID.IsZero() {
		s.trackPluginMetrics("update", "error", start)
		return nil, notFoundPluginManifest()
	}
	return updated, nil
}

// fetchPluginForDelete is a best-effort pre-delete read. Failures log a
// warning but do not block the delete itself, since the row removal is the
// primary intent of the call. The returned value may be nil.
func (s *PluginService) fetchPluginForDelete(ctx context.Context, pluginId *string) *dtos.PluginManifestResponse {
	existing, fetchErr := s.deps.PluginRepo.FindById(ctx, pluginId)
	if fetchErr != nil {
		logger.Warn(fmt.Sprintf("[SERVICE:Plugin] Failed to fetch plugin %s before delete: %v", *pluginId, fetchErr))
	}
	return existing
}

// invalidateAndFanoutPlugin invalidates the per-plugin entry plus the "all
// enabled" list in the local TieredCache, then publishes a FANOUT message
// so peer pods do the same. Action label flows through to the FANOUT
// payload for observability.
func (s *PluginService) invalidateAndFanoutPlugin(ctx context.Context, pluginId, action string) {
	s.deps.PluginLoader.Invalidate(ctx, pluginId)
	s.deps.PluginLoader.InvalidateAll(ctx)
	s.publishPluginInvalidate(ctx, pluginId, action)
}

// invalidateAfterPluginDelete is the post-delete variant: nil-safe so a
// failed pre-delete fetch simply skips invalidation rather than crashing
// the request. Mirrors invalidateAndFanoutPlugin's local + FANOUT pair.
func (s *PluginService) invalidateAfterPluginDelete(ctx context.Context, existing *dtos.PluginManifestResponse) {
	if existing == nil {
		return
	}
	s.invalidateAndFanoutPlugin(ctx, existing.PluginID, "delete")
}

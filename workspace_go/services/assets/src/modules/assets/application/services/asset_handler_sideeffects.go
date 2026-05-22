package services

import (
	ctx "context"
	"encoding/json"
	"fmt"
	"time"

	"assets/src/modules/assets/domain/entities"

	assetsContract "github.com/Mapex-Solutions/MapexOS/contracts/services/assets/assets"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// writeAssetMetadata writes the AssetReadModel to object storage (L2).
//
// The read model carries everything every downstream consumer needs:
// Router / JS-Executor / Events read the public asset state; the
// mapex-mqtt-broker plugin reads `Protocol.Mqtt.PasswordHash` and
// `CurrentCert.Serial` from the same payload to decide MQTT CONNECTs
// locally (bcrypt compare or cert serial match). There is no second
// bucket, no auth callout, no Redis credential cache — one read model,
// one L2.
//
// The adapter encapsulates all MinIO details. Template org id flows
// through so the read model exposes whether the template is a public
// (mapexos_public) one or org-scoped.
//
// Parameters:
//   - ctx: Context for cancellation and timeouts
//   - asset: The asset entity to publish
//
// Returns:
//   - templateOrgId: The resolved template organization ID
func (s *AssetService) writeAssetMetadata(ctx ctx.Context, asset *entities.Asset) string {
	return s.syncAssetL2(ctx, asset)
}

// getTemplateOrgId determines the organization ID for the asset's template.
//
// Logic:
//   - If template.IsSystem == true → "mapexos_public" (accessible to all orgs)
//   - Otherwise → asset's orgId (template belongs to same org as asset)
//
// Parameters:
//   - ctx: Context for cancellation and timeouts
//   - asset: The asset entity
//
// Returns:
//   - Template organization ID string
func (s *AssetService) getTemplateOrgId(ctx ctx.Context, asset *entities.Asset) string {
	const publicOrgId = "mapexos_public"

	if asset.AssetTemplateID.IsZero() {
		return publicOrgId
	}

	// Fetch template to check IsSystem flag
	templateId := asset.AssetTemplateID.Hex()
	template, err := s.deps.AssetTemplateRepo.FindById(ctx, &templateId)
	if err != nil || template == nil {
		// Fallback to asset's org if template not found
		return asset.OrgID.Hex()
	}

	if template.IsSystem {
		return publicOrgId
	}

	return asset.OrgID.Hex()
}

// deleteAssetMetadata removes the asset read model from object storage (L2 cache).
//
// This delegates to the AssetStoragePort which handles storage operations.
// The adapter encapsulates all MinIO details.
//
// Key format: {orgId}/{assetUUID}.json
//
// Parameters:
//   - ctx: Context for cancellation and timeouts
//   - asset: The asset entity being deleted
func (s *AssetService) deleteAssetMetadata(ctx ctx.Context, asset *entities.Asset) {
	s.deleteAssetL2(ctx, asset.OrgID.Hex(), asset.AssetUUID)
}

// publishAssetInvalidate publishes a FANOUT message for cache invalidation.
//
// Consuming services (Router, JS-Executor, Events, mapex-mqtt-broker
// plugin) subscribe to this subject and invalidate their TieredCache
// (L0/L1) when they receive a message. L2 (MinIO) is the source of
// truth and is already updated at this point.
//
// Payload format: { "orgId": "...", "assetUUID": "..." }
// Cache key format: {orgId}/{assetUUID}
//
// Parameters:
//   - ctx: Context for cancellation and timeouts
//   - asset: The asset entity to invalidate
func (s *AssetService) publishAssetInvalidate(ctx ctx.Context, asset *entities.Asset) {
	if asset == nil || asset.AssetUUID == "" {
		return
	}

	// Build payload with orgId + assetUUID for cache key construction
	invalidatePayload := assetsContract.AssetInvalidatePayload{
		OrgId:     asset.OrgID.Hex(),
		AssetUUID: asset.AssetUUID,
	}

	payload, err := json.Marshal(invalidatePayload)
	if err != nil {
		logger.Warn(fmt.Sprintf("[SERVICE:Asset] Failed to marshal FANOUT payload for %s: %v", asset.AssetUUID, err))
		return
	}

	if err := s.deps.NatsBus.PublishFanout(ctx, assetsContract.FanoutAssetSubject, payload); err != nil {
		logger.Warn(fmt.Sprintf("[SERVICE:Asset] Failed to publish FANOUT for %s: %v", asset.AssetUUID, err))
	}
}

// clearHealthStateOnDisable cleans up Redis health state and resets the
// Mongo healthStatus when the user toggles HealthMonitor.Enabled from true
// to false. Idempotent — does nothing if there's no transition.
//
// Called by updateAssetById (with before/after entities) AND deleteAssetById
// (which passes the deleted asset as `before` and `nil` as `after`).
//
// The Mongo healthStatus is reset to "unknown" via the same
// UpdateHealthStatusWithChangedAt call used by the heartbeat path — the
// healthStatusChangedAt timestamp moves to the disable moment so the
// badge UI reflects the change.
func (s *AssetService) clearHealthStateOnDisable(ctx ctx.Context, before *entities.Asset, after *entities.Asset) {
	wasEnabled := before != nil && before.HealthMonitor.IsActive()
	isEnabled := after != nil && after.HealthMonitor.IsActive()
	if !wasEnabled || isEnabled {
		return
	}

	orgId := before.OrgID.Hex()
	assetUUID := before.AssetUUID

	if err := s.deps.HealthLifecycle.ClearAssetState(ctx, orgId, assetUUID); err != nil {
		logger.Warn(fmt.Sprintf("[SERVICE:Asset] Failed to clear health state for %s: %v", assetUUID, err))
	}

	unknown := "unknown"
	if err := s.deps.AssetRepo.UpdateHealthStatusWithChangedAt(ctx, &assetUUID, unknown, time.Now().UTC()); err != nil {
		logger.Warn(fmt.Sprintf("[SERVICE:Asset] Failed to reset healthStatus for %s: %v", assetUUID, err))
	}
}

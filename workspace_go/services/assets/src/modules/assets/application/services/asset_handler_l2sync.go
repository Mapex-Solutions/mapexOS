package services

import (
	ctx "context"
	"encoding/json"
	"fmt"

	"assets/src/modules/assets/domain/entities"

	assetsAuthContract "github.com/Mapex-Solutions/MapexOS/contracts/services/assets/auth"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// syncAssetL2 orchestrates the two MinIO writes after a CRUD: the
// full read model (mapex-assets/{orgId}/{assetUUID}.json — existing
// behavior, consumed by Router / Events / JS-Executor) AND the slim
// auth projection (mapex-asset-auth/{assetUUID}.json — new bucket,
// consumed by the broker plugin only).
//
// On either failure the handler publishes a retry hint to the L2
// writes stream so the in-module fallback consumer can reconcile
// against current Mongo state once MinIO recovers. FANOUT is NOT
// emitted here — the caller in asset_handler_sideeffects.go emits
// it after this returns.
func (s *AssetService) syncAssetL2(ctx ctx.Context, asset *entities.Asset) string {
	templateOrgId := s.getTemplateOrgId(ctx, asset)

	fullErr := s.deps.AssetStoragePort.WriteAsset(ctx, asset, templateOrgId)
	if fullErr != nil {
		logger.Warn(fmt.Sprintf("[SERVICE:Asset] L2 full write failed for %s: %v", asset.AssetUUID, fullErr))
	}

	var authErr error
	if asset.Protocol.Type == "mqtt" && asset.Protocol.Mqtt != nil {
		projection := buildAuthProjection(asset)
		authErr = s.deps.AssetStoragePort.WriteAssetAuth(ctx, projection)
		if authErr != nil {
			logger.Warn(fmt.Sprintf("[SERVICE:Asset] L2 auth write failed for %s: %v", asset.AssetUUID, authErr))
		}
	}

	if fullErr != nil || authErr != nil {
		s.publishL2Retry(ctx, asset)
	}

	return templateOrgId
}

// deleteAssetL2 removes both projections from MinIO. Best-effort: a
// failure here doesn't propagate to the caller, mirroring the existing
// delete semantics (Mongo is the source of truth; stale L2 entries
// reconcile on the next CRUD).
func (s *AssetService) deleteAssetL2(ctx ctx.Context, orgId, assetUUID string) {
	if err := s.deps.AssetStoragePort.DeleteAsset(ctx, orgId, assetUUID); err != nil {
		logger.Warn(fmt.Sprintf("[SERVICE:Asset] L2 full delete failed for %s: %v", assetUUID, err))
	}
	if err := s.deps.AssetStoragePort.DeleteAssetAuth(ctx, assetUUID); err != nil {
		logger.Warn(fmt.Sprintf("[SERVICE:Asset] L2 auth delete failed for %s: %v", assetUUID, err))
	}
}

// publishL2Retry sends a retry hint to the L2 writes stream with
// `Nats-Msg-Id: asset:{id}` so NATS-native dedup (5s window)
// coalesces rapid successive failures on the same asset.
func (s *AssetService) publishL2Retry(ctx ctx.Context, asset *entities.Asset) {
	payload, err := json.Marshal(map[string]string{"assetId": asset.ID.Hex()})
	if err != nil {
		logger.Error(err, fmt.Sprintf("[SERVICE:Asset] L2 retry marshal failed for %s", asset.AssetUUID))
		return
	}
	msgId := "asset:" + asset.ID.Hex()
	if pubErr := s.deps.L2WritesPublisher.PublishRetry(ctx, assetsAuthContract.L2WritesAssetSubject, msgId, payload); pubErr != nil {
		logger.Error(pubErr, fmt.Sprintf("[SERVICE:Asset] L2 retry publish failed for %s", asset.AssetUUID))
	}
}

// fetchAssetByID looks up the current asset state from Mongo by id.
// Used by ProcessL2WriteRetry — the retry stream message carries only
// the id; the actual payload is rebuilt from current Mongo state so
// a stale event cannot overwrite newer data.
func (s *AssetService) fetchAssetByID(ctx ctx.Context, assetId string) (*entities.Asset, error) {
	asset, err := s.deps.AssetRepo.FindById(ctx, &assetId)
	if err != nil {
		return nil, fmt.Errorf("l2 retry: find asset %s: %w", assetId, err)
	}
	return asset, nil // (nil, nil) when not found — caller treats as "deleted, drop the retry"
}

// buildAuthProjection translates the persistence entity into the
// cross-service auth projection contract. Pure function — no
// repository access, no logging. Lives at file scope so the L3
// internal endpoint handler can reuse it without going through
// the service struct.
func buildAuthProjection(asset *entities.Asset) assetsAuthContract.AuthProjection {
	proj := assetsAuthContract.AuthProjection{
		AssetUUID: asset.AssetUUID,
		OrgId:     asset.OrgID.Hex(),
		Enabled:   asset.Enabled,
		Type:      "mqtt",
	}
	if asset.Protocol.Mqtt != nil {
		proj.AuthType = asset.Protocol.Mqtt.AuthType
		proj.PasswordHash = asset.Protocol.Mqtt.PasswordHash
	}
	if asset.CurrentCert != nil {
		proj.CurrentCertSerial = asset.CurrentCert.Serial
	}
	return proj
}

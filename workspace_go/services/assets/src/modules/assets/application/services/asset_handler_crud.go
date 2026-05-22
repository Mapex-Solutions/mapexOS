package services

import (
	ctx "context"
	"fmt"
	"time"

	"assets/src/modules/assets/application/constants"
	"assets/src/modules/assets/application/dtos"
	"assets/src/modules/assets/domain/entities"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
	"github.com/Mapex-Solutions/mapexGoKit/utils/mapper"
	"golang.org/x/crypto/bcrypt"
)

// bindOrgContextOnCreate copies orgId + pathKey from the request context
// onto the creation DTO so the entity inherits multi-tenant scoping.
func (s *AssetService) bindOrgContextOnCreate(rc *reqCtx.RequestContext, dto *dtos.AssetCreateDTO) {
	if rc.OrgContext != nil && *rc.OrgContext != "" {
		if orgObjectId, err := model.ToObjectID(*rc.OrgContext); err == nil {
			dto.OrgID = &orgObjectId
		}
	}
	if rc.OrgContextData != nil && rc.OrgContextData.PathKey != "" {
		pathKey := rc.OrgContextData.PathKey
		dto.PathKey = &pathKey
	}
}

// buildEntityFromDto maps the create DTO to an in-memory domain entity
// without persisting. The MQTT plaintext password (if any) is carried
// transiently on the entity by hashMqttPasswordIfNeeded — never written
// to Mongo as plaintext.
func (s *AssetService) buildEntityFromDto(dto *dtos.AssetCreateDTO) *entities.Asset {
	asset, _ := mapper.DtoToEntity[dtos.AssetCreateDTO, entities.Asset](dto)
	if dto.AssetTemplateID != "" {
		if templateId, err := model.ToObjectID(dto.AssetTemplateID); err == nil {
			asset.AssetTemplateID = templateId
		}
	}
	now := time.Now()
	asset.Created = now
	asset.Updated = now
	return asset
}

// persistNewAsset inserts the in-memory entity into Mongo. Returns the
// entity as Mongo materialized it (so the caller picks up the assigned
// `_id`).
func (s *AssetService) persistNewAsset(c ctx.Context, asset *entities.Asset) (*entities.Asset, error) {
	return s.deps.AssetRepo.Create(c, asset)
}

// hashMqttPasswordIfNeeded bcrypts the operator-supplied MQTT password
// onto the entity for MQTT-protocol assets. The plaintext is never
// persisted; the bcrypt hash is what the datasource-mqtt auth callout
// validates against on broker CONNECT. Returns BAD_REQUEST when an
// MQTT-protocol asset arrives without a password — defense in depth on
// top of the contract-level required-field validator.
func (s *AssetService) hashMqttPasswordIfNeeded(asset *entities.Asset, dto *dtos.AssetCreateDTO) error {
	if asset.Protocol.Type != "mqtt" || asset.Protocol.Mqtt == nil {
		return nil
	}
	// Password is optional on MQTT-protocol assets: cert-mode assets
	// are created without a password and gain auth later via the
	// `mqttcerts` module (operator clicks "Generate certificate" on
	// the asset details). The broker plugin fails closed if a CONNECT
	// arrives for an asset with neither passwordHash nor currentCert,
	// so the transient unauth state during the create-then-issue-cert
	// flow is safe.
	if dto.Protocol.Mqtt == nil || dto.Protocol.Mqtt.Password == "" {
		return nil
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(dto.Protocol.Mqtt.Password), constants.MqttPasswordBcryptCost)
	if err != nil {
		return fmt.Errorf("hash mqtt password: %w", err)
	}
	asset.Protocol.Mqtt.PasswordHash = string(hash)
	return nil
}

// fanoutCreateSideEffects writes the AssetReadModel to MinIO (L2) and
// invalidates the org counter cache. The read model includes
// PasswordHash + CurrentCert, which the mapex-mqtt-broker plugin reads
// from L2 on every CONNECT — so this write MUST land before the
// device's first CONNECT for password / cert auth to work.
func (s *AssetService) fanoutCreateSideEffects(c ctx.Context, rc *reqCtx.RequestContext, asset *entities.Asset) {
	s.writeAssetMetadata(c, asset)
	if rc.OrgContext != nil {
		counterKey := s.deps.CacheKeyBuilder.BuildCounterCacheKey(*rc.OrgContext)
		_ = s.deps.AppCache.Del(c, counterKey)
	}
}

// applyAssetPatch translates the update DTO into a Mongo $set map,
// converts AssetTemplateID, refreshes Updated, and runs FindByIdAndUpdate.
// When the patch carries a new plaintext MQTT password, this method
// bcrypt-hashes it and rewrites the field as `protocol.mqtt.passwordHash`
// so plaintext never lands in Mongo. The hash uses the same cost as
// the create path (constants.MqttPasswordBcryptCost).
func (s *AssetService) applyAssetPatch(c ctx.Context, assetId *string, dto *dtos.AssetUpdateDTO) (*entities.Asset, error) {
	fields, _ := mapper.DtoToMap(dto)
	if dto.AssetTemplateID != nil && *dto.AssetTemplateID != "" {
		if templateId, err := model.ToObjectID(*dto.AssetTemplateID); err == nil {
			fields["assetTemplateId"] = templateId
		}
	}
	if dto.Protocol != nil && dto.Protocol.Mqtt != nil && dto.Protocol.Mqtt.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(dto.Protocol.Mqtt.Password), constants.MqttPasswordBcryptCost)
		if err != nil {
			return nil, fmt.Errorf("hash mqtt password: %w", err)
		}
		fields["protocol.mqtt.passwordHash"] = string(hash)
	}
	delete(fields, "protocol.mqtt.password")
	fields["updated"] = time.Now()
	updated, _ := s.deps.AssetRepo.FindByIdAndUpdate(c, assetId, fields)
	if updated.ID.IsZero() {
		return nil, fmt.Errorf("asset not found after update")
	}
	return updated, nil
}

// fanoutUpdateSideEffects rewrites the MinIO read model, publishes
// FANOUT invalidation (which evicts the broker plugin's L1 + every
// other consumer's local cache), and clears health state when the user
// flipped HealthMonitor.Enabled from true to false.
func (s *AssetService) fanoutUpdateSideEffects(c ctx.Context, before, after *entities.Asset) {
	s.writeAssetMetadata(c, after)
	s.publishAssetInvalidate(c, after)
	s.clearHealthStateOnDisable(c, before, after)
}

// tearDownAssetCaches removes everything cached for one asset before
// the Mongo delete: MinIO read model, FANOUT invalidation broadcast,
// Redis health state, and the org counter cache. Cache-first order
// ensures retries after partial failure stay consistent.
func (s *AssetService) tearDownAssetCaches(c ctx.Context, asset *entities.Asset) {
	s.deleteAssetMetadata(c, asset)
	s.publishAssetInvalidate(c, asset)
	if err := s.deps.HealthLifecycle.ClearAssetState(c, asset.OrgID.Hex(), asset.AssetUUID); err != nil {
		logger.Warn(fmt.Sprintf("[SERVICE:Asset] Failed to clear health state on delete for %s: %v", asset.AssetUUID, err))
	}
	if !asset.OrgID.IsZero() {
		counterKey := s.deps.CacheKeyBuilder.BuildCounterCacheKey(asset.OrgID.Hex())
		_ = s.deps.AppCache.Del(c, counterKey)
	}
}

// recordAssetOp emits the count + duration metrics for one CRUD attempt.
// Centralized so every exit path stays observability-consistent.
func (s *AssetService) recordAssetOp(op, status string, start time.Time) {
	s.deps.Metrics.AssetOperations.WithLabelValues(op, status).Inc()
	s.deps.Metrics.AssetOperationDuration.WithLabelValues(op).Observe(time.Since(start).Seconds())
}

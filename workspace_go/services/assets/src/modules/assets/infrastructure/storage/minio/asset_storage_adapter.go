package minio

import (
	"context"
	"encoding/json"
	"fmt"

	"assets/src/modules/assets/application/converters"
	"assets/src/modules/assets/application/ports"
	"assets/src/modules/assets/domain/entities"

	assetsContract "github.com/Mapex-Solutions/MapexOS/contracts/services/assets/assets"
	assetsAuthContract "github.com/Mapex-Solutions/MapexOS/contracts/services/assets/auth"
	minioModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/minio"
)

// NewAssetStorageAdapter creates a new AssetStorageAdapter.
//
// Parameters:
//   - client:     Full read-model bucket client (key {orgId}/{assetUUID}.json)
//   - authClient: Slim auth projection bucket client (key {assetUUID}.json)
//
// Returns:
//   - ports.AssetStoragePort: The port interface implementation
func NewAssetStorageAdapter(client, authClient *minioModel.MinIOClient) ports.AssetStoragePort {
	return &AssetStorageAdapter{
		client:     client,
		authClient: authClient,
	}
}

// Compile-time check to ensure AssetStorageAdapter implements AssetStoragePort interface.
var _ ports.AssetStoragePort = (*AssetStorageAdapter)(nil)

// WriteAsset writes the asset read model to MinIO (L2 cache).
//
// This publishes the denormalized asset data for consumption by other services
// via TieredCache (Router, JS-Executor, Events).
//
// Key format: {orgId}/{assetUUID}.json
// Example: 507f1f77bcf86cd799439011/device-123-uuid.json
func (a *AssetStorageAdapter) WriteAsset(ctx context.Context, asset *entities.Asset, templateOrgId string) error {
	if asset == nil || asset.AssetUUID == "" {
		return nil
	}

	// Build AssetReadModel from entity (internal conversion)
	readModel := a.buildReadModel(asset, templateOrgId)

	// Serialize to JSON
	data, err := json.Marshal(readModel)
	if err != nil {
		return fmt.Errorf("failed to serialize asset read model: %w", err)
	}

	// Write to MinIO with orgId prefix for tenant isolation
	// Key format: {orgId}/{assetUUID}.json
	key := asset.OrgID.Hex() + "/" + asset.AssetUUID + ".json"
	if err := a.client.PutJSON(ctx, key, data); err != nil {
		return fmt.Errorf("failed to write asset to MinIO: %w", err)
	}

	return nil
}

// DeleteAsset removes the asset read model from MinIO (L2 cache).
//
// Called when an asset is deleted to ensure consuming services get cache miss.
//
// Key format: {orgId}/{assetUUID}.json
func (a *AssetStorageAdapter) DeleteAsset(ctx context.Context, orgId string, assetUUID string) error {
	if orgId == "" || assetUUID == "" {
		return nil
	}

	// Key format: {orgId}/{assetUUID}.json
	key := orgId + "/" + assetUUID + ".json"
	if err := a.client.Delete(ctx, key); err != nil {
		return fmt.Errorf("failed to delete asset from MinIO: %w", err)
	}

	return nil
}

// WriteAssetAuth writes the slim auth projection to the mapex-asset-auth
// bucket. The broker plugin reads this on every CONNECT lookup; payload
// is intentionally narrow (no full read model, no template metadata).
//
// Key format: {assetUUID}.json (flat — assetUUID is globally unique).
func (a *AssetStorageAdapter) WriteAssetAuth(ctx context.Context, projection assetsAuthContract.AuthProjection) error {
	if projection.AssetUUID == "" {
		return nil
	}
	data, err := json.Marshal(projection)
	if err != nil {
		return fmt.Errorf("failed to serialize auth projection: %w", err)
	}
	key := projection.AssetUUID + ".json"
	if err := a.authClient.PutJSON(ctx, key, data); err != nil {
		return fmt.Errorf("failed to write auth projection to MinIO: %w", err)
	}
	return nil
}

// DeleteAssetAuth removes the auth projection. Called when an asset
// is deleted so the broker's next L2 lookup misses (forcing an L3
// fallback that will return 404 and short-circuit the auth flow).
func (a *AssetStorageAdapter) DeleteAssetAuth(ctx context.Context, assetUUID string) error {
	if assetUUID == "" {
		return nil
	}
	key := assetUUID + ".json"
	if err := a.authClient.Delete(ctx, key); err != nil {
		return fmt.Errorf("failed to delete auth projection from MinIO: %w", err)
	}
	return nil
}

// buildReadModel converts an Asset entity to AssetReadModel for storage.
// This is an internal method used by WriteAsset for JSON serialization.
func (a *AssetStorageAdapter) buildReadModel(asset *entities.Asset, templateOrgId string) *assetsContract.AssetReadModel {
	description := ""
	if asset.Description != nil {
		description = *asset.Description
	}

	return &assetsContract.AssetReadModel{
		ID:                 asset.ID.Hex(),
		UUID:               asset.AssetUUID,
		OrgId:              asset.OrgID.Hex(),
		PathKey:            asset.PathKey,
		Enabled:            asset.Enabled,
		DebugEnabled:       asset.DebugEnabled,
		Name:               asset.Name,
		Description:        description,
		AssetTemplateID:    asset.AssetTemplateID.Hex(),
		AssetTemplateOrgID: templateOrgId,
		RouteGroupIds:      asset.RouteGroupIds,
		HealthMonitor:      converters.HealthMonitorEntityToContract(asset.HealthMonitor),
		HealthStatus:       asset.HealthStatus,
		Protocol:           a.convertProtocol(asset.Protocol),
		CurrentCert:        convertCurrentCert(asset.CurrentCert),
		Latitude:           asset.Latitude,
		Longitude:          asset.Longitude,
		Created:            asset.Created,
		Updated:            asset.Updated,
	}
}

// convertProtocol converts entity ProtocolType to contract ProtocolType.
// PasswordHash flows through here so the broker plugin's TieredCache (L2)
// has everything it needs to decide password-mode CONNECTs locally without
// a callout.
func (a *AssetStorageAdapter) convertProtocol(p entities.ProtocolType) *assetsContract.ProtocolType {
	if p.Type == "" {
		return nil
	}

	result := &assetsContract.ProtocolType{
		Type: p.Type,
	}

	if p.Mqtt != nil {
		result.Mqtt = &assetsContract.MqttConfig{
			ClientId:     p.Mqtt.ClientId,
			Username:     p.Mqtt.Username,
			AuthType:     p.Mqtt.AuthType,
			PasswordHash: p.Mqtt.PasswordHash,
		}
	}

	return result
}

// convertCurrentCert copies the active cert metadata so the broker plugin
// can validate cert-mode CONNECTs by serial-equality compare. Returns nil
// when the asset has no active cert (default-deny on the cert path).
func convertCurrentCert(c *entities.AssetCertificate) *assetsContract.AssetCertificate {
	if c == nil {
		return nil
	}
	return &assetsContract.AssetCertificate{
		Serial:      c.Serial,
		Fingerprint: c.Fingerprint,
		SubjectCN:   c.SubjectCN,
		IssuedAt:    c.IssuedAt,
		ExpiresAt:   c.ExpiresAt,
	}
}

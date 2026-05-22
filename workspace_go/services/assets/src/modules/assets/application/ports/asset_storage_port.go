package ports

import (
	"context"

	"assets/src/modules/assets/domain/entities"

	assetsAuthContract "github.com/Mapex-Solutions/MapexOS/contracts/services/assets/auth"
)

// AssetStoragePort defines the interface for asset read model storage operations.
//
// This port abstracts the object storage infrastructure (MinIO) following
// Hexagonal Architecture principles. The application layer depends on this
// interface, not on concrete storage implementations.
//
// The read model is a denormalized representation of asset data optimized
// for consumption by other services via TieredCache (L2 = MinIO).
//
// Implementations:
//   - AssetStorageAdapter (infrastructure/minio): MinIO-based implementation
type AssetStoragePort interface {
	// WriteAsset writes the asset read model to object storage (L2 cache).
	//
	// This publishes the denormalized asset data for consumption by other services
	// via TieredCache (Router, JS-Executor, Events).
	//
	// Parameters:
	//   - ctx: Context for controlling cancellation and timeouts
	//   - asset: The asset entity to publish
	//   - templateOrgId: The organization ID of the template (for JS-Executor template lookup)
	//                    Use "mapexos_public" for system templates, or the org's ID for private templates
	//
	// Returns:
	//   - error: nil on success, error on failure
	WriteAsset(ctx context.Context, asset *entities.Asset, templateOrgId string) error

	// DeleteAsset removes the asset read model from object storage (L2 cache).
	//
	// Called when an asset is deleted to ensure consuming services get cache miss.
	//
	// Key format: {orgId}/{assetUUID}.json
	//
	// Parameters:
	//   - ctx: Context for controlling cancellation and timeouts
	//   - orgId: The organization ID (hex string)
	//   - assetUUID: The asset UUID to delete
	//
	// Returns:
	//   - error: nil on success, error on failure
	DeleteAsset(ctx context.Context, orgId string, assetUUID string) error

	// WriteAssetAuth writes the slim auth projection to the
	// mapex-asset-auth bucket (key {assetUUID}.json). The projection
	// is the source the broker plugin reads on every CONNECT lookup.
	// Called alongside WriteAsset by the L2 sync helper; either
	// failure causes the side-effect handler to publish a retry
	// message to the L2 writes stream.
	WriteAssetAuth(ctx context.Context, projection assetsAuthContract.AuthProjection) error

	// DeleteAssetAuth removes the auth projection from the
	// mapex-asset-auth bucket. Called when an asset is deleted.
	// Best-effort like the existing DeleteAsset.
	DeleteAssetAuth(ctx context.Context, assetUUID string) error
}

package minio

import (
	minioModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/minio"
	"go.uber.org/dig"
)

/* AssetStorageProviderParams */

// AssetStorageProviderParams defines the dependencies for AssetStoragePort provider.
// Uses dig.In to enable automatic dependency injection with named parameters.
type AssetStorageProviderParams struct {
	dig.In

	// MinIOClient is injected with the "assets" name tag
	// Configured with bucket and key prefix for asset read models
	MinIOClient *minioModel.MinIOClient `name:"assets"`

	// MinIOAuthClient is injected with the "asset-auth" name tag.
	// Flat-key bucket consumed only by the broker plugin.
	MinIOAuthClient *minioModel.MinIOClient `name:"asset-auth"`
}

// AssetStorageAdapter implements AssetStoragePort using MinIO for object storage.
//
// The adapter handles two buckets:
//   - Full read-model bucket via `client` (key {orgId}/{assetUUID}.json)
//   - Slim auth projection bucket via `authClient` (key {assetUUID}.json)
type AssetStorageAdapter struct {
	client     *minioModel.MinIOClient
	authClient *minioModel.MinIOClient
}

package minio

import (
	"assets/src/modules/assets/application/ports"
)

// NewAssetStoragePort creates and returns an AssetStoragePort implementation.
//
// This provider creates an AssetStorageAdapter that implements the AssetStoragePort
// interface, following Hexagonal Architecture principles by returning the port
// interface instead of the concrete implementation.
//
// Parameters:
//   - params: Aggregated dependencies with named MinIO client (injected by dig)
//
// Returns:
//   - ports.AssetStoragePort: Port interface for asset storage operations
func NewAssetStoragePort(params AssetStorageProviderParams) ports.AssetStoragePort {
	return NewAssetStorageAdapter(params.MinIOClient, params.MinIOAuthClient)
}

package minio

import (
	minioModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/minio"
	"go.uber.org/dig"
)

/*
 * DEFINITION STORAGE ADAPTER TYPES
 * Infrastructure-level DI input struct and adapter struct for the MinIO
 * implementation of DefinitionStoragePort. Kept here (not in application/di)
 * so application code never references the concrete MinIO client.
 */

// DefinitionStorageProviderParams carries the concrete MinIO client to the
// NewDefinitionStoragePort provider. It lives in infrastructure because it
// references the concrete MinIO client type.
type DefinitionStorageProviderParams struct {
	dig.In

	// MinIOClient is injected with the "definitions" name tag.
	MinIOClient *minioModel.MinIOClient `name:"definitions"`
}

// DefinitionStorageAdapter implements DefinitionStoragePort using MinIO for object storage.
// Manages L2 cache for code node scripts and bytecodes.
type DefinitionStorageAdapter struct {
	client *minioModel.MinIOClient
}

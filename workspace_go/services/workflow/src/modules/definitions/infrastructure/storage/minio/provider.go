package minio

import (
	"workflow/src/modules/definitions/application/ports"
)

// NewDefinitionStoragePort creates and returns a DefinitionStoragePort implementation.
func NewDefinitionStoragePort(params DefinitionStorageProviderParams) ports.DefinitionStoragePort {
	return NewDefinitionStorageAdapter(params.MinIOClient)
}

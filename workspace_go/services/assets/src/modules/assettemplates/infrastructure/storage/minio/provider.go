package minio

import (
	"assets/src/modules/assettemplates/application/ports"
)

// NewTemplateStoragePort creates and returns a TemplateStoragePort implementation.
//
// This provider creates a TemplateStorageAdapter that implements the TemplateStoragePort
// interface, following Hexagonal Architecture principles by returning the port
// interface instead of the concrete implementation.
//
// Parameters:
//   - params: Aggregated dependencies with named MinIO client (injected by dig)
//
// Returns:
//   - ports.TemplateStoragePort: Port interface for template script storage operations
func NewTemplateStoragePort(params TemplateStorageProviderParams) ports.TemplateStoragePort {
	return NewTemplateStorageAdapter(params.MinIOClient)
}

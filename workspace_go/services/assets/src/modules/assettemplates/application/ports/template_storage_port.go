package ports

import (
	"context"

	"assets/src/modules/assettemplates/domain/entities"
)

// TemplateStoragePort defines the interface for template script storage operations.
//
// This port abstracts the object storage infrastructure (MinIO) following
// Hexagonal Architecture principles. The application layer depends on this
// interface, not on concrete storage implementations.
//
// Scripts are stored in L2 (MinIO) for consumption by other services
// via TieredCache (JS-Executor, Events).
//
// Implementations:
//   - TemplateStorageAdapter (infrastructure/minio): MinIO-based implementation
type TemplateStoragePort interface {
	// WriteScripts writes template scripts to object storage (L2 cache).
	//
	// This publishes the scripts (validator, conversion, test, processor)
	// for consumption by JS-Executor via TieredCache.
	//
	// Parameters:
	//   - ctx: Context for controlling cancellation and timeouts
	//   - template: The template entity containing scripts
	//
	// Returns:
	//   - error: nil on success, error on failure
	WriteScripts(ctx context.Context, template *entities.Assettemplate) error

	// DeleteScripts removes template scripts from object storage (L2 cache).
	//
	// Called when a template is deleted to ensure consuming services get cache miss.
	//
	// Key format: {orgId}/{templateId}.json
	// If IsSystem=true, orgId should be "mapexos_public"
	//
	// Parameters:
	//   - ctx: Context for controlling cancellation and timeouts
	//   - orgId: The organization ID (or "mapexos_public" for system templates)
	//   - templateId: The template ID to delete
	//
	// Returns:
	//   - error: nil on success, error on failure
	DeleteScripts(ctx context.Context, orgId string, templateId string) error
}

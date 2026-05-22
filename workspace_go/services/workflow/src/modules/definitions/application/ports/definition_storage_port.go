package ports

import (
	"context"
)

// DefinitionStoragePort abstracts the object storage operations for workflow definition scripts.
// Follows Hexagonal Architecture: application depends on port, not MinIO implementation.
//
// Manages L2 cache (MinIO) for:
//   - Script source:  {orgId}/{definitionId}/scripts/{nodeId}.json
//   - Bytecode:       {orgId}/{definitionId}/bytecode/{nodeId}.bin
type DefinitionStoragePort interface {
	// WriteNodeScript writes a code node's script source to MinIO (L2 cache).
	// Key format: {orgId}/{definitionId}/scripts/{nodeId}.json
	WriteNodeScript(ctx context.Context, orgId, definitionId, nodeId string, script []byte) error

	// DeleteNodeScript removes a code node's script source from MinIO.
	// Key format: {orgId}/{definitionId}/scripts/{nodeId}.json
	DeleteNodeScript(ctx context.Context, orgId, definitionId, nodeId string) error

	// DeleteNodeBytecode removes a code node's compiled bytecode from MinIO.
	// Called when script source changes (invalidates old bytecode).
	// Key format: {orgId}/{definitionId}/bytecode/{nodeId}.bin
	DeleteNodeBytecode(ctx context.Context, orgId, definitionId, nodeId string) error

	// DeleteAllNodeData removes all scripts and bytecodes for the given node IDs.
	// Used on definition deletion to clean up all L2 cache entries.
	DeleteAllNodeData(ctx context.Context, orgId, definitionId string, nodeIds []string) error
}

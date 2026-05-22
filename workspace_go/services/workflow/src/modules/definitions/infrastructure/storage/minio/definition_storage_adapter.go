package minio

import (
	"context"
	"fmt"

	"workflow/src/modules/definitions/application/ports"

	minioModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/minio"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

// NewDefinitionStorageAdapter creates a new DefinitionStorageAdapter.
func NewDefinitionStorageAdapter(client *minioModel.MinIOClient) ports.DefinitionStoragePort {
	return &DefinitionStorageAdapter{client: client}
}

// Compile-time check
var _ ports.DefinitionStoragePort = (*DefinitionStorageAdapter)(nil)

// WriteNodeScript writes a code node's script source to MinIO (L2 cache).
// Key format: {orgId}/{definitionId}/scripts/{nodeId}.json
func (a *DefinitionStorageAdapter) WriteNodeScript(ctx context.Context, orgId, definitionId, nodeId string, script []byte) error {
	if definitionId == "" || nodeId == "" {
		return nil
	}

	key := buildScriptKey(orgId, definitionId, nodeId)
	if err := a.client.PutJSON(ctx, key, script); err != nil {
		return fmt.Errorf("failed to write script for node %s: %w", nodeId, err)
	}
	return nil
}

// DeleteNodeScript removes a code node's script source from MinIO.
func (a *DefinitionStorageAdapter) DeleteNodeScript(ctx context.Context, orgId, definitionId, nodeId string) error {
	if definitionId == "" || nodeId == "" {
		return nil
	}

	key := buildScriptKey(orgId, definitionId, nodeId)
	if err := a.client.Delete(ctx, key); err != nil {
		return fmt.Errorf("failed to delete script for node %s: %w", nodeId, err)
	}
	return nil
}

// DeleteNodeBytecode removes a code node's compiled bytecode from MinIO.
func (a *DefinitionStorageAdapter) DeleteNodeBytecode(ctx context.Context, orgId, definitionId, nodeId string) error {
	if definitionId == "" || nodeId == "" {
		return nil
	}

	key := buildBytecodeKey(orgId, definitionId, nodeId)
	if err := a.client.Delete(ctx, key); err != nil {
		return fmt.Errorf("failed to delete bytecode for node %s: %w", nodeId, err)
	}
	return nil
}

// DeleteAllNodeData removes all scripts and bytecodes for the given node IDs.
func (a *DefinitionStorageAdapter) DeleteAllNodeData(ctx context.Context, orgId, definitionId string, nodeIds []string) error {
	for _, nodeId := range nodeIds {
		if err := a.DeleteNodeScript(ctx, orgId, definitionId, nodeId); err != nil {
			logger.Warn(fmt.Sprintf("[INFRA:DefinitionStorage] Failed to delete script %s: %v", nodeId, err))
		}
		if err := a.DeleteNodeBytecode(ctx, orgId, definitionId, nodeId); err != nil {
			logger.Warn(fmt.Sprintf("[INFRA:DefinitionStorage] Failed to delete bytecode %s: %v", nodeId, err))
		}
	}
	return nil
}

// buildScriptKey constructs the MinIO key for a code node's script source.
func buildScriptKey(orgId, definitionId, nodeId string) string {
	return fmt.Sprintf("%s/%s/scripts/%s.json", orgId, definitionId, nodeId)
}

// buildBytecodeKey constructs the MinIO key for a code node's compiled bytecode.
func buildBytecodeKey(orgId, definitionId, nodeId string) string {
	return fmt.Sprintf("%s/%s/bytecode/%s.bin", orgId, definitionId, nodeId)
}

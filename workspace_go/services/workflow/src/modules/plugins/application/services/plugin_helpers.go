package services

import (
	"context"
	"encoding/json"
	"fmt"

	"workflow/src/modules/plugins/application/constants"
	"workflow/src/modules/plugins/application/dtos"
	"workflow/src/modules/plugins/application/ports"
	"workflow/src/modules/plugins/domain/entities"

	"github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
)

/* PRIVATE METHODS — FANOUT */

// publishPluginInvalidate publishes a FANOUT message for cache invalidation.
// All pods subscribing to the FANOUT subject will invalidate L0 + L1 for the specified pluginId.
func (s *PluginService) publishPluginInvalidate(ctx context.Context, pluginId string, action string) {
	payload := ports.PluginInvalidatePayload{
		PluginID: pluginId,
		Action:   action,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		logger.Warn(fmt.Sprintf("[SERVICE:Plugin] Failed to marshal FANOUT payload for %s: %v", pluginId, err))
		return
	}

	if err := s.deps.NatsBus.PublishFanout(ctx, constants.FanoutPluginSubject, data); err != nil {
		logger.Warn(fmt.Sprintf("[SERVICE:Plugin] Failed to publish FANOUT for %s: %v", pluginId, err))
	}
}

/* PRIVATE HELPERS */

// buildUpdateMap converts a PluginManifestUpdate DTO to a map for MongoDB $set.
// Only includes fields that are non-nil (partial update).
func buildUpdateMap(dto *dtos.PluginManifestUpdate) map[string]any {
	fields := make(map[string]any)

	if dto.Name != nil {
		fields["name"] = *dto.Name
	}
	if dto.Version != nil {
		fields["version"] = *dto.Version
	}
	if dto.Category != nil {
		fields["category"] = *dto.Category
	}
	if dto.Icon != nil {
		fields["icon"] = *dto.Icon
	}
	if dto.Color != nil {
		fields["color"] = *dto.Color
	}
	if dto.Description != nil {
		fields["description"] = *dto.Description
	}
	if dto.Defaults != nil {
		fields["defaults"] = dto.Defaults
	}
	if dto.Credentials != nil {
		fields["credentials"] = dto.Credentials
	}
	if dto.NodeTypes != nil {
		sanitized := sanitizeNodeTypes(*dto.NodeTypes)
		fields["nodeTypes"] = sanitized
	}
	if dto.Metadata != nil {
		fields["metadata"] = dto.Metadata
	}
	if dto.Enabled != nil {
		fields["enabled"] = *dto.Enabled
	}
	if dto.FetchOptions != nil {
		fields["fetchOptions"] = dto.FetchOptions
	}
	if dto.Author != nil {
		fields["author"] = *dto.Author
	}
	if dto.Tags != nil {
		fields["tags"] = dto.Tags
	}

	return fields
}

// sanitizeNodeTypes strips script-type actions from nodeTypes received via API.
// Script actions can ONLY come from audited manifests — never from user API updates.
// This prevents arbitrary code injection via the plugin update endpoint.
func sanitizeNodeTypes(nodeTypes []entities.NodeTypeManifest) []entities.NodeTypeManifest {
	for i := range nodeTypes {
		// Sanitize operations
		for key, op := range nodeTypes[i].Operations {
			if op.Type == "script" {
				delete(nodeTypes[i].Operations, key)
				logger.Warn(fmt.Sprintf("[SERVICE:Plugin] Stripped script operation '%s' from node type '%s' — scripts only from audited manifests", key, nodeTypes[i].Type))
			}
		}
		// Sanitize hooks
		if nodeTypes[i].Hooks != nil {
			if nodeTypes[i].Hooks.Before != nil && nodeTypes[i].Hooks.Before.Type == "script" {
				nodeTypes[i].Hooks.Before = nil
				logger.Warn(fmt.Sprintf("[SERVICE:Plugin] Stripped script before hook from node type '%s'", nodeTypes[i].Type))
			}
			if nodeTypes[i].Hooks.After != nil && nodeTypes[i].Hooks.After.Type == "script" {
				nodeTypes[i].Hooks.After = nil
				logger.Warn(fmt.Sprintf("[SERVICE:Plugin] Stripped script after hook from node type '%s'", nodeTypes[i].Type))
			}
			if nodeTypes[i].Hooks.Destroy != nil && nodeTypes[i].Hooks.Destroy.Type == "script" {
				nodeTypes[i].Hooks.Destroy = nil
				logger.Warn(fmt.Sprintf("[SERVICE:Plugin] Stripped script destroy hook from node type '%s'", nodeTypes[i].Type))
			}
		}
	}
	return nodeTypes
}

package services

import (
	"sort"
	"strings"

	"workflow/src/modules/definitions/domain/entities"
)

/*
 * PLUGIN VALIDATION
 * Extracts required marketplace plugins from workflow nodes and computes
 * the definition status by comparing against the org's enabled plugins.
 */

// ExtractRequiredPlugins derives unique non-core plugin IDs from node types.
// Node types follow the pattern "{pluginId}/{action}" (e.g., "telegram/sendMessage").
// Core node types (prefixed with "core") are excluded.
func ExtractRequiredPlugins(nodes []entities.WorkflowNode) []string {
	seen := make(map[string]struct{})

	for _, node := range nodes {
		parts := strings.SplitN(node.Type, "/", 2)
		if len(parts) < 2 {
			continue
		}
		prefix := parts[0]
		if strings.HasPrefix(prefix, "core") {
			continue
		}
		seen[prefix] = struct{}{}
	}

	result := make([]string, 0, len(seen))
	for id := range seen {
		result = append(result, id)
	}
	sort.Strings(result)
	return result
}

// ComputeDefinitionStatus checks required plugins against enabled plugin IDs.
// Returns the computed status and the list of missing plugin IDs.
func ComputeDefinitionStatus(requiredPlugins []string, enabledPluginIDs []string) (entities.DefinitionStatus, []string) {
	if len(requiredPlugins) == 0 {
		return entities.StatusValid, nil
	}

	enabled := make(map[string]struct{}, len(enabledPluginIDs))
	for _, id := range enabledPluginIDs {
		enabled[id] = struct{}{}
	}

	var missing []string
	for _, req := range requiredPlugins {
		if _, ok := enabled[req]; !ok {
			missing = append(missing, req)
		}
	}

	if len(missing) > 0 {
		sort.Strings(missing)
		return entities.StatusPluginMissing, missing
	}

	return entities.StatusValid, nil
}

package dtos

import (
	contracts "github.com/Mapex-Solutions/MapexOS/contracts/services/workflow/plugins"
)

/*
 * Plugin DTO type aliases.
 * The wire-format authority lives in packages/contracts/services/workflow/plugins.
 * Domain entities under domain/entities/ also alias the same contract types,
 * keeping the entity files free of `json:"..."` tags per architecture rule.
 */
type (
	PluginManifestResponse = contracts.PluginManifest
	PluginManifestUpdate   = contracts.PluginManifestUpdate
	PluginIdDTO            = contracts.PluginId
	PluginQueryDTO         = contracts.PluginQuery
)

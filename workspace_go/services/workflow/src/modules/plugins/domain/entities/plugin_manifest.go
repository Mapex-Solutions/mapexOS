// Package entities holds the workflow plugins domain types.
//
// To comply with the architecture rule that domain entities under
// src/modules/{module}/domain/entities/*.go MUST NEVER carry json
// struct tags, every plugin manifest type is exposed here as a Go
// type alias to its authoritative wire-format definition in
// packages/contracts/services/workflow/plugins. The contract package
// is the JSON shape authority; persistence (bson) is co-located on the
// same struct in contracts because the Mongo document and the wire
// payload share a schema today. A future refactor can split the two
// behind a mapper without changing the public type surface used by
// services, repositories, and ports — they all reference these alias
// names.
package entities

import (
	contracts "github.com/Mapex-Solutions/MapexOS/contracts/services/workflow/plugins"
)

/*
 * HANDLE DEFINITIONS (same as workflow definition edges/handles)
 */

// HandleDef defines an input or output connection point on a node.
type HandleDef = contracts.HandleDef

/*
 * AVAILABLE OUTPUTS (describes what fields a node produces — informational for UI)
 */

// AvailableOutput describes a field available in the node's output for downstream use.
// Shown as hints in FieldSourceSelector when another node uses source type "nodeOutput".
type AvailableOutput = contracts.AvailableOutput

/*
 * UNIFIED ACTION CONTRACT
 * Used by: operations, fetchOptions, credential test, hooks.
 */

// HttpActionDef defines an HTTP request template.
type HttpActionDef = contracts.HttpActionDef

// MqttActionDef defines an MQTT publish template (future).
type MqttActionDef = contracts.MqttActionDef

// NatsActionDef defines a NATS publish template (future).
type NatsActionDef = contracts.NatsActionDef

// ScriptActionDef defines a JavaScript execution action.
// Always executed via V8 isolated-vm (JS Workflow Executor) at runtime.
// Backend rejects API updates with type "script" — scripts come ONLY from audited manifests.
type ScriptActionDef = contracts.ScriptActionDef

// ActionOutputDef defines how to extract/transform the response of an action.
type ActionOutputDef = contracts.ActionOutputDef

// ActionDef is the unified action contract.
type ActionDef = contracts.ActionDef

/*
 * FETCH OPTIONS (design-time dynamic data fetching for dropdowns)
 */

// FetchOptionsPagination defines how to paginate fetchOptions results.
type FetchOptionsPagination = contracts.FetchOptionsPagination

// FetchOptionsSearch defines server-side search support for fetchOptions.
type FetchOptionsSearch = contracts.FetchOptionsSearch

// FetchOptionsDef configures a dynamic options loader at the manifest level.
type FetchOptionsDef = contracts.FetchOptionsDef

/*
 * NODE PROPERTY SYSTEM (declarative form definitions)
 */

// PropertyRendering controls how a form field looks in the UI.
type PropertyRendering = contracts.PropertyRendering

// FetchOptionsRule defines which fetchOptions entry to use based on form state.
type FetchOptionsRule = contracts.FetchOptionsRule

// PropertyFetchOptions configures dynamic dropdown fetching for a fieldSource property.
type PropertyFetchOptions = contracts.PropertyFetchOptions

// PropertyOption represents a single option in a dropdown list.
type PropertyOption = contracts.PropertyOption

// DisplayOptions controls conditional visibility of a property based on other fields.
type DisplayOptions = contracts.DisplayOptions

// NodePropertyDef defines a single declarative form field for auto-generated config forms.
type NodePropertyDef = contracts.NodePropertyDef

/*
 * NODE HOOKS (lifecycle hooks per node type)
 */

// NodeHooks defines lifecycle hooks for a node type.
type NodeHooks = contracts.NodeHooks

/*
 * NODE TYPE MANIFEST (JSON-only — no Vue components)
 */

// NodeTimeoutDef specifies the default timeout for an async node type.
type NodeTimeoutDef = contracts.NodeTimeoutDef

// NodeTypeManifest defines a single node type within a plugin manifest.
type NodeTypeManifest = contracts.NodeTypeManifest

/*
 * CREDENTIAL DEFINITIONS
 */

// CredentialFieldDef defines a single credential input field.
type CredentialFieldDef = contracts.CredentialFieldDef

// CredentialDef describes one authentication method for a plugin.
type CredentialDef = contracts.CredentialDef

/*
 * PLUGIN METADATA (UI display in marketplace)
 */

// PluginMetadata contains visual metadata for marketplace display.
type PluginMetadata = contracts.PluginMetadata

/*
 * PLUGIN DEFAULTS (inherited by operations and fetchOptions)
 */

// PluginDefaults provides default values inherited by operations and fetchOptions.
type PluginDefaults = contracts.PluginDefaults

/*
 * PLUGIN MANIFEST (Root Aggregate — persisted in MongoDB)
 */

// PluginManifest is the root aggregate entity for an integration plugin.
// Multi-tenant visibility (same pattern as RouteGroup):
//   - IsTemplate=true: shared with child organizations (vendor templates)
//   - IsTemplate=false: local to a single organization (custom plugins)
type PluginManifest = contracts.PluginManifest

/*
 * PLUGIN MANIFEST UPDATE (Partial update — optional pointers)
 */

// PluginManifestUpdate is used for PATCH operations.
// Only non-nil fields are applied to the document.
type PluginManifestUpdate = contracts.PluginManifestUpdate

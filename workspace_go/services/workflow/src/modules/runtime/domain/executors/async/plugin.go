package async

import (
	"context"
	"fmt"
	"strings"
	"time"

	runtimePorts "workflow/src/modules/runtime/application/ports"
	defPorts "workflow/src/modules/definitions/application/ports"
	enginePorts "workflow/src/modules/engine/application/ports"
	pluginPorts "workflow/src/modules/plugins/application/ports"
	"workflow/src/modules/runtime/domain/entities"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
	"github.com/Mapex-Solutions/mapexGoKit/utils/templatereplace"
)

/*
 * PLUGIN EXECUTOR
 * Generic executor for all marketplace plugin nodes (non-core).
 * Resolves credentials, FieldValues, and templates before suspending.
 * The runtime dispatches the ready payload to the Triggers Service.
 */

// PluginExecutor handles execution of marketplace plugin nodes.
// Decrypts credentials, loads manifest, resolves FieldValues and templates,
// then suspends with a ready-to-execute payload in NodeState.
type PluginExecutor struct {
	vaultService runtimePorts.VaultPort
	pluginRepo        pluginPorts.PluginManifestRepository
	resolver          enginePorts.ValueResolverPort
}

// Compile-time check
var _ entities.NodeExecutor = (*PluginExecutor)(nil)

// NewPluginExecutor creates a new PluginExecutor with cross-module dependencies.
func NewPluginExecutor(
	vaultService runtimePorts.VaultPort,
	pluginRepo pluginPorts.PluginManifestRepository,
	resolver enginePorts.ValueResolverPort,
) entities.NodeExecutor {
	return &PluginExecutor{
		vaultService: vaultService,
		pluginRepo:        pluginRepo,
		resolver:          resolver,
	}
}

// NodeType returns "plugin" — this is a virtual type, never matched by prefix routing.
func (e *PluginExecutor) NodeType() string {
	return "plugin"
}

// Execute resolves all data and suspends with a ready payload for the Triggers Service.
func (e *PluginExecutor) Execute(ctx context.Context, execCtx *entities.NodeExecutionContext) (*entities.NodeExecutionResult, error) {
	cfg, ok := execCtx.ParsedConfig.(*entities.PluginNodeConfig)
	if !ok || cfg == nil {
		return nil, fmt.Errorf("plugin: missing or invalid config for node %s", execCtx.NodeID)
	}

	// 1. Extract pluginId from nodeType (e.g., "telegram/message" → "telegram")
	pluginID := extractPluginID(execCtx.NodeType)
	if pluginID == "" {
		return nil, fmt.Errorf("plugin: cannot extract pluginId from node type %s", execCtx.NodeType)
	}

	// 2. Load manifest
	manifest, err := e.pluginRepo.FindByPluginId(ctx, pluginID)
	if err != nil || manifest == nil {
		return nil, fmt.Errorf("plugin: manifest not found for plugin %s: %w", pluginID, err)
	}

	// 3. Find the nodeType manifest
	var nodeManifest *pluginPorts.NodeTypeManifest
	for i := range manifest.NodeTypes {
		if manifest.NodeTypes[i].Type == execCtx.NodeType {
			nodeManifest = &manifest.NodeTypes[i]
			break
		}
	}
	if nodeManifest == nil {
		return nil, fmt.Errorf("plugin: node type %s not found in manifest %s", execCtx.NodeType, pluginID)
	}

	// 4. Find the operation
	operation := cfg.Operation
	if operation == "" {
		return nil, fmt.Errorf("plugin: operation is required for node %s", execCtx.NodeID)
	}
	actionDef, ok := nodeManifest.Operations[operation]
	if !ok {
		return nil, fmt.Errorf("plugin: operation %s not found in node type %s", operation, execCtx.NodeType)
	}

	// 5. Decrypt credential (if present)
	credentialData := make(map[string]interface{})
	if cfg.CredentialID != "" {
		credentialData, err = e.vaultService.DecryptCredential(ctx, cfg.CredentialID)
		if err != nil {
			return nil, fmt.Errorf("plugin: failed to decrypt credential %s: %w", cfg.CredentialID, err)
		}
	}

	// 6. Resolve FieldValues in the raw config (state, event, input, nodeOutput → actual values)
	resolvedConfig, err := e.resolveConfigFieldValues(execCtx, cfg.RawConfig)
	if err != nil {
		return nil, fmt.Errorf("plugin: failed to resolve config in node %s: %w", execCtx.NodeID, err)
	}

	// 7. Build template contexts — all sources available for {{context.field}} resolution
	contexts := map[string]interface{}{
		"manifest": map[string]interface{}{
			"defaults": map[string]interface{}{
				"baseUrl": manifest.Defaults.BaseUrl,
			},
		},
		"credentials": credentialData,
		"config":      resolvedConfig,
		"wf": map[string]interface{}{
			"state": execCtx.State,
			"input": execCtx.ExternalInputs,
		},
		"event": execCtx.EventPayload,
	}
	if manifest.Defaults.Timeout != nil {
		contexts["manifest"].(map[string]interface{})["defaults"].(map[string]interface{})["timeout"] = *manifest.Defaults.Timeout
	}

	// 8. Resolve all {{context.field}} templates in the action
	resolvedAction := templatereplace.Resolve(actionToMap(&actionDef), contexts)

	// 9. Resolve templates in hooks
	var resolvedHooks interface{}
	if nodeManifest.Hooks != nil {
		resolvedHooks = templatereplace.Resolve(hooksToMap(nodeManifest.Hooks), contexts)
	}

	// 10. Suspend with ready payload
	expiresAt := CalculateExpiresAt(execCtx.Timeout, 30*time.Second)

	nodeState := map[string]interface{}{
		"waitType":     "callback",
		"pluginId":     pluginID,
		"nodeType":     execCtx.NodeType,
		"operation":    operation,
		"action":       resolvedAction,
		"expiresAt":    expiresAt,
		"enableOutput": IsEnableOutput(execCtx.Timeout),
	}
	if resolvedHooks != nil {
		nodeState["hooks"] = resolvedHooks
	}

	return &entities.NodeExecutionResult{
		OutputHandles: nil, // async — no output handles until resume
		NodeState:     nodeState,
	}, nil
}

// extractPluginID extracts the plugin identifier from a node type string.
// "telegram/message" → "telegram"
func extractPluginID(nodeType string) string {
	parts := strings.SplitN(nodeType, "/", 2)
	if len(parts) < 2 {
		return ""
	}
	return parts[0]
}

// resolveConfigFieldValues resolves FieldSourceValue objects in the raw config to actual values.
// Fields that are FieldSourceValue objects (e.g., {type: "state", value: "counter"})
// are resolved using the ValueResolver. Plain values pass through unchanged.
func (e *PluginExecutor) resolveConfigFieldValues(execCtx *entities.NodeExecutionContext, rawConfig map[string]interface{}) (map[string]interface{}, error) {
	resolved := make(map[string]interface{}, len(rawConfig))
	for key, val := range rawConfig {
		v, err := e.resolveFieldValue(execCtx, val)
		if err != nil {
			return nil, fmt.Errorf("field %q: %w", key, err)
		}
		resolved[key] = v
	}
	return resolved, nil
}

// resolveFieldValue resolves a single value — if it's a FieldSourceValue map, resolves it.
// Uses model.MapGetString to handle both map[string]interface{} and bson.M transparently.
func (e *PluginExecutor) resolveFieldValue(execCtx *entities.NodeExecutionContext, val interface{}) (interface{}, error) {
	// Try to convert val to map[string]interface{} (handles bson.M via model helper)
	m := model.ToMap(val)
	if m == nil {
		return val, nil
	}

	// Check if it's a FieldSourceValue (has "type" and "value" keys)
	fieldType := model.MapGetString(m, "type")
	fieldValue := model.MapGetString(m, "value")
	if fieldType == "" || fieldValue == "" {
		return val, nil
	}

	// Skip non-resolvable types — fetchOptions is effectively a literal at runtime
	if fieldType == "fetchOptions" || fieldType == "loadOptions" {
		fieldType = "literal"
	}

	fv := defPorts.FieldValue{
		Type:  defPorts.FieldValueType(fieldType),
		Value: fieldValue,
	}
	if nodeId := model.MapGetString(m, "nodeId"); nodeId != "" {
		fv.NodeID = nodeId
	}

	result, err := e.resolver.Resolve(&fv, execCtx.EventPayload, execCtx.State, execCtx.NodeOutputs, execCtx.ExternalInputs)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// actionToMap converts an ActionDef struct to map[string]interface{} for template resolution.
func actionToMap(action *pluginPorts.ActionDef) map[string]interface{} {
	result := map[string]interface{}{
		"type": action.Type,
	}

	if action.Http != nil {
		httpMap := map[string]interface{}{
			"method": action.Http.Method,
			"path":   action.Http.Path,
		}
		if action.Http.Headers != nil {
			headers := make(map[string]interface{}, len(action.Http.Headers))
			for k, v := range action.Http.Headers {
				headers[k] = v
			}
			httpMap["headers"] = headers
		}
		if action.Http.Body != nil {
			httpMap["body"] = action.Http.Body
		}
		if action.Http.Timeout != nil {
			httpMap["timeout"] = *action.Http.Timeout
		}
		result["http"] = httpMap
	}

	if action.Output != nil {
		result["output"] = map[string]interface{}{
			"dataPath":  action.Output.DataPath,
			"transform": action.Output.Transform,
		}
	}

	return result
}

// hooksToMap converts NodeHooks struct to map[string]interface{} for template resolution.
func hooksToMap(hooks *pluginPorts.NodeHooks) map[string]interface{} {
	result := make(map[string]interface{})
	if hooks.Before != nil {
		result["before"] = actionToMap(hooks.Before)
	}
	if hooks.After != nil {
		result["after"] = actionToMap(hooks.After)
	}
	if hooks.Destroy != nil {
		result["destroy"] = actionToMap(hooks.Destroy)
	}
	return result
}

import type { PluginNodeType, HandleDefinition, HandleOverrides, ResolvedHandles, NodeTimeoutConfig } from '@src/components/workflow/interfaces';

/** Timeout output handle injected when enableOutput is true */
const TIMEOUT_HANDLE: HandleDefinition = {
  id: 'timeout',
  label: 'Timeout',
  position: 'bottom',
  color: '#ff9800',
};

/**
 * Resolve dynamic handles for a node based on its type definition and config.
 * Calls resolveInputs/resolveOutputs if defined, otherwise returns static handles.
 * Applies user-defined label overrides from config.__handleOverrides.
 * Injects "timeout" output handle when node.timeout.enableOutput is true.
 *
 * @param {PluginNodeType} nodeType - Plugin node type definition
 * @param {Record<string, unknown>} config - Current node config
 * @param {NodeTimeoutConfig} [timeout] - Node-level timeout configuration
 * @returns {ResolvedHandles} Resolved input and output handles
 */
export function resolveNodeHandles(
  nodeType: PluginNodeType,
  config: Record<string, unknown>,
  timeout?: NodeTimeoutConfig,
): ResolvedHandles {
  const inputs = nodeType.resolveInputs
    ? nodeType.resolveInputs(config, nodeType.inputs)
    : nodeType.inputs;

  let outputs = nodeType.resolveOutputs
    ? nodeType.resolveOutputs(config, nodeType.outputs)
    : nodeType.outputs;

  // Inject timeout output handle when enableOutput is true
  if (timeout?.enableOutput && !outputs.some(h => h.id === 'timeout')) {
    outputs = [...outputs, TIMEOUT_HANDLE];
  }

  const overrides = (config.__handleOverrides as HandleOverrides) || {};

  return {
    inputs: applyOverrides(inputs, overrides),
    outputs: applyOverrides(outputs, overrides),
  };
}

/**
 * Apply user-defined label and position overrides to a list of handle definitions
 *
 * @param {HandleDefinition[]} handles - Original handle definitions
 * @param {HandleOverrides} overrides - User-defined overrides
 * @returns {HandleDefinition[]} Handles with overrides applied
 */
function applyOverrides(
  handles: HandleDefinition[],
  overrides: HandleOverrides,
): HandleDefinition[] {
  if (!Object.keys(overrides).length) return handles;

  return handles.map(handle => {
    const override = overrides[handle.id];
    if (!override) return handle;

    return {
      ...handle,
      ...(override.label ? { label: override.label } : {}),
      ...(override.position ? { position: override.position } : {}),
    };
  });
}

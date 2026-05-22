import type {
  WorkflowPlugin,
  PluginNodeType,
  PluginMetadata,
  PluginCredentialDefinition,
  NodePropertyDefinition,
  OperationDefinition,
  HandleDefinition,
} from '@src/components/workflow/interfaces';
import type { PluginResponse } from '@mapexos/schemas';
import { markRaw } from 'vue';
import { PLUGIN_CDN_BASE_URL } from '../components/PluginsTab/constants';
import { apis } from '@services/mapex';

/** Raw handle shape from Zod — wider than HandleDefinition (extra fields, `color?: string | undefined`) */
type HandleInput = { id: string; label: string; position: string; color?: string };

// ────────────────────────────────────────────────────────────────────────────
// Converter (Pure Function)
// ────────────────────────────────────────────────────────────────────────────

/**
 * Convert a PluginResponse from @mapexos/schemas into a WorkflowPlugin
 * ready for registration in pluginRegistryStore.
 *
 * This is a pure function — no fetch, no side effects.
 * Handles the shape differences between API response and internal WorkflowPlugin.
 *
 * Uses explicit casts where Zod-inferred types (wide strings, `string | undefined`)
 * diverge from the narrower local interfaces under `exactOptionalPropertyTypes`.
 *
 * @param {PluginResponse} manifest - Plugin response from API
 * @returns {WorkflowPlugin} Plugin object ready for pluginRegistryStore.registerPlugin()
 */
export function convertManifestToPlugin(manifest: PluginResponse): WorkflowPlugin {
  const nodeTypes: PluginNodeType[] = (manifest.nodeTypes ?? []).map((nt) => {
    const nodeType: PluginNodeType = {
      type: nt.type,
      label: nt.label,
      icon: nt.icon,
      color: nt.color,
      description: nt.description ?? '',
      inputs: adaptHandles(nt.inputs as HandleInput[]),
      outputs: adaptHandles(nt.outputs as HandleInput[]),
      configSchema: {},
      defaults: {},
    };

    // Zod schemas infer wide types (e.g. `type: string` vs `NodePropertyType`)
    // Cast to local interfaces which use narrower unions
    if (nt.properties !== undefined) nodeType.properties = nt.properties as unknown as NodePropertyDefinition[];
    if (nt.operations !== undefined) nodeType.operations = nt.operations as unknown as Record<string, OperationDefinition>;
    if (nt.availableOutputs !== undefined) nodeType.availableOutputs = nt.availableOutputs;
    if (nt.timeout !== undefined) nodeType.timeout = nt.timeout;

    return markRaw(nodeType);
  });

  const plugin: WorkflowPlugin = {
    id: manifest.pluginId,
    name: manifest.name,
    version: manifest.version,
    category: (manifest.category) ?? 'integrations',
    icon: manifest.icon,
    nodeTypes,
  };

  // Cast Zod-inferred types to narrower local interfaces (exactOptionalPropertyTypes)
  if (manifest.defaults != null) plugin.defaults = manifest.defaults as NonNullable<WorkflowPlugin['defaults']>;
  if (manifest.metadata !== undefined) plugin.metadata = manifest.metadata as unknown as PluginMetadata;
  if (manifest.credentials !== undefined) plugin.credentials = manifest.credentials as unknown as PluginCredentialDefinition[];
  return plugin;
}

// ────────────────────────────────────────────────────────────────────────────
// Fetch + Convert (for CDN / URL-based loading)
// ────────────────────────────────────────────────────────────────────────────

/**
 * Fetch a plugin manifest from CDN and convert to a WorkflowPlugin
 * ready for registration in the pluginRegistryStore.
 *
 * For direct API responses (already parsed JSON), use convertManifestToPlugin() instead.
 *
 * @param {string} manifestUrl - Relative manifest URL (e.g., 'telegram/manifest.json')
 * @returns {Promise<WorkflowPlugin>} Converted plugin object
 */
export async function loadManifest(manifestUrl: string): Promise<WorkflowPlugin> {
  const fullUrl = `${PLUGIN_CDN_BASE_URL}/${manifestUrl}`;
  const response = await fetch(fullUrl);

  if (!response.ok) {
    throw new Error(`Failed to fetch manifest: ${response.status} ${response.statusText}`);
  }

  const manifest: PluginResponse = await response.json();
  return convertManifestToPlugin(manifest);
}

// ────────────────────────────────────────────────────────────────────────────
// Boot Sequence
// ────────────────────────────────────────────────────────────────────────────

/**
 * Fetch enabled marketplace plugins from the API and register them.
 * Falls back silently to core-only if the API is unreachable.
 *
 * @param {(plugin: WorkflowPlugin) => void} registerPlugin - Registration callback from pluginRegistryStore
 * @returns {Promise<WorkflowPlugin[]>} Successfully registered marketplace plugins
 */
export async function bootMarketplacePlugins(
  registerPlugin: (plugin: WorkflowPlugin) => void,
): Promise<WorkflowPlugin[]> {
  let manifests: PluginResponse[];

  try {
    manifests = await apis.workflows.plugin.getEnabled();
  } catch {
    console.warn('[bootMarketplacePlugins] Failed to reach plugin API — using core-only');
    return [];
  }

  const registered: WorkflowPlugin[] = [];

  for (const manifest of manifests) {
    try {
      const plugin = convertManifestToPlugin(manifest);
      registerPlugin(plugin);
      registered.push(plugin);
    } catch (error) {
      console.error(
        `[bootMarketplacePlugins] Failed to convert plugin "${manifest.pluginId}":`,
        error,
      );
    }
  }

  return registered;
}

// ────────────────────────────────────────────────────────────────────────────
// Private helpers
// ────────────────────────────────────────────────────────────────────────────

/**
 * Adapt raw handle arrays to typed HandleDefinition[].
 *
 * @param {HandleInput[]} handles - Raw handles from API
 * @returns {HandleDefinition[]} Typed handle definitions
 */
function adaptHandles(handles: HandleInput[]): HandleDefinition[] {
  return handles.map((h) => {
    const handle: HandleDefinition = {
      id: h.id,
      label: h.label,
      position: h.position as 'top' | 'bottom' | 'left' | 'right',
    };

    if (h.color !== undefined) handle.color = h.color;

    return handle;
  });
}

import type {
  WorkflowPlugin,
  PluginNodeType,
  PluginMetadata,
  PluginCredentialDefinition,
  NodePropertyDefinition,
  OperationDefinition,
} from '@src/components/workflow/interfaces';

// ────────────────────────────────────────────────────────────────────────────
// Types
// ────────────────────────────────────────────────────────────────────────────

/**
 * Registry entry — lightweight metadata for listing plugins in the marketplace UI.
 * Fetched from `registry.json` on the CDN.
 */
export interface RegistryEntry {
  /** Plugin ID (e.g., 'telegram') */
  id: string;

  /** Display name */
  name: string;

  /** Plugin version (semver) */
  version: string;

  /** Category for grouping (e.g., 'messaging', 'ai') */
  category: string;

  /** Material Icons name */
  icon: string;

  /** Brand SVG icon path relative to CDN (e.g., 'telegram/icon.svg') */
  brandIcon?: string;

  /** Brand color hex */
  color?: string;

  /** Short description */
  description: string;

  /** Author name */
  author?: string;

  /** Searchable tags */
  tags?: string[];

  /** Relative path to full manifest (e.g., 'telegram/manifest.json') */
  manifestUrl: string;

  /** Docs URL */
  docsUrl?: string;

  /** Whether the plugin requires credentials */
  requiresCredentials?: boolean;

  /** Number of node types in this plugin */
  nodeCount?: number;
}

/**
 * Registry response — the top-level JSON served by the CDN.
 */
export interface PluginRegistry {
  /** Schema version identifier */
  $schema: string;

  /** Registry version */
  version: string;

  /** Last update timestamp */
  updatedAt?: string;

  /** Available plugins */
  plugins: RegistryEntry[];
}

/**
 * Options for ManifestLoader
 */
export interface ManifestLoaderOptions {
  /** Base URL of the plugin CDN (e.g., 'http://localhost:3099/plugins') */
  baseUrl: string;

  /** Request timeout in ms (default: 10000) */
  timeout?: number;
}

// ────────────────────────────────────────────────────────────────────────────
// Required fields validation
// ────────────────────────────────────────────────────────────────────────────

const REQUIRED_PLUGIN_FIELDS = ['name', 'version', 'nodeTypes'] as const;
const REQUIRED_NODE_FIELDS = ['type', 'label', 'icon', 'color', 'inputs', 'outputs'] as const;

// ────────────────────────────────────────────────────────────────────────────
// ManifestLoader
// ────────────────────────────────────────────────────────────────────────────

/**
 * Loads plugin manifests from a remote CDN and converts them
 * into WorkflowPlugin objects that can be registered in the plugin registry.
 *
 * JSON manifests don't contain Vue components or functions — the ManifestLoader
 * leaves `canvasComponent` and `configComponent` as undefined so the registry
 * falls back to GenericWorkflowNode and DynamicNodeForm respectively.
 */
export class ManifestLoader {
  private readonly baseUrl: string;
  private readonly timeout: number;

  /** In-memory cache: plugin ID → loaded WorkflowPlugin */
  private readonly cache = new Map<string, WorkflowPlugin>();

  constructor(options: ManifestLoaderOptions) {
    // Remove trailing slash
    this.baseUrl = options.baseUrl.replace(/\/+$/, '');
    this.timeout = options.timeout ?? 10_000;
  }

  /**
   * Fetch the plugin registry (list of all available plugins).
   *
   * @returns {Promise<PluginRegistry>} Registry with plugin entries
   * @throws {Error} If fetch fails or response is invalid
   */
  async fetchRegistry(): Promise<PluginRegistry> {
    const url = `${this.baseUrl}/registry.json`;
    const data = await this.fetchJson<PluginRegistry>(url);

    if (!Array.isArray(data.plugins)) {
      throw new Error(`[ManifestLoader] Invalid registry: "plugins" must be an array`);
    }

    return data;
  }

  /**
   * Load a single plugin manifest and convert it to a registrable WorkflowPlugin.
   * Results are cached by plugin ID — subsequent calls return the cached version.
   *
   * @param {RegistryEntry} entry - Registry entry with manifestUrl
   * @returns {Promise<WorkflowPlugin>} Plugin ready for registry.registerPlugin()
   * @throws {Error} If fetch fails or manifest is invalid
   */
  async loadManifest(entry: RegistryEntry): Promise<WorkflowPlugin> {
    // Return cached if available
    if (this.cache.has(entry.id)) {
      return this.cache.get(entry.id)!;
    }

    const url = `${this.baseUrl}/${entry.manifestUrl}`;
    const json = await this.fetchJson<Record<string, unknown>>(url);

    const plugin = this.jsonToPlugin(json);
    this.cache.set(plugin.id, plugin);

    return plugin;
  }

  /**
   * Convenience: fetch registry + load all manifests in parallel.
   *
   * @returns {Promise<WorkflowPlugin[]>} All plugins ready for registration
   */
  async loadAll(): Promise<WorkflowPlugin[]> {
    const registry = await this.fetchRegistry();
    const results = await Promise.allSettled(
      registry.plugins.map((entry) => this.loadManifest(entry)),
    );

    const plugins: WorkflowPlugin[] = [];
    for (const result of results) {
      if (result.status === 'fulfilled') {
        plugins.push(result.value);
      } else {
        console.error('[ManifestLoader] Failed to load plugin:', result.reason);
      }
    }

    return plugins;
  }

  /**
   * Clear the in-memory cache (e.g., when switching workspaces).
   */
  clearCache(): void {
    this.cache.clear();
  }

  // ──────────────────────────────────────────────────────────────────────────
  // Private
  // ──────────────────────────────────────────────────────────────────────────

  /**
   * Fetch JSON with timeout and error handling.
   *
   * @param {string} url - URL to fetch
   * @returns {Promise<T>} Parsed JSON
   */
  private async fetchJson<T>(url: string): Promise<T> {
    const controller = new AbortController();
    const timer = setTimeout(() => controller.abort(), this.timeout);

    try {
      const response = await fetch(url, { signal: controller.signal });

      if (!response.ok) {
        throw new Error(`HTTP ${response.status}: ${response.statusText}`);
      }

      return (await response.json()) as T;
    } catch (error) {
      if (error instanceof DOMException && error.name === 'AbortError') {
        throw new Error(`[ManifestLoader] Timeout fetching ${url} (${this.timeout}ms)`);
      }
      throw new Error(`[ManifestLoader] Failed to fetch ${url}: ${String(error)}`);
    } finally {
      clearTimeout(timer);
    }
  }

  /**
   * Convert raw JSON into a WorkflowPlugin object.
   * Validates required fields and adapts nodeTypes for registry compatibility.
   *
   * @param {Record<string, unknown>} json - Raw manifest JSON
   * @returns {WorkflowPlugin} Validated plugin object
   */
  private jsonToPlugin(json: Record<string, unknown>): WorkflowPlugin {
    // Validate required plugin fields
    for (const field of REQUIRED_PLUGIN_FIELDS) {
      if (!(field in json)) {
        throw new Error(`[ManifestLoader] Missing required field "${field}" in manifest`);
      }
    }

    // Support both `pluginId` (API shape) and `id` (legacy CDN shape)
    const pluginId = (json.pluginId as string | undefined) ?? (json.id as string | undefined);
    if (!pluginId) {
      throw new Error(`[ManifestLoader] Missing required field "pluginId" (or "id") in manifest`);
    }

    const rawNodeTypes = json.nodeTypes as Record<string, unknown>[];
    if (!Array.isArray(rawNodeTypes) || rawNodeTypes.length === 0) {
      throw new Error(`[ManifestLoader] "nodeTypes" must be a non-empty array`);
    }

    // Adapt each nodeType
    const nodeTypes: PluginNodeType[] = rawNodeTypes.map((raw) => this.adaptNodeType(raw));

    const plugin: WorkflowPlugin = {
      id: pluginId,
      name: json.name as string,
      version: json.version as string,
      category: (json.category as string) ?? 'custom',
      icon: (json.icon as string) ?? 'extension',
      nodeTypes,
    };

    // Only assign marketplace-specific fields when present (exactOptionalPropertyTypes)
    if (json.defaults !== undefined) plugin.defaults = json.defaults as { baseUrl?: string; timeout?: number };
    if (json.metadata !== undefined) plugin.metadata = json.metadata as PluginMetadata;
    if (json.credentials !== undefined) plugin.credentials = json.credentials as PluginCredentialDefinition[];
    return plugin;
  }

  /**
   * Adapt a raw nodeType JSON object into a PluginNodeType.
   * Leaves canvasComponent and configComponent as undefined
   * so the registry falls back to GenericWorkflowNode and DynamicNodeForm.
   *
   * @param {Record<string, unknown>} raw - Raw nodeType from JSON
   * @returns {PluginNodeType} Adapted node type
   */
  private adaptNodeType(raw: Record<string, unknown>): PluginNodeType {
    for (const field of REQUIRED_NODE_FIELDS) {
      if (!(field in raw)) {
        throw new Error(
          `[ManifestLoader] Missing required field "${field}" in nodeType "${typeof raw.type === 'string' ? raw.type : 'unknown'}"`,
        );
      }
    }

    const nodeType: PluginNodeType = {
      type: raw.type as string,
      label: raw.label as string,
      icon: raw.icon as string,
      color: raw.color as string,
      description: (raw.description as string) ?? '',
      inputs: raw.inputs as PluginNodeType['inputs'],
      outputs: raw.outputs as PluginNodeType['outputs'],
      configSchema: (raw.configSchema as Record<string, unknown>) ?? {},
    };

    // Only assign optional fields when present (exactOptionalPropertyTypes)
    if (raw.properties !== undefined) nodeType.properties = raw.properties as NodePropertyDefinition[];
    if (raw.defaults !== undefined) nodeType.defaults = raw.defaults as Record<string, unknown>;
    if (raw.operations !== undefined) nodeType.operations = raw.operations as Record<string, OperationDefinition>;
    if (raw.availableOutputs !== undefined) nodeType.availableOutputs = raw.availableOutputs as Array<{ path: string; description: string }>;
    if (raw.timeout !== undefined) nodeType.timeout = raw.timeout as { duration: number; unit: string; enableOutput: boolean };

    // canvasComponent: undefined  → Registry uses GenericWorkflowNode
    // configComponent: undefined  → Registry uses DynamicNodeForm (via properties[])
    // validate: undefined         → No custom validation

    return nodeType;
  }
}

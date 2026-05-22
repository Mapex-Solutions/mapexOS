import type { PluginCategory } from '@src/components/workflow/interfaces';

/**
 * Single entry from the plugin marketplace API (`GET /api/v1/plugins`).
 * Represents a plugin available for installation/enabling.
 */
export interface RegistryEntry {
  /** Unique plugin identifier (maps from API `pluginId`) */
  id: string;

  /** Display name */
  name: string;

  /** Semver version */
  version: string;

  /** Plugin category for grouping */
  category: PluginCategory;

  /** Material icon name (fallback) */
  icon: string;

  /** Relative path to brand SVG icon on CDN */
  brandIcon: string;

  /** Brand color (hex) */
  color: string;

  /** Short description for marketplace card */
  description: string;

  /** Author name */
  author: string;

  /** Search/filter tags */
  tags: string[];

  /** Relative path to full manifest.json (CDN) or null if API-only */
  manifestUrl: string;

  /** Link to plugin documentation */
  docsUrl: string;

  /** Whether the plugin requires credential configuration */
  requiresCredentials: boolean;

  /** Number of action node types in this plugin */
  nodeCount: number;

  /** Number of trigger types in this plugin */
  triggerCount: number;

  /** Whether the plugin is enabled for this organization */
  enabled: boolean;

  /** Whether this is a system plugin (cannot be disabled) */
  isSystem: boolean;
}

/**
 * Installed plugin summary used in the Installed sub-tab
 */
export interface InstalledPlugin {
  /** Plugin ID (e.g., 'telegram') */
  id: string;

  /** MongoDB ObjectID from the API (used for DELETE/UPDATE calls). Undefined for core plugins. */
  mongoId?: string;

  /** Display name */
  name: string;

  /** Semver version */
  version: string;

  /** Material icon name */
  icon: string;

  /** Relative path to brand SVG icon on CDN (from metadata.brandIcon) */
  brandIcon: string;

  /** Brand color (hex) */
  color: string;

  /** Short description */
  description: string;

  /** Author name */
  author: string;

  /** Search/filter tags */
  tags: string[];

  /** Plugin category */
  category: PluginCategory;

  /** Number of node types */
  nodeCount: number;

  /** Whether this is a core (built-in) plugin */
  isCore: boolean;
}

/**
 * Group of installed plugins under the same category
 */
export interface InstalledPluginGroup {
  /** Category key */
  category: PluginCategory;

  /** Translated category label */
  label: string;

  /** Plugins in this category */
  plugins: InstalledPlugin[];
}

/**
 * State for the plugin detail dialog — populated from either a RegistryEntry or a core WorkflowPlugin
 */
export interface PluginDetailState {
  /** Plugin ID */
  id: string;

  /** Display name */
  name: string;

  /** Author name */
  author: string;

  /** Semver version */
  version: string;

  /** Short description */
  description: string;

  /** Full URL to brand icon (or empty for icon-only fallback) */
  brandIconUrl: string;

  /** Material icon name — fallback when brandIconUrl is empty */
  icon: string;

  /** Brand color (hex) */
  color: string;

  /** Plugin category */
  category: PluginCategory;

  /** Search/filter tags */
  tags: string[];

  /** Whether this is a core plugin */
  isCore: boolean;

  /** Manifest URL for marketplace plugins (null for core) */
  manifestUrl: string | null;
}

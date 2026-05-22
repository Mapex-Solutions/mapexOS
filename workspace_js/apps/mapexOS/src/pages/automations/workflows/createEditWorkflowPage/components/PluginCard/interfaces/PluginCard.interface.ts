/**
 * Props for the PluginCard component
 */
export interface PluginCardProps {
  /** Plugin display name */
  name: string;

  /** Short description */
  description: string;

  /** Author name */
  author: string;

  /** Semver version string */
  version: string;

  /** Full URL to the brand icon (empty string = use icon fallback) */
  brandIconUrl: string;

  /** Material icon name — fallback when brandIconUrl is empty or fails to load */
  icon: string;

  /** Brand color (hex) for category badge */
  color: string;

  /** Translated category label */
  categoryLabel: string;

  /** Search/filter tags */
  tags: string[];

  /** Number of node types */
  nodeCount: number;

  /** Whether the plugin is already installed */
  installed: boolean;

  /** Whether an install operation is in progress */
  installing: boolean;

  /** Whether any install operation is running (disables all install buttons) */
  installDisabled: boolean;

  /** Label for the install button */
  installLabel: string;

  /** Label for the installing state */
  installingLabel: string;

  /** Label for the installed badge */
  installedLabel: string;

  /** Label for the details button */
  detailsLabel: string;
}

/**
 * Emits for the PluginCard component
 */
export interface PluginCardEmits {
  (e: 'install'): void;
  (e: 'details'): void;
}

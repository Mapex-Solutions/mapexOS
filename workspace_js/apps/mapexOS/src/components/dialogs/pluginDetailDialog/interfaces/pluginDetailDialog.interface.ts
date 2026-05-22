/**
 * Summary of a plugin node type for display in the detail dialog
 */
export interface PluginNodeTypeSummary {
  /** Fully qualified type (e.g., telegram/sendMessage) */
  type: string;

  /** Display label */
  label: string;

  /** Material icon name */
  icon: string;

  /** Brand color */
  color: string;

  /** Node description */
  description: string;

  /** Number of input handles */
  inputCount: number;

  /** Number of output handles */
  outputCount: number;
}

/**
 * Props for the PluginDetailDialog component
 */
export interface PluginDetailDialogProps {
  /** Whether the dialog is visible */
  modelValue: boolean;

  /** Plugin display name */
  name: string;

  /** Plugin author */
  author: string;

  /** Semver version */
  version: string;

  /** Short description */
  description: string;

  /** Full URL to the brand icon (empty string for icon-only fallback) */
  brandIconUrl: string;

  /** Material icon name — used as fallback when brandIconUrl is empty */
  icon: string;

  /** Brand color (hex) */
  color: string;

  /** Translated category label */
  categoryLabel: string;

  /** Search/filter tags */
  tags: string[];

  /** Whether the manifest is currently loading */
  loading: boolean;

  /** Node type summaries (derived from manifest by the parent) */
  nodeTypes: PluginNodeTypeSummary[];

  /** Whether the plugin is already installed */
  installed: boolean;

  /** Whether install is in progress for this plugin */
  installing: boolean;

  /** Whether any install is running globally */
  installDisabled: boolean;

  /** Label for the install button */
  installLabel: string;

  /** Label for the installing state */
  installingLabel: string;

  /** Label for the installed badge */
  installedLabel: string;

  /** Label for the node types section header */
  nodeTypesLabel: string;

  /** Label for loading state */
  loadingLabel: string;

  /** Label for input count badge */
  inputsLabel: string;

  /** Label for output count badge */
  outputsLabel: string;
}

/**
 * Emits for the PluginDetailDialog component
 */
export interface PluginDetailDialogEmits {
  (e: 'update:modelValue', value: boolean): void;
  (e: 'install'): void;
}

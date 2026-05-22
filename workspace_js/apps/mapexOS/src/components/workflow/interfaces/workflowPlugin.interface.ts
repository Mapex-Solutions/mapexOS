import type { Component } from 'vue';
import type { SourceType } from './fieldSource.interface';

// ────────────────────────────────────────────────────────────────────────────
// Declarative Node Property System
// ────────────────────────────────────────────────────────────────────────────

/**
 * Supported property types for declarative node forms.
 *
 * Core types (used by core plugins):
 * - string: text input
 * - number: numeric input
 * - boolean: toggle switch
 * - options: single-select dropdown
 * - json: JSON text editor
 *
 * Extended types (used by marketplace plugins):
 * - fieldSource: value with source type selector (event, state, literal, nodeOutput, etc.)
 * - multiOptions: multi-select dropdown
 * - collection: repeatable group of fields (add/remove rows)
 * - fixedCollection: named group of fields (collapsible section)
 * - dateTime: date/time picker
 * - hidden: stored in config but not rendered
 * - notice: read-only informational text
 */
export type NodePropertyType =
  | 'string'
  | 'number'
  | 'boolean'
  | 'options'
  | 'json'
  | 'fieldSource'
  | 'multiOptions'
  | 'collection'
  | 'fixedCollection'
  | 'dateTime'
  | 'hidden'
  | 'notice';

/**
 * Visual rendering options for a property.
 * Only affects how the field looks in the UI, not data.
 */
export interface PropertyRendering {
  /** Render input as textarea (for string type) */
  multiline?: boolean;

  /** Number of rows for textarea (requires multiline: true) */
  rows?: number;

  /** Render input as password field (masked) */
  password?: boolean;

  /** Editor language for code editor (e.g., 'json', 'javascript') */
  editor?: string;

  /** Allow multiple values (turns single input into array) */
  multipleValues?: boolean;

  /** Minimum value (for number type) */
  min?: number;

  /** Maximum value (for number type) */
  max?: number;

  /** Placeholder text for the input */
  placeholder?: string;

  /** Show only date without time (for dateTime type) */
  dateOnly?: boolean;
}

/**
 * Rule for selecting which fetchOptions entry to use based on form state.
 * The first matching rule wins. Use `when: {}` as fallback (always matches).
 */
export interface FetchOptionsRule {
  /** Condition based on other fields. {} = always matches (fallback) */
  when: Record<string, unknown[]>;

  /** References manifest.fetchOptions[key] */
  key: string;

  /** Display label in the source type selector */
  label: string;
}

/**
 * Configuration for fetching dynamic dropdown options from an external API.
 * Used by fieldSource properties to populate dropdowns at design-time.
 */
export interface PropertyFetchOptions {
  /** Rules to determine which fetchOptions entry to use */
  rules: FetchOptionsRule[];

  /** Reload when these sibling fields change. Values sent as dependsOn to backend */
  dependsOn?: string[];
}

/**
 * Declarative property definition for auto-generated node config forms.
 *
 * Core plugins use: name, displayName, type, default, hint, required, options, displayOptions.
 * Marketplace plugins additionally use: allowedSources, isSecret, rendering, fetchOptions.
 */
export interface NodePropertyDefinition {
  /** Property key (maps to config[name]) */
  name: string;

  /** Display label shown in the form */
  displayName: string;

  /** Property type — determines which form control is rendered */
  type: NodePropertyType;

  /** Default value when creating a new node */
  default: unknown;

  /** Help text shown below the field */
  hint?: string;

  /** Whether the field is required */
  required?: boolean;

  /** Options for 'options' and 'multiOptions' types */
  options?: { label: string; value: string | number | boolean }[];

  /** Conditional visibility — show only when other fields match */
  displayOptions?: { show?: Record<string, unknown[]> };

  /**
   * Allowed source types for 'fieldSource' type.
   * Controls which tabs appear in FieldSourceSelector.
   * @example ['literal', 'state', 'event', 'nodeOutput', 'fetchOptions']
   */
  allowedSources?: SourceType[];

  /**
   * Marks this field as containing sensitive data (API keys, tokens, passwords).
   * Fields with isSecret: true are rendered as masked inputs.
   */
  isSecret?: boolean;

  /** Visual rendering options — how the field looks */
  rendering?: PropertyRendering;

  /** Dynamic dropdown fetching config — where options come from */
  fetchOptions?: PropertyFetchOptions;

  /**
   * Sub-properties for 'collection' and 'fixedCollection' types.
   * Defines the schema of fields inside the group.
   */
  values?: NodePropertyDefinition[];

  /**
   * Notice type for 'notice' property type.
   * Controls the visual style of the informational text.
   */
  noticeType?: 'info' | 'warning' | 'success';
}

/**
 * Props for DynamicNodeForm component
 */
export interface DynamicNodeFormProps {
  /** Property definitions for the form */
  properties: NodePropertyDefinition[];

  /** Current config values */
  config: Record<string, unknown>;

  /** Node type string — used to rebuild config when operation changes */
  nodeType?: string;
}

/**
 * Emits for DynamicNodeForm component
 */
export interface DynamicNodeFormEmits {
  (e: 'update:config', config: Record<string, unknown>): void;
}

// ────────────────────────────────────────────────────────────────────────────
// Unified Action Contract
// ────────────────────────────────────────────────────────────────────────────

/**
 * HTTP action definition — request template with template variables.
 */
export interface HttpActionDef {
  method: 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE';
  path: string;
  headers?: Record<string, string>;
  body?: unknown;
  timeout?: number;
}

/**
 * Action output definition — how to extract/transform the response.
 */
export interface ActionOutputDef {
  dataPath?: string;
  valuePath?: string;
  labelPath?: string;
  transform?: string;
}

/**
 * Unified action contract. Used by operations, fetchOptions, credential test, and hooks.
 */
export interface ActionDef {
  type: 'http' | 'mqtt' | 'nats' | 'script';
  http?: HttpActionDef;
  output?: ActionOutputDef;
}

/**
 * Operation definition — maps operation value to an Action.
 */
export type OperationDefinition = ActionDef;

// ────────────────────────────────────────────────────────────────────────────
// Plugin System
// ────────────────────────────────────────────────────────────────────────────

/**
 * Plugin categories for catalog grouping.
 *
 * Core categories are used by built-in plugins.
 * Marketplace plugins can use ANY string as category (e.g., 'messaging', 'ai', 'payments').
 * The catalog UI dynamically creates groups for unknown categories.
 */
export type PluginCategory =
  | 'triggers'
  | 'logic'
  | 'state'
  | 'flow_control'
  | 'timers'
  | 'integrations'
  | 'observability'
  | 'annotations'
  | 'custom'
  | (string & {});  // Open union — accepts any string while preserving autocomplete

/**
 * Validation result from node config validation
 */
export interface ValidationResult {
  /** Whether the config is valid */
  valid: boolean;

  /** Error messages if invalid */
  errors: string[];
}

/**
 * Node type definition registered by a plugin.
 * Components are declared as lazy factories for code splitting.
 */
export interface PluginNodeType {
  /** Unique node type (e.g., 'core/condition', 'telegram/message') */
  type: string;

  /** Display label */
  label: string;

  /** Icon for catalog and canvas (Material Icons name) */
  icon: string;

  /** Color theme for the node card */
  color: string;

  /** Short description shown in catalog */
  description: string;

  /** Input handles */
  inputs: HandleDefinition[];

  /** Output handles */
  outputs: HandleDefinition[];

  /** JSON Schema for node configuration */
  configSchema: Record<string, unknown>;

  /** Vue component for canvas rendering (lazy import). Optional — GenericWorkflowNode is used if omitted. */
  canvasComponent?: Component | (() => Promise<{ default: Component }>);

  /** Vue component for config panel (lazy import). Optional — DynamicNodeForm is used when properties[] is defined. */
  configComponent?: Component | (() => Promise<{ default: Component }>);

  /** Vue component for fullscreen editor (lazy import, opens on double-click) */
  fullscreenComponent?: Component | (() => Promise<{ default: Component }>);

  /** Declarative property definitions for auto-generated config forms */
  properties?: NodePropertyDefinition[];

  /** Validation function */
  validate?: (config: Record<string, unknown>) => ValidationResult;

  /** Default config values */
  defaults?: Record<string, unknown>;

  /** Dynamic output handle resolver — called when config changes */
  resolveOutputs?: HandleResolver;

  /** Dynamic input handle resolver — called when config changes */
  resolveInputs?: HandleResolver;

  /** Whether the node can be deleted (default: true) */
  deletable?: boolean;

  /** Whether the node is hidden from the plugin catalog (default: false) */
  catalogHidden?: boolean;

  /** Whether clicking the node opens the config panel (default: true) */
  configurable?: boolean;

  /** Visual shape on canvas — 'square' (default) or 'circle' */
  shape?: 'square' | 'circle';

  /** Operation-to-action mappings. Keyed by operation value (e.g., 'sendText'). */
  operations?: Record<string, OperationDefinition>;

  /** Fields available in the node's output — shown as hints in FieldSourceSelector for nodeOutput */
  availableOutputs?: Array<{ path: string; description: string }>;

  /** Default async timeout for this node type. Users can override at the node level. */
  timeout?: { duration: number; unit: string; enableOutput: boolean };

  /** Default error handler config. Users can override at the node level. */
  errorHandler?: { enabled: boolean; maxAttempts: number; initialInterval: number; intervalUnit: string; backoffMultiplier: number };

  /** Called when a node of this type is added to the canvas */
  onNodeMount?: (nodeId: string, config: Record<string, unknown>) => void;

  /** Called when a node of this type is removed from the canvas */
  onNodeUnmount?: (nodeId: string) => void;

  /** Plugin ID — injected by the registry during registration. Read-only. */
  _pluginId?: string;
}

/**
 * Plugin registration interface.
 * Every plugin must implement this to register node types.
 */
export interface WorkflowPlugin {
  /** Unique plugin identifier */
  id: string;

  /** Display name */
  name: string;

  /** Plugin version (semver) */
  version: string;

  /** Plugin category for catalog grouping */
  category: PluginCategory;

  /** Icon (Material Icons name) */
  icon: string;

  /** Node types registered by this plugin (action nodes — outbound operations) */
  nodeTypes: PluginNodeType[];

  /** Called when the plugin is activated (registered in the editor) */
  onActivate?: (context: PluginActivationContext) => void;

  /** Called when the plugin is deactivated (unregistered from the editor) */
  onDeactivate?: () => void;

  // ── Marketplace plugin fields (optional — not used by core plugins) ──

  /** Default values inherited by operations and fetchOptions */
  defaults?: { baseUrl?: string; timeout?: number };

  /** Plugin metadata for UI display (brand icon, docs URL, etc.) */
  metadata?: PluginMetadata;

  /** Credential definitions — array of authentication methods */
  credentials?: PluginCredentialDefinition[];
}

// ────────────────────────────────────────────────────────────────────────────
// Marketplace Plugin Types
// ────────────────────────────────────────────────────────────────────────────

/**
 * Plugin metadata for UI display
 */
export interface PluginMetadata {
  /** Brand SVG icon path (relative to CDN, e.g., 'telegram.svg') */
  brandIcon?: string;

  /** Brand color (hex, e.g., '#0088CC') */
  color?: string;

  /** Documentation URL */
  docs?: string;
}

/**
 * Credential definition for a plugin.
 * Describes one authentication method and how to test it.
 * A plugin can have multiple credential definitions (e.g., API Key + OAuth2).
 */
export interface PluginCredentialDefinition {
  /** Credential type ID (e.g., 'telegramApi', 'apiKey', 'oauth2') */
  id: string;

  /** Display name (e.g., 'Telegram Bot API') */
  name: string;

  /** Credential fields (what the user needs to provide) */
  fields: CredentialFieldDefinition[];

  /** Test action to validate credentials — uses unified Action contract */
  test?: ActionDef;
}

/**
 * Individual credential field definition
 */
export interface CredentialFieldDefinition {
  /** Field key */
  name: string;

  /** Display label */
  displayName: string;

  /** Field type */
  type: 'string' | 'number' | 'boolean' | 'options';

  /** Whether this field is required */
  required?: boolean;

  /** Whether this field contains sensitive data (renders as password, encrypted in storage) */
  isSecret?: boolean;

  /** Help text */
  hint?: string;

  /** Default value */
  default?: unknown;

  /** Options for 'options' type */
  options?: { label: string; value: string | number }[];
}

/**
 * @deprecated Use ActionDef instead.
 */
export interface CredentialTestDefinition {
  method: 'GET' | 'POST';
  path: string;
  body?: Record<string, unknown>;
}

// ────────────────────────────────────────────────────────────────────────────
// Catalog
// ────────────────────────────────────────────────────────────────────────────

/**
 * Catalog category group for display
 */
export interface CatalogGroup {
  /** Category key */
  category: PluginCategory;

  /** Display label */
  label: string;

  /** Category icon */
  icon: string;

  /** Color for the category header (optional, used by marketplace categories) */
  color?: string;

  /** Node types in this category */
  nodeTypes: PluginNodeType[];
}

// ────────────────────────────────────────────────────────────────────────────
// Handle System
// ────────────────────────────────────────────────────────────────────────────

/**
 * Handle definition for node inputs/outputs
 */
export interface HandleDefinition {
  /** Handle ID (unique within node) */
  id: string;

  /** Display label */
  label: string;

  /** Handle position */
  position: 'top' | 'bottom' | 'left' | 'right';

  /** Data type for connection validation */
  dataType?: string;

  /** Max connections (null = unlimited) */
  maxConnections?: number | null;

  /** Handle dot color — CSS color value (e.g., '#4caf50'). Falls back to node edge color. */
  color?: string;
}

/**
 * Resolver function that dynamically generates handle definitions
 * based on node configuration.
 *
 * @param {Record<string, unknown>} config - Current node config
 * @param {HandleDefinition[]} staticHandles - Static handle definitions from the node type
 * @returns {HandleDefinition[]} Resolved handle definitions
 */
export type HandleResolver = (
  config: Record<string, unknown>,
  staticHandles: HandleDefinition[],
) => HandleDefinition[];

/**
 * Per-handle overrides stored in node config.
 * Allows users to rename labels and reposition handles.
 */
export interface HandleOverrides {
  [handleId: string]: {
    /** Custom label override */
    label?: string;
    /** Custom position override (top/bottom/left/right) */
    position?: 'top' | 'bottom' | 'left' | 'right';
  };
}

/**
 * Resolved handle result containing both inputs and outputs
 */
export interface ResolvedHandles {
  /** Resolved input handles */
  inputs: HandleDefinition[];

  /** Resolved output handles */
  outputs: HandleDefinition[];
}

// ────────────────────────────────────────────────────────────────────────────
// Lifecycle
// ────────────────────────────────────────────────────────────────────────────

/**
 * Disposable resource — call dispose() to release
 */
export interface Disposable {
  /** Release the resource */
  dispose: () => void;
}

/**
 * Context passed to plugin.onActivate().
 * Follows VS Code's ExtensionContext pattern with a disposal bag.
 */
export interface PluginActivationContext {
  /** Plugin ID */
  pluginId: string;

  /**
   * Disposal bag — push disposables here and they will be auto-disposed
   * when the plugin is deactivated.
   */
  subscriptions: Disposable[];

  /**
   * Register i18n translations for this plugin.
   * Messages are merged under the `wf.{pluginId}` namespace.
   *
   * @param {string} locale - Locale code (e.g., 'en-US', 'pt-BR')
   * @param {Record<string, unknown>} messages - Translation messages object
   */
  registerTranslations: (locale: string, messages: Record<string, unknown>) => void;
}

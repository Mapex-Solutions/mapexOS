/**
 * EventTracerPage Constants
 *
 * Configuration for the enterprise Event Tracer visualization.
 * Organized by the 3 phases of event processing.
 */

// ============================================================================
// Phase Configuration
// ============================================================================

/**
 * Phase identifiers
 */
export type PhaseType =
  | 'ingestion'
  | 'routing'

/**
 * Stage identifiers within phases
 */
export type StageType =
  | 'raw'
  | 'jsExec'
  | 'router'
  | 'directTrigger'
  | 'saveEvent'
  | 'lakeHouse'

/**
 * Phase configuration
 */
export interface PhaseConfig {
  /** Phase identifier */
  phase: PhaseType;
  /** Display label key (for i18n) */
  labelKey: string;
  /** Material icon name */
  icon: string;
  /** Color for the phase header */
  color: string;
  /** Whether this phase is currently implemented */
  implemented: boolean;
}

/**
 * Stage configuration
 */
export interface StageConfig {
  /** Stage identifier */
  stage: StageType;
  /** Display label key (for i18n) */
  labelKey: string;
  /** Material icon name */
  icon: string;
  /** Color for active/success state */
  color: string;
  /** Whether this stage is currently implemented */
  implemented: boolean;
}

// ============================================================================
// Phase Configurations
// ============================================================================

/**
 * Phase configurations in order of execution
 */
export const PHASES: PhaseConfig[] = [
  {
    phase: 'ingestion',
    labelKey: 'phases.ingestion',
    icon: 'input',
    color: 'blue-6',
    implemented: true,
  },
  {
    phase: 'routing',
    labelKey: 'phases.routing',
    icon: 'call_split',
    color: 'orange-6',
    implemented: true,
  },
];

// ============================================================================
// Stage Configurations by Phase
// ============================================================================

/**
 * Ingestion phase stages
 */
export const INGESTION_STAGES: StageConfig[] = [
  {
    stage: 'raw',
    labelKey: 'stages.raw',
    icon: 'terminal',
    color: 'blue-6',
    implemented: true,
  },
  {
    stage: 'jsExec',
    labelKey: 'stages.jsExec',
    icon: 'code',
    color: 'purple-6',
    implemented: true,
  },
];

/**
 * Routing phase stages (fan-out destinations)
 */
export const ROUTING_STAGES: StageConfig[] = [
  {
    stage: 'router',
    labelKey: 'stages.router',
    icon: 'call_split',
    color: 'orange-6',
    implemented: true,
  },
  {
    stage: 'directTrigger',
    labelKey: 'stages.directTrigger',
    icon: 'flash_on',
    color: 'amber-6',
    implemented: true,
  },
  {
    stage: 'saveEvent',
    labelKey: 'stages.saveEvent',
    icon: 'save',
    color: 'green-6',
    implemented: false, // TODO: Implement when save event logging is ready
  },
  {
    stage: 'lakeHouse',
    labelKey: 'stages.lakeHouse',
    icon: 'storage',
    color: 'indigo-6',
    implemented: false, // TODO: Implement when data lake is integrated
  },
];

// ============================================================================
// Trigger Type Icons and Colors
// ============================================================================

/**
 * Trigger type icon mapping
 */
export const TRIGGER_TYPE_ICONS: Record<string, string> = {
  http: 'http',
  mqtt: 'sensors',
  rabbitmq: 'hub',
  nats: 'dns',
  websocket: 'sync_alt',
  email: 'email',
  teams: 'groups',
  slack: 'tag',
};

/**
 * Trigger type color mapping
 */
export const TRIGGER_TYPE_COLORS: Record<string, string> = {
  http: 'blue-6',
  mqtt: 'purple-6',
  rabbitmq: 'orange-6',
  nats: 'cyan-6',
  websocket: 'teal-6',
  email: 'red-6',
  teams: 'indigo-6',
  slack: 'pink-6',
};

/**
 * Trigger category colors
 */
export const TRIGGER_CATEGORY_COLORS: Record<string, string> = {
  technical: 'blue-6',
  communication: 'green-6',
};

// ============================================================================
// Router Kind Icons and Colors
// ============================================================================

/**
 * Router kind icon mapping
 */
export const ROUTER_KIND_ICONS: Record<string, string> = {
  trigger: 'flash_on',
  save_event: 'save',
  lake_house: 'storage',
  notification: 'notifications',
};

/**
 * Router kind color mapping
 */
export const ROUTER_KIND_COLORS: Record<string, string> = {
  trigger: 'amber-6',
  save_event: 'green-6',
  lake_house: 'indigo-6',
  notification: 'pink-6',
};

// ============================================================================
// Status Colors
// ============================================================================

/**
 * Status color mapping
 */
export const STATUS_COLORS = {
  success: 'positive',
  failed: 'negative',
  pending: 'grey-5',
  notImplemented: 'grey-4',
  matched: 'positive',
  notMatched: 'warning',
} as const;

/**
 * Status icon mapping
 */
export const STATUS_ICONS = {
  success: 'check_circle',
  failed: 'error',
  pending: 'schedule',
  notImplemented: 'construction',
  matched: 'check',
  notMatched: 'remove',
} as const;

// ============================================================================
// Source Type Configuration
// ============================================================================

/**
 * Source type icon mapping
 */
export const SOURCE_TYPE_ICONS: Record<string, string> = {
  http_gateway: 'http',
  mqtt_gateway: 'sensors',
  lorawan_gateway: 'cell_tower',
  router: 'call_split',
};

/**
 * Source type color mapping
 */
export const SOURCE_TYPE_COLORS: Record<string, string> = {
  http_gateway: 'blue-6',
  mqtt_gateway: 'purple-6',
  lorawan_gateway: 'orange-6',
  router: 'orange-6',
};

// ============================================================================
// Validation Constants
// ============================================================================

/**
 * UUID validation regex pattern
 * Matches standard UUID format: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
 */
export const UUID_PATTERN = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i;

/**
 * Minimum characters to start searching
 */
export const MIN_SEARCH_LENGTH = 8;

/**
 * Default time window for searching related events (in milliseconds)
 */
export const DEFAULT_TIME_WINDOW_MS = 60000; // 1 minute

/**
 * Maximum items to fetch per stage
 */
export const MAX_ITEMS_PER_STAGE = 50;

// ============================================================================
// Legacy Constants (kept for backward compatibility)
// ============================================================================

/**
 */
export const PIPELINE_STAGES: StageConfig[] = [
  ...INGESTION_STAGES,
  {
    stage: 'router',
    labelKey: 'stages.router',
    icon: 'route',
    color: 'orange-6',
    implemented: true,
  },
  {
    stage: 'directTrigger',
    labelKey: 'stages.trigger',
    icon: 'flash_on',
    color: 'amber-6',
    implemented: true,
  },
];

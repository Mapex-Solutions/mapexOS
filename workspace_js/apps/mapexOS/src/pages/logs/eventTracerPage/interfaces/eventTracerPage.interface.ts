/**
 * EventTracerPage Interfaces
 *
 * Defines the data structures for the enterprise Event Tracer visualization.
 * The trace is organized into 3 phases:
 *   1. INGESTION: Raw Event → JS Executor
 *   2. ROUTING: Router fan-out to multiple destinations (Trigger, Rule Engine, Save Event, Data Lake)
 *   3. RULE ENGINE TRIGGERS: Triggers dispatched by the Rule Engine
 */

import type {
  EventsRawResponse,
  EventsJsExecResponse,
  EventsRouterResponse,

  EventsTriggerResponse,
} from '@mapexos/schemas';

// ============================================================================
// Base Stage Types
// ============================================================================

/**
 * Base stage data shared by all stages
 */
export interface BaseStageData {
  /** Whether this stage has data */
  hasData: boolean;
  /** Whether this stage was successful */
  success: boolean | null;
  /** Created date when this stage was processed */
  created: string | null;
  /** Execution time in milliseconds */
  durationMs: number | null;
  /** Error message (if any) */
  error: string | null;
}

// ============================================================================
// Phase 1: Ingestion Types
// ============================================================================

/**
 * Raw Event stage data
 */
export interface RawEventStageData extends BaseStageData {
  /** Raw event data from API */
  data: EventsRawResponse | null;
  /** Source type (http_gateway, mqtt_gateway, etc.) */
  source: string | null;
  /** Thread ID (data source identifier) */
  threadId: string | null;
}

/**
 * JS Executor stage data
 */
export interface JsExecStageData extends BaseStageData {
  /** JS Executor event data from API */
  data: EventsJsExecResponse | null;
  /** Failed at stage (if any) */
  failedAt: string | null;
}

/**
 * Ingestion phase containing Raw Event and JS Executor stages
 */
export interface IngestionPhase {
  /** Raw event data */
  raw: RawEventStageData;
  /** JS Executor data */
  jsExec: JsExecStageData;
  /** Total duration for ingestion phase */
  totalDurationMs: number;
  /** Whether all stages in this phase succeeded */
  allSuccess: boolean;
}

// ============================================================================
// Phase 2: Routing Types
// ============================================================================

/**
 * Router execution result with fan-out details
 */
export interface RouterStageData extends BaseStageData {
  /** Router event data from API */
  data: EventsRouterResponse | null;
  /** Route Group ID */
  routerId: string | null;
  /** Route Group name */
  name: string | null;
  /** Total number of routers in the group */
  totalRouters: number;
  /** Number of routers that matched conditions */
  matchedCount: number;
  /** Number of routers that published events */
  publishedCount: number;
  /** Parsed routers array from event JSON */
  routers: RouterResult[] | null;
}

/**
 * Individual router result within a Route Group
 */
export interface RouterResult {
  kind: string;
  /** Whether conditions matched */
  matched: boolean;
  /** Whether event was published */
  published: boolean;
  /** Condition evaluation results */
  conditions: ConditionResult[];
}

/**
 * Condition evaluation result
 */
export interface ConditionResult {
  /** Field that was evaluated */
  field: string;
  /** Operator used */
  operator: string;
  /** Expected value */
  expected: any;
  /** Actual value */
  actual: any;
  /** Whether condition passed */
  passed: boolean;
}

/**
 * Direct trigger from Router (not from Rule Engine)
 */
export interface DirectTriggerStageData extends BaseStageData {
  /** Trigger event data from API */
  data: EventsTriggerResponse | null;
  /** Trigger ID */
  triggerId: string | null;
  /** Trigger name */
  triggerName: string | null;
  /** Trigger type (http, mqtt, email, etc.) */
  triggerType: string | null;
  /** Category (technical, communication) */
  category: string | null;
}

/**
 * Save Event stage data
 */
export interface SaveEventStageData extends BaseStageData {
  /** Whether save event was triggered */
  triggered: boolean;
  /** Router result data for this destination */
  data: RouterResult | null;
}

/**
 * Data Lake stage data
 */
export interface LakeHouseStageData extends BaseStageData {
  /** Whether data lake was triggered */
  triggered: boolean;
  /** Router result data for this destination */
  data: RouterResult | null;
}

/**
 * Routing phase containing all fan-out destinations
 */
export interface RoutingPhase {
  /** Router execution data */
  router: RouterStageData;
  /** Direct triggers from router (not from rule engine) */
  directTriggers: DirectTriggerStageData[];
  /** Save event stage */
  saveEvent: SaveEventStageData;
  /** Data lake stage */
  lakeHouse: LakeHouseStageData;
  /** Number of destinations that were routed to */
  destinationsCount: number;
  /** Total duration for routing phase */
  totalDurationMs: number;
  /** Whether all stages in this phase succeeded */
  allSuccess: boolean;
}

// ============================================================================
// Complete Trace Result
// ============================================================================

/**
 * Complete event trace result with all phases
 */
export interface EventTraceResult {
  /** The event tracker ID being traced */
  eventTrackerId: string;
  /** The thread ID (data source identifier) */
  threadId: string;
  /** Created date of the raw event */
  created: string;

  /** Phase 1: Ingestion (Raw → JS Executor) */
  ingestion: IngestionPhase;

  /** Phase 2: Routing (Router → Fan-out destinations) */
  routing: RoutingPhase;

  /** Total execution time across all phases */
  totalExecutionTime: number;

  /** Whether all phases completed successfully */
  allSuccess: boolean;

  /** Summary statistics */
  summary: TraceSummary;
}

/**
 * Trace summary statistics
 */
export interface TraceSummary {
  /** Total stages with data */
  stagesWithData: number;
  /** Total stages that succeeded */
  stagesSucceeded: number;
  /** Total stages that failed */
  stagesFailed: number;
  /** Total triggers executed (direct + rule engine) */
  totalTriggers: number;
  /** First failure message (if any) */
  firstFailure: string | null;
}

// ============================================================================
// UI State Types
// ============================================================================

/**
 * Trace loading state
 */
export interface TraceLoadingState {
  /** Whether trace is loading */
  loading: boolean;
  /** Whether trace has been searched */
  hasSearched: boolean;
  /** Loading progress (0-100) */
  progress: number;
  /** Current loading stage message */
  loadingMessage: string;
}

/**
 * Expanded sections state for the visualization
 */
export interface ExpandedSectionsState {
  /** Ingestion phase expanded */
  ingestion: boolean;
  /** Routing phase expanded */
  routing: boolean;
}

// ============================================================================
// Legacy Types (kept for backward compatibility)
// ============================================================================

/**
 * @deprecated Use specific stage data types instead
 */
export interface TraceStageData {
  /** Stage type identifier */
  stage: string;
  /** Stage display label */
  label: string;
  /** Stage icon */
  icon: string;
  /** Stage color */
  color: string;
  /** Whether this stage has data */
  hasData: boolean;
  /** Whether this stage was successful */
  success: boolean | null;
  /** Created date when this stage was processed */
  created: string | null;
  /** Execution time in milliseconds (if applicable) */
  executionTime: number | null;
  /** Error message (if any) */
  error: string | null;
  /** Raw data for this stage */
  data: EventsRawResponse | EventsJsExecResponse | null;
}

/**
 * Filter state for event tracer page
 */
export interface EventTracerFilters {
  eventTrackerId?: string;
  startTime?: string;
  endTime?: string;
  includeChildren?: boolean;
}

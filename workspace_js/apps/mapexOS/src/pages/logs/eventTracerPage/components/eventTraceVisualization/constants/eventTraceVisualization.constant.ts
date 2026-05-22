/**
 * EventTraceVisualization Component Constants
 */

/**
 * Animation duration for phase expansion/collapse
 */
export const ANIMATION_DURATION_MS = 300;

/**
 * Default expanded state for phases
 */
export const DEFAULT_EXPANDED_STATE: Record<'ingestion' | 'routing', boolean> = {
  ingestion: true,
  routing: true,
};

/**
 * EventTraceVisualization Component Interfaces
 */

import type { EventTraceResult } from '../../../interfaces';

/**
 * Props for EventTraceVisualization component
 */
export interface EventTraceVisualizationProps {
  /** The trace result data to display */
  traceResult: EventTraceResult | null;
  /** Whether the trace is loading */
  loading: boolean;
}

/**
 * Emits for EventTraceVisualization component
 */
export interface EventTraceVisualizationEmits {
  /** Emitted when user clicks to view stage details */
  (e: 'view-details', data: { stage: string; data: any }): void;
  /** Emitted when user clicks on a trigger */
  (e: 'view-trigger', triggerId: string): void;
}

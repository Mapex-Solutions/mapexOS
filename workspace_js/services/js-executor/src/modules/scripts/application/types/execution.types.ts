import type { StandardizedPayload } from '@mapexos/schemas';

/**
 * Application layer types for script execution
 * These types represent use cases and application services concepts
 */

/**
 * Represents the result of script execution at the application level
 * This is what the application service returns to its clients
 */
export interface ScriptExecutionResult {
  /** Whether the execution was successful */
  success: boolean;
  /** Which step failed (if any) */
  failedAt?: 'decode' | 'validation' | 'transform' | null;
  /** Total execution time */
  totalExecutionTime?: number | null;
  /** Error message if execution failed */
  error?: string | null;
  /** Transformed payload */
  standardizedPayload?: StandardizedPayload;
  /** Asset UUID — enriched by executeScripts for consumer publishing */
  assetUUID?: string;
  /** Asset MongoDB ID — enriched by executeScripts for consumer publishing */
  assetId?: string;
  /** Debug enabled flag — enriched by executeScripts for consumer publishing */
  debugEnabled?: boolean;
}
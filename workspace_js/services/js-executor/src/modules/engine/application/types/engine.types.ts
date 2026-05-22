import type { StandardizedPayload } from '@mapexos/schemas';

/**
 * Application layer types for script engine operations
 * These types represent technical execution details at the application level
 */

/**
 * Represents the result of a single script execution
 * Used by the script engine for detailed execution tracking
 */
export interface SingleScriptResult {
  /** The name of the executed script */
  scriptName: string;
  /** The processed data returned by the script */
  data: any;
  /** Execution time in milliseconds */
  executionTime: number;
  /** Whether the execution was successful */
  success: boolean;
  /** Error message if execution failed */
  error?: string;
}

/**
 * Represents the complete pipeline execution result from the engine
 * This is the detailed technical result from the script engine
 */
export interface PipelineExecutionResult {
  /** Whether the entire pipeline was successful */
  success: boolean;
  /** The final processed payload after all steps */
  finalPayload?: StandardizedPayload;
  /** Which step failed (if any) */
  failedAt?: 'decode' | 'validation' | 'transform';
  /** Total pipeline execution time */
  totalPipelineTime?: number;
  /** Error */
  error?: string;
}

/**
 * Represents the result of error of a script execution, including details about the script itself and the execution
 */
export interface SanitizedError {
  /** The category of error (e.g., 'SyntaxError', 'ReferenceError', 'TypeError') */
  type: string;
  /** The core error message after sanitization, without isolated-vm references */
  message: string;
  /** The line number where the error occurred, if available */
  line?: number;
  /** The column number where the error occurred, if available */
  column?: number;
  /** A user-friendly error message with helpful tips and context */
  userFriendlyMessage: string;
  /** The original unmodified error message for internal debugging purposes */
  originalError?: string;
}

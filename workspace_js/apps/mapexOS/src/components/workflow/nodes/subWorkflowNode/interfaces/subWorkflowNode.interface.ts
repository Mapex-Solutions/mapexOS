/**
 * Input mapping: parent source → child workflow variable
 */
export interface InputMapping {
  childVariable: string;
  source: { type: string; value: string };
}

/**
 * Output mapping: child workflow output key → parent state variable.
 * Child output is also accessible via nodes.<nodeId>.output.<key> expressions.
 */
export interface OutputMapping {
  outputKey: string;
  targetVariable: string;
}

/**
 * Execution timeout configuration
 */
export interface TimeoutConfig {
  duration: number;
  unit: string;
}

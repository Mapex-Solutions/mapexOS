/**
 * Domain types for script engine entities
 * These represent core business concepts and should be framework-agnostic
 */

/**
 * Represents a set of scripts to be executed in the engine
 * This is the domain representation of executable scripts
 */
export interface ScriptSet {
  /** Script code for payload decoding */
  decode: string;
  /** Script code for payload validation */
  validation: string;
  /** Script code for payload transformation */
  transform: string;
}

/**
 * Domain types for script-related entities
 * These represent core business concepts and should be framework-agnostic
 */

/**
 * Represents the scripts configuration for an asset
 * This is a core domain concept that defines how data flows through the processing pipeline
 */
export interface AssetScripts {
  /** Script code for payload decoding */
  decode: string;
  /** Script code for payload validation */
  validation: string;
  /** Script code for payload transformation */
  transform: string;
}
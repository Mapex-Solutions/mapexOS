/**
 * Script Registry Interface (Port)
 *
 * Defines the contract for Script Registry implementations.
 * Following Hexagonal Architecture - this is the port that adapters implement.
 */

import type {
	CompiledScripts,
	ScriptRegistryConfig,
	ScriptRegistryStats,
	ScriptRegistryMap,
} from '../types';

/**
 * Script Registry Port
 *
 * Manages compiled Script Objects with TTL and LRU eviction.
 * Each Isolate has its own registry since Script objects are bound to their Isolate.
 *
 * @remarks
 * - Script Objects cannot be shared between Isolates (V8 limitation)
 * - TTL refreshes on every access (keeps hot templates in cache)
 * - LRU eviction when maxEntries is reached
 * - Automatic cleanup runs periodically to remove expired entries
 */
export interface ScriptRegistryPort {
	/**
	 * Get compiled scripts from the registry.
	 *
	 * @param isolateId - The isolate identifier
	 * @param templateId - The template identifier
	 * @returns The compiled scripts or undefined if not cached
	 *
	 * @remarks
	 * On HIT: Returns scripts and refreshes TTL (lastAccessedAt = now)
	 * On MISS: Returns undefined, caller should compile and call set()
	 */
	get(isolateId: number, templateId: string): CompiledScripts | undefined;

	/**
	 * Store compiled scripts in the registry.
	 *
	 * @param isolateId - The isolate identifier
	 * @param templateId - The template identifier
	 * @param scripts - The compiled scripts to store
	 *
	 * @remarks
	 * If maxEntries is reached, evicts the least recently used entry (LRU)
	 */
	set(isolateId: number, templateId: string, scripts: CompiledScripts): void;

	/**
	 * Invalidate a specific template across all isolates.
	 * Used when a template's scripts are updated.
	 *
	 * @param templateId - The template identifier to invalidate
	 * @returns Number of entries removed
	 */
	invalidate(templateId: string): number;

	/**
	 * Clear all entries for a specific isolate.
	 * Used when an isolate is recycled.
	 *
	 * @param isolateId - The isolate identifier
	 */
	clearIsolate(isolateId: number): void;

	/**
	 * Get current statistics.
	 *
	 * @returns Registry statistics for monitoring
	 */
	getStats(): ScriptRegistryStats;

	/**
	 * Get the internal registry map for an isolate.
	 * Used for advanced operations and testing.
	 *
	 * @param isolateId - The isolate identifier
	 * @returns The registry map or undefined if isolate not found
	 */
	getRegistryMap(isolateId: number): ScriptRegistryMap | undefined;

	/**
	 * Start the automatic cleanup timer.
	 * Should be called during service initialization.
	 */
	startCleanup(): void;

	/**
	 * Stop the automatic cleanup timer.
	 * Should be called during service shutdown.
	 */
	stopCleanup(): void;

	/**
	 * Get the current configuration.
	 */
	getConfig(): ScriptRegistryConfig;
}

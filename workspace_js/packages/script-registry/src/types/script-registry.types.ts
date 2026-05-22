/**
 * Script Registry Types
 *
 * Types for managing compiled Script Objects with TTL and LRU eviction.
 * Used to cache isolated-vm Script objects per Isolate.
 */

import type { Script } from 'isolated-vm';

/**
 * Compiled scripts for a template pipeline (decode, validation, transform)
 */
export interface CompiledScripts {
	/** Decode script (optional) */
	decode?: Script;
	/** Validation script (optional) */
	validation?: Script;
	/** Transform script (optional) */
	transform?: Script;
	/** Validator setup script - injects MapexValidator (optional) */
	validatorSetup?: Script;
	/** Combined single-run pipeline script (decode+validation+transform in one execution) */
	pipeline?: Script;
}

/**
 * Entry in the Script Registry with TTL management
 */
export interface ScriptRegistryEntry {
	/** Compiled scripts for this template */
	scripts: CompiledScripts;
	/** Timestamp of last access (for TTL refresh and LRU) */
	lastAccessedAt: number;
	/** Timestamp when entry was created */
	createdAt: number;
}

/**
 * Configuration for Script Registry
 */
export interface ScriptRegistryConfig {
	/** Max entries per isolate (default: 1000) */
	maxEntries: number;
	/** TTL in milliseconds (default: 900000 = 15 min) */
	ttlMs: number;
	/** Cleanup interval in ms (default: 60000 = 1 min) */
	cleanupIntervalMs: number;
}

/**
 * Default configuration for Script Registry
 */
export const DEFAULT_SCRIPT_REGISTRY_CONFIG: ScriptRegistryConfig = {
	maxEntries: 1000,
	ttlMs: 900_000, // 15 minutes
	cleanupIntervalMs: 60_000, // 1 minute
};

/**
 * Statistics for monitoring
 */
export interface ScriptRegistryStats {
	/** Total entries across all isolates */
	totalEntries: number;
	/** Cache hits (script already compiled) */
	hits: number;
	/** Cache misses (needed to compile) */
	misses: number;
	/** Entries evicted due to maxEntries limit (LRU) */
	evictions: number;
	/** Entries removed due to TTL expiration */
	ttlExpirations: number;
	/** Number of cleanup cycles executed */
	cleanupCycles: number;
}

/**
 * Script Registry map type - Maps templateId to entry
 */
export type ScriptRegistryMap = Map<string, ScriptRegistryEntry>;

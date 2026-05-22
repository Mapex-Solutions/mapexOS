/**
 * Script Registry Service
 *
 * Manages compiled Script Objects with TTL and LRU eviction.
 * Each Isolate has its own registry since Script objects are bound to their Isolate.
 *
 * Performance characteristics:
 * - get(): O(1) hash lookup + TTL refresh
 * - set(): O(1) hash insert + O(n) eviction if over limit
 * - cleanup(): O(n) scan for expired entries
 *
 * Memory management:
 * - TTL-based expiration removes stale entries
 * - LRU eviction keeps most recently used templates
 * - Per-isolate maps prevent cross-isolate pollution
 */

import type { ScriptRegistryPort } from '../interfaces';
import type {
	CompiledScripts,
	ScriptRegistryConfig,
	ScriptRegistryStats,
	ScriptRegistryEntry,
	ScriptRegistryMap,
} from '../types';
import { DEFAULT_SCRIPT_REGISTRY_CONFIG } from '../types';

/**
 * Script Registry Service Implementation
 *
 * @example
 * ```typescript
 * const registry = new ScriptRegistryService({
 *   maxEntries: 1000,
 *   ttlMs: 900_000,      // 15 min
 *   cleanupIntervalMs: 60_000, // 1 min
 * });
 *
 * registry.startCleanup();
 *
 * // Get or compile scripts
 * let scripts = registry.get(isolateId, templateId);
 * if (!scripts) {
 *   scripts = await compileScripts(isolate, userScripts);
 *   registry.set(isolateId, templateId, scripts);
 * }
 *
 * // On shutdown
 * registry.stopCleanup();
 * ```
 */
export class ScriptRegistryService implements ScriptRegistryPort {
	/** Per-isolate registries: Map<isolateId, Map<templateId, entry>> */
	private readonly registries: Map<number, ScriptRegistryMap> = new Map();

	/** Configuration */
	private readonly config: ScriptRegistryConfig;

	/** Statistics */
	private readonly stats: ScriptRegistryStats = {
		totalEntries: 0,
		hits: 0,
		misses: 0,
		evictions: 0,
		ttlExpirations: 0,
		cleanupCycles: 0,
	};

	/** Cleanup timer reference */
	private cleanupTimer: NodeJS.Timeout | null = null;

	/**
	 * Creates a new Script Registry Service.
	 *
	 * @param config - Optional configuration (uses defaults if not provided)
	 */
	constructor(config?: Partial<ScriptRegistryConfig>) {
		this.config = { ...DEFAULT_SCRIPT_REGISTRY_CONFIG, ...config };
	}

	/**
	 * Get compiled scripts from the registry.
	 * On HIT: Returns scripts and refreshes TTL.
	 * On MISS: Returns undefined.
	 */
	get(isolateId: number, templateId: string): CompiledScripts | undefined {
		const registry = this.registries.get(isolateId);
		if (!registry) {
			this.stats.misses++;
			return undefined;
		}

		const entry = registry.get(templateId);
		if (!entry) {
			this.stats.misses++;
			return undefined;
		}

		// Check if entry is expired
		const now = Date.now();
		if (now - entry.lastAccessedAt > this.config.ttlMs) {
			// Entry expired, remove it
			registry.delete(templateId);
			this.stats.totalEntries--;
			this.stats.ttlExpirations++;
			this.stats.misses++;
			return undefined;
		}

		// HIT - refresh TTL
		entry.lastAccessedAt = now;
		this.stats.hits++;
		return entry.scripts;
	}

	/**
	 * Store compiled scripts in the registry.
	 * Evicts LRU entry if maxEntries is reached.
	 */
	set(isolateId: number, templateId: string, scripts: CompiledScripts): void {
		// Get or create registry for this isolate
		let registry = this.registries.get(isolateId);
		if (!registry) {
			registry = new Map();
			this.registries.set(isolateId, registry);
		}

		// Check if we need to evict (LRU)
		if (registry.size >= this.config.maxEntries) {
			this.evictLRU(registry);
		}

		// Add new entry
		const now = Date.now();
		const entry: ScriptRegistryEntry = {
			scripts,
			lastAccessedAt: now,
			createdAt: now,
		};

		const isUpdate = registry.has(templateId);
		registry.set(templateId, entry);

		if (!isUpdate) {
			this.stats.totalEntries++;
		}
	}

	/**
	 * Invalidate a specific template across all isolates.
	 */
	invalidate(templateId: string): number {
		let removed = 0;

		for (const registry of this.registries.values()) {
			if (registry.delete(templateId)) {
				removed++;
				this.stats.totalEntries--;
			}
		}

		return removed;
	}

	/**
	 * Clear all entries for a specific isolate.
	 */
	clearIsolate(isolateId: number): void {
		const registry = this.registries.get(isolateId);
		if (registry) {
			this.stats.totalEntries -= registry.size;
			registry.clear();
			this.registries.delete(isolateId);
		}
	}

	/**
	 * Get current statistics.
	 */
	getStats(): ScriptRegistryStats {
		// Recalculate totalEntries for accuracy
		let total = 0;
		for (const registry of this.registries.values()) {
			total += registry.size;
		}
		this.stats.totalEntries = total;

		return { ...this.stats };
	}

	/**
	 * Get the internal registry map for an isolate.
	 */
	getRegistryMap(isolateId: number): ScriptRegistryMap | undefined {
		return this.registries.get(isolateId);
	}

	/**
	 * Start the automatic cleanup timer.
	 */
	startCleanup(): void {
		if (this.cleanupTimer) {
			return; // Already running
		}

		this.cleanupTimer = setInterval(() => {
			this.runCleanup();
		}, this.config.cleanupIntervalMs);

		// Don't prevent process from exiting
		this.cleanupTimer.unref();
	}

	/**
	 * Stop the automatic cleanup timer.
	 */
	stopCleanup(): void {
		if (this.cleanupTimer) {
			clearInterval(this.cleanupTimer);
			this.cleanupTimer = null;
		}
	}

	/**
	 * Get the current configuration.
	 */
	getConfig(): ScriptRegistryConfig {
		return { ...this.config };
	}

	/**
	 * Run cleanup cycle - remove expired entries.
	 * Called automatically by the cleanup timer.
	 */
	private runCleanup(): void {
		const now = Date.now();
		let expired = 0;

		for (const registry of this.registries.values()) {
			for (const [templateId, entry] of registry) {
				if (now - entry.lastAccessedAt > this.config.ttlMs) {
					registry.delete(templateId);
					expired++;
				}
			}
		}

		if (expired > 0) {
			this.stats.totalEntries -= expired;
			this.stats.ttlExpirations += expired;
		}

		this.stats.cleanupCycles++;
	}

	/**
	 * Evict the least recently used entry from a registry.
	 */
	private evictLRU(registry: ScriptRegistryMap): void {
		let oldestKey: string | null = null;
		let oldestTime = Infinity;

		for (const [key, entry] of registry) {
			if (entry.lastAccessedAt < oldestTime) {
				oldestTime = entry.lastAccessedAt;
				oldestKey = key;
			}
		}

		if (oldestKey) {
			registry.delete(oldestKey);
			this.stats.totalEntries--;
			this.stats.evictions++;
		}
	}
}

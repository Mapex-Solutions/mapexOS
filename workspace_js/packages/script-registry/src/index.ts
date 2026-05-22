/**
 * @mapexos/script-registry
 *
 * Script Registry with TTL and LRU eviction for isolated-vm Script Objects.
 *
 * Features:
 * - Per-isolate Script Object caching (V8 requirement)
 * - TTL-based expiration with refresh on access
 * - LRU eviction when maxEntries is reached
 * - Automatic cleanup timer
 * - Pool size auto-detection based on CPU cores
 *
 * @example
 * ```typescript
 * import {
 *   ScriptRegistryService,
 *   getPoolSize,
 *   type ScriptRegistryConfig,
 * } from '@mapexos/script-registry';
 *
 * // Create registry with custom config
 * const registry = new ScriptRegistryService({
 *   maxEntries: 1000,
 *   ttlMs: 900_000,      // 15 min
 *   cleanupIntervalMs: 60_000, // 1 min
 * });
 *
 * // Start cleanup timer
 * registry.startCleanup();
 *
 * // Get pool size (0 = auto based on CPU)
 * const poolSize = getPoolSize(0);
 *
 * // Use registry
 * let scripts = registry.get(isolateId, templateId);
 * if (!scripts) {
 *   scripts = await compile(...);
 *   registry.set(isolateId, templateId, scripts);
 * }
 *
 * // On shutdown
 * registry.stopCleanup();
 * ```
 *
 * @packageDocumentation
 */

// Types
export * from './types';

// Interfaces (Ports)
export * from './interfaces';

// Services
export * from './services';

// Utils
export * from './utils';

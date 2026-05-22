/**
 * Pool Size Utility
 *
 * Helper function to calculate optimal isolate pool size based on CPU cores.
 */

import { cpus } from 'os';

/** Minimum pool size (ensures at least some concurrency) */
const MIN_POOL_SIZE = 2;

/** Maximum pool size (prevents excessive resource usage) */
const MAX_POOL_SIZE = 8;

/**
 * Get optimal pool size based on configuration or CPU cores.
 *
 * @param configValue - Pool size from configuration (0 = auto)
 * @returns Optimal pool size
 *
 * @example
 * ```typescript
 * // Explicit value from ENV
 * getPoolSize(4)  // returns 4
 *
 * // Auto-detect based on CPU cores
 * getPoolSize(0)  // returns Math.min(Math.max(cpuCount, 2), 8)
 * ```
 *
 * @remarks
 * When configValue is 0 (auto):
 * - Uses number of CPU cores
 * - Minimum: 2 (ensures some concurrency even on single-core)
 * - Maximum: 8 (prevents excessive memory usage)
 */
export function getPoolSize(configValue: number): number {
	if (configValue > 0) {
		return configValue;
	}

	// Auto: use CPU cores with min/max bounds
	const cpuCount = cpus().length;
	return Math.max(MIN_POOL_SIZE, Math.min(cpuCount, MAX_POOL_SIZE));
}

/**
 * Get CPU count for informational purposes.
 *
 * @returns Number of CPU cores
 */
export function getCpuCount(): number {
	return cpus().length;
}

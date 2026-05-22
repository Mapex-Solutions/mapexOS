import type { Logger } from '@mapexos/microservices';
import type { TieredCacheClient, MinIOClient } from '@mapexos/infrastructure';
import type { ScriptEngineServicePort } from '@modules/engine/application/ports';
import type { PiscinaWorkerConfig } from '@modules/engine/infrastructure/worker';

import type { ScriptEngineMetrics, PiscinaPoolMetrics, PiscinaOptions } from '@modules/engine/application/types';
import { ScriptEngineService } from '@modules/engine/application/services';
import { TieredBytecodeCache } from '@modules/engine/infrastructure';

/**
 * Dependencies required for ScriptEngineService.
 *
 * This interface follows the Dependency Injection pattern from workspace_go,
 * explicitly declaring all dependencies needed by the ScriptEngineService.
 *
 * @remarks
 * ScriptEngineService orchestrates script execution via Piscina worker threads:
 * - **Piscina Workers**: Each worker owns its own V8 Isolate + compiled script cache
 * - **BytecodeCache (TieredCache)**: L0(RAM)/L1(Disk)/L2(MinIO) for compiled bytecode persistence
 * - **Main thread**: Only dispatches work and handles ack/nack (lightweight)
 */
export interface ScriptEngineServiceDependencies {
	/**
	 * Logger instance for service logging.
	 */
	logger: Logger;

	/**
	 * TieredCache client for caching compiled script bytecode.
	 * Uses L1 (Disk) + L2 (MinIO) for persistence and horizontal scaling.
	 */
	bytecodeCache: TieredCacheClient;

	/**
	 * MinIO client for L2 bytecode storage.
	 * Used to store compiled bytecode for horizontal scaling across pods.
	 */
	minioBytecodeClient: MinIOClient;

	/**
	 * Piscina pool options (workers count + worker file path).
	 */
	piscinaOptions: PiscinaOptions;

	/**
	 * Configuration passed to each worker thread via workerData.
	 */
	workerConfig: PiscinaWorkerConfig;

	/**
	 * Optional Prometheus metrics for ScriptEngineService instrumentation.
	 */
	engineMetrics?: ScriptEngineMetrics;

	/**
	 * Optional Prometheus metrics for Piscina pool instrumentation.
	 */
	poolMetrics?: PiscinaPoolMetrics;

}

/**
 * Factory function to create a ScriptEngineService instance.
 *
 * Creates:
 * - TieredBytecodeCache (Infrastructure): Handles bytecode storage in TieredCache + MinIO
 * - ScriptEngineService (Application): Orchestrates execution via Piscina workers
 *
 * @param deps - The dependencies required by ScriptEngineService
 * @returns A new instance of ScriptEngineService implementing ScriptEngineServicePort
 */
export function createScriptEngineService(deps: ScriptEngineServiceDependencies): ScriptEngineServicePort {
	// Create infrastructure layer (bytecode cache for persistence)
	const bytecodeCache = new TieredBytecodeCache(deps.bytecodeCache, deps.minioBytecodeClient, deps.logger);

	// Create application layer (orchestrator with Piscina workers)
	return new ScriptEngineService(
		deps.logger,
		bytecodeCache,
		deps.piscinaOptions,
		deps.workerConfig,
		deps.engineMetrics,
		deps.poolMetrics,
	);
}

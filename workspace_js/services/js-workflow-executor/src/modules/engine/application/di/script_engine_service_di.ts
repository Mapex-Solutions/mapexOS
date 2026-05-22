import type { Logger } from '@mapexos/microservices';
import type { ScriptEngineServicePort } from '@modules/engine/application/ports';
import type { PiscinaWorkerConfig } from '@modules/engine/infrastructure/worker';

import { ScriptEngineService } from '@modules/engine/application/services';
import type { ScriptEngineMetrics, PiscinaPoolMetrics, PiscinaOptions } from '@modules/engine/application/types';

/**
 * Dependencies required for ScriptEngineService.
 *
 * ScriptEngineService orchestrates workflow script execution via Piscina worker threads:
 * - **Piscina Workers**: Each worker owns its own V8 Isolate + compiled script cache
 * - **Main thread**: Only dispatches work and handles ack/nack (lightweight)
 */
export interface ScriptEngineServiceDependencies {
	/** Logger instance for service logging */
	logger: Logger;

	/** Piscina pool options (workers count + worker file path) */
	piscinaOptions: PiscinaOptions;

	/** Configuration passed to each worker thread via workerData */
	workerConfig: PiscinaWorkerConfig;

	/** Optional Prometheus metrics for ScriptEngineService instrumentation */
	engineMetrics?: ScriptEngineMetrics;

	/** Optional Prometheus metrics for Piscina pool instrumentation */
	poolMetrics?: PiscinaPoolMetrics;
}

/**
 * Factory function to create a ScriptEngineService instance.
 *
 * @param deps - The dependencies required by ScriptEngineService
 * @returns A new instance of ScriptEngineService implementing ScriptEngineServicePort
 */
export function createScriptEngineService(deps: ScriptEngineServiceDependencies): ScriptEngineServicePort {
	return new ScriptEngineService(
		deps.logger,
		deps.piscinaOptions,
		deps.workerConfig,
		deps.engineMetrics,
		deps.poolMetrics,
	);
}

import path from 'path';

import { container } from 'tsyringe';
import type { Logger } from '@mapexos/microservices';
import type { TieredCacheClient, MinIOClient } from '@mapexos/infrastructure';
import { LOGGER_TOKEN, ConfigModule, CONFIG_TOKEN } from '@mapexos/microservices';

import { createScriptEngineService, type ScriptEngineServiceDependencies } from '@modules/engine/application/di';
import type { ScriptEngineServicePort } from '@modules/engine/application/ports';
import { resolvePiscinaWorkers } from '@shared/constants';
import { mapexValidatorCode } from '@shared/utils';

import type { JsExecutorMetrics } from '@/bootstrap/metrics';
import { METRICS_TOKEN } from '@shared/constants';

/**
 * Token for ScriptEngineService in DI container
 */
export const SCRIPT_ENGINE_SERVICE_TOKEN = 'ScriptEngineService';

/**
 * Resolves the absolute path to the compiled Piscina worker file.
 *
 * In production (dist/): worker is compiled alongside the rest of the code.
 * In development (ts-node): worker is loaded via ts-node register hooks.
 */
function resolveWorkerPath(): string {
	const ext = __filename.endsWith('.ts') ? '.ts' : '.js';
	return path.resolve(__dirname, `infrastructure/worker/piscina-worker${ext}`);
}

/**
 * InitServices registers all services in the DI container
 * Following workspace_go pattern - called during Phase 2 of module initialization
 */
export function initServices(): void {
	const logger = container.resolve<Logger>(LOGGER_TOKEN);
	logger.debug('[MODULE:ENGINE] Registering services');

	// Resolve metrics for pool and engine instrumentation
	const metrics = container.resolve<JsExecutorMetrics>(METRICS_TOKEN);

	// Register ScriptEngineService factory
	container.register<ScriptEngineServicePort>(SCRIPT_ENGINE_SERVICE_TOKEN, {
		useFactory: (c) => {
			const config = c.resolve<ConfigModule>(CONFIG_TOKEN);
			const workers = resolvePiscinaWorkers(config);
			const workerPath = resolveWorkerPath();

			const deps: ScriptEngineServiceDependencies = {
				logger: c.resolve<Logger>(LOGGER_TOKEN),
				bytecodeCache: c.resolve<TieredCacheClient>('BytecodeCache'),
				minioBytecodeClient: c.resolve<MinIOClient>('MinIOBytecodeClient'),
				piscinaOptions: {
					workers, // CPU_LIMIT-1 workers — serves both HTTP and batch
					workerPath,
				},
				workerConfig: {
					memoryLimitMb: config.get('isolate_memory_limit_mb'),
					timeoutMs: config.get('worker_script_timeout_ms'),
					contextRecycleInterval: config.get('context_recycle_interval'),
					mapexValidatorCode,
				},
				engineMetrics: {
					scriptDuration: metrics.scriptDuration,
					scriptErrors: metrics.scriptErrors,
					compileDuration: metrics.compileDuration,
					bytecodeCache: metrics.bytecodeCache,
					scriptRegistry: metrics.scriptRegistry,
				},
				poolMetrics: {
					piscinaCompleted: metrics.piscinaCompleted,
					piscinaRunDuration: metrics.piscinaRunDuration,
					piscinaWaitDuration: metrics.piscinaWaitDuration,
					piscinaWorkers: metrics.piscinaWorkers,
				},
			};
			return createScriptEngineService(deps);
		},
	});

	logger.debug('[MODULE:ENGINE] Services registered');
}

/**
 * InitListeners initializes async components (Piscina Worker Pool)
 * Following workspace_go pattern - called during Phase 4 of module initialization
 */
export async function initListeners(): Promise<void> {
	const logger = container.resolve<Logger>(LOGGER_TOKEN);
	logger.info('[MODULE:ENGINE] Initializing Script Engine');

	// Resolve the service (triggers factory) and initialize
	const scriptEngine = container.resolve<ScriptEngineServicePort>(SCRIPT_ENGINE_SERVICE_TOKEN);
	await scriptEngine.initialize();

	logger.info('[MODULE:ENGINE] Script Engine initialized');
}

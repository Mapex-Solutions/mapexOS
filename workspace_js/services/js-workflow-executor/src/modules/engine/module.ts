import path from 'path';

import { container } from 'tsyringe';
import type { Logger } from '@mapexos/microservices';
import { LOGGER_TOKEN, ConfigModule, CONFIG_TOKEN } from '@mapexos/microservices';

import { createScriptEngineService, type ScriptEngineServiceDependencies } from '@modules/engine/application/di';
import type { ScriptEngineServicePort } from '@modules/engine/application/ports';
import { resolvePiscinaWorkers, METRICS_TOKEN, SCRIPT_ENGINE_SERVICE_TOKEN } from '@shared/constants';

import type { WorkflowExecutorMetrics } from '@/bootstrap/metrics';

/**
 * Resolves the absolute path to the compiled Piscina worker file.
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

	const metrics = container.resolve<WorkflowExecutorMetrics>(METRICS_TOKEN);

	container.register<ScriptEngineServicePort>(SCRIPT_ENGINE_SERVICE_TOKEN, {
		useFactory: (c) => {
			const config = c.resolve<ConfigModule>(CONFIG_TOKEN);
			const workers = resolvePiscinaWorkers(config);
			const workerPath = resolveWorkerPath();

			const deps: ScriptEngineServiceDependencies = {
				logger: c.resolve<Logger>(LOGGER_TOKEN),
				piscinaOptions: {
					workers,
					workerPath,
				},
				workerConfig: {
					memoryLimitMb: config.get('isolate_memory_limit_mb'),
					timeoutMs: config.get('worker_script_timeout_ms'),
					contextRecycleInterval: config.get('context_recycle_interval'),
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

	const scriptEngine = container.resolve<ScriptEngineServicePort>(SCRIPT_ENGINE_SERVICE_TOKEN);
	await scriptEngine.initialize();

	logger.info('[MODULE:ENGINE] Script Engine initialized');
}

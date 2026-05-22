import type { PiscinaWorkerInput, PiscinaWorkerOutput } from '@modules/engine/infrastructure/worker';

/**
 * Piscina pool statistics for monitoring
 */
export interface PiscinaStats {
	/** Total tasks completed */
	completed: number;
	/** Active worker thread count */
	threads: number;
	/** Average run time in microseconds */
	runTimeAvg: number;
	/** Average wait time in microseconds */
	waitTimeAvg: number;
	/** Pool utilization ratio (0-1) */
	utilization: number;
}

/**
 * Combined stats for pool monitoring
 */
export interface ScriptEngineStats {
	/** Piscina worker pool statistics */
	piscina: PiscinaStats;
}

/**
 * Port interface for Script Engine Service.
 *
 * Responsible for dispatching workflow code node scripts to Piscina worker threads.
 * Each worker owns its own V8 Isolate and compiled script cache.
 *
 * This module is completely independent and focused solely on script execution.
 * It has no knowledge of NATS, external APIs, or business orchestration logic.
 */
export interface ScriptEngineServicePort {
	/**
	 * Initialize the service — creates the Piscina worker thread pool.
	 */
	initialize(): Promise<void>;

	/**
	 * Shutdown the service and release resources.
	 */
	shutdown(): Promise<void>;

	/**
	 * Get combined statistics for monitoring.
	 */
	getPoolStats(): ScriptEngineStats;

	/**
	 * Executes a workflow code node script in a Piscina worker thread.
	 *
	 * The worker gets or compiles the script (cached per worker by cacheKey),
	 * creates a V8 Context with event, state, inputs, nodes,
	 * executes script extracting result = { output, statePatch },
	 * and returns result via structured clone.
	 *
	 * @param input - The workflow script execution input
	 * @returns The script execution result
	 */
	runWorkflowScript(input: PiscinaWorkerInput): Promise<PiscinaWorkerOutput>;
}

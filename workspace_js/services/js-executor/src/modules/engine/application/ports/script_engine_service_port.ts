import type { ScriptSet } from '@modules/engine/domain/types';
import type { PipelineExecutionResult } from '@modules/engine/application/types';
import type { BytecodeCacheContext } from './bytecode_cache_port';
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
 * Port interface for Script Engine Service (business logic contract).
 *
 * This interface follows the Hexagonal Architecture pattern, defining the contract
 * for script execution engine operations using isolated-vm via Piscina worker threads.
 *
 * The engine module is responsible for:
 * - Dispatching script execution to Piscina worker threads
 * - Each worker owns its own V8 Isolate and compiled script cache
 * - Managing the decode -> validate -> transform pipeline per event
 * - Detecting and propagating OOM errors for retry handling
 *
 * @remarks
 * This module is completely independent and focused solely on script execution.
 * It has no knowledge of NATS, external APIs, or business orchestration logic.
 * The scripts module depends on this module, not the other way around.
 */
export interface ScriptEngineServicePort {
	/**
	 * Initialize the service and its dependencies.
	 * Must be called before any script execution.
	 * Creates the Piscina worker thread pool.
	 */
	initialize(): Promise<void>;

	/**
	 * Shutdown the service and release resources.
	 * Waits for in-flight tasks and destroys the Piscina pool.
	 */
	shutdown(): Promise<void>;

	/**
	 * Get combined statistics for monitoring.
	 * Returns metrics for the Piscina worker pool.
	 */
	getPoolStats(): ScriptEngineStats;

	/**
	 * Executes the complete script pipeline (decode -> validation -> transform).
	 *
	 * This method dispatches to a Piscina worker thread that executes
	 * all three script phases in its own V8 Isolate:
	 * 1. Decode: Processes raw payload into structured data
	 * 2. Validation: Validates the decoded data against a schema
	 * 3. Transform: Transforms validated data into standardized format
	 *
	 * @param rawPayload - The initial payload to process through the pipeline
	 * @param userScripts - The set of scripts to execute for each pipeline phase
	 * @param cacheContext - Optional context for bytecode cache key generation (templateId + templateOrgId)
	 * @returns Promise resolving to the complete pipeline execution result including:
	 *          - success: Whether the pipeline completed successfully
	 *          - finalPayload: The transformed payload (if successful)
	 *          - failedAt: Which phase failed (if unsuccessful)
	 *          - error: Error details (if unsuccessful)
	 *          - totalPipelineTime: Total execution time in nanoseconds
	 */
	runScriptPipeline(
		rawPayload: any,
		userScripts: ScriptSet,
		cacheContext?: BytecodeCacheContext
	): Promise<PipelineExecutionResult>;

	/**
	 * Dispatches N events to the pool in parallel via Promise.all.
	 * Each event is processed by one worker thread independently.
	 * Returns results in the same order as the input array.
	 *
	 * Used by NATS consumers after preprocessing (parse, fetch, build).
	 * Workers are agnostic — only V8 exec, no NATS, no publishing.
	 *
	 * @param events - Array of worker inputs (scripts + payload per event)
	 * @returns Array of worker outputs in the same order as input
	 */
	runBatch(events: PiscinaWorkerInput[]): Promise<PiscinaWorkerOutput[]>;
}

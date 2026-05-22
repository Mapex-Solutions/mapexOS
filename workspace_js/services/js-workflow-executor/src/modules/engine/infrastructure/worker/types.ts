/**
 * Piscina Worker Types — Workflow JS Executor
 *
 * Defines the structured-clone boundary between main thread and worker threads.
 * These types must be serializable (no functions, no circular refs, no class instances).
 */

/**
 * Input sent from main thread to worker via piscina.run().
 *
 * Contains the workflow script and execution context (event, state, inputs, nodes).
 */
export interface PiscinaWorkerInput {
	/** The script source code to execute */
	script: string;

	/** Workflow code node cache key for per-worker script caching */
	cacheKey: string;

	/** Event payload from the trigger that started the workflow */
	event: Record<string, any>;

	/** Current workflow instance state */
	state: Record<string, any>;

	/** External inputs provided at trigger time */
	inputs: Record<string, any>;

	/** Outputs from previously executed nodes */
	nodes: Record<string, any>;

	/** Cached V8 bytecode for faster compilation (optional, from BytecodeCache) */
	cachedBytecode?: ArrayBuffer;

	/** Per-script timeout in ms (from node config). If not set, uses worker default from config. */
	timeoutMs?: number;
}

/**
 * Output returned from worker to main thread via piscina.run() return.
 */
export interface PiscinaWorkerOutput {
	/** Whether the script executed successfully */
	success: boolean;

	/** The result.output from the user script (saved to nodeOutputs[nodeId]) */
	output?: any;

	/** The result.statePatch from the user script (merged into instance state) */
	statePatch?: Record<string, any>;

	/** Script execution time in milliseconds */
	executionTime?: number;

	/** Error message (if failed) */
	error?: string;

	/** Whether failure was due to V8 OOM (isolate disposed) */
	isOOM?: boolean;

	/** Freshly produced V8 bytecode (returned on first compile, main thread stores in BytecodeCache) */
	newBytecode?: ArrayBuffer;
}

/**
 * Configuration passed to workers once via workerData.
 */
export interface PiscinaWorkerConfig {
	/** V8 isolate heap limit per worker in MB */
	memoryLimitMb: number;

	/** Script execution timeout in ms */
	timeoutMs: number;

	/** Recycle V8 context every N events to prevent memory leaks */
	contextRecycleInterval: number;
}

import type { Counter, Histogram, Gauge } from 'prom-client';

/**
 * Prometheus metrics accepted by ScriptEngineService (all optional).
 */
export interface ScriptEngineMetrics {
	/** V8 execution time per script type */
	scriptDuration: Histogram;
	/** Script execution errors by type */
	scriptErrors: Counter;
	/** Script compilation time (fresh vs bytecode) */
	compileDuration: Histogram;
	/** Bytecode cache hits per tier */
	bytecodeCache: Counter;
	/** ScriptRegistry compiled script cache lookups */
	scriptRegistry: Counter;
}

/**
 * Prometheus metrics for Piscina pool monitoring.
 */
export interface PiscinaPoolMetrics {
	/** Total tasks completed by Piscina workers */
	piscinaCompleted: Counter;
	/** Worker execution time (from piscina.run() call to return) */
	piscinaRunDuration: Histogram;
	/** Time waiting in queue before a worker picks up the task */
	piscinaWaitDuration: Histogram;
	/** Current active worker thread count */
	piscinaWorkers: Gauge;
}

/**
 * Options for Piscina pool creation.
 */
export interface PiscinaOptions {
	/** Number of worker threads (minThreads = maxThreads for stable pool) */
	workers: number;
	/** Absolute path to the compiled worker file */
	workerPath: string;
}

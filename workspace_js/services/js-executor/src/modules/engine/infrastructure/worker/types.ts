/**
 * Piscina Worker Types
 *
 * Defines the structured-clone boundary between main thread and worker threads.
 * These types must be serializable (no functions, no circular refs, no class instances).
 */

/**
 * Input sent from main thread to worker via piscina.run().
 *
 * Size: ~5KB typical (821B payload + 2-5KB scripts + 24B templateId)
 */
export interface PiscinaWorkerInput {
	/** The raw IoT event payload (already JSON-parsed by main thread) */
	rawPayload: any;

	/** Script source code strings for the pipeline */
	scripts: {
		decode?: string;
		validation?: string;
		transform?: string;
	};

	/** Template identifier for per-worker script caching */
	templateId: string;
}

/**
 * Output returned from worker to main thread via piscina.run() return.
 *
 * Size: ~1KB typical (success flag + finalPayload + timing)
 */
export interface PiscinaWorkerOutput {
	/** Whether the pipeline completed successfully */
	success: boolean;

	/** The final transformed payload (if successful) */
	finalPayload?: any;

	/** Which pipeline step failed (if any) */
	failedAt?: 'decode' | 'validation' | 'transform';

	/** Total pipeline execution time in milliseconds */
	totalPipelineTime?: number;

	/** Error message (if failed) */
	error?: string;

	/** Whether failure was due to V8 OOM (isolate disposed) */
	isOOM?: boolean;
}

/**
 * Configuration passed to workers once via workerData.
 *
 * Size: ~6KB (mostly mapexValidatorCode string)
 */
export interface PiscinaWorkerConfig {
	/** V8 isolate heap limit per worker in MB */
	memoryLimitMb: number;

	/** Script execution timeout in ms */
	timeoutMs: number;

	/** Recycle V8 context every N events to prevent memory leaks */
	contextRecycleInterval: number;

	/** Static MapexValidator JavaScript code injected before validation scripts */
	mapexValidatorCode: string;
}


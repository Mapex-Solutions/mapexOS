import type { ScriptEngineServicePort } from '@modules/engine/application/ports';
import type { PiscinaWorkerInput, PiscinaWorkerOutput, PiscinaWorkerConfig } from '@modules/engine/infrastructure/worker';
import type { Logger } from '@mapexos/microservices';
import type { ScriptEngineMetrics, PiscinaPoolMetrics, PiscinaOptions } from '@modules/engine/application/types';

import Piscina from 'piscina';

import { OOMError } from '@modules/engine/domain/errors';

/**
 * Application service that orchestrates workflow script execution via Piscina worker threads.
 *
 * - **Application Layer (this service)**: Dispatches workflow scripts to workers
 * - **Worker Thread (piscina-worker.ts)**: V8 isolate execution with script/bytecode caching
 * - **Infrastructure Layer (Piscina)**: Handles worker thread pool management
 *
 * The main thread is lightweight: only dispatches work and handles ack/nack.
 * Workers own their own V8 Isolates and execute workflow code node scripts.
 */
export class ScriptEngineService implements ScriptEngineServicePort {
	private piscina: Piscina | null = null;
	private isInitialized = false;
	private readonly metrics?: ScriptEngineMetrics;
	private readonly poolMetrics?: PiscinaPoolMetrics;

	constructor(
		private readonly logger: Logger,
		private readonly piscinaOptions: PiscinaOptions,
		private readonly workerConfig: PiscinaWorkerConfig,
		metrics?: ScriptEngineMetrics,
		poolMetrics?: PiscinaPoolMetrics,
	) {
		this.metrics = metrics;
		this.poolMetrics = poolMetrics;
	}

	/**
	 * Initialize the service — creates the Piscina worker pool.
	 */
	async initialize(): Promise<void> {
		if (this.isInitialized) {
			return;
		}

		this.logger.info('[SERVICE:ScriptEngine] Initializing Piscina worker pool...');

		const WORKER_BLOCKED_FLAGS = new Set(['--no-node-snapshot']);
		const commonExecArgv = [
			...process.execArgv.filter(f => !WORKER_BLOCKED_FLAGS.has(f)),
			'--disable-warning=ExperimentalWarning',
			'--disable-warning=MODULE_TYPELESS_PACKAGE_JSON',
		];

		this.piscina = new Piscina({
			filename: this.piscinaOptions.workerPath,
			minThreads: this.piscinaOptions.workers,
			maxThreads: this.piscinaOptions.workers,
			workerData: this.workerConfig,
			execArgv: commonExecArgv,
		});

		this.isInitialized = true;
		this.logger.info(
			`[SERVICE:ScriptEngine] Ready (${this.piscinaOptions.workers} workers, ` +
			`memory: ${this.workerConfig.memoryLimitMb}MB/worker, ` +
			`recycle: every ${this.workerConfig.contextRecycleInterval} events)`
		);
	}

	/**
	 * Shutdown the service and release resources.
	 */
	async shutdown(): Promise<void> {
		if (!this.isInitialized) {
			return;
		}

		this.logger.info('[SERVICE:ScriptEngine] Shutting down Piscina pool...');
		if (this.piscina) await this.piscina.destroy();
		this.piscina = null;
		this.isInitialized = false;
		this.logger.info('[SERVICE:ScriptEngine] Shutdown complete');
	}

	/**
	 * Get pool statistics for monitoring.
	 */
	getPoolStats() {
		if (!this.piscina) {
			return {
				piscina: { completed: 0, threads: 0, runTimeAvg: 0, waitTimeAvg: 0, utilization: 0 },
			};
		}

		const runTimeAvg = this.piscina.histogram?.runTime?.average ?? 0;
		const waitTimeAvg = this.piscina.histogram?.waitTime?.average ?? 0;

		if (this.poolMetrics) {
			this.poolMetrics.piscinaWorkers.set(this.piscinaOptions.workers);
		}

		return {
			piscina: {
				completed: this.piscina.completed,
				threads: this.piscinaOptions.workers,
				runTimeAvg: Math.round(runTimeAvg),
				waitTimeAvg: Math.round(waitTimeAvg),
				utilization: this.piscina.utilization,
			},
		};
	}

	/**
	 * Executes a workflow code node script by dispatching to a Piscina worker thread.
	 *
	 * @param input - The workflow script execution input
	 * @returns The script execution result
	 */
	async runWorkflowScript(input: PiscinaWorkerInput): Promise<PiscinaWorkerOutput> {
		if (!this.isInitialized) {
			await this.initialize();
		}

		const startTime = process.hrtime.bigint();

		try {
			const result: PiscinaWorkerOutput = await this.piscina!.run(input);

			this.observePiscinaMetrics(startTime);

			if (result.isOOM) {
				throw new OOMError(`Worker V8 OOM: ${result.error ?? 'isolate disposed'}`);
			}

			if (!result.success) {
				this.incScriptError('workflow_code', result.error);
			}

			return result;
		} catch (error) {
			if (error instanceof OOMError) {
				throw error;
			}

			const totalTime = this.elapsedMs(startTime);
			const errorMessage = error instanceof Error ? error.message : String(error);

			this.logger.error(`[SERVICE:ScriptEngine] Script execution failed after ${totalTime}ms: ${errorMessage}`);

			return {
				success: false,
				executionTime: totalTime,
				error: errorMessage,
			};
		}
	}

	/**
	 * Records Piscina pool metrics (run duration and completed count).
	 *
	 * @param startTime - High-resolution timestamp from process.hrtime.bigint()
	 */
	private observePiscinaMetrics(startTime: bigint): void {
		if (!this.poolMetrics) return;

		const elapsed = Number(process.hrtime.bigint() - startTime) / 1e9;
		this.poolMetrics.piscinaRunDuration.observe(elapsed);
		this.poolMetrics.piscinaCompleted.inc();
	}

	/**
	 * Increments the script error counter with a classified error type.
	 *
	 * @param scriptType - The script type label (e.g. 'workflow_code')
	 * @param error - The raw error message to classify
	 */
	private incScriptError(scriptType: string, error?: string): void {
		if (this.metrics) {
			const errorType = this.classifyError(error);
			this.metrics.scriptErrors.inc({ script_type: scriptType, error_type: errorType });
		}
	}

	/**
	 * Classifies an error message into a category for metrics labeling.
	 *
	 * @param error - The raw error message
	 * @returns The error category: 'syntax' | 'reference' | 'type' | 'timeout' | 'other'
	 */
	private classifyError(error?: string): string {
		if (!error) return 'other';
		const lower = error.toLowerCase();
		if (lower.includes('syntaxerror') || lower.includes('unexpected token')) return 'syntax';
		if (lower.includes('referenceerror') || lower.includes('is not defined')) return 'reference';
		if (lower.includes('typeerror')) return 'type';
		if (lower.includes('timeout') || lower.includes('script execution timed out')) return 'timeout';
		return 'other';
	}

	/**
	 * Calculates elapsed time in milliseconds from a high-resolution timestamp.
	 *
	 * @param startTime - High-resolution timestamp from process.hrtime.bigint()
	 * @returns Elapsed time in milliseconds (3 decimal places)
	 */
	private elapsedMs(startTime: bigint): number {
		return Number((Number(process.hrtime.bigint() - startTime) / 1e6).toFixed(3));
	}
}

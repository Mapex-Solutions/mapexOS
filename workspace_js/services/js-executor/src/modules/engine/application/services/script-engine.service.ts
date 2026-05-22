import type { ScriptSet } from '@modules/engine/domain/types';
import type { PipelineExecutionResult, ScriptEngineMetrics, PiscinaPoolMetrics, PiscinaOptions } from '@modules/engine/application/types';
import type { ScriptEngineServicePort, BytecodeCachePort, BytecodeCacheContext } from '@modules/engine/application/ports';
import type { PiscinaWorkerInput, PiscinaWorkerOutput, PiscinaWorkerConfig } from '@modules/engine/infrastructure/worker';

import type { Logger } from '@mapexos/microservices';

import Piscina from 'piscina';

import { OOMError } from '@modules/engine/domain/errors';

/**
 * Application service that orchestrates script execution via Piscina worker threads.
 *
 * Manages a SINGLE Piscina pool that serves both HTTP (single event) and
 * NATS (batch) use cases. Workers are agnostic — they receive scripts+payload,
 * execute V8, and return the result. No NATS, no cache, no publishing in workers.
 *
 * Architecture:
 * - Application Layer (this service): Dispatches to pool, collects results
 * - Domain Layer (ScriptExecutor in worker): PURE V8 execution
 * - Infrastructure Layer (BytecodeCache, Piscina): Caching and pool management
 */
export class ScriptEngineService implements ScriptEngineServicePort {
	private piscina: Piscina | null = null;
	private isInitialized = false;
	private readonly metrics?: ScriptEngineMetrics;
	private readonly poolMetrics?: PiscinaPoolMetrics;

	constructor(
		private readonly logger: Logger,
		private readonly bytecodeCache: BytecodeCachePort,
		private readonly piscinaOptions: PiscinaOptions,
		private readonly workerConfig: PiscinaWorkerConfig,
		metrics?: ScriptEngineMetrics,
		poolMetrics?: PiscinaPoolMetrics,
	) {
		this.metrics = metrics;
		this.poolMetrics = poolMetrics;
	}

	/**
	 * Initializes the Piscina worker pool.
	 * Called automatically on first use, or explicitly for eager initialization.
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
	 * Shuts down the pool and releases resources.
	 * Waits for in-flight tasks to complete before destroying.
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
	 * Returns pool statistics for monitoring.
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
	 * Executes a single event through the decode -> validate -> transform pipeline.
	 * Dispatches to one worker thread, waits for result, maps to domain type.
	 *
	 * Used by: HTTP test endpoints, single-event processing paths.
	 *
	 * @param rawPayload - The IoT event payload
	 * @param userScripts - Scripts for each pipeline step
	 * @param cacheContext - Template context for bytecode cache key
	 * @returns Pipeline execution result with transformed payload or error
	 */
	async runScriptPipeline(
		rawPayload: any,
		userScripts: ScriptSet,
		cacheContext?: BytecodeCacheContext
	): Promise<PipelineExecutionResult> {
		if (!this.isInitialized) {
			await this.initialize();
		}

		const pipelineStart = process.hrtime.bigint();
		const templateId = cacheContext?.templateId ?? 'default';

		const input: PiscinaWorkerInput = {
			rawPayload,
			scripts: {
				decode: userScripts.decode,
				validation: userScripts.validation,
				transform: userScripts.transform,
			},
			templateId,
		};

		try {
			const result: PiscinaWorkerOutput = await this.piscina!.run(input);

			this.observePiscinaMetrics(pipelineStart);

			if (result.isOOM) {
				throw new OOMError(`Worker V8 OOM: ${result.error ?? 'isolate disposed'}`);
			}

			if (result.success) {
				return {
					success: true,
					finalPayload: result.finalPayload,
					totalPipelineTime: result.totalPipelineTime,
				};
			}

			this.incScriptError(result.failedAt ?? 'unknown', result.error);

			return {
				success: false,
				failedAt: result.failedAt,
				error: result.error,
				totalPipelineTime: result.totalPipelineTime,
			};
		} catch (error) {
			if (error instanceof OOMError) {
				throw error;
			}

			const totalTime = this.elapsedMs(pipelineStart);
			const errorMessage = error instanceof Error ? error.message : String(error);

			this.logger.error(`[SERVICE:ScriptEngine] Pipeline execution failed after ${totalTime}ms: ${errorMessage}`);

			return {
				success: false,
				totalPipelineTime: totalTime,
				error: errorMessage,
			};
		}
	}

	/**
	 * Dispatches N events to the pool in parallel via Promise.all.
	 * Each event is processed by one worker thread independently.
	 * Returns results in the same order as the input array.
	 *
	 * Used by: NATS consumers (JsExecuteConsumer, MqttDataConsumer) for batch processing.
	 *
	 * @param events - Array of worker inputs (scripts + payload per event)
	 * @returns Array of worker outputs in the same order as input
	 */
	async runBatch(events: PiscinaWorkerInput[]): Promise<PiscinaWorkerOutput[]> {
		if (!this.isInitialized) {
			await this.initialize();
		}

		const start = process.hrtime.bigint();

		const results = await Promise.all(
			events.map(event =>
				(this.piscina!.run(event) as Promise<PiscinaWorkerOutput>).catch(err => ({
					success: false,
					error: err instanceof Error ? err.message : String(err),
					isOOM: err instanceof OOMError || (err instanceof Error && err.message.includes('disposed')),
				} as PiscinaWorkerOutput))
			),
		);

		const totalMs = Number(process.hrtime.bigint() - start) / 1e6;

		if (this.poolMetrics) {
			this.poolMetrics.piscinaCompleted.inc(events.length);
			this.poolMetrics.piscinaRunDuration.observe(totalMs / 1000);
		}

		return results;
	}

	/**
	 * Observes Piscina run/wait durations via Prometheus.
	 */
	private observePiscinaMetrics(startTime: bigint): void {
		if (!this.poolMetrics) return;

		const elapsed = Number(process.hrtime.bigint() - startTime) / 1e9;
		this.poolMetrics.piscinaRunDuration.observe(elapsed);
		this.poolMetrics.piscinaCompleted.inc();
	}

	/**
	 * Increments script error counter with parsed error type.
	 */
	private incScriptError(scriptType: string, error?: string): void {
		if (this.metrics) {
			const errorType = this.classifyError(error);
			this.metrics.scriptErrors.inc({ script_type: scriptType, error_type: errorType });
		}
	}

	/**
	 * Classifies an error message into a category for metrics.
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
	 * Calculates elapsed time in milliseconds from hrtime bigint.
	 */
	private elapsedMs(startTime: bigint): number {
		return Number((Number(process.hrtime.bigint() - startTime) / 1e6).toFixed(3));
	}
}

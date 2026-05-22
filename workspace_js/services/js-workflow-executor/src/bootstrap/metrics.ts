/**
 * Metrics Bootstrap — Initializes Prometheus metrics for the js-workflow-executor service.
 *
 * Creates a MetricsRegistry with namespace "wfexec" and defines custom metrics
 * covering script execution, cache tiers, Piscina workers, and NATS consumers.
 */

import type { Counter, Histogram, Gauge } from 'prom-client';

import { container } from 'tsyringe';
import { MetricsRegistry } from '@mapexos/microservices';

import { METRICS_TOKEN } from '@shared/constants';

/**
 * WorkflowExecutorMetrics bundles all Prometheus metrics for the js-workflow-executor service.
 */
export interface WorkflowExecutorMetrics {
	/** The underlying MetricsRegistry */
	registry: MetricsRegistry;

	/** A. Script Execution (P0) */

	/** Total workflow code executions by status */
	executionsTotal: Counter;
	/** End-to-end execution time (NATS receive → callback publish) */
	executionDuration: Histogram;
	/** V8 execution time per script */
	scriptDuration: Histogram;
	/** Script execution errors by type */
	scriptErrors: Counter;
	/** Script compilation time (fresh vs bytecode) */
	compileDuration: Histogram;

	/** B. Cache Tiers (P1) */

	/** Script source cache hits per tier */
	scriptSourceCache: Counter;
	/** Bytecode cache hits per tier */
	bytecodeCache: Counter;
	/** Piscina worker in-memory compiled script cache */
	scriptRegistry: Counter;

	/** C. Piscina Worker Pool (P1) */

	/** Total tasks completed by Piscina workers */
	piscinaCompleted: Counter;
	/** Worker execution time (piscina.run() call to return) */
	piscinaRunDuration: Histogram;
	/** Queue wait time before worker picks up the task */
	piscinaWaitDuration: Histogram;
	/** Current active worker thread count */
	piscinaWorkers: Gauge;

	/** D. NATS Consumer (P1) */

	/** Messages per batch */
	batchSize: Histogram;
	/** Pending messages in stream per consumer */
	natsConsumerLag: Gauge;
}

/** Histogram Buckets */

/** Execution duration: 1ms to 10s */
const EXECUTION_DURATION_BUCKETS = [0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10];

/** Script duration: 0.5ms to 10s */
const SCRIPT_DURATION_BUCKETS = [0.0005, 0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.5, 1, 5, 10];

/** Compile duration: 1ms to 500ms */
const COMPILE_DURATION_BUCKETS = [0.001, 0.002, 0.005, 0.01, 0.025, 0.05, 0.1, 0.5];

/** Piscina run duration: 0.1ms to 10s */
const PISCINA_RUN_BUCKETS = [0.0001, 0.0005, 0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1, 5, 10];

/** Piscina wait duration: 0.01ms to 100ms */
const PISCINA_WAIT_BUCKETS = [0.00001, 0.0001, 0.0005, 0.001, 0.005, 0.01, 0.05, 0.1];

/** Batch size: 1 to 1000 messages */
const BATCH_SIZE_BUCKETS = [1, 10, 50, 100, 250, 500, 1000];

/**
 * InitMetrics creates the MetricsRegistry and registers all custom metrics in the DI container.
 *
 * @returns WorkflowExecutorMetrics instance
 */
export function initMetrics(): WorkflowExecutorMetrics {
	const reg = new MetricsRegistry('wfexec');

	reg.enableDefaultMetrics();

	const metrics: WorkflowExecutorMetrics = {
		registry: reg,

		/** A. Script Execution */

		executionsTotal: reg.newCounterVec(
			{ name: 'executions_total', help: 'Total workflow code executions by status' },
			['status'],
		),

		executionDuration: reg.newHistogram(
			{ name: 'execution_duration_seconds', help: 'End-to-end execution time', buckets: EXECUTION_DURATION_BUCKETS },
		),

		scriptDuration: reg.newHistogramVec(
			{ name: 'script_duration_seconds', help: 'V8 execution time per script', buckets: SCRIPT_DURATION_BUCKETS },
			['script_type'],
		),

		scriptErrors: reg.newCounterVec(
			{ name: 'script_errors_total', help: 'Script execution errors by type' },
			['script_type', 'error_type'],
		),

		compileDuration: reg.newHistogramVec(
			{ name: 'compile_duration_seconds', help: 'Script compilation time', buckets: COMPILE_DURATION_BUCKETS },
			['source'],
		),

		/** B. Cache Tiers */

		scriptSourceCache: reg.newCounterVec(
			{ name: 'script_source_cache_total', help: 'Script source cache hits per tier' },
			['tier'],
		),

		bytecodeCache: reg.newCounterVec(
			{ name: 'bytecode_cache_total', help: 'Bytecode cache hits per tier' },
			['tier'],
		),

		scriptRegistry: reg.newCounterVec(
			{ name: 'script_registry_total', help: 'Worker in-memory compiled script cache lookups' },
			['result'],
		),

		/** C. Piscina Worker Pool */

		piscinaCompleted: reg.newCounter(
			{ name: 'piscina_completed_total', help: 'Total tasks completed by Piscina workers' },
		),

		piscinaRunDuration: reg.newHistogram(
			{ name: 'piscina_run_duration_seconds', help: 'Worker execution time', buckets: PISCINA_RUN_BUCKETS },
		),

		piscinaWaitDuration: reg.newHistogram(
			{ name: 'piscina_wait_duration_seconds', help: 'Queue wait time before worker', buckets: PISCINA_WAIT_BUCKETS },
		),

		piscinaWorkers: reg.newGauge(
			{ name: 'piscina_workers', help: 'Current active worker thread count' },
		),

		/** D. NATS Consumer */

		batchSize: reg.newHistogram(
			{ name: 'batch_size', help: 'Messages per NATS consumer batch', buckets: BATCH_SIZE_BUCKETS },
		),

		natsConsumerLag: reg.newGaugeVec(
			{ name: 'nats_consumer_lag', help: 'Pending messages in stream per consumer' },
			['consumer'],
		),
	};

	container.register(METRICS_TOKEN, { useValue: metrics });

	return metrics;
}

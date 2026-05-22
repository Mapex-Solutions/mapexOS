/**
 * Metrics Bootstrap — Initializes Prometheus metrics for the js-executor service.
 *
 * Creates a MetricsRegistry with namespace "jsexec" and defines custom metrics
 * covering event processing, script execution, cache tiers, Piscina workers, and NATS consumers.
 *
 * Follows the same pattern as Go's http_gateway metrics bootstrap
 * (workspace_go/services/http_gateway/src/bootstrap/metrics.go).
 */

import type { Counter, Histogram, Gauge } from 'prom-client';

import { container } from 'tsyringe';
import { MetricsRegistry } from '@mapexos/microservices';

import { METRICS_TOKEN } from '@shared/constants';

/**
 * JsExecutorMetrics bundles all Prometheus metrics for the js-executor service.
 */
export interface JsExecutorMetrics {
	/** The underlying MetricsRegistry */
	registry: MetricsRegistry;

	/** A. Event Processing (P0) */

	/** Total events by consumer and status */
	eventsProcessed: Counter;
	/** End-to-end event processing time */
	eventDuration: Histogram;
	/** Incoming payload size in bytes */
	payloadSize: Histogram;

	/** B. Script Execution (P0) */

	/** V8 execution time per script type */
	scriptDuration: Histogram;
	/** Script execution errors by type */
	scriptErrors: Counter;
	/** Script compilation time (fresh vs bytecode) */
	compileDuration: Histogram;

	/** C. Cache Tiers (P1) */

	/** Asset cache hits per tier */
	assetCache: Counter;
	/** Template cache hits per tier */
	templateCache: Counter;
	/** Bytecode cache hits per tier */
	bytecodeCache: Counter;
	/** ScriptRegistry in-memory compiled script cache */
	scriptRegistry: Counter;

	/** D. Piscina Worker Pool (P1) */

	/** Total tasks completed by Piscina workers */
	piscinaCompleted: Counter;
	/** Worker execution time (piscina.run() call to return) */
	piscinaRunDuration: Histogram;
	/** Queue wait time before worker picks up the task */
	piscinaWaitDuration: Histogram;
	/** Current active worker thread count */
	piscinaWorkers: Gauge;

	/** E. NATS Consumer (P1) */

	/** Messages per batch */
	batchSize: Histogram;
	/** Pending messages in stream per consumer */
	natsConsumerLag: Gauge;

	/** F. Heartbeat Gate (TKT-2026-0034) */

	/**
	 * Implicit heartbeats published by js-executor — fires only when
	 * healthMonitor.enabled=true AND heartbeatMode='implicit'.
	 */
	heartbeatsPublished: Counter;
	/**
	 * Heartbeat publishes skipped by the gate, labeled by skip reason:
	 *   - 'disabled'      → healthMonitor.enabled=false
	 *   - 'explicit_mode' → heartbeatMode='explicit'; liveness comes from a
	 *                       different path chosen by the asset's protocol
	 *                       (MQTT broker presence advisories OR
	 *                        HTTP POST /api/v1/heartbeat with body { assetUUID })
	 */
	heartbeatsSkipped: Counter;
}

/** Histogram Buckets */

/** Event duration: 1ms to 2.5s */
const EVENT_DURATION_BUCKETS = [0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5];

/** Script duration: 0.5ms to 100ms */
const SCRIPT_DURATION_BUCKETS = [0.0005, 0.001, 0.002, 0.005, 0.01, 0.025, 0.05, 0.1];

/** Compile duration: 1ms to 500ms */
const COMPILE_DURATION_BUCKETS = [0.001, 0.002, 0.005, 0.01, 0.025, 0.05, 0.1, 0.5];

/** Piscina run duration: 0.1ms to 500ms */
const PISCINA_RUN_BUCKETS = [0.0001, 0.0005, 0.001, 0.005, 0.01, 0.05, 0.1, 0.5];

/** Piscina wait duration: 0.01ms to 100ms */
const PISCINA_WAIT_BUCKETS = [0.00001, 0.0001, 0.0005, 0.001, 0.005, 0.01, 0.05, 0.1];

/** Payload size: 100B to 100KB */
const PAYLOAD_SIZE_BUCKETS = [100, 500, 1000, 5000, 10000, 50000, 100000];

/** Batch size: 1 to 1000 messages */
const BATCH_SIZE_BUCKETS = [1, 10, 50, 100, 250, 500, 1000];

/**
 * InitMetrics creates the MetricsRegistry and registers all custom metrics in the DI container.
 *
 * @returns JsExecutorMetrics instance
 */
export function initMetrics(): JsExecutorMetrics {
	const reg = new MetricsRegistry('jsexec');

	// Enable Node.js default metrics (heap, GC, event loop lag, active handles)
	reg.enableDefaultMetrics();

	const metrics: JsExecutorMetrics = {
		registry: reg,

		/** A. Event Processing */

		eventsProcessed: reg.newCounterVec(
			{ name: 'events_processed_total', help: 'Total events by consumer and status' },
			['consumer', 'status'],
		),

		eventDuration: reg.newHistogramVec(
			{ name: 'event_duration_seconds', help: 'End-to-end event processing time', buckets: EVENT_DURATION_BUCKETS },
			['consumer'],
		),

		payloadSize: reg.newHistogram(
			{ name: 'payload_size_bytes', help: 'Incoming payload size in bytes', buckets: PAYLOAD_SIZE_BUCKETS },
		),

		/** B. Script Execution */

		scriptDuration: reg.newHistogramVec(
			{ name: 'script_duration_seconds', help: 'V8 execution time per script type', buckets: SCRIPT_DURATION_BUCKETS },
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

		/** C. Cache Tiers */

		assetCache: reg.newCounterVec(
			{ name: 'asset_cache_total', help: 'Asset cache hits per tier' },
			['tier'],
		),

		templateCache: reg.newCounterVec(
			{ name: 'template_cache_total', help: 'Template cache hits per tier' },
			['tier'],
		),

		bytecodeCache: reg.newCounterVec(
			{ name: 'bytecode_cache_total', help: 'Bytecode cache hits per tier' },
			['tier'],
		),

		scriptRegistry: reg.newCounterVec(
			{ name: 'script_registry_total', help: 'ScriptRegistry compiled script cache lookups' },
			['result'],
		),

		/** D. Piscina Worker Pool */

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

		/** E. NATS Consumer */

		batchSize: reg.newHistogram(
			{ name: 'batch_size', help: 'Messages per NATS consumer batch', buckets: BATCH_SIZE_BUCKETS },
		),

		natsConsumerLag: reg.newGaugeVec(
			{ name: 'nats_consumer_lag', help: 'Pending messages in stream per consumer' },
			['consumer'],
		),

		/** F. Heartbeat Gate (TKT-2026-0034 explicit mode) */

		heartbeatsPublished: reg.newCounter(
			{ name: 'heartbeats_published_total', help: 'Implicit heartbeats published by js-executor (mode=implicit and asset enabled).' },
		),

		heartbeatsSkipped: reg.newCounterVec(
			{ name: 'heartbeats_skipped_total', help: 'Heartbeat publishes skipped by the gate, labeled by reason.' },
			['reason'],
		),
	};

	// Register in DI container
	container.register(METRICS_TOKEN, { useValue: metrics });

	return metrics;
}

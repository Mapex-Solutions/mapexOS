import os from 'os';

import type { ConfigModule } from '@mapexos/microservices';

/**
 * Shared constants for NATS consumers
 */

/** Service name for DLQ metadata */
export const SERVICE_NAME = 'js-executor';

/** Service type for DLQ metadata */
export const SERVICE_TYPE = 'js-executor';

/** Default retry policy with exponential backoff (1s, 5s, 30s, 2m, 10m) */
export const DEFAULT_RETRY_POLICY = {
	maxRetries: 5,
	backoff: [1000, 5000, 30000, 120000, 600000] as number[],
};

/** Default consumer fetch configuration (non-CPU-dependent) */
export const DEFAULT_CONSUMER_CONFIG = {
	fetchTimeout: 5000,
};

/** Multipliers for CPU-based auto-tuning (Piscina architecture) */
export const CPU_MULTIPLIERS = {
	/** chunk = workers × 8 (Piscina handles internal queueing) */
	chunk: 8,
	/** batch = CPU × 500 (larger batches since workers are faster) */
	batch: 500,
	/** maxAckPending = batch × 2 */
	ackPending: 2,
} as const;

/**
 * Resolves the effective CPU limit from config.
 * This is the single source of truth for all auto-tuning.
 *
 * Validates that CPU_LIMIT does not exceed the available CPUs.
 * If CPU_LIMIT > available CPUs, the service exits with error code 1
 * to prevent silent performance degradation from CPU oversubscription.
 *
 * @param config - ConfigModule instance
 * @returns CPU limit (minimum 1)
 */
export function resolveCpuLimit(config: ConfigModule): number {
	const cpuLimit = Math.max(1, config.get('cpu_limit') as number);
	const availableCpus = os.cpus().length;

	if (cpuLimit > availableCpus) {
		console.error(
			`[FATAL] CPU_LIMIT=${cpuLimit} exceeds available CPUs (${availableCpus}). ` +
			`This causes CPU oversubscription: ${cpuLimit - 1} worker threads + 1 main thread ` +
			`competing for ${availableCpus} cores, resulting in cache thrashing and degraded throughput. ` +
			`Set CPU_LIMIT<=${availableCpus} or increase container CPU allocation.`
		);
		process.exit(1);
	}

	return cpuLimit;
}

/**
 * Resolves the number of Piscina worker threads.
 * Priority: PISCINA_WORKERS (explicit) > CPU_LIMIT - 1 (auto, leave 1 CPU for main thread)
 *
 * @param config - ConfigModule instance
 * @returns Number of Piscina worker threads (minimum 1)
 */
export function resolvePiscinaWorkers(config: ConfigModule): number {
	const explicit = config.get('piscina_workers') as number;
	if (explicit > 0) return explicit;

	// Auto: CPU_LIMIT - 1 (leave 1 CPU for main thread), minimum 1
	return Math.max(1, resolveCpuLimit(config) - 1);
}

/**
 * Resolves the concurrency chunk size for internal batch processing.
 * Priority: CONCURRENCY_CHUNK_SIZE (explicit) > workers × 8 (auto)
 *
 * Large batches are fetched from NATS for network efficiency,
 * then processed in smaller sequential chunks.
 * With Piscina, each chunk item is dispatched to a worker thread.
 *
 * @param config - ConfigModule instance
 * @returns Chunk size for Promise.all batching
 */
export function resolveChunkSize(config: ConfigModule): number {
	const explicit = config.get('concurrency_chunk_size') as number;
	if (explicit > 0) return explicit;

	return resolvePiscinaWorkers(config) * CPU_MULTIPLIERS.chunk;
}

/**
 * Resolves the number of events per piscina.run() call for batch workers.
 * Priority: EVENTS_PER_WORKER (explicit) > default 500
 *
 * Each sub-batch is dispatched as a single piscina.run() call.
 * The batch worker processes all events in the sub-batch sequentially,
 * publishing to NATS directly and flushing once at the end.
 *
 * @param config - ConfigModule instance
 * @returns Number of events per worker dispatch
 */
export function resolveEventsPerWorker(config: ConfigModule): number {
	const explicit = config.get('events_per_worker') as number;
	if (explicit > 0) return explicit;

	return 500;
}

/**
 * Resolves consumer configuration from CPU_LIMIT with ENV overrides.
 * A value of 0 means "auto-calculate from CPU_LIMIT".
 *
 * @param config - ConfigModule instance
 * @returns Resolved consumer configuration
 */
export function resolveConsumerConfig(config: ConfigModule) {
	const cpuLimit = resolveCpuLimit(config);

	const envBatchSize = config.get('nats_consumer_batch_size') as number;
	const envFetchTimeout = config.get('nats_consumer_fetch_timeout') as number;
	const envMaxAckPending = config.get('nats_consumer_max_ack_pending') as number;

	const batchSize = envBatchSize > 0 ? envBatchSize : cpuLimit * CPU_MULTIPLIERS.batch;
	const fetchTimeout = envFetchTimeout > 0 ? envFetchTimeout : DEFAULT_CONSUMER_CONFIG.fetchTimeout;
	const maxAckPending = envMaxAckPending > 0 ? envMaxAckPending : batchSize * CPU_MULTIPLIERS.ackPending;

	return { batchSize, fetchTimeout, maxAckPending };
}

/**
 * Logs the full resolved configuration for diagnostics.
 *
 * @param config - ConfigModule instance
 * @returns Object with all resolved values for logging
 */
export function resolveAllTuning(config: ConfigModule) {
	const cpuLimit = resolveCpuLimit(config);
	const piscinaWorkers = resolvePiscinaWorkers(config);
	const chunkSize = resolveChunkSize(config);
	const eventsPerWorker = resolveEventsPerWorker(config);
	const consumerConfig = resolveConsumerConfig(config);

	return {
		cpuLimit,
		piscinaWorkers,
		chunkSize,
		eventsPerWorker,
		...consumerConfig,
	};
}

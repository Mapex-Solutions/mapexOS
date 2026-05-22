import os from 'os';

import type { ConfigModule } from '@mapexos/microservices';

/**
 * Shared constants for NATS consumers
 */

/** Service name for DLQ metadata */
export const SERVICE_NAME = 'js-workflow-executor';

/** Service type for DLQ metadata */
export const SERVICE_TYPE = 'js-workflow-executor';

/** Default retry policy with exponential backoff (1s, 5s, 30s, 2m, 10m) */
export const DEFAULT_RETRY_POLICY = {
	maxRetries: 5,
	backoff: [1000, 5000, 30000, 120000, 600000] as number[],
};

/** Default consumer fetch configuration */
export const DEFAULT_CONSUMER_CONFIG = {
	fetchTimeout: 5000,
};

/**
 * Resolves the effective CPU limit from config.
 * Validates that CPU_LIMIT does not exceed the available CPUs.
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
			`Set CPU_LIMIT<=${availableCpus} or increase container CPU allocation.`
		);
		process.exit(1);
	}

	return cpuLimit;
}

/**
 * Resolves the number of Piscina worker threads.
 * Priority: PISCINA_WORKERS (explicit) > CPU_LIMIT - 1 (auto)
 *
 * @param config - ConfigModule instance
 * @returns Number of Piscina worker threads (minimum 1)
 */
export function resolvePiscinaWorkers(config: ConfigModule): number {
	const explicit = config.get('piscina_workers') as number;
	if (explicit > 0) return explicit;

	return Math.max(1, resolveCpuLimit(config) - 1);
}

/**
 * Resolves consumer configuration from CPU_LIMIT with ENV overrides.
 *
 * @param config - ConfigModule instance
 * @returns Resolved consumer configuration
 */
export function resolveConsumerConfig(config: ConfigModule) {
	const cpuLimit = resolveCpuLimit(config);

	const envBatchSize = config.get('nats_consumer_batch_size') as number;
	const envFetchTimeout = config.get('nats_consumer_fetch_timeout') as number;
	const envMaxAckPending = config.get('nats_consumer_max_ack_pending') as number;

	const batchSize = envBatchSize > 0 ? envBatchSize : cpuLimit * 500;
	const fetchTimeout = envFetchTimeout > 0 ? envFetchTimeout : DEFAULT_CONSUMER_CONFIG.fetchTimeout;
	const maxAckPending = envMaxAckPending > 0 ? envMaxAckPending : batchSize * 2;

	return { batchSize, fetchTimeout, maxAckPending };
}

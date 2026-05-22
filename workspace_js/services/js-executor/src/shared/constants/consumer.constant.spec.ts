import os from 'os';

import {
	resolveCpuLimit,
	resolvePiscinaWorkers,
	resolveChunkSize,
	resolveEventsPerWorker,
	resolveConsumerConfig,
	resolveAllTuning,
	CPU_MULTIPLIERS,
	DEFAULT_CONSUMER_CONFIG,
} from './consumer.constant';

/** Mock ConfigModule — returns values from a config map */
function mockConfig(overrides: Record<string, unknown> = {}) {
	const defaults: Record<string, unknown> = {
		cpu_limit: 4,
		piscina_workers: 0,
		concurrency_chunk_size: 0,
		events_per_worker: 0,
		nats_consumer_batch_size: 0,
		nats_consumer_fetch_timeout: 0,
		nats_consumer_max_ack_pending: 0,
	};

	const merged = { ...defaults, ...overrides };

	return {
		get: (key: string) => merged[key],
	} as any;
}

describe('Consumer Constants', () => {
	const realCpuCount = os.cpus().length;

	describe('resolveCpuLimit', () => {
		it('should return CPU_LIMIT when within available CPUs', () => {
			const config = mockConfig({ cpu_limit: 4 });
			expect(resolveCpuLimit(config)).toBe(4);
		});

		it('should enforce minimum of 1', () => {
			const config = mockConfig({ cpu_limit: 0 });
			expect(resolveCpuLimit(config)).toBe(1);
		});

		it('should enforce minimum of 1 for negative values', () => {
			const config = mockConfig({ cpu_limit: -5 });
			expect(resolveCpuLimit(config)).toBe(1);
		});

		it('should accept CPU_LIMIT equal to available CPUs', () => {
			const config = mockConfig({ cpu_limit: realCpuCount });
			expect(resolveCpuLimit(config)).toBe(realCpuCount);
		});

		it('should call process.exit(1) when CPU_LIMIT exceeds available CPUs', () => {
			const exitSpy = jest.spyOn(process, 'exit').mockImplementation(() => undefined as never);
			const errorSpy = jest.spyOn(console, 'error').mockImplementation(() => {});

			const config = mockConfig({ cpu_limit: realCpuCount + 1 });
			resolveCpuLimit(config);

			expect(exitSpy).toHaveBeenCalledWith(1);
			expect(errorSpy).toHaveBeenCalledWith(
				expect.stringContaining(`[FATAL] CPU_LIMIT=${realCpuCount + 1} exceeds available CPUs (${realCpuCount})`)
			);

			exitSpy.mockRestore();
			errorSpy.mockRestore();
		});

		it('should include oversubscription details in fatal message', () => {
			const exitSpy = jest.spyOn(process, 'exit').mockImplementation(() => undefined as never);
			const errorSpy = jest.spyOn(console, 'error').mockImplementation(() => {});

			const overLimit = realCpuCount + 10;
			const config = mockConfig({ cpu_limit: overLimit });
			resolveCpuLimit(config);

			expect(errorSpy).toHaveBeenCalledWith(
				expect.stringContaining(`${overLimit - 1} worker threads + 1 main thread`)
			);
			expect(errorSpy).toHaveBeenCalledWith(
				expect.stringContaining(`Set CPU_LIMIT<=${realCpuCount}`)
			);

			exitSpy.mockRestore();
			errorSpy.mockRestore();
		});
	});

	describe('resolvePiscinaWorkers', () => {
		it('should return explicit PISCINA_WORKERS when set', () => {
			const config = mockConfig({ piscina_workers: 6 });
			expect(resolvePiscinaWorkers(config)).toBe(6);
		});

		it('should auto-calculate as CPU_LIMIT - 1', () => {
			const config = mockConfig({ cpu_limit: 8, piscina_workers: 0 });
			expect(resolvePiscinaWorkers(config)).toBe(7);
		});

		it('should enforce minimum of 1 worker', () => {
			const config = mockConfig({ cpu_limit: 1, piscina_workers: 0 });
			expect(resolvePiscinaWorkers(config)).toBe(1);
		});

		it('should enforce minimum of 1 worker for CPU_LIMIT=2', () => {
			const config = mockConfig({ cpu_limit: 2, piscina_workers: 0 });
			expect(resolvePiscinaWorkers(config)).toBe(1);
		});
	});

	describe('resolveChunkSize', () => {
		it('should return explicit CONCURRENCY_CHUNK_SIZE when set', () => {
			const config = mockConfig({ concurrency_chunk_size: 64 });
			expect(resolveChunkSize(config)).toBe(64);
		});

		it('should auto-calculate as workers × chunk multiplier', () => {
			const config = mockConfig({ cpu_limit: 8, piscina_workers: 0, concurrency_chunk_size: 0 });
			// workers = 8-1 = 7, chunk = 7 × 8 = 56
			expect(resolveChunkSize(config)).toBe(7 * CPU_MULTIPLIERS.chunk);
		});
	});

	describe('resolveEventsPerWorker', () => {
		it('should return explicit EVENTS_PER_WORKER when set', () => {
			const config = mockConfig({ events_per_worker: 1000 });
			expect(resolveEventsPerWorker(config)).toBe(1000);
		});

		it('should default to 500', () => {
			const config = mockConfig({ events_per_worker: 0 });
			expect(resolveEventsPerWorker(config)).toBe(500);
		});
	});

	describe('resolveConsumerConfig', () => {
		it('should auto-calculate all values from CPU_LIMIT', () => {
			const config = mockConfig({ cpu_limit: 4 });
			const result = resolveConsumerConfig(config);

			expect(result).toEqual({
				batchSize: 4 * CPU_MULTIPLIERS.batch,       // 2000
				fetchTimeout: DEFAULT_CONSUMER_CONFIG.fetchTimeout, // 5000
				maxAckPending: 4 * CPU_MULTIPLIERS.batch * CPU_MULTIPLIERS.ackPending, // 4000
			});
		});

		it('should use explicit overrides when set', () => {
			const config = mockConfig({
				cpu_limit: 4,
				nats_consumer_batch_size: 500,
				nats_consumer_fetch_timeout: 3000,
				nats_consumer_max_ack_pending: 1000,
			});

			const result = resolveConsumerConfig(config);

			expect(result).toEqual({
				batchSize: 500,
				fetchTimeout: 3000,
				maxAckPending: 1000,
			});
		});
	});

	describe('resolveAllTuning', () => {
		it('should return all resolved values', () => {
			const config = mockConfig({ cpu_limit: 4 });
			const result = resolveAllTuning(config);

			expect(result).toEqual({
				cpuLimit: 4,
				piscinaWorkers: 3,
				chunkSize: 3 * CPU_MULTIPLIERS.chunk,
				eventsPerWorker: 500,
				batchSize: 4 * CPU_MULTIPLIERS.batch,
				fetchTimeout: DEFAULT_CONSUMER_CONFIG.fetchTimeout,
				maxAckPending: 4 * CPU_MULTIPLIERS.batch * CPU_MULTIPLIERS.ackPending,
			});
		});
	});
});

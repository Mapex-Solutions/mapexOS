/**
 * Consumer Constants Unit Tests
 *
 * Tests CPU resolution, worker count calculation, and consumer config resolution.
 */

import type { ConfigModule } from '@mapexos/microservices';

import { resolveCpuLimit, resolvePiscinaWorkers, resolveConsumerConfig } from './consumer.constant';

// ─── Mock Helpers ────────────────────────────────────────────────────

const createMockConfig = (overrides: Record<string, any> = {}): ConfigModule => ({
	get: jest.fn((key: string) => {
		const defaults: Record<string, any> = {
			cpu_limit: 4,
			piscina_workers: 0,
			nats_consumer_batch_size: 0,
			nats_consumer_fetch_timeout: 0,
			nats_consumer_max_ack_pending: 0,
		};
		return overrides[key] ?? defaults[key] ?? 0;
	}),
} as unknown as ConfigModule);

// Mock os.cpus() to return a fixed number of CPUs
jest.mock('os', () => ({
	cpus: jest.fn(() => Array.from({ length: 16 }, () => ({}))),
}));

describe('resolveCpuLimit', () => {
	it('should return CPU_LIMIT from config', () => {
		const config = createMockConfig({ cpu_limit: 4 });

		const result = resolveCpuLimit(config);

		expect(result).toBe(4);
	});

	it('should return minimum 1 when CPU_LIMIT is 0', () => {
		const config = createMockConfig({ cpu_limit: 0 });

		const result = resolveCpuLimit(config);

		expect(result).toBe(1);
	});

	it('should return minimum 1 when CPU_LIMIT is negative', () => {
		const config = createMockConfig({ cpu_limit: -1 });

		const result = resolveCpuLimit(config);

		expect(result).toBe(1);
	});

	it('should exit process when CPU_LIMIT exceeds available CPUs', () => {
		const mockExit = jest.spyOn(process, 'exit').mockImplementation(() => undefined as never);
		const config = createMockConfig({ cpu_limit: 32 });

		resolveCpuLimit(config);

		expect(mockExit).toHaveBeenCalledWith(1);
		mockExit.mockRestore();
	});
});

describe('resolvePiscinaWorkers', () => {
	it('should return CPU_LIMIT - 1 by default', () => {
		const config = createMockConfig({ cpu_limit: 4 });

		const result = resolvePiscinaWorkers(config);

		expect(result).toBe(3);
	});

	it('should return explicit PISCINA_WORKERS when set', () => {
		const config = createMockConfig({ cpu_limit: 4, piscina_workers: 7 });

		const result = resolvePiscinaWorkers(config);

		expect(result).toBe(7);
	});

	it('should return minimum 1 when CPU_LIMIT is 1', () => {
		const config = createMockConfig({ cpu_limit: 1 });

		const result = resolvePiscinaWorkers(config);

		expect(result).toBe(1);
	});
});

describe('resolveConsumerConfig', () => {
	it('should auto-calculate batchSize from CPU_LIMIT', () => {
		const config = createMockConfig({ cpu_limit: 4 });

		const result = resolveConsumerConfig(config);

		expect(result.batchSize).toBe(2000); // 4 * 500
	});

	it('should use explicit batchSize when set', () => {
		const config = createMockConfig({ cpu_limit: 4, nats_consumer_batch_size: 100 });

		const result = resolveConsumerConfig(config);

		expect(result.batchSize).toBe(100);
	});

	it('should auto-calculate maxAckPending from batchSize', () => {
		const config = createMockConfig({ cpu_limit: 4 });

		const result = resolveConsumerConfig(config);

		expect(result.maxAckPending).toBe(4000); // 2000 * 2
	});

	it('should use explicit maxAckPending when set', () => {
		const config = createMockConfig({ cpu_limit: 4, nats_consumer_max_ack_pending: 500 });

		const result = resolveConsumerConfig(config);

		expect(result.maxAckPending).toBe(500);
	});

	it('should use default fetchTimeout when not explicitly set', () => {
		const config = createMockConfig({ cpu_limit: 2 });

		const result = resolveConsumerConfig(config);

		expect(result.fetchTimeout).toBe(5000);
	});

	it('should use explicit fetchTimeout when set', () => {
		const config = createMockConfig({ cpu_limit: 2, nats_consumer_fetch_timeout: 10000 });

		const result = resolveConsumerConfig(config);

		expect(result.fetchTimeout).toBe(10000);
	});
});

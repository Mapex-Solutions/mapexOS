/**
 * Shutdown Bootstrap Unit Tests
 *
 * Tests initShutdown: hook registration, priority order, Express/Redis/NATS cleanup.
 * Mocks: Express Server, DI container (Redis, NATS).
 */

import 'reflect-metadata';

import type { Logger } from '@mapexos/microservices';
import type { Server } from 'http';

import { container } from 'tsyringe';
import { NATS_CONNECTION_TOKEN, REDIS_SERVICE_TOKEN } from '@mapexos/infrastructure';

import { initShutdown } from './shutdown';

// ─── Mock Infrastructure ─────────────────────────────────────────────

jest.mock('@mapexos/infrastructure', () => ({
	closeNatsClient: jest.fn().mockResolvedValue(undefined),
	NATS_CONNECTION_TOKEN: 'NATS_CONNECTION_TOKEN',
	REDIS_SERVICE_TOKEN: 'REDIS_SERVICE_TOKEN',
}));

const { closeNatsClient } = jest.requireMock('@mapexos/infrastructure');

const createMockLogger = (): Logger => ({
	info: jest.fn(),
	debug: jest.fn(),
	warn: jest.fn(),
	error: jest.fn(),
	trace: jest.fn(),
	fatal: jest.fn(),
	child: jest.fn().mockReturnThis(),
} as unknown as Logger);

const createMockServer = (): Server => ({
	close: jest.fn((cb: (err?: Error) => void) => cb()),
} as unknown as Server);

const mockRedisService = { disconnect: jest.fn().mockResolvedValue(undefined) };
const mockNatsClient = { nc: { isClosed: () => false, close: jest.fn() } };

// ─── Tests ───────────────────────────────────────────────────────────

describe('initShutdown (js-executor)', () => {
	let logger: Logger;
	let server: Server;

	beforeEach(() => {
		jest.clearAllMocks();
		logger = createMockLogger();
		server = createMockServer();

		container.register(REDIS_SERVICE_TOKEN, { useValue: mockRedisService });
		container.register(NATS_CONNECTION_TOKEN, { useValue: mockNatsClient });
	});

	it('should create a ShutdownManager with hooks registered', () => {
		const sm = initShutdown(logger, server);
		expect(sm).toBeDefined();
	});

	it('should close Express server on shutdown (P0)', async () => {
		const sm = initShutdown(logger, server);
		await sm.executeShutdown(5000);

		expect(server.close).toHaveBeenCalled();
	});

	it('should disconnect Redis on shutdown (P5)', async () => {
		const sm = initShutdown(logger, server);
		await sm.executeShutdown(5000);

		expect(mockRedisService.disconnect).toHaveBeenCalled();
	});

	it('should close NATS on shutdown (P5)', async () => {
		const sm = initShutdown(logger, server);
		await sm.executeShutdown(5000);

		expect(closeNatsClient).toHaveBeenCalledWith(mockNatsClient);
	});

	it('should execute Express (P0) before Redis/NATS (P5)', async () => {
		const order: string[] = [];

		const serverWithTracking = {
			close: jest.fn((cb: (err?: Error) => void) => {
				order.push('express');
				cb();
			}),
		} as unknown as Server;

		mockRedisService.disconnect.mockImplementation(async () => { order.push('redis'); });
		(closeNatsClient as jest.Mock).mockImplementation(async () => { order.push('nats'); });

		const sm = initShutdown(logger, serverWithTracking);
		await sm.executeShutdown(5000);

		expect(order[0]).toBe('express');
		expect(order).toContain('redis');
		expect(order).toContain('nats');
	});

	it('should complete shutdown even if Redis fails', async () => {
		mockRedisService.disconnect.mockRejectedValue(new Error('Redis error'));

		const sm = initShutdown(logger, server);
		await sm.executeShutdown(5000);

		expect(closeNatsClient).toHaveBeenCalled();
		expect(logger.warn).toHaveBeenCalledWith(expect.stringContaining('redis failed'));
	});
});

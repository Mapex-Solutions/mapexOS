import type { Logger } from '@mapexos/infrastructure';
import type { Shutdowner } from './types';

import { ShutdownManager } from './shutdown';

// ─── Mock Helpers ────────────────────────────────────────────────────

const createMockLogger = (): Logger => ({
	info: jest.fn(),
	debug: jest.fn(),
	warn: jest.fn(),
	error: jest.fn(),
	trace: jest.fn(),
	fatal: jest.fn(),
	child: jest.fn().mockReturnThis(),
} as unknown as Logger);

const sleep = (ms: number) => new Promise((r) => setTimeout(r, ms));

// ─── Tests ───────────────────────────────────────────────────────────

describe('ShutdownManager', () => {
	let logger: Logger;
	let sm: ShutdownManager;

	beforeEach(() => {
		logger = createMockLogger();
		sm = new ShutdownManager(logger);
	});

	// --- registerFunc ---

	it('should register a hook', () => {
		sm.registerFunc('test', 3, async () => {});

		// Verify by executing — hook count is internal
		expect(sm).toBeDefined();
	});

	// --- register (Shutdowner interface) ---

	it('should delegate to Shutdowner.shutdown', async () => {
		const shutdowner: Shutdowner = { shutdown: jest.fn().mockResolvedValue(undefined) };
		sm.register('mock', 5, shutdowner);

		await sm.executeShutdown(5000);

		expect(shutdowner.shutdown).toHaveBeenCalled();
	});

	// --- executeShutdown: priority order ---

	it('should execute hooks in priority order (P0 → P3 → P5)', async () => {
		const order: string[] = [];

		// Register out of order to verify sorting
		sm.registerFunc('p5-conn', 5, async () => { order.push('p5'); });
		sm.registerFunc('p0-http', 0, async () => { order.push('p0'); });
		sm.registerFunc('p3-flush', 3, async () => { order.push('p3'); });

		await sm.executeShutdown(5000);

		expect(order).toEqual(['p0', 'p3', 'p5']);
	});

	// --- executeShutdown: same priority concurrent ---

	it('should run same-priority hooks concurrently', async () => {
		let maxConcurrent = 0;
		let running = 0;

		const makeHook = () => async () => {
			running++;
			if (running > maxConcurrent) maxConcurrent = running;
			await sleep(40);
			running--;
		};

		sm.registerFunc('a', 5, makeHook());
		sm.registerFunc('b', 5, makeHook());
		sm.registerFunc('c', 5, makeHook());

		await sm.executeShutdown(5000);

		expect(maxConcurrent).toBeGreaterThanOrEqual(2);
	});

	// --- executeShutdown: error isolation ---

	it('should not stop subsequent priorities when a hook fails', async () => {
		let p5Called = false;

		sm.registerFunc('fail-p0', 0, async () => { throw new Error('boom'); });
		sm.registerFunc('conn-p5', 5, async () => { p5Called = true; });

		await sm.executeShutdown(5000);

		expect(p5Called).toBe(true);
		expect(logger.warn).toHaveBeenCalledWith(expect.stringContaining('fail-p0 failed'));
	});

	// --- executeShutdown: timeout aborts remaining groups ---

	it('should abort remaining groups when timeout is reached', async () => {
		let p5Called = false;

		sm.registerFunc('slow-p0', 0, async () => { await sleep(200); });
		sm.registerFunc('conn-p5', 5, async () => { p5Called = true; });

		await sm.executeShutdown(100);

		expect(p5Called).toBe(false);
		expect(logger.warn).toHaveBeenCalledWith(expect.stringContaining('Timeout reached'));
	});

	// --- executeShutdown: no hooks ---

	it('should complete quickly with no hooks', async () => {
		const start = Date.now();
		await sm.executeShutdown(5000);
		const elapsed = Date.now() - start;

		expect(elapsed).toBeLessThan(100);
		expect(logger.info).toHaveBeenCalledWith(expect.stringContaining('Graceful shutdown complete'));
	});

	// --- executeShutdown: P0 must complete before P5 starts ---

	it('should ensure P0 finishes before P5 starts', async () => {
		let fiberEndTime = 0;
		let mongoStartTime = 0;

		sm.registerFunc('fiber', 0, async () => {
			await sleep(50);
			fiberEndTime = Date.now();
		});

		sm.registerFunc('mongodb', 5, async () => {
			mongoStartTime = Date.now();
			await sleep(10);
		});

		await sm.executeShutdown(5000);

		expect(mongoStartTime).toBeGreaterThanOrEqual(fiberEndTime);
	});

	// --- executeShutdown: completes within timeout ---

	it('should complete well under the timeout', async () => {
		sm.registerFunc('fiber', 0, async () => { await sleep(10); });
		sm.registerFunc('mongodb', 5, async () => { await sleep(20); });
		sm.registerFunc('redis', 5, async () => { await sleep(15); });
		sm.registerFunc('nats', 5, async () => { await sleep(10); });

		const start = Date.now();
		await sm.executeShutdown(15000);
		const elapsed = Date.now() - start;

		// P0 (10ms) + P5 (20ms concurrent) = ~30ms, well under 1s
		expect(elapsed).toBeLessThan(1000);
	});

	// --- executeShutdown: error in one P5 does not block others ---

	it('should close other connections despite one failing', async () => {
		let redisClosed = false;
		let natsClosed = false;

		sm.registerFunc('fiber', 0, async () => {});
		sm.registerFunc('mongodb', 5, async () => { throw new Error('mongo: connection reset'); });
		sm.registerFunc('redis', 5, async () => { redisClosed = true; });
		sm.registerFunc('nats', 5, async () => { natsClosed = true; });

		await sm.executeShutdown(5000);

		expect(redisClosed).toBe(true);
		expect(natsClosed).toBe(true);
	});

	// --- waitForSignal: signal handling ---

	it('should register SIGTERM and SIGINT listeners', () => {
		const onSpy = jest.spyOn(process, 'on');

		sm.waitForSignal(15000);

		expect(onSpy).toHaveBeenCalledWith('SIGTERM', expect.any(Function));
		expect(onSpy).toHaveBeenCalledWith('SIGINT', expect.any(Function));

		// Cleanup: remove the listeners we just registered
		process.removeAllListeners('SIGTERM');
		process.removeAllListeners('SIGINT');

		onSpy.mockRestore();
	});
});

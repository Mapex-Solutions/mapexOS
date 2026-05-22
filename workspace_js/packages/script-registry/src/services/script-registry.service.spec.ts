import { ScriptRegistryService } from './script-registry.service';
import type { CompiledScripts, ScriptRegistryConfig } from '../types';

/**
 * Mock Script object for testing
 */
const createMockScript = (name: string) => ({
	run: jest.fn().mockResolvedValue(`result-${name}`),
} as unknown as CompiledScripts['decode']);

/**
 * Create mock CompiledScripts for testing
 */
const createMockScripts = (templateId: string): CompiledScripts => ({
	decode: createMockScript(`decode-${templateId}`),
	validation: createMockScript(`validation-${templateId}`),
	transform: createMockScript(`transform-${templateId}`),
});

describe('ScriptRegistryService', () => {
	let registry: ScriptRegistryService;

	beforeEach(() => {
		jest.useFakeTimers();
		registry = new ScriptRegistryService({
			maxEntries: 3,
			ttlMs: 1000, // 1 second for testing
			cleanupIntervalMs: 500,
		});
	});

	afterEach(() => {
		registry.stopCleanup();
		jest.useRealTimers();
	});

	describe('constructor', () => {
		it('should use default config when no config provided', () => {
			const defaultRegistry = new ScriptRegistryService();
			const config = defaultRegistry.getConfig();

			expect(config.maxEntries).toBe(1000);
			expect(config.ttlMs).toBe(900_000);
			expect(config.cleanupIntervalMs).toBe(60_000);
		});

		it('should merge provided config with defaults', () => {
			const customRegistry = new ScriptRegistryService({ maxEntries: 500 });
			const config = customRegistry.getConfig();

			expect(config.maxEntries).toBe(500);
			expect(config.ttlMs).toBe(900_000); // default
		});
	});

	describe('get/set', () => {
		it('should return undefined for non-existent entry', () => {
			const result = registry.get(1, 'template-1');
			expect(result).toBeUndefined();
		});

		it('should store and retrieve scripts', () => {
			const scripts = createMockScripts('template-1');
			registry.set(1, 'template-1', scripts);

			const result = registry.get(1, 'template-1');
			expect(result).toEqual(scripts);
		});

		it('should track stats for hits and misses', () => {
			const scripts = createMockScripts('template-1');

			// Miss
			registry.get(1, 'template-1');
			expect(registry.getStats().misses).toBe(1);
			expect(registry.getStats().hits).toBe(0);

			// Set and hit
			registry.set(1, 'template-1', scripts);
			registry.get(1, 'template-1');
			expect(registry.getStats().hits).toBe(1);
			expect(registry.getStats().misses).toBe(1);
		});

		it('should refresh TTL on access', () => {
			const scripts = createMockScripts('template-1');
			registry.set(1, 'template-1', scripts);

			// Advance time but not past TTL
			jest.advanceTimersByTime(800);

			// Access should refresh TTL
			const result = registry.get(1, 'template-1');
			expect(result).toEqual(scripts);

			// Advance another 800ms (total 1600ms from set, but only 800ms from last access)
			jest.advanceTimersByTime(800);

			// Should still be valid because TTL was refreshed
			const result2 = registry.get(1, 'template-1');
			expect(result2).toEqual(scripts);
		});

		it('should return undefined for expired entries', () => {
			const scripts = createMockScripts('template-1');
			registry.set(1, 'template-1', scripts);

			// Advance past TTL
			jest.advanceTimersByTime(1100);

			const result = registry.get(1, 'template-1');
			expect(result).toBeUndefined();
			expect(registry.getStats().ttlExpirations).toBe(1);
		});

		it('should isolate entries by isolateId', () => {
			const scripts1 = createMockScripts('template-1');
			const scripts2 = createMockScripts('template-1-v2');

			registry.set(1, 'template-1', scripts1);
			registry.set(2, 'template-1', scripts2);

			expect(registry.get(1, 'template-1')).toEqual(scripts1);
			expect(registry.get(2, 'template-1')).toEqual(scripts2);
		});
	});

	describe('LRU eviction', () => {
		it('should evict LRU entry when maxEntries is reached', () => {
			const scripts1 = createMockScripts('template-1');
			const scripts2 = createMockScripts('template-2');
			const scripts3 = createMockScripts('template-3');
			const scripts4 = createMockScripts('template-4');

			// Add 3 entries (maxEntries = 3)
			registry.set(1, 'template-1', scripts1);
			jest.advanceTimersByTime(10);
			registry.set(1, 'template-2', scripts2);
			jest.advanceTimersByTime(10);
			registry.set(1, 'template-3', scripts3);

			// Access template-1 to make it recently used
			jest.advanceTimersByTime(10);
			registry.get(1, 'template-1');

			// Add 4th entry - should evict template-2 (LRU)
			jest.advanceTimersByTime(10);
			registry.set(1, 'template-4', scripts4);

			expect(registry.get(1, 'template-1')).toEqual(scripts1); // Still exists
			expect(registry.get(1, 'template-2')).toBeUndefined(); // Evicted
			expect(registry.get(1, 'template-3')).toEqual(scripts3); // Still exists
			expect(registry.get(1, 'template-4')).toEqual(scripts4); // New entry

			expect(registry.getStats().evictions).toBe(1);
		});
	});

	describe('invalidate', () => {
		it('should remove template from all isolates', () => {
			const scripts = createMockScripts('template-1');

			registry.set(1, 'template-1', scripts);
			registry.set(2, 'template-1', scripts);
			registry.set(3, 'template-1', scripts);

			const removed = registry.invalidate('template-1');

			expect(removed).toBe(3);
			expect(registry.get(1, 'template-1')).toBeUndefined();
			expect(registry.get(2, 'template-1')).toBeUndefined();
			expect(registry.get(3, 'template-1')).toBeUndefined();
		});

		it('should return 0 if template not found', () => {
			const removed = registry.invalidate('non-existent');
			expect(removed).toBe(0);
		});
	});

	describe('clearIsolate', () => {
		it('should clear all entries for specific isolate', () => {
			const scripts1 = createMockScripts('template-1');
			const scripts2 = createMockScripts('template-2');

			registry.set(1, 'template-1', scripts1);
			registry.set(1, 'template-2', scripts2);
			registry.set(2, 'template-1', scripts1);

			registry.clearIsolate(1);

			expect(registry.get(1, 'template-1')).toBeUndefined();
			expect(registry.get(1, 'template-2')).toBeUndefined();
			expect(registry.get(2, 'template-1')).toEqual(scripts1); // Other isolate unaffected
		});
	});

	describe('cleanup', () => {
		it('should remove expired entries during cleanup cycle', () => {
			const scripts1 = createMockScripts('template-1');
			const scripts2 = createMockScripts('template-2');

			registry.set(1, 'template-1', scripts1);
			registry.startCleanup();

			// Advance past TTL
			jest.advanceTimersByTime(1100);

			// Add new entry after TTL
			registry.set(1, 'template-2', scripts2);

			// Trigger cleanup
			jest.advanceTimersByTime(500);

			// template-1 should be cleaned up, template-2 should remain
			const stats = registry.getStats();
			expect(stats.cleanupCycles).toBeGreaterThan(0);
			expect(registry.get(1, 'template-2')).toEqual(scripts2);
		});

		it('should increment cleanupCycles stat', () => {
			registry.startCleanup();

			expect(registry.getStats().cleanupCycles).toBe(0);

			jest.advanceTimersByTime(500);
			expect(registry.getStats().cleanupCycles).toBe(1);

			jest.advanceTimersByTime(500);
			expect(registry.getStats().cleanupCycles).toBe(2);
		});

		it('should not start multiple cleanup timers', () => {
			registry.startCleanup();
			registry.startCleanup(); // Should be ignored

			jest.advanceTimersByTime(500);
			expect(registry.getStats().cleanupCycles).toBe(1);
		});

		it('should stop cleanup timer', () => {
			registry.startCleanup();
			jest.advanceTimersByTime(500);
			expect(registry.getStats().cleanupCycles).toBe(1);

			registry.stopCleanup();
			jest.advanceTimersByTime(1000);
			expect(registry.getStats().cleanupCycles).toBe(1); // No more cycles
		});
	});

	describe('getStats', () => {
		it('should return accurate statistics', () => {
			const scripts = createMockScripts('template-1');

			// Initial stats
			let stats = registry.getStats();
			expect(stats.totalEntries).toBe(0);
			expect(stats.hits).toBe(0);
			expect(stats.misses).toBe(0);

			// Add entry
			registry.set(1, 'template-1', scripts);
			stats = registry.getStats();
			expect(stats.totalEntries).toBe(1);

			// Hit
			registry.get(1, 'template-1');
			stats = registry.getStats();
			expect(stats.hits).toBe(1);

			// Miss
			registry.get(1, 'template-2');
			stats = registry.getStats();
			expect(stats.misses).toBe(1);
		});
	});

	describe('getRegistryMap', () => {
		it('should return undefined for non-existent isolate', () => {
			const map = registry.getRegistryMap(999);
			expect(map).toBeUndefined();
		});

		it('should return the internal map for existing isolate', () => {
			const scripts = createMockScripts('template-1');
			registry.set(1, 'template-1', scripts);

			const map = registry.getRegistryMap(1);
			expect(map).toBeDefined();
			expect(map?.size).toBe(1);
		});
	});
});

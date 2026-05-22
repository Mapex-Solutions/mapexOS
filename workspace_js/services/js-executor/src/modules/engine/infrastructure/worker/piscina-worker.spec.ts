/**
 * Piscina Worker Unit Tests
 *
 * Tests the V8 isolate execution pipeline in the worker thread context.
 * Verifies: decode→validate→transform pipeline, OOM recovery, script caching,
 * error handling, context recycling, and event loss prevention.
 *
 * Uses mocked isolated-vm to test logic without native module dependency.
 */

// ─── State tracked across mock calls ─────────────────────────────────

let mockCurrentPayload: any = {};
let mockIsDisposed = false;
let compileCallCount = 0;

/**
 * Extract user code from the IIFE wrapper that wrapScriptCode() generates.
 * The wrapper adds result-checking boilerplate that contains 'throw new Error' —
 * we need to test only against the USER code, not the wrapper.
 */
function extractUserCode(wrappedCode: string): string {
	// wrapScriptCode wraps as: (function() { USER_CODE \n if (typeof result === 'undefined') ... })();
	const match = wrappedCode.match(/\(function\(\)\s*\{([\s\S]*?)if \(typeof result/);
	return match ? match[1] : wrappedCode;
}

/**
 * Simulates script behavior based on the user code that was compiled.
 * Each compiled script remembers the code and returns appropriate results.
 */
function createMockScript(code: string) {
	const userCode = extractUserCode(code);

	return {
		runSync: jest.fn().mockImplementation(() => {
			// Check for error-producing patterns in the USER code only
			if (userCode.includes('while(true)')) {
				throw new Error('Script execution timed out');
			}
			if (userCode.includes('undefined.property') || userCode.includes('undefinedVar.prop')) {
				throw new Error('TypeError: Cannot read properties of undefined');
			}
			if (userCode.includes('throw new Error')) {
				const errMatch = userCode.match(/throw new Error\("([^"]+)"\)/);
				throw new Error(errMatch ? errMatch[1] : 'Script error');
			}
			// Check that user code defines a result variable
			if (!userCode.includes('var result') && !userCode.includes('let result') && !userCode.includes('const result')) {
				throw new Error('Script must define a "result" variable with the return value');
			}

			// Return current payload as JSON (simulates successful execution)
			return JSON.stringify(mockCurrentPayload);
		}),
	};
}

const mockRelease = jest.fn();
const mockDispose = jest.fn();
const mockSetSync = jest.fn().mockImplementation((_name: string, value: any) => {
	if (_name === 'payload') {
		mockCurrentPayload = value;
	}
});

jest.mock('isolated-vm', () => ({
	__esModule: true,
	default: {
		Isolate: jest.fn().mockImplementation(() => ({
			get isDisposed() { return mockIsDisposed; },
			dispose: mockDispose,
			createContextSync: jest.fn().mockImplementation(() => ({
				global: { setSync: mockSetSync },
				release: mockRelease,
			})),
			compileScriptSync: jest.fn().mockImplementation((code: string) => {
				compileCallCount++;
				return createMockScript(code);
			}),
		})),
		ExternalCopy: jest.fn().mockImplementation((data: any) => ({
			copyInto: jest.fn().mockReturnValue(JSON.parse(JSON.stringify(data))),
		})),
	},
}));

jest.mock('worker_threads', () => ({
	workerData: {
		memoryLimitMb: 64,
		timeoutMs: 5000,
		contextRecycleInterval: 100, // High to avoid interference
		mapexValidatorCode: '',
	},
}));

import processEvent from './piscina-worker';
import type { PiscinaWorkerInput } from './types';

describe('Piscina Worker', () => {
	const makeInput = (overrides?: Partial<PiscinaWorkerInput>): PiscinaWorkerInput => ({
		rawPayload: { temperature: 25.3, humidity: 60, deviceId: 'sensor-001' },
		scripts: {
			decode: '',
			validation: '',
			transform: 'var result = { data: payload, processed: true };',
		},
		templateId: 'template-001',
		...overrides,
	});

	beforeEach(() => {
		jest.clearAllMocks();
		mockIsDisposed = false;
		mockCurrentPayload = {};
		compileCallCount = 0;
	});

	describe('Pipeline Execution', () => {
		it('should execute transform script and return success', async () => {
			const input = makeInput();
			const result = await processEvent(input);

			expect(result.success).toBe(true);
			expect(result.finalPayload).toBeDefined();
			expect(result.totalPipelineTime).toBeGreaterThanOrEqual(0);
			expect(result.isOOM).toBeUndefined();
		});

		it('should execute full decode→validate→transform pipeline', async () => {
			const input = makeInput({
				templateId: 'full-pipeline',
				scripts: {
					decode: 'var result = { decoded: true, raw: payload };',
					validation: 'var result = payload;',
					transform: 'var result = { final: payload.decoded, data: payload.raw };',
				},
			});

			const result = await processEvent(input);
			expect(result.success).toBe(true);
			expect(result.finalPayload).toBeDefined();
		});

		it('should skip empty script steps', async () => {
			const input = makeInput({
				templateId: 'skip-empty',
				scripts: {
					decode: '',
					validation: '',
					transform: 'var result = { original: payload.temperature };',
				},
			});

			const result = await processEvent(input);
			expect(result.success).toBe(true);
		});

		it('should handle undefined script steps', async () => {
			const input = makeInput({
				templateId: 'undef-scripts',
				scripts: {
					transform: 'var result = payload;',
				},
			});

			const result = await processEvent(input);
			expect(result.success).toBe(true);
		});
	});

	describe('Event Loss Prevention', () => {
		it('should ALWAYS return a result (never throw) for valid input', async () => {
			for (let i = 0; i < 10; i++) {
				const input = makeInput({
					templateId: 'event-loss-test',
					rawPayload: { seq: i },
				});
				const result = await processEvent(input);
				expect(result).toBeDefined();
				expect(typeof result.success).toBe('boolean');
			}
		});

		it('should return error result (not throw) for script failures', async () => {
			const input = makeInput({
				templateId: 'error-no-throw',
				scripts: {
					transform: 'var x = undefined.property;',
				},
			});

			const result = await processEvent(input);
			expect(result.success).toBe(false);
			expect(result.failedAt).toBe('transform');
			expect(result.error).toBeDefined();
		});

		it('should continue processing after a failure', async () => {
			// First: fail
			const failResult = await processEvent(makeInput({
				templateId: 'fail-first',
				scripts: { transform: 'var x = undefinedVar.prop;' },
			}));
			expect(failResult.success).toBe(false);

			// Second: succeed with different template
			const successResult = await processEvent(makeInput({
				templateId: 'succeed-after',
				scripts: { transform: 'var result = { recovered: true };' },
			}));
			expect(successResult.success).toBe(true);
		});
	});

	describe('Script Caching', () => {
		it('should cache compiled scripts per templateId', async () => {
			const input = makeInput({
				templateId: 'cache-hit',
				scripts: { transform: 'var result = { cached: true };' },
			});

			// First call: compiles the transform script
			await processEvent(input);
			const countAfterFirst = compileCallCount;

			// Second call: should use cached compiled script — no new compile
			await processEvent(input);
			const countAfterSecond = compileCallCount;

			expect(countAfterSecond).toBe(countAfterFirst);
		});

		it('should maintain separate caches for different templates', async () => {
			await processEvent(makeInput({
				templateId: 'cache-A',
				scripts: { transform: 'var result = { type: "A" };' },
			}));
			const countAfterA = compileCallCount;

			await processEvent(makeInput({
				templateId: 'cache-B',
				scripts: { transform: 'var result = { type: "B" };' },
			}));
			const countAfterB = compileCallCount;

			// B should trigger new compilation
			expect(countAfterB).toBeGreaterThan(countAfterA);
		});
	});

	describe('OOM Handling', () => {
		it('should detect OOM and return isOOM=true', async () => {
			const ivm = require('isolated-vm').default;

			// Force isolate recreation by marking the current one as disposed
			mockIsDisposed = true;

			// Provide OOM mock with its own disposed tracking
			let oomDisposed = false;
			ivm.Isolate.mockImplementationOnce(() => ({
				get isDisposed() { return oomDisposed; },
				dispose: mockDispose,
				createContextSync: jest.fn().mockImplementation(() => ({
					global: { setSync: mockSetSync },
					release: mockRelease,
				})),
				compileScriptSync: jest.fn().mockImplementation(() => ({
					runSync: jest.fn().mockImplementation(() => {
						oomDisposed = true;
						throw new Error('isolate was disposed');
					}),
				})),
			}));

			const result = await processEvent(makeInput({
				templateId: 'oom-detect',
				scripts: { transform: 'var result = 1;' },
			}));

			expect(result.success).toBe(false);
			expect(result.isOOM).toBe(true);
		});

		it('should recover after OOM on next call', async () => {
			const ivm = require('isolated-vm').default;

			// Force isolate recreation by marking the current one as disposed
			mockIsDisposed = true;

			// First call: OOM — mock with its own disposed tracking
			let oomDisposed = false;
			ivm.Isolate.mockImplementationOnce(() => ({
				get isDisposed() { return oomDisposed; },
				dispose: mockDispose,
				createContextSync: jest.fn().mockImplementation(() => ({
					global: { setSync: mockSetSync },
					release: mockRelease,
				})),
				compileScriptSync: jest.fn().mockImplementation(() => ({
					runSync: jest.fn().mockImplementation(() => {
						oomDisposed = true;
						throw new Error('OOM');
					}),
				})),
			}));

			const oomResult = await processEvent(makeInput({
				templateId: 'oom-first',
				scripts: { transform: 'var result = 1;' },
			}));
			expect(oomResult.success).toBe(false);
			expect(oomResult.isOOM).toBe(true);

			// Reset — next call creates new isolate with default mock and succeeds
			mockIsDisposed = false;

			const recoverResult = await processEvent(makeInput({
				templateId: 'oom-recover',
				scripts: { transform: 'var result = { recovered: true };' },
			}));
			expect(recoverResult.success).toBe(true);
		});
	});

	describe('Context Release', () => {
		it('should release context after each event', async () => {
			await processEvent(makeInput({ templateId: 'ctx-release' }));
			expect(mockRelease).toHaveBeenCalled();
		});

		it('should release context even on script failure', async () => {
			await processEvent(makeInput({
				templateId: 'ctx-release-fail',
				scripts: { transform: 'var x = undefined.property;' },
			}));
			expect(mockRelease).toHaveBeenCalled();
		});
	});

	describe('Error Handling', () => {
		it('should report decode failure with failedAt="decode"', async () => {
			const result = await processEvent(makeInput({
				templateId: 'decode-fail',
				scripts: {
					decode: 'throw new Error("decode failed");',
					transform: 'var result = payload;',
				},
			}));
			expect(result.success).toBe(false);
			expect(result.failedAt).toBe('decode');
		});

		it('should report validation failure with failedAt="validation"', async () => {
			const result = await processEvent(makeInput({
				templateId: 'validation-fail',
				scripts: {
					validation: 'throw new Error("validation failed");',
					transform: 'var result = payload;',
				},
			}));
			expect(result.success).toBe(false);
			expect(result.failedAt).toBe('validation');
		});

		it('should report transform failure with failedAt="transform"', async () => {
			const result = await processEvent(makeInput({
				templateId: 'transform-fail',
				scripts: { transform: 'throw new Error("transform failed");' },
			}));
			expect(result.success).toBe(false);
			expect(result.failedAt).toBe('transform');
		});

		it('should handle timeout scripts gracefully', async () => {
			const result = await processEvent(makeInput({
				templateId: 'timeout',
				scripts: { transform: 'while(true) {} var result = 1;' },
			}));
			expect(result.success).toBe(false);
			expect(result.failedAt).toBe('transform');
		});
	});

	describe('Payload Isolation', () => {
		it('should inject payload via ExternalCopy for each event', async () => {
			const ivm = require('isolated-vm').default;

			await processEvent(makeInput({
				templateId: 'isolation',
				rawPayload: { secret: 'event-1-data' },
			}));

			expect(ivm.ExternalCopy).toHaveBeenCalledWith(
				expect.objectContaining({ secret: 'event-1-data' })
			);
		});
	});
});

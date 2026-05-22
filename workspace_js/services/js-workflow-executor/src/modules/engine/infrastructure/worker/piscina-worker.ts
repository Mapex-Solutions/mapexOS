/**
 * Piscina Worker — V8 Isolate execution for workflow code nodes.
 *
 * Each worker thread owns:
 *   - Its own V8 Isolate (created, recovered on OOM, context recycled)
 *   - Per-cacheKey compiled script cache (one-time compile per script per worker)
 *
 * Receives a single PiscinaWorkerInput via piscina.run().
 * Returns PiscinaWorkerOutput — minimal structured clone.
 *
 * IMPORTANT: This file must be self-contained — no DI, no tsyringe, no path aliases.
 * Only imports: isolated-vm, worker_threads.
 */
'use strict';

import ivm from 'isolated-vm';
import { workerData } from 'worker_threads';

import type { PiscinaWorkerInput, PiscinaWorkerOutput, PiscinaWorkerConfig } from './types';

/** Config from workerData (set once at worker creation) */

const config: PiscinaWorkerConfig = workerData ?? {
	memoryLimitMb: 32,
	timeoutMs: 10_000,
	contextRecycleInterval: 10_000,
};

/** Per-worker V8 state */

let isolate: ivm.Isolate | null = null;
let execCount = 0;

/** Compiled scripts cached by cacheKey (one-time compile per worker per script) */
const scriptCache = new Map<string, ivm.Script>();

/** V8 Isolate Management */

/**
 * Initialize or recover V8 isolate.
 * On OOM (isolate.isDisposed), creates a fresh one and clears script cache.
 */
function ensureIsolate(): void {
	if (isolate && !isolate.isDisposed) {
		return;
	}

	// Isolate is disposed (OOM) or first call — create fresh
	isolate = new ivm.Isolate({ memoryLimit: config.memoryLimitMb });
	scriptCache.clear();
	execCount = 0;
}

/**
 * Wraps user script code with IIFE result extraction.
 *
 * Supports two patterns:
 *   1. `result = { output, statePatch }` — primary pattern (bare assignment creates global)
 *   2. `return value` — convenience pattern (captured via inner IIFE, wrapped as { output: value, statePatch: {} })
 */
function wrapScriptCode(scriptCode: string): string {
	return `
		(function() {
			var __ret = (function() {
				${scriptCode}
			})();
			if (typeof result !== 'undefined') {
				return JSON.stringify(result);
			}
			if (__ret !== undefined) {
				if (typeof __ret !== 'object' || __ret === null || Array.isArray(__ret)) {
					throw new Error('return value must be an object (e.g. return { key: value }). To use output in other nodes, return an object with named properties.');
				}
				return JSON.stringify({ output: __ret, statePatch: {} });
			}
			throw new Error('Script must define a "result" variable with { output, statePatch } or return a value');
		})();
	`;
}

/**
 * Compile and cache a workflow script.
 * Scripts are compiled once per worker per cacheKey.
 * Returns the compiled script and optionally freshly produced bytecode (on first compile without cached bytecode).
 */
function getOrCompileScript(
	cacheKey: string,
	script: string,
	cachedBytecode?: ArrayBuffer,
): { script: ivm.Script; newBytecode?: ArrayBuffer } {
	const cached = scriptCache.get(cacheKey);
	if (cached) return { script: cached };

	const wrappedCode = wrapScriptCode(script);

	// Try compiling with cached bytecode (fast path ~1-5ms)
	if (cachedBytecode) {
		const cachedEC = new ivm.ExternalCopy(cachedBytecode);
		try {
			const compiled = isolate!.compileScriptSync(wrappedCode, {
				filename: cacheKey,
				cachedData: cachedEC,
			});
			scriptCache.set(cacheKey, compiled);
			return { script: compiled };
		} catch {
			// Bytecode invalid (V8 version mismatch, corrupted) — fall through to fresh compile
		} finally {
			try { (cachedEC as any).dispose?.(); } catch { /* ignore */ }
		}
	}

	// Fresh compile — produce bytecode for caching (~10-50ms)
	const compiled = isolate!.compileScriptSync(wrappedCode, {
		filename: cacheKey,
		produceCachedData: true,
	});
	const newBytecode = (compiled as any).cachedData?.copy() as ArrayBuffer | undefined;

	scriptCache.set(cacheKey, compiled);
	return { script: compiled, newBytecode };
}

/**
 * Sanitizes isolated-vm errors for user-friendly messages.
 */
function sanitizeError(error: Error): string {
	let msg = (error.message || 'Unknown error')
		.replace(/\[<isolated-vm[^>]*>\]/g, '')
		.replace(/isolated-vm/gi, '')
		.replace(/compileWithCache/g, '')
		.replace(/\(<isolated-vm boundary>\)/g, '')
		.trim();

	// Script result hint
	if (msg.includes('Script must define a "result" variable') || msg.includes('Script must define "result"')) {
		return 'The script must either define result = { output, statePatch } or return a value';
	}

	return msg;
}

/** Main Export: Execute a workflow code node script */

/**
 * Executes a workflow code node script in an isolated V8 context.
 *
 * Injects: event, state, inputs, nodes into the V8 context.
 * Expects the script to define: result = { output, statePatch }
 *
 * @param input - PiscinaWorkerInput (script + cacheKey + context data)
 * @returns PiscinaWorkerOutput with result or error details
 */
async function executeWorkflowScript(input: PiscinaWorkerInput): Promise<PiscinaWorkerOutput> {
	const start = process.hrtime.bigint();

	try {
		ensureIsolate();
	} catch {
		return {
			success: false,
			error: 'Failed to create V8 isolate',
			isOOM: true,
			executionTime: elapsed(start),
		};
	}

	let context: ivm.Context | null = null;

	try {
		execCount++;

		// Get or compile the script (uses cached bytecode when available)
		const { script: compiled, newBytecode } = getOrCompileScript(
			input.cacheKey,
			input.script,
			input.cachedBytecode,
		);

		// Create a fresh context for this execution
		context = isolate!.createContextSync();
		const jail = context.global;

		// Inject workflow context globals
		jail.setSync('event', new ivm.ExternalCopy(input.event).copyInto());
		jail.setSync('state', new ivm.ExternalCopy(input.state).copyInto());
		jail.setSync('inputs', new ivm.ExternalCopy(input.inputs).copyInto());
		jail.setSync('nodes', new ivm.ExternalCopy(input.nodes).copyInto());

		// Execute the script
		const timeout = input.timeoutMs ?? config.timeoutMs;
		const jsonResult = compiled.runSync(context, { timeout });
		const parsed = JSON.parse(jsonResult);

		return {
			success: true,
			output: parsed.output,
			statePatch: parsed.statePatch,
			executionTime: elapsed(start),
			newBytecode,
		};
	} catch (err) {
		const isOOM = !!(isolate && isolate.isDisposed);
		if (isOOM) {
			isolate = null;
			scriptCache.clear();
		}

		const error = err instanceof Error ? sanitizeError(err) : String(err);
		return {
			success: false,
			error,
			isOOM,
			executionTime: elapsed(start),
		};
	} finally {
		// Release context to free memory
		if (context) {
			try { context.release(); } catch { /* already released */ }
		}

		// Recycle isolate periodically to prevent memory leaks
		if (execCount > 0 && execCount % config.contextRecycleInterval === 0) {
			recycleIsolate();
		}
	}
}

/**
 * Recycle isolate by disposing and clearing cache.
 * Next call to ensureIsolate() will create a fresh one.
 */
function recycleIsolate(): void {
	if (isolate && !isolate.isDisposed) {
		try { isolate.dispose(); } catch { /* ignore */ }
	}
	isolate = null;
	scriptCache.clear();
	execCount = 0;
}

/**
 * Calculates elapsed time in milliseconds from a high-resolution timestamp.
 *
 * @param start - High-resolution timestamp from process.hrtime.bigint()
 * @returns Elapsed time in milliseconds
 */
function elapsed(start: bigint): number {
	return Number(process.hrtime.bigint() - start) / 1e6;
}

export default executeWorkflowScript;

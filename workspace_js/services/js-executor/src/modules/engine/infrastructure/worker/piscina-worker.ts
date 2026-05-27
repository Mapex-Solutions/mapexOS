/**
 * Piscina Worker — V8 Isolate execution in a dedicated thread.
 *
 * Each worker thread owns:
 *   - Its own V8 Isolate (created, recovered on OOM, context recycled)
 *   - Per-template compiled script cache (one-time compile per template per worker)
 *
 * Receives a single PiscinaWorkerInput via piscina.run().
 * Returns PiscinaWorkerOutput — minimal structured clone.
 *
 * Based on the proven benchmark worker (scripts/piscina-bench/worker.js),
 * adapted for production with the full decode→validate→transform pipeline.
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
	mapexValidatorCode: '',
};

/** Per-worker V8 state */

let isolate: ivm.Isolate | null = null;
let execCount = 0;

/** Compiled scripts cached per template ID (one-time compile per worker per template) */
const scriptCache = new Map<string, CompiledTemplateScripts>();

interface CompiledTemplateScripts {
	decode?: ivm.Script;
	validation?: ivm.Script;
	validatorSetup?: ivm.Script;
	transform?: ivm.Script;
}

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
 * Replicates ScriptExecutor.wrapScriptCode() logic.
 */
function wrapScriptCode(scriptCode: string): string {
	return `
		(function() {
			${scriptCode}
			if (typeof result === 'undefined') {
				throw new Error('Script must define a "result" variable with the return value');
			}
			return JSON.stringify(result);
		})();
	`;
}

/**
 * Compile and cache scripts for a template.
 * Scripts are compiled once per worker per template — cached by templateId.
 */
function getOrCompileScripts(
	templateId: string,
	scripts: PiscinaWorkerInput['scripts'],
): CompiledTemplateScripts {
	let cached = scriptCache.get(templateId);
	if (cached) return cached;

	cached = {};

	if (scripts.decode?.trim()) {
		cached.decode = isolate!.compileScriptSync(wrapScriptCode(scripts.decode), { filename: 'decode' });
	}

	if (scripts.validation?.trim()) {
		cached.validation = isolate!.compileScriptSync(wrapScriptCode(scripts.validation), { filename: 'validation' });
		// Compile validator setup (MapexValidator injection)
		if (config.mapexValidatorCode) {
			cached.validatorSetup = isolate!.compileScriptSync(config.mapexValidatorCode, { filename: 'validatorSetup' });
		}
	}

	if (scripts.transform?.trim()) {
		cached.transform = isolate!.compileScriptSync(wrapScriptCode(scripts.transform), { filename: 'transform' });
	}

	scriptCache.set(templateId, cached);
	return cached;
}

/**
 * Sanitizes isolated-vm errors for user-friendly messages.
 * Replicates ScriptExecutor.sanitizeIsolatedVmError() logic.
 */
function sanitizeError(error: Error): string {
	let msg = (error.message || 'Unknown error')
		.replace(/\[<isolated-vm[^>]*>\]/g, '')
		.replace(/isolated-vm/gi, '')
		.replace(/compileWithCache/g, '')
		.replace(/\(<isolated-vm boundary>\)/g, '')
		.trim();

	// Script "result" variable hint
	if (msg.includes('Script must define a "result" variable')) {
		return 'The script must define a variable called "result" with the return value';
	}

	return msg;
}

/**
 * Run a compiled script in a context and return the parsed result.
 * Returns null on failure (error written to failRef).
 */
function runScript(
	script: ivm.Script,
	context: ivm.Context,
	scriptName: string,
): { data: any } | { error: string } {
	try {
		const jsonResult = script.runSync(context, { timeout: config.timeoutMs });
		// A non-string runSync return means the user's code exited the wrapper IIFE
		// early (e.g. typed `return X` instead of `const result = X`). The
		// JSON.parse on a non-string would surface a cryptic "undefined is not
		// valid JSON" — translate it here into something actionable.
		if (typeof jsonResult !== 'string') {
			return { error: `Error into script ${scriptName}: scripts must assign to a "result" variable (e.g. \`const result = X;\`) — do not use \`return X\` at the top level.` };
		}
		return { data: JSON.parse(jsonResult) };
	} catch (err) {
		const error = err instanceof Error ? err : new Error(String(err));
		return { error: `Error into script ${scriptName}: ${sanitizeError(error)}` };
	}
}

/** Main Export: Process a single event */

/**
 * Processes a single event through the decode→validate→transform pipeline.
 *
 * @param input - PiscinaWorkerInput (payload + scripts + templateId)
 * @returns PiscinaWorkerOutput with result or error details
 */
async function processEvent(input: PiscinaWorkerInput): Promise<PiscinaWorkerOutput> {
	const start = process.hrtime.bigint();

	try {
		ensureIsolate();
	} catch {
		return {
			success: false,
			error: 'Failed to create V8 isolate',
			isOOM: true,
			totalPipelineTime: elapsed(start),
		};
	}

	let context: ivm.Context | null = null;

	try {
		execCount++;

		// Get or compile scripts for this template
		const compiled = getOrCompileScripts(input.templateId, input.scripts);

		// Create a fresh context for this execution
		context = isolate!.createContextSync();
		const jail = context.global;

		// Inject payload
		jail.setSync('payload', new ivm.ExternalCopy(input.rawPayload).copyInto());

		let currentData = input.rawPayload;

		/** Decode */
		if (compiled.decode) {
			const result = runScript(compiled.decode, context, 'payloadDecode');
			if ('error' in result) {
				const isOOM = checkOOM();
				return { success: false, failedAt: 'decode', error: result.error, ...(isOOM && { isOOM: true }), totalPipelineTime: elapsed(start) };
			}
			currentData = result.data;
			jail.setSync('payload', new ivm.ExternalCopy(currentData).copyInto());
		}

		/** Validation */
		if (compiled.validation) {
			// Inject MapexValidator first
			if (compiled.validatorSetup) {
				compiled.validatorSetup.runSync(context, { timeout: config.timeoutMs });
			}

			const result = runScript(compiled.validation, context, 'payloadValidation');
			if ('error' in result) {
				const isOOM = checkOOM();
				return { success: false, failedAt: 'validation', error: result.error, ...(isOOM && { isOOM: true }), totalPipelineTime: elapsed(start) };
			}
			currentData = result.data;
			jail.setSync('payload', new ivm.ExternalCopy(currentData).copyInto());
		}

		/** Transform */
		if (compiled.transform) {
			const result = runScript(compiled.transform, context, 'payloadTransform');
			if ('error' in result) {
				const isOOM = checkOOM();
				return { success: false, failedAt: 'transform', error: result.error, ...(isOOM && { isOOM: true }), totalPipelineTime: elapsed(start) };
			}
			currentData = result.data;
		}

		return {
			success: true,
			finalPayload: currentData,
			totalPipelineTime: elapsed(start),
		};
	} catch (err) {
		const isOOM = !!(isolate && isolate.isDisposed);
		if (isOOM) {
			// Isolate died from OOM — clear state for recovery on next call
			isolate = null;
			scriptCache.clear();
		}

		const error = err instanceof Error ? sanitizeError(err) : String(err);
		return {
			success: false,
			error: `Pipeline error: ${error}`,
			isOOM,
			totalPipelineTime: elapsed(start),
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
 * Check if isolate was disposed (OOM) and clean up state.
 * Returns true if OOM was detected — caller MUST propagate isOOM flag.
 */
function checkOOM(): boolean {
	if (isolate && isolate.isDisposed) {
		isolate = null;
		scriptCache.clear();
		return true;
	}
	return false;
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

/** Calculates elapsed time in milliseconds from a hrtime bigint start. */
function elapsed(start: bigint): number {
	return Number(process.hrtime.bigint() - start) / 1e6;
}

export default processEvent;

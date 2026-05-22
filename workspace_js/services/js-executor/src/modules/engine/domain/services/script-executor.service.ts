import type { Script, Isolate, Context } from 'isolated-vm';

import type { SingleScriptResult, SanitizedError } from '@modules/engine/application/types';

import { ExternalCopy } from 'isolated-vm';

/**
 * Domain service for script execution in isolated-vm.
 *
 * This service contains PURE business logic for script execution without any I/O operations.
 * It is responsible for:
 * - Compiling scripts (optionally using cached bytecode)
 * - Running compiled scripts in a sandbox
 * - Sanitizing and formatting error messages
 *
 * @remarks
 * This service has NO external dependencies (no Redis, no file system, no network).
 * All I/O operations are handled by the Application layer (ScriptEngineService)
 * or Infrastructure layer (BytecodeCache).
 */
export class ScriptExecutor {
	/** Memory limit in MB for each isolate */
	private readonly memoryLimitMb: number;
	/** Timeout in ms for script execution */
	private readonly timeoutMs: number;

	constructor(memoryLimitMb = 32, timeoutMs = 10_000) {
		this.memoryLimitMb = memoryLimitMb;
		this.timeoutMs = timeoutMs;
	}

	/**
	 * Compiles a script, optionally using cached bytecode for faster compilation.
	 *
	 * @param isolate - The isolated-vm isolate instance
	 * @param code - The wrapped script code to compile
	 * @param cachedBytecode - Optional cached bytecode to speed up compilation
	 * @returns Promise resolving to { script, bytecode } where bytecode is present on fresh compile
	 */
	async compileScript(
		isolate: Isolate,
		code: string,
		cachedBytecode?: Buffer
	): Promise<{ script: Script; bytecode?: ArrayBuffer }> {
		if (cachedBytecode) {
			// Compile using cached bytecode
			const cachedEC = new ExternalCopy(this.bufferToArrayBuffer(cachedBytecode));
			try {
				const script = await isolate.compileScript(code, { cachedData: cachedEC });
				return { script };
			} catch {
				// Cache unusable, compile fresh and return new bytecode
				const script = await isolate.compileScript(code, { produceCachedData: true });
				const bytecode = (script as any).cachedData?.copy();
				return { script, bytecode };
			} finally {
				try {
					(cachedEC as any).dispose?.();
				} catch {
					// Ignore disposal errors
				}
			}
		}

		// Fresh compile
		const script = await isolate.compileScript(code, { produceCachedData: true });
		const bytecode = (script as any).cachedData?.copy();
		return { script, bytecode };
	}

	/**
	 * Creates a context for script execution and injects the payload.
	 *
	 * @param isolate - The isolate to create context in
	 * @param inputData - The data to inject as 'payload'
	 * @returns Promise resolving to the created Context
	 */
	async createContext(isolate: Isolate, inputData: any): Promise<Context> {
		const context = await isolate.createContext();
		const jail = context.global;
		await jail.set('payload', new ExternalCopy(inputData).copyInto());
		return context;
	}

	/**
	 * Runs a compiled script in a context and returns the result.
	 *
	 * @param scriptName - The name of the script for logging
	 * @param script - The compiled Script object
	 * @param context - The context with payload already injected
	 * @param startTime - The start time for execution time calculation
	 * @returns Promise resolving to the script execution result
	 */
	async runCompiledScript(
		scriptName: string,
		script: Script,
		context: Context,
		startTime: bigint
	): Promise<SingleScriptResult> {
		try {
			const jsonResult = await script.run(context, { timeout: this.timeoutMs });
			const endTime = process.hrtime.bigint();

			return {
				scriptName,
				data: JSON.parse(jsonResult),
				executionTime: this.diffMs(endTime, startTime),
				success: true,
			};
		} catch (error) {
			const endTime = process.hrtime.bigint();
			const sanitizedError = this.sanitizeIsolatedVmError(
				error instanceof Error ? error : new Error(String(error))
			);

			return {
				scriptName,
				data: null,
				executionTime: this.diffMs(endTime, startTime),
				success: false,
				error: `Error into script ${scriptName}: ${sanitizedError.userFriendlyMessage}`,
			};
		}
	}

	/**
	 * Runs a validation script by first executing the validator setup.
	 *
	 * @param scriptName - The name of the script for logging
	 * @param script - The compiled user Script object
	 * @param validatorSetupScript - The compiled MapexValidator setup Script
	 * @param context - The context with payload already injected
	 * @param startTime - The start time for execution time calculation
	 * @returns Promise resolving to the script execution result
	 */
	async runValidationScript(
		scriptName: string,
		script: Script,
		validatorSetupScript: Script,
		context: Context,
		startTime: bigint
	): Promise<SingleScriptResult> {
		try {
			// First execute validator setup
			await validatorSetupScript.run(context);
			// Then run user script
			return this.runCompiledScript(scriptName, script, context, startTime);
		} catch (error) {
			const endTime = process.hrtime.bigint();
			const sanitizedError = this.sanitizeIsolatedVmError(
				error instanceof Error ? error : new Error(String(error))
			);

			return {
				scriptName,
				data: null,
				executionTime: this.diffMs(endTime, startTime),
				success: false,
				error: `Error into script ${scriptName}: ${sanitizedError.userFriendlyMessage}`,
			};
		}
	}

	/**
	 * Wraps user script code with result extraction logic.
	 *
	 * Uses an IIFE (Immediately Invoked Function Expression) to isolate
	 * each script's scope, preventing variable collisions when multiple
	 * scripts run in the same Context (decode -> validation -> transform).
	 *
	 * @param scriptCode - The raw user script code
	 * @returns The wrapped code ready for compilation
	 */
	wrapScriptCode(scriptCode: string): string {
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
	 * Calculates time difference in milliseconds.
	 *
	 * @param endTime - End time from process.hrtime.bigint()
	 * @param startTime - Start time from process.hrtime.bigint()
	 * @returns Time difference in milliseconds (rounded)
	 */
	diffMs(endTime: bigint, startTime: bigint): number {
		const execTime = Number((Number(endTime - startTime) / 1e6).toFixed(3));
		return Math.round(execTime);
	}

	/**
	 * Converts Buffer to ArrayBuffer for isolated-vm compatibility.
	 *
	 * @param buf - Buffer to convert
	 * @returns The converted ArrayBuffer
	 */
	private bufferToArrayBuffer(buf: Buffer): ArrayBuffer {
		return buf.buffer.slice(buf.byteOffset, buf.byteOffset + buf.byteLength) as ArrayBuffer;
	}

	/**
	 * Sanitizes isolated-vm errors by removing internal references
	 * and creating user-friendly messages.
	 *
	 * @param error - The error to sanitize
	 * @returns The sanitized error with user-friendly message
	 */
	private sanitizeIsolatedVmError(error: Error): SanitizedError {
		const originalMessage = error.message || 'Erro desconhecido';

		// Remove isolated-vm specific references
		let cleanMessage = originalMessage
			.replace(/\[<isolated-vm[^>]*>\]/g, '')
			.replace(/isolated-vm/gi, '')
			.replace(/compileWithCache/g, '')
			.replace(/\(<isolated-vm boundary>\)/g, '')
			.trim();

		// Extract useful information from the error
		const syntaxErrorMatch = cleanMessage.match(/SyntaxError: (.+)/);
		const runtimeErrorMatch = cleanMessage.match(/ReferenceError: (.+)|TypeError: (.+)|Error: (.+)/);
		const lineColumnMatch = cleanMessage.match(/(\d+):(\d+)/);

		let errorType = 'Error';
		let coreMessage = cleanMessage;

		if (syntaxErrorMatch) {
			errorType = 'SyntaxError';
			coreMessage = syntaxErrorMatch[1];
		} else if (runtimeErrorMatch) {
			errorType = error.name || 'RuntimeError';
			coreMessage = runtimeErrorMatch[1] || runtimeErrorMatch[2] || runtimeErrorMatch[3];
		}

		const line = lineColumnMatch ? parseInt(lineColumnMatch[1]) : undefined;
		const column = lineColumnMatch ? parseInt(lineColumnMatch[2]) : undefined;

		return {
			type: errorType,
			message: coreMessage,
			line,
			column,
			userFriendlyMessage: this.createUserFriendlyMessage(errorType, coreMessage, line, column),
			originalError: originalMessage,
		};
	}

	/**
	 * Creates user-friendly messages based on error type.
	 *
	 * @param errorType - The type of error
	 * @param message - The error message
	 * @param line - Optional line number
	 * @param column - Optional column number
	 * @returns User-friendly error message with tips
	 */
	private createUserFriendlyMessage(
		errorType: string,
		message: string,
		line?: number,
		column?: number
	): string {
		let friendlyMessage = '';
		const locationInfo = line && column ? ` (line ${line}, column ${column})` : '';

		switch (errorType) {
			case 'SyntaxError':
				friendlyMessage = `Erro de sintaxe${locationInfo}: ${message}`;

				// Add specific hints based on the message
				if (message.includes('Invalid shorthand property initializer')) {
					friendlyMessage += '\n💡 Tip: Use ":" instead of "=" for object properties. Ex: { name: "value" }';
				} else if (message.includes('Unexpected token')) {
					friendlyMessage += '\n💡 Tip: Check for unclosed parentheses, braces, or quotes';
				} else if (message.includes('Unexpected end of input')) {
					friendlyMessage += '\n💡 Tip: You might be missing a closing code block';
				} else if (message.includes('Unexpected identifier')) {
					friendlyMessage += '\n💡 Tip: Check if you forgot to add commas or semicolons';
				}
				break;

			case 'ReferenceError':
				friendlyMessage = `Erro de referência${locationInfo}: ${message}`;
				friendlyMessage += '\n💡 Tip: Check if all variables were declared before use';
				break;

			case 'TypeError':
				friendlyMessage = `Erro de tipo${locationInfo}: ${message}`;
				friendlyMessage += '\n💡 Tip: Check if you are calling methods on the correct data type';
				break;

			default:
				friendlyMessage = `Erro${locationInfo}: ${message}`;

				// Check if it's the specific script error
				if (message.includes('Script must define a "result" variable')) {
					friendlyMessage = 'The script must define a variable called "result" with the return value';
					friendlyMessage += '\n💡 Tip: Add "const result = yourValue;" at the end of the script';
				}
				break;
		}

		return friendlyMessage;
	}
}

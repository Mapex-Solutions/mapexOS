import type { Logger } from '@mapexos/microservices';

import type { ScriptEngineServicePort } from '@modules/engine/application/ports';
import type { BytecodeCachePort } from '@modules/engine/application/ports';
import type { WorkflowScriptServicePort, WorkflowScriptInput, WorkflowScriptCallback, CallbackPublisherPort, ScriptSourceCachePort } from '../ports';
import type { PiscinaWorkerOutput } from '@modules/engine/infrastructure/worker';

/**
 * Application service for workflow code node script execution.
 *
 * Fetches script source from TieredCache, dispatches to Piscina worker (V8 isolate),
 * and publishes callback result to WORKFLOW-RESUME.
 *
 * Idempotency handled externally by Runtime MsgId dedup (publish-side) and CAS on NATS KV (consumer-side).
 */
export class WorkflowScriptService implements WorkflowScriptServicePort {
	constructor(
		private readonly logger: Logger,
		private readonly scriptSourceCache: ScriptSourceCachePort,
		private readonly bytecodeCache: BytecodeCachePort,
		private readonly scriptEngine: ScriptEngineServicePort,
		private readonly callbackPublisher: CallbackPublisherPort,
	) {}

	/**
	 * Executes a workflow code node script.
	 *
	 * @param input - The workflow script execution input
	 * @returns The Piscina worker output
	 */
	async execute(input: WorkflowScriptInput): Promise<PiscinaWorkerOutput> {
		const { orgId, workflowId, nodeId, instanceId, callbackSubject, executionToken } = input;

		// Fetch script source and bytecode in parallel
		const sourceCacheKey = `${orgId}/${workflowId}/scripts/${nodeId}`;
		const bytecodeCacheKey = this.bytecodeCache.buildCacheKey({ orgId, workflowId, nodeId });

		const [scriptSource, cachedBytecode] = await Promise.all([
			this.fetchScriptSource(sourceCacheKey),
			this.fetchBytecode(bytecodeCacheKey),
		]);

		if (!scriptSource) {
			const errorResult: PiscinaWorkerOutput = { success: false, error: `Script source not found for ${sourceCacheKey}` };
			const callback = this.buildCallback(errorResult, { instanceId, nodeId, executionToken }, 'SCRIPT_NOT_FOUND');
			await this.publishCallback(callbackSubject, callback);
			return errorResult;
		}

		// Dispatch to Piscina worker (with bytecode for faster cold compile)
		const result = await this.scriptEngine.runWorkflowScript({
			script: scriptSource,
			cacheKey: sourceCacheKey,
			event: input.eventPayload,
			state: input.state,
			inputs: input.inputs,
			nodes: input.nodeOutputs,
			cachedBytecode: cachedBytecode ?? undefined,
			timeoutMs: input.timeout ? input.timeout * 1000 : undefined,
		});

		// Store freshly produced bytecode for future cold starts (fire-and-forget)
		if (result.newBytecode) {
			this.bytecodeCache.set(bytecodeCacheKey, result.newBytecode).catch((err) => {
				const msg = err instanceof Error ? err.message : String(err);
				this.logger.warn(`[SERVICE:WorkflowScript] Failed to store bytecode for ${bytecodeCacheKey}: ${msg}`);
			});
		}

		// Build callback and publish
		const callback = this.buildCallback(result, { instanceId, nodeId, executionToken });
		await this.publishCallback(callbackSubject, callback);

		return result;
	}

	/**
	 * Invalidates cached script source + bytecode for specific nodes (L0 + L1 only).
	 * Called on FANOUT with granular nodeIds from Go workflow service.
	 * L2 (MinIO) cleanup is handled by the Go service.
	 *
	 * @param orgId - Organization ID
	 * @param definitionId - Workflow definition ID
	 * @param nodeIds - Node IDs to invalidate
	 */
	async invalidateNodes(orgId: string, definitionId: string, nodeIds: string[]): Promise<void> {
		// Invalidate script source from L0 + L1
		for (const nodeId of nodeIds) {
			const scriptKey = `${orgId}/${definitionId}/scripts/${nodeId}`;
			this.scriptSourceCache.invalidate(scriptKey);
		}

		// Invalidate bytecode from L0 + L1
		await this.bytecodeCache.invalidateNodes(orgId, definitionId, nodeIds);

		this.logger.info(
			`[SERVICE:WorkflowScript] Invalidated ${nodeIds.length} node(s) for definition ${definitionId} from L0/L1`
		);
	}

	/**
	 * Invalidates cached script source for a workflow definition (L0 + L1 only).
	 * Fallback when nodeIds are not available — relies on TTL expiry.
	 *
	 * @param orgId - Organization ID
	 * @param workflowId - Workflow definition ID
	 */
	async invalidateWorkflow(orgId: string, workflowId: string): Promise<void> {
		await this.bytecodeCache.invalidateWorkflow(orgId, workflowId);
		this.logger.info(`[SERVICE:WorkflowScript] Workflow ${workflowId} invalidation — L0/L1 entries will expire by TTL`);
	}

	/**
	 * Fetches script source from cache via port.
	 * Returns the script string or null if not found.
	 */
	private async fetchScriptSource(cacheKey: string): Promise<string | null> {
		try {
			const script = await this.scriptSourceCache.get(cacheKey);

			if (!script) {
				this.logger.warn(`[SERVICE:WorkflowScript] Script source cache miss for ${cacheKey}`);
				return null;
			}

			this.logger.debug(`[SERVICE:WorkflowScript] Script source loaded for ${cacheKey}`);
			return script;
		} catch (error) {
			const errorMessage = error instanceof Error ? error.message : String(error);
			this.logger.error(`[SERVICE:WorkflowScript] Failed to fetch script source for ${cacheKey}: ${errorMessage}`);
			return null;
		}
	}

	/**
	 * Fetches cached V8 bytecode from BytecodeCache.
	 * Returns the bytecode ArrayBuffer or null if not cached.
	 */
	private async fetchBytecode(cacheKey: string): Promise<ArrayBuffer | null> {
		try {
			const buffer = await this.bytecodeCache.get(cacheKey);
			if (!buffer) return null;

			this.logger.debug(`[SERVICE:WorkflowScript] Bytecode cache hit for ${cacheKey}`);
			return buffer.buffer.slice(buffer.byteOffset, buffer.byteOffset + buffer.byteLength) as ArrayBuffer;
		} catch {
			return null;
		}
	}

	/**
	 * Builds a WorkflowScriptCallback from the worker result.
	 *
	 * @param result - Piscina worker output
	 * @param identity - Instance, node, and execution token
	 * @param errorCode - Override error code (defaults to SCRIPT_ERROR)
	 */
	private buildCallback(
		result: PiscinaWorkerOutput,
		identity: Pick<WorkflowScriptCallback, 'instanceId' | 'nodeId' | 'executionToken'>,
		errorCode = 'SCRIPT_ERROR',
	): WorkflowScriptCallback {
		const base = { ...identity };

		if (result.success) {
			return { ...base, status: 'success', output: result.output, statePatch: result.statePatch };
		}

		return { ...base, status: 'error', error: { code: errorCode, message: result.error ?? 'Unknown script error' } };
	}

	/**
	 * Publishes callback result to WORKFLOW-RESUME via CallbackPublisherPort.
	 */
	private async publishCallback(callbackSubject: string, callback: WorkflowScriptCallback): Promise<void> {
		try {
			await this.callbackPublisher.publishCallback(callbackSubject, callback);
			if (callback.status === 'error') {
				this.logger.warn(
					`[SERVICE:WorkflowScript] Script error for instance ${callback.instanceId} node ${callback.nodeId}: ${callback.error?.code} — ${callback.error?.message}`
				);
			} else {
				this.logger.debug(
					`[SERVICE:WorkflowScript] Script success for instance ${callback.instanceId} node ${callback.nodeId}`
				);
			}
		} catch (error) {
			const errorMessage = error instanceof Error ? error.message : String(error);
			this.logger.error(
				`[SERVICE:WorkflowScript] Failed to publish callback for instance ${callback.instanceId}: ${errorMessage}`
			);
		}
	}
}

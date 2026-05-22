import type { PiscinaWorkerOutput } from '@modules/engine/infrastructure/worker';

/**
 * Input for workflow script execution.
 * Received from NATS WORKFLOW-JS-CODE consumer.
 */
export interface WorkflowScriptInput {
	/** Organization ID (multi-tenant isolation) */
	orgId: string;
	/** Hierarchical path key for multi-tenant DLQ filtering */
	pathKey: string;
	/** Workflow definition ID */
	workflowId: string;
	/** Node ID within the workflow */
	nodeId: string;
	/** Workflow instance ID */
	instanceId: string;
	/** NATS subject to publish callback result */
	callbackSubject: string;
	/** Execution token for callback validation (optional, backward compatible) */
	executionToken?: string;
	/** Trigger event payload */
	eventPayload: Record<string, any>;
	/** Current workflow instance state */
	state: Record<string, any>;
	/** External inputs provided at trigger time */
	inputs: Record<string, any>;
	/** Outputs from previous nodes */
	nodeOutputs: Record<string, any>;
	/** Script execution timeout in seconds (from node config). If not set, uses worker default. */
	timeout?: number;
}

/**
 * Callback result published to WORKFLOW-RESUME after script execution.
 */
export interface WorkflowScriptCallback {
	/** Workflow instance ID */
	instanceId: string;
	/** Node ID that executed */
	nodeId: string;
	/** Execution token echoed from input (optional, backward compatible) */
	executionToken?: string;
	/** Execution status */
	status: 'success' | 'error';
	/** Script output (when success) */
	output?: any;
	/** State patch to merge into instance state (when success) */
	statePatch?: Record<string, any>;
	/** Error details (when error) */
	error?: { code: string; message: string };
}

/**
 * Port interface for Workflow Script Service.
 *
 * Responsible for fetching script source from TieredCache,
 * dispatching to ScriptEngineService (Piscina workers),
 * and publishing callback result to WORKFLOW-RESUME.
 */
export interface WorkflowScriptServicePort {
	/**
	 * Executes a workflow code node script.
	 *
	 * @param input - The workflow script execution input
	 * @returns The Piscina worker output
	 */
	execute(input: WorkflowScriptInput): Promise<PiscinaWorkerOutput>;

	/**
	 * Invalidates cached script source + bytecode for specific nodes (L0 + L1 only).
	 * Called on FANOUT with granular nodeIds from Go workflow service.
	 * L2 (MinIO) cleanup is handled by the Go service.
	 *
	 * @param orgId - Organization ID
	 * @param definitionId - Workflow definition ID
	 * @param nodeIds - Node IDs to invalidate
	 */
	invalidateNodes(orgId: string, definitionId: string, nodeIds: string[]): Promise<void>;

	/**
	 * Invalidates cached script source for a workflow definition (L0 + L1 only).
	 * Fallback when nodeIds are not available — relies on TTL expiry.
	 *
	 * @param orgId - Organization ID
	 * @param workflowId - Workflow definition ID
	 */
	invalidateWorkflow(orgId: string, workflowId: string): Promise<void>;
}

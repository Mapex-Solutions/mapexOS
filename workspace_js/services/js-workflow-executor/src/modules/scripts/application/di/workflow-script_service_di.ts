import type { Logger } from '@mapexos/microservices';
import type { ScriptEngineServicePort, BytecodeCachePort } from '@modules/engine/application/ports';
import type {
	WorkflowScriptServicePort,
	CallbackPublisherPort,
	ScriptSourceCachePort,
} from '@modules/scripts/application/ports';
import { WorkflowScriptService } from '@modules/scripts/application/services';

/**
 * Dependencies required for WorkflowScriptService.
 * All fields are PORT interfaces (never concrete classes).
 */
export interface WorkflowScriptServiceDependencies {
	logger: Logger;
	scriptSourceCache: ScriptSourceCachePort;
	bytecodeCache: BytecodeCachePort;
	scriptEngine: ScriptEngineServicePort;
	callbackPublisher: CallbackPublisherPort;
}

/**
 * Factory constructs WorkflowScriptService with typed dependencies.
 * Returns the PORT interface, not the concrete class.
 */
export function createWorkflowScriptService(deps: WorkflowScriptServiceDependencies): WorkflowScriptServicePort {
	return new WorkflowScriptService(
		deps.logger,
		deps.scriptSourceCache,
		deps.bytecodeCache,
		deps.scriptEngine,
		deps.callbackPublisher,
	);
}

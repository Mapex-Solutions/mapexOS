import type { Logger } from '@mapexos/microservices';
import type { NatsBus } from '@mapexos/infrastructure';
import type { WorkflowScriptServicePort } from '@modules/scripts/application/ports';

/**
 * Dependencies for DefinitionInvalidateConsumer
 */
export interface DefinitionInvalidateConsumerDeps {
	/** NATS bus for messaging */
	natsBus: NatsBus;
	/** Logger instance */
	logger: Logger;
	/** Workflow script service for cache invalidation */
	scriptService: WorkflowScriptServicePort;
	/** Service name for fanout consumer identification */
	serviceName: string;
}

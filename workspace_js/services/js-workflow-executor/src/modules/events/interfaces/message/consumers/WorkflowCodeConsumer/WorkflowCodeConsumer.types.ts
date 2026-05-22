import type { Histogram, Counter } from 'prom-client';
import type { Logger, ConfigModule } from '@mapexos/microservices';
import type { NatsBus } from '@mapexos/infrastructure';
import type { WorkflowScriptServicePort } from '@modules/scripts/application/ports';

/**
 * Dependencies for WorkflowCodeConsumer
 */
export interface WorkflowCodeConsumerDeps {
	/** NATS bus for messaging */
	natsBus: NatsBus;
	/** Logger instance */
	logger: Logger;
	/** Workflow script service for processing messages */
	scriptService: WorkflowScriptServicePort;
	/** Config module for ENV-based consumer tuning */
	config: ConfigModule;
	/** Optional execution duration histogram */
	executionDuration?: Histogram;
	/** Optional executions total counter */
	executionsTotal?: Counter;
	/** Optional batch size histogram */
	batchSize?: Histogram;
}

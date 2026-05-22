import { container } from 'tsyringe';

import type { Logger, ConfigModule } from '@mapexos/microservices';
import type { NatsBus } from '@mapexos/infrastructure';

import { LOGGER_TOKEN, CONFIG_TOKEN } from '@mapexos/microservices';
import { NATS_BUS_TOKEN } from '@mapexos/infrastructure';

import type { WorkflowScriptServicePort } from '@modules/scripts/application/ports';
import { WORKFLOW_SCRIPT_SERVICE_TOKEN } from '@modules/scripts/module';

import type { WorkflowExecutorMetrics } from '@/bootstrap/metrics';
import { METRICS_TOKEN } from '@shared/constants';

import {
	initWorkflowCodeConsumer,
	initDefinitionInvalidateConsumer,
} from '@modules/events/interfaces/message';

import { streamName, subject } from '@shared/configuration/naming';

/** FANOUT stream name and subjects — resolves at module load via the canonical helpers. */
const FANOUT_STREAM = streamName('FANOUT', '');
const FANOUT_SUBJECTS = [subject('fanout', '') + '>'];

/**
 * InitListeners starts NATS event listeners for the events module.
 * Following workspace_go pattern — called during Phase 4 of module initialization.
 *
 * Consumer Types:
 * - Queue consumer (WorkflowCode): Load-balanced, one instance processes each message
 * - FANOUT consumer (DefinitionInvalidate): Broadcast, all instances receive all messages
 */
export async function initListeners(): Promise<void> {
	const logger = container.resolve<Logger>(LOGGER_TOKEN);
	const natsBus = container.resolve<NatsBus>(NATS_BUS_TOKEN);
	const scriptService = container.resolve<WorkflowScriptServicePort>(WORKFLOW_SCRIPT_SERVICE_TOKEN);
	const config = container.resolve<ConfigModule>(CONFIG_TOKEN);
	const serviceName = config.get('service_name') as string;

	const metrics = container.resolve<WorkflowExecutorMetrics>(METRICS_TOKEN);

	// Queue consumer: workflow code execution requests (load-balanced)
	void initWorkflowCodeConsumer({
		natsBus,
		logger,
		scriptService,
		config,
		executionDuration: metrics.executionDuration,
		executionsTotal: metrics.executionsTotal,
		batchSize: metrics.batchSize,
	});

	// Ensure FANOUT stream exists before subscribing
	await natsBus.ensureFanoutStream({
		name: FANOUT_STREAM,
		subjects: FANOUT_SUBJECTS,
		maxAge: 5 * 60 * 1000, // 5 minutes
		maxMsgs: 10000,
		description: 'FANOUT stream for cache invalidation events',
	});

	// FANOUT consumer: definition invalidation (broadcast to all pods)
	void initDefinitionInvalidateConsumer({ natsBus, logger, scriptService, serviceName });

	logger.debug('[MODULE:EVENTS] Listeners registered');
}

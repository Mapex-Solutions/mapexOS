import { container } from 'tsyringe';

import type { Logger } from '@mapexos/microservices';
import type { TieredCacheClient, NatsBus, MinIOClient } from '@mapexos/infrastructure';
import { LOGGER_TOKEN } from '@mapexos/microservices';
import { NATS_BUS_TOKEN } from '@mapexos/infrastructure';
import { NatsCallbackPublisher, TieredScriptSourceCacheAdapter } from '@modules/scripts/infrastructure';

import type { ScriptEngineServicePort, BytecodeCachePort } from '@modules/engine/application/ports';
import { TieredBytecodeCache } from '@modules/engine/infrastructure';
import { SCRIPT_ENGINE_SERVICE_TOKEN } from '@shared/constants';
import type { WorkflowScriptServicePort, ScriptSourceCachePort } from '@modules/scripts/application/ports';
import { WorkflowScriptService } from '@modules/scripts/application/services';

/**
 * Token for WorkflowScriptService in DI container
 */
export const WORKFLOW_SCRIPT_SERVICE_TOKEN = 'WorkflowScriptService';

/**
 * InitServices registers WorkflowScriptService in the DI container.
 * Called during Phase 2 of module initialization.
 */
export function initServices(): void {
	const logger = container.resolve<Logger>(LOGGER_TOKEN);
	logger.debug('[MODULE:SCRIPTS] Registering services');

	container.register<WorkflowScriptServicePort>(WORKFLOW_SCRIPT_SERVICE_TOKEN, {
		useFactory: (c) => {
			const bytecodeTieredCache = c.resolve<TieredCacheClient>('BytecodeCache');
			const minioClient = c.resolve<MinIOClient>('MinIOWorkflowsClient');
			const serviceLogger = c.resolve<Logger>(LOGGER_TOKEN);

			const scriptSourceCache: ScriptSourceCachePort = new TieredScriptSourceCacheAdapter(
				c.resolve<TieredCacheClient>('ScriptSourceCache'),
			);

			const bytecodeCache: BytecodeCachePort = new TieredBytecodeCache(
				bytecodeTieredCache,
				minioClient,
				serviceLogger,
			);

			const callbackPublisher = new NatsCallbackPublisher(c.resolve<NatsBus>(NATS_BUS_TOKEN));

			return new WorkflowScriptService(
				serviceLogger,
				scriptSourceCache,
				bytecodeCache,
				c.resolve<ScriptEngineServicePort>(SCRIPT_ENGINE_SERVICE_TOKEN),
				callbackPublisher,
			);
		},
	});

	logger.debug('[MODULE:SCRIPTS] Services registered');
}

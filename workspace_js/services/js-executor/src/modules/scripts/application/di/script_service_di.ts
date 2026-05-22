import type { Logger } from '@mapexos/microservices';
import type { ScriptEngineServicePort } from '@modules/engine/application/ports';
import type { ScriptServicePort } from '@modules/scripts/application/ports';
import type { AssetCachePort } from '@modules/scripts/application/ports/asset_cache_port';
import type { TemplateCachePort } from '@modules/scripts/application/ports/template_cache_port';
import type { EventPublisherPort } from '@modules/scripts/application/ports/event_publisher_port';
import type { ScriptServiceMetrics } from '@modules/scripts/application/types';
import { ScriptService } from '@modules/scripts/application/services';

/**
 * Dependencies required for ScriptService.
 *
 * Uses port interfaces (AssetCachePort, TemplateCachePort, EventPublisherPort)
 * instead of concrete types.
 */
export interface ScriptServiceDependencies {
	/** Logger instance for service logging. */
	logger: Logger;

	/** Script engine service for executing V8 scripts. */
	scriptEngine: ScriptEngineServicePort;

	/** Asset cache port for fetching asset read models. */
	assetCachePort: AssetCachePort;

	/** Template cache port for fetching template read models. */
	templateCachePort: TemplateCachePort;

	/** Event publisher port for publishing results, logs, and heartbeats. */
	eventPublisher: EventPublisherPort;

	/** Optional Prometheus metrics for instrumentation. */
	metrics?: ScriptServiceMetrics;
}

/**
 * Factory function to create a ScriptService instance.
 *
 * @param deps - Dependencies required by ScriptService
 * @returns ScriptServicePort implementation
 */
export function createScriptService(deps: ScriptServiceDependencies): ScriptServicePort {
	return new ScriptService(
		deps.logger,
		deps.scriptEngine,
		deps.assetCachePort,
		deps.templateCachePort,
		deps.eventPublisher,
		deps.metrics,
	);
}

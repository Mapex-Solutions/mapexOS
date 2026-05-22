import type { Logger, Application, ConfigModule } from '@mapexos/microservices';
import type { TieredCacheClient } from '@mapexos/infrastructure';

import type { ScriptServicePort } from '@modules/scripts/application/ports';
import type { EventPublisherPort } from '@modules/scripts/application/ports/event_publisher_port';
import type { AssetCachePort } from '@modules/scripts/application/ports/asset_cache_port';
import type { TemplateCachePort } from '@modules/scripts/application/ports/template_cache_port';
import type { ScriptEngineServicePort } from '@modules/engine/application/ports';

import { container } from 'tsyringe';
import { LOGGER_TOKEN, APP_TOKEN, CONFIG_TOKEN, apiKeyAuthMiddleware } from '@mapexos/microservices';
import { NATS_BUS_TOKEN } from '@mapexos/infrastructure';
import type { NatsBus } from '@mapexos/infrastructure';

import { createScriptService, type ScriptServiceDependencies } from '@modules/scripts/application/di';
import { NatsEventPublisherAdapter, AssetCacheAdapter, TemplateCacheAdapter } from '@modules/scripts/infrastructure';
import { SCRIPT_ENGINE_SERVICE_TOKEN } from '@modules/engine/module';
import { registerRoutes, registerInternalRoutes } from '@modules/scripts/interfaces/http';
import { SCRIPT_SERVICE_TOKEN, METRICS_TOKEN } from '@shared/constants';

import type { JsExecutorMetrics } from '@/bootstrap/metrics';

/** DI tokens for adapters — used by events/module.ts to inject into consumers */
export const EVENT_PUBLISHER_TOKEN = 'EventPublisher';
export const ASSET_CACHE_PORT_TOKEN = 'AssetCachePort';
export const TEMPLATE_CACHE_PORT_TOKEN = 'TemplateCachePort';

/**
 * InitServices registers all services and adapters in the DI container.
 * Following workspace_go pattern — called during Phase 2 of module initialization.
 */
export function initServices(): void {
	const logger = container.resolve<Logger>(LOGGER_TOKEN);
	logger.debug('[MODULE:SCRIPTS] Registering services and adapters');

	const metrics = container.resolve<JsExecutorMetrics>(METRICS_TOKEN);

	// Register infrastructure adapters (ports → concrete implementations)
	container.register<EventPublisherPort>(EVENT_PUBLISHER_TOKEN, {
		useFactory: (c) => {
			const natsBus = c.resolve<NatsBus>(NATS_BUS_TOKEN);
			return new NatsEventPublisherAdapter(natsBus);
		},
	});

	container.register<AssetCachePort>(ASSET_CACHE_PORT_TOKEN, {
		useFactory: (c) => {
			const cache = c.resolve<TieredCacheClient>('AssetCache');
			return new AssetCacheAdapter(cache);
		},
	});

	container.register<TemplateCachePort>(TEMPLATE_CACHE_PORT_TOKEN, {
		useFactory: (c) => {
			const cache = c.resolve<TieredCacheClient>('TemplateCache');
			return new TemplateCacheAdapter(cache);
		},
	});

	// Register ScriptService (uses ports, not concrete types)
	container.register<ScriptServicePort>(SCRIPT_SERVICE_TOKEN, {
		useFactory: (c) => {
			const deps: ScriptServiceDependencies = {
				logger: c.resolve<Logger>(LOGGER_TOKEN),
				scriptEngine: c.resolve<ScriptEngineServicePort>(SCRIPT_ENGINE_SERVICE_TOKEN),
				assetCachePort: c.resolve<AssetCachePort>(ASSET_CACHE_PORT_TOKEN),
				templateCachePort: c.resolve<TemplateCachePort>(TEMPLATE_CACHE_PORT_TOKEN),
				eventPublisher: c.resolve<EventPublisherPort>(EVENT_PUBLISHER_TOKEN),
				metrics: {
					eventsProcessed: metrics.eventsProcessed,
					eventDuration: metrics.eventDuration,
					payloadSize: metrics.payloadSize,
					assetCache: metrics.assetCache,
					templateCache: metrics.templateCache,
					heartbeatsPublished: metrics.heartbeatsPublished,
					heartbeatsSkipped: metrics.heartbeatsSkipped,
				},
			};
			return createScriptService(deps);
		},
	});

	logger.debug('[MODULE:SCRIPTS] Services and adapters registered');
}

/**
 * InitInterfaces registers HTTP routes.
 * Following workspace_go pattern — called during Phase 3 of module initialization.
 */
export function initInterfaces(): void {
	const logger = container.resolve<Logger>(LOGGER_TOKEN);
	const config = container.resolve<ConfigModule>(CONFIG_TOKEN);
	const app = container.resolve<Application>(APP_TOKEN);
	const scriptService = container.resolve<ScriptServicePort>(SCRIPT_SERVICE_TOKEN);

	app.use('/api/v1/scripts', registerRoutes(scriptService));

	const internalAPIKey = config.get('internal_api_key') as string;
	app.use('/internal/templates', apiKeyAuthMiddleware(internalAPIKey), registerInternalRoutes(scriptService));

	logger.debug('[MODULE:SCRIPTS] Routes registered (public + internal)');
}

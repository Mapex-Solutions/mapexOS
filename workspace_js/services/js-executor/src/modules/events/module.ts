import { container } from 'tsyringe';

import type { Logger, ConfigModule } from '@mapexos/microservices';
import type { NatsBus } from '@mapexos/infrastructure';

import { LOGGER_TOKEN, CONFIG_TOKEN } from '@mapexos/microservices';
import { NATS_BUS_TOKEN } from '@mapexos/infrastructure';

import type { ScriptServicePort } from '@modules/scripts/application/ports';
import type { AssetCachePort } from '@modules/scripts/application/ports/asset_cache_port';
import type { TemplateCachePort } from '@modules/scripts/application/ports/template_cache_port';

import { SCRIPT_SERVICE_TOKEN, METRICS_TOKEN, resolveAllTuning } from '@shared/constants';
import { ASSET_CACHE_PORT_TOKEN, TEMPLATE_CACHE_PORT_TOKEN } from '@modules/scripts/module';

import type { JsExecutorMetrics } from '@/bootstrap/metrics';

import {
	initJsExecuteConsumer,
	initMqttDataConsumer,
	initAssetInvalidateConsumer,
	initTemplateInvalidateConsumer,
} from '@modules/events/interfaces/message';

import { streamName, subject } from '@shared/configuration/naming';

/** FANOUT stream name and subjects — resolves at module load via the canonical helpers. */
const FANOUT_STREAM = streamName('FANOUT', '');
const FANOUT_SUBJECTS = [subject('fanout', '') + '>'];

/**
 * InitListeners starts NATS event listeners for the events module.
 *
 * Consumer Types:
 * - Queue consumers (JsExecute, MqttData): Load-balanced, full lifecycle
 * - FANOUT consumers (AssetInvalidate, TemplateInvalidate): Broadcast, cache invalidation
 */
export async function initListeners(): Promise<void> {
	const logger = container.resolve<Logger>(LOGGER_TOKEN);
	const natsBus = container.resolve<NatsBus>(NATS_BUS_TOKEN);
	const scriptService = container.resolve<ScriptServicePort>(SCRIPT_SERVICE_TOKEN);
	const assetCachePort = container.resolve<AssetCachePort>(ASSET_CACHE_PORT_TOKEN);
	const templateCachePort = container.resolve<TemplateCachePort>(TEMPLATE_CACHE_PORT_TOKEN);
	const config = container.resolve<ConfigModule>(CONFIG_TOKEN);
	const serviceName = config.get('service_name') as string;

	const metrics = container.resolve<JsExecutorMetrics>(METRICS_TOKEN);
	const tuning = resolveAllTuning(config);

	logger.info({ tuning }, '[MODULE:EVENTS] Auto-tuning resolved from CPU_LIMIT');

	// Queue consumers (validate + delegate to service)
	void initJsExecuteConsumer({
		natsBus, logger, scriptService, config,
		batchSize: metrics.batchSize,
		eventsProcessed: metrics.eventsProcessed,
		payloadSize: metrics.payloadSize,
	});

	void initMqttDataConsumer({
		natsBus, logger, scriptService, config,
		batchSize: metrics.batchSize,
		eventsProcessed: metrics.eventsProcessed,
		payloadSize: metrics.payloadSize,
	});

	// Ensure FANOUT stream exists before subscribing
	await natsBus.ensureFanoutStream({
		name: FANOUT_STREAM,
		subjects: FANOUT_SUBJECTS,
		maxAge: 5 * 60 * 1000,
		maxMsgs: 10000,
		description: 'FANOUT stream for cache invalidation events',
	});

	// FANOUT consumers (broadcast, cache invalidation via adapters)
	void initAssetInvalidateConsumer({ natsBus, logger, assetCachePort, serviceName });
	void initTemplateInvalidateConsumer({ natsBus, logger, templateCachePort, serviceName });

	logger.debug('[MODULE:EVENTS] Listeners registered');
}

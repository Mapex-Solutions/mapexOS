import type { Histogram, Counter } from 'prom-client';
import type { Logger, ConfigModule } from '@mapexos/microservices';
import type { NatsBus } from '@mapexos/infrastructure';
import type { ScriptServicePort } from '@modules/scripts/application/ports';

/**
 * Dependencies for JsExecuteConsumer.
 *
 * Consumer validates and delegates to ScriptService — no engine, no publisher, no cache.
 */
export interface JsExecuteConsumerDeps {
	/** NATS bus for consumer registration */
	natsBus: NatsBus;
	/** Logger instance */
	logger: Logger;
	/** Script service for batch processing */
	scriptService: ScriptServicePort;
	/** Config module for ENV-based consumer tuning */
	config: ConfigModule;
	/** Optional batch size histogram for Prometheus metrics */
	batchSize?: Histogram;
	/** Optional events processed counter for Prometheus metrics */
	eventsProcessed?: Counter;
	/** Optional payload size histogram for Prometheus metrics */
	payloadSize?: Histogram;
}

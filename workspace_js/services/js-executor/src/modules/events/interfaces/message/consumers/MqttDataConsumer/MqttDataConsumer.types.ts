import type { Histogram, Counter } from 'prom-client';
import type { Logger, ConfigModule } from '@mapexos/microservices';
import type { NatsBus } from '@mapexos/infrastructure';
import type { ScriptServicePort } from '@modules/scripts/application/ports';

/**
 * Dependencies for MqttDataConsumer.
 *
 * Consumer validates and delegates to ScriptService — no engine, no publisher.
 */
export interface MqttDataConsumerDeps {
	/** NATS bus for consumer registration */
	natsBus: NatsBus;
	/** Logger instance */
	logger: Logger;
	/** Script service for batch processing */
	scriptService: ScriptServicePort;
	/** Config module for ENV-based consumer tuning */
	config: ConfigModule;
	/** Optional batch size histogram */
	batchSize?: Histogram;
	/** Optional events processed counter */
	eventsProcessed?: Counter;
	/** Optional payload size histogram */
	payloadSize?: Histogram;
}

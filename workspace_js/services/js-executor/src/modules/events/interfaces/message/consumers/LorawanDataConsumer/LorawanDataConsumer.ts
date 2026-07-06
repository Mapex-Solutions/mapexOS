import type { Message } from '@mapexos/infrastructure';
import type { LorawanDataConsumerDeps } from './LorawanDataConsumer.types';

import { DeliverPolicy } from '@mapexos/infrastructure';

import { SERVICE_NAME, SERVICE_TYPE, DEFAULT_RETRY_POLICY, resolveConsumerConfig } from '@shared/constants';
import { LORAWAN_DATA_STREAM, LORAWAN_DATA_SUBJECT, LORAWAN_DATA_DURABLE, LORAWAN_DATA_EVENT_TYPE } from './constants';

/**
 * LorawanData consumer — receives LoRaWAN uplinks and delegates to ScriptService.
 *
 * Flow: receive → call service.handleLorawanBatch() → ACK/Nack per result.
 *
 * @param deps - Consumer dependencies
 */
export async function initLorawanDataConsumer(deps: LorawanDataConsumerDeps): Promise<void> {
	const { natsBus, logger, scriptService, config } = deps;
	const consumerConfig = resolveConsumerConfig(config);

	logger.info({ config: consumerConfig }, '[CONSUMER:LorawanData] Initializing');

	await natsBus.startConsumer({
		stream: LORAWAN_DATA_STREAM,
		subject: LORAWAN_DATA_SUBJECT,
		durable: LORAWAN_DATA_DURABLE,
		deliverPolicy: DeliverPolicy.All,
		batchSize: consumerConfig.batchSize,
		fetchTimeout: consumerConfig.fetchTimeout,
		maxAckPending: consumerConfig.maxAckPending,
		retryPolicy: DEFAULT_RETRY_POLICY,
		dlqPolicy: {
			serviceName: SERVICE_NAME,
			serviceType: SERVICE_TYPE,
			eventType: LORAWAN_DATA_EVENT_TYPE,
		},
		batchMessageHandlerV2: async (messages: Message[]) => {
			deps.batchSize?.observe(messages.length);

			const results = await scriptService.handleLorawanBatch(messages);

			for (const result of results) {
				const msg = messages[result.index];

				if (result.success) {
					msg.ack();
					deps.eventsProcessed?.inc({ consumer: 'lorawan_data', status: 'success' });
				} else if (result.isPermanent) {
					await msg.reject(result.error ?? 'Invalid message');
					deps.eventsProcessed?.inc({ consumer: 'lorawan_data', status: 'rejected' });
				} else if (result.isOOM) {
					msg.nack(new Error(result.error ?? 'V8 OOM')).catch(() => {});
					deps.eventsProcessed?.inc({ consumer: 'lorawan_data', status: 'failure' });
				} else {
					msg.nack(new Error(result.error ?? 'Processing error')).catch(() => {});
					deps.eventsProcessed?.inc({ consumer: 'lorawan_data', status: 'failure' });
				}
			}
		},
	});

	logger.info('[CONSUMER:LorawanData] Initialized');
}

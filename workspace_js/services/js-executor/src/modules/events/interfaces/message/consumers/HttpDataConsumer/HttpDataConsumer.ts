import type { Message } from '@mapexos/infrastructure';
import type { HttpDataConsumerDeps } from './HttpDataConsumer.types';

import { DeliverPolicy } from '@mapexos/infrastructure';

import { SERVICE_NAME, SERVICE_TYPE, DEFAULT_RETRY_POLICY, resolveConsumerConfig } from '@shared/constants';
import { HTTP_DATA_STREAM, HTTP_DATA_SUBJECT, HTTP_DATA_DURABLE, HTTP_DATA_EVENT_TYPE } from './constants';

/**
 * HttpData consumer — receives HTTP datasource events and delegates to ScriptService.
 *
 * Flow: receive → call service.handleHttpBatch() → ACK/Nack per result.
 *
 * @param deps - Consumer dependencies
 */
export async function initHttpDataConsumer(deps: HttpDataConsumerDeps): Promise<void> {
	const { natsBus, logger, scriptService, config } = deps;
	const consumerConfig = resolveConsumerConfig(config);

	logger.info({ config: consumerConfig }, '[CONSUMER:HttpData] Initializing');

	await natsBus.startConsumer({
		stream: HTTP_DATA_STREAM,
		subject: HTTP_DATA_SUBJECT,
		durable: HTTP_DATA_DURABLE,
		deliverPolicy: DeliverPolicy.All,
		batchSize: consumerConfig.batchSize,
		fetchTimeout: consumerConfig.fetchTimeout,
		maxAckPending: consumerConfig.maxAckPending,
		retryPolicy: DEFAULT_RETRY_POLICY,
		dlqPolicy: {
			serviceName: SERVICE_NAME,
			serviceType: SERVICE_TYPE,
			eventType: HTTP_DATA_EVENT_TYPE,
		},
		batchMessageHandlerV2: async (messages: Message[]) => {
			deps.batchSize?.observe(messages.length);

			const results = await scriptService.handleHttpBatch(messages);

			for (const result of results) {
				const msg = messages[result.index];

				if (result.success) {
					msg.ack();
					deps.eventsProcessed?.inc({ consumer: 'http_data', status: 'success' });
				} else if (result.isPermanent) {
					await msg.reject(result.error ?? 'Invalid message');
					deps.eventsProcessed?.inc({ consumer: 'http_data', status: 'rejected' });
				} else if (result.isOOM) {
					msg.nack(new Error(result.error ?? 'V8 OOM')).catch(() => {});
					deps.eventsProcessed?.inc({ consumer: 'http_data', status: 'failure' });
				} else {
					msg.nack(new Error(result.error ?? 'Processing error')).catch(() => {});
					deps.eventsProcessed?.inc({ consumer: 'http_data', status: 'failure' });
				}
			}
		},
	});

	logger.info('[CONSUMER:HttpData] Initialized');
}

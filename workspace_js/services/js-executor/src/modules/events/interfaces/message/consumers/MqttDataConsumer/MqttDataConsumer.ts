import type { Message } from '@mapexos/infrastructure';
import type { MqttDataConsumerDeps } from './MqttDataConsumer.types';

import { DeliverPolicy } from '@mapexos/infrastructure';

import { SERVICE_NAME, SERVICE_TYPE, DEFAULT_RETRY_POLICY, resolveConsumerConfig } from '@shared/constants';
import { MQTT_DATA_STREAM, MQTT_DATA_SUBJECT, MQTT_DATA_DURABLE, MQTT_DATA_EVENT_TYPE } from './constants';

/**
 * MqttData consumer — receives MQTT telemetry and delegates to ScriptService.
 *
 * Flow: receive → call service.handleMqttBatch() → ACK/Nack per result.
 *
 * @param deps - Consumer dependencies
 */
export async function initMqttDataConsumer(deps: MqttDataConsumerDeps): Promise<void> {
	const { natsBus, logger, scriptService, config } = deps;
	const consumerConfig = resolveConsumerConfig(config);

	logger.info({ config: consumerConfig }, '[CONSUMER:MqttData] Initializing');

	await natsBus.startConsumer({
		stream: MQTT_DATA_STREAM,
		subject: MQTT_DATA_SUBJECT,
		durable: MQTT_DATA_DURABLE,
		deliverPolicy: DeliverPolicy.All,
		batchSize: consumerConfig.batchSize,
		fetchTimeout: consumerConfig.fetchTimeout,
		maxAckPending: consumerConfig.maxAckPending,
		retryPolicy: DEFAULT_RETRY_POLICY,
		dlqPolicy: {
			serviceName: SERVICE_NAME,
			serviceType: SERVICE_TYPE,
			eventType: MQTT_DATA_EVENT_TYPE,
		},
		batchMessageHandlerV2: async (messages: Message[]) => {
			deps.batchSize?.observe(messages.length);

			const results = await scriptService.handleMqttBatch(messages);

			for (const result of results) {
				const msg = messages[result.index];

				if (result.success) {
					msg.ack();
					deps.eventsProcessed?.inc({ consumer: 'mqtt_data', status: 'success' });
				} else if (result.isPermanent) {
					await msg.reject(result.error ?? 'Invalid message');
					deps.eventsProcessed?.inc({ consumer: 'mqtt_data', status: 'rejected' });
				} else if (result.isOOM) {
					msg.nack(new Error(result.error ?? 'V8 OOM')).catch(() => {});
					deps.eventsProcessed?.inc({ consumer: 'mqtt_data', status: 'failure' });
				} else {
					msg.nack(new Error(result.error ?? 'Processing error')).catch(() => {});
					deps.eventsProcessed?.inc({ consumer: 'mqtt_data', status: 'failure' });
				}
			}
		},
	});

	logger.info('[CONSUMER:MqttData] Initialized');
}

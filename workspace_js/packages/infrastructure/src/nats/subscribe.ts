import type { NatsClient, NatsSubscribeInternalOptions } from './types';

import { consumerOpts } from 'nats';
import { getInfraLogger } from '../logger';
import { AckPolicy, DeliverPolicy, StorageType, RetentionPolicy } from '@src/nats';

/**
 * Ensures that a consumer exists for the given stream and durable name.
 * If the consumer does not exist, it creates a new one with the specified options.
 *
 * @param client - The NatsClient instance used to interact with the NATS server.
 * @param options - The options for the subscription, including stream, durable name, and other settings.
 * @throws Will throw an error if the stream or durable name is not provided, or if there is a failure in creating or retrieving the consumer.
 * @returns A promise that resolves when the consumer is successfully created or retrieved.
 */
async function createOrGetConsumer(client: NatsClient, options: NatsSubscribeInternalOptions) {
	if (!options.stream || !options.durable) {
		throw new Error('Stream and durable name are required for pull subscriptions');
	}

	try {
		const jsm = await client.nc.jetstreamManager();

		// First, ensure the stream exists
		try {
			await jsm.streams.info(options.stream);
		} catch {
			// Stream doesn't exist, create it
			await jsm.streams.add({
				name: options.stream,
				subjects: [options.subject],

				// Add other stream configuration as needed
				retention: RetentionPolicy.Workqueue,
				storage: StorageType.File,
			});
		}

		// Now try to get existing consumer
		try {
			await jsm.consumers.info(options.stream, options.durable);
			return;
		} catch {
			// Consumer doesn't exist, create it
			await jsm.consumers.add(options.stream, {
				durable_name: options.durable,
				ack_policy: (options.ackPolicy ?? AckPolicy.Explicit),
				deliver_policy: (options.deliverPolicy ?? DeliverPolicy.New),
				filter_subject: options.filterSubject ?? options.subject,
				max_ack_pending: 128,
			});
		}
	} catch (error) {
		throw new Error(`Failed to create/get stream/consumer: ${error}`);
	}
}

/**
 * Handles a push subscription to a NATS subject, processing messages as they arrive.
 *
 * @param client - The NatsClient instance used to interact with the NATS server.
 * @param options - The options for the subscription, including subject, durable name, queue group, and message handler.
 * @returns A function that, when called, will unsubscribe from the NATS subject and clean up resources.
 */
async function handlePushSubscription(client: NatsClient, options: NatsSubscribeInternalOptions) {
	const co = consumerOpts();
	
	if (options.durable) co.durable(options.durable);
	if (options.queueGroup) co.queue(options.queueGroup);
	if (options.stream) co.bindStream(options.stream);
	co.manualAck();

	const subscription = await client.js.subscribe(options.subject, co);

	const processMessages = async () => {
		for await (const message of subscription) {
			try {
				await options.handler(message.data);
				message.ack();
			} catch (error) {
				getInfraLogger().error({ err: error }, '[INFRA:NATS] Push handler error');
				try {
					message.term?.();
				} catch {
				}
			}
		}
	};

	processMessages().catch(() => {
	});

	return async () => {
		try {
			await (subscription as any).drain?.();
		} catch {
			(subscription as any).unsubscribe?.();
		}
	};
}

/**
 * Handles a pull subscription to a NATS subject, processing messages in batches.
 * This function ensures that a consumer is created or retrieved for the specified stream and durable name,
 * and then processes messages in a loop until the subscription is stopped.
 *
 * @param client - The NatsClient instance used to interact with the NATS server.
 * @param options - The options for the subscription, including stream, durable name, subject, and message handler.
 *                  - `stream`: The name of the stream to subscribe to.
 *                  - `durable`: The durable name for the consumer.
 *                  - `subject`: The subject to subscribe to.
 *                  - `handler`: A function to handle incoming messages.
 *                  - `maxMessages`: Optional. The maximum number of messages to process in a batch. Defaults to 10.
 *                  - `expires`: Optional. The time in milliseconds to wait for messages before timing out. Defaults to 5000.
 * @returns A function that, when called, will stop the pull subscription and clean up resources.
 */
async function handlePullSubscription(client: NatsClient, options: NatsSubscribeInternalOptions) {
  await createOrGetConsumer(client, options);

  const consumers = (client.js as any).consumers;
  if (!consumers) throw new Error('Consumer API not available');

  const consumer = await consumers.get(options.stream, options.durable);

  // Read pacing configuration (with safe defaults)
  const batchSize = options.maxMessages ?? 50;
  const timeoutMs = options.expires ?? 5000;

  let isRunning = true;

  const sleep = (ms: number) => new Promise<void>(r => setTimeout(r, ms));

  const processMessages = async () => {
    while (isRunning) {
      try {
        const iterator = await consumer.fetch({
          max_messages: batchSize,
          expires: timeoutMs,
        });

        for await (const message of iterator) {
          if (!isRunning) break;
          try {
            await options.handler(message.data, message.subject);
            message.ack();
          } catch (error) {
            getInfraLogger().error({ err: error }, '[INFRA:NATS] Pull handler error');

            // If you want retry behavior, remove term() to let ack_wait trigger redelivery
            try { message.term?.(); } catch { /* noop */ }
          }
        }
      } catch (error) {
        if (isRunning) {
          getInfraLogger().warn({ err: error }, '[INFRA:NATS] Pull subscription error, retrying...');
          await sleep(1000);
        }
      }
    }
  };

  processMessages().catch(() => { /* swallow loop rejection on shutdown */ });

  return async () => {
    isRunning = false;
    try { await consumer.close?.(); } catch { /* noop */ }
  };
}


/**
 * Subscribes to a NATS subject using either a push or pull subscription based on the provided options.
 *
 * @param client - The NatsClient instance used to interact with the NATS server.
 * @param options - The options for the subscription, including:
 *                  - `subject`: The subject to subscribe to. This is required.
 *                  - `stream`: The name of the stream to subscribe to. This is required for both push and pull subscriptions.
 *                  - `durable`: The durable name for the consumer. This is required for both push and pull subscriptions.
 *                  - `pull`: A boolean indicating whether to use a pull subscription. If true, a pull subscription is used; otherwise, a push subscription is used.
 *                  - Additional options specific to the type of subscription.
 * @returns A promise that resolves to a function. When this function is called, it will unsubscribe from the NATS subject and clean up resources.
 * @throws Will throw an error if the subject, stream, or durable name is not provided.
 */
export async function subscribe(client: NatsClient, options: NatsSubscribeInternalOptions) {
	if (!options.subject) {
		throw new Error('Subject is required');
	}

	if (!options.stream) {
		throw new Error('Stream name is required for both push and pull subscriptions');
	}

	if (!options.durable) {
		throw new Error('Durable name is required for both push and pull subscriptions');
	}

	return options.pull
		? handlePullSubscription(client, options)
		: handlePushSubscription(client, options);
}
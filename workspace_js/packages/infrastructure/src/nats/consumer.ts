import type { NatsClient, ConsumerOptions, DLQMessage, ConsumerHandle } from './types';
import type { JsMsg } from 'nats';

import { AckPolicy, DeliverPolicy, StorageType, RetentionPolicy } from 'nats';

import { Message } from './types';
import { DLQ_STREAM, DLQ_SUBJECT, DLQ_MAX_AGE_NANOS } from './constants';
import { sleep } from '@mapexos/utils';
import { getInfraLogger } from '../logger';

/** 10 minutes in nanoseconds — default dedup window for new streams */
const DEFAULT_DUPLICATE_WINDOW_NS = 10 * 60 * 1_000_000_000;

/**
 * Ensures that the stream exists for the given options.
 * If the stream does not exist, it creates a new one with dedup window enabled.
 *
 * @param client - The NatsClient instance
 * @param options - Consumer options
 */
async function ensureStream(client: NatsClient, options: ConsumerOptions): Promise<void> {
	const jsm = await client.nc.jetstreamManager();

	try {
		await jsm.streams.info(options.stream);
	} catch {
		// Stream doesn't exist, create it
		await jsm.streams.add({
			name: options.stream,
			subjects: [options.subject],
			retention: RetentionPolicy.Workqueue,
			storage: StorageType.File,
			duplicate_window: DEFAULT_DUPLICATE_WINDOW_NS,
		});
	}
}

/**
 * Ensures a JetStream stream has deduplication configured.
 * Updates existing streams with the specified duplicate_window.
 *
 * Call this at service startup for streams that receive fire-and-forget publishes
 * (publishCore + flush) to enable Nats-Msg-Id deduplication.
 *
 * @param client - NatsClient instance
 * @param streamName - Name of the JetStream stream
 * @param duplicateWindowMs - Deduplication window in milliseconds (default: 600000 = 10 min)
 */
export async function ensureStreamDedup(
	client: NatsClient,
	streamName: string,
	duplicateWindowMs: number = 10 * 60 * 1000,
): Promise<void> {
	const jsm = await client.nc.jetstreamManager();
	const dupWindowNanos = duplicateWindowMs * 1_000_000;

	try {
		const info = await jsm.streams.info(streamName);
		if (info.config.duplicate_window !== dupWindowNanos) {
			await jsm.streams.update(streamName, {
				...info.config,
				duplicate_window: dupWindowNanos,
			});
			getInfraLogger().info({ stream: streamName, windowMs: duplicateWindowMs }, '[INFRA:NATS] Stream dedup window updated');
		}
	} catch {
		getInfraLogger().debug({ stream: streamName }, '[INFRA:NATS] Stream not found — dedup will be set on creation');
	}
}

/**
 * Ensures that a consumer exists for the given stream and durable name.
 * If the consumer does not exist, it creates a new one.
 *
 * @param client - The NatsClient instance
 * @param options - Consumer options
 */
async function ensureConsumer(client: NatsClient, options: ConsumerOptions): Promise<void> {
	const jsm = await client.nc.jetstreamManager();

	// Ensure stream exists first
	await ensureStream(client, options);

	// Try to get existing consumer
	try {
		await jsm.consumers.info(options.stream, options.durable);
		return;
	} catch {
		// Consumer doesn't exist, create it
		await jsm.consumers.add(options.stream, {
			durable_name: options.durable,
			ack_policy: AckPolicy.Explicit,
			deliver_policy: options.deliverPolicy ?? DeliverPolicy.All,
			filter_subject: options.filterSubject ?? options.subject,
			max_ack_pending: options.maxAckPending ?? (options.batchSize ? options.batchSize * 2 : 128),
		});
	}
}

/**
 * Ensures DLQ stream exists for dead letter messages
 *
 * @param client - The NatsClient instance
 */
async function ensureDLQStream(client: NatsClient): Promise<void> {
	const jsm = await client.nc.jetstreamManager();

	try {
		await jsm.streams.info(DLQ_STREAM);
	} catch {
		// DLQ stream doesn't exist, create it
		await jsm.streams.add({
			name: DLQ_STREAM,
			subjects: [DLQ_SUBJECT],
			retention: RetentionPolicy.Limits,
			storage: StorageType.File,
			max_age: DLQ_MAX_AGE_NANOS,
		});
	}
}

/**
 * Publishes a message to the DLQ stream
 *
 * @param client - The NatsClient instance
 * @param dlqMsg - The DLQ message to publish
 */
async function publishToDLQ(client: NatsClient, dlqMsg: DLQMessage): Promise<void> {
	const data = new TextEncoder().encode(JSON.stringify(dlqMsg));
	await client.js.publish(DLQ_SUBJECT, data);
	getInfraLogger().info({ dlqId: dlqMsg.id }, '[INFRA:NATS] Message sent to DLQ');
}

/**
 * Starts a NATS consumer with retry/DLQ support (V2 API).
 *
 * This function creates a pull-based consumer that fetches messages in batches
 * and provides Message wrappers with Ack/Nack/Reject methods for service-level
 * control over message lifecycle.
 *
 * ## Architecture:
 *
 * ```
 * NATS Fetch (batch of N messages)
 *            ↓
 * BatchMessageHandlerV2 receives ALL messages at once
 *            ↓
 * Service processes messages and calls:
 *   - msg.ack() for success
 *   - msg.nack(err) for retry with backoff
 *   - msg.reject(reason) for immediate DLQ
 *            ↓
 * Message handles ACK/NAK/DLQ internally
 * ```
 *
 * ## Usage:
 *
 * ```typescript
 * const consumer = await startConsumer(natsBus.client, {
 *   stream: 'MY-STREAM',
 *   subject: 'my.subject',
 *   durable: 'my-consumer',
 *   retryPolicy: {
 *     maxRetries: 5,
 *     backoff: [1000, 5000, 30000, 120000, 600000],
 *   },
 *   dlqPolicy: {
 *     serviceName: 'my-service',
 *     serviceType: 'processor',
 *     eventType: 'my.event',
 *   },
 *   batchMessageHandlerV2: async (messages) => {
 *     for (const msg of messages) {
 *       try {
 *         const data = JSON.parse(new TextDecoder().decode(msg.data));
 *         msg.orgId = data.orgId; // MANDATORY for multi-tenant filtering
 *         msg.pathKey = data.pathKey;
 *         await processMessage(data);
 *         msg.ack();
 *       } catch (err) {
 *         await msg.nack(err);
 *       }
 *     }
 *   },
 * });
 *
 * // To stop:
 * await consumer.stop();
 * ```
 *
 * @param client - The NatsClient instance
 * @param options - Consumer options with retry/DLQ policies
 * @returns Consumer handle with stop function
 */
export async function startConsumer(
	client: NatsClient,
	options: ConsumerOptions,
): Promise<ConsumerHandle> {
	// Validate required options
	if (!options.stream) {
		throw new Error('[NATS Consumer] Stream name is required');
	}
	if (!options.subject) {
		throw new Error('[NATS Consumer] Subject is required');
	}
	if (!options.durable) {
		throw new Error('[NATS Consumer] Durable name is required');
	}
	if (!options.batchMessageHandlerV2) {
		throw new Error('[NATS Consumer] batchMessageHandlerV2 is required');
	}

	// Ensure consumer exists
	await ensureConsumer(client, options);

	// Ensure DLQ stream exists if DLQ policy is configured
	if (options.dlqPolicy) {
		await ensureDLQStream(client);
	}

	// Get consumer handle
	const consumers = (client.js as any).consumers;
	if (!consumers) {
		throw new Error('[NATS Consumer] Consumer API not available');
	}
	const consumer = await consumers.get(options.stream, options.durable);

	// Read pacing configuration
	const batchSize = options.batchSize ?? 50;
	const fetchTimeoutMs = options.fetchTimeout ?? 5000;

	let isRunning = true;

	/**
	 * Wrap raw JsMsg into Message with retry/DLQ support
	 */
	const wrapMessage = (jsMsg: JsMsg): Message => {
		return new Message(
			jsMsg,
			options.stream,
			options.durable,
			options.retryPolicy,
			options.dlqPolicy,
			(dlqMsg: DLQMessage) => publishToDLQ(client, dlqMsg),
		);
	};

	/**
	 * Fetch a batch of messages from NATS
	 */
	const fetchBatch = async (): Promise<Message[]> => {
		const iterator = await consumer.fetch({
			max_messages: batchSize,
			expires: fetchTimeoutMs,
		});
		const messages: Message[] = [];
		for await (const jsMsg of iterator) {
			if (!isRunning) break;
			messages.push(wrapMessage(jsMsg));
		}
		return messages;
	};

	/**
	 * Handle errors from batch handler — NAK all unhandled messages
	 */
	const handleBatchError = async (err: unknown, messages: Message[]): Promise<void> => {
		getInfraLogger().error({ err }, '[INFRA:NATS] Handler error');
		const errorMsg = err instanceof Error ? err.message : String(err);
		for (const msg of messages) {
			try {
				await msg.nack(errorMsg);
			} catch {
				// Already handled, ignore
			}
		}
	};

	/**
	 * Main processing loop — Always-on double-buffer
	 *
	 * The next fetch is ALWAYS in-flight, even when the stream is empty.
	 * NATS server-side `expires` (default 5s) acts as a long-poll, holding
	 * the request until data arrives — no client-side polling, no hammering.
	 *
	 *   Stream HAS DATA → fetch(N+1) already in-flight → 0ms idle between batches
	 *   Stream EMPTY    → fetch blocks server-side (expires=5s) → zero CPU
	 *   No modes, no state variable, no if/else — always the same code path.
	 */
	const processMessages = async (): Promise<void> => {
		let nextFetch = fetchBatch();

		while (isRunning) {
			try {
				const messages = await nextFetch;

				nextFetch = fetchBatch();

				if (messages.length === 0 || !isRunning) continue;

				try {
					await options.batchMessageHandlerV2!(messages);
				} catch (err) {
					await handleBatchError(err, messages);
				}
			} catch (err) {
				if (isRunning) {
					getInfraLogger().warn({ err }, '[INFRA:NATS] Fetch error, retrying...');
					await sleep(1000);
					nextFetch = fetchBatch();
				}
			}
		}
	};

	// Start processing loop
	processMessages().catch((err) => {
		getInfraLogger().error({ err }, '[INFRA:NATS] Processing loop error');
	});

	getInfraLogger().info({ durable: options.durable }, '[INFRA:NATS] Consumer started');

	// Return consumer handle
	return {
		stop: async (): Promise<void> => {
			isRunning = false;
			try {
				await consumer.close?.();
			} catch {
				// Ignore close errors
			}
			getInfraLogger().info({ durable: options.durable }, '[INFRA:NATS] Consumer stopped');
		},
	};
}

import type { Publisher, Subscriber } from './ports';
import type { NatsClient, NatsSubscribeInternalOptions, ConsumerOptions, ConsumerHandle, FanoutSubscribeOptions, FanoutSubscription, FanoutStreamConfig } from './types';

import { consumerOpts, createInbox, RetentionPolicy, StorageType } from 'nats';
import { publish as _publish, publishCore as _publishCore, flushConnection as _flushConnection } from './publish';
import { subscribe as _subscribe } from './subscribe';
import { startConsumer as _startConsumer } from './consumer';
import { getInfraLogger } from '../logger';

export class NatsBus implements Publisher, Subscriber {
	constructor(private c: NatsClient) {
	}

	/**
	 * Get the underlying NATS client for direct access
	 */
	get client(): NatsClient {
		return this.c;
	}

	/**
	 * Measures the round-trip time to the NATS server using the native PING/PONG protocol.
	 *
	 * @returns RTT in milliseconds
	 * @throws If the connection is closed
	 */
	async ping(): Promise<number> {
		if (!this.c.nc || this.c.nc.isClosed()) {
			throw new Error('nats: connection is closed');
		}
		const rtt = await this.c.nc.rtt();
		return rtt;
	}

	/**
	 * Returns true if the NATS connection is active and not closed.
	 */
	isConnected(): boolean {
		return !!this.c.nc && !this.c.nc.isClosed();
	}

	/**
	 * Publishes a message to a specified subject.
	 *
	 * @param subject - The subject to which the message will be published.
	 * @param payload - The message payload, which can be any type. If not a Uint8Array, it will be encoded as JSON.
	 * @param headers - Optional headers to include with the message, represented as a record of string key-value pairs.
	 * @returns A promise that resolves when the message has been successfully published.
	 */
	async publish(subject: string, payload: unknown, headers?: Record<string, string>): Promise<void> {
		const data = payload instanceof Uint8Array ? payload : new TextEncoder().encode(JSON.stringify(payload));
		await _publish(this.c, { subject, data, headers });
	}

	/**
	 * Publishes a message using core NATS (fire-and-forget, no JetStream ACK).
	 *
	 * Enqueues the message in the TCP buffer. Call flush() after a batch
	 * of publishCore() calls to guarantee delivery to the NATS server.
	 *
	 * For JetStream deduplication, pass msgId (sets Nats-Msg-Id header).
	 * JetStream discards duplicates within the stream's duplicate_window.
	 *
	 * @param subject - The subject to publish to
	 * @param payload - The message payload (auto-encoded as JSON if not Uint8Array)
	 * @param headers - Optional headers
	 * @param msgId - Optional Nats-Msg-Id for JetStream deduplication
	 */
	publishCore(subject: string, payload: unknown, headers?: Record<string, string>, msgId?: string): void {
		const data = payload instanceof Uint8Array ? payload : new TextEncoder().encode(JSON.stringify(payload));
		_publishCore(this.c.nc, { subject, data, headers, msgId });
	}

	/**
	 * Flushes all pending core NATS publishes to the server.
	 * Single TCP roundtrip — call after a batch of publishCore() calls.
	 */
	async flush(): Promise<void> {
		await _flushConnection(this.c.nc);
	}

	/**
	 * Subscribes to a specified subject and processes incoming messages using a handler function.
	 *
	 * For pull-based subscriptions, a stream name is required.
	 * For push-based subscriptions, stream name is optional but recommended.
	 *
	 * @param opts - An object containing subscription options.
	 * @param opts.stream - Stream name (required for pull subscriptions, recommended for push).
	 * @param opts.subject - The subject to subscribe to.
	 * @param opts.durable - Optional durable name for the subscription, allowing message persistence.
	 * @param opts.queueGroup - Optional queue group name for load balancing message processing.
	 * @param opts.pull - Optional flag indicating if the subscription should be a pull-based subscription.
	 * @param opts.handler - A function to handle incoming messages, receiving the message data as a Uint8Array.
	 * @returns A promise that resolves to a function which, when called, will unsubscribe from the subject.
	 */
	async subscribe(opts: NatsSubscribeInternalOptions): Promise<() => Promise<void>> {
		// Validate stream requirement for pull subscriptions

		if (opts.pull && !opts.stream) {
			throw new Error('Pull subscriptions require a stream name');
		}

		return _subscribe(this.c, opts);
	}

	/**
	 * Starts a consumer with V2 API (retry/DLQ support).
	 *
	 * This is the recommended method for production use. It provides:
	 * - Message wrapper with Ack/Nack/Reject methods
	 * - Automatic retry with exponential backoff
	 * - Dead Letter Queue support
	 * - Service-level control over message lifecycle
	 *
	 * ## Usage:
	 *
	 * ```typescript
	 * const consumer = await natsBus.startConsumer({
	 *   stream: 'MY-STREAM',
	 *   subject: 'my.subject',
	 *   durable: 'my-consumer',
	 *   retryPolicy: { maxRetries: 5, backoff: [1000, 5000, 30000] },
	 *   dlqPolicy: { serviceName: 'my-service', serviceType: 'processor', eventType: 'my.event' },
	 *   batchMessageHandlerV2: async (messages) => {
	 *     for (const msg of messages) {
	 *       try {
	 *         const data = JSON.parse(new TextDecoder().decode(msg.data));
	 *         msg.orgId = data.orgId;
	 *         msg.pathKey = data.pathKey;
	 *         await processMessage(data);
	 *         msg.ack();
	 *       } catch (err) {
	 *         await msg.nack(err);
	 *       }
	 *     }
	 *   },
	 * });
	 * ```
	 *
	 * @param opts - Consumer options with retry/DLQ policies
	 * @returns Consumer handle with stop function
	 */
	async startConsumer(opts: ConsumerOptions): Promise<ConsumerHandle> {
		return _startConsumer(this.c, opts);
	}

	/**
	 * Ensures a FANOUT stream exists with appropriate settings for ephemeral consumers.
	 * Call this BEFORE subscribeFanout to ensure the stream exists.
	 *
	 * FANOUT streams use:
	 * - MemoryStorage for low latency
	 * - LimitsPolicy retention (time/count based expiration)
	 * - Short max_age (default 5 minutes)
	 *
	 * @param config - FANOUT stream configuration
	 */
	async ensureFanoutStream(config: FanoutStreamConfig): Promise<void> {
		if (!config.name) {
			throw new Error('Stream name is required');
		}
		if (!config.subjects || config.subjects.length === 0) {
			throw new Error('At least one subject is required');
		}

		const jsm = await this.c.nc.jetstreamManager();

		// Check if stream already exists
		try {
			await jsm.streams.info(config.name);
			getInfraLogger().debug({ stream: config.name }, '[INFRA:NATS] FANOUT stream already exists');
			return;
		} catch {
			// Stream doesn't exist, create it
		}

		// Set defaults
		const maxAge = config.maxAge ?? 5 * 60 * 1000; // 5 minutes in ms
		const maxMsgs = config.maxMsgs ?? 10000;
		const maxBytes = config.maxBytes ?? 10 * 1024 * 1024; // 10MB

		await jsm.streams.add({
			name: config.name,
			description: config.description,
			subjects: config.subjects,
			retention: RetentionPolicy.Limits,
			max_age: maxAge * 1000000, // Convert ms to nanoseconds
			max_msgs: maxMsgs,
			max_bytes: maxBytes,
			storage: StorageType.Memory, // Memory for low latency
			num_replicas: 1,
			discard: 'old' as any,
		});

		getInfraLogger().info({ stream: config.name }, '[INFRA:NATS] FANOUT stream created');
	}

	/**
	 * Subscribes to a FANOUT subject using an ephemeral consumer.
	 * All instances receive all messages (broadcast pattern).
	 *
	 * IMPORTANT: Call ensureFanoutStream() BEFORE this method to ensure
	 * the stream exists.
	 *
	 * FANOUT Pattern:
	 * - Each service instance receives a copy of the message (no queue group)
	 * - Used for cache invalidation across all replicas
	 * - Ephemeral consumer (not durable) - created fresh on each startup
	 * - DeliverNew policy - only receives new messages (no replay)
	 *
	 * @param opts - FANOUT subscription options
	 * @returns FanoutSubscription with stop() method
	 */
	async subscribeFanout(opts: FanoutSubscribeOptions): Promise<FanoutSubscription> {
		if (!opts.handler) {
			throw new Error('Handler is required');
		}
		if (!opts.serviceName) {
			throw new Error('Service name is required');
		}
		if (!opts.subject) {
			throw new Error('Subject is required');
		}
		if (!opts.stream) {
			throw new Error('Stream is required');
		}

		// Create unique consumer name with random suffix for ephemeral consumer
		const randomSuffix = Math.random().toString(36).substring(2, 10);
		const consumerName = `${opts.serviceName}-fanout-${randomSuffix}`;

		getInfraLogger().debug({ consumer: consumerName, subject: opts.subject, stream: opts.stream }, '[INFRA:NATS] Creating FANOUT subscription');

		const co = consumerOpts();
		co.bindStream(opts.stream);
		co.deliverNew(); // Only receive new messages (not replay)
		co.manualAck();
		co.ackWait(30 * 1000); // 30 seconds ack wait
		co.inactiveEphemeralThreshold(5 * 60 * 1000); // 5 minutes inactive threshold
		co.consumerName(consumerName);

		// Push consumers require a deliver subject - use an inbox for ephemeral delivery
		co.deliverTo(createInbox());

		const subscription = await this.c.js.subscribe(opts.subject, co);

		// Process messages asynchronously
		const processMessages = async () => {
			for await (const message of subscription) {
				try {
					await opts.handler(message.data);
				} catch (err) {
					getInfraLogger().warn({ subject: opts.subject, err }, '[INFRA:NATS] FANOUT handler error');
				}
				// Always ack FANOUT messages (fire-and-forget)
				message.ack();
			}
		};

		processMessages().catch(() => {
			// Swallow loop rejection on shutdown
		});

		getInfraLogger().debug({ consumer: consumerName }, '[INFRA:NATS] FANOUT subscription active');

		return {
			stop: async () => {
				try {
					await (subscription as any).drain?.();
				} catch {
					(subscription as any).unsubscribe?.();
				}
			},
		};
	}
}
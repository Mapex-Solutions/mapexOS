import type {
	JetStreamClient,
	NatsConnection,
	ConnectionOptions,
	JsMsg,
} from 'nats';

import { AckPolicy, DeliverPolicy } from 'nats';

import { DEFAULT_RETRY_BACKOFF, DEFAULT_MAX_RETRIES } from './constants';
import { getInfraLogger } from '../logger';

export interface NatsClient {
	js: JetStreamClient;
	nc: NatsConnection;
}

export type NatsConnectionOptions = ConnectionOptions;

export interface NatsPublishOptions {
	subject: string;
	data: Uint8Array;
	headers?: Record<string, string>;
	/** Nats-Msg-Id header for JetStream deduplication.
	 *  When set, JetStream discards duplicates within the stream's duplicate_window. */
	msgId?: string;
}

export interface NatsSubscribeInternalOptions {
	subject: string;
	stream?: string;
	durable?: string;
	queueGroup?: string;
	pull?: boolean;
	handler: (data: Uint8Array, subject?: string) => Promise<void> | void;
	maxMessages?: number;     // batch size
  expires?: number;         // server wait-to-fill (ms)
	deliverPolicy?: DeliverPolicy;
	ackPolicy?: AckPolicy;
	filterSubject?: string;
}

/**
 * V2 Types - Retry/DLQ Support
 */

/**
 * RetryPolicy defines retry behavior with exponential backoff
 */
export interface RetryPolicy {
	/** Maximum number of retries before sending to DLQ (default: 5) */
	maxRetries: number;
	/** Backoff delays in milliseconds for each retry attempt */
	backoff: number[];
}

/**
 * DLQPolicy defines Dead Letter Queue configuration
 */
export interface DLQPolicy {
	/** Service name for DLQ message metadata */
	serviceName: string;
	/** Service type for DLQ filtering (e.g., "processor", "router") */
	serviceType: string;
	/** Event type for DLQ filtering (e.g., "js.execute", "route.execute") */
	eventType: string;
}

/**
 * DLQMessage represents a message sent to the Dead Letter Queue
 */
export interface DLQMessage {
	id: string;
	orgId: string;
	pathKey: string;
	serviceName: string;
	serviceType: string;
	eventType: string;
	originalSubject: string;
	originalStream: string;
	originalData: string;
	originalHeaders: Record<string, string>;
	lastError: string;
	errorCount: number;
	firstDelivery: string;
	lastDelivery: string;
	totalDeliveries: number;
	consumerName: string;
	sentToDLQAt: string;
}

/**
 * Message wraps a NATS JsMsg with retry/DLQ-aware methods.
 * Service layer uses these methods to control message lifecycle.
 */
export class Message {
	/** Raw message data */
	readonly data: Uint8Array;
	/** Message headers */
	readonly headers: Record<string, string[]>;
	/** Message subject */
	readonly subject: string;
	/** Number of delivery attempts */
	readonly deliveryCount: number;
	/** Stream sequence number — stable across redeliveries, used for deterministic dedup IDs */
	readonly streamSequence: number;
	/** Message timestamp */
	readonly timestamp: Date;
	/** Tenant context - must be set by service after parsing payload */
	orgId: string = '';
	/** Hierarchical path key - must be set by service after parsing payload */
	pathKey: string = '';

	private readonly msg: JsMsg;
	private readonly retryPolicy?: RetryPolicy;
	private readonly dlqPolicy?: DLQPolicy;
	private readonly publishToDLQ: (dlqMsg: DLQMessage) => Promise<void>;
	private readonly stream: string;
	private readonly consumerName: string;
	private handled: boolean = false;
	private lastError: string = '';
	private firstDelivery: Date;

	constructor(
		msg: JsMsg,
		stream: string,
		consumerName: string,
		retryPolicy: RetryPolicy | undefined,
		dlqPolicy: DLQPolicy | undefined,
		publishToDLQ: (dlqMsg: DLQMessage) => Promise<void>,
	) {
		this.msg = msg;
		this.stream = stream;
		this.consumerName = consumerName;
		this.retryPolicy = retryPolicy;
		this.dlqPolicy = dlqPolicy;
		this.publishToDLQ = publishToDLQ;

		this.data = msg.data;
		this.subject = msg.subject;
		this.deliveryCount = msg.info?.redeliveryCount ?? 0;
		this.streamSequence = msg.info?.streamSequence ?? 0;
		this.timestamp = new Date();

		// Parse headers
		this.headers = {};
		if (msg.headers) {
			for (const [key, values] of msg.headers) {
				this.headers[key] = values;
			}
		}

		// Estimate first delivery time based on redelivery count
		this.firstDelivery = new Date();
		if (this.deliveryCount > 0 && this.retryPolicy?.backoff) {
			const totalBackoff = this.retryPolicy.backoff
				.slice(0, this.deliveryCount)
				.reduce((sum, delay) => sum + delay, 0);
			this.firstDelivery = new Date(Date.now() - totalBackoff);
		}
	}

	/**
	 * Acknowledge the message - processed successfully
	 */
	ack(): void {
		if (this.handled) return;
		this.handled = true;
		this.msg.ack();
	}

	/**
	 * Negative acknowledge - retry with backoff or send to DLQ
	 * @param error - Error that caused the failure
	 */
	async nack(error: Error | string): Promise<void> {
		if (this.handled) return;
		this.handled = true;

		this.lastError = error instanceof Error ? error.message : error;

		// Check if we should send to DLQ
		const maxRetries = this.retryPolicy?.maxRetries ?? DEFAULT_MAX_RETRIES;
		if (this.deliveryCount >= maxRetries) {
			await this.sendToDLQ(this.lastError);
			this.msg.ack(); // ACK to remove from queue after DLQ
			return;
		}

		// Calculate backoff delay
		const backoffDelays = this.retryPolicy?.backoff ?? DEFAULT_RETRY_BACKOFF;
		const delayIndex = Math.min(this.deliveryCount, backoffDelays.length - 1);
		const delayMs = backoffDelays[delayIndex];

		// NAK with delay for server-side retry
		this.msg.nak(delayMs);
	}

	/**
	 * Reject the message - invalid data, send to DLQ immediately
	 * @param reason - Reason for rejection
	 */
	async reject(reason: string): Promise<void> {
		if (this.handled) return;
		this.handled = true;

		this.lastError = reason;
		await this.sendToDLQ(reason);
		this.msg.ack(); // ACK to remove from queue after DLQ
	}

	/**
	 * Terminate the message - unrecoverable error, no DLQ
	 */
	term(): void {
		if (this.handled) return;
		this.handled = true;

		try {
			this.msg.term();
		} catch {
			// Fallback to ack if term not available
			this.msg.ack();
		}
	}

	/**
	 * Send message to Dead Letter Queue
	 */
	private async sendToDLQ(errorMsg: string): Promise<void> {
		if (!this.dlqPolicy) {
			// No DLQ policy - just log and return
			getInfraLogger().warn('[INFRA:NATS] No DLQ policy configured, message will be lost');
			return;
		}

		const headersObj: Record<string, string> = {};
		for (const [key, values] of Object.entries(this.headers)) {
			headersObj[key] = values.join(', ');
		}

		const dlqMessage: DLQMessage = {
			id: crypto.randomUUID(),
			orgId: this.orgId,
			pathKey: this.pathKey,
			serviceName: this.dlqPolicy.serviceName,
			serviceType: this.dlqPolicy.serviceType,
			eventType: this.dlqPolicy.eventType,
			originalSubject: this.subject,
			originalStream: this.stream,
			originalData: new TextDecoder().decode(this.data),
			originalHeaders: headersObj,
			lastError: errorMsg,
			errorCount: this.deliveryCount + 1,
			firstDelivery: this.firstDelivery.toISOString(),
			lastDelivery: new Date().toISOString(),
			totalDeliveries: this.deliveryCount + 1,
			consumerName: this.consumerName,
			sentToDLQAt: new Date().toISOString(),
		};

		try {
			await this.publishToDLQ(dlqMessage);
		} catch (err) {
			getInfraLogger().error({ err }, '[INFRA:NATS] Failed to publish to DLQ');
		}
	}
}

/**
 * Consumer handle returned by startConsumer
 */
export interface ConsumerHandle {
	/** Stop the consumer gracefully */
	stop: () => Promise<void>;
}

/**
 * FANOUT Types - Ephemeral broadcast subscriptions
 */

/**
 * FanoutHandler is a simple handler for FANOUT messages
 * Returns void - errors are logged but not retried
 */
export type FanoutHandler = (data: Uint8Array) => void | Promise<void>;

/**
 * FanoutSubscription represents an active FANOUT subscription
 */
export interface FanoutSubscription {
	/** Stop the subscription */
	stop: () => Promise<void>;
}

/**
 * FanoutSubscribeOptions for subscribeFanout
 */
export interface FanoutSubscribeOptions {
	/** Stream name (required) */
	stream: string;
	/** Service name for ephemeral consumer naming (required) */
	serviceName: string;
	/** Subject to subscribe to (required) */
	subject: string;
	/** Handler function for incoming messages */
	handler: FanoutHandler;
}

/**
 * FanoutStreamConfig for ensureFanoutStream
 * Configures a FANOUT stream with appropriate settings for ephemeral consumers.
 */
export interface FanoutStreamConfig {
	/** Stream name (required) */
	name: string;
	/** Subject patterns for the stream (required). Example: ['fanout.>'] */
	subjects: string[];
	/** Maximum age for messages in milliseconds (default: 5 minutes) */
	maxAge?: number;
	/** Maximum number of messages (default: 10000) */
	maxMsgs?: number;
	/** Maximum total size in bytes (default: 10MB) */
	maxBytes?: number;
	/** Optional stream description */
	description?: string;
}

/**
 * ConsumerOptions for V2 API with retry/DLQ support
 */
export interface ConsumerOptions {
	/** Stream name (required) */
	stream: string;
	/** Subject to subscribe to (required) */
	subject: string;
	/** Durable consumer name (required) */
	durable: string;
	/** Queue group for load balancing */
	queueGroup?: string;
	/** Batch size for fetching messages (default: 50) */
	batchSize?: number;
	/** Fetch timeout in milliseconds — server-side long-poll (default: 5000) */
	fetchTimeout?: number;
	/** Maximum unacknowledged messages (default: auto-calculated from batchSize * 2) */
	maxAckPending?: number;
	/** Delivery policy (default: All) */
	deliverPolicy?: DeliverPolicy;
	/** Filter subject for wildcards */
	filterSubject?: string;
	/** Retry policy with exponential backoff */
	retryPolicy?: RetryPolicy;
	/** Dead Letter Queue policy */
	dlqPolicy?: DLQPolicy;

	/**
	 * V2 batch handler - receives all messages with Ack/Nack/Reject control.
	 * Service is responsible for calling msg.ack(), msg.nack(err), or msg.reject(reason).
	 */
	batchMessageHandlerV2?: (messages: Message[]) => Promise<void> | void;

}
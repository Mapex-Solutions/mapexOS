import type { NatsClient, NatsPublishOptions } from './types';
import type { NatsConnection } from 'nats';
import { headers as natsHeaders } from 'nats';

/**
 * Builds NATS headers from NatsPublishOptions.
 * Sets Nats-Msg-Id if msgId is provided (JetStream deduplication).
 */
function buildHeaders(opts: NatsPublishOptions) {
	const hdrs = natsHeaders();
	if (opts.headers) {
		for (const [k, v] of Object.entries(opts.headers)) hdrs.set(k, v);
	}
	if (opts.msgId) {
		hdrs.set('Nats-Msg-Id', opts.msgId);
	}
	return hdrs;
}

/**
 * Publishes a message using JetStream (acknowledged publish).
 * Waits for the server ACK — guarantees the message was persisted.
 *
 * Use this for critical messages where you need confirmation of persistence.
 * For high-throughput fire-and-forget, use publishCore() + flushConnection().
 *
 * @param c - The NATS client instance (needs JetStream)
 * @param opts - Publish options including subject, data, optional headers, and optional msgId for dedup
 */
export async function publish(c: NatsClient, opts: NatsPublishOptions): Promise<void> {
	const hdrs = buildHeaders(opts);
	await c.js.publish(opts.subject, opts.data, { headers: hdrs });
}

/**
 * Publishes a message using core NATS (fire-and-forget, no JetStream ACK).
 *
 * The message is enqueued in the client TCP buffer and sent on the next flush.
 * JetStream still captures it if a stream matches the subject.
 *
 * Use with flushConnection() after a batch of publishes to guarantee
 * all messages reached the NATS server.
 *
 * Safety: Set opts.msgId (Nats-Msg-Id header) for JetStream deduplication.
 * If the publisher crashes and messages are re-published on retry,
 * JetStream discards duplicates within the stream's duplicate_window.
 *
 * @param nc - The raw NATS connection (not JetStream)
 * @param opts - Publish options including subject, data, optional headers, and optional msgId for dedup
 */
export function publishCore(nc: NatsConnection, opts: NatsPublishOptions): void {
	const hdrs = buildHeaders(opts);
	nc.publish(opts.subject, opts.data, { headers: hdrs });
}

/**
 * Flushes all pending core NATS publishes to the server.
 *
 * After calling publishCore() N times, call flushConnection() once
 * to guarantee all messages have been delivered to the NATS server.
 *
 * This is a single TCP roundtrip regardless of how many messages are pending.
 *
 * @param nc - The raw NATS connection
 */
export async function flushConnection(nc: NatsConnection): Promise<void> {
	await nc.flush();
}

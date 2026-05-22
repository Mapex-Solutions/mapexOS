/**
 * OOMError represents an Out-of-Memory error from a V8 Isolate.
 *
 * This error is TRANSIENT — the isolate has been recycled and a fresh one is available.
 * Throwing this error signals to the message handler that the event should be NACK'd
 * for retry, NOT acknowledged as processed.
 *
 * Flow: V8 OOM → isolate disposed → OOMError thrown → msg.nack() → retry with backoff
 */
export class OOMError extends Error {
	constructor(message: string) {
		super(message);
		this.name = 'OOMError';
	}
}

import type { Counter, Histogram } from 'prom-client';

/**
 * Prometheus metrics accepted by ScriptService (all optional).
 */
export interface ScriptServiceMetrics {
	/** Total events by consumer and status */
	eventsProcessed: Counter;
	/** End-to-end event processing time */
	eventDuration: Histogram;
	/** Incoming payload size in bytes */
	payloadSize: Histogram;
	/** Asset cache hits per tier */
	assetCache: Counter;
	/** Template cache hits per tier */
	templateCache: Counter;
	/**
	 * Implicit heartbeats published by js-executor (heartbeatMode='implicit'
	 * AND healthMonitor.enabled=true).
	 */
	heartbeatsPublished?: Counter;
	/**
	 * Heartbeat publishes skipped by the gate, labeled by reason:
	 *   - 'disabled'      → healthMonitor.enabled=false
	 *   - 'explicit_mode' → heartbeatMode='explicit'
	 */
	heartbeatsSkipped?: Counter<'reason'>;
}

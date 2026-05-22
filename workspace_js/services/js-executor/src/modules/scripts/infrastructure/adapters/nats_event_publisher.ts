import type { NatsBus } from '@mapexos/infrastructure';
import type {
	EventPublisherPort,
	PublishResultParams,
	PublishRawEventParams,
	PublishExecutionLogParams,
	PublishHeartbeatParams,
} from '../../application/ports';

import { subject } from '@shared/configuration/naming';

// Subjects use the platform's env-prefixed naming convention so the
// same binary serves multiple environments on a shared NATS cluster
// (e.g. dev.mapexos.events.raw vs prod.mapexos.events.raw). All four
// JetStream streams (events-raw, route-execute, events-logs, assets-
// heartbeat) are provisioned by nats-init with these env-prefixed
// filters — publishing without the prefix produces silent black-hole
// messages NATS never delivers anywhere.
const ROUTE_EXECUTE_SUBJECT = subject('route', 'execute');
const EVENTS_RAW_SUBJECT = subject('events', 'raw');
const EVENTS_LOGS_JSEXECUTOR_SUBJECT = subject('events', 'logs.jsexecutor');
const ASSET_HEARTBEAT_SUBJECT_PREFIX = subject('asset', 'heartbeat');

/**
 * Top-level eventSource discriminator stamped on every mapexos.route.execute
 * payload produced by this adapter. The router uses it to decide which
 * route-group collection on the asset cache to use (assetEvent →
 * asset.RouteGroupIds). Symmetric to healthmonitor's `'healthStatus'`.
 */
const EVENT_SOURCE_ASSET = 'assetEvent';

/**
 * NatsEventPublisherAdapter implements EventPublisherPort using NatsBus.
 *
 * All publishes use publishCore (core NATS, fire-and-forget) for maximum throughput.
 * Safety guarantee: consumer ACKs the original message ONLY after flush(),
 * so if the service crashes before ACK, NATS redelivers the original event.
 * Dedup IDs (Nats-Msg-Id) prevent duplicate processing on redelivery.
 *
 * Subjects:
 *   - mapexos.route.execute: Router pipeline (transformed events) — tags
 *     eventSource='assetEvent' so the Router can discriminate from health-status
 *     events (which tag 'healthStatus').
 *   - mapexos.events.raw: Debug raw event logging (ClickHouse events_raw)
 *   - mapexos.events.logs.jsexecutor: JS execution logs (ClickHouse events_jsexecutor)
 *   - mapexos.asset.heartbeat.{orgId}: Health monitoring heartbeat
 */
export class NatsEventPublisherAdapter implements EventPublisherPort {
	constructor(private readonly natsBus: NatsBus) {}

	/**
	 * Publishes transformed event to route.execute for the Router pipeline.
	 * Sets `eventSource='assetEvent'` at the top level — explicit producer-side
	 * contract that the Router reads to pick the default route-group source
	 * (asset.RouteGroupIds) rather than the HealthMonitor-scoped ones.
	 * Dedup ID: {eventTrackerId}-route
	 *
	 * @param params - Route result parameters
	 */
	publishResult(params: PublishResultParams): void {
		console.log('[TRACE:Publisher] publishResult subject=' + ROUTE_EXECUTE_SUBJECT + ' trackId=' + params.eventTrackerId);
		this.natsBus.publishCore(
			ROUTE_EXECUTE_SUBJECT,
			{
				eventSource: EVENT_SOURCE_ASSET,
				assetUUID: params.assetUUID,
				assetId: params.assetId,
				orgId: params.orgId,
				pathKey: params.pathKey,
				eventTrackerId: params.eventTrackerId,
				dataSource: params.dataSource,
				event: params.event,
			},
			undefined,
			`${params.eventTrackerId}-route`,
		);
	}

	/**
	 * Publishes raw event to events.raw for debug logging.
	 * Dedup ID: {eventTrackerId}-raw
	 *
	 * @param params - Raw event parameters
	 */
	publishRawEvent(params: PublishRawEventParams): void {
		console.log('[TRACE:Publisher] publishRawEvent subject=' + EVENTS_RAW_SUBJECT + ' trackId=' + params.eventTrackerId);
		this.natsBus.publishCore(
			EVENTS_RAW_SUBJECT,
			{
				eventTrackerId: params.eventTrackerId || 'n/a',
				threadId: params.assetUUID || 'n/a',
				orgId: params.orgId || '',
				pathKey: params.pathKey || '',
				event: params.event,
				source: mapSourceType(params.sourceType),
				created: params.timestamp,
				name: params.name || '',
				description: params.description || '',
				success: true,
				error: '',
			},
			undefined,
			`${params.eventTrackerId}-raw`,
		);
	}

	/**
	 * Publishes JS execution log to events.logs.jsexecutor.
	 * Dedup ID: {eventTrackerId}-jslog
	 *
	 * @param params - Execution log parameters
	 */
	publishExecutionLog(params: PublishExecutionLogParams): void {
		this.natsBus.publishCore(
			EVENTS_LOGS_JSEXECUTOR_SUBJECT,
			{
				eventTrackerId: params.eventTrackerId || '',
				created: params.timestamp,
				threadId: params.assetUUID || 'n/a',
				orgId: params.orgId || '',
				pathKey: params.pathKey || '',
				name: params.name || '',
				description: params.description || '',
				execution: {
					...params.execution,
					totalExecutionTime: Math.round(params.execution.totalExecutionTime),
				},
				event: params.event,
			},
			undefined,
			`${params.eventTrackerId}-jslog`,
		);
	}

	/**
	 * Publishes a heartbeat to asset.heartbeat.{orgId} for health monitoring.
	 * This is the SOLE producer of IMPLICIT heartbeats — gated upstream by
	 * the per-event flow when asset.healthMonitor.heartbeatMode is 'implicit'.
	 *
	 * Other liveness paths do NOT pass through this method:
	 *   - Explicit-mode MQTT-protocol assets: liveness is signaled by the
	 *     NATS broker itself ($SYS.ACCOUNT.*.CONNECT/DISCONNECT advisories
	 *     consumed by the assets/healthmonitor module). No code path here.
	 *   - Explicit-mode HTTP-protocol assets: the device POSTs to
	 *     /api/v1/heartbeat?ds=... and the http_gateway publishes directly
	 *     to mapexos.asset.heartbeat.{orgId}.
	 *
	 * The downstream consumer (assets/healthmonitor heartbeat handler) is
	 * origin-agnostic: this implicit publish and the http_gateway explicit
	 * publish share the same subject and payload shape.
	 *
	 * Fire-and-forget — included in the existing flush() roundtrip.
	 *
	 * @param params - Heartbeat parameters
	 */
	publishHeartbeat(params: PublishHeartbeatParams): void {
		this.natsBus.publishCore(
			`${ASSET_HEARTBEAT_SUBJECT_PREFIX}.${params.orgId}`,
			{
				orgId: params.orgId,
				assetUUID: params.assetUUID,
				pathKey: params.pathKey,
				ts: Math.floor(Date.now() / 1000),
			},
		);
	}

	/**
	 * Flushes all buffered core NATS publishes in a single TCP roundtrip.
	 * Must be called after all publish calls and before ACK/Nack.
	 */
	async flush(): Promise<void> {
		await this.natsBus.flush();
	}
}

/**
 * Maps sourceType to gateway name for events logging.
 *
 * @param sourceType - Protocol source type (http/mqtt/lorawan)
 * @returns Gateway name for ClickHouse events_raw table
 */
function mapSourceType(sourceType: string): string {
	switch (sourceType) {
		case 'http': return 'http_gateway';
		case 'mqtt': return 'mqtt_gateway';
		case 'lorawan': return 'lorawan_gateway';
		default: return 'unknown_gateway';
	}
}

/**
 * EventPublisherPort defines the contract for publishing IoT event processing results.
 *
 * Implemented by NatsEventPublisherAdapter (infrastructure layer).
 * Used by consumers (interface layer) after service returns domain results.
 *
 * This port replaces the direct NATS publishing that was previously in ScriptService.
 * Each method produces the exact same payload format as before — zero functional changes.
 */
export interface EventPublisherPort {
	/**
	 * Publishes transformed event to route.execute for the Router pipeline.
	 * Only called on successful script execution.
	 *
	 * @param params - Route result parameters matching the Router's expected payload
	 */
	publishResult(params: PublishResultParams): void;

	/**
	 * Publishes raw event to events.raw for debug logging.
	 * Called BEFORE script processing when debugEnabled=true.
	 *
	 * @param params - Raw event parameters for ClickHouse events_raw table
	 */
	publishRawEvent(params: PublishRawEventParams): void;

	/**
	 * Publishes JS execution log to events.logs.jsexecutor.
	 * Called ALWAYS on failure, and on success only when debugEnabled=true.
	 *
	 * @param params - Execution log parameters for ClickHouse events_jsexecutor table
	 */
	publishExecutionLog(params: PublishExecutionLogParams): void;

	/**
	 * Publishes a fire-and-forget heartbeat to ASSET-HEARTBEAT stream.
	 * Called for each successful script execution.
	 * Included in the existing flush() TCP roundtrip.
	 *
	 * @param params - Heartbeat parameters (orgId, assetUUID, pathKey)
	 */
	publishHeartbeat(params: PublishHeartbeatParams): void;

	/**
	 * Flushes all buffered publishes in a single TCP roundtrip.
	 * Must be called after all publish calls and before ACK/Nack.
	 */
	flush(): Promise<void>;
}

/** Parameters for publishing a heartbeat to ASSET-HEARTBEAT stream. */
export interface PublishHeartbeatParams {
	orgId: string;
	assetUUID: string;
	pathKey: string;
}

/** Parameters for publishing a successful route result to route.execute. */
export interface PublishResultParams {
	assetUUID: string;
	assetId: string;
	orgId: string;
	pathKey: string;
	eventTrackerId: string;
	dataSource: DataSourcePayload;
	event: unknown;
}

/** Parameters for publishing a raw event to events.raw (debug). */
export interface PublishRawEventParams {
	eventTrackerId: string;
	assetUUID: string;
	orgId: string;
	pathKey: string;
	name: string;
	description: string;
	event: unknown;
	sourceType: string;
	timestamp: string;
}

/** Parameters for publishing a JS execution log to events.logs.jsexecutor. */
export interface PublishExecutionLogParams {
	eventTrackerId: string;
	assetUUID: string;
	orgId: string;
	pathKey: string;
	name: string;
	description: string;
	execution: ExecutionLogPayload;
	event: unknown;
	timestamp: string;
}

/** DataSource metadata included in route.execute publish payloads. */
export interface DataSourcePayload {
	id: string;
	orgId: string;
	pathKey: string;
	name: string;
	description: string;
}

/** Execution result details for the JS execution log. */
export interface ExecutionLogPayload {
	success: boolean;
	failedAt: string;
	totalExecutionTime: number;
	error: string;
}

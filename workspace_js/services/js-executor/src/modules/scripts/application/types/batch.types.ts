import type { Logger } from '@mapexos/microservices';
import type { ScriptEngineServicePort } from '@modules/engine/application/ports';
import type { EventPublisherPort } from '@modules/scripts/application/ports/event_publisher_port';
import type { AssetCachePort } from '@modules/scripts/application/ports/asset_cache_port';
import type { TemplateCachePort } from '@modules/scripts/application/ports/template_cache_port';
import type { DataSource } from '@modules/scripts/application/types';
import type { ScriptServiceMetrics } from '@modules/scripts/application/types/script_service.types';

/** Common input contract — normalized from both MQTT and HTTP sources */
export interface BatchMessageInput {
	/** Original message index in the batch (for result correlation) */
	index: number;
	/** Organization ID */
	orgId: string;
	/** Resolved asset UUID */
	assetUUID: string;
	/** Raw event payload */
	event: any;
	/** Source type for logging */
	sourceType: 'mqtt' | 'http' | 'lorawan';
	/** Event tracker ID for dedup/tracing */
	eventTrackerId: string;
	/** DataSource metadata (HTTP has full dataSource, MQTT builds minimal) */
	dataSource?: DataSource;
}

/** Per-message result returned by processBatch */
export interface BatchMessageResult {
	/** Original message index (correlates with input) */
	index: number;
	/** Whether processing succeeded */
	success: boolean;
	/** Error message if failed */
	error?: string;
	/** True if V8 OOM — consumer should nack for retry */
	isOOM?: boolean;
	/** True if permanent error (bad input) — consumer should reject */
	isPermanent?: boolean;
}

/** Internal dependencies passed to handler functions (since TS has no partial classes) */
export interface ScriptServiceInternalDeps {
	logger: Logger;
	scriptEngine: ScriptEngineServicePort;
	assetCachePort: AssetCachePort;
	templateCachePort: TemplateCachePort;
	eventPublisher: EventPublisherPort;
	metrics?: ScriptServiceMetrics;
}

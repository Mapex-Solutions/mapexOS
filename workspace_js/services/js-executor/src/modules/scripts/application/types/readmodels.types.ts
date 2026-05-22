/**
 * Read Model Types for TieredCache
 *
 * These types match the Go structs in workspace_go/packages/contracts/
 * and represent the data stored in MinIO (L2 cache).
 */

/**
 * AssetReadModel represents the asset data stored in MinIO.
 * Key format: assets/{assetUUID}.json
 *
 * This is a CQRS Read Model - a denormalized projection optimized for reads.
 * Matches: workspace_go/packages/contracts/services/assets/assets/dto.go
 */
export interface AssetReadModel {
	/** MongoDB ObjectId */
	id: string;

	/** Device identifier (devEUI, deviceId, etc) */
	uuid: string;

	/** Organization ID for tenant isolation */
	orgId: string;

	/** Hierarchical path for range queries */
	pathKey: string;

	/** Device enabled status */
	enabled: boolean;

	/** Debug logging enabled - when true, sends execution logs to events.raw */
	debugEnabled: boolean;

	/** Asset name */
	name: string;

	/** Asset description */
	description?: string;

	/** Template ID for fetching scripts */
	assetTemplateId?: string;

	/** Template's organization ID (for template cache lookup) */
	assetTemplateOrgId?: string;

	/** Route group IDs for event routing */
	routeGroupIds?: string[];

	/**
	 * Health monitoring configuration.
	 *
	 * heartbeatMode controls who emits heartbeats:
	 *   - 'implicit' (default, missing = same): js-executor emits a heartbeat
	 *     for every data event the device sends.
	 *   - 'explicit': js-executor SKIPS implicit publishes; liveness is
	 *     captured via a path chosen by the asset's protocol:
	 *       • MQTT-protocol assets: automatic via NATS broker presence
	 *         ($SYS.ACCOUNT.*.CONNECT/DISCONNECT advisories).
	 *       • HTTP-protocol assets: device POSTs to /api/v1/heartbeat?ds=…
	 *         with body { assetUUID }.
	 */
	healthMonitor?: {
		enabled: boolean;
		thresholdMinutes: number;
		requiredMisses: number;
		heartbeatMode?: 'implicit' | 'explicit';
		offlineRouteGroupIds?: string[];
		onlineRouteGroupIds?: string[];
	};

	/** Protocol configuration */
	protocol?: {
		type: 'http' | 'mqtt' | 'lorawan';
		mqtt?: {
			clientId: string;
			username: string;
			/**
			 * Bcrypt hash of the device MQTT password. Present in the
			 * read-model L2 payload because the mapex-mqtt-broker plugin
			 * uses it for local CONNECT bcrypt-compare. Optional because
			 * it is empty for cert-only assets and for non-MQTT
			 * protocols. JS-Executor does not consume this field — it
			 * is here to keep the wire shape stable across consumers.
			 */
			passwordHash?: string;
		};
	};

	/**
	 * Active MQTT device certificate metadata. Present when the asset
	 * is authenticating MQTT CONNECTs via mTLS instead of (or in
	 * addition to) password mode. Consumed by the mapex-mqtt-broker
	 * plugin for serial-equality cert validation; JS-Executor ignores it.
	 */
	currentCert?: {
		serial: string;
		fingerprint: string;
		subjectCN: string;
		issuedAt: string;
		expiresAt: string;
	};

	/** Geolocation */
	latitude?: number;
	longitude?: number;

	/** Timestamps */
	created: string;
	updated: string;
}

/**
 * TemplateReadModel represents the template data stored in MinIO.
 * Key format: templates/{templateId}.json
 *
 * Contains the scripts for JS-Executor pipeline.
 * Matches: workspace_go/packages/contracts/services/assets/assets_templates/dto.go
 */
export interface TemplateReadModel {
	/** MongoDB ObjectId */
	id: string;

	/** Template name */
	name: string;

	/** Template enabled status */
	enabled: boolean;

	/** Template description */
	description?: string;

	/** Organization ID */
	orgId?: string;

	/** Path to extract asset UUID from payload */
	assetIdPath?: string;

	/** Scripts (the main data needed by JS-Executor) */

	/** Decode script (scriptProcessor) - processes raw input */
	scriptProcessor?: string;

	/** Validation script - validates decoded data */
	scriptValidator?: string;

	/** Transform script (scriptConversion) - converts to standardized format */
	scriptConversion: string;

	/** Test script for debugging */
	scriptTest?: string;

	/** Classification */

	categoryId?: string;
	categoryName?: string;
	manufacturerId?: string;
	manufacturerName?: string;
	modelId?: string;
	modelName?: string;
	version?: string;

	/** Template flags */

	isSystem?: boolean;
	isTemplate?: boolean;

	/** Available fields for Rule autocomplete */
	availableFields?: string[];

	/** Dynamic fields for typed event storage */
	dynamicFields?: Array<{
		field: string;
		value?: string;
		type: 'string' | 'number' | 'bool' | 'date' | 'geo';
		latitudePath?: string;
		longitudePath?: string;
	}>;

	/** Timestamps */
	created: string;
	updated: string;
}

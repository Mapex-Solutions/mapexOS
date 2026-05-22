import { z, StringAndNotBeEmpty, StringAndNotBeEmptyOrOptional, IsBoolean, IsStringDateFormat, IsMongoId, IsString, IsNumber, NumberIntAndPositive } from '@mapexos/validations';

/**
 * Protocol configuration schemas
 */
const ZodNoneConfigSchema = z.object({});

/**
 * MQTT identity schema.
 *
 * `authType` is the platform-declared credential mode the broker
 * enforces — password XOR cert. The asset is bound to ONE mode at a
 * time; the broker rejects any CONNECT that presents the wrong shape
 * (cert on a password-mode asset, missing cert on a cert-mode asset).
 *
 * `password` is plaintext, request-only and only meaningful when
 * `authType` is `password` — present on create / update bodies, NEVER
 * on response shapes (the platform bcrypts it server-side and the
 * read-model exposes only the hash, consumed by the broker plugin
 * for local bcrypt-compare on every password-mode CONNECT).
 *
 * In `cert` mode the broker reads the asset by its cert Subject CN
 * (`{orgId}:{assetUUID}`) — the username field on the wire can be
 * empty. The platform still stores a canonical username on the asset
 * so the broker has a stable lookup key when the device chooses to
 * send one anyway.
 */
/**
 * Cert TTL configuration — operator-declared validity window applied
 * when the platform signs this asset's MQTT device cert. Only
 * meaningful when authType=cert. Value/unit are split so the wizard
 * can render a (number, dropdown) pair without parsing strings.
 *
 * Bounds enforced server-side (1 day .. 10 years total). The UI mirrors
 * the unit list (`day | week | month | year`) for surface-level
 * validation; backend is the source of truth.
 */
const ZodCertTTLConfigSchema = z.object({
	value: z.number().int().min(1).max(3650),
	unit: z.enum(['day', 'week', 'month', 'year']),
});

const ZodMqttConfigSchema = z.object({
	clientId: StringAndNotBeEmpty,
	username: StringAndNotBeEmpty,
	authType: z.enum(['password', 'cert']),
	password: z.string().min(8).optional(),
	passwordHash: z.string().optional(),
	certTTL: ZodCertTTLConfigSchema.optional(),
});

const ZodProtocolTypeSchema = z.object({
	type: z.enum(['http', 'mqtt', 'lorawan']),
	http: ZodNoneConfigSchema.optional(),
	mqtt: ZodMqttConfigSchema.optional(),
});

/**
 * Active MQTT device cert metadata. Present on AssetResponse when the
 * asset is in certificate auth mode and has an issued cert; absent
 * otherwise. The UI uses presence to derive the auth-mode radio and
 * to gate "Generate / Rotate cert" affordances. PEM bytes are NEVER
 * exposed here — only the metadata fields kept on the entity.
 */
const ZodAssetCertificateSchema = z.object({
	serial: IsString,
	fingerprint: IsString,
	subjectCN: IsString,
	issuedAt: IsStringDateFormat,
	expiresAt: IsStringDateFormat,
});

/**
 * Health Monitor configuration schema
 *
 * All fields are field-level optional to align with the Go contract
 * (packages/contracts/services/assets/assets/dto.go::HealthMonitorConfig — all
 * pointers, all `omitempty`). When `enabled` is false or omitted, callers may
 * send a minimal payload like `{ enabled: false }` (or omit healthMonitor
 * entirely).
 *
 * When `enabled` is true, the superRefine enforces:
 *   - thresholdMinutes is provided (≥ 10).
 *
 * Empty arrays are valid — Asset will be monitored without router actions
 * (online/offline transitions still write to Mongo healthStatus and
 * ClickHouse asset_status_history; no publish to mapexos.route.execute).
 *
 * heartbeatMode chooses how the platform learns the device is alive:
 *   - 'implicit' (default): js-executor emits a heartbeat for every data
 *     event the device sends.
 *   - 'explicit': js-executor SKIPS implicit publishes; the platform learns
 *     liveness through one of two paths chosen by the asset's protocol:
 *       • MQTT-protocol assets — automatic via the device's MQTT connection.
 *         The NATS broker emits $SYS.ACCOUNT.*.CONNECT/DISCONNECT advisories
 *         that the platform consumes; nothing for the device firmware to
 *         implement beyond keeping the MQTT connection open.
 *       • HTTP-protocol assets — the device POSTs to
 *         /api/v1/heartbeat?ds={dataSourceId} with body { assetUUID } at
 *         its own cadence (same auth chain as /api/v1/events).
 *
 * Server-side Create/Update mirrors this rule via validateHealthMonitorConfig
 * (asset_helpers.go), so client-side validation is fast feedback only.
 *
 * Note: requiredMisses keeps `.default(3)` because the create/edit form
 * relies on a sensible default for the form-initial-state path; the default
 * is applied AFTER undefined checks, so parse({ enabled: false }) still
 * succeeds without callers needing to supply requiredMisses.
 */
export const ZodHealthMonitorConfigSchema = z
	.object({
		enabled: IsBoolean.optional(),
		thresholdMinutes: IsNumber.int().min(10).optional(),
		requiredMisses: IsNumber.int().min(1).optional().default(3),
		heartbeatMode: z.enum(['implicit', 'explicit']).optional().default('implicit'),
		offlineRouteGroupIds: z.array(IsString).max(3).optional().default([]),
		onlineRouteGroupIds: z.array(IsString).max(3).optional().default([]),
	})
	.superRefine((data, ctx) => {
		if (data.enabled === true) {
			if (data.thresholdMinutes === undefined) {
				ctx.addIssue({
					code: z.ZodIssueCode.custom,
					path: ['thresholdMinutes'],
					message:
						'thresholdMinutes is required when health monitoring is enabled',
				});
			}
		}
	});

/**
 * Asset ID parameter schema (for URL params - MongoDB ObjectID)
 */
export const ZodAssetIdSchema = z.object({
	assetId: IsMongoId,
});

/**
 * Asset UUID parameter schema (for device identifier)
 */
export const ZodAssetUUIDSchema = z.object({
	assetUUID: StringAndNotBeEmpty.min(5),
});

/**
 * Asset Create schema - Used for creating new assets
 *
 * NOTE: Category and AssetType removed - they come from AssetTemplate
 * Classification (manufacturer, model, category) is managed at template level
 */
export const ZodAssetCreateSchema = z.object({
	name: StringAndNotBeEmpty,
	enabled: IsBoolean,
	debugEnabled: IsBoolean.optional(),
	description: IsString.max(500).optional(),

	assetUUID: StringAndNotBeEmpty.min(5),
	assetTemplateId: IsMongoId,

	orgId: IsMongoId,
	routeGroupIds: z.array(IsMongoId).min(1).max(3),

	healthMonitor: ZodHealthMonitorConfigSchema.optional(),

	protocol: ZodProtocolTypeSchema,
	latitude: IsNumber.optional(),
	longitude: IsNumber.optional(),
}).refine((data) => {
	// Validate MQTT config when type is mqtt
	if (data.protocol.type === 'mqtt' && !data.protocol.mqtt) {
		return false;
	}
	return true;
}, {
	message: "MQTT configuration is required when protocol type is 'mqtt'",
	path: ['protocol', 'mqtt'],
}).refine((data) => {
	// Validate routeGroupIds uniqueness
	const uniqueIds = new Set(data.routeGroupIds);
	return uniqueIds.size === data.routeGroupIds.length;
}, {
	message: "RouteGroupIds must be unique (no duplicates)",
	path: ['routeGroupIds'],
});

/**
 * Asset Update schema - Used for updating existing assets (all fields optional)
 *
 * NOTE: Category and AssetType removed - they come from AssetTemplate
 */
export const ZodAssetUpdateSchema = z.object({
	name: StringAndNotBeEmpty.optional(),
	enabled: IsBoolean.optional(),
	debugEnabled: IsBoolean.optional(),
	description: IsString.max(500).optional(),

	assetUUID: StringAndNotBeEmpty.optional(),
	assetTemplateId: IsMongoId.optional(),

	orgId: IsMongoId.optional(),
	routeGroupIds: z.array(IsMongoId).min(1).max(3).optional(),

	healthMonitor: ZodHealthMonitorConfigSchema.optional(),

	protocol: ZodProtocolTypeSchema.optional(),
	latitude: IsNumber.optional(),
	longitude: IsNumber.optional(),
}).refine((data) => {
	// Validate MQTT config when type is mqtt
	if (data.protocol?.type === 'mqtt' && !data.protocol.mqtt) {
		return false;
	}
	return true;
}, {
	message: "MQTT configuration is required when protocol type is 'mqtt'",
	path: ['protocol', 'mqtt'],
}).refine((data) => {
	// Validate routeGroupIds uniqueness if provided
	if (data.routeGroupIds) {
		const uniqueIds = new Set(data.routeGroupIds);
		return uniqueIds.size === data.routeGroupIds.length;
	}
	return true;
}, {
	message: "RouteGroupIds must be unique (no duplicates)",
	path: ['routeGroupIds'],
});

/**
 * Asset Response schema - Returned from API
 *
 * NOTE: Template classification fields (categoryId, manufacturerId, modelId + names)
 * are populated via MongoDB $lookup aggregation from asset_templates collection
 */
export const ZodAssetResponseSchema = z.object({
	id: IsMongoId.optional(),
	name: StringAndNotBeEmpty.optional(),
	enabled: IsBoolean.optional(),
	debugEnabled: IsBoolean.optional(),
	description: StringAndNotBeEmptyOrOptional,

	assetUUID: StringAndNotBeEmptyOrOptional,
	assetTemplateId: StringAndNotBeEmptyOrOptional,

	// Template classification data (populated from AssetTemplate lookup)
	assetTemplateName: StringAndNotBeEmptyOrOptional,
	assetIdPath: StringAndNotBeEmptyOrOptional, // Path to extract asset UUID from payload (from template)
	categoryId: IsMongoId.optional(),
	categoryName: StringAndNotBeEmptyOrOptional,
	manufacturerId: IsMongoId.optional(),
	manufacturerName: StringAndNotBeEmptyOrOptional,
	modelId: IsMongoId.optional(),
	modelName: StringAndNotBeEmptyOrOptional,
	version: StringAndNotBeEmptyOrOptional,

	orgId: StringAndNotBeEmptyOrOptional,
	pathKey: StringAndNotBeEmptyOrOptional,
	customerId: StringAndNotBeEmptyOrOptional,
	routeGroupIds: z.array(IsString).optional(),
	routeGroupNames: z.array(IsString).optional(),

	healthMonitor: ZodHealthMonitorConfigSchema.optional(),
	healthStatus: z.enum(['online', 'offline', 'unknown']).optional(),
	healthStatusChangedAt: IsStringDateFormat.optional().nullable(),
	lastSeenAt: IsStringDateFormat.optional().nullable(),

	protocol: ZodProtocolTypeSchema.optional(),
	latitude: IsNumber.optional(),
	longitude: IsNumber.optional(),

	currentCert: ZodAssetCertificateSchema.optional(),

	created: IsStringDateFormat.optional(),
	updated: IsStringDateFormat.optional(),
});

/**
 * Internal API schemas - Used for MS-to-MS communication
 */

// Internal Get/Delete - orgId in X-Org-Id header, assetId in params
export const ZodAssetInternalIdSchema = z.object({
	assetId: IsMongoId,
});

// Internal Update - data in body, assetId in params, orgId in header
export const ZodAssetInternalUpdateSchema = z.object({
	data: ZodAssetUpdateSchema,
});

// Internal Get Scripts - assetUUID in params
export const ZodAssetUUIDParamSchema = z.object({
	assetUUID: StringAndNotBeEmpty.min(5),
});

// Asset Scripts Response - returned from internal getScripts endpoint
export const ZodAssetScriptsResponseSchema = z.object({
	id: IsString,
	name: IsString,
	assetUUID: IsString,
	assetTemplateId: IsString,
	scriptProcessor: StringAndNotBeEmptyOrOptional,
	scriptValidator: IsString,
	scriptConversion: IsString,
});

/**
 * Asset Query schema - Used for filtering/pagination
 *
 * NOTE: category and assetType removed (they don't exist in Asset entity)
 * Use categoryId, manufacturerId, modelId for template-based filtering
 */
export const ZodAssetQuerySchema = z.object({
	// BaseQueryDTO fields
	projection: StringAndNotBeEmptyOrOptional,
	page: NumberIntAndPositive.optional(),
	perPage: NumberIntAndPositive.max(100).optional(),
	sort: StringAndNotBeEmptyOrOptional,
	includeChildren: IsBoolean.optional(),

	// Module-specific filters
	name: IsString.max(100).optional(),
	enabled: IsBoolean.optional(),
	assetUUID: IsString.max(100).optional(),
	assetTemplateId: IsMongoId.optional(),

	// Classification filters (applied to asset_templates via $lookup)
	categoryId: IsMongoId.optional(),
	manufacturerId: IsMongoId.optional(),
	modelId: IsMongoId.optional(),

	// Health monitoring filter
	healthStatus: z.enum(['online', 'offline', 'unknown']).optional(),
});

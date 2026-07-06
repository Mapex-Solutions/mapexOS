import { z, StringAndNotBeEmpty, StringAndNotBeEmptyOrOptional, IsMongoId, IsString, IsStringDateFormat, NumberIntAndPositive } from '@mapexos/validations';

/**
 * OTA remote firmware update — HTTP DTO schemas.
 *
 * Mirror of Go: packages/contracts/services/ota/dtos/ota_dtos.go
 * Served by the assets service under /api/v1/ota.
 *
 * The edge-side contracts (DownlinkEnvelope, OTAUpdateCommand, the status
 * advisory) are Go-only mirror pairs (mapexGoKit/contracts) — no TS consumer,
 * so they are intentionally NOT declared here.
 */

/**
 * Plan lifecycle statuses (the platform state machine).
 */
export const ZodOTAPlanStatusSchema = z.enum([
	'SCHEDULED',
	'IN_PROGRESS',
	'COMPLETED',
	'CLOSED',
	'CANCELED',
	'ABORTED',
]);

/**
 * Per-device execution states. Terminal: UPDATED | FAILED | TIMED_OUT.
 */
export const ZodOTAExecutionStateSchema = z.enum([
	'QUEUED',
	'INITIATED',
	'DOWNLOADING',
	'DOWNLOADED',
	'VERIFIED',
	'UPDATING',
	'UPDATED',
	'FAILED',
	'TIMED_OUT',
]);

/**
 * Firmware artifact statuses. Only READY/ACTIVE are referenceable by a plan.
 */
export const ZodOTAFirmwareStatusSchema = z.enum([
	'PENDING_UPLOAD',
	'READY',
	'ACTIVE',
	'DEPRECATED',
	'REVOKED',
	'ABANDONED',
	'PURGED',
]);

/**
 * Starts a presigned firmware upload. `version` is the target version the
 * device receives in the OTA command and reports after applying; `sha256` is
 * the hex digest the platform validates against the stored object on finalize.
 */
export const ZodFirmwareInitRequestSchema = z.object({
	targetTemplateId: IsMongoId,
	version: StringAndNotBeEmpty.max(64),
	filename: StringAndNotBeEmpty.max(255),
	size: NumberIntAndPositive,
	sha256: StringAndNotBeEmpty.max(128),
});

/**
 * The artifact id + the short-TTL presigned PUT URL the client uses to upload
 * the .bin directly to object storage (bytes never traverse the Asset MS).
 */
export const ZodFirmwareInitResponseSchema = z.object({
	firmwareId: IsString,
	uploadUrl: IsString,
});

/**
 * Pacing and auto-abort of a plan rollout.
 */
export const ZodOTARolloutConfigSchema = z.object({
	ratePerMinute: NumberIntAndPositive,
	abortThresholdPct: z.number().int().min(0).max(100),
	abortMinExecuted: z.number().int().min(0),
});

/**
 * Creates an OTA plan migrating the selected assets from the source template
 * to the target template carried by the firmware. The asset's template
 * switches only when the device reports updated.
 */
export const ZodOTAPlanCreateRequestSchema = z.object({
	firmwareId: IsMongoId,
	sourceTemplateId: IsMongoId,
	assetIds: z.array(IsMongoId).min(1),
	// Optional: absent means "run immediately" (the platform stamps now)
	startAt: IsStringDateFormat.optional(),
	maxTime: IsStringDateFormat,
	name: StringAndNotBeEmpty.max(100),
	description: IsString.max(500).optional(),
	rolloutConfig: ZodOTARolloutConfigSchema,
});

/**
 * Partial plan update (only provided fields are applied).
 */
export const ZodOTAPlanUpdateRequestSchema = z.object({
	name: StringAndNotBeEmptyOrOptional,
	description: IsString.max(500).optional(),
});

/**
 * Per-plan aggregate counters (the plan holds only counters; per-device
 * progress lives on the executions).
 */
export const ZodOTAPlanCountersSchema = z.object({
	total: z.number().int(),
	succeeded: z.number().int(),
	failed: z.number().int(),
	timedOut: z.number().int(),
});

/**
 * Denormalized firmware artifact detail — present only on the plan DETAIL
 * response (GetPlan), so the UI shows the file facts without a second fetch.
 */
export const ZodOTAFirmwareInfoSchema = z.object({
	id: IsString,
	version: IsString,
	filename: IsString,
	size: z.number().int(),
	checksum: IsString,
	checksumAlgorithm: IsString,
	status: ZodOTAFirmwareStatusSchema,
	targetTemplateId: IsString,
});

/**
 * Presigned firmware download response. `url` is a short-TTL GET URL the browser
 * uses to fetch the artifact directly from object storage; `expiresAt` is when
 * it stops working.
 */
export const ZodOTAFirmwareDownloadResponseSchema = z.object({
	url: IsString,
	filename: IsString,
	expiresAt: IsStringDateFormat,
});

/**
 * Plan read model for list/detail responses. `firmware` is present only on the
 * detail response.
 */
export const ZodOTAPlanResponseSchema = z.object({
	id: IsString,
	orgId: IsString,
	name: IsString,
	description: IsString.optional().default(''),
	firmwareId: IsString,
	sourceTemplateId: IsString,
	targetTemplateId: IsString,
	status: ZodOTAPlanStatusSchema,
	startAt: IsString,
	maxTime: IsString,
	counters: ZodOTAPlanCountersSchema,
	firmware: ZodOTAFirmwareInfoSchema.optional(),
	created: IsString,
	updated: IsString,
});

/**
 * Per-device execution read model.
 */
export const ZodOTAExecutionResponseSchema = z.object({
	id: IsString,
	planId: IsString,
	assetId: IsString,
	protocol: IsString.optional().default(''),
	state: ZodOTAExecutionStateSchema,
	percentage: z.number().int(),
	error: IsString.optional(),
	attempts: z.number().int(),
	queuedAt: IsString,
	updated: IsString,
});

/**
 * Query for the plans list.
 */
export const ZodOTAPlanQuerySchema = z.object({
	page: NumberIntAndPositive.optional(),
	perPage: NumberIntAndPositive.max(100).optional(),
	sort: StringAndNotBeEmptyOrOptional,
	name: IsString.max(100).optional(),
	status: ZodOTAPlanStatusSchema.optional(),
});

/**
 * Query for a plan's executions list.
 */
export const ZodOTAExecutionQuerySchema = z.object({
	page: NumberIntAndPositive.optional(),
	perPage: NumberIntAndPositive.max(100).optional(),
	sort: StringAndNotBeEmptyOrOptional,
	state: ZodOTAExecutionStateSchema.optional(),
	assetId: IsMongoId.optional(),
});

/**
 * Path params.
 */
export const ZodOTAPlanIdSchema = z.object({
	planId: IsMongoId,
});

export const ZodOTAFirmwareIdSchema = z.object({
	firmwareId: IsMongoId,
});

import { z, StringAndNotBeEmpty, NumberIntAndPositive, IsBoolean, IsMongoId, IsString, IsStringDateFormat } from '@mapexos/validations';

/**
 * Zod schema for retention policy upsert (PUT body).
 * Maps to Go: RetentionPolicyUpsert in contracts/services/events/retention/dto.go
 */
export const ZodRetentionPolicyUpsertSchema = z.object({
	/** ClickHouse table identifier (e.g., "events", "eventsRaw") */
	type: StringAndNotBeEmpty,

	/** Human-readable label for the policy */
	name: StringAndNotBeEmpty,

	/** Number of days to retain data */
	retentionDays: NumberIntAndPositive,

	/** Whether this retention policy is active */
	enabled: IsBoolean.optional(),
});

/**
 * Zod schema for retention policy API response.
 * Maps to Go: RetentionPolicyResponse in contracts/services/events/retention/dto.go
 */
export const ZodRetentionPolicyResponseSchema = z.object({
	/** MongoDB ObjectId */
	id: IsMongoId.optional(),

	/** Human-readable label */
	name: IsString.optional(),

	/** ClickHouse table identifier */
	type: IsString.optional(),

	/** Number of days to retain data */
	retentionDays: z.number().int().min(0).optional(),

	/** Organization ID */
	orgId: IsString.optional(),

	/** Organization hierarchical path key */
	pathKey: IsString.optional(),

	/** Whether this retention policy is active */
	enabled: IsBoolean.optional(),

	/** Creation timestamp */
	created: IsStringDateFormat.optional(),

	/** Last update timestamp */
	updated: IsStringDateFormat.optional(),
});

/**
 * Zod schema for retention policy query parameters (GET query).
 * Maps to Go: RetentionPolicyQuery in contracts/services/events/retention/dto.go
 */
export const ZodRetentionPolicyQuerySchema = z.object({
	/** Page number for pagination */
	page: NumberIntAndPositive.optional(),

	/** Items per page */
	perPage: NumberIntAndPositive.optional(),

	/** Sort field */
	sort: IsString.optional(),

	/** Include children organizations */
	includeChildren: IsBoolean.optional(),

	/** Filter by type (comma-separated for $in) */
	type: IsString.optional(),
});

/**
 * Zod schema for retention policy route params.
 * Maps to Go: RetentionPolicyParams in contracts/services/events/retention/dto.go
 */
export const ZodRetentionPolicyParamsSchema = z.object({
	/** MongoDB ObjectId of the retention policy */
	retentionPolicyId: IsMongoId,
});

/**
 * Zod schema for paginated retention policy list response.
 */
export const ZodRetentionPolicyPaginatedResultSchema = z.object({
	/** List of retention policies */
	items: z.array(ZodRetentionPolicyResponseSchema),

	/** Pagination metadata */
	pagination: z.object({
		page: z.number(),
		perPage: z.number(),
		totalItems: z.number(),
	}),
});

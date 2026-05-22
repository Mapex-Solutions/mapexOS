import { z, StringAndNotBeEmpty, StringAndNotBeEmptyOrOptional, IsBoolean, IsString, IsNumber, NumberIntAndPositive, IsStringDateFormat, IsMongoId, IsUrl } from '@mapexos/validations';

/**
 * Auth type schemas - nested authentication configurations
 */
const ZodAuthAPIKeySchema = z.object({
	type: z.enum(['header', 'query']),
	fieldName: StringAndNotBeEmpty,
	key: IsString.min(20),
});

const ZodAuthJWTSchema = z.object({
	secret: IsString.min(6),
	algorithms: z.enum(['HS256', 'HS512']),
	headerName: IsString.min(1).optional(),
});

const ZodAuthIPWhitelistSchema = z.object({
	cidrs: z.array(IsString).min(1),
});

const ZodAuthOAuth2Schema = z.object({
	jwksURL: IsUrl,
	clientId: IsString.min(3).optional(),
	clientSecret: IsString.min(6).optional(),
});

const ZodAuthNoneSchema = z.object({});

/**
 * DataSource Auth schema - discriminated union based on type
 */
const ZodDataSourceAuthSchema = z.object({
	type: z.enum(['apiKey', 'jwt', 'ip_whitelist', 'oauth2', 'none']),
	apiKey: ZodAuthAPIKeySchema.optional(),
	jwt: ZodAuthJWTSchema.optional(),
	ipWhitelist: ZodAuthIPWhitelistSchema.optional(),
	oauth2: ZodAuthOAuth2Schema.optional(),
	none: ZodAuthNoneSchema.optional(),
}).refine((data) => {
	// Validate that corresponding field is provided based on type
	switch (data.type) {
		case 'apiKey':
			return !!data.apiKey;
		case 'jwt':
			return !!data.jwt;
		case 'ip_whitelist':
			return !!data.ipWhitelist;
		case 'oauth2':
			return !!data.oauth2;
		case 'none':
			return true; // none doesn't require additional data
		default:
			return false;
	}
}, {
	message: "Auth configuration must match the specified type",
	path: ['type'],
});

/**
 * Working Hours schema
 */
const ZodWorkingHoursSchema = z.object({
	enabled: IsBoolean,
	days: z.array(IsNumber.int().min(0).max(6)).optional(),
	startAt: IsString.optional(),
	endAt: IsString.optional(),
	timeZone: IsString.optional(),
});

/**
 * Rate Limit schema
 */
const ZodRateLimitSchema = z.object({
	type: z.enum(['second', 'minute', 'hour']),
	value: NumberIntAndPositive,
	burstCapacity: NumberIntAndPositive,
	actionOnExceed: z.enum(['drop', 'queue']),
});

/**
 * Asset Bind schema
 */
export const ZodAssetBindSchema = z.object({
	type: z.enum(['fixedAssetId', 'uuidField']),
	data: z.object({
		uuidField: z.array(IsString).min(1).optional(),
		assetId: StringAndNotBeEmptyOrOptional,
	}),
}).refine((data) => {
	// Validate that corresponding field is provided based on type
	if (data.type === 'fixedAssetId') {
		return !!data.data.assetId;
	} else if (data.type === 'uuidField') {
		return !!data.data.uuidField && data.data.uuidField.length > 0;
	}
	return false;
}, {
	message: "Asset bind configuration must match the specified type",
	path: ['type'],
});

/**
 * DataSource ID parameter schema (for URL params - MongoDB ObjectID)
 */
export const ZodDataSourceIdSchema = z.object({
	dataSourceId: IsMongoId,
});

/**
 * DataSource Create schema - Used for creating new data sources
 */
export const ZodDataSourceCreateSchema = z.object({
	name: StringAndNotBeEmpty,
	enabled: IsBoolean,
	description: IsString.max(500).optional(),
	mode: z.enum(['pull', 'push', 'X']),
	protocol: z.enum(['mqtt', 'http']),
	workingHours: ZodWorkingHoursSchema.optional(),
	rateLimit: ZodRateLimitSchema.optional(),
	auth: ZodDataSourceAuthSchema,
	assetBind: ZodAssetBindSchema,
});

/**
 * DataSource Update schema - Used for updating existing data sources (all fields optional)
 */
export const ZodDataSourceUpdateSchema = z.object({
	name: StringAndNotBeEmpty.optional(),
	enabled: IsBoolean.optional(),
	description: IsString.max(500).optional(),
	mode: z.enum(['pull', 'push', 'X']).optional(),
	protocol: z.enum(['mqtt', 'http']).optional(),
	workingHours: ZodWorkingHoursSchema.optional(),
	rateLimit: ZodRateLimitSchema.optional(),
	auth: ZodDataSourceAuthSchema.optional(),
	assetBind: ZodAssetBindSchema.optional(),
});

/**
 * DataSource Response schema - Returned from API
 */
export const ZodDataSourceResponseSchema = z.object({
	id: IsMongoId.optional(),
	name: IsString.optional(),
	enabled: IsBoolean.optional(),
	description: IsString.optional(),
	mode: IsString.optional(),
	protocol: IsString.optional(),
	workingHours: ZodWorkingHoursSchema.optional(),
	rateLimit: ZodRateLimitSchema.optional(),
	auth: ZodDataSourceAuthSchema.optional(),
	assetBind: ZodAssetBindSchema.optional(),
	created: IsStringDateFormat.optional(),
	updated: IsStringDateFormat.optional(),
});

/**
 * DataSource Query schema - Used for filtering/pagination
 */
export const ZodDataSourceQuerySchema = z.object({
	// BaseQueryDTO fields
	projection: IsString.optional(),
	page: IsNumber.int().min(1).optional(),
	perPage: IsNumber.int().min(1).max(100).optional(),
	sort: IsString.optional(),
	includeChildren: IsBoolean.optional(),

	// Module-specific filters
	name: IsString.max(100).optional(),
	enabled: IsBoolean.optional(),
	mode: z.enum(['pull', 'push']).optional(),
	auth: z.enum(['apiKey', 'jwt', 'ip_whitelist', 'oauth2', 'none']).optional(),
	assetBind: z.enum(['fixedAssetId', 'uuidField']).optional(),
});

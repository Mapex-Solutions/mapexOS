import { z, StringAndNotBeEmpty, StringAndNotBeEmptyOrOptional, IsBoolean, IsStringDateFormat, IsMongoId, IsString, IsNumber, NumberIntAndPositive } from '@mapexos/validations';

/**
 * Organization ID parameter schema (for URL params)
 */
export const ZodOrganizationIdSchema = z.object({
	organizationId: IsMongoId,
});

/**
 * Address schema
 */
export const ZodAddressSchema = z.object({
	city: StringAndNotBeEmpty,
	state: StringAndNotBeEmpty,
	country: StringAndNotBeEmpty,
	zipCode: StringAndNotBeEmpty,
});

/**
 * AuthConfig schema
 */
export const ZodAuthConfigSchema = z.object({
	providerType: z.enum(['keycloak', 'internal']),
	issuerUrl: StringAndNotBeEmptyOrOptional,
	clientId: StringAndNotBeEmptyOrOptional,
	jwtClaimMappings: z.record(IsString, IsString).optional(),
	metadata: z.record(IsString, z.any()).optional(),
});

/**
 * AccessPolicy schema
 */
export const ZodAccessPolicySchema = z.object({
	rolePolicy: z.enum(['strict', 'merge']),
	defaultScope: z.enum(['local', 'recursive']),
});

/**
 * Organization Create schema
 */
export const ZodOrganizationCreateSchema = z.object({
	name: IsString.min(3).max(150),
	type: z.enum(['vendor', 'customer', 'site', 'building', 'floor', 'zone']),
	parentOrgId: IsMongoId.optional(),
	logo: IsString.url().optional(),
	enabled: IsBoolean,
	address: ZodAddressSchema.optional(), // Optional - only for customer and site types
	phone: StringAndNotBeEmptyOrOptional, // E164 validated by backend, optional per Go DTO
	authConfig: ZodAuthConfigSchema,
	accessPolicy: ZodAccessPolicySchema,
});

/**
 * Organization Update schema
 */
export const ZodOrganizationUpdateSchema = z.object({
	name: IsString.min(3).max(150).optional(),
	type: z.enum(['vendor', 'customer', 'site', 'building', 'floor', 'zone']).optional(),
	parentOrgId: IsMongoId.optional(),
	logo: IsString.url().optional(),
	enabled: IsBoolean.optional(),
	address: ZodAddressSchema.optional(),
	phone: StringAndNotBeEmptyOrOptional,
	authConfig: ZodAuthConfigSchema.optional(),
	accessPolicy: ZodAccessPolicySchema.optional(),
});

/**
 * Organization Query schema - Used for filtering/pagination
 */
export const ZodOrganizationQuerySchema = z.object({
	// BaseQueryDTO fields
	projection: StringAndNotBeEmptyOrOptional,
	page: NumberIntAndPositive.optional(),
	perPage: NumberIntAndPositive.max(100).optional(),
	sort: StringAndNotBeEmptyOrOptional,
	includeChildren: IsBoolean.optional(),

	// Module-specific filters
	type: z.enum(['vendor', 'customer', 'site', 'building', 'floor', 'zone']).optional(),
	parentOrgId: IsMongoId.optional(),
	name: IsString.max(150).optional(),
	enabled: IsBoolean.optional(),
	depth: IsNumber.int().min(0).max(10).optional(),
});

/**
 * Organization Response schema
 */
export const ZodOrganizationResponseSchema = z.object({
	id: IsMongoId.optional(),
	name: StringAndNotBeEmptyOrOptional,
	type: StringAndNotBeEmptyOrOptional,
	parentOrgId: StringAndNotBeEmptyOrOptional,
	code: StringAndNotBeEmptyOrOptional,
	pathKey: StringAndNotBeEmptyOrOptional,
	depth: IsNumber.optional(),
	customerId: StringAndNotBeEmptyOrOptional,
	childCount: IsNumber.optional(),
	logo: StringAndNotBeEmptyOrOptional,
	enabled: IsBoolean.optional(),
	address: ZodAddressSchema.optional(),
	phone: StringAndNotBeEmptyOrOptional,
	authConfig: ZodAuthConfigSchema.optional(),
	accessPolicy: ZodAccessPolicySchema.optional(),
	created: IsStringDateFormat.optional(),
	updated: IsStringDateFormat.optional(),
});

/**
 * Tree Query schema - Cursor-based pagination for tree navigation
 */
export const ZodTreeQuerySchema = z.object({
	cursor: IsString.length(24).optional(), // ObjectID hex string
	direction: z.enum(['next', 'previous']).optional(),
	limit: IsNumber.int().min(1).max(300).optional(),
	sortAsc: IsBoolean.optional(),
});

/**
 * Tree Item schema - Minimal organization data for tree UI
 */
export const ZodTreeItemSchema = z.object({
	id: IsMongoId,
	name: StringAndNotBeEmpty,
	type: IsString,
});

/**
 * Tree Response schema - Paginated tree items with cursor metadata
 */
export const ZodTreeResponseSchema = z.object({
	items: z.array(ZodTreeItemSchema),
	cursor: z.object({
		current: StringAndNotBeEmptyOrOptional,
		next: StringAndNotBeEmptyOrOptional,
		previous: StringAndNotBeEmptyOrOptional,
		hasNext: IsBoolean,
		hasPrevious: IsBoolean,
	}),
});

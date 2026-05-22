import { z, StringAndNotBeEmpty, StringAndNotBeEmptyOrOptional, IsBoolean, IsStringDateFormat, IsMongoId, IsString, NumberIntAndPositive } from '@mapexos/validations';

/**
 * AuthProvider schemas
 */
export const ZodAuthProviderSchema = z.object({
	type: z.enum(['internal', 'google', 'github', 'microsoft', 'keycloak']),
	externalId: StringAndNotBeEmptyOrOptional,
	metadata: z.record(IsString, z.any()).optional(),
});

export const ZodAuthProviderUpdateSchema = z.object({
	type: z.enum(['internal', 'google', 'github', 'microsoft', 'keycloak']).optional(),
	externalId: StringAndNotBeEmptyOrOptional,
	metadata: z.record(IsString, z.any()).optional(),
});

/**
 * User ID parameter schema (for URL params)
 */
export const ZodUserIdSchema = z.object({
	userId: IsMongoId,
});

/**
 * User Create schema
 * AuthProvider removed - V1 always uses internal auth.
 * Next version: auth provider will be determined by the customer's Organization.AuthConfig
 */
export const ZodUserCreateSchema = z.object({
	email: IsString.email().max(254),
	password: IsString.min(8).max(72).optional(),
	changePasswordNextLogin: IsBoolean,
	firstName: IsString.min(2).max(100),
	lastName: IsString.min(2).max(100),
	phone: StringAndNotBeEmptyOrOptional, // E164 format validated by backend
	jobTitle: IsString.max(120).optional(),
	enabled: IsBoolean,
	avatar: IsString.url().optional(),
	startTour: IsBoolean,
});

/**
 * User Update schema
 * AuthProvider removed - V1 always uses internal auth.
 * Next version: auth provider will be determined by the customer's Organization.AuthConfig
 */
export const ZodUserUpdateSchema = z.object({
	email: IsString.email().max(254).optional(),
	password: IsString.min(8).max(72).optional(),
	changePasswordNextLogin: IsBoolean.optional(),
	firstName: IsString.min(2).max(100).optional(),
	lastName: IsString.min(2).max(100).optional(),
	phone: StringAndNotBeEmptyOrOptional,
	jobTitle: IsString.max(120).optional(),
	enabled: IsBoolean.optional(),
	avatar: IsString.url().optional(),
	startTour: IsBoolean.optional(),
});

/**
 * User Query schema - Used for filtering/pagination
 */
export const ZodUserQuerySchema = z.object({
	// BaseQueryDTO fields
	projection: StringAndNotBeEmptyOrOptional,
	page: NumberIntAndPositive.optional(),
	perPage: NumberIntAndPositive.max(100).optional(),
	sort: StringAndNotBeEmptyOrOptional,
	includeChildren: IsBoolean.optional(),

	// Module-specific filters
	email: IsString.max(254).optional(),
	firstName: IsString.max(100).optional(),
	lastName: IsString.max(100).optional(),
	enabled: IsBoolean.optional(),
});

/**
 * User Group Info schema (for detail view)
 */
export const ZodUserGroupInfoSchema = z.object({
	id: IsMongoId,
	name: StringAndNotBeEmpty,
	description: StringAndNotBeEmptyOrOptional,
});

/**
 * User Membership Info schema (for detail view)
 */
export const ZodUserMembershipInfoSchema = z.object({
	orgId: IsMongoId,
	orgName: StringAndNotBeEmpty,
	orgType: StringAndNotBeEmpty,
	scope: z.enum(['local', 'recursive']),
	roleNames: z.array(IsString),
	via: StringAndNotBeEmpty, // "direct" or "Group: {groupName}"
});

/**
 * User Response schema
 */
export const ZodUserResponseSchema = z.object({
	id: IsMongoId.optional(),
	email: StringAndNotBeEmptyOrOptional,
	changePasswordNextLogin: IsBoolean.optional(),
	authProvider: ZodAuthProviderSchema.optional(),
	firstName: StringAndNotBeEmptyOrOptional,
	lastName: StringAndNotBeEmptyOrOptional,
	phone: StringAndNotBeEmptyOrOptional,
	jobTitle: StringAndNotBeEmptyOrOptional,
	enabled: IsBoolean.optional(),
	avatar: StringAndNotBeEmptyOrOptional,
	startTour: IsBoolean.optional(),
	created: IsStringDateFormat.optional(),
	updated: IsStringDateFormat.optional(),

	// Enriched fields (populated by service layer)
	groupsCount: z.number().int().nonnegative().optional(),
	groups: z.array(ZodUserGroupInfoSchema).optional(),
	memberships: z.array(ZodUserMembershipInfoSchema).optional(),
});

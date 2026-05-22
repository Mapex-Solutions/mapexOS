import { z, StringAndNotBeEmpty, StringAndNotBeEmptyOrOptional, IsBoolean, IsStringDateFormat, IsMongoId, IsString, NumberIntAndPositive } from '@mapexos/validations';

/**
 * Group ID parameter schema (for URL params)
 */
export const ZodGroupIdSchema = z.object({
	groupId: IsMongoId,
});

/**
 * Group Create schema
 */
export const ZodGroupCreateSchema = z.object({
	name: IsString.min(3).max(150),
	description: IsString.max(500).optional(),
	enabled: IsBoolean,
	orgId: IsMongoId.optional(),
	roleIds: z.array(IsMongoId).min(1, 'At least one role is required'),
});

/**
 * Group Update schema
 */
export const ZodGroupUpdateSchema = z.object({
	name: IsString.min(3).max(150).optional(),
	description: IsString.max(500).optional(),
	enabled: IsBoolean.optional(),
});

/**
 * Group Query schema - Used for filtering/pagination
 */
export const ZodGroupQuerySchema = z.object({
	// BaseQueryDTO fields
	projection: StringAndNotBeEmptyOrOptional,
	page: NumberIntAndPositive.optional(),
	perPage: NumberIntAndPositive.max(100).optional(),
	sort: StringAndNotBeEmptyOrOptional,
	includeChildren: IsBoolean.optional(),

	// Module-specific filters
	name: IsString.max(150).optional(),
	enabled: IsBoolean.optional(),
	memberId: IsMongoId.optional(),
});

/**
 * Group Response schema
 */
export const ZodGroupResponseSchema = z.object({
	id: IsMongoId.optional(),
	name: StringAndNotBeEmptyOrOptional,
	description: StringAndNotBeEmptyOrOptional,
	membersCount: z.number().int().nonnegative().optional(),
	roleIds: z.array(IsMongoId).optional(),
	enabled: IsBoolean.optional(),
	orgId: StringAndNotBeEmptyOrOptional,
	pathKey: StringAndNotBeEmptyOrOptional,
	created: IsStringDateFormat.optional(),
	updated: IsStringDateFormat.optional(),
});

/**
 * Group Members Query schema - Used for paginated member listing
 */
export const ZodGroupMembersQuerySchema = z.object({
	page: NumberIntAndPositive.optional(),
	perPage: NumberIntAndPositive.max(100).optional(),
});

/**
 * Group Member Response schema
 */
export const ZodGroupMemberResponseSchema = z.object({
	id: IsMongoId.optional(),
	userId: IsMongoId.optional(),
	groupId: IsMongoId.optional(),
	orgId: IsMongoId.optional(),
	addedAt: IsStringDateFormat.optional(),
	addedBy: IsMongoId.optional(),
	created: IsStringDateFormat.optional(),
	// Denormalized user info (may be populated by service)
	userEmail: StringAndNotBeEmptyOrOptional,
	userFirstName: StringAndNotBeEmptyOrOptional,
	userLastName: StringAndNotBeEmptyOrOptional,
});

/**
 * Group Member Add schema - Used to add a member to a group
 */
export const ZodGroupMemberAddSchema = z.object({
	userId: IsMongoId,
});

/**
 * Group Member ID schema - Used for remove member operations
 */
export const ZodGroupMemberIdSchema = z.object({
	groupId: IsMongoId,
	userId: IsMongoId,
});

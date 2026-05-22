import { z, StringAndNotBeEmptyOrOptional, IsBoolean, IsStringDateFormat, IsMongoId, IsString, NumberIntAndPositive } from '@mapexos/validations';

/**
 * Membership ID parameter schema (for URL params)
 */
export const ZodMembershipIdSchema = z.object({
	membershipId: IsMongoId,
});

/**
 * Membership Create schema
 */
export const ZodMembershipCreateSchema = z.object({
	assigneeType: z.enum(['user', 'group']),
	assigneeId: IsMongoId,
	orgId: IsMongoId,
	roleIds: z.array(IsMongoId).min(1),
	scope: z.enum(['local', 'recursive']),
	enabled: IsBoolean,
});

/**
 * Membership Update schema
 */
export const ZodMembershipUpdateSchema = z.object({
	roleIds: z.array(IsMongoId).min(1).optional(),
	scope: z.enum(['local', 'recursive']).optional(),
	enabled: IsBoolean.optional(),
});

/**
 * Membership Query schema - Used for filtering/pagination
 */
export const ZodMembershipQuerySchema = z.object({
	// BaseQueryDTO fields
	projection: StringAndNotBeEmptyOrOptional,
	page: NumberIntAndPositive.optional(),
	perPage: NumberIntAndPositive.max(100).optional(),
	sort: StringAndNotBeEmptyOrOptional,
	includeChildren: IsBoolean.optional(),

	// Module-specific filters
	assigneeType: z.enum(['user', 'group']).optional(),
	assigneeId: IsMongoId.optional(),
	userId: IsMongoId.optional(),
	roleId: IsMongoId.optional(),
	scope: z.enum(['local', 'recursive']).optional(),
	enabled: IsBoolean.optional(),
});

/**
 * Membership Response schema
 */
export const ZodMembershipResponseSchema = z.object({
	id: IsMongoId.optional(),
	assigneeType: z.enum(['user', 'group']).optional(),
	assigneeId: IsMongoId.optional(),
	orgId: IsMongoId.optional(),
	orgPathKey: StringAndNotBeEmptyOrOptional,
	customerId: IsMongoId.optional(),
	roleIds: z.array(IsMongoId).optional(),
	scope: z.enum(['local', 'recursive']).optional(),
	enabled: IsBoolean.optional(),
	created: IsStringDateFormat.optional(),
	updated: IsStringDateFormat.optional(),
});


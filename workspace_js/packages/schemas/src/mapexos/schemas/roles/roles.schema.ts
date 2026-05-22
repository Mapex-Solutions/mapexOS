import { z, StringAndNotBeEmpty, StringAndNotBeEmptyOrOptional, IsBoolean, IsStringDateFormat, IsMongoId, IsString, NumberIntAndPositive } from '@mapexos/validations';

/**
 * Role ID parameter schema (for URL params)
 */
export const ZodRoleIdSchema = z.object({
	roleId: IsMongoId,
});

/**
 * Role Create schema
 */
export const ZodRoleCreateSchema = z.object({
	name: IsString.min(3).max(100),
	description: IsString.max(500).optional(),
	permissions: z.array(StringAndNotBeEmpty).min(1),
	isSystem: IsBoolean,
	isTemplate: IsBoolean.optional(),
	orgId: IsMongoId.optional(),
	pathKey: StringAndNotBeEmpty,
	scope: z.enum(['global', 'local']),
}).refine((data) => {
	// If organization role (isSystem=false and orgId exists), scope must be valid
	if (!data.isSystem && data.orgId) {
		if (data.scope !== 'global' && data.scope !== 'local') {
			return false;
		}
	}
	return true;
}, {
	message: "scope must be 'global' or 'local' for organization roles",
	path: ['scope'],
});

/**
 * Role Update schema
 */
export const ZodRoleUpdateSchema = z.object({
	name: IsString.min(3).max(100).optional(),
	description: IsString.max(500).optional(),
	permissions: z.array(StringAndNotBeEmpty).min(1).optional(),
});

/**
 * Role Query schema - Used for filtering/pagination
 */
export const ZodRoleQuerySchema = z.object({
	// BaseQueryDTO fields
	projection: StringAndNotBeEmptyOrOptional,
	page: NumberIntAndPositive.optional(),
	perPage: NumberIntAndPositive.max(100).optional(),
	sort: StringAndNotBeEmptyOrOptional,
	includeChildren: IsBoolean.optional(),

	// Module-specific filters
	name: IsString.max(100).optional(),
	isSystem: IsBoolean.optional(),
	isTemplate: IsBoolean.optional(),
	scope: z.enum(['global', 'local']).optional(),
	permission: StringAndNotBeEmptyOrOptional,
});

/**
 * Role Response schema
 */
export const ZodRoleResponseSchema = z.object({
	id: IsMongoId.optional(),
	name: StringAndNotBeEmptyOrOptional,
	description: StringAndNotBeEmptyOrOptional,
	permissions: z.array(IsString).optional(),
	isSystem: IsBoolean.optional(),
	isTemplate: IsBoolean.optional(),
	orgId: StringAndNotBeEmptyOrOptional,
	pathKey: StringAndNotBeEmptyOrOptional,
	scope: StringAndNotBeEmptyOrOptional,
	created: IsStringDateFormat.optional(),
	updated: IsStringDateFormat.optional(),
});

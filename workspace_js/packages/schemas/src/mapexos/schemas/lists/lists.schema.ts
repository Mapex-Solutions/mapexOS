import { z, StringAndNotBeEmpty, IsBoolean, IsStringDateFormat, IsMongoId, IsRecord, IsString, NumberIntAndPositive } from '@mapexos/validations';

/**
 * List ID parameter schema (for URL params)
 */
export const ZodListIdSchema = z.object({
	listId: IsMongoId,
});

/**
 * List Create schema
 */
export const ZodListCreateSchema = z.object({
	type: StringAndNotBeEmpty.max(100),
	name: StringAndNotBeEmpty.max(254),
	value: StringAndNotBeEmpty.max(254),
	enabled: IsBoolean,
	parentId: IsMongoId.optional(),
	metadata: IsRecord.optional(),
	isSystem: IsBoolean,
	isTemplate: IsBoolean,
	orgId: IsMongoId.optional(),
	pathKey: IsString.optional(),
});

/**
 * List Update schema
 */
export const ZodListUpdateSchema = z.object({
	name: IsString.max(254).optional(),
	value: IsString.max(254).optional(),
	enabled: IsBoolean.optional(),
	parentId: IsMongoId.optional(),
	metadata: IsRecord.optional(),
});

/**
 * List Query schema - Used for filtering/pagination
 */
export const ZodListQuerySchema = z.object({
	// BaseQueryDTO fields
	projection: IsString.optional(),
	page: NumberIntAndPositive.optional(),
	perPage: NumberIntAndPositive.max(100).optional(),
	sort: IsString.optional(),
	includeChildren: IsBoolean.optional(),

	// Module-specific filters
	type: IsString.optional(),
	isSystem: IsBoolean.optional(),
	isTemplate: IsBoolean.optional(),
	name: IsString.max(254).optional(),
	enabled: IsBoolean.optional(),
	parentId: IsMongoId.optional(),
});

/**
 * List Response schema
 */
export const ZodListResponseSchema = z.object({
	id: IsMongoId.optional(),
	type: IsString.optional(),
	name: IsString.optional(),
	value: IsString.optional(),
	enabled: IsBoolean.optional(),
	metadata: IsRecord.optional(),
	parentId: IsMongoId.optional(),
	parentName: IsString.optional(),
	parentType: IsString.optional(),
	isSystem: IsBoolean.optional(),
	isTemplate: IsBoolean.optional(),
	orgId: IsMongoId.optional(),
	pathKey: IsString.optional(),
	created: IsStringDateFormat.optional(),
	updated: IsStringDateFormat.optional(),
});

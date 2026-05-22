import { z, StringAndNotBeEmpty, StringAndBeEmptyOrOptional, IsBoolean, IsMongoId, IsString, NumberIntAndPositive } from '@mapexos/validations';

/**
 * Dynamic Field schema - for typed event storage and querying
 */
export const ZodDynamicFieldSchema = z.object({
	fieldId: NumberIntAndPositive.optional(),
	field: StringAndNotBeEmpty,
	value: StringAndBeEmptyOrOptional,
	type: z.enum(['string', 'number', 'bool', 'date', 'geo']),
	status: NumberIntAndPositive.optional(),
	latitudePath: StringAndBeEmptyOrOptional,
	longitudePath: StringAndBeEmptyOrOptional,
});

/**
 * Asset Template ID parameter schema (for URL params - MongoDB ObjectID)
 */
export const ZodAssetTemplateIdSchema = z.object({
	assetTemplateId: IsMongoId,
});

/**
 * Asset Template Create schema - Used for creating new asset templates
 */
export const ZodAssetTemplateCreateSchema = z.object({
	name: StringAndNotBeEmpty,
	enabled: IsBoolean,
	description: StringAndNotBeEmpty.max(500).optional(),

	// Asset Classification (FLAT structure) - IDs are source of truth, Names are for UI performance
	categoryId: StringAndBeEmptyOrOptional,
	categoryName: IsString.max(254).optional(),
	manufacturerId: StringAndBeEmptyOrOptional,
	manufacturerName: IsString.max(254).optional(),
	modelId: StringAndBeEmptyOrOptional,
	modelName: IsString.max(254).optional(),
	version: IsString.max(100).optional(),

	// Template visibility flags
	isSystem: IsBoolean.optional(),
	isTemplate: IsBoolean.optional(),

	assetIdPath: StringAndNotBeEmpty,

	scriptProcessor: StringAndBeEmptyOrOptional,
	scriptValidator: StringAndBeEmptyOrOptional, 
	scriptConversion: StringAndNotBeEmpty,
	scriptTest: StringAndNotBeEmpty,

	// Available Fields (for Rule autocomplete support)
	availableFields: z.array(IsString).optional(),

	// Dynamic Fields (for typed event storage and querying)
	dynamicFields: z.array(ZodDynamicFieldSchema).optional(),

	// Multi-tenant fields (populated automatically by coverage middleware)
	orgId: StringAndBeEmptyOrOptional,
	pathKey: StringAndBeEmptyOrOptional,

	created: StringAndBeEmptyOrOptional,
	updated: StringAndBeEmptyOrOptional,
});

/**
 * Asset Template Update schema - Used for updating existing asset templates (all fields optional)
 */
export const ZodAssetTemplateUpdateSchema = z.object({
	name: StringAndBeEmptyOrOptional,
	enabled: IsBoolean.optional(),
	description: StringAndNotBeEmpty.max(500).optional(),

	// Asset Classification (FLAT structure)
	categoryId: StringAndBeEmptyOrOptional,
	categoryName: IsString.max(254).optional(),
	manufacturerId: StringAndBeEmptyOrOptional,
	manufacturerName: IsString.max(254).optional(),
	modelId: StringAndBeEmptyOrOptional,
	modelName: IsString.max(254).optional(),
	version: IsString.max(100).optional(),

	// Template visibility flags
	isSystem: IsBoolean.optional(),
	isTemplate: IsBoolean.optional(),

	assetIdPath: StringAndBeEmptyOrOptional,

	scriptTest: StringAndBeEmptyOrOptional,
	scriptProcessor: StringAndBeEmptyOrOptional,
	scriptValidator: StringAndBeEmptyOrOptional,
	scriptConversion: StringAndBeEmptyOrOptional,

	// Available Fields (for Rule autocomplete support)
	availableFields: z.array(IsString).optional(),

	// Dynamic Fields (for typed event storage and querying)
	dynamicFields: z.array(ZodDynamicFieldSchema).optional(),

	created: StringAndBeEmptyOrOptional,
	updated: StringAndBeEmptyOrOptional,
});

/**
 * Asset Template Query schema - Used for filtering and pagination
 */
export const ZodAssetTemplateQuerySchema = z.object({
	projection: StringAndBeEmptyOrOptional,
	page: NumberIntAndPositive.optional(),
	perPage: NumberIntAndPositive.max(100).optional(),
	sort: StringAndBeEmptyOrOptional,
	includeChildren: IsBoolean.optional(),

	// Module-specific filters
	name: StringAndNotBeEmpty.max(150).optional(),
	enabled: IsBoolean.optional(),
	isSystem: IsBoolean.optional(),
	isTemplate: IsBoolean.optional(),
	categoryId: StringAndBeEmptyOrOptional,
	manufacturerId: StringAndBeEmptyOrOptional,
	modelId: StringAndBeEmptyOrOptional,
});

/**
 * Asset Template Response schema - Represents API response structure
 */
export const ZodAssetTemplateResponseSchema = z.object({
	id: StringAndBeEmptyOrOptional,
	name: StringAndBeEmptyOrOptional,
	enabled: IsBoolean.optional(),
	description: StringAndBeEmptyOrOptional,

	// Asset Classification (FLAT structure)
	categoryId: StringAndBeEmptyOrOptional,
	categoryName: StringAndBeEmptyOrOptional,
	manufacturerId: StringAndBeEmptyOrOptional,
	manufacturerName: StringAndBeEmptyOrOptional,
	modelId: StringAndBeEmptyOrOptional,
	modelName: StringAndBeEmptyOrOptional,
	version: StringAndBeEmptyOrOptional,

	assetIdPath: StringAndBeEmptyOrOptional,
	orgId: StringAndBeEmptyOrOptional,

	// Template visibility flags
	isSystem: IsBoolean.optional(),
	isTemplate: IsBoolean.optional(),

	scriptTest: StringAndBeEmptyOrOptional,
	scriptProcessor: StringAndBeEmptyOrOptional,
	scriptValidator: StringAndBeEmptyOrOptional,
	scriptConversion: StringAndBeEmptyOrOptional,

	// Available Fields (for Rule autocomplete support)
	availableFields: z.array(IsString).optional(),

	// Dynamic Fields (for typed event storage and querying)
	dynamicFields: z.array(ZodDynamicFieldSchema).optional(),

	created: StringAndBeEmptyOrOptional,
	updated: StringAndBeEmptyOrOptional,
});

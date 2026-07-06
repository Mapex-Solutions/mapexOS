import { z, StringAndNotBeEmpty, StringAndBeEmptyOrOptional, IsString, IsNumber, IsMongoId, NumberIntAndPositive } from '@mapexos/validations';

/**
 * Template migration wire contracts.
 * Mirror of Go: packages/contracts/services/assets/assets_templates/migrations_dto.go
 * (MigrationPlanCreateRequest, MigrationPlanUpdateRequest, MigrationPlanResponse,
 * MigrationExecutionResponse).
 */

/**
 * Migration Plan Create schema - creates a plan migrating the selected assets
 * from their source template to the target template.
 */
export const ZodMigrationPlanCreateSchema = z.object({
	name: StringAndNotBeEmpty,
	description: StringAndBeEmptyOrOptional,
	fromTemplateId: IsMongoId,
	toTemplateId: IsMongoId,
	assetIds: z.array(IsMongoId).min(1, 'At least one asset is required'),
	scheduleAt: z.iso.datetime().nullable().optional(),
	batchSize: NumberIntAndPositive.optional(),
});

/**
 * Migration Plan Update schema - partial-update body (all fields optional).
 */
export const ZodMigrationPlanUpdateSchema = z.object({
	name: StringAndBeEmptyOrOptional,
	description: StringAndBeEmptyOrOptional,
	scheduleAt: z.iso.datetime().nullable().optional(),
});

/**
 * Migration Plan status enum.
 */
export const ZodMigrationPlanStatusSchema = z.enum([
	'pending',
	'scheduled',
	'running',
	'complete',
	'completed_with_errors',
	'cancelled',
]);

/**
 * Migration Plan Response schema - read model for a plan in list/detail responses.
 */
export const ZodMigrationPlanResponseSchema = z.object({
	id: IsString,
	orgId: IsString,
	name: IsString,
	description: IsString,
	fromTemplateId: IsString,
	fromTemplateName: IsString,
	toTemplateId: IsString,
	toTemplateName: IsString,
	scheduleAt: IsString,
	batchSize: IsNumber,
	status: ZodMigrationPlanStatusSchema,
	total: IsNumber,
	migrated: IsNumber,
	failed: IsNumber,
	progressPct: IsNumber,
	created: IsString,
	updated: IsString,
});

/**
 * Migration Plan Query schema - list filtering and pagination params.
 */
export const ZodMigrationPlanQuerySchema = z.object({
	page: NumberIntAndPositive.optional(),
	perPage: NumberIntAndPositive.max(100).optional(),
	sort: StringAndBeEmptyOrOptional,
	name: StringAndBeEmptyOrOptional,
	status: ZodMigrationPlanStatusSchema.optional(),
});

/**
 * Migration Plan Params schema - path params for single-plan endpoints.
 */
export const ZodMigrationPlanParamsSchema = z.object({
	id: IsMongoId,
});

/**
 * Migration Execution status enum.
 */
export const ZodMigrationExecutionStatusSchema = z.enum([
	'pending',
	'migrated',
	'failed',
	'cancelled',
]);

/**
 * Migration Execution Query schema - per-plan execution filtering and pagination.
 */
export const ZodMigrationExecutionQuerySchema = z.object({
	page: NumberIntAndPositive.optional(),
	perPage: NumberIntAndPositive.max(100).optional(),
	status: ZodMigrationExecutionStatusSchema.optional(),
	assetId: StringAndBeEmptyOrOptional,
});

/**
 * Migration Execution Response schema - read model for a single per-asset
 * migration execution.
 */
export const ZodMigrationExecutionResponseSchema = z.object({
	id: IsString,
	planId: IsString,
	assetId: IsString,
	orgId: IsString,
	status: ZodMigrationExecutionStatusSchema,
	error: StringAndBeEmptyOrOptional,
	attempts: IsNumber,
	created: IsString,
	updated: IsString,
});

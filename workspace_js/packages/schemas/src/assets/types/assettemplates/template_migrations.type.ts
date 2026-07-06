import { z } from 'zod';
import {
	ZodMigrationPlanCreateSchema,
	ZodMigrationPlanUpdateSchema,
	ZodMigrationPlanStatusSchema,
	ZodMigrationPlanQuerySchema,
	ZodMigrationPlanParamsSchema,
	ZodMigrationPlanResponseSchema,
	ZodMigrationExecutionStatusSchema,
	ZodMigrationExecutionQuerySchema,
	ZodMigrationExecutionResponseSchema,
} from '../../schemas/assettemplates';

/**
 * TypeScript types inferred from Zod schemas for Template Migrations
 */

export type MigrationPlanCreateRequest = z.infer<typeof ZodMigrationPlanCreateSchema>;
export type MigrationPlanUpdateRequest = z.infer<typeof ZodMigrationPlanUpdateSchema>;
export type MigrationPlanStatus = z.infer<typeof ZodMigrationPlanStatusSchema>;
export type MigrationPlanQuery = z.infer<typeof ZodMigrationPlanQuerySchema>;
export type MigrationPlanParams = z.infer<typeof ZodMigrationPlanParamsSchema>;
export type MigrationPlanResponse = z.infer<typeof ZodMigrationPlanResponseSchema>;
export type MigrationExecutionStatus = z.infer<typeof ZodMigrationExecutionStatusSchema>;
export type MigrationExecutionQuery = z.infer<typeof ZodMigrationExecutionQuerySchema>;
export type MigrationExecutionResponse = z.infer<typeof ZodMigrationExecutionResponseSchema>;

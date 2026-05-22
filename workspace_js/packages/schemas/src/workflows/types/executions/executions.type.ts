import { z } from 'zod';
import {
	ZodPathEntrySchema,
	ZodErrorInfoSchema,
	ZodExecutionIdSchema,
	ZodExecutionQuerySchema,
	ZodSignalRequestSchema,
	ZodExecutionResponseSchema,
} from '@/workflows/schemas/executions/executions.schema';
import { ExecutionStatusEnum } from '@/workflows/enums';

// Supporting types
export type PathEntry = z.infer<typeof ZodPathEntrySchema>;
export type ErrorInfo = z.infer<typeof ZodErrorInfoSchema>;

// DTO types
export type ExecutionId = z.infer<typeof ZodExecutionIdSchema>;
export type ExecutionQuery = z.infer<typeof ZodExecutionQuerySchema>;
export type SignalRequest = z.infer<typeof ZodSignalRequestSchema>;
export type ExecutionResponse = z.infer<typeof ZodExecutionResponseSchema>;

// Re-export enum as type
export type ExecutionStatus = ExecutionStatusEnum;

// Re-export enum object
export { ExecutionStatusEnum };

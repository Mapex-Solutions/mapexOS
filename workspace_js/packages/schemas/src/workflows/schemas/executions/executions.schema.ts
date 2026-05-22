import { z, StringAndNotBeEmpty, StringAndNotBeEmptyOrOptional, IsMongoId, IsString, IsNumber, NumberIntAndPositive } from '@mapexos/validations';

import { ExecutionStatusEnum } from '@/workflows/enums';

/**
 * ============================================================================
 * SUPPORTING SCHEMAS
 * ============================================================================
 */

/**
 * PathEntry schema - Single step in the execution path
 */
export const ZodPathEntrySchema = z.object({
	nodeId: StringAndNotBeEmpty,
	nodeType: StringAndNotBeEmpty,
	status: StringAndNotBeEmpty,
	enteredAt: StringAndNotBeEmptyOrOptional,
	exitedAt: StringAndNotBeEmptyOrOptional,
	durationMs: IsNumber.int().min(0).optional().default(0),
	outputHandle: IsString.optional().default(''),
	error: IsString.optional(),
});

/**
 * ErrorInfo schema - Error details when execution fails
 */
export const ZodErrorInfoSchema = z.object({
	code: StringAndNotBeEmpty,
	message: StringAndNotBeEmpty,
	nodeId: StringAndNotBeEmpty,
	nodeType: StringAndNotBeEmpty,
	timestamp: StringAndNotBeEmptyOrOptional,
	stackTrace: IsString.optional().default(''),
});

/**
 * ============================================================================
 * EXECUTION DTOs
 * ============================================================================
 */

/**
 * Execution ID parameter schema (for URL params - UUID string)
 */
export const ZodExecutionIdSchema = z.object({
	executionId: StringAndNotBeEmpty,
});

/**
 * Execution Query schema - Used for filtering and pagination
 */
export const ZodExecutionQuerySchema = z.object({
	instanceId: IsMongoId.optional(),
	definitionId: IsMongoId.optional(),
	status: IsString.optional(), // Accepts single or comma-separated (e.g., "running,waiting")
	page: NumberIntAndPositive.optional(),
	perPage: NumberIntAndPositive.max(100).optional(),
});

/**
 * Signal Request schema - Payload for sending a signal to a waiting execution
 */
export const ZodSignalRequestSchema = z.object({
	signalName: StringAndNotBeEmpty,
	data: z.record(IsString, z.any()).optional(),
});

/**
 * Execution Response schema - Used for API responses
 */
export const ZodExecutionResponseSchema = z.object({
	_id: IsString.optional(),
	workflowUUID: IsString.optional(),
	instanceId: IsMongoId.optional(),
	definitionId: IsMongoId.optional(),
	workflowName: StringAndNotBeEmptyOrOptional,
	orgId: IsMongoId.optional(),
	pathKey: IsString.optional(),
	eventTrackerId: IsString.optional(),
	triggerSource: IsString.optional(),
	version: IsNumber.int().optional(),
	status: z.nativeEnum(ExecutionStatusEnum).optional(),
	activeNodeIds: z.array(IsString).optional(),
	state: z.record(IsString, z.any()).optional(),
	eventPayload: z.record(IsString, z.any()).optional(),
	executionPath: z.array(ZodPathEntrySchema).optional(),
	nodeOutputs: z.record(IsString, z.any()).optional(),
	errorInfo: ZodErrorInfoSchema.optional(),
	parentExecutionId: IsString.optional(),
	depth: IsNumber.int().optional(),
	retentionDays: IsNumber.int().optional(),
	startedAt: StringAndNotBeEmptyOrOptional,
	completedAt: StringAndNotBeEmptyOrOptional,
	created: StringAndNotBeEmptyOrOptional,
	updated: StringAndNotBeEmptyOrOptional,
});

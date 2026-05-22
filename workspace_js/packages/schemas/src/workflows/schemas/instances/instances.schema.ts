import { z, StringAndNotBeEmpty, StringAndNotBeEmptyOrOptional, IsBoolean, IsMongoId, IsString, IsNumber, NumberIntAndPositive } from '@mapexos/validations';

/**
 * ============================================================================
 * INSTANCE CONFIG DTOs
 * ============================================================================
 */

/**
 * Instance ID parameter schema (for URL params - MongoDB ObjectID)
 */
export const ZodInstanceIdSchema = z.object({
	instanceId: IsMongoId,
});

/**
 * Instance Create schema - Body for creating an instance config
 */
export const ZodInstanceCreateSchema = z.object({
	definitionId: IsMongoId,
	definitionVersion: IsNumber.int().min(1),
	definitionName: IsString.optional().default(''),
	name: StringAndNotBeEmpty,
	description: IsString.optional().default(''),
	pathKey: IsString.optional().default(''),
	externalInputs: z.record(IsString, z.any()).optional(),
	isTemplate: IsBoolean.optional().default(false),
	uniqueExecution: IsBoolean.optional().default(false),
	workflowUUID: IsString.optional().default(''),
});

/**
 * Instance Update schema - Body for updating an instance config
 */
export const ZodInstanceUpdateSchema = z.object({
	name: StringAndNotBeEmpty.optional(),
	description: IsString.optional(),
	externalInputs: z.record(IsString, z.any()).optional(),
	isTemplate: IsBoolean.optional(),
	uniqueExecution: IsBoolean.optional(),
	workflowUUID: IsString.optional(),
	enabled: IsBoolean.optional(),
});

/**
 * Instance Query schema - Used for filtering and pagination
 */
export const ZodInstanceQuerySchema = z.object({
	definitionId: IsMongoId.optional(),
	name: IsString.max(100).optional(),
	enabled: IsBoolean.optional(),
	uniqueExecution: IsBoolean.optional(),
	projection: IsString.optional(),
	page: NumberIntAndPositive.optional(),
	perPage: NumberIntAndPositive.max(100).optional(),
});

/**
 * Execute Request schema - Optional body for executing a workflow instance
 */
export const ZodExecuteRequestSchema = z.object({
	eventPayload: z.record(IsString, z.any()).optional(),
	workflowUUID: IsString.optional(),
});

/**
 * Execute Response schema - Response from executing a workflow instance
 */
export const ZodExecuteResponseSchema = z.object({
	workflowUUID: IsString,
	status: IsString,
	errorInfo: z.object({
		code: IsString,
		message: IsString,
		nodeId: IsString.optional(),
		nodeType: IsString.optional(),
	}).optional().nullable(),
});

export interface ExecuteRequest {
	eventPayload?: Record<string, any>;
	workflowUUID?: string;
}

export interface ExecuteResponse {
	workflowUUID: string;
	status: string;
	errorInfo?: {
		code: string;
		message: string;
		nodeId?: string;
		nodeType?: string;
	} | null;
}

/**
 * Instance Response schema - API response for instance config
 */
export const ZodInstanceResponseSchema = z.object({
	_id: IsMongoId.optional(),
	definitionId: IsMongoId.optional(),
	definitionVersion: IsNumber.int().optional(),
	definitionName: StringAndNotBeEmptyOrOptional,
	name: StringAndNotBeEmptyOrOptional,
	description: StringAndNotBeEmptyOrOptional,
	orgId: IsMongoId.optional(),
	pathKey: StringAndNotBeEmptyOrOptional,
	externalInputs: z.record(IsString, z.any()).optional(),
	isSystem: IsBoolean.optional(),
	isTemplate: IsBoolean.optional(),
	uniqueExecution: IsBoolean.optional(),
	workflowUUID: StringAndNotBeEmptyOrOptional,
	enabled: IsBoolean.optional(),
	created: StringAndNotBeEmptyOrOptional,
	updated: StringAndNotBeEmptyOrOptional,
});

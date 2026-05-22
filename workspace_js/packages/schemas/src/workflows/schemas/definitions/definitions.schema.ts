import { z, StringAndNotBeEmpty, StringAndNotBeEmptyOrOptional, IsBoolean, IsMongoId, IsString, IsNumber, NumberIntAndPositive } from '@mapexos/validations';

import {
	FieldValueTypeEnum,
	VariableTypeEnum,
	GroupLogicOperatorEnum,
	DefinitionStatusEnum,
} from '@/workflows/enums';

/**
 * ============================================================================
 * BUILDING BLOCK SCHEMAS
 * ============================================================================
 */

/**
 * Position schema - Canvas coordinates for a node
 */
export const ZodPositionSchema = z.object({
	x: IsNumber,
	y: IsNumber,
});

/**
 * FieldValue schema - Dynamic reference to event, state, variable, literal, node output, or engine value
 */
export const ZodFieldValueSchema = z.object({
	type: z.nativeEnum(FieldValueTypeEnum),
	value: IsString,
	mode: IsString.optional(),
	nodeId: IsString.optional(),
});

/**
 * WorkflowNode schema - A single node in the workflow graph
 */
/**
 * NodeTimeoutConfig schema - Async timeout configuration at node level
 */
export const ZodNodeTimeoutSchema = z.object({
	duration: NumberIntAndPositive,
	unit: IsString, // "seconds" | "minutes" | "hours" | "days"
	enableOutput: IsBoolean.optional().default(false),
});

/**
 * NodeErrorHandlerConfig schema - Retry policy at node level
 */
export const ZodNodeErrorHandlerSchema = z.object({
	enabled: IsBoolean,
	maxAttempts: IsNumber.int().min(1).max(10),
	initialInterval: IsNumber.int().min(1),
	intervalUnit: IsString,
	backoffMultiplier: IsNumber.min(1).max(10),
});

export const ZodWorkflowNodeSchema = z.object({
	id: StringAndNotBeEmpty,
	type: StringAndNotBeEmpty,
	label: IsString.optional().default(''),
	position: ZodPositionSchema,
	config: z.record(IsString, z.any()).optional().default({}),
	timeout: ZodNodeTimeoutSchema.optional(),
	errorHandler: ZodNodeErrorHandlerSchema.optional(),
	parentNodeId: IsString.optional().default(''),
});

/**
 * WorkflowEdge schema - A connection between two nodes
 */
export const ZodWorkflowEdgeSchema = z.object({
	id: StringAndNotBeEmpty,
	source: StringAndNotBeEmpty,
	sourceHandle: IsString.optional().default(''),
	target: StringAndNotBeEmpty,
	targetHandle: IsString.optional().default(''),
	label: IsString.optional().default(''),
	pathOffsetX: IsNumber.optional().default(0),
	pathOffsetY: IsNumber.optional().default(0),
});

/**
 * WorkflowVariable schema - User-defined variable with type and default
 */
export const ZodWorkflowVariableSchema = z.object({
	field: StringAndNotBeEmpty,
	type: z.nativeEnum(VariableTypeEnum),
	defaultValue: z.any().optional(),
	description: IsString.optional().default(''),
	durable: IsBoolean.optional().default(false),
});

/**
 * CaptureField schema - Event field to capture into workflow state
 */
export const ZodCaptureFieldSchema = z.object({
	field: StringAndNotBeEmpty,
	type: z.nativeEnum(VariableTypeEnum),
	description: IsString.optional().default(''),
});

/**
 * ExternalSignal schema - Named signal that this workflow can wait for
 */
export const ZodExternalSignalSchema = z.object({
	name: StringAndNotBeEmpty,
	description: IsString.optional().default(''),
});

/**
 * RetryPolicy schema - Global retry configuration for the workflow
 */
export const ZodRetryPolicySchema = z.object({
	enabled: IsBoolean.optional().default(false),
	maxAttempts: IsNumber.int().min(1).max(10).optional().default(3),
	initialInterval: IsString.optional().default('1s'),
	backoffMultiplier: IsNumber.min(1).max(10).optional().default(2),
	maxInterval: IsString.optional().default('60s'),
	nonRetryableErrors: z.array(IsString).optional().default([]),
});

/**
 * CanvasViewport schema - Saved viewport position and zoom level
 */
export const ZodCanvasViewportSchema = z.object({
	x: IsNumber.optional().default(0),
	y: IsNumber.optional().default(0),
	zoom: IsNumber.min(0.1).max(4).optional().default(1),
});

/**
 * DefinitionMetadata schema - UI metadata for the workflow editor
 */
export const ZodDefinitionMetadataSchema = z.object({
	canvasViewport: ZodCanvasViewportSchema.optional(),
});

/**
 * ============================================================================
 * DEFINITION DTOs
 * ============================================================================
 */

/**
 * Definition ID parameter schema (for URL params - MongoDB ObjectID)
 */
export const ZodDefinitionIdSchema = z.object({
	workflowId: IsMongoId,
});

/**
 * Definition Create schema - Used for creating new workflow definitions
 */
export const ZodDefinitionCreateSchema = z.object({
	name: IsString.min(1).max(255),
	description: IsString.max(1000).optional().default(''),
	enabled: IsBoolean.optional().default(false),
	isTemplate: IsBoolean.optional().default(false),
	timezone: ZodFieldValueSchema.optional(),
	retryPolicy: ZodRetryPolicySchema.optional(),
	states: z.array(ZodWorkflowVariableSchema).optional().default([]),
	captureFields: z.array(ZodCaptureFieldSchema).optional().default([]),
	externalInputs: z.array(z.object({
		field: IsString,
		label: IsString,
		icon: IsString.optional().default(''),
		type: IsString,
		description: IsString.optional().default(''),
		defaultValue: z.any().optional(),
		required: IsBoolean.optional().default(false),
		assetTemplateId: IsString.optional().default(''),
		fieldPath: IsString.optional().default(''),
	})).optional().default([]),
	externalSignals: z.array(ZodExternalSignalSchema).optional().default([]),
	nodes: z.array(ZodWorkflowNodeSchema).min(1),
	edges: z.array(ZodWorkflowEdgeSchema).optional().default([]),
	installedPlugins: z.array(IsString).optional().default([]),
	metadata: ZodDefinitionMetadataSchema.optional(),
	scope: IsString.optional(),
});

/**
 * Definition Update schema - Used for updating existing definitions (all fields optional)
 */
export const ZodDefinitionUpdateSchema = z.object({
	name: IsString.min(1).max(255).optional(),
	description: IsString.max(1000).optional(),
	enabled: IsBoolean.optional(),
	isTemplate: IsBoolean.optional(),
	timezone: ZodFieldValueSchema.optional(),
	retryPolicy: ZodRetryPolicySchema.optional(),
	states: z.array(ZodWorkflowVariableSchema).optional(),
	captureFields: z.array(ZodCaptureFieldSchema).optional(),
	externalInputs: z.array(z.object({
		field: IsString,
		label: IsString,
		icon: IsString.optional().default(''),
		type: IsString,
		description: IsString.optional().default(''),
		defaultValue: z.any().optional(),
		required: IsBoolean.optional().default(false),
		assetTemplateId: IsString.optional().default(''),
		fieldPath: IsString.optional().default(''),
	})).optional(),
	externalSignals: z.array(ZodExternalSignalSchema).optional(),
	nodes: z.array(ZodWorkflowNodeSchema).optional(),
	edges: z.array(ZodWorkflowEdgeSchema).optional(),
	installedPlugins: z.array(IsString).optional(),
	metadata: ZodDefinitionMetadataSchema.optional(),
	scope: IsString.optional(),
});

/**
 * Definition Query schema - Used for filtering and pagination
 */
export const ZodDefinitionQuerySchema = z.object({
	name: IsString.max(100).optional(),
	enabled: IsBoolean.optional(),
	status: z.nativeEnum(DefinitionStatusEnum).optional(),
	isTemplate: IsBoolean.optional(),
	definitionVersion: IsNumber.int().min(1).optional(),
	page: NumberIntAndPositive.optional(),
	perPage: NumberIntAndPositive.max(100).optional(),
	projection: IsString.optional(),
});

/**
 * Definition Response schema - Used for API responses
 */
export const ZodDefinitionResponseSchema = z.object({
	_id: IsMongoId.optional(),
	orgId: IsMongoId.optional(),
	name: StringAndNotBeEmptyOrOptional,
	description: StringAndNotBeEmptyOrOptional,
	enabled: IsBoolean.optional(),
	isTemplate: IsBoolean.optional(),
	definitionVersion: IsNumber.int().optional(),
	timezone: ZodFieldValueSchema.optional(),
	retryPolicy: ZodRetryPolicySchema.optional(),
	states: z.array(ZodWorkflowVariableSchema).optional(),
	captureFields: z.array(ZodCaptureFieldSchema).optional(),
	externalInputs: z.array(z.object({
		field: IsString,
		label: IsString,
		icon: IsString.optional().default(''),
		type: IsString,
		description: IsString.optional().default(''),
		defaultValue: z.any().optional(),
		required: IsBoolean.optional().default(false),
		assetTemplateId: IsString.optional().default(''),
		fieldPath: IsString.optional().default(''),
	})).optional(),
	externalSignals: z.array(ZodExternalSignalSchema).optional(),
	nodes: z.array(ZodWorkflowNodeSchema).optional(),
	edges: z.array(ZodWorkflowEdgeSchema).optional(),
	installedPlugins: z.array(IsString).optional(),
	missingPlugins: z.array(IsString).optional(),
	status: z.nativeEnum(DefinitionStatusEnum).optional(),
	metadata: ZodDefinitionMetadataSchema.optional(),
	pathKey: StringAndNotBeEmptyOrOptional,
	scope: StringAndNotBeEmptyOrOptional,
	created: StringAndNotBeEmptyOrOptional,
	updated: StringAndNotBeEmptyOrOptional,
});

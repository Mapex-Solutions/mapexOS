import { z } from 'zod';
import {
	ZodPositionSchema,
	ZodFieldValueSchema,
	ZodNodeErrorHandlerSchema,
	ZodWorkflowNodeSchema,
	ZodWorkflowEdgeSchema,
	ZodWorkflowVariableSchema,
	ZodCaptureFieldSchema,
	ZodExternalSignalSchema,
	ZodRetryPolicySchema,
	ZodCanvasViewportSchema,
	ZodDefinitionMetadataSchema,
	ZodDefinitionIdSchema,
	ZodDefinitionCreateSchema,
	ZodDefinitionUpdateSchema,
	ZodDefinitionQuerySchema,
	ZodDefinitionResponseSchema,
} from '@/workflows/schemas/definitions/definitions.schema';
import {
	FieldValueTypeEnum,
	VariableTypeEnum,
	GroupLogicOperatorEnum,
} from '@/workflows/enums';

// Building block types
export type Position = z.infer<typeof ZodPositionSchema>;
export type FieldValue = z.infer<typeof ZodFieldValueSchema>;
export type WorkflowNode = z.infer<typeof ZodWorkflowNodeSchema>;
export type WorkflowEdge = z.infer<typeof ZodWorkflowEdgeSchema>;
export type WorkflowVariable = z.infer<typeof ZodWorkflowVariableSchema>;
export type CaptureField = z.infer<typeof ZodCaptureFieldSchema>;
export type ExternalSignal = z.infer<typeof ZodExternalSignalSchema>;
export type RetryPolicy = z.infer<typeof ZodRetryPolicySchema>;
export type NodeErrorHandler = z.infer<typeof ZodNodeErrorHandlerSchema>;
export type CanvasViewport = z.infer<typeof ZodCanvasViewportSchema>;
export type DefinitionMetadata = z.infer<typeof ZodDefinitionMetadataSchema>;

// DTO types
export type DefinitionId = z.infer<typeof ZodDefinitionIdSchema>;
export type DefinitionCreate = z.infer<typeof ZodDefinitionCreateSchema>;
export type DefinitionUpdate = z.infer<typeof ZodDefinitionUpdateSchema>;
export type DefinitionQuery = z.infer<typeof ZodDefinitionQuerySchema>;
export type DefinitionResponse = z.infer<typeof ZodDefinitionResponseSchema>;

// Re-export enums as types
export type FieldValueType = FieldValueTypeEnum;
export type VariableType = VariableTypeEnum;
export type GroupLogicOperator = GroupLogicOperatorEnum;

// Re-export enum objects
export { FieldValueTypeEnum, VariableTypeEnum, GroupLogicOperatorEnum };

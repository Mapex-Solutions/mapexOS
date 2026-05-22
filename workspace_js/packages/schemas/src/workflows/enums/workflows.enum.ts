/**
 * Enum for field value source types in workflow conditions and configs
 */
export enum FieldValueTypeEnum {
	EVENT = 'event',
	STATE = 'state',
	VARIABLE = 'variable',
	LITERAL = 'literal',
	NODE_OUTPUT = 'node_output',
	ENGINE = 'engine',
}

/**
 * Enum for workflow variable types
 */
export enum VariableTypeEnum {
	STRING = 'string',
	NUMBER = 'number',
	BOOLEAN = 'boolean',
	JSON = 'json',
}

/**
 * Enum for condition group logic operators
 */
export enum GroupLogicOperatorEnum {
	AND = 'AND',
	OR = 'OR',
	NAND = 'NAND',
	NOR = 'NOR',
}

/**
 * Enum for workflow definition status (computed by backend)
 */
export enum DefinitionStatusEnum {
	VALID = 'valid',
	PLUGIN_MISSING = 'plugin_missing',
	INVALID = 'invalid',
}

/**
 * Enum for workflow execution status
 */
export enum ExecutionStatusEnum {
	CREATED = 'created',
	RUNNING = 'running',
	WAITING = 'waiting',
	COMPLETED = 'completed',
	FAILED = 'failed',
	CANCELLED = 'cancelled',
}

/**
 * @deprecated Use ExecutionStatusEnum instead
 */
export const InstanceStatusEnum = ExecutionStatusEnum;
export type InstanceStatusEnum = ExecutionStatusEnum;

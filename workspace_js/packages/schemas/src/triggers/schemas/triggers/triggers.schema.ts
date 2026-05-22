import { z, StringAndNotBeEmpty, StringAndNotBeEmptyOrOptional, IsBoolean, IsMongoId, IsString, IsNumber, NumberIntAndPositive } from '@mapexos/validations';
import type { z as zType } from 'zod';
import {
	TriggerTypeEnum,
	TriggerCategoryEnum,
	HttpMethodEnum,
	MqttQosEnum,
	RabbitmqPublishModeEnum,
	RabbitmqExchangeTypeEnum,
} from '@/triggers/enums';

/**
 * ============================================================================
 * TECHNICAL TRIGGER CONFIG SCHEMAS
 * ============================================================================
 */

/**
 * HTTP Config schema - HTTP/HTTPS request configuration
 */
const ZodHttpConfigSchema = z.object({
	endpoint: IsString.url('Endpoint must be a valid URL'),
	method: z.nativeEnum(HttpMethodEnum),
	headers: z.record(IsString, IsString).optional(),
	body: z.record(IsString, z.any()).optional(),
	timeout: IsNumber.int().min(1000).max(300000).optional(), // 1s to 5min
});

/**
 * MQTT Config schema - MQTT broker configuration
 */
const ZodMqttConfigSchema = z.object({
	broker: StringAndNotBeEmpty,
	port: IsNumber.int().min(1).max(65535),
	topic: StringAndNotBeEmpty,
	qos: z.nativeEnum(MqttQosEnum),
	username: StringAndNotBeEmptyOrOptional,
	password: StringAndNotBeEmptyOrOptional,
	clientId: StringAndNotBeEmptyOrOptional,
	message: z.record(IsString, z.any()).optional(),
	useTLS: IsBoolean.optional(),
});

/**
 * RabbitMQ Config schema - RabbitMQ messaging configuration
 */
const ZodRabbitmqConfigSchema = z.object({
	host: StringAndNotBeEmpty,
	port: IsNumber.int().min(1).max(65535),
	vhost: StringAndNotBeEmptyOrOptional,
	username: StringAndNotBeEmpty,
	password: StringAndNotBeEmpty,
	publishMode: z.nativeEnum(RabbitmqPublishModeEnum),
	exchange: StringAndNotBeEmptyOrOptional,
	exchangeType: z.nativeEnum(RabbitmqExchangeTypeEnum).optional(),
	routingKey: StringAndNotBeEmptyOrOptional,
	queue: StringAndNotBeEmptyOrOptional,
	message: z.record(IsString, z.any()).optional(),
	useTLS: IsBoolean.optional(),
}).refine((data) => {
	// Validate exchange mode requirements
	if (data.publishMode === RabbitmqPublishModeEnum.EXCHANGE) {
		return !!data.exchange && !!data.exchangeType;
	}
	// Validate queue mode requirements
	if (data.publishMode === RabbitmqPublishModeEnum.QUEUE) {
		return !!data.queue;
	}
	return true;
}, {
	message: "Exchange and exchangeType are required when publishMode is 'exchange', or queue is required when publishMode is 'queue'",
	path: ['publishMode'],
});

/**
 * NATS Config schema - NATS messaging configuration
 */
const ZodNatsConfigSchema = z.object({
	server: StringAndNotBeEmpty,
	subject: StringAndNotBeEmpty,
	username: StringAndNotBeEmptyOrOptional,
	password: StringAndNotBeEmptyOrOptional,
	token: StringAndNotBeEmptyOrOptional,
	message: z.record(IsString, z.any()).optional(),
	useTLS: IsBoolean.optional(),
});

/**
 * WebSocket Config schema - WebSocket message configuration
 */
const ZodWebsocketConfigSchema = z.object({
	url: IsString.url('URL must be a valid WebSocket URL'),
	message: z.record(IsString, z.any()).optional(),
	headers: z.record(IsString, IsString).optional(),
});

/**
 * ============================================================================
 * COMMUNICATION TRIGGER CONFIG SCHEMAS
 * ============================================================================
 */

/**
 * Email Config schema - Email notification configuration
 */
const ZodEmailConfigSchema = z.object({
	smtpHost: StringAndNotBeEmpty,
	smtpPort: IsNumber.int().min(1).max(65535),
	username: StringAndNotBeEmptyOrOptional,
	password: StringAndNotBeEmptyOrOptional,
	fromAddr: StringAndNotBeEmpty.min(3, 'fromAddr must be at least 3 characters'),
	to: StringAndNotBeEmpty.min(3, 'to must be at least 3 characters'),
	cc: StringAndNotBeEmptyOrOptional,
	bcc: StringAndNotBeEmptyOrOptional,
	subject: StringAndNotBeEmpty,
	body: StringAndNotBeEmptyOrOptional,
	htmlBody: StringAndNotBeEmptyOrOptional,
});

/**
 * Teams Config schema - Microsoft Teams webhook configuration
 */
const ZodTeamsConfigSchema = z.object({
	webhookUrl: IsString.url('Webhook URL must be a valid URL'),
	title: StringAndNotBeEmpty,
	text: StringAndNotBeEmpty,
	themeColor: IsString.regex(/^[0-9A-Fa-f]{6}$/, 'Theme color must be a valid hex color (without #)').optional(),
});

/**
 * Slack Config schema - Slack webhook configuration
 */
const ZodSlackConfigSchema = z.object({
	webhookUrl: IsString.url('Webhook URL must be a valid URL'),
	channel: StringAndNotBeEmptyOrOptional,
	username: StringAndNotBeEmptyOrOptional,
	iconEmoji: StringAndNotBeEmptyOrOptional,
	message: StringAndNotBeEmpty,
});

/**
 * ============================================================================
 * TRIGGER CONFIG (UNION TYPE)
 * ============================================================================
 */

/**
 * Trigger Config schema - Union type for all trigger configurations
 * Only ONE config field should be populated based on triggerType
 */
const ZodTriggerConfigSchema = z.object({
	// Technical triggers
	http: ZodHttpConfigSchema.optional(),
	mqtt: ZodMqttConfigSchema.optional(),
	rabbitmq: ZodRabbitmqConfigSchema.optional(),
	nats: ZodNatsConfigSchema.optional(),
	websocket: ZodWebsocketConfigSchema.optional(),

	// Communication triggers
	email: ZodEmailConfigSchema.optional(),
	teams: ZodTeamsConfigSchema.optional(),
	slack: ZodSlackConfigSchema.optional(),
});

/**
 * Helper function to validate config matches triggerType
 */
function validateConfigMatchesTriggerType(
	triggerType: string,
	config: zType.infer<typeof ZodTriggerConfigSchema>
): boolean {
	// Count populated fields
	const populatedFields = Object.keys(config).filter(key => config[key as keyof typeof config] !== undefined);

	// Must have exactly ONE field populated
	if (populatedFields.length !== 1) {
		return false;
	}

	// The populated field must match the triggerType
	const populatedField = populatedFields[0];
	return populatedField === triggerType;
}

/**
 * ============================================================================
 * TRIGGER DTOs
 * ============================================================================
 */

/**
 * Trigger ID parameter schema (for URL params - MongoDB ObjectID)
 */
export const ZodTriggerIdSchema = z.object({
	triggerId: IsMongoId,
});

/**
 * Trigger Create schema - Used for creating new triggers
 */
export const ZodTriggerCreateSchema = z.object({
	name: IsString.min(3).max(150),
	description: IsString.max(500).optional(),
	triggerType: z.nativeEnum(TriggerTypeEnum),
	category: z.nativeEnum(TriggerCategoryEnum),
	enabled: IsBoolean,

	// Template visibility flags
	isSystem: IsBoolean.optional().default(false),
	isTemplate: IsBoolean.optional().default(false),

	config: ZodTriggerConfigSchema,
	orgId: IsMongoId.optional(),
	pathKey: StringAndNotBeEmptyOrOptional,
}).refine((data) => {
	// Validate that config matches triggerType
	return validateConfigMatchesTriggerType(data.triggerType, data.config);
}, {
	message: "Config must have exactly one field populated that matches the triggerType",
	path: ['config'],
});

/**
 * Trigger Update schema - Used for updating existing triggers (all fields optional)
 */
export const ZodTriggerUpdateSchema = z.object({
	name: IsString.min(3).max(150).optional(),
	description: IsString.max(500).optional(),
	triggerType: z.nativeEnum(TriggerTypeEnum).optional(),
	category: z.nativeEnum(TriggerCategoryEnum).optional(),
	enabled: IsBoolean.optional(),

	// Template visibility flags (optional for updates)
	isSystem: IsBoolean.optional(),
	isTemplate: IsBoolean.optional(),

	config: ZodTriggerConfigSchema.optional(),
}).refine((data) => {
	// Only validate config if both triggerType and config are provided
	if (data.triggerType && data.config) {
		return validateConfigMatchesTriggerType(data.triggerType, data.config);
	}
	return true;
}, {
	message: "Config must match the triggerType when both are provided",
	path: ['config'],
});

/**
 * Trigger Query schema - Used for filtering and pagination
 */
export const ZodTriggerQuerySchema = z.object({
	id: IsMongoId.optional(),
	name: IsString.max(100).optional(),
	triggerType: z.nativeEnum(TriggerTypeEnum).optional(),
	category: z.nativeEnum(TriggerCategoryEnum).optional(),
	enabled: IsBoolean.optional(),
	orgId: IsMongoId.optional(),
	pathKey: StringAndNotBeEmptyOrOptional,

	// Template filters (for querying system/template triggers)
	isSystem: IsBoolean.optional(),
	isTemplate: IsBoolean.optional(),
	includeChildren: IsBoolean.optional(),

	// Pagination
	page: NumberIntAndPositive.optional(),
	pageSize: NumberIntAndPositive.max(100).optional(),
	sort: StringAndNotBeEmptyOrOptional,
});

/**
 * Trigger Response schema - Used for API responses
 */
export const ZodTriggerResponseSchema = z.object({
	id: IsMongoId.optional(),
	name: StringAndNotBeEmptyOrOptional,
	description: StringAndNotBeEmptyOrOptional,
	triggerType: z.nativeEnum(TriggerTypeEnum).optional(),
	category: z.nativeEnum(TriggerCategoryEnum).optional(),
	enabled: IsBoolean.optional(),

	// Template visibility flags
	isSystem: IsBoolean.optional(),
	isTemplate: IsBoolean.optional(),

	config: ZodTriggerConfigSchema.optional(),
	orgId: IsMongoId.optional(),
	pathKey: StringAndNotBeEmptyOrOptional,
	created: StringAndNotBeEmptyOrOptional,
	updated: StringAndNotBeEmptyOrOptional,
});

/**
 * Export config schemas for individual use if needed
 */
export const TriggerConfigSchemas = {
	http: ZodHttpConfigSchema,
	mqtt: ZodMqttConfigSchema,
	rabbitmq: ZodRabbitmqConfigSchema,
	nats: ZodNatsConfigSchema,
	websocket: ZodWebsocketConfigSchema,
	email: ZodEmailConfigSchema,
	teams: ZodTeamsConfigSchema,
	slack: ZodSlackConfigSchema,
} as const;

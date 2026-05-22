import { z } from 'zod';
import {
	ZodTriggerIdSchema,
	ZodTriggerCreateSchema,
	ZodTriggerUpdateSchema,
	ZodTriggerResponseSchema,
	ZodTriggerQuerySchema,
	TriggerConfigSchemas,
} from '@/triggers/schemas/triggers/triggers.schema';
import {
	TriggerTypeEnum,
	TriggerCategoryEnum,
	HttpMethodEnum,
	MqttQosEnum,
	RabbitmqPublishModeEnum,
	RabbitmqExchangeTypeEnum,
} from '@/triggers/enums';

// Export inferred types for main DTOs
export type TriggerId = z.infer<typeof ZodTriggerIdSchema>;
export type TriggerCreate = z.infer<typeof ZodTriggerCreateSchema>;
export type TriggerUpdate = z.infer<typeof ZodTriggerUpdateSchema>;
export type TriggerResponse = z.infer<typeof ZodTriggerResponseSchema>;
export type TriggerQuery = z.infer<typeof ZodTriggerQuerySchema>;

// Export inferred types for individual config schemas
export type HttpConfig = z.infer<typeof TriggerConfigSchemas.http>;
export type MqttConfig = z.infer<typeof TriggerConfigSchemas.mqtt>;
export type RabbitmqConfig = z.infer<typeof TriggerConfigSchemas.rabbitmq>;
export type NatsConfig = z.infer<typeof TriggerConfigSchemas.nats>;
export type WebsocketConfig = z.infer<typeof TriggerConfigSchemas.websocket>;
export type EmailConfig = z.infer<typeof TriggerConfigSchemas.email>;
export type TeamsConfig = z.infer<typeof TriggerConfigSchemas.teams>;
export type SlackConfig = z.infer<typeof TriggerConfigSchemas.slack>;

// TriggerConfig union type
export type TriggerConfig = {
	http?: HttpConfig;
	mqtt?: MqttConfig;
	rabbitmq?: RabbitmqConfig;
	nats?: NatsConfig;
	websocket?: WebsocketConfig;
	email?: EmailConfig;
	teams?: TeamsConfig;
	slack?: SlackConfig;
};

// Re-export enums as types
export type TriggerType = TriggerTypeEnum;
export type TriggerCategory = TriggerCategoryEnum;
export type HttpMethod = HttpMethodEnum;
export type MqttQos = MqttQosEnum;
export type RabbitmqPublishMode = RabbitmqPublishModeEnum;
export type RabbitmqExchangeType = RabbitmqExchangeTypeEnum;

// Re-export enum objects
export { TriggerTypeEnum, TriggerCategoryEnum, HttpMethodEnum, MqttQosEnum, RabbitmqPublishModeEnum, RabbitmqExchangeTypeEnum };

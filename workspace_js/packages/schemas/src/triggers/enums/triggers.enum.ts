/**
 * Enum for trigger types
 */
export enum TriggerTypeEnum {
	HTTP = 'http',
	MQTT = 'mqtt',
	RABBITMQ = 'rabbitmq',
	NATS = 'nats',
	WEBSOCKET = 'websocket',
	EMAIL = 'email',
	TEAMS = 'teams',
	SLACK = 'slack',
}

/**
 * Enum for trigger categories
 */
export enum TriggerCategoryEnum {
	TECHNICAL = 'technical',
	COMMUNICATION = 'communication',
}

/**
 * Enum for HTTP methods
 */
export enum HttpMethodEnum {
	GET = 'GET',
	POST = 'POST',
	PUT = 'PUT',
	PATCH = 'PATCH',
	DELETE = 'DELETE',
}

/**
 * Enum for MQTT QoS levels
 */
export enum MqttQosEnum {
	AT_MOST_ONCE = 0,
	AT_LEAST_ONCE = 1,
	EXACTLY_ONCE = 2,
}

/**
 * Enum for RabbitMQ publish modes
 */
export enum RabbitmqPublishModeEnum {
	EXCHANGE = 'exchange',
	QUEUE = 'queue',
}

/**
 * Enum for RabbitMQ exchange types
 */
export enum RabbitmqExchangeTypeEnum {
	DIRECT = 'direct',
	FANOUT = 'fanout',
	TOPIC = 'topic',
	HEADERS = 'headers',
}

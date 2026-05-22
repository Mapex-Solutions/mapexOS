import type { ConfigModule, Logger } from '@mapexos/microservices';

import { container } from 'tsyringe';

import {
	createRedisClient,
	RedisService as BaseRedisService,
	REDIS_CONNECTION_TOKEN,
	REDIS_SERVICE_TOKEN,
} from '@mapexos/infrastructure';

import { RedisService } from '@shared/services';

// InitRedis registers Redis client and services in DI container.
/** Initializes Redis client with lock configuration, registers in DI container. */
export function initRedis(configModule: ConfigModule, logger: Logger) {
	const redisConfig = {
		host: configModule.get('redis_host'),
		port: configModule.get('redis_port'),
		username: configModule.get('redis_username') || undefined,
		password: configModule.get('redis_password') || undefined,
		db: configModule.get('redis_db'),
		lock: {
			driftFactor: configModule.get('redis_lock_drift_factor'),
			retryCount: configModule.get('redis_lock_retry_count'),
			retryDelay: configModule.get('redis_lock_retry_delay'),
			retryJitter: configModule.get('redis_lock_retry_jitter'),
		}
	};

	const redisClient = createRedisClient(redisConfig);
	const baseRedisService = new BaseRedisService(redisClient, redisConfig);
	const redisService = new RedisService(logger, baseRedisService);

	container.register(REDIS_CONNECTION_TOKEN, { useValue: redisClient });
	container.register(REDIS_SERVICE_TOKEN, { useValue: baseRedisService });
	container.register('RedisService', { useValue: redisService });
}

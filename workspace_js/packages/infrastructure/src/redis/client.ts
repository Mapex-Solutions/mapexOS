import Redis from 'ioredis';
import { RedisConfig } from './types';
import { getInfraLogger } from '../logger';

export function createRedisClient(config: RedisConfig): Redis {
	const redis = new Redis({
		host: config.host,
		port: config.port,
		username: config.username || undefined,
		password: config.password || undefined,
		db: config.db || 0,
		maxRetriesPerRequest: 3,
		connectTimeout: 60000,
		commandTimeout: 5000,
		lazyConnect: true,
	});

	redis.on('connect', () => {
		getInfraLogger().info('[INFRA:REDIS] Connected');
	});

	redis.on('error', (err) => {
		getInfraLogger().error({ err }, '[INFRA:REDIS] Connection error');
	});

	redis.on('close', () => {
		getInfraLogger().info('[INFRA:REDIS] Connection closed');
	});

	redis.on('reconnecting', () => {
		getInfraLogger().info('[INFRA:REDIS] Reconnecting...');
	});

	return redis;
}
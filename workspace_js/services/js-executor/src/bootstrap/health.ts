import type { Application, ConfigModule } from '@mapexos/microservices';

import { container } from 'tsyringe';

import {
	RedisService as BaseRedisService,
	NatsBus,
	MinIOClient,
	REDIS_SERVICE_TOKEN,
	NATS_BUS_TOKEN,
} from '@mapexos/infrastructure';

import {
	registerHealthRoutes,
	RedisAdapter,
	NATSAdapter,
	MinIOAdapter,
} from '@mapexos/microservices';

// InitHealth registers the /health endpoint with all infrastructure checkers.
/** Registers health check endpoint with Redis, NATS, and MinIO checkers. */
export function initHealth(app: Application, configModule: ConfigModule) {
	const serviceName = configModule.get('service_name') as string;
	const serviceVersion = configModule.get('service_version') as string;

	const redisService = container.resolve<BaseRedisService>(REDIS_SERVICE_TOKEN);
	const natsBus = container.resolve<NatsBus>(NATS_BUS_TOKEN);
	const minioAssetsClient = container.resolve<MinIOClient>('MinIOAssetsClient');
	const minioTemplatesClient = container.resolve<MinIOClient>('MinIOTemplatesClient');
	const minioBytecodeClient = container.resolve<MinIOClient>('MinIOBytecodeClient');

	registerHealthRoutes(app, {
		serviceName,
		version: serviceVersion,
		cacheTTL: 10_000,
		timeout: 5_000,
	}, [
		{ checker: new RedisAdapter(redisService, 'app'), critical: true },
		{ checker: new NATSAdapter(natsBus, 'core'), critical: true },
		{ checker: new MinIOAdapter(minioAssetsClient, 'assets'), critical: false },
		{ checker: new MinIOAdapter(minioTemplatesClient, 'templates'), critical: false },
		{ checker: new MinIOAdapter(minioBytecodeClient, 'bytecode'), critical: false },
	]);
}

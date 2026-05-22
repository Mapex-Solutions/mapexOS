import type { Application, ConfigModule } from '@mapexos/microservices';

import { container } from 'tsyringe';

import {
	NatsBus,
	MinIOClient,
	NATS_BUS_TOKEN,
} from '@mapexos/infrastructure';

import {
	registerHealthRoutes,
	NATSAdapter,
	MinIOAdapter,
} from '@mapexos/microservices';

// InitHealth registers the /health endpoint with all infrastructure checkers.
/** Registers health check endpoint with NATS and MinIO infrastructure checkers. */
export function initHealth(app: Application, configModule: ConfigModule) {
	const serviceName = configModule.get('service_name') as string;
	const serviceVersion = configModule.get('service_version') as string;

	const natsBus = container.resolve<NatsBus>(NATS_BUS_TOKEN);
	const minioWorkflowsClient = container.resolve<MinIOClient>('MinIOWorkflowsClient');

	registerHealthRoutes(app, {
		serviceName,
		version: serviceVersion,
		cacheTTL: 10_000,
		timeout: 5_000,
	}, [
		{ checker: new NATSAdapter(natsBus, 'core'), critical: true },
		{ checker: new MinIOAdapter(minioWorkflowsClient, 'workflows'), critical: false },
	]);
}

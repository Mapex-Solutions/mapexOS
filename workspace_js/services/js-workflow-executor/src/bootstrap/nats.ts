import type { ConfigModule } from '@mapexos/microservices';

import { container } from 'tsyringe';

import {
	createNatsClient,
	NatsBus,
	NATS_CONNECTION_TOKEN,
	NATS_BUS_TOKEN,
} from '@mapexos/infrastructure';

// InitNATS registers NATS client and Bus in DI container.
/** Initializes NATS client and NatsBus, registers in DI container. */
export async function initNATS(configModule: ConfigModule) {
	const natsConfig = configModule.getNatsConfig();
	const natsClient = await createNatsClient(natsConfig);
	const natsBus = new NatsBus(natsClient);

	container.register(NATS_CONNECTION_TOKEN, { useValue: natsClient });
	container.register(NATS_BUS_TOKEN, { useValue: natsBus });
}

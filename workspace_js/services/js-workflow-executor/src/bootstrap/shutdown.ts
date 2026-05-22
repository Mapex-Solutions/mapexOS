import type { Logger } from '@mapexos/microservices';
import type { Server } from 'http';

import { container } from 'tsyringe';

import { ShutdownManager } from '@mapexos/microservices';
import { closeNatsClient, NATS_CONNECTION_TOKEN } from '@mapexos/infrastructure';

import type { NatsClient } from '@mapexos/infrastructure';

// InitShutdown creates a ShutdownManager and registers graceful shutdown hooks
// for all infrastructure components.
//
//   P0 — Express HTTP (stop accepting, drain in-flight requests)
//   P5 — Connections: NATS (close connection)
/** Registers graceful shutdown hooks for HTTP and NATS. */
export function initShutdown(logger: Logger, server: Server): ShutdownManager {
	const sm = new ShutdownManager(logger);

	// P0: HTTP — stop accepting new requests, drain in-flight
	sm.registerFunc('express', 0, async () => {
		await new Promise<void>((resolve, reject) => {
			server.close((err) => (err ? reject(err) : resolve()));
		});
	});

	// P5: NATS — close connection
	sm.registerFunc('nats', 5, async () => {
		const natsClient = container.resolve<NatsClient>(NATS_CONNECTION_TOKEN);
		await closeNatsClient(natsClient);
	});

	return sm;
}

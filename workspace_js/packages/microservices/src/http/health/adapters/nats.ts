/**
 * NATS Health Adapter
 * Same structure as workspace_go/packages/microservices/http/health/adapters/nats.go
 */

import type { HealthStatus, NatsBus } from '@mapexos/infrastructure';
import type { Checker } from '../types';

export class NATSAdapter implements Checker {
	constructor(
		private readonly bus: NatsBus,
		private readonly instanceName: string,
	) {}

	/** Returns the adapter identifier for this NATS instance. */
	name(): string {
		return `nats:${this.instanceName}`;
	}

	/**
	 * Pings NATS using the native PING/PONG protocol and returns the health status.
	 * @returns {Promise<HealthStatus>} The current health status with RTT latency
	 */
	async check(): Promise<HealthStatus> {
		try {
			const rtt = await this.bus.ping();
			return {
				service: `nats:${this.instanceName}`,
				connected: true,
				latencyMs: rtt,
				lastCheckAt: new Date(),
			};
		} catch (err) {
			return {
				service: `nats:${this.instanceName}`,
				connected: false,
				latencyMs: 0,
				lastCheckAt: new Date(),
				errorMessage: err instanceof Error ? err.message : String(err),
			};
		}
	}
}

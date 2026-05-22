/**
 * Redis Health Adapter
 * Same structure as workspace_go/packages/microservices/http/health/adapters/redis.go
 */

import type { HealthStatus, RedisService } from '@mapexos/infrastructure';
import type { Checker } from '../types';

export class RedisAdapter implements Checker {
	constructor(
		private readonly client: RedisService,
		private readonly instanceName: string,
	) {}

	/** Returns the adapter identifier for this Redis instance. */
	name(): string {
		return `redis:${this.instanceName}`;
	}

	/**
	 * Pings Redis and returns the health status with measured latency.
	 * @returns {Promise<HealthStatus>} The current health status
	 */
	async check(): Promise<HealthStatus> {
		const start = Date.now();
		try {
			await this.client.ping();
			const latency = Date.now() - start;
			return {
				service: `redis:${this.instanceName}`,
				connected: true,
				latencyMs: latency,
				lastCheckAt: new Date(),
			};
		} catch (err) {
			const latency = Date.now() - start;
			return {
				service: `redis:${this.instanceName}`,
				connected: false,
				latencyMs: latency,
				lastCheckAt: new Date(),
				errorMessage: err instanceof Error ? err.message : String(err),
			};
		}
	}
}

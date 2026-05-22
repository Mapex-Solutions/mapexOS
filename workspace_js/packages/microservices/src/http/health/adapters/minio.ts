/**
 * MinIO Health Adapter
 * Same structure as workspace_go/packages/microservices/http/health/adapters/minio.go
 */

import type { HealthStatus, MinIOClient } from '@mapexos/infrastructure';
import type { Checker } from '../types';

export class MinIOAdapter implements Checker {
	constructor(
		private readonly client: MinIOClient,
		private readonly instanceName: string,
	) {}

	/** Returns the adapter identifier for this MinIO instance. */
	name(): string {
		return `minio:${this.instanceName}`;
	}

	/**
	 * Pings MinIO via listBuckets and returns the health status with measured latency.
	 * @returns {Promise<HealthStatus>} The current health status
	 */
	async check(): Promise<HealthStatus> {
		const start = Date.now();
		try {
			await this.client.Ping();
			const latency = Date.now() - start;
			return {
				service: `minio:${this.instanceName}`,
				connected: true,
				latencyMs: latency,
				lastCheckAt: new Date(),
			};
		} catch (err) {
			const latency = Date.now() - start;
			return {
				service: `minio:${this.instanceName}`,
				connected: false,
				latencyMs: latency,
				lastCheckAt: new Date(),
				errorMessage: err instanceof Error ? err.message : String(err),
			};
		}
	}
}

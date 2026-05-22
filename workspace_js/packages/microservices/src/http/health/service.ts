/**
 * Health Service
 * Same structure as workspace_go/packages/microservices/http/health/service.go
 */

import type { HealthConfig, CheckerConfig, HealthResponse, CheckDetail } from './types';

/** Aggregates health checkers and caches results. */
export class HealthService {
	private readonly cfg: HealthConfig;
	private readonly checkers: CheckerConfig[];
	private readonly startedAt: number;

	private cachedResp: HealthResponse | null = null;
	private lastCheckAt = 0;

	constructor(cfg: HealthConfig, checkers: CheckerConfig[]) {
		this.cfg = {
			...cfg,
			cacheTTL: cfg.cacheTTL || 10_000,
			timeout: cfg.timeout || 5_000,
		};
		this.checkers = checkers;
		this.startedAt = Date.now();
	}

	/** Runs all health checkers (with cache) and returns the aggregated result. */
	async check(): Promise<HealthResponse> {
		const now = Date.now();

		if (this.cachedResp && (now - this.lastCheckAt) < this.cfg.cacheTTL) {
			return this.cachedResp;
		}

		const checks: Record<string, CheckDetail> = {};

		const results = await Promise.allSettled(
			this.checkers.map(async (cc) => {
				const timeoutPromise = new Promise<never>((_, reject) =>
					setTimeout(() => reject(new Error('health check timeout')), this.cfg.timeout)
				);

				try {
					const hs = await Promise.race([cc.checker.check(), timeoutPromise]);
					return { name: cc.checker.name(), critical: cc.critical, hs };
				} catch (err) {
					return {
						name: cc.checker.name(),
						critical: cc.critical,
						hs: {
							connected: false,
							service: cc.checker.name(),
							latencyMs: 0,
							lastCheckAt: new Date(),
							errorMessage: err instanceof Error ? err.message : String(err),
						},
					};
				}
			})
		);

		for (const result of results) {
			if (result.status === 'fulfilled') {
				const { name, critical, hs } = result.value;
				checks[name] = {
					connected: hs.connected,
					critical,
					latencyMs: hs.latencyMs,
					errorMessage: hs.errorMessage,
				};
			}
		}

		let status = 'healthy';
		for (const d of Object.values(checks)) {
			if (!d.connected && d.critical) {
				status = 'unhealthy';
				break;
			}
			if (!d.connected) {
				status = 'degraded';
			}
		}

		const uptimeMs = Date.now() - this.startedAt;
		const uptimeSec = Math.floor(uptimeMs / 1000);
		const hours = Math.floor(uptimeSec / 3600);
		const minutes = Math.floor((uptimeSec % 3600) / 60);
		const seconds = uptimeSec % 60;
		const uptime = `${hours}h${minutes}m${seconds}s`;

		const checkTime = new Date();
		const resp: HealthResponse = {
			status,
			service: this.cfg.serviceName,
			version: this.cfg.version,
			uptime,
			timestamp: checkTime,
			lastCheckAt: checkTime,
			checks,
		};

		this.cachedResp = resp;
		this.lastCheckAt = Date.now();

		return resp;
	}
}

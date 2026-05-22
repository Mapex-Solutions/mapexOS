/**
 * Health Check Types
 * Same structure as workspace_go/packages/microservices/http/health/types.go
 */

import type { HealthStatus } from '@mapexos/infrastructure';

/** Configuration for the health check service. */
export interface HealthConfig {
	serviceName: string;
	version: string;
	cacheTTL: number;   // milliseconds
	timeout: number;    // milliseconds
}

/** Pairs a health checker with its criticality level. */
export interface CheckerConfig {
	checker: Checker;
	critical: boolean;
}

/** Interface that infrastructure adapters must implement. */
export interface Checker {
	name(): string;
	check(): Promise<HealthStatus>;
}

/** JSON payload returned by the /health endpoint. */
export interface HealthResponse {
	status: string;     // "healthy", "degraded", "unhealthy"
	service: string;
	version: string;
	uptime: string;
	timestamp: Date;
	lastCheckAt: Date;
	checks: Record<string, CheckDetail>;
}

/** Status of an individual dependency in the health response. */
export interface CheckDetail {
	connected: boolean;
	critical: boolean;
	latencyMs?: number;
	errorMessage?: string;
}

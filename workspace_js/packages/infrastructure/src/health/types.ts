/**
 * HealthStatus represents the health state of an infrastructure dependency.
 * Same structure as workspace_go/packages/infrastructure/common/ports/health.ports.go
 */
export interface HealthStatus {
	connected: boolean;
	service: string;
	latencyMs?: number;
	lastCheckAt: Date;
	errorMessage?: string;
}

/**
 * Health Routes
 * Same structure as workspace_go/packages/microservices/http/health/routes.go
 */

import type { Application } from 'express';
import type { HealthConfig, CheckerConfig } from './types';
import { HealthService } from './service';
import { healthHandler } from './handler';

/** Registers the /health endpoint on the given Express app. */
export function registerHealthRoutes(
	app: Application,
	cfg: HealthConfig,
	checkers: CheckerConfig[],
): HealthService {
	const service = new HealthService(cfg, checkers);
	app.get('/health', healthHandler(service));
	return service;
}

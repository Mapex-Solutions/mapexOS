import type { Application } from '@mapexos/microservices';

import express from 'express';
import cors from 'cors';
import { container } from 'tsyringe';

import { APP_TOKEN } from '@mapexos/microservices';

import type { WorkflowExecutorMetrics } from './metrics';
import { METRICS_TOKEN } from '@shared/constants';

// InitExpress creates and registers the Express instance with global middlewares.
/** Initializes Express app with CORS, JSON middleware, and Prometheus metrics endpoint. */
export function initExpress(): Application {
	const app: Application = express();

	container.register(APP_TOKEN, { useValue: app });

	// Global middlewares
	app.use(cors());
	app.use(express.json());

	// Register Prometheus metrics endpoint (GET /metrics)
	const metrics = container.resolve<WorkflowExecutorMetrics>(METRICS_TOKEN);
	metrics.registry.registerEndpoint(app);

	return app;
}

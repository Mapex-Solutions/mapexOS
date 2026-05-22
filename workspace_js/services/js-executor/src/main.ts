import 'reflect-metadata';

import { initConfig, initMetrics, initNATS, initRedis, initCache, initExpress, initHealth, initShutdown } from './bootstrap';
import { initModule } from '@modules/app/module';

// Main function to start the JS Executor service
/** Application entry point — bootstraps all infrastructure and modules. */
async function bootstrap() {

	/**
	* Initialize configuration and logger
	 */
	const { configModule, logger } = initConfig();

	/**
	* Initialize Prometheus metrics registry (17 custom + Node.js defaults)
	 */
	initMetrics();

	/**
	* Initialize all infrastructure providers
	 */
	await initNATS(configModule);
	initRedis(configModule, logger);
	await initCache(configModule, logger);

	/**
	* Create Express instance with global middlewares
	 */
	const app = initExpress();

	/**
	* Initialize health check endpoint (before business modules, no auth required)
	 */
	initHealth(app, configModule);

	/**
	* Initialize all business modules (services, consumers, routes)
	 */
	await initModule();

	/**
	* Start the HTTP server (non-blocking)
	 */
	const httpPort = configModule.get('http_port');
	const server = app.listen(httpPort, () => {
		logger.info({ port: httpPort }, '[APP:BOOTSTRAP] Server started');
	});

	/*
	 * Create shutdown manager and register infrastructure hooks,
	 * then block until SIGTERM/SIGINT for graceful shutdown
	 */
	const sm = initShutdown(logger, server);
	sm.waitForSignal(15_000);
}

bootstrap();

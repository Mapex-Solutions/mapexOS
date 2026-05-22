import path from 'path';
import type { Application } from 'express';

import type { Logger } from '@shared/logger';
import type { MarketplaceServicePort } from './application/ports';
import { createMarketplaceService } from './application/di';
import { MARKETPLACE_PORT, PUBLIC_DIR_RELATIVE } from './application/constants';
import { registerMarketplaceRoutes } from './interfaces/http/routes';

/**
 * Bootstraps the Marketplace module: creates the service and wires its HTTP
 * interfaces into the Express application.
 *
 * @param app - Express application instance.
 * @param logger - Logger instance for module-level diagnostics.
 * @returns The initialized MarketplaceServicePort.
 */
export function initMarketplaceModule(app: Application, logger: Logger): MarketplaceServicePort {
	logger.info('[MODULE:Marketplace] Initializing');

	const publicDirectory = path.join(__dirname, '..', '..', PUBLIC_DIR_RELATIVE);

	const service = createMarketplaceService({
		logger,
		publicDirectory,
		port: MARKETPLACE_PORT,
	});

	registerMarketplaceRoutes(app, service);

	logger.info('[MODULE:Marketplace] Initialized');
	return service;
}

import type { Application } from 'express';

import type { MarketplaceServicePort } from '@modules/marketplace/application/ports';
import { registerStaticHandler } from '../handlers';

/**
 * Registers all marketplace-related HTTP routes on the Express app.
 * Currently, the marketplace mock serves only static plugin manifests;
 * this module is the single integration point for future endpoints.
 *
 * @param app - Express application instance.
 * @param service - Marketplace service port providing configuration.
 */
export function registerMarketplaceRoutes(app: Application, service: MarketplaceServicePort): void {
	registerStaticHandler(app, service);
}

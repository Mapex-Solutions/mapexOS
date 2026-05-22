import type { Application, RequestHandler } from 'express';
import express from 'express';
import cors from 'cors';

import type { MarketplaceServicePort } from '@modules/marketplace/application/ports';

/**
 * Registers middleware and static-file handlers on the Express app.
 *
 * The marketplace mock serves the public directory verbatim so clients can
 * fetch `/plugins/registry.json`, `/plugins/{pluginId}/manifest.json`, etc.
 *
 * @param app - Express application instance.
 * @param service - Marketplace service port providing the public directory path.
 */
export function registerStaticHandler(app: Application, service: MarketplaceServicePort): void {
	const corsMiddleware: RequestHandler = cors();
	app.use(corsMiddleware);
	app.use(express.static(service.getPublicDirectory()));
}

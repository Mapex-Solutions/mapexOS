import express from 'express';

import { createConsoleLogger } from '@shared/logger';
import { initMarketplaceModule } from '@modules/marketplace/module';

/**
 * Bootstrap: wires the Marketplace module into an Express app and starts
 * the HTTP server. All runtime configuration lives inside the module.
 */
function bootstrap(): void {
	const logger = createConsoleLogger();
	const app = express();

	const service = initMarketplaceModule(app, logger);
	const port = service.getPort();

	app.listen(port, () => {
		logger.info(`[SERVICE:Marketplace] Listening at http://localhost:${port}`);
		logger.info(`[SERVICE:Marketplace] Registry: http://localhost:${port}/plugins/registry.json`);
		logger.info(`[SERVICE:Marketplace] Telegram: http://localhost:${port}/plugins/telegram/manifest.json`);
	});
}

bootstrap();

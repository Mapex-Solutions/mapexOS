import { container } from 'tsyringe';
import type { Logger, Application } from '@mapexos/microservices';
import { LOGGER_TOKEN, APP_TOKEN } from '@mapexos/microservices';

import { registerRoutes } from '@modules/app/interfaces/http';
import { Modules } from '@shared/configuration';

/**
 * InitModule orchestrates all module initializations in phases
 * Following workspace_go pattern:
 * Phase 1: Repositories
 * Phase 2: Services
 * Phase 3: Interfaces (HTTP routes)
 * Phase 4: Listeners (NATS consumers)
 */
export async function initModule(): Promise<void> {
	const logger = container.resolve<Logger>(LOGGER_TOKEN);
	const app = container.resolve<Application>(APP_TOKEN);

	logger.info('[MODULE:APP] Initializing');

	// Register APP module routes
	app.use('/', registerRoutes());

	logger.info('[MODULE:APP] Initialized');
	logger.info({ total: Modules.length, modules: Modules.map(m => m.name) }, '[MODULE:DOMAIN] Initializing domain modules');

	// PHASE 1: Initialize all repositories
	logger.info('[MODULE:PHASE1] Initializing repositories');
	for (const mod of Modules) {
		if (!mod.lazy && mod.initRepositories) {
			logger.debug({ module: mod.name }, '[MODULE:PHASE1] Initializing repositories');
			mod.initRepositories();
		}
	}

	// PHASE 2: Initialize all services
	logger.info('[MODULE:PHASE2] Initializing services');
	for (const mod of Modules) {
		if (!mod.lazy && mod.initServices) {
			logger.debug({ module: mod.name }, '[MODULE:PHASE2] Initializing services');
			mod.initServices();
		}
	}

	// PHASE 3: Initialize all interfaces (HTTP routes)
	logger.info('[MODULE:PHASE3] Initializing interfaces');
	for (const mod of Modules) {
		if (!mod.lazy && mod.initInterfaces) {
			logger.debug({ module: mod.name }, '[MODULE:PHASE3] Initializing interfaces');
			mod.initInterfaces();
		}
	}

	// PHASE 4: Initialize listeners (NATS consumers)
	// Some listeners may be async (e.g., ensureFanoutStream)
	logger.info('[MODULE:PHASE4] Initializing listeners');
	for (const mod of Modules) {
		if (!mod.lazy && mod.initListeners) {
			logger.debug({ module: mod.name }, '[MODULE:PHASE4] Initializing listeners');
			await mod.initListeners();
		}
	}

	logger.info('[MODULE:DOMAIN] All modules initialized');
}

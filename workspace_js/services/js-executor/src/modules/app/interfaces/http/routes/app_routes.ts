import { Router } from '@mapexos/microservices';

/**
 * Registers and returns a group of HTTP routes for the application.
 * Note: /health is registered in bootstrap/health.ts with infrastructure checks.
 *
 * @returns {Router} A Router instance containing the registered routes.
 */
export function registerRoutes(): Router {
	/** Router instance */
	const group = Router();

	/** Return the router group */
	return group;
}

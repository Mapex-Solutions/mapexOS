import type { Logger } from '@shared/logger';
import type { MarketplaceServicePort } from '../ports';
import { MarketplaceService } from '../services';

/**
 * Dependencies required for MarketplaceService.
 * All fields use port interfaces or platform primitives — never concrete classes.
 */
export interface MarketplaceServiceDependencies {
	logger: Logger;
	/** Absolute path to the public directory that hosts plugin manifests. */
	publicDirectory: string;
	/** Port the HTTP server binds to. */
	port: number;
}

/**
 * Factory constructs the service with typed dependencies.
 * Returns the PORT interface, not the concrete class.
 */
export function createMarketplaceService(deps: MarketplaceServiceDependencies): MarketplaceServicePort {
	return new MarketplaceService(deps);
}

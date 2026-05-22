import type { MarketplaceServicePort } from '../ports';
import type { MarketplaceServiceDependencies } from '../di';

/**
 * MarketplaceService orchestrates the plugin marketplace mock.
 * It is a thin wrapper that exposes configuration resolved from the module
 * bootstrap (public directory path, port) via the port interface.
 */
export class MarketplaceService implements MarketplaceServicePort {
	constructor(private readonly deps: MarketplaceServiceDependencies) {}

	getPublicDirectory(): string {
		return this.deps.publicDirectory;
	}

	getPort(): number {
		return this.deps.port;
	}
}

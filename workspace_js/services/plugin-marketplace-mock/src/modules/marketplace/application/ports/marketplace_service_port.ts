/**
 * MarketplaceServicePort defines the contract for the plugin marketplace service.
 *
 * Responsibilities:
 * - Resolve the absolute path of the public directory (where static manifests live)
 * - Report the configured HTTP port
 *
 * This mock serves static files; business logic is intentionally minimal.
 */
export interface MarketplaceServicePort {
	/** Returns the absolute path to the public directory containing plugin manifests. */
	getPublicDirectory(): string;

	/** Returns the HTTP port the server should bind to. */
	getPort(): number;
}

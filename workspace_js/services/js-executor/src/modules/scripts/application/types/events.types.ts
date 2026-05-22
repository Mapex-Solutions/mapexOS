/**
 * Event types for script service message processing.
 *
 * These types define the structure of domain events consumed by the ScriptService.
 */

/**
 * Event published by asset service when assetTemplateId changes.
 *
 * This event triggers cache invalidation for the affected asset's scripts,
 * ensuring that the next execution fetches fresh scripts from the API.
 */
export interface AssetScriptsUpdatedEvent {
	/** The MongoDB ID of the asset */
	assetId: string;
	/** The UUID of the asset (device identifier) */
	assetUUID: string;
	/** The new asset template ID */
	assetTemplateId: string;
	/** The organization ID (for tenant context) */
	orgId?: string;
	/** The path key (for DLQ filtering) */
	pathKey?: string;
}

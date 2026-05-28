/**
 * ListsListPage Interfaces
 */

/**
 * Filter state for lists list page.
 *
 * Hierarchy filtering goes through `parentId` — the cascade picker resolves
 * the deepest selected ancestor (manufacturer if set, otherwise category),
 * then forwards its id as `parentId` to the backend.
 *
 * `category` and `manufacturer` levels are UI-only state used by the
 * cascade dropdowns; they don't reach the backend directly.
 */
export interface ListsListPageFilters {
	name: string | undefined;
	type: string | undefined;
	parentId: string | undefined;
	isSystem: boolean | undefined;
	isTemplate: boolean | undefined;
	includeChildren: boolean | undefined;
}

/**
 * Lookup option exposed by the cascade selects (q-select map-options shape)
 */
export interface ListCascadeOption {
	id: string;
	name: string;
}

/**
 * Column visibility state for lists list page
 */
export interface ListsListPageColumnVisibility {
	organization: boolean;
	description: boolean;
	parent: boolean;
	type: boolean;
	items: boolean;
	isTemplate: boolean;
	created: boolean;
	source: boolean;
	scope: boolean;
}

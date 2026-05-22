/**
 * ListsListPage Interfaces
 */

/**
 * Filter state for lists list page
 */
export interface ListsListPageFilters {
	name: string | undefined;
	category: 'iot' | undefined;
	type: string | undefined;
	isSystem: boolean | undefined;
	isTemplate: boolean | undefined;
	includeChildren: boolean | undefined;
}

/**
 * Column visibility state for lists list page
 */
export interface ListsListPageColumnVisibility {
	organization: boolean;
	description: boolean;
	category: boolean;
	type: boolean;
	items: boolean;
	isTemplate: boolean;
	created: boolean;
	source: boolean;
	scope: boolean;
}

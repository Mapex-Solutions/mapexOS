/**
 * ListsListPage Constants
 */

/**
 * Default number of items per page
 */
export const DEFAULT_ITEMS_PER_PAGE = 15;

/**
 * Default column visibility state for lists list
 */
export const COLUMN_VISIBILITY_DEFAULTS = {
	organization: true,
	description: true,
	parent: true,
	type: true,
	items: true,
	isTemplate: true,
	created: true,
	source: true,
	scope: true,
} as const;

/**
 * API projection fields for lists query
 */
export const LISTS_PROJECTION = 'name,value,type,parentId,parentName,parentType,enabled,isSystem,isTemplate,created,orgId' as const;

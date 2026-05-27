/**
 * CreateEditListPage Constants
 */

import type { ListFormData, ListType } from '../interfaces';

/**
 * Total number of steps in the form (Basic + Review)
 */
export const TOTAL_STEPS = 2;

/**
 * Step numbers enum for readability
 */
export const STEP = {
	BASIC_INFO: 1,
	REVIEW: 2,
} as const;

/**
 * Field length limits (mirror ZodListCreateSchema)
 */
export const NAME_MAX_LENGTH = 254;
export const VALUE_MAX_LENGTH = 254;

/**
 * Allowed pattern for "value" — lowercase letters, numbers and underscore.
 * Matches the convention seen in the seed (e.g. "iot", "milesight_iot", "am102").
 */
export const VALUE_PATTERN = /^[a-z0-9_]+$/;

/**
 * Parent type required when creating each list type.
 * `null` means the type has no parent (top of the tree).
 */
export const PARENT_TYPE_FOR: Record<ListType, ListType | null> = {
	asset_category: null,
	asset_manufacturer: 'asset_category',
	asset_model: 'asset_manufacturer',
};

/**
 * Initial form data values for a brand new list item
 */
export const INITIAL_LIST_FORM_DATA: ListFormData = {
	type: null,
	parentId: null,
	name: '',
	value: '',
	enabled: true,
	isTemplate: false,
};

/**
 * Icon used in the page header per type — falls back to "list" before a type is picked.
 */
export const TYPE_ICON: Record<ListType, string> = {
	asset_category: 'folder_special',
	asset_manufacturer: 'precision_manufacturing',
	asset_model: 'memory',
};

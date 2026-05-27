/**
 * CreateEditListPage Interfaces
 */

/**
 * Allowed list types. Mirrors backend enum and matches the
 * lookups seeded in mapex_iam/lists.json.
 */
export type ListType = 'asset_category' | 'asset_manufacturer' | 'asset_model';

/**
 * List form data structure - matches ZodListCreateSchema fields the UI exposes
 */
export interface ListFormData {
	/** List type (defines its place in the tree) */
	type: ListType | null;

	/** Parent list ID (null only allowed for asset_category) */
	parentId: string | null;

	/** Display name (1-254 chars) */
	name: string;

	/** Machine value used in payloads (1-254 chars, snake-case recommended) */
	value: string;

	/** Whether the list item is enabled */
	enabled: boolean;

	/** Whether the item is shared with child organizations as template */
	isTemplate: boolean;
}

/**
 * Parent option exposed by the cascade select
 */
export interface ParentOption {
	id: string;
	name: string;
	value: string;
}

/**
 * Props for Step1Basic component
 */
export interface Step1BasicProps {
	/** Current form data */
	modelValue: ListFormData;

	/** Whether the page is in edit mode (locks type) */
	isEditMode?: boolean;

	/** Whether the loaded item is a system list (read-only) */
	isSystemList?: boolean;

	/** Loading state for the parent options */
	loadingParents?: boolean;

	/** Available parent options for the currently chosen type */
	parentOptions: ParentOption[];
}

/**
 * Emits for Step1Basic component
 */
export interface Step1BasicEmits {
	(e: 'update:modelValue', value: Partial<ListFormData>): void;
}

/**
 * Props for Step2Review component
 */
export interface Step2ReviewProps {
	/** Current form data */
	modelValue: ListFormData;

	/** Whether the page is in edit mode */
	isEditMode?: boolean;

	/** Resolved parent name for display */
	parentName?: string | null;
}

/**
 * Emits for Step2Review component
 */
export interface Step2ReviewEmits {
	(e: 'editSection', stepNumber: number): void;
}

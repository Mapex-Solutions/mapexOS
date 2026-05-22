/**
 * Configuration for each filter field.
 */
import type { InputFilterType } from '../components/inputFilter/interfaces';

export interface FilterListOption {
	/** Label shown for the option */
	label: string;
	/** Value for the option */
	value: any;
	/** Optional icon name for the option */
	icon?: string;
	/** Optional color for the icon */
	color?: string;
}

export interface FilterListItem {
	/** Unique key for the filter value */
	key: string;
	/** Type of field: 'input', 'select', 'multiselect', 'daterange', or 'user-select' */
	type: 'input' | 'select' | 'multiselect' | 'daterange' | 'user-select';
	/** Label shown on the field */
	label: string;
	/** Icon name for the field */
	icon: string;
	/** Options when type is 'select' or 'multiselect' */
	options?: FilterListOption[];
	/** Optional grid classes (e.g. 'col-12 col-md-4') */
	grid?: string;
	/** Whether the field is clearable */
	clearable?: boolean;
	/** Whether the field is disabled (readonly) - controlled by parent */
	disabled?: boolean;
	/** Default value for the field */
	defaultValue?: any;
	/** Input mask pattern for Quasar (e.g., '########' for numbers only) */
	mask?: string | undefined;
	/** Input type (e.g., 'text', 'number', 'tel') - only for 'input' type */
	inputType?: InputFilterType | undefined;
	/** Placeholder text - only for 'input' type */
	placeholder?: string | undefined;
}

export interface FilterListProps {
	/** Array of filter field configurations */
	items: FilterListItem[];
	/** Enable auto-apply with debounce (default: false) */
	autoApply?: boolean;
	/** Show include children toggle in header (default: false) */
	showIncludeChildren?: boolean;
	/** Initial value for include children toggle (default: false) */
	includeChildrenInitial?: boolean;
	/** Array of field keys to watch for real-time changes (for business logic in parent) */
	watchFields?: string[];
}

export interface FilterListEmitEvents {
	/** Emitted when user clicks Search or auto-apply triggers */
	(e: 'apply', values: Record<string, any>): void;

	/** Emitted when user clicks Reset */
	(e: 'reset'): void;

	/** Emitted when a watched field changes (for real-time business logic) */
	(e: 'fieldChange', field: string, value: any): void;
}

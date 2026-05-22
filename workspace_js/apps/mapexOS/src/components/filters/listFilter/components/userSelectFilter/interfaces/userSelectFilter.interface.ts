/**
 * UserSelectFilter props interface
 */
export interface UserSelectFilterProps {
	/** Current selected user ID */
	modelValue: string | null;
	/** Label for the field */
	label: string;
	/** Icon name for the field */
	icon: string;
	/** Whether the field is clearable */
	clearable?: boolean | undefined;
	/** Whether the field is disabled */
	disabled?: boolean | undefined;
	/** Placeholder text */
	placeholder?: string | undefined;
}

/**
 * UserSelectFilter emits interface
 */
export interface UserSelectFilterEmits {
	(e: 'update:modelValue', value: string | null): void;
	(e: 'enter'): void;
}

/**
 * Selected user display info
 */
export interface SelectedUserInfo {
	id: string;
	name: string;
	email?: string | undefined;
}

/** Valid input types for q-input */
export type InputFilterType = 'text' | 'number' | 'email' | 'password' | 'search' | 'tel' | 'url' | 'time' | 'date' | 'datetime-local' | 'textarea' | 'file';

export interface InputFilterProps {
	/** Current value */
	modelValue: string;
	/** Field label */
	label: string;
	/** Icon name */
	icon: string;
	/** Whether the field is clearable */
	clearable?: boolean;
	/** Whether the field is disabled */
	disabled?: boolean;
	/** Debounce time in ms (for auto-apply) */
	debounce?: number;
	/** Input mask pattern (e.g., '########' for numbers only) */
	mask?: string | undefined;
	/** Input type (e.g., 'text', 'number', 'tel') */
	type?: InputFilterType | undefined;
	/** Placeholder text */
	placeholder?: string | undefined;
}

export interface InputFilterEmits {
	(e: 'update:modelValue', value: string): void;
	(e: 'enter'): void;
}

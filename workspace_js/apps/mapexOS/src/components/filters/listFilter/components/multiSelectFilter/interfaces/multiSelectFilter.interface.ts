export interface MultiSelectFilterOption {
	label: string;
	value: any;
	icon?: string;
	color?: string;
}

export interface MultiSelectFilterProps {
	/** Current value (array) */
	modelValue: any[];
	/** Field label */
	label: string;
	/** Icon name */
	icon: string;
	/** Select options */
	options: MultiSelectFilterOption[];
	/** Whether the field is clearable */
	clearable?: boolean;
	/** Whether the field is disabled */
	disabled?: boolean;
}

export interface MultiSelectFilterEmits {
	(e: 'update:modelValue', value: any[]): void;
	(e: 'enter'): void;
}

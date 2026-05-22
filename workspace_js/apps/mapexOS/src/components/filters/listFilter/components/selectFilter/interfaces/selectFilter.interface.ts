export interface SelectFilterOption {
	label: string;
	value: any;
	icon?: string;
	color?: string;
}

export interface SelectFilterProps {
	/** Current value */
	modelValue: any;
	/** Field label */
	label: string;
	/** Icon name */
	icon: string;
	/** Select options */
	options: SelectFilterOption[];
	/** Whether the field is clearable */
	clearable?: boolean;
	/** Whether the field is disabled */
	disabled?: boolean;
}

export interface SelectFilterEmits {
	(e: 'update:modelValue', value: any): void;
	(e: 'enter'): void;
}

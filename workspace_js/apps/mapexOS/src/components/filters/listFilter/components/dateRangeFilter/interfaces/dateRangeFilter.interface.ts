export interface DateRangeFilterProps {
	/** Current value (date range object with from/to) */
	modelValue: { from?: string; to?: string } | null;
	/** Field label */
	label: string;
	/** Icon name */
	icon: string;
	/** Whether the field is clearable */
	clearable?: boolean;
	/** Whether the field is disabled */
	disabled?: boolean;
}

export interface DateRangeFilterEmits {
	(e: 'update:modelValue', value: { from?: string; to?: string } | null): void;
	(e: 'enter'): void;
}

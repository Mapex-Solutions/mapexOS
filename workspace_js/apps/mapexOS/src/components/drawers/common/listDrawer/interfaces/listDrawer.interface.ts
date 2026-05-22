export interface ListDrawerProps {
	/** Controls drawer open/close state */
	modelValue: boolean;
	/** List ID to fetch details */
	listId?: string;
}

export interface ListDrawerEmits {
	(e: 'update:modelValue', value: boolean): void;
}

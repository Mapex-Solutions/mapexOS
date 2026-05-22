export interface HttpDataSourceDetailsDrawerProps {
	modelValue: boolean;
	dataSourceId?: string | undefined;
}

export interface HttpDataSourceDetailsDrawerEmits {
	(e: 'update:modelValue', value: boolean): void;
	(e: 'edit', dataSourceId: string): void;
}

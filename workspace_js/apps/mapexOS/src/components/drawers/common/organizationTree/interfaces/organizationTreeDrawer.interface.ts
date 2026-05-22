export interface OrganizationTreeDrawerProps {
	modelValue: boolean;
}

export interface OrganizationTreeDrawerEmits {
	(e: 'update:modelValue', value: boolean): void;
	(e: 'select', orgId: string): void;
}

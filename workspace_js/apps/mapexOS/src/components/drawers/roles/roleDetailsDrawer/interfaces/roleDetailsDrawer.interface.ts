export interface RoleDetailsDrawerProps {
	modelValue: boolean;
	roleId: string | null;
}

export interface RoleDetailsDrawerEmits {
	(e: 'update:modelValue', value: boolean): void;
	(e: 'edit', roleId: string): void;
}

export interface PermissionGroup {
	resource: string;
	actions: string[];
	count: number;
}

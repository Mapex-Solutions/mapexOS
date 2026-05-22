/**
 * GroupDetailsDrawer props interface
 */
export interface GroupDetailsDrawerProps {
	/** Whether the drawer is visible */
	modelValue: boolean;
	/** ID of the group to display */
	groupId: string | null;
}

/**
 * GroupDetailsDrawer emits interface
 */
export interface GroupDetailsDrawerEmits {
	(e: 'update:modelValue', value: boolean): void;
	(e: 'edit', groupId: string): void;
}

/**
 * Member display item
 */
export interface MemberDisplayItem {
	id: string;
	name: string;
	email?: string;
	avatar?: string;
}

/**
 * Role display item
 */
export interface RoleDisplayItem {
	id: string;
	name: string;
	description?: string;
	isSystem?: boolean;
}

import type { GroupDetailTab } from '../interfaces';

/**
 * Tab identifiers
 */
export const TAB = {
	INFO: 'info',
	ROLES: 'roles',
	MEMBERS: 'members',
} as const;

/**
 * Default tab on page load
 */
export const DEFAULT_TAB = TAB.INFO;

/**
 * Tab configuration factory (requires translations)
 *
 * @param {any} t - Translation composable
 * @returns {GroupDetailTab[]} Array of tab configurations
 */
export function getTabsConfig(t: any): GroupDetailTab[] {
	return [
		{
			name: TAB.INFO,
			label: t.tabs.info.value,
			icon: 'info',
		},
		{
			name: TAB.ROLES,
			label: t.tabs.roles.value,
			icon: 'admin_panel_settings',
		},
		{
			name: TAB.MEMBERS,
			label: t.tabs.members.value,
			icon: 'people',
		},
	];
}

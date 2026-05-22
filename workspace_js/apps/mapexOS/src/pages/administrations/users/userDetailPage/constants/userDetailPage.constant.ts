import type { UserDetailTab } from '../interfaces';

/**
 * Tab identifiers
 */
export const TAB = {
  PROFILE: 'profile',
  ACCESS: 'access',
  GROUPS: 'groups',
} as const;

/**
 * Default tab on page load
 */
export const DEFAULT_TAB = TAB.PROFILE;

/**
 * Tab configuration factory (requires translations)
 *
 * @param {any} t - Translation composable
 * @returns {UserDetailTab[]} Array of tab configurations
 */
export function getTabsConfig(t: any): UserDetailTab[] {
  return [
    {
      name: TAB.PROFILE,
      label: t.tabs.profile.value,
      icon: 'person',
    },
    {
      name: TAB.ACCESS,
      label: t.tabs.access.value,
      icon: 'admin_panel_settings',
    },
    {
      name: TAB.GROUPS,
      label: t.tabs.groups.value,
      icon: 'groups',
    },
  ];
}

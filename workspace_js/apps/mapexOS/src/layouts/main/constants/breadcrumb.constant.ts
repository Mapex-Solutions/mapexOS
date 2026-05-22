import type { RouteHierarchyEntry, MenuItem } from '../interfaces';

/** Action segment display config (add, edit, detail) */
export const ACTION_BREADCRUMBS: Record<string, { label: string; icon: string }> = {
  add: { label: 'New', icon: 'add' },
  edit: { label: 'Edit', icon: 'edit' },
  detail: { label: 'Details', icon: 'visibility' },
};

/** Default icon for unmapped segments */
export const DEFAULT_BREADCRUMB_ICON = 'chevron_right';

/**
 * Build route hierarchy map from a menu list.
 * Maps each route path to its parent category and page info.
 * Dashboard (/home) is excluded — the static icon in AppHeader handles it.
 *
 * @param {MenuItem[]} menuList - Translated menu items
 * @returns {Map<string, RouteHierarchyEntry>} Route hierarchy map
 */
export function buildRouteHierarchy(menuList: MenuItem[]): Map<string, RouteHierarchyEntry> {
  const map = new Map<string, RouteHierarchyEntry>();

  for (const item of menuList) {
    // Skip dashboard — static icon in AppHeader handles it
    if (item.to === '/home') continue;

    // Top-level item without children (e.g., LakeHouse)
    if (item.to && !item.children) {
      map.set(item.to, {
        page: { label: item.label ?? '', icon: item.icon ?? '' },
      });
    }

    // Items with children — each child gets the parent category
    if (item.children) {
      const parent = { label: item.label ?? '', icon: item.icon ?? '' };

      for (const child of item.children) {
        if (child.separator || !child.to) continue;

        map.set(child.to, {
          parent,
          page: { label: child.label ?? '', icon: child.icon ?? '' },
        });
      }
    }
  }

  // Routes not in menu (accessed via user menu, etc.)
  map.set('/my_profile', {
    page: { label: 'My Profile', icon: 'person' },
  });

  return map;
}

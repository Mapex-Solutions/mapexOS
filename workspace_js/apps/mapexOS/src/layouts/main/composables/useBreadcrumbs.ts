import type { BreadcrumbItem, MenuItem } from '../interfaces';
import type { ComputedRef } from 'vue';
import { computed } from 'vue';
import { useRoute } from 'vue-router';

import { buildRouteHierarchy, ACTION_BREADCRUMBS } from '../constants/breadcrumb.constant';

/** UUID pattern to detect dynamic IDs */
const UUID_PATTERN = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i;

/** Numeric ID pattern */
const NUMERIC_ID_PATTERN = /^\d+$/;

/** MongoDB ObjectId pattern */
const OBJECT_ID_PATTERN = /^[0-9a-f]{24}$/i;

/**
 * Check if segment is a dynamic ID (UUID, numeric, ObjectId)
 *
 * @param {string} segment - Route segment to check
 * @returns {boolean} True if segment is a dynamic ID
 */
function isDynamicId(segment: string): boolean {
  return (
    UUID_PATTERN.test(segment) ||
    NUMERIC_ID_PATTERN.test(segment) ||
    OBJECT_ID_PATTERN.test(segment)
  );
}

/**
 * Composable for generating breadcrumbs from current route.
 *
 * Uses buildRouteHierarchy (derived from translated menu) to always show
 * the full category context — enterprise pattern:
 *
 * ```
 * Dashboard > Administration > Users > New
 * Dashboard > Data Sources > HTTP
 * Dashboard > Logs & Executions > Event Tracer
 * Dashboard > Automation > Business Rules > Edit
 * ```
 *
 * @param {ComputedRef<MenuItem[]>} menuList - Reactive translated menu items
 * @returns {{ breadcrumbs: ComputedRef<BreadcrumbItem[]> }} Breadcrumbs computed ref
 */
export function useBreadcrumbs(menuList: ComputedRef<MenuItem[]>) {
  const route = useRoute();

  const breadcrumbs = computed<BreadcrumbItem[]>(() => {
    const hierarchy = buildRouteHierarchy(menuList.value);
    const segments = route.path.split('/').filter(Boolean);
    const items: BreadcrumbItem[] = [];

    // Find the longest matching base path in hierarchy
    let basePath = '';
    let actionStartIndex = segments.length;
    let accumulatedPath = '';

    for (let i = 0; i < segments.length; i++) {
      accumulatedPath += `/${segments[i]}`;

      if (hierarchy.has(accumulatedPath)) {
        basePath = accumulatedPath;
        actionStartIndex = i + 1;
      }
    }

    if (!basePath) return items;

    const entry = hierarchy.get(basePath)!;

    // 1. Parent category (non-clickable label)
    if (entry.parent) {
      items.push({
        label: entry.parent.label,
        icon: entry.parent.icon,
      });
    }

    // 2. Page (clickable link)
    items.push({
      label: entry.page.label,
      icon: entry.page.icon,
      to: basePath,
    });

    // 3. Action segments after base path (non-clickable)
    for (let i = actionStartIndex; i < segments.length; i++) {
      const segment = segments[i];
      if (!segment) continue;

      if (isDynamicId(segment)) continue;

      const action = ACTION_BREADCRUMBS[segment];
      if (action) {
        items.push({
          label: action.label,
          icon: action.icon,
        });
      }
    }

    return items;
  });

  return {
    breadcrumbs,
  };
}

/**
 * Sidebar menu item (leaf node or parent group)
 */
export interface MenuItem {
  /** Material icon name */
  icon?: string;
  /** Display label */
  label?: string;
  /** Route path */
  to?: string;
  /** Required permissions (ANY match grants access) */
  permissions?: string[];
  /** Visual separator flag */
  separator?: boolean;
  /** Nested children (for expandable groups) */
  children?: MenuItem[];
}

/**
 * Breadcrumb item for the AppHeader navigation trail
 */
export interface BreadcrumbItem {
  /** Display label */
  label: string;
  /** Material icon name */
  icon: string;
  /** Route path (clickable if defined) */
  to?: string | undefined;
}

/**
 * Route hierarchy entry mapping a route to its parent category and page info
 * Used by useBreadcrumbs to generate enterprise-pattern breadcrumbs
 */
export interface RouteHierarchyEntry {
  /** Parent category (displayed as non-clickable label) */
  parent?: { label: string; icon: string };
  /** Page info (displayed as clickable breadcrumb) */
  page: { label: string; icon: string };
}

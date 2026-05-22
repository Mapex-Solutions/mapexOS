// RolesListPage Interfaces

/**
 * Filter state for roles list page
 */
export interface RolesListPageFilters {
  name: string | undefined;
  isSystem: boolean | undefined;
  scope: 'global' | 'local' | undefined;
  permission: string | undefined;
  includeChildren: boolean | undefined;
  isTemplate: boolean | undefined;
}

/**
 * Column visibility state for roles list page
 */
export interface RolesListPageColumnVisibility {
  organization: boolean;
  description: boolean;
  permissions: boolean;
  scope: boolean;
  isTemplate: boolean;
  created: boolean;
}

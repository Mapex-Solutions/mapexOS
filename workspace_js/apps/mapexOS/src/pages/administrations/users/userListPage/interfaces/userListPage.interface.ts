// UserListPage Interfaces

/**
 * Filter state for users list page
 */
export interface UserListPageFilters {
  email: string | undefined;
  firstName: string | undefined;
  lastName: string | undefined;
  enabled: boolean | undefined;
  includeChildren: boolean | undefined;
}

/**
 * Column visibility state for users list page
 */
export interface UserListPageColumnVisibility {
  organization: boolean;
  email: boolean;
  jobTitle: boolean;
  groups: boolean;
  created: boolean;
}

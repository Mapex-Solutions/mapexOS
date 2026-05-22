// CustomersListPage Interfaces

/**
 * Filter state for customers list page
 */
export interface CustomersListPageFilters {
  name: string | undefined;
  enabled: boolean | undefined;
  includeChildren: boolean | undefined;
  type: 'customer' | 'site' | 'building' | 'floor' | 'zone' | undefined;
}

/**
 * Column visibility state for customers list page
 */
export interface CustomersListPageColumnVisibility {
  organization: boolean;
  address: boolean;
  created: boolean;
}

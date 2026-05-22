// GroupsListPage Interfaces

export interface GroupsListPageFilters {
  name: string | undefined;
  enabled: boolean | undefined;
  memberId: string | undefined;
  includeChildren: boolean | undefined;
}

export interface GroupsListPageColumnVisibility {
  organization: boolean;
  description: boolean;
  membersCount: boolean;
  created: boolean;
}

export interface GroupsListPageState {
  groupsList: any[];
  loading: boolean;
  error: string | undefined;
  currentPage: number;
  itemsPerPage: number;
  totalPages: number;
  totalItems: number;
  filters: GroupsListPageFilters;
  columnVisibilityState: GroupsListPageColumnVisibility;
}

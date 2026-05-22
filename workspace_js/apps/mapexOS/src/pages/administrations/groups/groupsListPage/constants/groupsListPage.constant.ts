// GroupsListPage Constants

export const GROUPS_LIST_PAGE_DEFAULTS = {
  ITEMS_PER_PAGE: 15,
  INITIAL_PAGE: 1,
} as const;

export const GROUPS_COLUMN_VISIBILITY_DEFAULTS = {
  organization: true,
  description: true,
  membersCount: true,
  created: true,
} as const;

export const GROUPS_FILTER_DEFAULTS = {
  name: undefined,
  enabled: undefined,
  memberId: undefined,
  includeChildren: undefined,
} as const;

export const GROUPS_PROJECTION = 'name,description,membersCount,enabled,orgId,created' as const;

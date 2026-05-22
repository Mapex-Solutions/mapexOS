// HttpDataSourcesListPage Constants

export const HTTP_DATASOURCES_LIST_PAGE_DEFAULTS = {
  ITEMS_PER_PAGE: 15,
  INITIAL_PAGE: 1,
} as const;

export const HTTP_DATASOURCES_COLUMN_VISIBILITY_DEFAULTS = {
  organization: true,
  assetBind: true,
  auth: true,
  mode: true,
} as const;

export const HTTP_DATASOURCES_FILTER_DEFAULTS = {
  name: undefined,
  enabled: undefined,
  includeChildren: undefined,
  mode: undefined,
  auth: undefined,
  assetBind: undefined,
} as const;

export const HTTP_DATASOURCES_PROJECTION = 'name,description,auth,assetBind,enabled,mode,orgId' as const;

export const LIST_CARD_ACTION = [
  { eventName: 'edit', icon: 'edit', label: 'Edit Data Source', color: 'grey-7' },
  { eventName: 'disable', icon: 'play_arrow', label: 'Enable / Disable Data Source', color: 'grey-7' },
  { eventName: 'delete', icon: 'delete', label: 'Delete Data Source', color: 'grey-7' },
];

export const LIST_BOTTOM_ACTION = {
  label: 'VIEW DATA SOURCE',
  icon: 'settings_input_antenna',
  color: 'primary',
};
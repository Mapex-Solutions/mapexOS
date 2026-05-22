/**
 * AssetsListPage Constants
 */

/**
 * Default number of items per page
 */
export const DEFAULT_ITEMS_PER_PAGE = 15;

/**
 * Assets list page defaults for pagination and items per page
 */
export const ASSETS_LIST_PAGE_DEFAULTS = {
  ITEMS_PER_PAGE: 15,
  INITIAL_PAGE: 1,
} as const;

/**
 * Default column visibility state for assets list
 */
export const ASSETS_COLUMN_VISIBILITY_DEFAULTS = {
  uuid: true,
  protocol: true,
  manufacturerModel: true,
  debugEnabled: true,
  healthStatus: true,
  organization: true,
} as const;

/**
 * Default filter values for assets list
 */
export const ASSETS_FILTER_DEFAULTS = {
  name: undefined,
  assetUUID: undefined,
  status: undefined,
  categoryId: undefined,
  manufacturerId: undefined,
  modelId: undefined,
  includeChildren: undefined,
} as const;

/**
 * API projection fields for assets query
 */
export const ASSETS_PROJECTION = 'name,description,category,protocol,assetUUID,enabled,debugEnabled,orgId,manufacturerName,modelName,healthStatus,healthStatusChangedAt,lastSeenAt' as const;

/**
 * List types for dynamic filters
 */
export const LIST_TYPE = {
  ASSET_CATEGORY: 'asset_category',
  ASSET_MANUFACTURER: 'asset_manufacturer',
  ASSET_MODEL: 'asset_model',
} as const;

/**
 * Watch fields for cascading filter changes
 */
export const WATCH_FIELDS = ['category', 'manufacture'];

/**
 * Default parameters for cascading filter API calls
 */
export const CASCADING_FILTER_DEFAULTS = {
  page: 1,
  perPage: 100,
} as const;

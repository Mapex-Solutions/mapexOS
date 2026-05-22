/**
 * AssetTemplateListPage Constants
 */

/**
 * Default number of items per page
 */
export const DEFAULT_ITEMS_PER_PAGE = 15;

/**
 * Asset template list page defaults for pagination and items per page
 */
export const ASSET_TEMPLATE_LIST_PAGE_DEFAULTS = {
  ITEMS_PER_PAGE: 15,
  INITIAL_PAGE: 1,
} as const;

/**
 * Default column visibility state for asset templates list
 */
export const ASSET_TEMPLATE_COLUMN_VISIBILITY_DEFAULTS = {
  organization: true,
  manufacturerModel: true,
  version: true,
  isSystem: true,
  isTemplate: true,
} as const;

/**
 * Default filter values for asset templates list
 */
export const ASSET_TEMPLATE_FILTER_DEFAULTS = {
  name: undefined,
  status: undefined,
  isSystem: undefined,
  isTemplate: undefined,
  categoryId: undefined,
  manufacturerId: undefined,
  modelId: undefined,
  includeChildren: undefined,
} as const;

/**
 * API projection fields for asset templates query
 */
export const ASSET_TEMPLATE_PROJECTION = 'name,description,categoryName,manufacturerName,modelName,version,enabled,isSystem,isTemplate,orgId' as const;

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

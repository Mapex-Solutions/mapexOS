// AssetsListPage Interfaces

import type { AssetResponse } from '@mapexos/schemas';

/**
 * Filter state for assets list page
 */
export interface AssetsListPageFilters {
  name: string | undefined;
  assetUUID: string | undefined;
  status: boolean | undefined;
  categoryId: string | undefined;
  manufacturerId: string | undefined;
  modelId: string | undefined;
  includeChildren: boolean | undefined;
}

/**
 * Column visibility state for assets list page
 */
export interface AssetsListPageColumnVisibility {
  uuid: boolean;
  protocol: boolean;
  manufacturerModel: boolean;
  debugEnabled: boolean;
  healthStatus: boolean;
  organization: boolean;
}

/**
 * Dynamic filter options for cascading select filters
 */
export interface DynamicFilterOptions {
  label: string;
  value: string;
}

/**
 * Asset with enriched organization name
 */
export interface EnrichedAsset extends AssetResponse {
  organizationName?: string;
}

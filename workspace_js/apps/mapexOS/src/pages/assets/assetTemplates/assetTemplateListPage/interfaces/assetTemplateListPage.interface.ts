// AssetTemplateListPage Interfaces

import type { AssetTemplateResponse } from '@mapexos/schemas';

/**
 * Filter state for asset templates list page
 */
export interface AssetTemplateListPageFilters {
  name: string | undefined;
  status: boolean | undefined;
  isSystem: boolean | undefined;
  isTemplate: boolean | undefined;
  categoryId: string | undefined;
  manufacturerId: string | undefined;
  modelId: string | undefined;
  includeChildren: boolean | undefined;
}

/**
 * Column visibility state for asset templates list page
 */
export interface AssetTemplateListPageColumnVisibility {
  organization: boolean;
  manufacturerModel: boolean;
  version: boolean;
  isSystem: boolean;
  isTemplate: boolean;
}

/**
 * Dynamic filter options for cascading select filters
 */
export interface DynamicFilterOptions {
  label: string;
  value: string;
}

/**
 * Asset template with enriched organization name and computed fields
 */
export interface EnrichedAssetTemplate extends AssetTemplateResponse {
  organizationName?: string;
  hasPreprocessor?: boolean;
  hasValidation?: boolean;
  hasConversion?: boolean;
  manufacturer?: string;
  deviceModel?: string;
}

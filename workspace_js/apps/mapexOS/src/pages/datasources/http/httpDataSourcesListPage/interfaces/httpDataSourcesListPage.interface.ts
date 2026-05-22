// HttpDataSourcesListPage Interfaces

import type { DataSourceResponse } from '@mapexos/schemas';

/**
 * EnrichedDataSource Interface
 * Extends DataSourceResponse with additional organization name for display purposes
 */
export interface EnrichedDataSource extends DataSourceResponse {
  organizationName?: string;
}

/**
 * HttpDataSourcesListPageFilters Interface
 * Defines the available filter options for the HTTP Data Sources list page
 */
export interface HttpDataSourcesListPageFilters {
  /** Filter by data source name */
  name: string | undefined;
  /** Filter by mode (push or pull) */
  mode: 'push' | 'pull' | undefined;
  /** Filter by enabled status */
  enabled: boolean | undefined;
  /** Filter by authentication type */
  auth: string | undefined;
  /** Filter by asset binding */
  assetBind: string | undefined;
  /** Include child organization data sources */
  includeChildren: boolean | undefined;
}

/**
 * HttpDataSourcesListPageColumnVisibility Interface
 * Controls which columns are visible in the data sources list
 */
export interface HttpDataSourcesListPageColumnVisibility {
  /** Show organization column */
  organization: boolean;
  /** Show asset binding column */
  assetBind: boolean;
  /** Show authentication column */
  auth: boolean;
  /** Show mode column */
  mode: boolean;
}

/**
 * HttpDataSourcesListPageState Interface
 * Complete state management for the HTTP Data Sources list page
 */
export interface HttpDataSourcesListPageState {
  /** List of data sources */
  dataSourcesList: DataSourceResponse[];
  /** Loading state */
  loading: boolean;
  /** Error message if any */
  error: string | undefined;
  /** Current page number */
  currentPage: number;
  /** Number of items per page */
  itemsPerPage: number;
  /** Total number of pages */
  totalPages: number;
  /** Total number of items */
  totalItems: number;
  /** Active filters */
  filters: HttpDataSourcesListPageFilters;
  /** Column visibility state */
  columnVisibilityState: HttpDataSourcesListPageColumnVisibility;
}

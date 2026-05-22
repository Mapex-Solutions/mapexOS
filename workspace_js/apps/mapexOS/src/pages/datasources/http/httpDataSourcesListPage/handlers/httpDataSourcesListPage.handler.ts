// HttpDataSourcesListPage Handlers

import type { Ref } from 'vue';
import type { DataSourceResponse, DataSourceQuery } from '@mapexos/schemas';
import type { ListHeaderMenuColumn } from '@components/headers';
import type {
  HttpDataSourcesListPageFilters,
  HttpDataSourcesListPageColumnVisibility,
  EnrichedDataSource,
} from '../interfaces';

import { apis } from '@services/mapex';
import { notifySuccess, notifyFail, dialogDelete } from '@utils/alert';
import { useOrganizationStore } from '@stores/organization';
import { useLogger } from '@composables/useLogger';
import { HTTP_DATASOURCES_PROJECTION } from '../constants';

const logger = useLogger('httpDataSourcesListPageHandler');

/**
 * Fetch data sources from API with current filters and pagination
 * @param {Ref<HttpDataSourcesListPageFilters>} filters - Current filter state
 * @param {Ref<number>} currentPage - Current page number
 * @param {Ref<number>} itemsPerPage - Items per page
 * @param {Ref<EnrichedDataSource[]>} dataSourcesList - List of data sources
 * @param {Ref<number>} totalPages - Total number of pages
 * @param {Ref<number>} totalItems - Total number of items
 * @param {Ref<boolean>} loading - Loading state
 * @param {Ref<string | undefined>} error - Error message
 * @returns {Promise<void>}
 */
export async function fetchDataSourcesHandler(
  filters: Ref<HttpDataSourcesListPageFilters>,
  currentPage: Ref<number>,
  itemsPerPage: Ref<number>,
  dataSourcesList: Ref<EnrichedDataSource[]>,
  totalPages: Ref<number>,
  totalItems: Ref<number>,
  loading: Ref<boolean>,
  error: Ref<string | undefined>,
): Promise<void> {
  if (!apis.httpGateway) {
    error.value = 'HTTP Gateway API not initialized';
    logger.error('HTTP Gateway API not initialized');
    return;
  }

  try {
    loading.value = true;
    error.value = undefined;

    // Build query parameters from filters and pagination state
    const queryParams: Record<string, any> = {
      page: currentPage.value,
      perPage: itemsPerPage.value,
      projection: HTTP_DATASOURCES_PROJECTION,
    };

    // Add active filters conditionally (only if they have values)
    if (filters.value.name) {
      queryParams.name = filters.value.name;
    }
    if (typeof filters.value.enabled === 'boolean') {
      queryParams.enabled = filters.value.enabled;
    }
    if (typeof filters.value.includeChildren === 'boolean') {
      queryParams.includeChildren = filters.value.includeChildren;
    }
    if (filters.value.mode) {
      queryParams.mode = filters.value.mode;
    }
    if (filters.value.auth) {
      queryParams.auth = filters.value.auth;
    }
    if (filters.value.assetBind) {
      queryParams.assetBind = filters.value.assetBind;
    }

    const response = await apis.httpGateway.datasource.list(queryParams as DataSourceQuery);

    // Enrich data sources with organization name
    const orgStore = useOrganizationStore();
    const enrichedDataSources = (response.items || []).map((dataSource: any) => {
      const organization = orgStore.flatList.find((org: any) => org.id === dataSource.orgId);
      return {
        ...dataSource,
        organizationName: organization?.name || 'Unknown',
      };
    });

    dataSourcesList.value = enrichedDataSources;

    // Update pagination state from response
    if (response.pagination) {
      totalPages.value = response.pagination.totalPages || 1;
      totalItems.value = response.pagination.totalItems || 0;
    }
  } catch (err: any) {
    logger.error('Error fetching data sources:', err);
    const errorMsg = err.message || 'Failed to fetch data sources';
    error.value = errorMsg;
    notifyFail({ message: errorMsg });
  } finally {
    loading.value = false;
  }
}

/**
 * Handle filter apply event from ListFilter component
 * @param {Record<string, any>} appliedFilters - Applied filter values
 * @param {Ref<HttpDataSourcesListPageFilters>} filters - Current filter state
 * @param {Ref<number>} currentPage - Current page number
 * @param {() => void} fetchCallback - Callback to fetch data sources
 * @returns {void}
 */
export function handleFilterApplyHandler(
  appliedFilters: Record<string, any>,
  filters: Ref<HttpDataSourcesListPageFilters>,
  currentPage: Ref<number>,
  fetchCallback: () => void,
): void {
  // Map filter values from ListFilter to API query parameters
  filters.value.name = appliedFilters.name || undefined;
  filters.value.enabled = appliedFilters.enabled;
  filters.value.includeChildren = appliedFilters.includeChildren;
  filters.value.mode = appliedFilters.mode || undefined;
  filters.value.auth = appliedFilters.auth || undefined;
  filters.value.assetBind = appliedFilters.assetBind || undefined;

  // Reset to first page when filters change
  currentPage.value = 1;

  // Fetch with new filters
  fetchCallback();
}

/**
 * Handle pagination navigation
 * @param {number} page - New page number
 * @param {Ref<number>} currentPage - Current page number
 * @param {() => void} fetchCallback - Callback to fetch data sources
 * @returns {void}
 */
export function handlePageChangeHandler(
  page: number,
  currentPage: Ref<number>,
  fetchCallback: () => void,
): void {
  currentPage.value = page;
  fetchCallback();
}

/**
 * Handle items per page change
 * @param {number} newValue - New items per page value
 * @param {Ref<number>} itemsPerPage - Items per page state
 * @param {Ref<number>} currentPage - Current page number
 * @param {() => void} fetchCallback - Callback to fetch data sources
 * @returns {void}
 */
export function handleItemsPerPageChangeHandler(
  newValue: number,
  itemsPerPage: Ref<number>,
  currentPage: Ref<number>,
  fetchCallback: () => void,
): void {
  itemsPerPage.value = newValue;
  currentPage.value = 1; // Reset to first page
  fetchCallback();
}

/**
 * Update menu columns when changed
 * @param {ListHeaderMenuColumn[]} columns - Updated column visibility states
 * @param {Ref<HttpDataSourcesListPageColumnVisibility>} columnVisibilityState - Column visibility state
 * @returns {void}
 */
export function handleColumnsUpdateHandler(
  columns: ListHeaderMenuColumn[],
  columnVisibilityState: Ref<HttpDataSourcesListPageColumnVisibility>,
): void {
  columns.forEach(col => {
    if (col.key === 'organization') columnVisibilityState.value.organization = col.visible;
    if (col.key === 'assetBind') columnVisibilityState.value.assetBind = col.visible;
    if (col.key === 'auth') columnVisibilityState.value.auth = col.visible;
    if (col.key === 'mode') columnVisibilityState.value.mode = col.visible;
  });
}

/**
 * View data source details - opens drawer on single click
 * @param {DataSourceResponse} dataSource - Data source to view
 * @param {Ref<string | undefined>} selectedDataSourceId - Selected data source ID state
 * @param {Ref<boolean>} detailsDrawerOpen - Details drawer open state
 * @returns {void}
 */
export function viewDetailsHandler(
  dataSource: DataSourceResponse,
  selectedDataSourceId: Ref<string | undefined>,
  detailsDrawerOpen: Ref<boolean>,
): void {
  if (dataSource.id) {
    selectedDataSourceId.value = dataSource.id;
    detailsDrawerOpen.value = true;
  }
}

/**
 * Edit data source - navigate to edit page
 * @param {any} dataSource - Data source to edit
 * @param {any} router - Vue router instance
 * @returns {void}
 */
export function editDataSourceHandler(dataSource: any, router: any): void {
  void router.push(`/data_sources/http/edit/${dataSource.id}`);
}

/**
 * Confirm delete data source with dialog
 * @param {DataSourceResponse} dataSource - Data source to delete
 * @param {any} t - Translation composable
 * @param {(dataSource: DataSourceResponse) => Promise<void>} deleteCallback - Delete callback function
 * @returns {Promise<void>}
 */
export async function confirmDeleteHandler(
  dataSource: DataSourceResponse,
  t: any,
  deleteCallback: (dataSource: DataSourceResponse) => Promise<void>,
): Promise<void> {
  const confirmed = await dialogDelete({
    title: t.deleteDialog.title.value,
    message: t.deleteDialog.message(dataSource.name || ''),
  });

  if (confirmed) {
    await deleteCallback(dataSource);
  }
}

/**
 * Delete data source from API and update list
 * @param {DataSourceResponse} dataSource - Data source to delete
 * @param {Ref<EnrichedDataSource[]>} dataSourcesList - List of data sources
 * @param {Ref<number>} totalItems - Total number of items
 * @param {Ref<number>} currentPage - Current page number
 * @param {any} t - Translation composable
 * @param {() => Promise<void>} fetchCallback - Callback to fetch data sources
 * @returns {Promise<void>}
 */
export async function deleteDataSourceHandler(
  dataSource: DataSourceResponse,
  dataSourcesList: Ref<EnrichedDataSource[]>,
  totalItems: Ref<number>,
  currentPage: Ref<number>,
  t: any,
  fetchCallback: () => Promise<void>,
): Promise<void> {
  if (!apis.httpGateway) {
    notifyFail({ message: t.errors.apiNotInitialized.value });
    return;
  }

  if (!dataSource.id) {
    notifyFail({ message: t.notifications.deleteFailed.value });
    return;
  }

  try {
    await apis.httpGateway.datasource.delete({ dataSourceId: dataSource.id });

    // Update total items count after successful deletion
    totalItems.value = Math.max(0, totalItems.value - 1);

    // Remove from local list
    const index = dataSourcesList.value.findIndex(r => r.id === dataSource.id);
    if (index !== -1) {
      dataSourcesList.value.splice(index, 1);
    }

    // If current page is now empty and not the first page, go to previous page
    if (dataSourcesList.value.length === 0 && currentPage.value > 1) {
      currentPage.value -= 1;
      await fetchCallback();
    }

    notifySuccess({ message: t.notifications.deleteSuccess.value });
  } catch (err: any) {
    logger.error('Error deleting data source:', err);
    notifyFail({ message: err.message || t.notifications.deleteFailed.value });
  }
}

/**
 * RouteGroupsListPage Handlers
 */

import type { Ref } from 'vue';
import type { ListHeaderMenuColumn } from '@components/headers';
import type {
  RouteGroupsListPageFilters,
  RouteGroupsListPageColumnVisibility,
  EnrichedRouteGroup,
} from '../interfaces';

import { apis } from '@services/mapex';
import { notifySuccess, notifyFail, notifyWarning, dialogDelete } from '@utils/alert';
import { cleanQueryParams } from '@utils/query';
import { useOrganizationStore } from '@stores/organization';
import { useLogger } from '@composables/useLogger';
import { ROUTE_GROUPS_PROJECTION } from '../constants';

const logger = useLogger('routeGroupsListPageHandler');

/**
 * Fetch route groups from API with current filters and pagination
 * @param {Ref<RouteGroupsListPageFilters>} filters - Current filter state
 * @param {Ref<number>} currentPage - Current page number
 * @param {Ref<number>} itemsPerPage - Items per page
 * @param {Ref<EnrichedRouteGroup[]>} routeGroupsList - List of route groups
 * @param {Ref<number>} totalPages - Total number of pages
 * @param {Ref<number>} totalItems - Total number of items
 * @param {Ref<boolean>} hasNext - Has next page flag
 * @param {Ref<boolean>} hasPrev - Has previous page flag
 * @param {Ref<boolean>} loading - Loading state
 * @param {Ref<string | undefined>} error - Error message
 * @returns {Promise<void>}
 */
export async function fetchRouteGroupsHandler(
  filters: Ref<RouteGroupsListPageFilters>,
  currentPage: Ref<number>,
  itemsPerPage: Ref<number>,
  routeGroupsList: Ref<EnrichedRouteGroup[]>,
  totalPages: Ref<number>,
  totalItems: Ref<number>,
  hasNext: Ref<boolean>,
  hasPrev: Ref<boolean>,
  loading: Ref<boolean>,
  error: Ref<string | undefined>,
): Promise<void> {
  if (!apis.router) {
    error.value = 'Router API not initialized';
    logger.error('Router API not initialized');
    return;
  }

  try {
    loading.value = true;
    error.value = undefined;

    // Build query parameters from filters and pagination state
    const queryParams: Record<string, any> = {
      page: currentPage.value,
      perPage: itemsPerPage.value,
      projection: ROUTE_GROUPS_PROJECTION,
    };

    // Add active filters conditionally (only if they have values)
    if (filters.value.name) {
      queryParams.name = filters.value.name;
    }
    if (typeof filters.value.enabled === 'boolean') {
      queryParams.enabled = filters.value.enabled;
    }
    if (typeof filters.value.isTemplate === 'boolean') {
      queryParams.isTemplate = filters.value.isTemplate;
    }
    if (typeof filters.value.includeChildren === 'boolean') {
      queryParams.includeChildren = filters.value.includeChildren;
    }

    // Clean undefined values to avoid sending "undefined" as string in URL
    const cleanedParams = cleanQueryParams(queryParams);

    const response = await apis.router.routegroup.list(cleanedParams);

    // Safely access response data with proper null checks
    const routeGroupsData = response?.items || [];

    // Enrich route groups with computed fields and organization name
    const orgStore = useOrganizationStore();
    const enrichedRouteGroups = routeGroupsData.map((routeGroup: any) => {
      const organization = orgStore.flatList.find((org: any) => org.id === routeGroup.orgId);
      return {
        ...routeGroup,
        organizationName: organization?.name || 'Unknown',
        routersCount: routeGroup.routers?.length || 0,
      };
    });

    routeGroupsList.value = enrichedRouteGroups;

    // Update pagination state from response
    if (response.pagination) {
      totalItems.value = response.pagination.totalItems || 0;
      totalPages.value = response.pagination.totalPages || 1;
      hasNext.value = currentPage.value < totalPages.value;
      hasPrev.value = currentPage.value > 1;
    }
  } catch (err: any) {
    logger.error('Error fetching route groups:', err);
    const errorMsg = err.message || 'Failed to fetch route groups';
    error.value = errorMsg;
    notifyFail({ message: errorMsg });
  } finally {
    loading.value = false;
  }
}

/**
 * Handle filter apply event from ListFilter component
 * @param {Record<string, any>} appliedFilters - Applied filter values
 * @param {Ref<RouteGroupsListPageFilters>} filters - Current filter state
 * @param {Ref<RouteGroupsListPageColumnVisibility>} columnVisibilityState - Column visibility state
 * @param {Ref<number>} currentPage - Current page number
 * @param {() => void} fetchCallback - Callback to fetch route groups
 * @returns {void}
 */
export function handleFilterApplyHandler(
  appliedFilters: Record<string, any>,
  filters: Ref<RouteGroupsListPageFilters>,
  columnVisibilityState: Ref<RouteGroupsListPageColumnVisibility>,
  currentPage: Ref<number>,
  fetchCallback: () => void,
): void {
  filters.value.name = appliedFilters.name || undefined;
  filters.value.enabled = appliedFilters.enabled;
  filters.value.isTemplate = appliedFilters.isTemplate;
  filters.value.includeChildren = appliedFilters.includeChildren;

  // Auto-hide columns to prevent horizontal scroll when includeChildren is active
  if (appliedFilters.includeChildren === true) {
    columnVisibilityState.value.routers = false;
    columnVisibilityState.value.isTemplate = false;
  } else {
    // Restore columns when includeChildren is disabled
    columnVisibilityState.value.routers = true;
    columnVisibilityState.value.isTemplate = true;
  }

  // Reset to first page when filters change
  currentPage.value = 1;

  // Fetch with new filters
  fetchCallback();
}

/**
 * Handle pagination navigation
 * @param {number} page - New page number
 * @param {Ref<number>} currentPage - Current page number
 * @param {() => void} fetchCallback - Callback to fetch route groups
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
 * @param {() => void} fetchCallback - Callback to fetch route groups
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
 * @param {Ref<RouteGroupsListPageColumnVisibility>} columnVisibilityState - Column visibility state
 * @returns {void}
 */
export function handleColumnsUpdateHandler(
  columns: ListHeaderMenuColumn[],
  columnVisibilityState: Ref<RouteGroupsListPageColumnVisibility>,
): void {
  columns.forEach(col => {
    if (col.key === 'organization') columnVisibilityState.value.organization = col.visible;
    if (col.key === 'routers') columnVisibilityState.value.routers = col.visible;
    if (col.key === 'isTemplate') columnVisibilityState.value.isTemplate = col.visible;
  });
}

/**
 * Check if user can edit/delete a route group
 * Rules:
 * - isTemplate = true: Can only edit/delete if orgId matches current organization
 * - isTemplate = false: Can always edit/delete (local resource)
 * @param {EnrichedRouteGroup} routeGroup - Route group to check
 * @param {string | null} selectedOrganizationId - Currently selected organization ID
 * @returns {boolean}
 */
export function canModifyRouteGroupHandler(
  routeGroup: EnrichedRouteGroup,
  selectedOrganizationId: string | null,
): boolean {
  // Shared templates can only be modified by the owner organization
  if (routeGroup.isTemplate) {
    return routeGroup.orgId === selectedOrganizationId;
  }

  // Local resources can always be modified
  return true;
}

/**
 * View route group details
 * @param {EnrichedRouteGroup} routeGroup - Route group to view
 * @param {any} t - Translations object
 * @returns {void}
 */
export function viewDetailsHandler(routeGroup: EnrichedRouteGroup, t: any): void {
  if (!routeGroup.id) {
    notifyFail({ message: t.errors.idMissing.value });
    return;
  }
  // TODO: Implement details drawer
  logger.debug('View route group details:', routeGroup);
}

/**
 * Edit route group - navigate to edit page
 * @param {EnrichedRouteGroup} routeGroup - Route group to edit
 * @param {string | null} selectedOrganizationId - Currently selected organization ID
 * @param {any} router - Vue router instance
 * @param {any} t - Translations object
 * @returns {void}
 */
export function editRouteGroupHandler(
  routeGroup: EnrichedRouteGroup,
  selectedOrganizationId: string | null,
  router: any,
  t: any,
): void {
  if (!routeGroup.id) {
    notifyFail({ message: t.errors.idMissing.value });
    return;
  }

  if (!canModifyRouteGroupHandler(routeGroup, selectedOrganizationId)) {
    notifyWarning({ message: t.notifications.sharedEdit.value });
    return;
  }

  void router.push(`/routing/route_groups/edit/${routeGroup.id}`);
}

/**
 * Confirm delete route group with dialog
 * @param {EnrichedRouteGroup} routeGroup - Route group to delete
 * @param {string | null} selectedOrganizationId - Currently selected organization ID
 * @param {any} t - Translations object
 * @param {(routeGroup: EnrichedRouteGroup) => Promise<void>} deleteCallback - Delete callback function
 * @returns {Promise<void>}
 */
export async function confirmDeleteHandler(
  routeGroup: EnrichedRouteGroup,
  selectedOrganizationId: string | null,
  t: any,
  deleteCallback: (routeGroup: EnrichedRouteGroup) => Promise<void>,
): Promise<void> {
  if (!canModifyRouteGroupHandler(routeGroup, selectedOrganizationId)) {
    notifyWarning({ message: t.notifications.sharedDelete.value });
    return;
  }

  const routeGroupName = routeGroup.name || 'this route group';
  const confirmed = await dialogDelete({
    title: t.dialog.confirmDelete.title.value,
    message: t.dialog.confirmDelete.message(routeGroupName),
  });

  if (confirmed) {
    await deleteCallback(routeGroup);
  }
}

/**
 * Delete route group from API and update list
 * @param {EnrichedRouteGroup} routeGroup - Route group to delete
 * @param {Ref<EnrichedRouteGroup[]>} routeGroupsList - List of route groups
 * @param {Ref<number>} totalItems - Total number of items
 * @param {Ref<number>} currentPage - Current page number
 * @param {any} t - Translations object
 * @param {() => Promise<void>} fetchCallback - Callback to fetch route groups
 * @returns {Promise<void>}
 */
export async function deleteRouteGroupHandler(
  routeGroup: EnrichedRouteGroup,
  routeGroupsList: Ref<EnrichedRouteGroup[]>,
  totalItems: Ref<number>,
  currentPage: Ref<number>,
  t: any,
  fetchCallback: () => Promise<void>,
): Promise<void> {
  if (!apis.router) {
    notifyFail({ message: t.errors.apiNotInitialized.value });
    return;
  }

  if (!routeGroup.id) {
    notifyFail({ message: t.errors.idMissing.value });
    return;
  }

  try {
    await apis.router.routegroup.delete({ routeGroupId: routeGroup.id });

    // Update total items count after successful deletion
    totalItems.value = Math.max(0, totalItems.value - 1);

    // Remove from local list
    const index = routeGroupsList.value.findIndex((r: EnrichedRouteGroup) => r.id === routeGroup.id);
    if (index !== -1) {
      routeGroupsList.value.splice(index, 1);
    }

    // If current page is now empty and not the first page, go to previous page
    if (routeGroupsList.value.length === 0 && currentPage.value > 1) {
      currentPage.value -= 1;
      await fetchCallback();
    }

    notifySuccess({ message: t.notifications.deleted.value });
  } catch (err: any) {
    logger.error('Error deleting route group:', err);
    notifyFail({ message: err.message || 'Failed to delete route group' });
  }
}

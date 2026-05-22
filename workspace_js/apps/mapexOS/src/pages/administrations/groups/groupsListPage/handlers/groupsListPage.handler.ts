// GroupsListPage Handlers

import type { Ref } from 'vue';
import type { GroupResponse, GroupQuery } from '@mapexos/schemas';
import type { ListHeaderMenuColumn } from '@components/headers';
import type {
  GroupsListPageFilters,
  GroupsListPageColumnVisibility,
} from '../interfaces';

import { apis } from '@services/mapex';
import { notifySuccess, notifyFail, notifyWarning, dialogDelete } from '@utils/alert';
import { useOrganizationStore } from '@stores/organization';
import { useLogger } from '@composables/useLogger';
import { GROUPS_PROJECTION } from '../constants';

const logger = useLogger('groupsListPageHandler');

/**
 * Fetch groups from API with current filters and pagination
 */
export async function fetchGroupsHandler(
  filters: Ref<GroupsListPageFilters>,
  currentPage: Ref<number>,
  itemsPerPage: Ref<number>,
  groupsList: Ref<GroupResponse[]>,
  totalPages: Ref<number>,
  totalItems: Ref<number>,
  loading: Ref<boolean>,
  error: Ref<string | undefined>,
) {
  if (!apis.mapexOS?.groups) {
    error.value = 'Groups API not initialized';
    return;
  }

  try {
    loading.value = true;
    error.value = undefined;

    // Build query parameters from filters and pagination state
    const queryParams: Record<string, any> = {
      page: currentPage.value,
      perPage: itemsPerPage.value,
      projection: GROUPS_PROJECTION,
    };

    // Add active filters conditionally (only if they have values)
    if (filters.value.name) {
      queryParams.name = filters.value.name;
    }
    if (typeof filters.value.enabled === 'boolean') {
      queryParams.enabled = filters.value.enabled;
    }
    if (filters.value.memberId) {
      queryParams.memberId = filters.value.memberId;
    }
    if (typeof filters.value.includeChildren === 'boolean') {
      queryParams.includeChildren = filters.value.includeChildren;
    }

    const response = await apis.mapexOS.groups.list(queryParams as GroupQuery);

    // Enrich groups with organization name
    const orgStore = useOrganizationStore();
    const enrichedGroups = (response.items || []).map((group: any) => {
      const organization = orgStore.flatList.find((org: any) => org.id === group.orgId);
      return {
        ...group,
        organizationName: organization?.name || 'Unknown',
      };
    });

    groupsList.value = enrichedGroups;

    // Update pagination state from response
    if (response.pagination) {
      totalPages.value = response.pagination.totalPages || 1;
      totalItems.value = response.pagination.totalItems || 0;
    }
  } catch (err: any) {
    logger.error('Error fetching groups:', err);
    const errorMsg = err.message || 'Failed to fetch groups';
    error.value = errorMsg;
    notifyFail({ message: errorMsg });
  } finally {
    loading.value = false;
  }
}

/**
 * Handle filter apply event from ListFilter component
 */
export function handleFilterApplyHandler(
  appliedFilters: Record<string, any>,
  filters: Ref<GroupsListPageFilters>,
  columnVisibilityState: Ref<GroupsListPageColumnVisibility>,
  currentPage: Ref<number>,
  fetchCallback: () => void,
) {
  // Map filter values from ListFilter to API query parameters
  filters.value.name = appliedFilters.name || undefined;
  filters.value.enabled = appliedFilters.enabled;
  filters.value.memberId = appliedFilters.memberId || undefined;
  filters.value.includeChildren = appliedFilters.includeChildren;

  // Auto-hide columns to prevent horizontal scroll when includeChildren is active
  if (appliedFilters.includeChildren === true) {
    columnVisibilityState.value.description = false;  // Hide description
    columnVisibilityState.value.membersCount = false;      // Hide members
    columnVisibilityState.value.created = false;      // Hide created
  } else {
    // Restore columns when includeChildren is disabled
    columnVisibilityState.value.description = true;
    columnVisibilityState.value.membersCount = true;
    columnVisibilityState.value.created = true;
  }

  // Reset to first page when filters change
  currentPage.value = 1;

  // Fetch with new filters
  fetchCallback();
}

/**
 * Handle pagination navigation
 */
export function handlePageChangeHandler(
  page: number,
  currentPage: Ref<number>,
  fetchCallback: () => void,
) {
  currentPage.value = page;
  fetchCallback();
}

/**
 * Handle items per page change
 */
export function handleItemsPerPageChangeHandler(
  newValue: number,
  itemsPerPage: Ref<number>,
  currentPage: Ref<number>,
  fetchCallback: () => void,
) {
  itemsPerPage.value = newValue;
  currentPage.value = 1; // Reset to first page
  fetchCallback();
}

/**
 * Update menu columns when changed
 */
export function handleColumnsUpdateHandler(
  columns: ListHeaderMenuColumn[],
  columnVisibilityState: Ref<GroupsListPageColumnVisibility>,
) {
  columns.forEach(col => {
    if (col.key === 'organization') columnVisibilityState.value.organization = col.visible;
    if (col.key === 'description') columnVisibilityState.value.description = col.visible;
    if (col.key === 'membersCount') columnVisibilityState.value.membersCount = col.visible;
    if (col.key === 'created') columnVisibilityState.value.created = col.visible;
  });
}

/**
 * Check if user can edit/delete a group
 * Rules:
 * - Can edit/delete if orgId matches current organization
 */
export function canModifyGroupHandler(group: any): boolean {
  const orgStore = useOrganizationStore();

  // Groups can only be modified by the owner organization
  return group.orgId === orgStore.selectedOrganizationId;
}

/**
 * View group details
 */
export function viewDetailsHandler(group: any) {
  logger.debug('View group details:', group);
  // TODO: Implement details drawer or navigation
}

/**
 * Edit group - Navigate to edit page
 */
export function editGroupHandler(group: any, t: any, router: any) {
  if (!canModifyGroupHandler(group)) {
    notifyWarning({ message: t.notifications.cannotEdit.value });
    return;
  }

  if (!group.id) {
    notifyFail({ message: t.errors.idMissing.value });
    return;
  }

  logger.debug('Navigating to edit group:', group.id);
  void router.push(`/groups/edit/${group.id}`);
}

/**
 * Confirm delete group
 */
export async function confirmDeleteHandler(
  group: GroupResponse,
  t: any,
  deleteCallback: (group: GroupResponse) => Promise<void>,
) {
  if (!canModifyGroupHandler(group)) {
    notifyWarning({ message: t.notifications.cannotDelete.value });
    return;
  }

  const groupName = group.name || 'this group';
  const confirmed = await dialogDelete({
    title: t.dialog.deleteTitle.value,
    message: t.messages.confirmDelete(groupName),
  });

  if (confirmed) {
    await deleteCallback(group);
  }
}

/**
 * Delete group
 */
export async function deleteGroupHandler(
  group: GroupResponse,
  groupsList: Ref<GroupResponse[]>,
  totalItems: Ref<number>,
  currentPage: Ref<number>,
  t: any,
  fetchCallback: () => Promise<void>,
) {
  if (!apis.mapexOS?.groups) {
    notifyFail({ message: t.errors.apiNotInitialized.value });
    return;
  }

  if (!group.id) {
    notifyFail({ message: t.errors.idMissing.value });
    return;
  }

  try {
    await apis.mapexOS.groups.delete({ groupId: group.id });

    // Update total items count after successful deletion
    totalItems.value = Math.max(0, totalItems.value - 1);

    // Remove from local list
    const index = groupsList.value.findIndex(r => r.id === group.id);
    if (index !== -1) {
      groupsList.value.splice(index, 1);
    }

    // If current page is now empty and not the first page, go to previous page
    if (groupsList.value.length === 0 && currentPage.value > 1) {
      currentPage.value -= 1;
      await fetchCallback();
    }

    notifySuccess({ message: t.messages.deletedSuccessfully.value });
  } catch (err: any) {
    notifyFail({ message: err.message || 'Failed to delete group' });
  }
}

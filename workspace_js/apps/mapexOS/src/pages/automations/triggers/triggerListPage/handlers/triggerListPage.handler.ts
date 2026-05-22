import type { TriggerResponse } from '@mapexos/schemas';

import { apis } from '@services/mapex';
import { handleApiError } from '@utils/error';

/**
 * Trigger list page filters state
 * Maps to TriggerQuery fields from backend
 */
export interface TriggerListPageFilters {
  /** Search by trigger name (partial match) */
  name: string | undefined;

  /** Filter by status (active/inactive) */
  status: boolean | undefined;

  /** Include triggers from child organizations */
  includeChildren: boolean | undefined;

  /** Filter by category (technical/communication) */
  category: string | undefined;

  /** Filter by trigger type */
  triggerType: string | undefined;
}

/**
 * Result returned from fetchTriggersHandler
 */
export interface FetchTriggersResult {
  /** List of triggers from API */
  triggers: TriggerResponse[];

  /** Total number of pages */
  totalPages: number;

  /** Total number of items across all pages */
  totalItems: number;
}

/**
 * Fetch triggers from API with current filters and pagination
 *
 * @param {TriggerListPageFilters} filters - Active filters
 * @param {number} currentPage - Current page number
 * @param {number} itemsPerPage - Items per page
 * @returns {Promise<FetchTriggersResult>} Triggers and pagination info
 */
export async function fetchTriggersHandler(
  filters: TriggerListPageFilters,
  currentPage: number,
  itemsPerPage: number
): Promise<FetchTriggersResult> {
  try {
    // Build query params
    const queryParams: Record<string, any> = {
      page: currentPage,
      perPage: itemsPerPage,
    };

    // Add filters conditionally
    if (filters.name) {
      queryParams.name = filters.name;
    }
    if (typeof filters.status === 'boolean') {
      queryParams.enabled = filters.status;
    }
    if (typeof filters.includeChildren === 'boolean') {
      queryParams.includeChildren = filters.includeChildren;
    }
    if (filters.category) {
      queryParams.category = filters.category;
    }
    if (filters.triggerType) {
      queryParams.triggerType = filters.triggerType;
    }

    // Call API
    const response = await apis.triggers.trigger.list(queryParams);

    // Map response to frontend structure
    return {
      triggers: response.items || [],
      totalPages: response.pagination.totalPages,
      totalItems: response.pagination.totalItems,
    };
  } catch (error: any) {
    handleApiError(error, {
      defaultMessage: 'Failed to load triggers',
      timeout: 5000,
    });

    // Return empty result on error
    return {
      triggers: [],
      totalPages: 0,
      totalItems: 0,
    };
  }
}

/**
 * Handle filter apply event from ListFilter component
 * Updates filter state and resets pagination to page 1
 *
 * @param {Record<string, any>} appliedFilters - Filters from ListFilter component
 * @param {TriggerListPageFilters} filters - Current filters object to update
 * @param {any} currentPage - Current page ref to reset
 * @param {() => void} fetchCallback - Callback to fetch data
 */
export function handleFilterApplyHandler(
  appliedFilters: Record<string, any>,
  filters: TriggerListPageFilters,
  currentPage: any,
  fetchCallback: () => void
): void {
  // Update filters (Row 1: Common Filters)
  filters.name = appliedFilters.name || undefined;
  filters.status = appliedFilters.status;
  filters.includeChildren = appliedFilters.includeChildren;

  // Update filters (Row 2: Domain-Specific Filters)
  filters.category = appliedFilters.category || undefined;
  filters.triggerType = appliedFilters.triggerType || undefined;

  // Reset to first page
  currentPage.value = 1;

  // Fetch data
  fetchCallback();
}

/**
 * Handle page change event from pagination component
 * Updates current page and refetches data
 *
 * @param {number} page - New page number (1-indexed)
 * @param {any} currentPage - Current page ref to update
 * @param {() => void} fetchCallback - Callback to fetch data
 */
export function handlePageChangeHandler(
  page: number,
  currentPage: any,
  fetchCallback: () => void
): void {
  currentPage.value = page;
  fetchCallback();
}

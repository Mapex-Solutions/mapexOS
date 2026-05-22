/** TYPE IMPORTS (ALL types first, grouped) */
import type { RouteGroupResponse } from '@mapexos/schemas';

/** VUE IMPORTS */
import { ref, computed } from 'vue';

/** SERVICES */
import { apis } from '@services/mapex';

/** UTILS */
import { handleApiError } from '@utils/error';
import { cleanQueryParams } from '@utils/query';

/**
 * Route group filters interface
 * Used for filtering route groups in list and selector views
 */
export interface RouteGroupFilters {
  /** Search by route group name */
  name?: string | undefined;
  /** Filter by status (active/inactive) */
  status?: boolean | undefined;
  /** Filter by system/custom templates */
  isSystem?: boolean | undefined;
  /** Filter by template/instance */
  isTemplate?: boolean | undefined;
  /** Include route groups from child organizations */
  includeChildren?: boolean | undefined;
  /** Filter route groups by router kinds (strict: every router must be in this set) */
  kinds?: string[] | undefined;
}

/**
 * Pagination state interface
 */
export interface PaginationState {
  /** Current page number (1-indexed) */
  currentPage: number;
  /** Items per page */
  perPage: number;
  /** Total number of pages */
  totalPages: number;
  /** Total number of items */
  totalItems: number;
  /** Has next page */
  hasNext: boolean;
  /** Has previous page */
  hasPrev: boolean;
}

/**
 * Composable for managing Route Groups (Routers)
 * Provides centralized state and methods for fetching, filtering, and paginating route groups
 *
 * Features:
 * - Fetch route groups with filters and pagination
 * - Infinite scroll support
 * - Loading states management
 *
 * @returns {Object} Routers composable methods and state
 */
export function useRouters() {
  /** STATE */
  const routeGroups = ref<RouteGroupResponse[]>([]);
  const isLoading = ref(false);
  const isLoadingMore = ref(false);

  const filters = ref<RouteGroupFilters>({
    name: undefined,
    status: undefined,
    isSystem: undefined,
    isTemplate: undefined,
    includeChildren: undefined,
    kinds: undefined,
  });

  const pagination = ref<PaginationState>({
    currentPage: 1,
    perPage: 30,
    totalPages: 1,
    totalItems: 0,
    hasNext: false,
    hasPrev: false,
  });

  /** COMPUTED */
  const hasActiveFilters = computed(() => {
    return !!(
      filters.value.name ||
      typeof filters.value.status === 'boolean' ||
      typeof filters.value.isSystem === 'boolean' ||
      typeof filters.value.isTemplate === 'boolean' ||
      typeof filters.value.includeChildren === 'boolean' ||
      (filters.value.kinds && filters.value.kinds.length > 0)
    );
  });

  /** FUNCTIONS */

  /**
   * Fetch route groups from API with current filters and pagination
   *
   * @param {boolean} append - If true, append results to existing list (for infinite scroll)
   * @returns {Promise<void>}
   */
  async function fetchRouteGroups(append = false): Promise<void> {
    if (!apis.router) {
      throw new Error('Router API not initialized');
    }

    if (append) {
      isLoadingMore.value = true;
    } else {
      isLoading.value = true;
      routeGroups.value = [];
    }

    try {
      const queryParams: Record<string, any> = {
        page: pagination.value.currentPage,
        perPage: pagination.value.perPage,
        sort: 'name:asc',
      };

      if (filters.value.name) queryParams.name = filters.value.name;
      if (typeof filters.value.status === 'boolean') queryParams.enabled = filters.value.status;
      if (typeof filters.value.includeChildren === 'boolean') queryParams.includeChildren = filters.value.includeChildren;
      if (filters.value.kinds && filters.value.kinds.length > 0) {
        queryParams.kinds = filters.value.kinds;
      }

      const cleanedParams = cleanQueryParams(queryParams);
      const response = await apis.router.routegroup.list(cleanedParams);

      if (append) {
        routeGroups.value.push(...response.items);
      } else {
        routeGroups.value = response.items;
      }

      pagination.value.totalItems = response.pagination.totalItems;
      pagination.value.totalPages = response.pagination.totalPages;
      pagination.value.hasNext = pagination.value.currentPage < pagination.value.totalPages;
      pagination.value.hasPrev = pagination.value.currentPage > 1;

    } catch (error: any) {
      handleApiError(error, {
        defaultMessage: 'Failed to fetch route groups',
        timeout: 5000,
      });
    } finally {
      isLoading.value = false;
      isLoadingMore.value = false;
    }
  }

  /**
   * Apply filters and reset pagination
   * @param {Partial<RouteGroupFilters>} newFilters - Filters to apply
   * @returns {Promise<void>}
   */
  async function applyFilters(newFilters: Partial<RouteGroupFilters>): Promise<void> {
    filters.value = { ...filters.value, ...newFilters };
    pagination.value.currentPage = 1;
    await fetchRouteGroups();
  }

  /**
   * Clear all filters and reset to defaults
   * @returns {Promise<void>}
   */
  async function clearFilters(): Promise<void> {
    filters.value = {
      name: undefined,
      status: undefined,
      includeChildren: undefined,
    };
    pagination.value.currentPage = 1;
    await fetchRouteGroups();
  }

  /**
   * Go to specific page
   * @param {number} page - Page number to navigate to
   * @returns {Promise<void>}
   */
  async function goToPage(page: number): Promise<void> {
    pagination.value.currentPage = page;
    await fetchRouteGroups();
  }

  /**
   * Load next page (infinite scroll)
   * @returns {Promise<void>}
   */
  async function loadMore(): Promise<void> {
    if (pagination.value.hasNext && !isLoadingMore.value) {
      pagination.value.currentPage++;
      await fetchRouteGroups(true);
    }
  }

  /**
   * Set items per page and reset to first page
   * @param {number} perPage - New items per page value
   * @returns {Promise<void>}
   */
  async function setItemsPerPage(perPage: number): Promise<void> {
    pagination.value.perPage = perPage;
    pagination.value.currentPage = 1;
    await fetchRouteGroups();
  }

  /**
   * Reset all state to defaults
   */
  function reset(): void {
    routeGroups.value = [];
    filters.value = {
      name: undefined,
      status: undefined,
      isSystem: undefined,
      isTemplate: undefined,
      includeChildren: undefined,
    };
    pagination.value = {
      currentPage: 1,
      perPage: 30,
      totalPages: 1,
      totalItems: 0,
      hasNext: false,
      hasPrev: false,
    };
  }

  return {
    // State
    routeGroups,
    isLoading,
    isLoadingMore,
    filters,
    pagination,

    // Computed
    hasActiveFilters,

    // Methods
    fetchRouteGroups,
    applyFilters,
    clearFilters,
    goToPage,
    loadMore,
    setItemsPerPage,
    reset,
  };
}

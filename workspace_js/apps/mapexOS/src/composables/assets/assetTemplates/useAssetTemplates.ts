/** TYPE IMPORTS (ALL types first, grouped) */
import type { AssetTemplateResponse } from '@mapexos/schemas';

/** VUE IMPORTS */
import { ref, computed } from 'vue';

/** COMPOSABLES */
import { useLogger } from '@composables/useLogger';

/** SERVICES */
import { apis } from '@services/mapex';

/** UTILS */
import { handleApiError } from '@utils/error';
import { cleanQueryParams } from '@utils/query';

const logger = useLogger('useAssetTemplates');

/**
 * Asset Template filters interface
 * Used for filtering asset templates in list and selector views
 */
export interface AssetTemplateFilters {
  /** Search by template name */
  name?: string | undefined;
  /** Filter by status (active/inactive) */
  status?: boolean | undefined;
  /** Filter by system flag */
  isSystem?: boolean | undefined;
  /** Filter by template flag */
  isTemplate?: boolean | undefined;
  /** Filter by category ID */
  categoryId?: string | undefined;
  /** Filter by manufacturer ID */
  manufacturerId?: string | undefined;
  /** Filter by model ID */
  modelId?: string | undefined;
  /** Include templates from child organizations */
  includeChildren?: boolean | undefined;
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
 * Cascading filter option interface
 * Used for category, manufacturer, and model dropdowns
 */
export interface CascadingFilterOption {
  /** Display label */
  label: string;
  /** Value (typically ID) */
  value: string;
}

/**
 * Composable for managing Asset Templates
 * Provides centralized state and methods for fetching, filtering, and paginating asset templates
 *
 * Features:
 * - Fetch asset templates with filters and pagination
 * - Cascading filters (category → manufacturer → model)
 * - Infinite scroll support
 * - Loading states management
 *
 * @returns {Object} Asset templates composable methods and state
 */
export function useAssetTemplates() {
  /** STATE */

  /**
   * List of fetched asset templates
   */
  const templates = ref<AssetTemplateResponse[]>([]);

  /**
   * Loading state for initial fetch
   */
  const isLoading = ref(false);

  /**
   * Loading state for infinite scroll / load more
   */
  const isLoadingMore = ref(false);

  /**
   * Filter state
   */
  const filters = ref<AssetTemplateFilters>({
    name: undefined,
    status: undefined,
    isSystem: undefined,
    isTemplate: undefined,
    categoryId: undefined,
    manufacturerId: undefined,
    modelId: undefined,
    includeChildren: undefined,
  });

  /**
   * Pagination state
   */
  const pagination = ref<PaginationState>({
    currentPage: 1,
    perPage: 30,
    totalPages: 1,
    totalItems: 0,
    hasNext: false,
    hasPrev: false,
  });

  /**
   * Cascading filter options
   */
  const categoryOptions = ref<CascadingFilterOption[]>([]);
  const manufacturerOptions = ref<CascadingFilterOption[]>([]);
  const modelOptions = ref<CascadingFilterOption[]>([]);

  /**
   * Loading states for cascading filters
   */
  const loadingCategories = ref(false);
  const loadingManufacturers = ref(false);
  const loadingModels = ref(false);

  /** COMPUTED */

  /**
   * Check if there are any active filters
   */
  const hasActiveFilters = computed(() => {
    return !!(
      filters.value.name ||
      typeof filters.value.status === 'boolean' ||
      typeof filters.value.isSystem === 'boolean' ||
      typeof filters.value.isTemplate === 'boolean' ||
      filters.value.categoryId ||
      filters.value.manufacturerId ||
      filters.value.modelId ||
      typeof filters.value.includeChildren === 'boolean'
    );
  });

  /** FUNCTIONS */

  /**
   * Fetch asset templates from API with current filters and pagination
   *
   * @param {boolean} append - If true, append results to existing list (for infinite scroll)
   * @returns {Promise<void>}
   */
  async function fetchTemplates(append = false): Promise<void> {
    if (!apis.assets) {
      throw new Error('Assets API not initialized');
    }

    // Set loading state
    if (append) {
      isLoadingMore.value = true;
    } else {
      isLoading.value = true;
      templates.value = [];
    }

    try {
      // Build query parameters
      const queryParams: Record<string, any> = {
        page: pagination.value.currentPage,
        perPage: pagination.value.perPage,
        sort: 'name:asc',
      };

      // Add active filters to query params
      if (filters.value.name) queryParams.name = filters.value.name;
      if (typeof filters.value.status === 'boolean') queryParams.enabled = filters.value.status;
      if (typeof filters.value.isSystem === 'boolean') queryParams.isSystem = filters.value.isSystem;
      if (typeof filters.value.isTemplate === 'boolean') queryParams.isTemplate = filters.value.isTemplate;
      if (filters.value.categoryId) queryParams.categoryId = filters.value.categoryId;
      if (filters.value.manufacturerId) queryParams.manufacturerId = filters.value.manufacturerId;
      if (filters.value.modelId) queryParams.modelId = filters.value.modelId;
      if (typeof filters.value.includeChildren === 'boolean') queryParams.includeChildren = filters.value.includeChildren;

      // Clean undefined values
      const cleanedParams = cleanQueryParams(queryParams);

      // Fetch from API
      const response = await apis.assets.assetTemplate.list(cleanedParams);

      // Update templates list
      if (append) {
        templates.value.push(...response.items);
      } else {
        templates.value = response.items;
      }

      // Update pagination state
      pagination.value.totalItems = response.pagination.totalItems;
      pagination.value.totalPages = response.pagination.totalPages;
      pagination.value.hasNext = pagination.value.currentPage < pagination.value.totalPages;
      pagination.value.hasPrev = pagination.value.currentPage > 1;

    } catch (error: any) {
      handleApiError(error, {
        defaultMessage: 'Failed to fetch asset templates',
        timeout: 5000,
      });
    } finally {
      isLoading.value = false;
      isLoadingMore.value = false;
    }
  }

  /**
   * Load categories from API for cascading filters
   * @returns {Promise<void>}
   */
  async function loadCategories(): Promise<void> {
    if (!apis.mapexOS?.lists) {
      return;
    }

    try {
      loadingCategories.value = true;
      const response = await apis.mapexOS.lists.list({
        type: 'asset_category',
        page: 1,
        perPage: 100,
      });

      categoryOptions.value = response.items.map((item: any) => ({
        label: item.name,
        value: item.id,
      }));
    } catch (error: any) {
      logger.error('Error loading categories:', error);
      categoryOptions.value = [];
    } finally {
      loadingCategories.value = false;
    }
  }

  /**
   * Load manufacturers based on selected category
   * @param {string} categoryId - Category ID to fetch manufacturers for
   * @returns {Promise<void>}
   */
  async function loadManufacturers(categoryId?: string): Promise<void> {
    const targetCategoryId = categoryId || filters.value.categoryId;

    if (!apis.mapexOS?.lists || !targetCategoryId) {
      manufacturerOptions.value = [];
      return;
    }

    try {
      loadingManufacturers.value = true;
      const response = await apis.mapexOS.lists.list({
        type: 'asset_manufacturer',
        parentId: targetCategoryId,
        page: 1,
        perPage: 100,
      });

      manufacturerOptions.value = response.items.map((item: any) => ({
        label: item.name,
        value: item.id,
      }));
    } catch (error: any) {
      logger.error('Error loading manufacturers:', error);
      manufacturerOptions.value = [];
    } finally {
      loadingManufacturers.value = false;
    }
  }

  /**
   * Load models based on selected manufacturer
   * @param {string} manufacturerId - Manufacturer ID to fetch models for
   * @returns {Promise<void>}
   */
  async function loadModels(manufacturerId?: string): Promise<void> {
    const targetManufacturerId = manufacturerId || filters.value.manufacturerId;

    if (!apis.mapexOS?.lists || !targetManufacturerId) {
      modelOptions.value = [];
      return;
    }

    try {
      loadingModels.value = true;
      const response = await apis.mapexOS.lists.list({
        type: 'asset_model',
        parentId: targetManufacturerId,
        page: 1,
        perPage: 100,
      });

      modelOptions.value = response.items.map((item: any) => ({
        label: item.name,
        value: item.id,
      }));
    } catch (error: any) {
      logger.error('Error loading models:', error);
      modelOptions.value = [];
    } finally {
      loadingModels.value = false;
    }
  }

  /**
   * Handle category change in cascading filters
   * Resets dependent filters and loads manufacturers
   * @param {string | undefined} categoryId - New category ID
   * @returns {Promise<void>}
   */
  async function handleCategoryChange(categoryId: string | undefined): Promise<void> {
    filters.value.categoryId = categoryId;
    filters.value.manufacturerId = undefined;
    filters.value.modelId = undefined;
    manufacturerOptions.value = [];
    modelOptions.value = [];

    if (categoryId) {
      await loadManufacturers(categoryId);
    }
  }

  /**
   * Handle manufacturer change in cascading filters
   * Resets dependent filters and loads models
   * @param {string | undefined} manufacturerId - New manufacturer ID
   * @returns {Promise<void>}
   */
  async function handleManufacturerChange(manufacturerId: string | undefined): Promise<void> {
    filters.value.manufacturerId = manufacturerId;
    filters.value.modelId = undefined;
    modelOptions.value = [];

    if (manufacturerId) {
      await loadModels(manufacturerId);
    }
  }

  /**
   * Apply filters and reset pagination
   * @param {Partial<AssetTemplateFilters>} newFilters - Filters to apply
   * @returns {Promise<void>}
   */
  async function applyFilters(newFilters: Partial<AssetTemplateFilters>): Promise<void> {
    filters.value = { ...filters.value, ...newFilters };
    pagination.value.currentPage = 1;
    await fetchTemplates();
  }

  /**
   * Clear all filters and reset to defaults
   * @returns {Promise<void>}
   */
  async function clearFilters(): Promise<void> {
    filters.value = {
      name: undefined,
      status: undefined,
      isSystem: undefined,
      isTemplate: undefined,
      categoryId: undefined,
      manufacturerId: undefined,
      modelId: undefined,
      includeChildren: undefined,
    };
    manufacturerOptions.value = [];
    modelOptions.value = [];
    pagination.value.currentPage = 1;
    await fetchTemplates();
  }

  /**
   * Go to specific page
   * @param {number} page - Page number to navigate to
   * @returns {Promise<void>}
   */
  async function goToPage(page: number): Promise<void> {
    pagination.value.currentPage = page;
    await fetchTemplates();
  }

  /**
   * Load next page (infinite scroll)
   * @returns {Promise<void>}
   */
  async function loadMore(): Promise<void> {
    if (pagination.value.hasNext && !isLoadingMore.value) {
      pagination.value.currentPage++;
      await fetchTemplates(true);
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
    await fetchTemplates();
  }

  /**
   * Reset all state to defaults
   */
  function reset(): void {
    templates.value = [];
    filters.value = {
      name: undefined,
      status: undefined,
      isSystem: undefined,
      isTemplate: undefined,
      categoryId: undefined,
      manufacturerId: undefined,
      modelId: undefined,
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
    categoryOptions.value = [];
    manufacturerOptions.value = [];
    modelOptions.value = [];
  }

  return {
    // State
    templates,
    isLoading,
    isLoadingMore,
    filters,
    pagination,
    categoryOptions,
    manufacturerOptions,
    modelOptions,
    loadingCategories,
    loadingManufacturers,
    loadingModels,

    // Computed
    hasActiveFilters,

    // Methods
    fetchTemplates,
    loadCategories,
    loadManufacturers,
    loadModels,
    handleCategoryChange,
    handleManufacturerChange,
    applyFilters,
    clearFilters,
    goToPage,
    loadMore,
    setItemsPerPage,
    reset,
  };
}

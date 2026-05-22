/** TYPE IMPORTS (ALL types first, grouped) */
import type { AssetResponse } from '@mapexos/schemas';

/** VUE IMPORTS */
import { ref, computed } from 'vue';

/** COMPOSABLES */
import { useLogger } from '@composables/useLogger';

/** SERVICES */
import { apis } from '@services/mapex';

/** UTILS */
import { handleApiError } from '@utils/error';
import { cleanQueryParams } from '@utils/query';

const logger = useLogger('useAssets');

/**
 * Asset filters interface
 * Used for filtering assets in list and selector views
 */
export interface AssetFilters {
  /** Search by asset name */
  name?: string | undefined;
  /** Search by asset UUID */
  assetUUID?: string | undefined;
  /** Filter by status (active/inactive) */
  status?: boolean | undefined;
  /** Filter by category ID */
  categoryId?: string | undefined;
  /** Filter by manufacturer ID */
  manufacturerId?: string | undefined;
  /** Filter by model ID */
  modelId?: string | undefined;
  /** Include assets from child organizations */
  includeChildren?: boolean | undefined;
  /** Filter by asset template ID */
  assetTemplateId?: string | undefined;
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
 * Composable for managing Assets
 * Provides centralized state and methods for fetching, filtering, and paginating assets
 *
 * Features:
 * - Fetch assets with filters and pagination
 * - Cascading filters (category → manufacturer → model)
 * - Infinite scroll support
 * - Loading states management
 *
 * @returns {Object} Assets composable methods and state
 */
export function useAssets() {
  /** STATE */
  const assets = ref<AssetResponse[]>([]);
  const isLoading = ref(false);
  const isLoadingMore = ref(false);

  const filters = ref<AssetFilters>({
    name: undefined,
    assetUUID: undefined,
    status: undefined,
    categoryId: undefined,
    manufacturerId: undefined,
    modelId: undefined,
    includeChildren: undefined,
    assetTemplateId: undefined,
  });

  const pagination = ref<PaginationState>({
    currentPage: 1,
    perPage: 30,
    totalPages: 1,
    totalItems: 0,
    hasNext: false,
    hasPrev: false,
  });

  const categoryOptions = ref<CascadingFilterOption[]>([]);
  const manufacturerOptions = ref<CascadingFilterOption[]>([]);
  const modelOptions = ref<CascadingFilterOption[]>([]);

  const loadingCategories = ref(false);
  const loadingManufacturers = ref(false);
  const loadingModels = ref(false);

  /** COMPUTED */
  const hasActiveFilters = computed(() => {
    return !!(
      filters.value.name ||
      filters.value.assetUUID ||
      typeof filters.value.status === 'boolean' ||
      filters.value.categoryId ||
      filters.value.manufacturerId ||
      filters.value.modelId ||
      typeof filters.value.includeChildren === 'boolean'
    );
  });

  /** FUNCTIONS */

  /**
   * Fetch assets from API with current filters and pagination
   *
   * @param {boolean} append - If true, append results to existing list (for infinite scroll)
   * @returns {Promise<void>}
   */
  async function fetchAssets(append = false): Promise<void> {
    if (!apis.assets) {
      throw new Error('Assets API not initialized');
    }

    if (append) {
      isLoadingMore.value = true;
    } else {
      isLoading.value = true;
      assets.value = [];
    }

    try {
      const queryParams: Record<string, any> = {
        page: pagination.value.currentPage,
        perPage: pagination.value.perPage,
        sort: 'name:asc',
      };

      if (filters.value.name) queryParams.name = filters.value.name;
      if (filters.value.assetUUID) queryParams.assetUUID = filters.value.assetUUID;
      if (typeof filters.value.status === 'boolean') queryParams.enabled = filters.value.status;
      if (filters.value.categoryId) queryParams.categoryId = filters.value.categoryId;
      if (filters.value.manufacturerId) queryParams.manufacturerId = filters.value.manufacturerId;
      if (filters.value.modelId) queryParams.modelId = filters.value.modelId;
      if (typeof filters.value.includeChildren === 'boolean') queryParams.includeChildren = filters.value.includeChildren;
      if (filters.value.assetTemplateId) queryParams.assetTemplateId = filters.value.assetTemplateId;

      const cleanedParams = cleanQueryParams(queryParams);
      const response = await apis.assets.asset.list(cleanedParams);

      if (append) {
        assets.value.push(...response.items);
      } else {
        assets.value = response.items;
      }

      pagination.value.totalItems = response.pagination.totalItems;
      pagination.value.totalPages = response.pagination.totalPages;
      pagination.value.hasNext = pagination.value.currentPage < pagination.value.totalPages;
      pagination.value.hasPrev = pagination.value.currentPage > 1;

    } catch (error: any) {
      handleApiError(error, {
        defaultMessage: 'Failed to fetch assets',
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
   * @param {Partial<AssetFilters>} newFilters - Filters to apply
   * @returns {Promise<void>}
   */
  async function applyFilters(newFilters: Partial<AssetFilters>): Promise<void> {
    filters.value = { ...filters.value, ...newFilters };
    pagination.value.currentPage = 1;
    await fetchAssets();
  }

  /**
   * Clear all filters and reset to defaults
   * @returns {Promise<void>}
   */
  async function clearFilters(): Promise<void> {
    filters.value = {
      name: undefined,
      assetUUID: undefined,
      status: undefined,
      categoryId: undefined,
      manufacturerId: undefined,
      modelId: undefined,
      includeChildren: undefined,
    };
    manufacturerOptions.value = [];
    modelOptions.value = [];
    pagination.value.currentPage = 1;
    await fetchAssets();
  }

  /**
   * Go to specific page
   * @param {number} page - Page number to navigate to
   * @returns {Promise<void>}
   */
  async function goToPage(page: number): Promise<void> {
    pagination.value.currentPage = page;
    await fetchAssets();
  }

  /**
   * Load next page (infinite scroll)
   * @returns {Promise<void>}
   */
  async function loadMore(): Promise<void> {
    if (pagination.value.hasNext && !isLoadingMore.value) {
      pagination.value.currentPage++;
      await fetchAssets(true);
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
    await fetchAssets();
  }

  /**
   * Reset all state to defaults
   */
  function reset(): void {
    assets.value = [];
    filters.value = {
      name: undefined,
      assetUUID: undefined,
      status: undefined,
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
    assets,
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
    fetchAssets,
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

import { describe, it, expect, beforeEach, vi } from 'vitest';
import { setActivePinia, createPinia } from 'pinia';
import { useAssets } from './useAssets';
import { apis } from '@services/mapex';

vi.mock('@composables/useLogger', () => ({
  useLogger: () => ({
    info: vi.fn(),
    warn: vi.fn(),
    error: vi.fn(),
    debug: vi.fn(),
  }),
}));

vi.mock('@utils/error', () => ({
  handleApiError: vi.fn(),
}));

vi.mock('@utils/query', () => ({
  cleanQueryParams: (params: Record<string, any>) => params,
}));

/**
 * Helper to build a fake paginated API response
 */
function makePaginatedResponse(items: any[], totalItems: number, totalPages: number) {
  return {
    items,
    pagination: { totalItems, totalPages },
  };
}

describe('useAssets', () => {
  beforeEach(() => {
    setActivePinia(createPinia());
    vi.clearAllMocks();
  });

  describe('initial state', () => {
    it('returns empty assets and default pagination', () => {
      const { assets, isLoading, isLoadingMore, pagination, filters } = useAssets();

      expect(assets.value).toEqual([]);
      expect(isLoading.value).toBe(false);
      expect(isLoadingMore.value).toBe(false);
      expect(pagination.value.currentPage).toBe(1);
      expect(pagination.value.perPage).toBe(30);
      expect(pagination.value.totalItems).toBe(0);
      expect(pagination.value.hasNext).toBe(false);
      expect(pagination.value.hasPrev).toBe(false);
      expect(filters.value.name).toBeUndefined();
    });
  });

  describe('hasActiveFilters', () => {
    it('returns false when no filters set', () => {
      const { hasActiveFilters } = useAssets();
      expect(hasActiveFilters.value).toBe(false);
    });

    it('returns true when name filter set', () => {
      const { hasActiveFilters, filters } = useAssets();
      filters.value.name = 'sensor';
      expect(hasActiveFilters.value).toBe(true);
    });

    it('returns true when status filter set to false', () => {
      const { hasActiveFilters, filters } = useAssets();
      filters.value.status = false;
      expect(hasActiveFilters.value).toBe(true);
    });

    it('returns true when categoryId filter set', () => {
      const { hasActiveFilters, filters } = useAssets();
      filters.value.categoryId = 'cat-1';
      expect(hasActiveFilters.value).toBe(true);
    });

    it('returns true when includeChildren is boolean', () => {
      const { hasActiveFilters, filters } = useAssets();
      filters.value.includeChildren = true;
      expect(hasActiveFilters.value).toBe(true);
    });
  });

  describe('fetchAssets', () => {
    it('throws when assets API is not initialized', async () => {
      (apis as any).assets = undefined;
      const { fetchAssets } = useAssets();
      await expect(fetchAssets()).rejects.toThrow('Assets API not initialized');
    });

    it('sets isLoading during fetch and populates assets', async () => {
      const mockItems = [{ id: '1', name: 'Asset A' }, { id: '2', name: 'Asset B' }];
      (apis as any).assets = {
        asset: {
          list: vi.fn().mockResolvedValue(makePaginatedResponse(mockItems, 2, 1)),
        },
      };

      const { fetchAssets, assets, isLoading, pagination } = useAssets();

      const promise = fetchAssets();
      // isLoading should be true while fetching
      expect(isLoading.value).toBe(true);

      await promise;

      expect(isLoading.value).toBe(false);
      expect(assets.value).toEqual(mockItems);
      expect(pagination.value.totalItems).toBe(2);
      expect(pagination.value.totalPages).toBe(1);
      expect(pagination.value.hasNext).toBe(false);
    });

    it('sets isLoadingMore when appending', async () => {
      (apis as any).assets = {
        asset: {
          list: vi.fn().mockResolvedValue(makePaginatedResponse([], 0, 1)),
        },
      };

      const { fetchAssets, isLoadingMore } = useAssets();
      const promise = fetchAssets(true);
      expect(isLoadingMore.value).toBe(true);
      await promise;
      expect(isLoadingMore.value).toBe(false);
    });

    it('clears assets on non-append fetch', async () => {
      (apis as any).assets = {
        asset: {
          list: vi.fn().mockResolvedValue(makePaginatedResponse([{ id: '1' }], 1, 1)),
        },
      };

      const { fetchAssets, assets } = useAssets();
      assets.value = [{ id: 'old' }];

      const promise = fetchAssets();
      // Assets should be cleared immediately
      expect(assets.value).toEqual([]);
      await promise;
    });

    it('computes hasNext and hasPrev correctly', async () => {
      (apis as any).assets = {
        asset: {
          list: vi.fn().mockResolvedValue(makePaginatedResponse([], 90, 3)),
        },
      };

      const { fetchAssets, pagination } = useAssets();

      pagination.value.currentPage = 2;
      await fetchAssets();

      expect(pagination.value.hasNext).toBe(true);
      expect(pagination.value.hasPrev).toBe(true);
    });

    it('handles API error gracefully', async () => {
      const { handleApiError } = await import('@utils/error');
      (apis as any).assets = {
        asset: {
          list: vi.fn().mockRejectedValue(new Error('Network error')),
        },
      };

      const { fetchAssets, isLoading } = useAssets();
      await fetchAssets();

      expect(isLoading.value).toBe(false);
      expect(handleApiError).toHaveBeenCalled();
    });
  });

  describe('pagination methods', () => {
    beforeEach(() => {
      (apis as any).assets = {
        asset: {
          list: vi.fn().mockResolvedValue(makePaginatedResponse([], 90, 3)),
        },
      };
    });

    it('goToPage sets currentPage and fetches', async () => {
      const { goToPage, pagination } = useAssets();
      await goToPage(3);
      expect(pagination.value.currentPage).toBe(3);
      expect(apis.assets.asset.list).toHaveBeenCalled();
    });

    it('loadMore increments page and appends when hasNext', async () => {
      const { loadMore, pagination } = useAssets();
      pagination.value.hasNext = true;
      pagination.value.currentPage = 1;

      await loadMore();
      expect(pagination.value.currentPage).toBe(2);
    });

    it('loadMore does nothing when no next page', async () => {
      const { loadMore, pagination } = useAssets();
      pagination.value.hasNext = false;
      pagination.value.currentPage = 1;

      await loadMore();
      expect(pagination.value.currentPage).toBe(1);
    });

    it('setItemsPerPage resets to page 1', async () => {
      const { setItemsPerPage, pagination } = useAssets();
      pagination.value.currentPage = 5;

      await setItemsPerPage(50);
      expect(pagination.value.perPage).toBe(50);
      expect(pagination.value.currentPage).toBe(1);
    });
  });

  describe('filter methods', () => {
    beforeEach(() => {
      (apis as any).assets = {
        asset: {
          list: vi.fn().mockResolvedValue(makePaginatedResponse([], 0, 1)),
        },
      };
    });

    it('applyFilters merges filters and resets page to 1', async () => {
      const { applyFilters, filters, pagination } = useAssets();
      pagination.value.currentPage = 3;

      await applyFilters({ name: 'sensor', status: true });
      expect(filters.value.name).toBe('sensor');
      expect(filters.value.status).toBe(true);
      expect(pagination.value.currentPage).toBe(1);
    });

    it('clearFilters resets all filters and dependent options', async () => {
      const { clearFilters, filters, manufacturerOptions, modelOptions } = useAssets();
      filters.value.name = 'test';
      filters.value.categoryId = 'cat-1';
      manufacturerOptions.value = [{ label: 'M', value: 'm-1' }];
      modelOptions.value = [{ label: 'X', value: 'x-1' }];

      await clearFilters();
      expect(filters.value.name).toBeUndefined();
      expect(filters.value.categoryId).toBeUndefined();
      expect(manufacturerOptions.value).toEqual([]);
      expect(modelOptions.value).toEqual([]);
    });
  });

  describe('cascading filters', () => {
    it('handleCategoryChange resets manufacturer and model', async () => {
      (apis as any).mapexOS = {
        lists: {
          list: vi.fn().mockResolvedValue({
            items: [{ name: 'Manufacturer A', id: 'mfr-1' }],
          }),
        },
      };

      const { handleCategoryChange, filters, manufacturerOptions, modelOptions } = useAssets();
      filters.value.manufacturerId = 'old-mfr';
      filters.value.modelId = 'old-model';

      await handleCategoryChange('cat-1');

      expect(filters.value.categoryId).toBe('cat-1');
      expect(filters.value.manufacturerId).toBeUndefined();
      expect(filters.value.modelId).toBeUndefined();
      expect(modelOptions.value).toEqual([]);
      expect(manufacturerOptions.value).toEqual([{ label: 'Manufacturer A', value: 'mfr-1' }]);
    });

    it('handleCategoryChange with undefined clears dependent options', async () => {
      const { handleCategoryChange, filters, manufacturerOptions } = useAssets();
      await handleCategoryChange(undefined);

      expect(filters.value.categoryId).toBeUndefined();
      expect(manufacturerOptions.value).toEqual([]);
    });

    it('handleManufacturerChange resets model', async () => {
      (apis as any).mapexOS = {
        lists: {
          list: vi.fn().mockResolvedValue({
            items: [{ name: 'Model X', id: 'mdl-1' }],
          }),
        },
      };

      const { handleManufacturerChange, filters, modelOptions } = useAssets();
      filters.value.modelId = 'old-model';

      await handleManufacturerChange('mfr-1');

      expect(filters.value.manufacturerId).toBe('mfr-1');
      expect(filters.value.modelId).toBeUndefined();
      expect(modelOptions.value).toEqual([{ label: 'Model X', value: 'mdl-1' }]);
    });

    it('loadCategories populates categoryOptions', async () => {
      (apis as any).mapexOS = {
        lists: {
          list: vi.fn().mockResolvedValue({
            items: [
              { name: 'Sensors', id: 'cat-1' },
              { name: 'Actuators', id: 'cat-2' },
            ],
          }),
        },
      };

      const { loadCategories, categoryOptions, loadingCategories } = useAssets();
      const promise = loadCategories();
      expect(loadingCategories.value).toBe(true);
      await promise;
      expect(loadingCategories.value).toBe(false);
      expect(categoryOptions.value).toHaveLength(2);
      expect(categoryOptions.value[0]).toEqual({ label: 'Sensors', value: 'cat-1' });
    });

    it('loadCategories handles missing API gracefully', async () => {
      (apis as any).mapexOS = undefined;

      const { loadCategories, categoryOptions } = useAssets();
      await loadCategories();
      expect(categoryOptions.value).toEqual([]);
    });

    it('loadManufacturers clears options when no categoryId', async () => {
      (apis as any).mapexOS = { lists: { list: vi.fn() } };

      const { loadManufacturers, manufacturerOptions } = useAssets();
      manufacturerOptions.value = [{ label: 'old', value: 'old' }];

      await loadManufacturers(undefined);
      expect(manufacturerOptions.value).toEqual([]);
    });
  });

  describe('reset', () => {
    it('resets all state to initial values', () => {
      const { reset, assets, filters, pagination, categoryOptions, manufacturerOptions, modelOptions } = useAssets();

      assets.value = [{ id: '1' }];
      filters.value.name = 'test';
      pagination.value.currentPage = 5;
      categoryOptions.value = [{ label: 'C', value: 'c-1' }];

      reset();

      expect(assets.value).toEqual([]);
      expect(filters.value.name).toBeUndefined();
      expect(pagination.value.currentPage).toBe(1);
      expect(pagination.value.perPage).toBe(30);
      expect(categoryOptions.value).toEqual([]);
      expect(manufacturerOptions.value).toEqual([]);
      expect(modelOptions.value).toEqual([]);
    });
  });
});

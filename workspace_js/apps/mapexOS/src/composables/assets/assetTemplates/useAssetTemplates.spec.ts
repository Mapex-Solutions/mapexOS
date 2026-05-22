import { describe, it, expect, beforeEach, vi } from 'vitest';
import { setActivePinia, createPinia } from 'pinia';
import { useAssetTemplates } from './useAssetTemplates';
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

describe('useAssetTemplates', () => {
  beforeEach(() => {
    setActivePinia(createPinia());
    vi.clearAllMocks();
  });

  describe('initial state', () => {
    it('returns empty templates and default pagination', () => {
      const { templates, isLoading, isLoadingMore, pagination, filters } = useAssetTemplates();

      expect(templates.value).toEqual([]);
      expect(isLoading.value).toBe(false);
      expect(isLoadingMore.value).toBe(false);
      expect(pagination.value.currentPage).toBe(1);
      expect(pagination.value.perPage).toBe(30);
      expect(pagination.value.totalItems).toBe(0);
      expect(filters.value.name).toBeUndefined();
      expect(filters.value.isSystem).toBeUndefined();
      expect(filters.value.isTemplate).toBeUndefined();
    });
  });

  describe('hasActiveFilters', () => {
    it('returns false when no filters set', () => {
      const { hasActiveFilters } = useAssetTemplates();
      expect(hasActiveFilters.value).toBe(false);
    });

    it('returns true when name filter set', () => {
      const { hasActiveFilters, filters } = useAssetTemplates();
      filters.value.name = 'template';
      expect(hasActiveFilters.value).toBe(true);
    });

    it('returns true when isSystem filter is boolean', () => {
      const { hasActiveFilters, filters } = useAssetTemplates();
      filters.value.isSystem = false;
      expect(hasActiveFilters.value).toBe(true);
    });

    it('returns true when isTemplate filter is boolean', () => {
      const { hasActiveFilters, filters } = useAssetTemplates();
      filters.value.isTemplate = true;
      expect(hasActiveFilters.value).toBe(true);
    });

    it('returns true when status filter is boolean', () => {
      const { hasActiveFilters, filters } = useAssetTemplates();
      filters.value.status = true;
      expect(hasActiveFilters.value).toBe(true);
    });
  });

  describe('fetchTemplates', () => {
    it('throws when assets API is not initialized', async () => {
      (apis as any).assets = undefined;
      const { fetchTemplates } = useAssetTemplates();
      await expect(fetchTemplates()).rejects.toThrow('Assets API not initialized');
    });

    it('sets isLoading during fetch and populates templates', async () => {
      const mockItems = [{ id: '1', name: 'Template A' }];
      (apis as any).assets = {
        assetTemplate: {
          list: vi.fn().mockResolvedValue(makePaginatedResponse(mockItems, 1, 1)),
        },
      };

      const { fetchTemplates, templates, isLoading, pagination } = useAssetTemplates();

      const promise = fetchTemplates();
      expect(isLoading.value).toBe(true);

      await promise;

      expect(isLoading.value).toBe(false);
      expect(templates.value).toEqual(mockItems);
      expect(pagination.value.totalItems).toBe(1);
    });

    it('sets isLoadingMore when appending', async () => {
      (apis as any).assets = {
        assetTemplate: {
          list: vi.fn().mockResolvedValue(makePaginatedResponse([], 0, 1)),
        },
      };

      const { fetchTemplates, isLoadingMore } = useAssetTemplates();
      const promise = fetchTemplates(true);
      expect(isLoadingMore.value).toBe(true);
      await promise;
      expect(isLoadingMore.value).toBe(false);
    });

    it('computes hasNext and hasPrev correctly', async () => {
      (apis as any).assets = {
        assetTemplate: {
          list: vi.fn().mockResolvedValue(makePaginatedResponse([], 60, 3)),
        },
      };

      const { fetchTemplates, pagination } = useAssetTemplates();
      pagination.value.currentPage = 2;
      await fetchTemplates();

      expect(pagination.value.hasNext).toBe(true);
      expect(pagination.value.hasPrev).toBe(true);
    });

    it('handles API error gracefully', async () => {
      const { handleApiError } = await import('@utils/error');
      (apis as any).assets = {
        assetTemplate: {
          list: vi.fn().mockRejectedValue(new Error('fail')),
        },
      };

      const { fetchTemplates, isLoading } = useAssetTemplates();
      await fetchTemplates();

      expect(isLoading.value).toBe(false);
      expect(handleApiError).toHaveBeenCalled();
    });
  });

  describe('pagination methods', () => {
    beforeEach(() => {
      (apis as any).assets = {
        assetTemplate: {
          list: vi.fn().mockResolvedValue(makePaginatedResponse([], 60, 3)),
        },
      };
    });

    it('goToPage sets currentPage and fetches', async () => {
      const { goToPage, pagination } = useAssetTemplates();
      await goToPage(2);
      expect(pagination.value.currentPage).toBe(2);
    });

    it('loadMore increments page when hasNext', async () => {
      const { loadMore, pagination } = useAssetTemplates();
      pagination.value.hasNext = true;
      pagination.value.currentPage = 1;

      await loadMore();
      expect(pagination.value.currentPage).toBe(2);
    });

    it('loadMore does nothing when no next page', async () => {
      const { loadMore, pagination } = useAssetTemplates();
      pagination.value.hasNext = false;

      await loadMore();
      expect(pagination.value.currentPage).toBe(1);
    });

    it('setItemsPerPage resets to page 1', async () => {
      const { setItemsPerPage, pagination } = useAssetTemplates();
      pagination.value.currentPage = 5;

      await setItemsPerPage(50);
      expect(pagination.value.perPage).toBe(50);
      expect(pagination.value.currentPage).toBe(1);
    });
  });

  describe('filter methods', () => {
    beforeEach(() => {
      (apis as any).assets = {
        assetTemplate: {
          list: vi.fn().mockResolvedValue(makePaginatedResponse([], 0, 1)),
        },
      };
    });

    it('applyFilters merges filters and resets page to 1', async () => {
      const { applyFilters, filters, pagination } = useAssetTemplates();
      pagination.value.currentPage = 3;

      await applyFilters({ name: 'temp', isSystem: true });
      expect(filters.value.name).toBe('temp');
      expect(filters.value.isSystem).toBe(true);
      expect(pagination.value.currentPage).toBe(1);
    });

    it('clearFilters resets all filters and cascading options', async () => {
      const { clearFilters, filters, manufacturerOptions, modelOptions } = useAssetTemplates();
      filters.value.name = 'test';
      filters.value.categoryId = 'cat-1';
      manufacturerOptions.value = [{ label: 'M', value: 'm-1' }];

      await clearFilters();
      expect(filters.value.name).toBeUndefined();
      expect(filters.value.categoryId).toBeUndefined();
      expect(filters.value.isSystem).toBeUndefined();
      expect(filters.value.isTemplate).toBeUndefined();
      expect(manufacturerOptions.value).toEqual([]);
      expect(modelOptions.value).toEqual([]);
    });
  });

  describe('cascading filters', () => {
    it('handleCategoryChange resets manufacturer and model', async () => {
      (apis as any).mapexOS = {
        lists: {
          list: vi.fn().mockResolvedValue({
            items: [{ name: 'Mfr A', id: 'mfr-1' }],
          }),
        },
      };

      const { handleCategoryChange, filters, manufacturerOptions, modelOptions } = useAssetTemplates();
      filters.value.manufacturerId = 'old-mfr';
      filters.value.modelId = 'old-model';

      await handleCategoryChange('cat-1');

      expect(filters.value.categoryId).toBe('cat-1');
      expect(filters.value.manufacturerId).toBeUndefined();
      expect(filters.value.modelId).toBeUndefined();
      expect(modelOptions.value).toEqual([]);
      expect(manufacturerOptions.value).toEqual([{ label: 'Mfr A', value: 'mfr-1' }]);
    });

    it('handleManufacturerChange resets model', async () => {
      (apis as any).mapexOS = {
        lists: {
          list: vi.fn().mockResolvedValue({
            items: [{ name: 'Model X', id: 'mdl-1' }],
          }),
        },
      };

      const { handleManufacturerChange, filters, modelOptions } = useAssetTemplates();
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
            items: [{ name: 'Sensors', id: 'cat-1' }],
          }),
        },
      };

      const { loadCategories, categoryOptions, loadingCategories } = useAssetTemplates();
      const promise = loadCategories();
      expect(loadingCategories.value).toBe(true);
      await promise;
      expect(loadingCategories.value).toBe(false);
      expect(categoryOptions.value).toEqual([{ label: 'Sensors', value: 'cat-1' }]);
    });
  });

  describe('reset', () => {
    it('resets all state to initial values', () => {
      const { reset, templates, filters, pagination, categoryOptions } = useAssetTemplates();

      templates.value = [{ id: '1' }];
      filters.value.name = 'test';
      pagination.value.currentPage = 5;
      categoryOptions.value = [{ label: 'C', value: 'c-1' }];

      reset();

      expect(templates.value).toEqual([]);
      expect(filters.value.name).toBeUndefined();
      expect(pagination.value.currentPage).toBe(1);
      expect(categoryOptions.value).toEqual([]);
    });
  });
});

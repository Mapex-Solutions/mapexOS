import { describe, it, expect, beforeEach, vi } from 'vitest';
import { setActivePinia, createPinia } from 'pinia';
import { useRouters } from './useRouters';
import { apis } from '@services/mapex';

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

describe('useRouters', () => {
  beforeEach(() => {
    setActivePinia(createPinia());
    vi.clearAllMocks();
  });

  describe('initial state', () => {
    it('returns empty routeGroups and default pagination', () => {
      const { routeGroups, isLoading, isLoadingMore, pagination, filters } = useRouters();

      expect(routeGroups.value).toEqual([]);
      expect(isLoading.value).toBe(false);
      expect(isLoadingMore.value).toBe(false);
      expect(pagination.value.currentPage).toBe(1);
      expect(pagination.value.perPage).toBe(30);
      expect(pagination.value.totalItems).toBe(0);
      expect(filters.value.name).toBeUndefined();
      expect(filters.value.status).toBeUndefined();
      expect(filters.value.isSystem).toBeUndefined();
      expect(filters.value.isTemplate).toBeUndefined();
      expect(filters.value.includeChildren).toBeUndefined();
    });
  });

  describe('hasActiveFilters', () => {
    it('returns false when no filters set', () => {
      const { hasActiveFilters } = useRouters();
      expect(hasActiveFilters.value).toBe(false);
    });

    it('returns true when name filter set', () => {
      const { hasActiveFilters, filters } = useRouters();
      filters.value.name = 'route';
      expect(hasActiveFilters.value).toBe(true);
    });

    it('returns true when status filter is boolean false', () => {
      const { hasActiveFilters, filters } = useRouters();
      filters.value.status = false;
      expect(hasActiveFilters.value).toBe(true);
    });

    it('returns true when isSystem filter is boolean', () => {
      const { hasActiveFilters, filters } = useRouters();
      filters.value.isSystem = true;
      expect(hasActiveFilters.value).toBe(true);
    });

    it('returns true when isTemplate filter is boolean', () => {
      const { hasActiveFilters, filters } = useRouters();
      filters.value.isTemplate = false;
      expect(hasActiveFilters.value).toBe(true);
    });

    it('returns true when includeChildren is boolean', () => {
      const { hasActiveFilters, filters } = useRouters();
      filters.value.includeChildren = true;
      expect(hasActiveFilters.value).toBe(true);
    });
  });

  describe('fetchRouteGroups', () => {
    it('throws when router API is not initialized', async () => {
      (apis as any).router = undefined;
      const { fetchRouteGroups } = useRouters();
      await expect(fetchRouteGroups()).rejects.toThrow('Router API not initialized');
    });

    it('sets isLoading during fetch and populates routeGroups', async () => {
      const mockItems = [{ id: '1', name: 'Route A' }];
      (apis as any).router = {
        routegroup: {
          list: vi.fn().mockResolvedValue(makePaginatedResponse(mockItems, 1, 1)),
        },
      };

      const { fetchRouteGroups, routeGroups, isLoading, pagination } = useRouters();

      const promise = fetchRouteGroups();
      expect(isLoading.value).toBe(true);

      await promise;

      expect(isLoading.value).toBe(false);
      expect(routeGroups.value).toEqual(mockItems);
      expect(pagination.value.totalItems).toBe(1);
    });

    it('appends results when append=true', async () => {
      const page1 = [{ id: '1', name: 'A' }];
      const page2 = [{ id: '2', name: 'B' }];

      const listFn = vi
        .fn()
        .mockResolvedValueOnce(makePaginatedResponse(page1, 2, 2))
        .mockResolvedValueOnce(makePaginatedResponse(page2, 2, 2));

      (apis as any).router = { routegroup: { list: listFn } };

      const { fetchRouteGroups, routeGroups, pagination } = useRouters();

      await fetchRouteGroups();
      expect(routeGroups.value).toHaveLength(1);
      expect(routeGroups.value[0]!.id).toBe('1');

      pagination.value.currentPage = 2;
      await fetchRouteGroups(true);
      expect(routeGroups.value).toHaveLength(2);
      expect(routeGroups.value[0]!.id).toBe('1');
      expect(routeGroups.value[1]!.id).toBe('2');
    });

    it('sets isLoadingMore when appending', async () => {
      (apis as any).router = {
        routegroup: {
          list: vi.fn().mockResolvedValue(makePaginatedResponse([], 0, 1)),
        },
      };

      const { fetchRouteGroups, isLoadingMore } = useRouters();
      const promise = fetchRouteGroups(true);
      expect(isLoadingMore.value).toBe(true);
      await promise;
      expect(isLoadingMore.value).toBe(false);
    });

    it('computes hasNext and hasPrev correctly', async () => {
      (apis as any).router = {
        routegroup: {
          list: vi.fn().mockResolvedValue(makePaginatedResponse([], 90, 3)),
        },
      };

      const { fetchRouteGroups, pagination } = useRouters();
      pagination.value.currentPage = 2;
      await fetchRouteGroups();

      expect(pagination.value.hasNext).toBe(true);
      expect(pagination.value.hasPrev).toBe(true);
    });

    it('handles API error gracefully', async () => {
      const { handleApiError } = await import('@utils/error');
      (apis as any).router = {
        routegroup: {
          list: vi.fn().mockRejectedValue(new Error('Network error')),
        },
      };

      const { fetchRouteGroups, isLoading } = useRouters();
      await fetchRouteGroups();

      expect(isLoading.value).toBe(false);
      expect(handleApiError).toHaveBeenCalled();
    });
  });

  describe('pagination methods', () => {
    beforeEach(() => {
      (apis as any).router = {
        routegroup: {
          list: vi.fn().mockResolvedValue(makePaginatedResponse([], 90, 3)),
        },
      };
    });

    it('goToPage sets currentPage and fetches', async () => {
      const { goToPage, pagination } = useRouters();
      await goToPage(3);
      expect(pagination.value.currentPage).toBe(3);
    });

    it('loadMore increments page when hasNext', async () => {
      const { loadMore, pagination } = useRouters();
      pagination.value.hasNext = true;
      pagination.value.currentPage = 1;

      await loadMore();
      expect(pagination.value.currentPage).toBe(2);
    });

    it('loadMore does nothing when no next page', async () => {
      const { loadMore, pagination } = useRouters();
      pagination.value.hasNext = false;
      pagination.value.currentPage = 1;

      await loadMore();
      expect(pagination.value.currentPage).toBe(1);
    });

    it('setItemsPerPage resets to page 1', async () => {
      const { setItemsPerPage, pagination } = useRouters();
      pagination.value.currentPage = 5;

      await setItemsPerPage(50);
      expect(pagination.value.perPage).toBe(50);
      expect(pagination.value.currentPage).toBe(1);
    });
  });

  describe('filter methods', () => {
    beforeEach(() => {
      (apis as any).router = {
        routegroup: {
          list: vi.fn().mockResolvedValue(makePaginatedResponse([], 0, 1)),
        },
      };
    });

    it('applyFilters merges filters and resets page to 1', async () => {
      const { applyFilters, filters, pagination } = useRouters();
      pagination.value.currentPage = 3;

      await applyFilters({ name: 'route', status: true });
      expect(filters.value.name).toBe('route');
      expect(filters.value.status).toBe(true);
      expect(pagination.value.currentPage).toBe(1);
    });

    it('clearFilters resets all filters', async () => {
      const { clearFilters, filters } = useRouters();
      filters.value.name = 'test';
      filters.value.status = true;

      await clearFilters();
      expect(filters.value.name).toBeUndefined();
      expect(filters.value.status).toBeUndefined();
      expect(filters.value.includeChildren).toBeUndefined();
    });
  });

  describe('reset', () => {
    it('resets all state to initial values', () => {
      const { reset, routeGroups, filters, pagination } = useRouters();

      routeGroups.value = [{ id: '1' }];
      filters.value.name = 'test';
      pagination.value.currentPage = 5;

      reset();

      expect(routeGroups.value).toEqual([]);
      expect(filters.value.name).toBeUndefined();
      expect(filters.value.isSystem).toBeUndefined();
      expect(filters.value.isTemplate).toBeUndefined();
      expect(pagination.value.currentPage).toBe(1);
      expect(pagination.value.perPage).toBe(30);
      expect(pagination.value.totalItems).toBe(0);
    });
  });
});

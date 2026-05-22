import { describe, it, expect, vi, beforeEach } from 'vitest';
import { ref } from 'vue';
import type {
  RouteGroupsListPageFilters,
  RouteGroupsListPageColumnVisibility,
  EnrichedRouteGroup,
} from '../interfaces';
import {
  fetchRouteGroupsHandler,
  handleFilterApplyHandler,
  handlePageChangeHandler,
  handleItemsPerPageChangeHandler,
  handleColumnsUpdateHandler,
  canModifyRouteGroupHandler,
  editRouteGroupHandler,
  deleteRouteGroupHandler,
} from './routeGroupsListPage.handler';
import { apis } from '@services/mapex';

vi.mock('@utils/alert', () => ({
  notifySuccess: vi.fn(),
  notifyFail: vi.fn(),
  notifyWarning: vi.fn(),
  dialogDelete: vi.fn(),
}));
vi.mock('@utils/query', () => ({
  cleanQueryParams: vi.fn((params: any) => params),
}));
vi.mock('@stores/organization', () => ({
  useOrganizationStore: () => ({
    flatList: [
      { id: 'org-1', name: 'Org One' },
      { id: 'org-2', name: 'Org Two' },
    ],
    selectedOrganizationId: 'org-1',
  }),
}));
vi.mock('@composables/useLogger', () => ({
  useLogger: () => ({
    warn: vi.fn(),
    error: vi.fn(),
    debug: vi.fn(),
  }),
}));

function makeFilters(overrides: Partial<RouteGroupsListPageFilters> = {}): RouteGroupsListPageFilters {
  return {
    name: undefined,
    enabled: undefined,
    isTemplate: undefined,
    includeChildren: undefined,
    ...overrides,
  };
}

function makeColumnVisibility(overrides: Partial<RouteGroupsListPageColumnVisibility> = {}): RouteGroupsListPageColumnVisibility {
  return {
    organization: true,
    routers: true,
    isTemplate: true,
    ...overrides,
  };
}

describe('routeGroupsListPage.handler', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    // Ensure list, delete methods exist
    (apis.router.routegroup as any).list = vi.fn();
    (apis.router.routegroup as any).delete = vi.fn();
  });

  // ── fetchRouteGroupsHandler ──────────────────────────────────────────
  describe('fetchRouteGroupsHandler', () => {
    it('should fetch route groups and enrich with organization name', async () => {
      vi.mocked(apis.router.routegroup.list).mockResolvedValue({
        items: [
          { id: 'rg-1', name: 'RG1', orgId: 'org-1', routers: [{}] },
          { id: 'rg-2', name: 'RG2', orgId: 'org-2', routers: [] },
        ],
        pagination: { totalItems: 2, totalPages: 1 },
      } as any);

      const filters = ref(makeFilters());
      const currentPage = ref(1);
      const itemsPerPage = ref(10);
      const routeGroupsList = ref<EnrichedRouteGroup[]>([]);
      const totalPages = ref(0);
      const totalItems = ref(0);
      const hasNext = ref(false);
      const hasPrev = ref(false);
      const loading = ref(false);
      const error = ref<string | undefined>(undefined);

      await fetchRouteGroupsHandler(
        filters, currentPage, itemsPerPage, routeGroupsList,
        totalPages, totalItems, hasNext, hasPrev, loading, error,
      );

      expect(loading.value).toBe(false);
      expect(routeGroupsList.value).toHaveLength(2);
      expect(routeGroupsList.value[0]!.organizationName).toBe('Org One');
      expect(routeGroupsList.value[0]!.routersCount).toBe(1);
      expect(routeGroupsList.value[1]!.organizationName).toBe('Org Two');
      expect(totalItems.value).toBe(2);
    });

    it('should set error when API is not initialized', async () => {
      // Temporarily remove the API
      const origRouter = apis.router;
      (apis as any).router = undefined;

      const filters = ref(makeFilters());
      const currentPage = ref(1);
      const itemsPerPage = ref(10);
      const routeGroupsList = ref<EnrichedRouteGroup[]>([]);
      const totalPages = ref(0);
      const totalItems = ref(0);
      const hasNext = ref(false);
      const hasPrev = ref(false);
      const loading = ref(false);
      const error = ref<string | undefined>(undefined);

      await fetchRouteGroupsHandler(
        filters, currentPage, itemsPerPage, routeGroupsList,
        totalPages, totalItems, hasNext, hasPrev, loading, error,
      );

      expect(error.value).toBe('Router API not initialized');

      // Restore
      (apis as any).router = origRouter;
    });

    it('should include filters in query params when set', async () => {
      vi.mocked(apis.router.routegroup.list).mockResolvedValue({
        items: [],
        pagination: { totalItems: 0, totalPages: 1 },
      } as any);

      const filters = ref(makeFilters({ name: 'test', enabled: true, isTemplate: false }));
      const currentPage = ref(1);
      const itemsPerPage = ref(10);
      const routeGroupsList = ref<EnrichedRouteGroup[]>([]);
      const totalPages = ref(0);
      const totalItems = ref(0);
      const hasNext = ref(false);
      const hasPrev = ref(false);
      const loading = ref(false);
      const error = ref<string | undefined>(undefined);

      await fetchRouteGroupsHandler(
        filters, currentPage, itemsPerPage, routeGroupsList,
        totalPages, totalItems, hasNext, hasPrev, loading, error,
      );

      const callArgs = vi.mocked(apis.router.routegroup.list).mock.calls[0]![0] as any;
      expect(callArgs.name).toBe('test');
      expect(callArgs.enabled).toBe(true);
      expect(callArgs.isTemplate).toBe(false);
    });

    it('should handle API errors and set error state', async () => {
      vi.mocked(apis.router.routegroup.list).mockRejectedValue(new Error('Server error'));

      const filters = ref(makeFilters());
      const currentPage = ref(1);
      const itemsPerPage = ref(10);
      const routeGroupsList = ref<EnrichedRouteGroup[]>([]);
      const totalPages = ref(0);
      const totalItems = ref(0);
      const hasNext = ref(false);
      const hasPrev = ref(false);
      const loading = ref(false);
      const error = ref<string | undefined>(undefined);

      await fetchRouteGroupsHandler(
        filters, currentPage, itemsPerPage, routeGroupsList,
        totalPages, totalItems, hasNext, hasPrev, loading, error,
      );

      expect(error.value).toBe('Server error');
      expect(loading.value).toBe(false);
    });
  });

  // ── handleFilterApplyHandler ─────────────────────────────────────────
  describe('handleFilterApplyHandler', () => {
    it('should update filters, reset page, and call fetch', () => {
      const filters = ref(makeFilters());
      const columnVisibility = ref(makeColumnVisibility());
      const currentPage = ref(3);
      const fetchCallback = vi.fn();

      handleFilterApplyHandler(
        { name: 'test', enabled: true, isTemplate: false, includeChildren: undefined },
        filters,
        columnVisibility,
        currentPage,
        fetchCallback,
      );

      expect(filters.value.name).toBe('test');
      expect(filters.value.enabled).toBe(true);
      expect(currentPage.value).toBe(1);
      expect(fetchCallback).toHaveBeenCalledOnce();
    });

    it('should hide columns when includeChildren is true', () => {
      const filters = ref(makeFilters());
      const columnVisibility = ref(makeColumnVisibility());
      const currentPage = ref(1);
      const fetchCallback = vi.fn();

      handleFilterApplyHandler(
        { includeChildren: true },
        filters,
        columnVisibility,
        currentPage,
        fetchCallback,
      );

      expect(columnVisibility.value.routers).toBe(false);
      expect(columnVisibility.value.isTemplate).toBe(false);
    });

    it('should restore columns when includeChildren is not true', () => {
      const filters = ref(makeFilters());
      const columnVisibility = ref(makeColumnVisibility({ routers: false, isTemplate: false }));
      const currentPage = ref(1);
      const fetchCallback = vi.fn();

      handleFilterApplyHandler(
        { includeChildren: false },
        filters,
        columnVisibility,
        currentPage,
        fetchCallback,
      );

      expect(columnVisibility.value.routers).toBe(true);
      expect(columnVisibility.value.isTemplate).toBe(true);
    });
  });

  // ── handlePageChangeHandler ──────────────────────────────────────────
  describe('handlePageChangeHandler', () => {
    it('should set current page and call fetch', () => {
      const currentPage = ref(1);
      const fetchCallback = vi.fn();

      handlePageChangeHandler(5, currentPage, fetchCallback);

      expect(currentPage.value).toBe(5);
      expect(fetchCallback).toHaveBeenCalledOnce();
    });
  });

  // ── handleItemsPerPageChangeHandler ──────────────────────────────────
  describe('handleItemsPerPageChangeHandler', () => {
    it('should set items per page, reset to page 1, and call fetch', () => {
      const itemsPerPage = ref(10);
      const currentPage = ref(3);
      const fetchCallback = vi.fn();

      handleItemsPerPageChangeHandler(25, itemsPerPage, currentPage, fetchCallback);

      expect(itemsPerPage.value).toBe(25);
      expect(currentPage.value).toBe(1);
      expect(fetchCallback).toHaveBeenCalledOnce();
    });
  });

  // ── handleColumnsUpdateHandler ───────────────────────────────────────
  describe('handleColumnsUpdateHandler', () => {
    it('should update column visibility based on columns array', () => {
      const columnVisibility = ref(makeColumnVisibility());

      handleColumnsUpdateHandler(
        [
          { key: 'organization', visible: false, label: '' },
          { key: 'routers', visible: false, label: '' },
          { key: 'isTemplate', visible: true, label: '' },
        ],
        columnVisibility,
      );

      expect(columnVisibility.value.organization).toBe(false);
      expect(columnVisibility.value.routers).toBe(false);
      expect(columnVisibility.value.isTemplate).toBe(true);
    });
  });

  // ── canModifyRouteGroupHandler ───────────────────────────────────────
  describe('canModifyRouteGroupHandler', () => {
    it('should return true for non-template route groups', () => {
      const rg: EnrichedRouteGroup = { id: 'rg-1', isTemplate: false, orgId: 'org-other' };
      expect(canModifyRouteGroupHandler(rg, 'org-1')).toBe(true);
    });

    it('should return true for templates owned by current org', () => {
      const rg: EnrichedRouteGroup = { id: 'rg-1', isTemplate: true, orgId: 'org-1' };
      expect(canModifyRouteGroupHandler(rg, 'org-1')).toBe(true);
    });

    it('should return false for templates not owned by current org', () => {
      const rg: EnrichedRouteGroup = { id: 'rg-1', isTemplate: true, orgId: 'org-other' };
      expect(canModifyRouteGroupHandler(rg, 'org-1')).toBe(false);
    });
  });

  // ── editRouteGroupHandler ────────────────────────────────────────────
  describe('editRouteGroupHandler', () => {
    it('should navigate to edit page for modifiable route group', () => {
      const router = { push: vi.fn() };
      const rg: EnrichedRouteGroup = { id: 'rg-1', isTemplate: false };

      editRouteGroupHandler(rg, 'org-1', router, { notifications: { sharedEdit: { value: '' } } });

      expect(router.push).toHaveBeenCalledWith('/routing/route_groups/edit/rg-1');
    });

    it('should warn when trying to edit a shared template from another org', () => {
      const router = { push: vi.fn() };
      const rg: EnrichedRouteGroup = { id: 'rg-1', isTemplate: true, orgId: 'org-other' };

      editRouteGroupHandler(rg, 'org-1', router, { notifications: { sharedEdit: { value: 'Cannot edit' } } });

      expect(router.push).not.toHaveBeenCalled();
    });
  });

  // ── deleteRouteGroupHandler ──────────────────────────────────────────
  describe('deleteRouteGroupHandler', () => {
    it('should delete and remove from local list', async () => {
      (apis.router.routegroup as any).delete = vi.fn().mockResolvedValue({});

      const routeGroupsList = ref<EnrichedRouteGroup[]>([
        { id: 'rg-1', name: 'RG1' },
        { id: 'rg-2', name: 'RG2' },
      ]);
      const totalItems = ref(2);
      const currentPage = ref(1);
      const mockT = { notifications: { deleted: { value: 'Deleted' } } };
      const fetchCallback = vi.fn();

      await deleteRouteGroupHandler(
        { id: 'rg-1', name: 'RG1' },
        routeGroupsList,
        totalItems,
        currentPage,
        mockT,
        fetchCallback,
      );

      expect(routeGroupsList.value).toHaveLength(1);
      expect(routeGroupsList.value[0]!.id).toBe('rg-2');
      expect(totalItems.value).toBe(1);
    });

    it('should go to previous page when current page becomes empty', async () => {
      (apis.router.routegroup as any).delete = vi.fn().mockResolvedValue({});

      const routeGroupsList = ref<EnrichedRouteGroup[]>([{ id: 'rg-1', name: 'RG1' }]);
      const totalItems = ref(1);
      const currentPage = ref(2);
      const mockT = { notifications: { deleted: { value: 'Deleted' } } };
      const fetchCallback = vi.fn().mockResolvedValue(undefined);

      await deleteRouteGroupHandler(
        { id: 'rg-1', name: 'RG1' },
        routeGroupsList,
        totalItems,
        currentPage,
        mockT,
        fetchCallback,
      );

      expect(currentPage.value).toBe(1);
      expect(fetchCallback).toHaveBeenCalledOnce();
    });
  });
});

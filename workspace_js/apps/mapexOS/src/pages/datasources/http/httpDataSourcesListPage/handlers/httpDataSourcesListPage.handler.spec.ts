import { describe, it, expect, vi, beforeEach } from 'vitest';
import { ref } from 'vue';
import type {
  HttpDataSourcesListPageFilters,
  HttpDataSourcesListPageColumnVisibility,
  EnrichedDataSource,
} from '../interfaces';
import {
  fetchDataSourcesHandler,
  handleFilterApplyHandler,
  handlePageChangeHandler,
  handleItemsPerPageChangeHandler,
  handleColumnsUpdateHandler,
  viewDetailsHandler,
  editDataSourceHandler,
  deleteDataSourceHandler,
} from './httpDataSourcesListPage.handler';
import { apis } from '@services/mapex';

vi.mock('@utils/alert', () => ({
  notifySuccess: vi.fn(),
  notifyFail: vi.fn(),
  dialogDelete: vi.fn(),
}));
vi.mock('@stores/organization', () => ({
  useOrganizationStore: () => ({
    flatList: [
      { id: 'org-1', name: 'Org One' },
      { id: 'org-2', name: 'Org Two' },
    ],
  }),
}));
vi.mock('@composables/useLogger', () => ({
  useLogger: () => ({
    warn: vi.fn(),
    error: vi.fn(),
    debug: vi.fn(),
  }),
}));

function makeFilters(overrides: Partial<HttpDataSourcesListPageFilters> = {}): HttpDataSourcesListPageFilters {
  return {
    name: undefined,
    mode: undefined,
    enabled: undefined,
    auth: undefined,
    assetBind: undefined,
    includeChildren: undefined,
    ...overrides,
  };
}

function makeColumnVisibility(overrides: Partial<HttpDataSourcesListPageColumnVisibility> = {}): HttpDataSourcesListPageColumnVisibility {
  return {
    organization: true,
    assetBind: true,
    auth: true,
    mode: true,
    ...overrides,
  };
}

describe('httpDataSourcesListPage.handler', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    // Ensure httpGateway API mock exists
    (apis as any).httpGateway = {
      datasource: {
        list: vi.fn(),
        delete: vi.fn(),
      },
    };
  });

  // ── fetchDataSourcesHandler ──────────────────────────────────────────
  describe('fetchDataSourcesHandler', () => {
    it('should fetch data sources and enrich with organization name', async () => {
      (apis.httpGateway as any).datasource.list.mockResolvedValue({
        items: [
          { id: 'ds-1', name: 'DS1', orgId: 'org-1' },
          { id: 'ds-2', name: 'DS2', orgId: 'org-2' },
        ],
        pagination: { totalPages: 1, totalItems: 2 },
      });

      const filters = ref(makeFilters());
      const currentPage = ref(1);
      const itemsPerPage = ref(10);
      const dataSourcesList = ref<EnrichedDataSource[]>([]);
      const totalPages = ref(0);
      const totalItems = ref(0);
      const loading = ref(false);
      const error = ref<string | undefined>(undefined);

      await fetchDataSourcesHandler(
        filters, currentPage, itemsPerPage, dataSourcesList,
        totalPages, totalItems, loading, error,
      );

      expect(loading.value).toBe(false);
      expect(dataSourcesList.value).toHaveLength(2);
      expect(dataSourcesList.value[0]!.organizationName).toBe('Org One');
      expect(totalItems.value).toBe(2);
    });

    it('should set error when API is not initialized', async () => {
      (apis as any).httpGateway = undefined;

      const filters = ref(makeFilters());
      const currentPage = ref(1);
      const itemsPerPage = ref(10);
      const dataSourcesList = ref<EnrichedDataSource[]>([]);
      const totalPages = ref(0);
      const totalItems = ref(0);
      const loading = ref(false);
      const error = ref<string | undefined>(undefined);

      await fetchDataSourcesHandler(
        filters, currentPage, itemsPerPage, dataSourcesList,
        totalPages, totalItems, loading, error,
      );

      expect(error.value).toBe('HTTP Gateway API not initialized');
    });

    it('should include all filter params when set', async () => {
      (apis as any).httpGateway = {
        datasource: {
          list: vi.fn().mockResolvedValue({ items: [], pagination: { totalPages: 1, totalItems: 0 } }),
        },
      };

      const filters = ref(makeFilters({
        name: 'test',
        enabled: true,
        mode: 'push',
        auth: 'bearer',
        assetBind: 'asset-1',
        includeChildren: true,
      }));
      const currentPage = ref(1);
      const itemsPerPage = ref(10);
      const dataSourcesList = ref<EnrichedDataSource[]>([]);
      const totalPages = ref(0);
      const totalItems = ref(0);
      const loading = ref(false);
      const error = ref<string | undefined>(undefined);

      await fetchDataSourcesHandler(
        filters, currentPage, itemsPerPage, dataSourcesList,
        totalPages, totalItems, loading, error,
      );

      const callArgs = (apis.httpGateway as any).datasource.list.mock.calls[0][0];
      expect(callArgs.name).toBe('test');
      expect(callArgs.enabled).toBe(true);
      expect(callArgs.mode).toBe('push');
      expect(callArgs.auth).toBe('bearer');
      expect(callArgs.assetBind).toBe('asset-1');
      expect(callArgs.includeChildren).toBe(true);
    });

    it('should handle API errors and set error state', async () => {
      (apis.httpGateway as any).datasource.list.mockRejectedValue(new Error('Network error'));

      const filters = ref(makeFilters());
      const currentPage = ref(1);
      const itemsPerPage = ref(10);
      const dataSourcesList = ref<EnrichedDataSource[]>([]);
      const totalPages = ref(0);
      const totalItems = ref(0);
      const loading = ref(false);
      const error = ref<string | undefined>(undefined);

      await fetchDataSourcesHandler(
        filters, currentPage, itemsPerPage, dataSourcesList,
        totalPages, totalItems, loading, error,
      );

      expect(error.value).toBe('Network error');
      expect(loading.value).toBe(false);
    });
  });

  // ── handleFilterApplyHandler ─────────────────────────────────────────
  describe('handleFilterApplyHandler', () => {
    it('should update filters, reset page, and call fetch', () => {
      const filters = ref(makeFilters());
      const currentPage = ref(3);
      const fetchCallback = vi.fn();

      handleFilterApplyHandler(
        { name: 'test', enabled: true, mode: 'pull' },
        filters,
        currentPage,
        fetchCallback,
      );

      expect(filters.value.name).toBe('test');
      expect(filters.value.enabled).toBe(true);
      expect(filters.value.mode).toBe('pull');
      expect(currentPage.value).toBe(1);
      expect(fetchCallback).toHaveBeenCalledOnce();
    });
  });

  // ── handlePageChangeHandler ──────────────────────────────────────────
  describe('handlePageChangeHandler', () => {
    it('should set current page and call fetch', () => {
      const currentPage = ref(1);
      const fetchCallback = vi.fn();

      handlePageChangeHandler(4, currentPage, fetchCallback);

      expect(currentPage.value).toBe(4);
      expect(fetchCallback).toHaveBeenCalledOnce();
    });
  });

  // ── handleItemsPerPageChangeHandler ──────────────────────────────────
  describe('handleItemsPerPageChangeHandler', () => {
    it('should set items per page, reset to page 1, and call fetch', () => {
      const itemsPerPage = ref(10);
      const currentPage = ref(5);
      const fetchCallback = vi.fn();

      handleItemsPerPageChangeHandler(50, itemsPerPage, currentPage, fetchCallback);

      expect(itemsPerPage.value).toBe(50);
      expect(currentPage.value).toBe(1);
      expect(fetchCallback).toHaveBeenCalledOnce();
    });
  });

  // ── handleColumnsUpdateHandler ───────────────────────────────────────
  describe('handleColumnsUpdateHandler', () => {
    it('should update column visibility from columns array', () => {
      const columnVisibility = ref(makeColumnVisibility());

      handleColumnsUpdateHandler(
        [
          { key: 'organization', visible: false, label: '' },
          { key: 'assetBind', visible: false, label: '' },
          { key: 'auth', visible: true, label: '' },
          { key: 'mode', visible: false, label: '' },
        ],
        columnVisibility,
      );

      expect(columnVisibility.value.organization).toBe(false);
      expect(columnVisibility.value.assetBind).toBe(false);
      expect(columnVisibility.value.auth).toBe(true);
      expect(columnVisibility.value.mode).toBe(false);
    });
  });

  // ── viewDetailsHandler ───────────────────────────────────────────────
  describe('viewDetailsHandler', () => {
    it('should set selected data source and open drawer', () => {
      const selectedId = ref<string | undefined>(undefined);
      const drawerOpen = ref(false);

      viewDetailsHandler({ id: 'ds-1' }, selectedId, drawerOpen);

      expect(selectedId.value).toBe('ds-1');
      expect(drawerOpen.value).toBe(true);
    });

    it('should not open drawer when data source has no id', () => {
      const selectedId = ref<string | undefined>(undefined);
      const drawerOpen = ref(false);

      viewDetailsHandler({}, selectedId, drawerOpen);

      expect(selectedId.value).toBeUndefined();
      expect(drawerOpen.value).toBe(false);
    });
  });

  // ── editDataSourceHandler ────────────────────────────────────────────
  describe('editDataSourceHandler', () => {
    it('should navigate to edit page', () => {
      const router = { push: vi.fn() };

      editDataSourceHandler({ id: 'ds-1' }, router);

      expect(router.push).toHaveBeenCalledWith('/data_sources/http/edit/ds-1');
    });
  });

  // ── deleteDataSourceHandler ──────────────────────────────────────────
  describe('deleteDataSourceHandler', () => {
    it('should delete and remove from local list', async () => {
      (apis.httpGateway as any).datasource.delete.mockResolvedValue({});

      const dataSourcesList = ref<EnrichedDataSource[]>([
        { id: 'ds-1', name: 'DS1' },
        { id: 'ds-2', name: 'DS2' },
      ]);
      const totalItems = ref(2);
      const currentPage = ref(1);
      const mockT = {
        notifications: {
          deleteSuccess: { value: 'Deleted' },
          deleteFailed: { value: 'Failed' },
        },
      };
      const fetchCallback = vi.fn();

      await deleteDataSourceHandler(
        { id: 'ds-1', name: 'DS1' },
        dataSourcesList,
        totalItems,
        currentPage,
        mockT,
        fetchCallback,
      );

      expect(dataSourcesList.value).toHaveLength(1);
      expect((dataSourcesList.value[0] as any).id).toBe('ds-2');
      expect(totalItems.value).toBe(1);
    });

    it('should go to previous page when current page becomes empty', async () => {
      (apis.httpGateway as any).datasource.delete.mockResolvedValue({});

      const dataSourcesList = ref<EnrichedDataSource[]>([{ id: 'ds-1', name: 'DS1' }]);
      const totalItems = ref(1);
      const currentPage = ref(3);
      const mockT = {
        notifications: {
          deleteSuccess: { value: 'Deleted' },
          deleteFailed: { value: 'Failed' },
        },
      };
      const fetchCallback = vi.fn().mockResolvedValue(undefined);

      await deleteDataSourceHandler(
        { id: 'ds-1', name: 'DS1' },
        dataSourcesList,
        totalItems,
        currentPage,
        mockT,
        fetchCallback,
      );

      expect(currentPage.value).toBe(2);
      expect(fetchCallback).toHaveBeenCalledOnce();
    });

    it('should not delete when API is not initialized', async () => {
      (apis as any).httpGateway = undefined;

      const dataSourcesList = ref<EnrichedDataSource[]>([]);
      const totalItems = ref(0);
      const currentPage = ref(1);
      const mockT = {
        notifications: {
          deleteSuccess: { value: 'Deleted' },
          deleteFailed: { value: 'Failed' },
        },
      };
      const fetchCallback = vi.fn();

      await deleteDataSourceHandler(
        { id: 'ds-1' },
        dataSourcesList,
        totalItems,
        currentPage,
        mockT,
        fetchCallback,
      );

      expect(fetchCallback).not.toHaveBeenCalled();
    });
  });
});

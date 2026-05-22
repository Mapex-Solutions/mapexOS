import { describe, it, expect, vi, beforeEach } from 'vitest';
import { ref } from 'vue';
import type {
  GroupsListPageFilters,
  GroupsListPageColumnVisibility,
} from '../interfaces';
import {
  fetchGroupsHandler,
  handleFilterApplyHandler,
  handlePageChangeHandler,
  handleItemsPerPageChangeHandler,
  handleColumnsUpdateHandler,
  canModifyGroupHandler,
  editGroupHandler,
  deleteGroupHandler,
} from './groupsListPage.handler';
import { apis } from '@services/mapex';

vi.mock('@utils/alert', () => ({
  notifySuccess: vi.fn(),
  notifyFail: vi.fn(),
  notifyWarning: vi.fn(),
  dialogDelete: vi.fn(),
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

function makeFilters(overrides: Partial<GroupsListPageFilters> = {}): GroupsListPageFilters {
  return {
    name: undefined,
    enabled: undefined,
    memberId: undefined,
    includeChildren: undefined,
    ...overrides,
  };
}

function makeColumnVisibility(overrides: Partial<GroupsListPageColumnVisibility> = {}): GroupsListPageColumnVisibility {
  return {
    organization: true,
    description: true,
    membersCount: true,
    created: true,
    ...overrides,
  };
}

describe('groupsListPage.handler', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  // ── fetchGroupsHandler ───────────────────────────────────────────────
  describe('fetchGroupsHandler', () => {
    it('should fetch groups and enrich with organization name', async () => {
      vi.mocked(apis.mapexOS.groups.list).mockResolvedValue({
        items: [
          { id: 'g-1', name: 'Group1', orgId: 'org-1' },
          { id: 'g-2', name: 'Group2', orgId: 'org-2' },
        ],
        pagination: { totalPages: 1, totalItems: 2 },
      } as any);

      const filters = ref(makeFilters());
      const currentPage = ref(1);
      const itemsPerPage = ref(10);
      const groupsList = ref<any[]>([]);
      const totalPages = ref(0);
      const totalItems = ref(0);
      const loading = ref(false);
      const error = ref<string | undefined>(undefined);

      await fetchGroupsHandler(
        filters, currentPage, itemsPerPage, groupsList,
        totalPages, totalItems, loading, error,
      );

      expect(loading.value).toBe(false);
      expect(groupsList.value).toHaveLength(2);
      expect(groupsList.value[0].organizationName).toBe('Org One');
      expect(groupsList.value[1].organizationName).toBe('Org Two');
      expect(totalItems.value).toBe(2);
    });

    it('should set error when API is not initialized', async () => {
      const origMapexOS = apis.mapexOS;
      (apis as any).mapexOS = undefined;

      const filters = ref(makeFilters());
      const currentPage = ref(1);
      const itemsPerPage = ref(10);
      const groupsList = ref<any[]>([]);
      const totalPages = ref(0);
      const totalItems = ref(0);
      const loading = ref(false);
      const error = ref<string | undefined>(undefined);

      await fetchGroupsHandler(
        filters, currentPage, itemsPerPage, groupsList,
        totalPages, totalItems, loading, error,
      );

      expect(error.value).toBe('Groups API not initialized');

      // Restore
      (apis as any).mapexOS = origMapexOS;
    });

    it('should include name and memberId filters when set', async () => {
      vi.mocked(apis.mapexOS.groups.list).mockResolvedValue({
        items: [],
        pagination: { totalPages: 1, totalItems: 0 },
      } as any);

      const filters = ref(makeFilters({ name: 'devs', memberId: 'user-1', enabled: true }));
      const currentPage = ref(1);
      const itemsPerPage = ref(10);
      const groupsList = ref<any[]>([]);
      const totalPages = ref(0);
      const totalItems = ref(0);
      const loading = ref(false);
      const error = ref<string | undefined>(undefined);

      await fetchGroupsHandler(
        filters, currentPage, itemsPerPage, groupsList,
        totalPages, totalItems, loading, error,
      );

      const callArgs = vi.mocked(apis.mapexOS.groups.list).mock.calls[0]![0] as any;
      expect(callArgs.name).toBe('devs');
      expect(callArgs.memberId).toBe('user-1');
      expect(callArgs.enabled).toBe(true);
    });

    it('should handle API errors', async () => {
      vi.mocked(apis.mapexOS.groups.list).mockRejectedValue(new Error('Server down'));

      const filters = ref(makeFilters());
      const currentPage = ref(1);
      const itemsPerPage = ref(10);
      const groupsList = ref<any[]>([]);
      const totalPages = ref(0);
      const totalItems = ref(0);
      const loading = ref(false);
      const error = ref<string | undefined>(undefined);

      await fetchGroupsHandler(
        filters, currentPage, itemsPerPage, groupsList,
        totalPages, totalItems, loading, error,
      );

      expect(error.value).toBe('Server down');
      expect(loading.value).toBe(false);
    });
  });

  // ── handleFilterApplyHandler ─────────────────────────────────────────
  describe('handleFilterApplyHandler', () => {
    it('should update filters, reset page, and call fetch', () => {
      const filters = ref(makeFilters());
      const columnVisibility = ref(makeColumnVisibility());
      const currentPage = ref(4);
      const fetchCallback = vi.fn();

      handleFilterApplyHandler(
        { name: 'test', enabled: false, memberId: 'u-1' },
        filters,
        columnVisibility,
        currentPage,
        fetchCallback,
      );

      expect(filters.value.name).toBe('test');
      expect(filters.value.enabled).toBe(false);
      expect(filters.value.memberId).toBe('u-1');
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

      expect(columnVisibility.value.description).toBe(false);
      expect(columnVisibility.value.membersCount).toBe(false);
      expect(columnVisibility.value.created).toBe(false);
    });

    it('should restore columns when includeChildren is disabled', () => {
      const filters = ref(makeFilters());
      const columnVisibility = ref(makeColumnVisibility({
        description: false,
        membersCount: false,
        created: false,
      }));
      const currentPage = ref(1);
      const fetchCallback = vi.fn();

      handleFilterApplyHandler(
        { includeChildren: false },
        filters,
        columnVisibility,
        currentPage,
        fetchCallback,
      );

      expect(columnVisibility.value.description).toBe(true);
      expect(columnVisibility.value.membersCount).toBe(true);
      expect(columnVisibility.value.created).toBe(true);
    });
  });

  // ── handlePageChangeHandler ──────────────────────────────────────────
  describe('handlePageChangeHandler', () => {
    it('should set current page and call fetch', () => {
      const currentPage = ref(1);
      const fetchCallback = vi.fn();

      handlePageChangeHandler(3, currentPage, fetchCallback);

      expect(currentPage.value).toBe(3);
      expect(fetchCallback).toHaveBeenCalledOnce();
    });
  });

  // ── handleItemsPerPageChangeHandler ──────────────────────────────────
  describe('handleItemsPerPageChangeHandler', () => {
    it('should set items per page, reset to page 1, and call fetch', () => {
      const itemsPerPage = ref(10);
      const currentPage = ref(5);
      const fetchCallback = vi.fn();

      handleItemsPerPageChangeHandler(25, itemsPerPage, currentPage, fetchCallback);

      expect(itemsPerPage.value).toBe(25);
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
          { key: 'description', visible: false, label: '' },
          { key: 'membersCount', visible: true, label: '' },
          { key: 'created', visible: false, label: '' },
        ],
        columnVisibility,
      );

      expect(columnVisibility.value.organization).toBe(false);
      expect(columnVisibility.value.description).toBe(false);
      expect(columnVisibility.value.membersCount).toBe(true);
      expect(columnVisibility.value.created).toBe(false);
    });
  });

  // ── canModifyGroupHandler ────────────────────────────────────────────
  describe('canModifyGroupHandler', () => {
    it('should return true when group orgId matches selected organization', () => {
      expect(canModifyGroupHandler({ orgId: 'org-1' })).toBe(true);
    });

    it('should return false when group orgId does not match', () => {
      expect(canModifyGroupHandler({ orgId: 'org-other' })).toBe(false);
    });
  });

  // ── editGroupHandler ─────────────────────────────────────────────────
  describe('editGroupHandler', () => {
    it('should navigate to edit page for modifiable group', () => {
      const router = { push: vi.fn() };

      editGroupHandler(
        { id: 'g-1', orgId: 'org-1' },
        { notifications: { cannotEdit: { value: '' } } },
        router,
      );

      expect(router.push).toHaveBeenCalledWith('/groups/edit/g-1');
    });

    it('should warn when trying to edit a group from another org', () => {
      const router = { push: vi.fn() };

      editGroupHandler(
        { id: 'g-1', orgId: 'org-other' },
        { notifications: { cannotEdit: { value: 'Cannot edit' } } },
        router,
      );

      expect(router.push).not.toHaveBeenCalled();
    });

    it('should notify fail when group ID is missing', () => {
      const router = { push: vi.fn() };

      editGroupHandler(
        { orgId: 'org-1' },
        { notifications: { cannotEdit: { value: '' } } },
        router,
      );

      expect(router.push).not.toHaveBeenCalled();
    });
  });

  // ── deleteGroupHandler ───────────────────────────────────────────────
  describe('deleteGroupHandler', () => {
    it('should delete group and remove from local list', async () => {
      vi.mocked(apis.mapexOS.groups.delete).mockResolvedValue({} as any);

      const groupsList = ref<any[]>([
        { id: 'g-1', name: 'Group1' },
        { id: 'g-2', name: 'Group2' },
      ]);
      const totalItems = ref(2);
      const currentPage = ref(1);
      const mockT = { messages: { deletedSuccessfully: { value: 'Deleted' } } };
      const fetchCallback = vi.fn();

      await deleteGroupHandler(
        { id: 'g-1', name: 'Group1' },
        groupsList,
        totalItems,
        currentPage,
        mockT,
        fetchCallback,
      );

      expect(groupsList.value).toHaveLength(1);
      expect(groupsList.value[0].id).toBe('g-2');
      expect(totalItems.value).toBe(1);
    });

    it('should go to previous page when current page becomes empty', async () => {
      vi.mocked(apis.mapexOS.groups.delete).mockResolvedValue({} as any);

      const groupsList = ref<any[]>([{ id: 'g-1', name: 'Group1' }]);
      const totalItems = ref(1);
      const currentPage = ref(2);
      const mockT = { messages: { deletedSuccessfully: { value: 'Deleted' } } };
      const fetchCallback = vi.fn().mockResolvedValue(undefined);

      await deleteGroupHandler(
        { id: 'g-1', name: 'Group1' },
        groupsList,
        totalItems,
        currentPage,
        mockT,
        fetchCallback,
      );

      expect(currentPage.value).toBe(1);
      expect(fetchCallback).toHaveBeenCalledOnce();
    });

    it('should not delete when group ID is missing', async () => {
      const groupsList = ref<any[]>([]);
      const totalItems = ref(0);
      const currentPage = ref(1);
      const mockT = { messages: { deletedSuccessfully: { value: 'Deleted' } } };
      const fetchCallback = vi.fn();

      await deleteGroupHandler(
        { name: 'NoId' },
        groupsList,
        totalItems,
        currentPage,
        mockT,
        fetchCallback,
      );

      expect(apis.mapexOS.groups.delete).not.toHaveBeenCalled();
    });
  });
});

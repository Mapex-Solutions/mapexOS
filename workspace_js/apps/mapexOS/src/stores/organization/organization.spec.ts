import { describe, it, expect, beforeEach, vi } from 'vitest';
import { setActivePinia, createPinia } from 'pinia';
import { useOrganizationStore } from './index';
import type { OrganizationCoverageItem } from './types';

/** Mock useLogger — not globally mocked in setup.ts */
vi.mock('@composables/useLogger', () => ({
  useLogger: () => ({
    info: vi.fn(),
    warn: vi.fn(),
    error: vi.fn(),
    debug: vi.fn(),
  }),
}));

/** Mock treeBuilder utility */
vi.mock('@utils/organization/treeBuilder', () => ({
  buildOrganizationTree: vi.fn((items: OrganizationCoverageItem[]) =>
    items.map(item => ({ ...item, depth: 0, enabled: true, children: [] })),
  ),
  getParentPathKey: vi.fn((pathKey: string) => {
    const parts = pathKey.split('.');
    if (parts.length <= 1) return null;
    return parts.slice(0, -1).join('.');
  }),
}));

/** Shared mock instance so assertions work across module boundaries */
const mockPermStore = {
  fetchPermissions: vi.fn().mockResolvedValue(undefined),
  clearPermissions: vi.fn(),
};

/** Mock permission store — organization actions depend on it */
vi.mock('@stores/permission', () => ({
  usePermissionStore: () => mockPermStore,
}));

// ────────────────────────────────────────────────────────────────────────
// Test fixtures
// ────────────────────────────────────────────────────────────────────────

const MOCK_VENDOR: OrganizationCoverageItem = {
  id: 'vendor-1',
  name: 'Mapex Global',
  type: 'vendor',
  pathKey: 'mapex',
  scope: 'recursive',
  membershipId: 'mem-1',
  roleIds: ['role-admin'],
};

const MOCK_CUSTOMER: OrganizationCoverageItem = {
  id: 'customer-1',
  name: 'Acme Corp',
  type: 'customer',
  pathKey: 'mapex.acme',
  scope: 'inherited',
  membershipId: 'mem-2',
  roleIds: ['role-user'],
};

const MOCK_SITE: OrganizationCoverageItem = {
  id: 'site-1',
  name: 'Acme HQ',
  type: 'site',
  pathKey: 'mapex.acme.hq',
  scope: 'inherited',
  membershipId: 'mem-2',
  roleIds: ['role-user'],
};

const MOCK_BUILDING: OrganizationCoverageItem = {
  id: 'building-1',
  name: 'Building A',
  type: 'building',
  pathKey: 'mapex.acme.hq.bldgA',
  scope: 'inherited',
  membershipId: 'mem-2',
  roleIds: ['role-user'],
};

const ALL_ORGS = [MOCK_VENDOR, MOCK_CUSTOMER, MOCK_SITE, MOCK_BUILDING];

const MOCK_COVERAGE_RESPONSE = {
  lastUpdated: '2026-03-17T10:00:00Z',
  organizations: ALL_ORGS,
};

describe('OrganizationStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia());
    vi.clearAllMocks();
    localStorage.clear();
  });

  // ────────────────────────────────────────────────────────────────────────
  // State
  // ────────────────────────────────────────────────────────────────────────

  describe('state', () => {
    it('has coverage null by default', () => {
      const store = useOrganizationStore();
      expect(store.coverage).toBeNull();
    });

    it('has empty treeNodes by default', () => {
      const store = useOrganizationStore();
      expect(store.treeNodes).toEqual([]);
    });

    it('has empty flatList by default', () => {
      const store = useOrganizationStore();
      expect(store.flatList).toEqual([]);
    });

    it('has loading false by default', () => {
      const store = useOrganizationStore();
      expect(store.loading).toBe(false);
    });

    it('has error null by default', () => {
      const store = useOrganizationStore();
      expect(store.error).toBeNull();
    });

    it('has lastUpdated null by default', () => {
      const store = useOrganizationStore();
      expect(store.lastUpdated).toBeNull();
    });

    it('has selectedOrganizationId null by default', () => {
      const store = useOrganizationStore();
      expect(store.selectedOrganizationId).toBeNull();
    });

    it('has selectedOrganizationName null by default', () => {
      const store = useOrganizationStore();
      expect(store.selectedOrganizationName).toBeNull();
    });
  });

  // ────────────────────────────────────────────────────────────────────────
  // Getters
  // ────────────────────────────────────────────────────────────────────────

  describe('getters', () => {
    describe('getOrganizationById', () => {
      it('returns the org when found', () => {
        const store = useOrganizationStore();
        store.flatList = [...ALL_ORGS];

        expect(store.getOrganizationById('vendor-1')).toEqual(MOCK_VENDOR);
      });

      it('returns undefined when not found', () => {
        const store = useOrganizationStore();
        store.flatList = [...ALL_ORGS];

        expect(store.getOrganizationById('nonexistent')).toBeUndefined();
      });
    });

    describe('getOrganizationsByType', () => {
      it('returns orgs matching the type', () => {
        const store = useOrganizationStore();
        store.flatList = [...ALL_ORGS];

        const customers = store.getOrganizationsByType('customer');
        expect(customers).toHaveLength(1);
        expect(customers[0]!.id).toBe('customer-1');
      });

      it('returns empty array when no match', () => {
        const store = useOrganizationStore();
        store.flatList = [...ALL_ORGS];

        expect(store.getOrganizationsByType('zone')).toHaveLength(0);
      });
    });

    describe('vendors', () => {
      it('returns only vendor-type tree nodes', () => {
        const store = useOrganizationStore();
        store.treeNodes = ALL_ORGS.map(o => ({ ...o, depth: 0, enabled: true, children: [] }));

        expect(store.vendors).toHaveLength(1);
        expect(store.vendors[0]!.type).toBe('vendor');
      });
    });

    describe('hasAccessTo', () => {
      it('returns true when org exists in flatList', () => {
        const store = useOrganizationStore();
        store.flatList = [...ALL_ORGS];

        expect(store.hasAccessTo('vendor-1')).toBe(true);
      });

      it('returns false when org does not exist', () => {
        const store = useOrganizationStore();
        store.flatList = [...ALL_ORGS];

        expect(store.hasAccessTo('nonexistent')).toBe(false);
      });
    });

    describe('selectedOrganization', () => {
      it('returns null when no org is selected', () => {
        const store = useOrganizationStore();

        expect(store.selectedOrganization).toBeNull();
      });

      it('returns the selected org details', () => {
        const store = useOrganizationStore();
        store.flatList = [...ALL_ORGS];
        store.selectedOrganizationId = 'customer-1';

        expect(store.selectedOrganization).toEqual(MOCK_CUSTOMER);
      });
    });

    describe('isDataStale', () => {
      it('returns true when lastUpdated is null', () => {
        const store = useOrganizationStore();

        expect(store.isDataStale).toBe(true);
      });

      it('returns false when lastUpdated is recent', () => {
        const store = useOrganizationStore();
        store.lastUpdated = new Date().toISOString();

        expect(store.isDataStale).toBe(false);
      });

      it('returns true when lastUpdated is older than 5 minutes', () => {
        const store = useOrganizationStore();
        const sixMinutesAgo = new Date(Date.now() - 6 * 60 * 1000);
        store.lastUpdated = sixMinutesAgo.toISOString();

        expect(store.isDataStale).toBe(true);
      });
    });

    describe('isSelected', () => {
      it('returns true for the selected org', () => {
        const store = useOrganizationStore();
        store.selectedOrganizationId = 'vendor-1';

        expect(store.isSelected('vendor-1')).toBe(true);
      });

      it('returns false for a different org', () => {
        const store = useOrganizationStore();
        store.selectedOrganizationId = 'vendor-1';

        expect(store.isSelected('customer-1')).toBe(false);
      });
    });

    describe('totalCount', () => {
      it('returns the number of orgs in flatList', () => {
        const store = useOrganizationStore();
        store.flatList = [...ALL_ORGS];

        expect(store.totalCount).toBe(4);
      });

      it('returns 0 when flatList is empty', () => {
        const store = useOrganizationStore();

        expect(store.totalCount).toBe(0);
      });
    });

    describe('isVendor', () => {
      it('returns true when selected org is vendor type', () => {
        const store = useOrganizationStore();
        store.flatList = [...ALL_ORGS];
        store.selectedOrganizationId = 'vendor-1';

        expect(store.isVendor).toBe(true);
      });

      it('returns false when selected org is not vendor', () => {
        const store = useOrganizationStore();
        store.flatList = [...ALL_ORGS];
        store.selectedOrganizationId = 'customer-1';

        expect(store.isVendor).toBe(false);
      });

      it('returns false when no org is selected', () => {
        const store = useOrganizationStore();

        expect(store.isVendor).toBe(false);
      });
    });

    describe('isCustomer', () => {
      it('returns true when selected org is customer type', () => {
        const store = useOrganizationStore();
        store.flatList = [...ALL_ORGS];
        store.selectedOrganizationId = 'customer-1';

        expect(store.isCustomer).toBe(true);
      });

      it('returns false when selected org is not customer', () => {
        const store = useOrganizationStore();
        store.flatList = [...ALL_ORGS];
        store.selectedOrganizationId = 'vendor-1';

        expect(store.isCustomer).toBe(false);
      });
    });

    describe('isSite', () => {
      it('returns true when selected org is site type', () => {
        const store = useOrganizationStore();
        store.flatList = [...ALL_ORGS];
        store.selectedOrganizationId = 'site-1';

        expect(store.isSite).toBe(true);
      });

      it('returns true when selected org is building type', () => {
        const store = useOrganizationStore();
        store.flatList = [...ALL_ORGS];
        store.selectedOrganizationId = 'building-1';

        expect(store.isSite).toBe(true);
      });

      it('returns false when selected org is vendor', () => {
        const store = useOrganizationStore();
        store.flatList = [...ALL_ORGS];
        store.selectedOrganizationId = 'vendor-1';

        expect(store.isSite).toBe(false);
      });
    });
  });

  // ────────────────────────────────────────────────────────────────────────
  // Actions
  // ────────────────────────────────────────────────────────────────────────

  describe('actions', () => {
    describe('fetchCoverage', () => {
      it('calls API and populates state on success', async () => {
        const { apis } = await import('@src/services/mapex');

        (apis.mapexOS as any).auth = {
          getUserCoverage: vi.fn().mockResolvedValue(MOCK_COVERAGE_RESPONSE),
        };

        const store = useOrganizationStore();
        await store.fetchCoverage(true);

        expect(apis.mapexOS.auth.getUserCoverage).toHaveBeenCalled();
        expect(store.coverage).toEqual(MOCK_COVERAGE_RESPONSE);
        expect(store.flatList).toEqual(ALL_ORGS);
        expect(store.lastUpdated).toBe('2026-03-17T10:00:00Z');
        expect(store.treeNodes).toHaveLength(4);
        expect(store.loading).toBe(false);
        expect(store.error).toBeNull();
      });

      it('builds tree from organizations', async () => {
        const { apis } = await import('@src/services/mapex');
        const { buildOrganizationTree } = await import('@utils/organization/treeBuilder');

        (apis.mapexOS as any).auth = {
          getUserCoverage: vi.fn().mockResolvedValue(MOCK_COVERAGE_RESPONSE),
        };

        const store = useOrganizationStore();
        await store.fetchCoverage(true);

        expect(buildOrganizationTree).toHaveBeenCalledWith(ALL_ORGS);
      });

      it('sets error when API returns null', async () => {
        const { apis } = await import('@src/services/mapex');

        (apis.mapexOS as any).auth = {
          getUserCoverage: vi.fn().mockResolvedValue(null),
        };

        const store = useOrganizationStore();

        await expect(store.fetchCoverage(true)).rejects.toThrow('No coverage data received');
        expect(store.error).toBe('No coverage data received');
        expect(store.loading).toBe(false);
      });

      it('sets error when API returns empty organizations', async () => {
        const { apis } = await import('@src/services/mapex');

        (apis.mapexOS as any).auth = {
          getUserCoverage: vi.fn().mockResolvedValue({
            lastUpdated: '2026-03-17T10:00:00Z',
            organizations: [],
          }),
        };

        const store = useOrganizationStore();

        await expect(store.fetchCoverage(true)).rejects.toThrow('User has no organization access');
        expect(store.error).toBe('User has no organization access');
      });

      it('sets error when API throws', async () => {
        const { apis } = await import('@src/services/mapex');

        (apis.mapexOS as any).auth = {
          getUserCoverage: vi.fn().mockRejectedValue(new Error('Network error')),
        };

        const store = useOrganizationStore();

        await expect(store.fetchCoverage(true)).rejects.toThrow('Network error');
        expect(store.error).toBe('Network error');
        expect(store.loading).toBe(false);
      });

      it('skips fetch when data is fresh and not forcing', async () => {
        const { apis } = await import('@src/services/mapex');
        const mockFn = vi.fn();
        (apis.mapexOS as any).auth = { getUserCoverage: mockFn };

        const store = useOrganizationStore();
        store.coverage = MOCK_COVERAGE_RESPONSE;
        store.lastUpdated = new Date().toISOString();

        await store.fetchCoverage(false);

        expect(mockFn).not.toHaveBeenCalled();
      });

      it('fetches when forced even if data is fresh', async () => {
        const { apis } = await import('@src/services/mapex');
        const mockFn = vi.fn().mockResolvedValue(MOCK_COVERAGE_RESPONSE);
        (apis.mapexOS as any).auth = { getUserCoverage: mockFn };

        const store = useOrganizationStore();
        store.coverage = MOCK_COVERAGE_RESPONSE;
        store.lastUpdated = new Date().toISOString();

        await store.fetchCoverage(true);

        expect(mockFn).toHaveBeenCalled();
      });
    });

    describe('initializeAfterLogin', () => {
      it('fetches coverage and selects initial org', async () => {
        const { apis } = await import('@src/services/mapex');

        (apis.mapexOS as any).auth = {
          getUserCoverage: vi.fn().mockResolvedValue(MOCK_COVERAGE_RESPONSE),
        };

        const store = useOrganizationStore();
        const selectedId = await store.initializeAfterLogin();

        expect(store.flatList).toHaveLength(4);
        expect(store.selectedOrganizationId).toBeTruthy();
        expect(selectedId).toBeTruthy();
      });

      it('auto-selects single org when only one exists', async () => {
        const { apis } = await import('@src/services/mapex');

        (apis.mapexOS as any).auth = {
          getUserCoverage: vi.fn().mockResolvedValue({
            lastUpdated: '2026-03-17T10:00:00Z',
            organizations: [MOCK_VENDOR],
          }),
        };

        const store = useOrganizationStore();
        const selectedId = await store.initializeAfterLogin();

        expect(selectedId).toBe('vendor-1');
        expect(store.selectedOrganizationId).toBe('vendor-1');
      });

      it('restores last selected org from localStorage', async () => {
        const { apis } = await import('@src/services/mapex');
        localStorage.setItem('selectedOrgId', 'customer-1');

        (apis.mapexOS as any).auth = {
          getUserCoverage: vi.fn().mockResolvedValue(MOCK_COVERAGE_RESPONSE),
        };

        const store = useOrganizationStore();
        const selectedId = await store.initializeAfterLogin();

        expect(selectedId).toBe('customer-1');
      });

      it('falls back to recursive scope org when no stored preference', async () => {
        const { apis } = await import('@src/services/mapex');

        (apis.mapexOS as any).auth = {
          getUserCoverage: vi.fn().mockResolvedValue(MOCK_COVERAGE_RESPONSE),
        };

        const store = useOrganizationStore();
        const selectedId = await store.initializeAfterLogin();

        // MOCK_VENDOR has scope 'recursive'
        expect(selectedId).toBe('vendor-1');
      });

      it('saves selected org to localStorage', async () => {
        const { apis } = await import('@src/services/mapex');

        (apis.mapexOS as any).auth = {
          getUserCoverage: vi.fn().mockResolvedValue({
            lastUpdated: '2026-03-17T10:00:00Z',
            organizations: [MOCK_VENDOR],
          }),
        };

        const store = useOrganizationStore();
        await store.initializeAfterLogin();

        expect(localStorage.getItem('selectedOrgId')).toBe('vendor-1');
      });

      it('fetches permissions after selecting org', async () => {
        const { apis } = await import('@src/services/mapex');

        (apis.mapexOS as any).auth = {
          getUserCoverage: vi.fn().mockResolvedValue({
            lastUpdated: '2026-03-17T10:00:00Z',
            organizations: [MOCK_VENDOR],
          }),
        };

        const store = useOrganizationStore();
        await store.initializeAfterLogin();

        expect(mockPermStore.fetchPermissions).toHaveBeenCalledWith(true);
      });
    });

    describe('selectOrganization', () => {
      it('sets selectedOrganizationId and name', () => {
        const store = useOrganizationStore();
        store.flatList = [...ALL_ORGS];

        store.selectOrganization('customer-1');

        expect(store.selectedOrganizationId).toBe('customer-1');
        expect(store.selectedOrganizationName).toBe('Acme Corp');
      });

      it('saves to localStorage', () => {
        const store = useOrganizationStore();
        store.flatList = [...ALL_ORGS];

        store.selectOrganization('customer-1');

        expect(localStorage.getItem('selectedOrgId')).toBe('customer-1');
        expect(localStorage.getItem('selectedOrgName')).toBe('Acme Corp');
      });

      it('clears localStorage when setting null', () => {
        localStorage.setItem('selectedOrgId', 'vendor-1');
        localStorage.setItem('selectedOrgName', 'Mapex Global');

        const store = useOrganizationStore();
        store.selectOrganization(null);

        expect(store.selectedOrganizationId).toBeNull();
        expect(store.selectedOrganizationName).toBeNull();
        expect(localStorage.getItem('selectedOrgId')).toBeNull();
        expect(localStorage.getItem('selectedOrgName')).toBeNull();
      });

      it('does not set when org is not in flatList', () => {
        const store = useOrganizationStore();
        store.flatList = [...ALL_ORGS];
        store.selectedOrganizationId = 'vendor-1';

        store.selectOrganization('nonexistent-org');

        // Should remain unchanged
        expect(store.selectedOrganizationId).toBe('vendor-1');
      });

      it('refreshes permissions when selecting an org', () => {
        const store = useOrganizationStore();
        store.flatList = [...ALL_ORGS];

        store.selectOrganization('vendor-1');

        expect(mockPermStore.fetchPermissions).toHaveBeenCalledWith(true);
      });
    });

    describe('addOrganizationToTree', () => {
      it('adds org to flatList and rebuilds tree', async () => {
        const { buildOrganizationTree } = await import('@utils/organization/treeBuilder');

        const store = useOrganizationStore();
        store.flatList = [MOCK_VENDOR, MOCK_CUSTOMER];

        store.addOrganizationToTree({
          id: 'site-new',
          name: 'New Site',
          type: 'site',
          pathKey: 'mapex.acme.newsite',
        });

        expect(store.flatList).toHaveLength(3);
        expect(store.flatList[2]!.id).toBe('site-new');
        expect(store.flatList[2]!.scope).toBe('inherited'); // inherited from MOCK_CUSTOMER
        expect(store.flatList[2]!.membershipId).toBe('mem-2');
        expect(store.flatList[2]!.roleIds).toEqual(['role-user']);
        expect(buildOrganizationTree).toHaveBeenCalled();
      });

      it('does nothing when parent pathKey is not found', () => {
        const store = useOrganizationStore();
        store.flatList = [MOCK_VENDOR];

        const originalLength = store.flatList.length;

        store.addOrganizationToTree({
          id: 'orphan',
          name: 'Orphan',
          type: 'site',
          pathKey: 'nonexistent.parent.orphan',
        });

        expect(store.flatList).toHaveLength(originalLength);
      });

      it('does nothing when pathKey has no parent (root)', () => {
        const store = useOrganizationStore();
        store.flatList = [MOCK_VENDOR];

        const originalLength = store.flatList.length;

        store.addOrganizationToTree({
          id: 'root',
          name: 'Root',
          type: 'vendor',
          pathKey: 'root',
        });

        expect(store.flatList).toHaveLength(originalLength);
      });
    });

    describe('clearCoverage', () => {
      it('resets all state to defaults', () => {
        const store = useOrganizationStore();
        store.coverage = MOCK_COVERAGE_RESPONSE;
        store.treeNodes = ALL_ORGS.map(o => ({ ...o, depth: 0, enabled: true, children: [] }));
        store.flatList = [...ALL_ORGS];
        store.selectedOrganizationId = 'vendor-1';
        store.selectedOrganizationName = 'Mapex Global';
        store.lastUpdated = '2026-03-17T10:00:00Z';
        store.error = 'some error';
        localStorage.setItem('selectedOrgId', 'vendor-1');
        localStorage.setItem('selectedOrgName', 'Mapex Global');

        store.clearCoverage();

        expect(store.coverage).toBeNull();
        expect(store.treeNodes).toEqual([]);
        expect(store.flatList).toEqual([]);
        expect(store.selectedOrganizationId).toBeNull();
        expect(store.selectedOrganizationName).toBeNull();
        expect(store.lastUpdated).toBeNull();
        expect(store.error).toBeNull();
        expect(localStorage.getItem('selectedOrgId')).toBeNull();
        expect(localStorage.getItem('selectedOrgName')).toBeNull();
      });

      it('clears permission store', () => {
        const store = useOrganizationStore();
        store.clearCoverage();

        expect(mockPermStore.clearPermissions).toHaveBeenCalled();
      });
    });
  });
});

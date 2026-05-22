import { describe, it, expect, beforeEach, vi } from 'vitest';
import { setActivePinia, createPinia } from 'pinia';
import { usePermissionStore } from './index';

/** Mock useLogger — not globally mocked in setup.ts */
vi.mock('@composables/useLogger', () => ({
  useLogger: () => ({
    info: vi.fn(),
    warn: vi.fn(),
    error: vi.fn(),
    debug: vi.fn(),
  }),
}));

/** Mock organization store — permission actions depend on it */
vi.mock('@stores/organization', () => ({
  useOrganizationStore: () => ({
    selectedOrganizationId: 'org-123',
  }),
}));

describe('PermissionStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia());
  });

  // ────────────────────────────────────────────────────────────────────────
  // State
  // ────────────────────────────────────────────────────────────────────────

  describe('state', () => {
    it('has empty permissions by default', () => {
      const store = usePermissionStore();
      expect(store.permissions).toEqual([]);
    });

    it('has loading false by default', () => {
      const store = usePermissionStore();
      expect(store.loading).toBe(false);
    });

    it('has error null by default', () => {
      const store = usePermissionStore();
      expect(store.error).toBeNull();
    });

    it('has version 0 by default', () => {
      const store = usePermissionStore();
      expect(store.version).toBe(0);
    });

    it('has forOrganizationId null by default', () => {
      const store = usePermissionStore();
      expect(store.forOrganizationId).toBeNull();
    });

    it('has lastFetched null by default', () => {
      const store = usePermissionStore();
      expect(store.lastFetched).toBeNull();
    });
  });

  // ────────────────────────────────────────────────────────────────────────
  // Getters
  // ────────────────────────────────────────────────────────────────────────

  describe('getters', () => {
    describe('hasPermission', () => {
      it('returns true when exact permission exists', () => {
        const store = usePermissionStore();
        store.permissions = ['users.list', 'users.create'];

        expect(store.hasPermission('users.list')).toBe(true);
      });

      it('returns false when permission does not exist', () => {
        const store = usePermissionStore();
        store.permissions = ['users.create'];

        expect(store.hasPermission('users.list')).toBe(false);
      });

      it('returns false when permissions array is empty', () => {
        const store = usePermissionStore();

        expect(store.hasPermission('users.list')).toBe(false);
      });

      it('matches resource wildcard (users.*) against users.list', () => {
        const store = usePermissionStore();
        store.permissions = ['users.*'];

        expect(store.hasPermission('users.list')).toBe(true);
      });

      it('matches resource wildcard against nested permissions', () => {
        const store = usePermissionStore();
        store.permissions = ['events.*'];

        expect(store.hasPermission('events.raw.list')).toBe(true);
      });

      it('does not match wildcard for different resource', () => {
        const store = usePermissionStore();
        store.permissions = ['users.*'];

        expect(store.hasPermission('roles.list')).toBe(false);
      });

      it('matches root wildcard (mapex.*) against anything', () => {
        const store = usePermissionStore();
        store.permissions = ['mapex.*'];

        expect(store.hasPermission('users.list')).toBe(true);
        expect(store.hasPermission('roles.create')).toBe(true);
        expect(store.hasPermission('events.raw.list')).toBe(true);
      });

      it('matches admin wildcards (admin_vendor.*, admin_customer.*, admin.*)', () => {
        const store = usePermissionStore();

        store.permissions = ['admin_vendor.*'];
        expect(store.hasPermission('users.list')).toBe(true);

        store.permissions = ['admin_customer.*'];
        expect(store.hasPermission('roles.create')).toBe(true);

        store.permissions = ['admin.*'];
        expect(store.hasPermission('anything.here')).toBe(true);
      });
    });

    describe('hasAnyPermission', () => {
      it('returns true if any permission matches', () => {
        const store = usePermissionStore();
        store.permissions = ['users.list'];

        expect(store.hasAnyPermission(['users.list', 'roles.create'])).toBe(true);
      });

      it('returns false if none match', () => {
        const store = usePermissionStore();
        store.permissions = ['users.list'];

        expect(store.hasAnyPermission(['roles.create', 'groups.delete'])).toBe(false);
      });

      it('returns true for empty required array', () => {
        const store = usePermissionStore();

        expect(store.hasAnyPermission([])).toBe(true);
      });
    });

    describe('hasAllPermissions', () => {
      it('returns true when all permissions match', () => {
        const store = usePermissionStore();
        store.permissions = ['users.list', 'roles.create'];

        expect(store.hasAllPermissions(['users.list', 'roles.create'])).toBe(true);
      });

      it('returns false when some are missing', () => {
        const store = usePermissionStore();
        store.permissions = ['users.list'];

        expect(store.hasAllPermissions(['users.list', 'roles.create'])).toBe(false);
      });

      it('returns true for empty required array', () => {
        const store = usePermissionStore();

        expect(store.hasAllPermissions([])).toBe(true);
      });
    });

    describe('isLoaded', () => {
      it('returns false when permissions are empty', () => {
        const store = usePermissionStore();

        expect(store.isLoaded).toBe(false);
      });

      it('returns false when loading is true', () => {
        const store = usePermissionStore();
        store.permissions = ['users.list'];
        store.loading = true;

        expect(store.isLoaded).toBe(false);
      });

      it('returns true when permissions exist and not loading', () => {
        const store = usePermissionStore();
        store.permissions = ['users.list'];
        store.loading = false;

        expect(store.isLoaded).toBe(true);
      });
    });

    describe('isStale', () => {
      it('returns true when lastFetched is null', () => {
        const store = usePermissionStore();

        expect(store.isStale).toBe(true);
      });

      it('returns false when lastFetched is recent', () => {
        const store = usePermissionStore();
        store.lastFetched = new Date().toISOString();

        expect(store.isStale).toBe(false);
      });

      it('returns true when lastFetched is older than 5 minutes', () => {
        const store = usePermissionStore();
        const sixMinutesAgo = new Date(Date.now() - 6 * 60 * 1000);
        store.lastFetched = sixMinutesAgo.toISOString();

        expect(store.isStale).toBe(true);
      });
    });
  });

  // ────────────────────────────────────────────────────────────────────────
  // Actions
  // ────────────────────────────────────────────────────────────────────────

  describe('actions', () => {
    describe('fetchPermissions', () => {
      it('calls API and sets permissions on success', async () => {
        const { apis } = await import('@src/services/mapex');
        const mockResponse = {
          permissions: ['users.list', 'users.create', 'roles.*'],
          version: 3,
        };

        // Add auth.getMyPermissions mock since it is not in global setup
        (apis.mapexOS as any).auth = {
          getMyPermissions: vi.fn().mockResolvedValue(mockResponse),
        };

        const store = usePermissionStore();
        await store.fetchPermissions(true);

        expect(apis.mapexOS.auth.getMyPermissions).toHaveBeenCalled();
        expect(store.permissions).toEqual(['users.list', 'users.create', 'roles.*']);
        expect(store.version).toBe(3);
        expect(store.loading).toBe(false);
        expect(store.error).toBeNull();
        expect(store.forOrganizationId).toBe('org-123');
        expect(store.lastFetched).toBeTruthy();
      });

      it('sets error message when API throws', async () => {
        const { apis } = await import('@src/services/mapex');

        (apis.mapexOS as any).auth = {
          getMyPermissions: vi.fn().mockRejectedValue(new Error('Network error')),
        };

        const store = usePermissionStore();
        await store.fetchPermissions(true);

        expect(store.permissions).toEqual([]);
        expect(store.loading).toBe(false);
        expect(store.error).toBe('Network error');
      });

      it('sets generic error when API returns null', async () => {
        const { apis } = await import('@src/services/mapex');

        (apis.mapexOS as any).auth = {
          getMyPermissions: vi.fn().mockResolvedValue(null),
        };

        const store = usePermissionStore();
        await store.fetchPermissions(true);

        expect(store.error).toBe('No permission data received');
        expect(store.loading).toBe(false);
      });

      it('skips fetch when data is fresh and same org (no force)', async () => {
        const { apis } = await import('@src/services/mapex');
        const mockFn = vi.fn();
        (apis.mapexOS as any).auth = { getMyPermissions: mockFn };

        const store = usePermissionStore();
        store.permissions = ['users.list'];
        store.forOrganizationId = 'org-123';
        store.lastFetched = new Date().toISOString();

        await store.fetchPermissions(false);

        expect(mockFn).not.toHaveBeenCalled();
      });

      it('fetches when forced even if data is fresh', async () => {
        const { apis } = await import('@src/services/mapex');
        const mockResponse = { permissions: ['new.perm'], version: 5 };
        const mockFn = vi.fn().mockResolvedValue(mockResponse);
        (apis.mapexOS as any).auth = { getMyPermissions: mockFn };

        const store = usePermissionStore();
        store.permissions = ['old.perm'];
        store.forOrganizationId = 'org-123';
        store.lastFetched = new Date().toISOString();

        await store.fetchPermissions(true);

        expect(mockFn).toHaveBeenCalled();
        expect(store.permissions).toEqual(['new.perm']);
      });
    });

    describe('clearPermissions', () => {
      it('resets all state to defaults', () => {
        const store = usePermissionStore();
        store.permissions = ['users.list'];
        store.version = 5;
        store.loading = true;
        store.error = 'some error';
        store.forOrganizationId = 'org-123';
        store.lastFetched = new Date().toISOString();

        store.clearPermissions();

        expect(store.permissions).toEqual([]);
        expect(store.version).toBe(0);
        expect(store.loading).toBe(false);
        expect(store.error).toBeNull();
        expect(store.forOrganizationId).toBeNull();
        expect(store.lastFetched).toBeNull();
      });
    });
  });
});

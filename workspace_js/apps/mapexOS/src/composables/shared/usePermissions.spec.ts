import { describe, it, expect, beforeEach } from 'vitest';
import { setActivePinia, createPinia } from 'pinia';
import { usePermissions } from './usePermissions';
import { usePermissionStore } from '@stores/permission';

describe('usePermissions', () => {
  beforeEach(() => {
    setActivePinia(createPinia());
  });

  /**
   * Seed the permission store with a set of permissions.
   */
  function seedPermissions(perms: string[]) {
    const store = usePermissionStore();
    store.$patch({ permissions: perms, loading: false });
  }

  describe('hasPermission', () => {
    it('returns a computed that reflects store.hasPermission', () => {
      seedPermissions(['assets.list', 'assets.create']);

      const { hasPermission } = usePermissions();
      expect(hasPermission('assets.list').value).toBe(true);
      expect(hasPermission('assets.delete').value).toBe(false);
    });

    it('reacts to permission changes in the store', () => {
      seedPermissions([]);

      const { hasPermission } = usePermissions();
      const result = hasPermission('assets.list');
      expect(result.value).toBe(false);

      // Add permission
      const store = usePermissionStore();
      store.$patch({ permissions: ['assets.list'] });
      expect(result.value).toBe(true);
    });
  });

  describe('hasAnyPermission', () => {
    it('returns true if user has at least one of the permissions', () => {
      seedPermissions(['assets.list']);

      const { hasAnyPermission } = usePermissions();
      expect(hasAnyPermission(['assets.list', 'assets.delete']).value).toBe(true);
    });

    it('returns false if user has none of the permissions', () => {
      seedPermissions(['assets.list']);

      const { hasAnyPermission } = usePermissions();
      expect(hasAnyPermission(['users.create', 'users.delete']).value).toBe(false);
    });

    it('returns true for empty array', () => {
      seedPermissions([]);

      const { hasAnyPermission } = usePermissions();
      expect(hasAnyPermission([]).value).toBe(true);
    });
  });

  describe('hasAllPermissions', () => {
    it('returns true if user has all permissions', () => {
      seedPermissions(['assets.list', 'assets.create']);

      const { hasAllPermissions } = usePermissions();
      expect(hasAllPermissions(['assets.list', 'assets.create']).value).toBe(true);
    });

    it('returns false if user is missing at least one permission', () => {
      seedPermissions(['assets.list']);

      const { hasAllPermissions } = usePermissions();
      expect(hasAllPermissions(['assets.list', 'assets.create']).value).toBe(false);
    });

    it('returns true for empty array', () => {
      seedPermissions([]);

      const { hasAllPermissions } = usePermissions();
      expect(hasAllPermissions([]).value).toBe(true);
    });
  });

  describe('CRUD shorthands', () => {
    it('canList checks resource.list', () => {
      seedPermissions(['assets.list']);

      const { canList } = usePermissions();
      expect(canList('assets').value).toBe(true);
      expect(canList('users').value).toBe(false);
    });

    it('canCreate checks resource.create', () => {
      seedPermissions(['assets.create']);

      const { canCreate } = usePermissions();
      expect(canCreate('assets').value).toBe(true);
      expect(canCreate('users').value).toBe(false);
    });

    it('canRead checks resource.read', () => {
      seedPermissions(['assets.read']);

      const { canRead } = usePermissions();
      expect(canRead('assets').value).toBe(true);
      expect(canRead('users').value).toBe(false);
    });

    it('canUpdate checks resource.update', () => {
      seedPermissions(['assets.update']);

      const { canUpdate } = usePermissions();
      expect(canUpdate('assets').value).toBe(true);
      expect(canUpdate('users').value).toBe(false);
    });

    it('canDelete checks resource.delete', () => {
      seedPermissions(['assets.delete']);

      const { canDelete } = usePermissions();
      expect(canDelete('assets').value).toBe(true);
      expect(canDelete('users').value).toBe(false);
    });
  });

  describe('wildcard support', () => {
    it('mapex.* grants all permissions', () => {
      seedPermissions(['mapex.*']);

      const { hasPermission, canCreate, canDelete } = usePermissions();
      expect(hasPermission('anything.here').value).toBe(true);
      expect(canCreate('assets').value).toBe(true);
      expect(canDelete('users').value).toBe(true);
    });

    it('resource.* grants all actions on that resource', () => {
      seedPermissions(['assets.*']);

      const { canList, canCreate, canRead, canUpdate, canDelete } = usePermissions();
      expect(canList('assets').value).toBe(true);
      expect(canCreate('assets').value).toBe(true);
      expect(canRead('assets').value).toBe(true);
      expect(canUpdate('assets').value).toBe(true);
      expect(canDelete('assets').value).toBe(true);
      expect(canList('users').value).toBe(false);
    });
  });

  describe('loading and loaded state', () => {
    it('permissionsLoading reflects store.loading', () => {
      const store = usePermissionStore();
      store.$patch({ loading: true });

      const { permissionsLoading } = usePermissions();
      expect(permissionsLoading.value).toBe(true);

      store.$patch({ loading: false });
      expect(permissionsLoading.value).toBe(false);
    });

    it('permissionsLoaded reflects store.isLoaded', () => {
      const store = usePermissionStore();
      store.$patch({ permissions: ['a.b'], loading: false });

      const { permissionsLoaded } = usePermissions();
      expect(permissionsLoaded.value).toBe(true);
    });

    it('permissionsLoaded is false when loading', () => {
      const store = usePermissionStore();
      store.$patch({ permissions: ['a.b'], loading: true });

      const { permissionsLoaded } = usePermissions();
      expect(permissionsLoaded.value).toBe(false);
    });

    it('permissionsLoaded is false when permissions are empty', () => {
      const store = usePermissionStore();
      store.$patch({ permissions: [], loading: false });

      const { permissionsLoaded } = usePermissions();
      expect(permissionsLoaded.value).toBe(false);
    });
  });
});

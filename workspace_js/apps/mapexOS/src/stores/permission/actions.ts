import type { PermissionState } from './types';
import { apis } from '@src/services/mapex';
import { useOrganizationStore } from '@stores/organization';
import { useLogger } from '@composables/useLogger';

const logger = useLogger('PermissionStore');

export const actions = {
  /**
   * Fetch permissions from backend for the current organization context.
   * Skips if same org, data is fresh (<5 min), and not forced.
   *
   * @param {boolean} force - Skip cache and always fetch
   */
  async fetchPermissions(force = false) {
    const store = this as PermissionState & typeof actions;
    const orgStore = useOrganizationStore();
    const orgId = orgStore.selectedOrganizationId;

    // Check staleness
    const isStale = !store.lastFetched ||
      (new Date(store.lastFetched).getTime() < Date.now() - 5 * 60 * 1000);

    // Skip if same org, not stale, and not forced
    if (!force && store.forOrganizationId === orgId && !isStale && store.permissions.length > 0) {
      return;
    }

    store.loading = true;
    store.error = null;

    try {
      const data = await apis.mapexOS?.auth.getMyPermissions();

      if (!data) {
        throw new Error('No permission data received');
      }

      store.permissions = data.permissions;
      store.version = data.version;
      store.forOrganizationId = orgId;
      store.lastFetched = new Date().toISOString();

      logger.info(`Permissions loaded: ${data.permissions.length} permissions, version=${data.version}, org=${orgId}`);
    } catch (error) {
      logger.error('Error fetching permissions', error);
      store.error = error instanceof Error ? error.message : 'Failed to fetch permissions';
    } finally {
      store.loading = false;
    }
  },

  /**
   * Clear all permission data (on logout or org context loss)
   */
  clearPermissions() {
    const store = this as PermissionState & typeof actions;
    store.permissions = [];
    store.version = 0;
    store.loading = false;
    store.error = null;
    store.forOrganizationId = null;
    store.lastFetched = null;
  },
};

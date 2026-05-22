import type { OrganizationState, OrganizationCoverageItem, OrganizationType } from './types';
import { apis } from '@src/services/mapex';
import { buildOrganizationTree, getParentPathKey } from '@utils/organization/treeBuilder';
import { usePermissionStore } from '@stores/permission';
import { useLogger } from '@composables/useLogger';

const logger = useLogger('OrganizationStore');

/**
 * Select initial organization with intelligent strategy
 * Priority order:
 * 1. If user has only 1 org → auto-select it
 * 2. Restore last selected org from localStorage (if user still has access)
 * 3. Prioritize org with scope "recursive" (vendor)
 * 4. Fallback: first vendor (type === 'vendor')
 * 5. Last fallback: first org in list
 */
function selectInitialOrganization(organizations: OrganizationCoverageItem[]): string {
  if (organizations.length === 0) {
    throw new Error('No organizations available');
  }

  // 1. If only 1 org → auto-select (simplest case)
  if (organizations.length === 1) {
    return organizations[0]!.id;
  }

  // 2. Restore last selected org from localStorage (HIGHEST PRIORITY after single org)
  const lastOrgId = localStorage.getItem('selectedOrgId');
  if (lastOrgId) {
    const hasAccess = organizations.some(org => org.id === lastOrgId);
    if (hasAccess) {
      return lastOrgId;
    }
    // User lost access to previously selected org, clear localStorage
    localStorage.removeItem('selectedOrgId');
    localStorage.removeItem('selectedOrgName');
  }

  // 3. Prioritize org with scope "recursive" (usually vendor)
  const recursiveOrg = organizations.find(org => org.scope === 'recursive');
  if (recursiveOrg) {
    return recursiveOrg.id;
  }

  // 4. Fallback: first vendor
  const vendor = organizations.find(org => org.type === 'vendor');
  if (vendor) {
    return vendor.id;
  }

  // 5. Last fallback: first org in list
  return organizations[0]!.id;
}

export const actions = {
  /**
   * Fetch user's organization coverage from API
   * Builds the tree structure and stores in state
   */
  async fetchCoverage(force = false) {
    const store = this as OrganizationState & typeof actions;

    // Check if data is stale (> 5 minutes)
    const isStale = !store.lastUpdated ||
      (new Date(store.lastUpdated).getTime() < Date.now() - 5 * 60 * 1000);

    // Skip if data is fresh and not forcing
    if (!force && !isStale && store.coverage) {
      return;
    }

    store.loading = true;
    store.error = null;

    try {
      const data = await apis.mapexOS?.auth.getUserCoverage();

      if (!data) {
        throw new Error('No coverage data received');
      }

      if (!data.organizations || data.organizations.length === 0) {
        throw new Error('User has no organization access');
      }

      // Store raw coverage data
      store.coverage = data;
      store.flatList = data.organizations;
      store.lastUpdated = data.lastUpdated;

      // Build tree structure using utility function
      store.treeNodes = buildOrganizationTree(data.organizations);

    } catch (error) {
      logger.error('Error fetching organization coverage', error);
      store.error = error instanceof Error ? error.message : 'Failed to fetch coverage';
      throw error;
    } finally {
      store.loading = false;
    }
  },

  /**
   * Initialize after user login
   * Fetches coverage and selects initial organization
   */
  async initializeAfterLogin() {
    const store = this as OrganizationState & typeof actions;

    // 1. Fetch coverage
    await store.fetchCoverage(true); // Force fresh data after login

    // 2. Select initial organization intelligently
    const selectedId = selectInitialOrganization(store.flatList);

    // 3. Set as selected
    store.selectOrganization(selectedId);

    // 4. Fetch permissions for the selected organization
    const permStore = usePermissionStore();
    await permStore.fetchPermissions(true);

    // 5. Save preference to localStorage
    localStorage.setItem('selectedOrgId', selectedId);

    return selectedId;
  },

  /**
   * Refresh coverage data (forces new fetch)
   */
  async refreshCoverage() {
    return this.fetchCoverage(true);
  },

  /**
   * Select an organization (set context)
   * This is the main way to switch between organizations
   */
  selectOrganization(orgId: string | null) {
    const store = this as OrganizationState & typeof actions;

    // Validate org exists in flatList
    if (orgId && !store.flatList.some(org => org.id === orgId)) {
      logger.warn(`Organization ${orgId} not found in user's coverage`);
      return;
    }

    // Get organization details
    const org = orgId ? store.flatList.find(o => o.id === orgId) : null;

    // Update state with both id and name
    store.selectedOrganizationId = orgId;
    store.selectedOrganizationName = org?.name || null;

    // Save to localStorage
    if (orgId) {
      localStorage.setItem('selectedOrgId', orgId);
      if (org?.name) {
        localStorage.setItem('selectedOrgName', org.name);
      }
    } else {
      localStorage.removeItem('selectedOrgId');
      localStorage.removeItem('selectedOrgName');
    }

    // Refresh permissions for the new organization context
    if (orgId) {
      const permStore = usePermissionStore();
      void permStore.fetchPermissions(true);
    }
  },

  /**
   * Add a newly created organization to the tree without re-fetching from backend.
   * Inherits membership data (scope, membershipId, roleIds) from the parent org.
   *
   * @param org - The created organization data (from API response)
   */
  addOrganizationToTree(org: { id: string; name: string; type: string; pathKey: string }) {
    const store = this as OrganizationState & typeof actions;

    // Find parent in flatList to inherit membership data
    const parentPathKey = getParentPathKey(org.pathKey);
    if (!parentPathKey) return;

    const parent = store.flatList.find(o => o.pathKey === parentPathKey);
    if (!parent) {
      logger.warn(`Parent org with pathKey=${parentPathKey} not found in coverage — skipping local tree update`);
      return;
    }

    // Create coverage item inheriting parent's membership context
    const newItem: OrganizationCoverageItem = {
      id: org.id,
      name: org.name,
      type: org.type as OrganizationType,
      pathKey: org.pathKey,
      scope: parent.scope,
      membershipId: parent.membershipId,
      roleIds: [...parent.roleIds],
    };

    // Add to flatList and rebuild tree
    store.flatList.push(newItem);
    store.treeNodes = buildOrganizationTree(store.flatList);

    logger.debug(`Added org=${org.id} (${org.type}) to tree under parent pathKey=${parentPathKey}`);
  },

  /**
   * Clear all organization data (on logout)
   */
  clearCoverage() {
    const store = this as OrganizationState & typeof actions;
    store.coverage = null;
    store.treeNodes = [];
    store.flatList = [];
    store.selectedOrganizationId = null;
    store.selectedOrganizationName = null;
    store.lastUpdated = null;
    store.error = null;

    // Clear localStorage
    localStorage.removeItem('selectedOrgId');
    localStorage.removeItem('selectedOrgName');

    // Clear permission cache
    const permStore = usePermissionStore();
    permStore.clearPermissions();
  },
};

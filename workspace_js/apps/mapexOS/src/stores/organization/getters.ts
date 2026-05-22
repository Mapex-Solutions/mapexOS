import type { OrganizationState, OrganizationTreeNode, OrganizationType } from './types';

export const getters = {
  /**
   * Get organization by ID
   */
  getOrganizationById: (state: OrganizationState) => (id: string) => {
    return state.flatList.find(org => org.id === id);
  },

  /**
   * Get organizations by type
   */
  getOrganizationsByType: (state: OrganizationState) => (type: OrganizationType) => {
    return state.flatList.filter(org => org.type === type);
  },

  /**
   * Get vendor (root) organizations
   */
  vendors: (state: OrganizationState): OrganizationTreeNode[] => {
    return state.treeNodes.filter(org => org.type === 'vendor');
  },

  /**
   * Check if user has access to specific organization
   */
  hasAccessTo: (state: OrganizationState) => (orgId: string): boolean => {
    return state.flatList.some(org => org.id === orgId);
  },

  /**
   * Get selected organization details
   */
  selectedOrganization: (state: OrganizationState) => {
    if (!state.selectedOrganizationId) return null;
    return state.flatList.find(org => org.id === state.selectedOrganizationId);
  },

  /**
   * Check if data is stale (> 5 minutes)
   */
  isDataStale: (state: OrganizationState): boolean => {
    if (!state.lastUpdated) return true;
    const fiveMinutesAgo = Date.now() - 5 * 60 * 1000;
    return new Date(state.lastUpdated).getTime() < fiveMinutesAgo;
  },

  /**
   * Check if org is currently selected
   */
  isSelected: (state: OrganizationState) => (orgId: string): boolean => {
    return state.selectedOrganizationId === orgId;
  },

  /**
   * Get total count of organizations
   */
  totalCount: (state: OrganizationState): number => {
    return state.flatList.length;
  },

  /**
   * Check if selected organization is a vendor
   */
  isVendor: (state: OrganizationState): boolean => {
    if (!state.selectedOrganizationId) return false;
    const org = state.flatList.find(org => org.id === state.selectedOrganizationId);
    return org?.type === 'vendor';
  },

  /**
   * Check if selected organization is a customer
   */
  isCustomer: (state: OrganizationState): boolean => {
    if (!state.selectedOrganizationId) return false;
    const org = state.flatList.find(org => org.id === state.selectedOrganizationId);
    return org?.type === 'customer';
  },

  /**
   * Check if selected organization is a site or deeper (site, building, floor, zone)
   */
  isSite: (state: OrganizationState): boolean => {
    if (!state.selectedOrganizationId) return false;
    const org = state.flatList.find(org => org.id === state.selectedOrganizationId);
    return org?.type === 'site' || org?.type === 'building' || org?.type === 'floor' || org?.type === 'zone';
  },
};

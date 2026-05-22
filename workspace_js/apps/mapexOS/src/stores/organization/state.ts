import type { OrganizationState } from './types';

export const state = (): OrganizationState => ({
  coverage: null,
  treeNodes: [],
  flatList: [],

  loading: false,
  error: null,
  lastUpdated: null,

  selectedOrganizationId: null,
  selectedOrganizationName: null,
});

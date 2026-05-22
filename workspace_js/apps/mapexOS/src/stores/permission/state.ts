import type { PermissionState } from './types';

export const state = (): PermissionState => ({
  permissions: [],
  version: 0,
  loading: false,
  error: null,
  forOrganizationId: null,
  lastFetched: null,
});

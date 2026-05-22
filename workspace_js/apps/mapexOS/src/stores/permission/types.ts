/**
 * Permission store state interface
 */
export interface PermissionState {
  /** Resolved permission strings for the current user/org context */
  permissions: string[];

  /** Cache version from backend (round-robin 1-100) */
  version: number;

  /** Loading state while fetching permissions */
  loading: boolean;

  /** Error message if fetch fails */
  error: string | null;

  /** Organization ID these permissions were fetched for */
  forOrganizationId: string | null;

  /** ISO timestamp of last successful fetch (staleness check) */
  lastFetched: string | null;
}

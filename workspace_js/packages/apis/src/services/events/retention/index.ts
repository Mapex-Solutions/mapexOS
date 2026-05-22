import type { AxiosInstance } from 'axios';
import type { GetToken } from '@src/common';

import { userApi } from './user.api';

/**
 * Creates Retention module API with user endpoints.
 *
 * User methods (JWT auth) are spread to the root level.
 *
 * @param http - Axios instance for HTTP requests
 * @param getToken - Function to retrieve JWT token
 * @returns Object with user API methods
 *
 * @example
 * const retentionApi = createRetentionModuleApi(http, getToken);
 *
 * // User methods (JWT) - spread to root
 * await retentionApi.listRetentionPolicies({ page: 1, perPage: 20 });
 * await retentionApi.upsertRetentionPolicy({ type: 'events', name: 'Events', retentionDays: 90 });
 */
export function createRetentionModuleApi(http: AxiosInstance, getToken: GetToken | undefined) {
	return {
		...userApi(http, getToken),  // Spread user methods to root
	};
}

export type RetentionModuleApiMethods = ReturnType<typeof createRetentionModuleApi>;

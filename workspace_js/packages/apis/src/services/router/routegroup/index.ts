import type { AxiosInstance } from 'axios';
import type { GetToken } from '@src/common';

import { userApi } from './user.api';

/**
 * Creates RouteGroup module API with user endpoints.
 *
 * User methods (JWT auth) are spread to the root level.
 *
 * @param http - Axios instance for HTTP requests
 * @param getToken - Function to retrieve JWT token
 * @returns Object with user API methods
 *
 * @example
 * const routeGroupApi = createRouteGroupApi(http, getToken);
 *
 * // User methods (JWT) - spread to root
 * await routeGroupApi.list({ page: 1, perPage: 20 });
 * await routeGroupApi.create({ name: 'Route Group 1', version: '1.0.0', ... });
 * await routeGroupApi.getById({ routeGroupId: '...' });
 */
export function createRouteGroupApi(http: AxiosInstance, getToken: GetToken | undefined) {
	return {
		...userApi(http, getToken),  // Spread user methods to root
	};
}

export type RouteGroupApiMethods = ReturnType<typeof createRouteGroupApi>;

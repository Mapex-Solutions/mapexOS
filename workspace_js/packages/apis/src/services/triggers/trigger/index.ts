import type { AxiosInstance } from 'axios';
import type { GetToken } from '@src/common';

import { userApi } from './user.api';

/**
 * Creates Trigger module API with user endpoints.
 *
 * User methods (JWT auth) are spread to the root level.
 *
 * @param http - Axios instance for HTTP requests
 * @param getToken - Function to retrieve JWT token
 * @returns Object with user API methods
 *
 * @example
 * const triggerApi = createTriggerApi(http, getToken);
 *
 * // User methods (JWT) - spread to root
 * await triggerApi.list({ page: 1, perPage: 20 });
 * await triggerApi.create({ name: 'HTTP Trigger', triggerType: 'http', ... });
 * await triggerApi.getById({ triggerId: '...' });
 */
export function createTriggerApi(http: AxiosInstance, getToken: GetToken | undefined) {
	return {
		...userApi(http, getToken),  // Spread user methods to root
	};
}

export type TriggerApiMethods = ReturnType<typeof createTriggerApi>;

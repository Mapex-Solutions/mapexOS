import type { AxiosInstance } from 'axios';
import type { GetToken } from '@src/common';

import { userApi } from './user.api';

/**
 * Creates Events module API with user endpoints.
 *
 * User methods (JWT auth) are spread to the root level.
 *
 * @param http - Axios instance for HTTP requests
 * @param getToken - Function to retrieve JWT token
 * @returns Object with user API methods
 *
 * @example
 * const eventsApi = createEventsModuleApi(http, getToken);
 *
 * // User methods (JWT) - spread to root
 * await eventsApi.listRaw({ page: 1, perPage: 20 });
 * await eventsApi.listRaw({ threadId: 'datasource-id', startTime: '2024-01-01T00:00:00Z' });
 */
export function createEventsModuleApi(http: AxiosInstance, getToken: GetToken | undefined) {
	return {
		...userApi(http, getToken),  // Spread user methods to root
	};
}

export type EventsModuleApiMethods = ReturnType<typeof createEventsModuleApi>;

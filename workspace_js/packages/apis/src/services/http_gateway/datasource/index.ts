import type { AxiosInstance } from 'axios';
import type { GetToken } from '@src/common';

import { userApi } from './user.api';

/**
 * Creates DataSource module API with user endpoints.
 *
 * User methods (JWT auth) are spread to the root level.
 *
 * @param http - Axios instance for HTTP requests
 * @param getToken - Function to retrieve JWT token
 * @returns Object with user API methods
 *
 * @example
 * const dataSourceApi = createDataSourceApi(http, getToken);
 *
 * // User methods (JWT) - spread to root
 * await dataSourceApi.list({ page: 1, perPage: 20 });
 * await dataSourceApi.create({ name: 'Data Source 1', mode: 'push', ... });
 * await dataSourceApi.getById({ dataSourceId: '...' });
 */
export function createDataSourceApi(http: AxiosInstance, getToken: GetToken | undefined) {
	return {
		...userApi(http, getToken),  // Spread user methods to root
	};
}

export type DataSourceApiMethods = ReturnType<typeof createDataSourceApi>;

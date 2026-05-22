import type { AxiosInstance } from 'axios';
import type { GetToken } from '@src/common';

import { userApi } from './user.api';

/**
 * Creates List module API.
 *
 * List methods (JWT auth) are spread to the root level.
 *
 * @param http - Axios instance for HTTP requests
 * @param getToken - Function to retrieve JWT token
 * @returns Object with list API methods
 *
 * @example
 * const listApi = createListApi(http, getToken);
 *
 * // List methods (JWT)
 * await listApi.list({ page: 1, perPage: 20, type: 'assetType' });
 * await listApi.create({ type: 'assetType', name: 'Temperature Sensor', value: 'temp_sensor', isSystem: false, orgId: '...' });
 * await listApi.getById({ listId: '...' });
 * await listApi.update({ listId: '...' }, { name: 'Updated Name' });
 * await listApi.delete({ listId: '...' });
 */
export function createListApi(http: AxiosInstance, getToken: GetToken | undefined) {
	return {
		...userApi(http, getToken),
	};
}

export type ListApiMethods = ReturnType<typeof createListApi>;

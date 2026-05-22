import type { AxiosInstance } from 'axios';
import type { GetToken } from '@src/common';

import { userApi } from './user.api';

/**
 * Creates Group module API.
 *
 * Group methods (JWT auth) are spread to the root level.
 *
 * @param http - Axios instance for HTTP requests
 * @param getToken - Function to retrieve JWT token
 * @returns Object with group API methods
 *
 * @example
 * const groupApi = createGroupApi(http, getToken);
 *
 * // Group methods (JWT)
 * await groupApi.list({ page: 1, perPage: 20 });
 * await groupApi.create({ name: 'Admins', enabled: true, isSystem: false, orgId: '...' });
 * await groupApi.getById({ groupId: '...' });
 * await groupApi.update({ groupId: '...' }, { name: 'New Name' });
 * await groupApi.delete({ groupId: '...' });
 */
export function createGroupApi(http: AxiosInstance, getToken: GetToken | undefined) {
	return {
		...userApi(http, getToken),
	};
}

export type GroupApiMethods = ReturnType<typeof createGroupApi>;

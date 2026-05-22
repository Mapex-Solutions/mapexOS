import type { AxiosInstance } from 'axios';
import type { GetToken } from '@src/common';

import { userApi } from './user.api';

/**
 * Creates Role module API.
 *
 * Role methods (JWT auth) are spread to the root level.
 *
 * @param http - Axios instance for HTTP requests
 * @param getToken - Function to retrieve JWT token
 * @returns Object with role API methods
 *
 * @example
 * const roleApi = createRoleApi(http, getToken);
 *
 * // Role methods (JWT)
 * await roleApi.list({ page: 1, perPage: 20 });
 * await roleApi.create({ name: 'Admin', permissions: ['user:read', 'user:write'], isSystem: false, orgId: '...', scope: 'global' });
 * await roleApi.getById({ roleId: '...' });
 * await roleApi.update({ roleId: '...' }, { name: 'Super Admin' });
 * await roleApi.delete({ roleId: '...' });
 */
export function createRoleApi(http: AxiosInstance, getToken: GetToken | undefined) {
	return {
		...userApi(http, getToken),
	};
}

export type RoleApiMethods = ReturnType<typeof createRoleApi>;

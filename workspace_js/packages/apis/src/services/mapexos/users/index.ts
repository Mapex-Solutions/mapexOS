import type { AxiosInstance } from 'axios';
import type { GetToken } from '@src/common';

import { userApi } from './user.api';

/**
 * Creates User module API.
 *
 * User methods (JWT auth) are spread to the root level.
 *
 * @param http - Axios instance for HTTP requests
 * @param getToken - Function to retrieve JWT token
 * @returns Object with user API methods
 *
 * @example
 * const userApi = createUserApi(http, getToken);
 *
 * // User methods (JWT)
 * await userApi.me();
 * await userApi.updateMe({ firstName: 'John' });
 * await userApi.list({ page: 1, perPage: 20 });
 * await userApi.create({ email: 'user@example.com', ... });
 * await userApi.getById({ userId: '...' });
 * await userApi.update({ userId: '...' }, { firstName: 'Jane' });
 * await userApi.delete({ userId: '...' });
 */
export function createUserApi(http: AxiosInstance, getToken: GetToken | undefined) {
	return {
		...userApi(http, getToken),
	};
}

export type UserApiMethods = ReturnType<typeof createUserApi>;

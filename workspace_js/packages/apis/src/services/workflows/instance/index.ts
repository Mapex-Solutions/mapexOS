import type { AxiosInstance } from 'axios';
import type { GetToken } from '@src/common';

import { userApi } from './user.api';

/**
 * Creates Instance module API with user endpoints.
 *
 * User methods (JWT auth) are spread to the root level.
 *
 * @param http - Axios instance for HTTP requests
 * @param getToken - Function to retrieve JWT token
 * @returns Object with user API methods
 */
export function createInstanceApi(http: AxiosInstance, getToken: GetToken | undefined) {
	return {
		...userApi(http, getToken),
	};
}

export type InstanceApiMethods = ReturnType<typeof createInstanceApi>;

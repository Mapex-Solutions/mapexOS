import type { AxiosInstance } from 'axios';
import type { GetToken } from '@src/common';

import { userApi } from './scripts.api';

/**
 * Creates Scripts module API with user endpoints.
 *
 * User methods (JWT auth) are spread to the root level.
 *
 * @param http - Axios instance for HTTP requests
 * @param getToken - Function to retrieve JWT token
 * @returns Object with user API methods
 *
 * @example
 * const scriptsApi = createScriptsApi(http, getToken);
 *
 * // User methods (JWT) - spread to root
 * await scriptsApi.test({ decode: '...', validation: '...', transform: '...', event: {...} });
 */
export function createScriptsApi(http: AxiosInstance, getToken: GetToken | undefined) {
	return {
		...userApi(http, getToken),  // Spread user methods to root
	};
}

export type ScriptsApiMethods = ReturnType<typeof createScriptsApi>;

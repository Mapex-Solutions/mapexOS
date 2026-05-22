import type { AxiosInstance } from 'axios';
import type { GetToken, SessionsConfig } from '@src/common';

import { authApi } from './user.api';

/**
 * Creates a user API instance that provides access to user-related endpoints.
 *
 *
 * @param http - The Axios instance used for making HTTP requests
 * @param getToken - A function that retrieves authentication tokens for authorized requests
 * @returns An object containing user API methods, internal API methods, admin API methods,
 *          and the original HTTP client instance
 */
export function createAuthApi(
	http: AxiosInstance,
	getToken: GetToken | undefined,
	sessionsConfig?: SessionsConfig
) {

	return {
		...authApi(http, getToken, sessionsConfig),
	};
}

/**
 * Interface for API initialization configuration.
 */
export type AuthApiModules = ReturnType<typeof createAuthApi>;


import type { Login, LoginResponse, OrganizationCoverageResponse, PermissionsResponse } from '@mapexos/schemas';

import { AxiosInstance } from 'axios';
import { ZodLoginSchema } from '@mapexos/schemas';

import { createApiFactory, GetToken, SessionsConfig } from '@src/common';

/**
 * Creates an authentication API with predefined methods for login, logout, and token refresh.
 *
 * @param http - An instance of Axios used for making HTTP requests.
 * @param getToken - A function to retrieve the current authentication token, or undefined if not needed.
 * @param sessionsConfig - Optional configuration for session management.
 * @returns An object containing methods for authentication operations.
 */
export function authApi(
	http: AxiosInstance,
	getToken: GetToken | undefined,
	sessionsConfig?: SessionsConfig
) {
	const factory = createApiFactory(http, sessionsConfig);

	return factory({
		basePath: '/auth',
		useAuthJWT: true,
		getToken,
		methods: {
			login: {
				method: 'POST',
				path: '/login',
				bodyParams: {} as Login,
				bodySchema: ZodLoginSchema,
				responseType: {} as LoginResponse
			},
			logout: {
				method: 'POST',
				path: '/logout',
			},
			refreshToken: {
				method: 'POST',
				path: '/refresh',
			},
			getUserCoverage: {
				method: 'GET',
				path: '/users/me/coverage',
				responseType: {} as OrganizationCoverageResponse
			},
			getMyPermissions: {
				method: 'GET',
				path: '/me/permissions',
				responseType: {} as PermissionsResponse
			}
		},
	});
}

/**
 * Interface for API initialization configuration.
 */
export type AuthApiMethods = ReturnType<typeof authApi>;
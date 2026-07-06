import type { AxiosInstance } from 'axios';
import type { GetToken } from '@src/common';
import { userApi } from './user.api';

/**
 * Creates the complete Template Migration module API.
 *
 * Like assetTemplate, it only exposes user-facing endpoints, so the methods
 * are returned directly without nesting.
 *
 * @param http - Axios instance for HTTP requests
 * @param getToken - Function to retrieve JWT token
 * @returns Object containing template migration API methods
 */
export function createMigrationApi(http: AxiosInstance, getToken: GetToken | undefined) {
	return userApi(http, getToken);
}

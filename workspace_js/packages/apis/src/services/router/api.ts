import { isEmpty } from 'lodash';
import { ApiConfig, ApiInterceptors, SessionsConfig, GetToken } from '@src/common';
import { createHttp } from '@src/tools';

import { createRouteGroupApi } from './routegroup';

/**
 * Creates complete Router Service API with all modules.
 *
 * Currently includes:
 * - routegroup: RouteGroup CRUD operations
 *
 * @param config - API configuration (baseURL, timeout, headers, etc.)
 * @param sessionsConfig - Optional session configuration (interceptors, getToken)
 * @returns Object with HTTP instance and all module APIs
 *
 * @example
 * const routerApi = createRouterApi({
 *   baseURL: 'http://localhost:5003',
 *   timeout: 30000
 * }, {
 *   getToken: () => localStorage.getItem('token'),
 *   interceptors: myInterceptors
 * });
 *
 * // Usage:
 * await routerApi.routegroup.list({ page: 1, perPage: 20 });
 * await routerApi.routegroup.create({ name: 'Route Group 1', version: '1.0.0', ... });
 */
export function createRouterApi(config: ApiConfig, sessionsConfig?: SessionsConfig) {
	/** Global/Local interceptors for the router service */
	const interceptors = !isEmpty(sessionsConfig?.interceptors)
		? sessionsConfig?.interceptors
		: config?.interceptors || {} as ApiInterceptors;

	const getToken: GetToken | undefined = sessionsConfig?.getToken;

	/** Create single HTTP instance for the entire service */
	const http = createHttp(config, interceptors);

	/** Initialize all modules with the same HTTP instance */
	const routegroup = createRouteGroupApi(http, getToken);

	return {
		http,       // Expose HTTP instance
		routegroup, // RouteGroup module (CRUD operations)
	};
}

/**
 * Type for the complete Router Service API
 */
export type RouterApi = ReturnType<typeof createRouterApi>;

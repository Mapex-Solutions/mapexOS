import { isEmpty } from 'lodash';
import { ApiConfig, ApiInterceptors, SessionsConfig, GetToken } from '@src/common';
import { createHttp } from '@src/tools';

import { createDataSourceApi } from './datasource';
import { createHeartbeatApi } from './heartbeat';

/**
 * Creates complete HTTP Gateway Service API with all modules.
 *
 * Currently includes:
 * - datasource: DataSource CRUD operations
 * - heartbeat:  HTTP-protocol asset liveness ping (per-DataSource auth)
 *
 * @param config - API configuration (baseURL, timeout, headers, etc.)
 * @param sessionsConfig - Optional session configuration (interceptors, getToken)
 * @returns Object with HTTP instance and all module APIs
 *
 * @example
 * const httpGatewayApi = createHttpGatewayApi({
 *   baseURL: 'http://localhost:5001',
 *   timeout: 30000
 * }, {
 *   getToken: () => localStorage.getItem('token'),
 *   interceptors: myInterceptors
 * });
 *
 * // Usage:
 * await httpGatewayApi.datasource.list({ page: 1, perPage: 20 });
 * await httpGatewayApi.heartbeat.post('ds-id', { assetUUID: 'a' }, { headers: { 'x-api-key': 'k' } });
 */
export function createHttpGatewayApi(config: ApiConfig, sessionsConfig?: SessionsConfig) {
	/** Global/Local interceptors for the http_gateway service */
	const interceptors = !isEmpty(sessionsConfig?.interceptors)
		? sessionsConfig?.interceptors
		: config?.interceptors || {} as ApiInterceptors;

	const getToken: GetToken | undefined = sessionsConfig?.getToken;

	/** Create single HTTP instance for the entire service */
	const http = createHttp(config, interceptors);

	/** Initialize all modules with the same HTTP instance */
	const datasource = createDataSourceApi(http, getToken);
	const heartbeat = createHeartbeatApi(http);

	return {
		http,       // Expose HTTP instance
		datasource, // DataSource module (CRUD operations)
		heartbeat,  // Heartbeat module (per-DataSource auth)
	};
}

/**
 * Type for the complete HTTP Gateway Service API
 */
export type HttpGatewayApi = ReturnType<typeof createHttpGatewayApi>;

import { isEmpty } from 'lodash';
import { ApiConfig, ApiInterceptors, SessionsConfig, GetToken } from '@src/common';
import { createHttp } from '@src/tools';

import { createScriptsApi } from './scripts';

/**
 * Creates complete JSExecutor Service API with all modules.
 *
 * Currently includes:
 * - scripts: Script testing operations
 *
 * @param config - API configuration (baseURL, timeout, headers, etc.)
 * @param sessionsConfig - Optional session configuration (interceptors, getToken)
 * @returns Object with HTTP instance and all module APIs
 *
 * @example
 * const jsExecutorApi = createJsExecutorApi({
 *   baseURL: 'http://localhost:5003',
 *   timeout: 30000
 * }, {
 *   getToken: () => localStorage.getItem('token'),
 *   interceptors: myInterceptors
 * });
 *
 * // Usage:
 * await jsExecutorApi.scripts.test({ decode: '...', validation: '...', transform: '...', event: {...} });
 */
export function createJsExecutorApi(config: ApiConfig, sessionsConfig?: SessionsConfig) {
	/** Global/Local interceptors for the js_executor service */
	const interceptors = !isEmpty(sessionsConfig?.interceptors)
		? sessionsConfig?.interceptors
		: config?.interceptors || {} as ApiInterceptors;

	const getToken: GetToken | undefined = sessionsConfig?.getToken;

	/** Create single HTTP instance for the entire service */
	const http = createHttp(config, interceptors);

	/** Initialize all modules with the same HTTP instance */
	const scripts = createScriptsApi(http, getToken);

	return {
		http,       // Expose HTTP instance
		scripts,    // Scripts module (test operations)
	};
}

/**
 * Type for the complete JSExecutor Service API
 */
export type JsExecutorApi = ReturnType<typeof createJsExecutorApi>;

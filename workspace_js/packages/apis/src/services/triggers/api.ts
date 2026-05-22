import { isEmpty } from 'lodash';
import { ApiConfig, ApiInterceptors, SessionsConfig, GetToken } from '@src/common';
import { createHttp } from '@src/tools';

import { createTriggerApi } from './trigger';

/**
 * Creates complete Triggers Service API with all modules.
 *
 * Currently includes:
 * - trigger: Trigger CRUD operations
 *
 * @param config - API configuration (baseURL, timeout, headers, etc.)
 * @param sessionsConfig - Optional session configuration (interceptors, getToken)
 * @returns Object with HTTP instance and all module APIs
 *
 * @example
 * const triggersApi = createTriggersApi({
 *   baseURL: 'http://localhost:5004',
 *   timeout: 30000
 * }, {
 *   getToken: () => localStorage.getItem('token'),
 *   interceptors: myInterceptors
 * });
 *
 * // Usage:
 * await triggersApi.trigger.list({ page: 1, perPage: 20 });
 * await triggersApi.trigger.create({ name: 'HTTP Webhook', triggerType: 'http', ... });
 */
export function createTriggersApi(config: ApiConfig, sessionsConfig?: SessionsConfig) {
	/** Global/Local interceptors for the triggers service */
	const interceptors = !isEmpty(sessionsConfig?.interceptors)
		? sessionsConfig?.interceptors
		: config?.interceptors || {} as ApiInterceptors;

	const getToken: GetToken | undefined = sessionsConfig?.getToken;

	/** Create single HTTP instance for the entire service */
	const http = createHttp(config, interceptors);

	/** Initialize all modules with the same HTTP instance */
	const trigger = createTriggerApi(http, getToken);

	return {
		http,    // Expose HTTP instance
		trigger, // Trigger module (CRUD operations)
	};
}

/**
 * Type for the complete Triggers Service API
 */
export type TriggersApi = ReturnType<typeof createTriggersApi>;

import { isEmpty } from 'lodash';
import { ApiConfig, ApiInterceptors, SessionsConfig, GetToken } from '@src/common';
import { createHttp } from '@src/tools';

import { createEventsModuleApi } from './events';
import { createRetentionModuleApi } from './retention';

/**
 * Creates complete Events Service API with all modules.
 *
 * Currently includes:
 * - events: Raw events query operations (ClickHouse)
 * - retention: Retention policy CRUD operations (MongoDB)
 *
 * @param config - API configuration (baseURL, timeout, headers, etc.)
 * @param sessionsConfig - Optional session configuration (interceptors, getToken)
 * @returns Object with HTTP instance and all module APIs
 *
 * @example
 * const eventsApi = createEventsApi({
 *   baseURL: 'http://localhost:5004',
 *   timeout: 30000
 * }, {
 *   getToken: () => localStorage.getItem('token'),
 *   interceptors: myInterceptors
 * });
 *
 * // Usage:
 * await eventsApi.events.listRaw({ page: 1, perPage: 20 });
 * await eventsApi.events.listRaw({ threadId: 'ds-123', startTime: '2024-01-01T00:00:00Z' });
 */
export function createEventsApi(config: ApiConfig, sessionsConfig?: SessionsConfig) {
	/** Global/Local interceptors for the events service */
	const interceptors = !isEmpty(sessionsConfig?.interceptors)
		? sessionsConfig?.interceptors
		: config?.interceptors || {} as ApiInterceptors;

	const getToken: GetToken | undefined = sessionsConfig?.getToken;

	/** Create single HTTP instance for the entire service */
	const http = createHttp(config, interceptors);

	/** Initialize all modules with the same HTTP instance */
	const events = createEventsModuleApi(http, getToken);
	const retention = createRetentionModuleApi(http, getToken);

	return {
		http,       // Expose HTTP instance
		events,     // Events module (raw events queries)
		retention,  // Retention module (retention policy CRUD)
	};
}

/**
 * Type for the complete Events Service API
 */
export type EventsApi = ReturnType<typeof createEventsApi>;

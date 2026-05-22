import { isEmpty } from 'lodash';
import { ApiConfig, ApiInterceptors, SessionsConfig, GetToken } from '@src/common';
import { createHttp } from '@src/tools';

import { createCredentialApi } from './credential';
import { createDefinitionApi } from './definition';
import { createExecutionApi } from './execution';
import { createInstanceApi } from './instance';
import { createPluginApi } from './plugin';

/**
 * Creates complete Workflows Service API with all modules.
 *
 * Includes:
 * - credential: Credential instance CRUD + test
 * - definition: Workflow definition CRUD operations
 * - execution: Workflow execution list/get/cancel/signal
 * - instance: Workflow instance config CRUD
 * - plugin: Plugin manifest CRUD + getEnabled
 *
 * @param config - API configuration (baseURL, timeout, headers, etc.)
 * @param sessionsConfig - Optional session configuration (interceptors, getToken)
 * @returns Object with HTTP instance and all module APIs
 *
 * @example
 * const workflowsApi = createWorkflowsApi({
 *   baseURL: 'http://localhost:5007',
 *   timeout: 30000
 * }, {
 *   getToken: () => localStorage.getItem('token'),
 *   interceptors: myInterceptors
 * });
 *
 * // Usage:
 * await workflowsApi.definition.list({ page: 1, perPage: 20 });
 * await workflowsApi.instance.create({ name: 'My Instance', definitionId: '...', ... });
 * await workflowsApi.execution.list({ instanceId: '...' });
 * await workflowsApi.execution.cancel({ executionId: '...' });
 */
export function createWorkflowsApi(config: ApiConfig, sessionsConfig?: SessionsConfig) {
	/** Global/Local interceptors for the workflows service */
	const interceptors = !isEmpty(sessionsConfig?.interceptors)
		? sessionsConfig?.interceptors
		: config?.interceptors || {} as ApiInterceptors;

	const getToken: GetToken | undefined = sessionsConfig?.getToken;

	/** Create single HTTP instance for the entire service */
	const http = createHttp(config, interceptors);

	/** Initialize all modules with the same HTTP instance */
	const credential = createCredentialApi(http, getToken);
	const definition = createDefinitionApi(http, getToken);
	const execution = createExecutionApi(http, getToken);
	const instance = createInstanceApi(http, getToken);
	const plugin = createPluginApi(http, getToken);

	return {
		http,       // Expose HTTP instance
		credential, // Credential module (CRUD + test)
		definition, // Definition module (CRUD operations)
		execution,  // Execution module (list/get/cancel/signal)
		instance,   // Instance module (config CRUD)
		plugin,     // Plugin module (manifest CRUD + getEnabled)
	};
}

/**
 * Type for the complete Workflows Service API
 */
export type WorkflowsApi = ReturnType<typeof createWorkflowsApi>;

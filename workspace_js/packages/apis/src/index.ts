
/**
 * Mapex APIs module
*/
import type { MapexOSApi, AssetsApi, EventsApi, RouterApi, TriggersApi, HttpGatewayApi, JsExecutorApi, WorkflowsApi, VaultApiMethods } from './services'

import { ApiInitConfig, SessionsConfig } from './common'
import { createMapexOSApi, createAssetsApi, createEventsApi, createRouterApi, createTriggersApi, createHttpGatewayApi, createJsExecutorApi, createWorkflowsApi, createVaultApi } from './services'

/**
 * Export all interfaces from the module.
 */
export * from './common'

/**
 * Initializes the Mapex APIs with the provided configuration.
 *
 * @example
 * const apiConfig = {
 *   interceptors: myInterceptors,
 *   getToken: myGetTokenFunction,
 *   notifications: notificationConfig,
 *   ...Others API configurations
 * };
 * const apis = initializeApis(apiConfig);
 * console.log(apis.mapexOS); // Access the mapexOS API
 *
 * @param {ApiInitConfig} params - The configuration parameters for initializing the APIs, including interceptors and token retrieval function.
 * @returns An object containing the initialized APIs based on the provided configuration.
 */
export function initializeApis(params: ApiInitConfig) {
  const { interceptors, getToken } = params
  const sessionParams = { interceptors, getToken } as SessionsConfig

  /** Export the initialized API's */
  const apis: {
    mapexOS?: MapexOSApi
    assets?: AssetsApi
    events?: EventsApi
    router?: RouterApi
    triggers?: TriggersApi
    httpGateway?: HttpGatewayApi
    jsExecutor?: JsExecutorApi,
    workflows?: WorkflowsApi,
    vault?: VaultApiMethods,
    /** @deprecated Rule engine removed - stub for compilation */
    rules?: any,
    /** @deprecated Rule engine removed - stub for compilation */
    businessRules?: any,
  } = {}

  /** Initialize API wrappers */
  if (params.mapexOS) apis.mapexOS = createMapexOSApi(params.mapexOS, sessionParams)
  if (params.assets) apis.assets = createAssetsApi(params.assets, sessionParams)
  if (params.events) apis.events = createEventsApi(params.events, sessionParams)
  if (params.router) apis.router = createRouterApi(params.router, sessionParams)
  if (params.triggers) apis.triggers = createTriggersApi(params.triggers, sessionParams)
  if (params.httpGateway) apis.httpGateway = createHttpGatewayApi(params.httpGateway, sessionParams)
  if (params.jsExecutor) apis.jsExecutor = createJsExecutorApi(params.jsExecutor, sessionParams)
  if (params.workflows) apis.workflows = createWorkflowsApi(params.workflows, sessionParams)
  if (params.vault) apis.vault = createVaultApi(params.vault, sessionParams)

  return apis
}

/**
 * Interface for API initialization configuration.
 */
export type ApiWrappers = ReturnType<typeof initializeApis>

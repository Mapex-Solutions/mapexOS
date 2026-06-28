import { isEmpty } from 'lodash';
import type { ApiConfig, ApiInterceptors, SessionsConfig } from '@src/common';
import { createHttp } from '@src/tools';
import { createMarketplaceApi } from './marketplace';

/**
 * Creates the Workflow Plugins Marketplace Service API.
 *
 * Public, no-auth catalog of installable workflow plugins served by the
 * marketplace host. Exposes a single `marketplace` module:
 * - list: paginated, server-side-filtered plugin catalog
 * - facets: category/capability filter options
 * - get: the wrapped raw plugin manifest for one vendor/slug
 * - assetUrl: absolute asset URL builder (e.g. brand icon)
 */
export function createWorkflowPluginsMarketplaceApi(config: ApiConfig, sessionsConfig?: SessionsConfig) {
  const interceptors = !isEmpty(sessionsConfig?.interceptors)
    ? sessionsConfig?.interceptors
    : config?.interceptors || {} as ApiInterceptors;

  const http = createHttp(config, interceptors);

  const marketplace = createMarketplaceApi(http);

  return {
    http,
    marketplace,
  };
}

export type WorkflowPluginsMarketplaceApiMethods = ReturnType<typeof createWorkflowPluginsMarketplaceApi>;

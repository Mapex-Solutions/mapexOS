import { isEmpty } from 'lodash';
import type { ApiConfig, ApiInterceptors, SessionsConfig } from '@src/common';
import { createHttp } from '@src/tools';
import { createMarketplaceApi } from './marketplace';

/**
 * Creates the Asset Templates Marketplace Service API.
 *
 * Public, no-auth catalog of installable asset templates served by the
 * marketplace host. Exposes a single `marketplace` module:
 * - list: paginated, server-side-filtered asset template catalog
 * - facets: category/vendor/model/version filter options
 * - get: the wrapped asset template bundle for one vendor/slug
 * - assetUrl: absolute asset URL builder (e.g. template image)
 */
export function createAssetTemplatesMarketplaceApi(config: ApiConfig, sessionsConfig?: SessionsConfig) {
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

export type AssetTemplatesMarketplaceApiMethods = ReturnType<typeof createAssetTemplatesMarketplaceApi>;

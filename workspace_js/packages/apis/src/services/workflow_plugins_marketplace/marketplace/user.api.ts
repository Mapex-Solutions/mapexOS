import type { AxiosInstance, AxiosResponse } from 'axios';
import type {
  PluginCatalogQuery,
  PluginCatalogListResponse,
  PluginFacetsResponse,
  PluginManifestResponse,
} from '@mapexos/schemas';
import { createApiFactory } from '@src/common';

/**
 * Creates the Workflow Plugins Marketplace API.
 *
 * This is a PUBLIC, no-auth service — the same public host that serves
 * `/api/v1/devices`. No JWT is attached (`useAuthJWT` is omitted), so the
 * shared session `getToken` is never invoked for these calls.
 *
 * The factory's default response handler unwraps Axios `response.data.data`.
 * These endpoints already wrap their payload in a `{status,errors,data}`
 * envelope, so each method overrides `afterRequest` to return the whole
 * envelope (`response.data`) — callers read `res.data.items` for the list
 * and unwrap `res.data` for the manifest, matching the documented contract.
 *
 * The configured baseURL already includes the `/api/v1/workflow_plugins`
 * prefix, so `basePath` is empty and method paths are appended directly.
 *
 * @param http - Axios instance for HTTP requests (baseURL is the full marketplace base)
 * @returns Object with marketplace catalog methods plus an asset URL builder
 */
export function userApi(http: AxiosInstance) {
  const factory = createApiFactory(http);

  const client = factory({
    basePath: '',
    methods: {
      // LIST - GET /
      list: {
        method: 'GET',
        path: '',
        queryParams: {} as PluginCatalogQuery,
        responseType: {} as PluginCatalogListResponse,
        afterRequest: (response: AxiosResponse) => response.data as PluginCatalogListResponse,
      },

      // FACETS - GET /facets
      facets: {
        method: 'GET',
        path: '/facets',
        responseType: {} as PluginFacetsResponse,
        afterRequest: (response: AxiosResponse) => response.data as PluginFacetsResponse,
      },

      // GET MANIFEST - GET /:vendor/:slug (wrapped — caller unwraps .data)
      get: {
        method: 'GET',
        path: '/:vendor/:slug',
        pathParams: {} as { vendor: string; slug: string },
        responseType: {} as PluginManifestResponse,
        afterRequest: (response: AxiosResponse) => response.data as PluginManifestResponse,
      },
    },
  });

  return {
    ...client,

    /**
     * Builds an absolute URL for a plugin asset (e.g. `icon.svg`), suitable
     * for direct use in an `<img src>`. No request is made.
     *
     * @param vendor - Plugin vendor key
     * @param slug - Plugin slug
     * @param assetPath - Asset path under the plugin (e.g. `icon.svg`)
     * @returns Absolute URL to the asset
     */
    assetUrl(vendor: string, slug: string, assetPath: string): string {
      const base = (http.defaults.baseURL ?? '').replace(/\/$/, '');
      return `${base}/${vendor}/${slug}/assets/${assetPath}`;
    },
  };
}

export type WorkflowPluginsMarketplaceUserApiMethods = ReturnType<typeof userApi>;

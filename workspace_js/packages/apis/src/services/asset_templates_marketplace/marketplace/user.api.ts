import type { AxiosInstance, AxiosResponse } from 'axios';
import type {
  AssetTemplateCatalogQuery,
  AssetTemplateCatalogListResponse,
  AssetTemplateFacetsResponse,
  AssetTemplateFacetsQuery,
  AssetTemplateBundleResponse,
} from '@mapexos/schemas';
import { createApiFactory } from '@src/common';

/**
 * Creates the Asset Templates Marketplace API.
 *
 * This is a PUBLIC, no-auth service — the same public host that serves
 * `/api/v1/devices`. No JWT is attached (`useAuthJWT` is omitted), so the
 * shared session `getToken` is never invoked for these calls.
 *
 * The factory's default response handler unwraps Axios `response.data.data`.
 * These endpoints already wrap their payload in a `{status,errors,data}`
 * envelope, so each method overrides `afterRequest` to return the whole
 * envelope (`response.data`) — callers read `res.data.items` for the list
 * and unwrap `res.data` for the bundle, matching the documented contract.
 *
 * The configured baseURL already includes the `/api/v1/asset_templates`
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
        queryParams: {} as AssetTemplateCatalogQuery,
        responseType: {} as AssetTemplateCatalogListResponse,
        afterRequest: (response: AxiosResponse) => response.data as AssetTemplateCatalogListResponse,
      },

      // FACETS - GET /facets?vendor=&model=&lang=
      // The drill-down levels are computed server-side: passing a vendor
      // narrows the returned models, and passing a model narrows the versions.
      facets: {
        method: 'GET',
        path: '/facets',
        queryParams: {} as AssetTemplateFacetsQuery,
        responseType: {} as AssetTemplateFacetsResponse,
        afterRequest: (response: AxiosResponse) => response.data as AssetTemplateFacetsResponse,
      },

      // GET BUNDLE - GET /:vendor/:slug (wrapped — caller unwraps .data)
      get: {
        method: 'GET',
        path: '/:vendor/:slug',
        pathParams: {} as { vendor: string; slug: string },
        responseType: {} as AssetTemplateBundleResponse,
        afterRequest: (response: AxiosResponse) => response.data as AssetTemplateBundleResponse,
      },
    },
  });

  return {
    ...client,

    /**
     * Builds an absolute URL for a template asset (e.g. the template image),
     * suitable for direct use in an `<img src>`. No request is made.
     *
     * @param vendor - Template vendor key
     * @param slug - Template slug
     * @param assetPath - Asset path under the template (e.g. `image.svg`)
     * @returns Absolute URL to the asset
     */
    assetUrl(vendor: string, slug: string, assetPath: string): string {
      const base = (http.defaults.baseURL ?? '').replace(/\/$/, '');
      return `${base}/${vendor}/${slug}/assets/${assetPath}`;
    },
  };
}

export type AssetTemplatesMarketplaceUserApiMethods = ReturnType<typeof userApi>;

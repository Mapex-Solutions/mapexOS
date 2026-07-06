import type { AxiosInstance } from 'axios';

import { userApi } from './user.api';

/**
 * Creates the Marketplace module API.
 *
 * The marketplace endpoints are public (no auth), so no token is threaded
 * through. User methods are spread to the root level.
 *
 * @param http - Axios instance for HTTP requests
 * @returns Object with marketplace catalog methods
 */
export function createMarketplaceApi(http: AxiosInstance) {
  return {
    ...userApi(http),
  };
}

export type MarketplaceApiMethods = ReturnType<typeof createMarketplaceApi>;

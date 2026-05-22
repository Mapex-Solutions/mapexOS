import type { AxiosInstance } from 'axios';
import type { GetToken } from '@src/common';

import { userApi } from './user.api';
import { internalApi } from './internal.api';

/**
 * Creates Asset module API combining user and internal endpoints.
 *
 * User methods (JWT auth) are spread to the root level.
 * Internal methods (API Key auth) are nested under the `internal` property.
 *
 * @param http - Axios instance for HTTP requests
 * @param getToken - Function to retrieve JWT token
 * @returns Object with user API methods (spread) and internal API methods (nested)
 *
 * @example
 * const assetApi = createAssetApi(http, getToken);
 *
 * // User methods (JWT) - spread to root
 * await assetApi.list({ page: 1, perPage: 20 });
 * await assetApi.create({ name: 'Sensor 001', ... });
 * await assetApi.getById({ assetId: '...' });
 *
 * // Internal methods (API Key) - nested under .internal
 * await assetApi.internal.getScriptsByAssetUUID({ assetUUID: 'dev-eui-12345' });
 */
export function createAssetApi(http: AxiosInstance, getToken: GetToken | undefined) {
	return {
		...userApi(http, getToken),  // Spread user methods to root
		internal: internalApi(http), // Internal methods under .internal
	};
}

export type AssetApiMethods = ReturnType<typeof createAssetApi>;

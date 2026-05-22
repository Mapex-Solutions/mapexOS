import type { AxiosInstance } from 'axios';
import type { GetToken } from '@src/common';
import { userApi } from './user.api';

/**
 * Creates complete Asset Template module API.
 *
 * Unlike the asset module, assetTemplate only has user-facing CRUD operations
 * with no internal endpoints, so methods are exported directly without nesting.
 *
 * @param http - Axios instance for HTTP requests
 * @param getToken - Function to retrieve JWT token
 * @returns Object containing asset template API methods
 *
 * @example
 * const assetTemplateApi = createAssetTemplateApi(http, getToken);
 *
 * // Usage - all methods at root level:
 * await assetTemplateApi.list({ page: 1, perPage: 20 });
 * await assetTemplateApi.create({ name: 'LoRaWAN Template', ... });
 * await assetTemplateApi.getById({ assetTemplateId: '507f...' });
 */
export function createAssetTemplateApi(http: AxiosInstance, getToken: GetToken | undefined) {
	return userApi(http, getToken);
}

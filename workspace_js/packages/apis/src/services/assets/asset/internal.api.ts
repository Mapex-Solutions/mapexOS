import type { AssetUUIDParam, AssetScriptsResponse } from '@mapexos/schemas';
import { AxiosInstance } from 'axios';
import { createApiFactory } from '@src/common';
import { ZodAssetUUIDParamSchema } from '@mapexos/schemas';

/**
 * Creates Asset internal API with API Key endpoints.
 * Used for inter-service communication (microservice-to-microservice).
 *
 * @param http - Axios instance for HTTP requests
 * @returns Object containing asset internal API methods
 */
export function internalApi(http: AxiosInstance) {
	const factory = createApiFactory(http);

	return factory({
		basePath: '/api/internal/v1/assets',
		useAuthJWT: false, // API Key authentication
		methods: {
			// GET SCRIPTS BY ASSET UUID - GET /scripts/:assetUUID
			getScriptsByAssetUUID: {
				method: 'GET',
				path: '/scripts/:assetUUID',
				pathParams: {} as AssetUUIDParam,
				paramSchema: ZodAssetUUIDParamSchema,
				responseType: {} as AssetScriptsResponse,
			},
		},
	});
}

export type InternalApiMethods = ReturnType<typeof internalApi>;

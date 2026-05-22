import type {
	AssetId,
	AssetCreate,
	AssetUpdate,
	AssetQuery,
	AssetResponse,
	OrgContext,
	PaginatedResponse,
} from '@mapexos/schemas';

import type { GetToken } from '@src/common';
import { AxiosInstance } from 'axios';
import { createApiFactory } from '@src/common';
import {
	ZodAssetIdSchema,
	ZodAssetCreateSchema,
	ZodAssetUpdateSchema,
	ZodAssetQuerySchema,
	ZodOrgContextSchema,
} from '@mapexos/schemas';

/**
 * Creates Asset user API with external CRUD endpoints.
 *
 * @param http - Axios instance for HTTP requests
 * @param getToken - Function to retrieve JWT token
 * @returns Object containing asset external API methods
 */
export function userApi(http: AxiosInstance, getToken: GetToken | undefined) {
	const factory = createApiFactory(http);

	return factory({
		basePath: '/api/v1/assets',
		useAuthJWT: true,
		getToken,
		methods: {
			// LIST - GET /
			list: {
				method: 'GET',
				path: '',
				queryParams: {} as AssetQuery,
				querySchema: ZodAssetQuerySchema,
				responseType: {} as PaginatedResponse<AssetResponse>,
			},

			// COUNTER - GET /counter
			counter: {
				method: 'GET',
				path: '/counter',
				responseType: {} as { count: number },
			},

			// CREATE - POST /
			// Body includes `protocol.mqtt.password` (plaintext, min 8) for
			// MQTT-protocol assets. The server bcrypt-hashes it before
			// persisting; the plaintext is never readable afterwards.
			create: {
				method: 'POST',
				path: '',
				bodyParams: {} as AssetCreate,
				bodySchema: ZodAssetCreateSchema,
				responseType: {} as AssetResponse,
			},

			// GET BY ID - GET /:assetId
			getById: {
				method: 'GET',
				path: '/:assetId',
				pathParams: {} as AssetId,
				paramSchema: ZodAssetIdSchema,
				responseType: {} as AssetResponse,
			},

			// UPDATE - PATCH /:assetId
			// Operators may include `protocol.mqtt.password` to rotate the
			// MQTT password — the server bcrypts on the way in, same as
			// the create path. Omit it to leave the existing hash alone.
			update: {
				method: 'PATCH',
				path: '/:assetId',
				pathParams: {} as AssetId,
				paramSchema: ZodAssetIdSchema,
				bodyParams: {} as AssetUpdate,
				bodySchema: ZodAssetUpdateSchema,
				responseType: {} as AssetResponse,
			},

			// DELETE - DELETE /:assetId
			delete: {
				method: 'DELETE',
				path: '/:assetId',
				pathParams: {} as AssetId,
				paramSchema: ZodAssetIdSchema,
				responseType: {} as { success: boolean },
			},
		},
	});
}

export type UserApiMethods = ReturnType<typeof userApi>;

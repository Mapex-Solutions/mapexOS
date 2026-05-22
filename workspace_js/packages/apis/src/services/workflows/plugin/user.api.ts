import type {
	PluginId,
	PluginCreate,
	PluginUpdate,
	PluginQuery,
	PluginResponse,
	PaginatedResponse,
} from '@mapexos/schemas';
import type { GetToken } from '@src/common';
import { AxiosInstance } from 'axios';
import { createApiFactory } from '@src/common';
import {
	ZodPluginIdSchema,
	ZodPluginCreateSchema,
	ZodPluginUpdateSchema,
	ZodPluginQuerySchema,
} from '@mapexos/schemas';

/**
 * Creates Workflow Plugin user API with external CRUD endpoints.
 *
 * @param http - Axios instance for HTTP requests
 * @param getToken - Function to retrieve JWT token
 * @returns Object containing plugin external API methods
 */
export function userApi(http: AxiosInstance, getToken: GetToken | undefined) {
	const factory = createApiFactory(http);

	return factory({
		basePath: '/api/v1/plugins',
		useAuthJWT: true,
		getToken,
		methods: {
			// LIST - GET /
			list: {
				method: 'GET',
				path: '',
				queryParams: {} as PluginQuery,
				querySchema: ZodPluginQuerySchema,
				responseType: {} as PaginatedResponse<PluginResponse>,
			},

			// CREATE - POST /
			create: {
				method: 'POST',
				path: '',
				bodyParams: {} as PluginCreate,
				bodySchema: ZodPluginCreateSchema,
				responseType: {} as PluginResponse,
			},

			// GET ENABLED - GET /enabled
			getEnabled: {
				method: 'GET',
				path: '/enabled',
				responseType: {} as PluginResponse[],
			},

			// GET BY ID - GET /:id
			getById: {
				method: 'GET',
				path: '/:id',
				pathParams: {} as PluginId,
				paramSchema: ZodPluginIdSchema,
				responseType: {} as PluginResponse,
			},

			// UPDATE - PATCH /:id
			update: {
				method: 'PATCH',
				path: '/:id',
				pathParams: {} as PluginId,
				bodyParams: {} as PluginUpdate,
				paramSchema: ZodPluginIdSchema,
				bodySchema: ZodPluginUpdateSchema,
				responseType: {} as PluginResponse,
			},

			// DELETE - DELETE /:id
			delete: {
				method: 'DELETE',
				path: '/:id',
				pathParams: {} as PluginId,
				paramSchema: ZodPluginIdSchema,
				responseType: {} as { success: boolean },
			},
		},
	});
}

export type UserApiMethods = ReturnType<typeof userApi>;

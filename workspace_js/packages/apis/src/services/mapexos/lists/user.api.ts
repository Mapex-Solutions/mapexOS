import type {
	ListId,
	ListCreate,
	ListUpdate,
	ListQuery,
	ListResponse,
	PaginatedResponse,
} from '@mapexos/schemas';
import type { GetToken } from '@src/common';
import { AxiosInstance } from 'axios';
import { createApiFactory } from '@src/common';
import {
	ZodListIdSchema,
	ZodListCreateSchema,
	ZodListUpdateSchema,
	ZodListQuerySchema,
} from '@mapexos/schemas';

/**
 * Creates List API with external CRUD endpoints.
 *
 * @param http - Axios instance for HTTP requests
 * @param getToken - Function to retrieve JWT token
 * @returns Object containing list API methods
 */
export function userApi(http: AxiosInstance, getToken: GetToken | undefined) {
	const factory = createApiFactory(http);

	return factory({
		basePath: '/api/v1/lists',
		useAuthJWT: true,
		getToken,
		methods: {
			// LIST - GET /
			list: {
				method: 'GET',
				path: '',
				queryParams: {} as ListQuery,
				querySchema: ZodListQuerySchema,
				responseType: {} as PaginatedResponse<ListResponse>,
			},

			// CREATE - POST /
			create: {
				method: 'POST',
				path: '',
				bodyParams: {} as ListCreate,
				bodySchema: ZodListCreateSchema,
				responseType: {} as ListResponse,
			},

			// GET BY ID - GET /:listId
			getById: {
				method: 'GET',
				path: '/:listId',
				pathParams: {} as ListId,
				paramSchema: ZodListIdSchema,
				responseType: {} as ListResponse,
			},

			// UPDATE - PATCH /:listId
			update: {
				method: 'PATCH',
				path: '/:listId',
				pathParams: {} as ListId,
				bodyParams: {} as ListUpdate,
				paramSchema: ZodListIdSchema,
				bodySchema: ZodListUpdateSchema,
				responseType: {} as ListResponse,
			},

			// DELETE - DELETE /:listId
			delete: {
				method: 'DELETE',
				path: '/:listId',
				pathParams: {} as ListId,
				paramSchema: ZodListIdSchema,
				responseType: {} as { success: boolean },
			},
		},
	});
}

export type UserApiMethods = ReturnType<typeof userApi>;

import type {
	RoleId,
	RoleCreate,
	RoleUpdate,
	RoleQuery,
	RoleResponse,
	PaginatedResponse,
} from '@mapexos/schemas';
import type { GetToken } from '@src/common';
import { AxiosInstance } from 'axios';
import { createApiFactory } from '@src/common';
import {
	ZodRoleIdSchema,
	ZodRoleCreateSchema,
	ZodRoleUpdateSchema,
	ZodRoleQuerySchema,
} from '@mapexos/schemas';

/**
 * Creates Role API with external CRUD endpoints.
 *
 * @param http - Axios instance for HTTP requests
 * @param getToken - Function to retrieve JWT token
 * @returns Object containing role API methods
 */
export function userApi(http: AxiosInstance, getToken: GetToken | undefined) {
	const factory = createApiFactory(http);

	return factory({
		basePath: '/api/v1/roles',
		useAuthJWT: true,
		getToken,
		methods: {
			// LIST - GET /
			list: {
				method: 'GET',
				path: '',
				queryParams: {} as RoleQuery,
				querySchema: ZodRoleQuerySchema,
				responseType: {} as PaginatedResponse<RoleResponse>,
			},

			// CREATE - POST /
			create: {
				method: 'POST',
				path: '',
				bodyParams: {} as RoleCreate,
				bodySchema: ZodRoleCreateSchema,
				responseType: {} as RoleResponse,
			},

			// GET BY ID - GET /:roleId
			getById: {
				method: 'GET',
				path: '/:roleId',
				pathParams: {} as RoleId,
				paramSchema: ZodRoleIdSchema,
				responseType: {} as RoleResponse,
			},

			// UPDATE - PATCH /:roleId
			update: {
				method: 'PATCH',
				path: '/:roleId',
				pathParams: {} as RoleId,
				bodyParams: {} as RoleUpdate,
				paramSchema: ZodRoleIdSchema,
				bodySchema: ZodRoleUpdateSchema,
				responseType: {} as RoleResponse,
			},

			// DELETE - DELETE /:roleId
			delete: {
				method: 'DELETE',
				path: '/:roleId',
				pathParams: {} as RoleId,
				paramSchema: ZodRoleIdSchema,
				responseType: {} as { success: boolean },
			},
		},
	});
}

export type UserApiMethods = ReturnType<typeof userApi>;

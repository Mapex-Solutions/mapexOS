import type {
	RouteGroupId,
	RouteGroupCreate,
	RouteGroupUpdate,
	RouteGroupQuery,
	RouteGroupResponse,
	PaginatedResponse,
} from '@mapexos/schemas';
import type { GetToken } from '@src/common';
import { AxiosInstance } from 'axios';
import { createApiFactory } from '@src/common';
import {
	ZodRouteGroupIdSchema,
	ZodRouteGroupCreateSchema,
	ZodRouteGroupUpdateSchema,
	ZodRouteGroupQuerySchema,
} from '@mapexos/schemas';

/**
 * Creates RouteGroup user API with external CRUD endpoints.
 *
 * @param http - Axios instance for HTTP requests
 * @param getToken - Function to retrieve JWT token
 * @returns Object containing route group external API methods
 */
export function userApi(http: AxiosInstance, getToken: GetToken | undefined) {
	const factory = createApiFactory(http);

	return factory({
		basePath: '/api/v1/route_groups',
		useAuthJWT: true,
		getToken,
		methods: {
			// LIST - GET /
			list: {
				method: 'GET',
				path: '',
				queryParams: {} as RouteGroupQuery,
				querySchema: ZodRouteGroupQuerySchema,
				responseType: {} as PaginatedResponse<RouteGroupResponse>,
			},

			// COUNTER - GET /counter
			counter: {
				method: 'GET',
				path: '/counter',
				responseType: {} as { count: number },
			},

			// CREATE - POST /
			create: {
				method: 'POST',
				path: '',
				bodyParams: {} as RouteGroupCreate,
				bodySchema: ZodRouteGroupCreateSchema,
				responseType: {} as RouteGroupResponse,
			},

			// GET BY ID - GET /:routeGroupId
			getById: {
				method: 'GET',
				path: '/:routeGroupId',
				pathParams: {} as RouteGroupId,
				paramSchema: ZodRouteGroupIdSchema,
				responseType: {} as RouteGroupResponse,
			},

			// UPDATE - PATCH /:routeGroupId
			update: {
				method: 'PATCH',
				path: '/:routeGroupId',
				pathParams: {} as RouteGroupId,
				bodyParams: {} as RouteGroupUpdate,
				paramSchema: ZodRouteGroupIdSchema,
				bodySchema: ZodRouteGroupUpdateSchema,
				responseType: {} as RouteGroupResponse,
			},

			// DELETE - DELETE /:routeGroupId
			delete: {
				method: 'DELETE',
				path: '/:routeGroupId',
				pathParams: {} as RouteGroupId,
				paramSchema: ZodRouteGroupIdSchema,
				responseType: {} as { success: boolean },
			},
		},
	});
}

export type UserApiMethods = ReturnType<typeof userApi>;

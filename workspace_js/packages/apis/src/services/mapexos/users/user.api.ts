import type {
	UserId,
	UserCreate,
	UserUpdate,
	UserQuery,
	UserResponse,
	PaginatedResponse,
} from '@mapexos/schemas';
import type { GetToken } from '@src/common';
import { AxiosInstance } from 'axios';
import { createApiFactory } from '@src/common';
import {
	ZodUserIdSchema,
	ZodUserCreateSchema,
	ZodUserUpdateSchema,
	ZodUserQuerySchema,
} from '@mapexos/schemas';

/**
 * Creates User API with external CRUD endpoints.
 *
 * @param http - Axios instance for HTTP requests
 * @param getToken - Function to retrieve JWT token
 * @returns Object containing user API methods
 */
export function userApi(http: AxiosInstance, getToken: GetToken | undefined) {
	const factory = createApiFactory(http);

	return factory({
		basePath: '/api/v1/users',
		useAuthJWT: true,
		getToken,
		methods: {
			// GET MY INFO - GET /me
			me: {
				method: 'GET',
				path: '/me',
				responseType: {} as UserResponse,
			},

			// UPDATE MY INFO - PATCH /me
			updateMe: {
				method: 'PATCH',
				path: '/me',
				bodyParams: {} as UserUpdate,
				bodySchema: ZodUserUpdateSchema,
				responseType: {} as UserResponse,
			},

			// DISABLE MY TOUR - PATCH /me/tour
			disableMyTour: {
				method: 'PATCH',
				path: '/me/tour',
				responseType: {} as UserResponse,
			},

			// LIST - GET /
			list: {
				method: 'GET',
				path: '',
				queryParams: {} as UserQuery,
				querySchema: ZodUserQuerySchema,
				responseType: {} as PaginatedResponse<UserResponse>,
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
				bodyParams: {} as UserCreate,
				bodySchema: ZodUserCreateSchema,
				responseType: {} as UserResponse,
			},

			// GET BY ID - GET /:userId
			getById: {
				method: 'GET',
				path: '/:userId',
				pathParams: {} as UserId,
				paramSchema: ZodUserIdSchema,
				responseType: {} as UserResponse,
			},

			// UPDATE - PATCH /:userId
			update: {
				method: 'PATCH',
				path: '/:userId',
				pathParams: {} as UserId,
				bodyParams: {} as UserUpdate,
				paramSchema: ZodUserIdSchema,
				bodySchema: ZodUserUpdateSchema,
				responseType: {} as UserResponse,
			},

			// DELETE - DELETE /:userId
			delete: {
				method: 'DELETE',
				path: '/:userId',
				pathParams: {} as UserId,
				paramSchema: ZodUserIdSchema,
				responseType: {} as { success: boolean },
			},
		},
	});
}

export type UserApiMethods = ReturnType<typeof userApi>;

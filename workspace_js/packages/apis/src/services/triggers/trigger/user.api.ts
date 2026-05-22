import type {
	TriggerId,
	TriggerCreate,
	TriggerUpdate,
	TriggerQuery,
	TriggerResponse,
	PaginatedResponse,
} from '@mapexos/schemas';
import type { GetToken } from '@src/common';
import { AxiosInstance } from 'axios';
import { createApiFactory } from '@src/common';
import {
	ZodTriggerIdSchema,
	ZodTriggerCreateSchema,
	ZodTriggerUpdateSchema,
	ZodTriggerQuerySchema,
} from '@mapexos/schemas';

/**
 * Creates Trigger user API with external CRUD endpoints.
 *
 * @param http - Axios instance for HTTP requests
 * @param getToken - Function to retrieve JWT token
 * @returns Object containing trigger external API methods
 */
export function userApi(http: AxiosInstance, getToken: GetToken | undefined) {
	const factory = createApiFactory(http);

	return factory({
		basePath: '/api/v1/triggers',
		useAuthJWT: true,
		getToken,
		methods: {
			// LIST - GET /
			list: {
				method: 'GET',
				path: '',
				queryParams: {} as TriggerQuery,
				querySchema: ZodTriggerQuerySchema,
				responseType: {} as PaginatedResponse<TriggerResponse>,
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
				bodyParams: {} as TriggerCreate,
				bodySchema: ZodTriggerCreateSchema,
				responseType: {} as TriggerResponse,
			},

			// GET BY ID - GET /:id
			getById: {
				method: 'GET',
				path: '/:triggerId',
				pathParams: {} as TriggerId,
				paramSchema: ZodTriggerIdSchema,
				responseType: {} as TriggerResponse,
			},

			// UPDATE - PATCH /:id
			update: {
				method: 'PATCH',
				path: '/:triggerId',
				pathParams: {} as TriggerId,
				bodyParams: {} as TriggerUpdate,
				paramSchema: ZodTriggerIdSchema,
				bodySchema: ZodTriggerUpdateSchema,
				responseType: {} as TriggerResponse,
			},

			// DELETE - DELETE /:id
			delete: {
				method: 'DELETE',
				path: '/:triggerId',
				pathParams: {} as TriggerId,
				paramSchema: ZodTriggerIdSchema,
				responseType: {} as { success: boolean },
			},
		},
	});
}

export type UserApiMethods = ReturnType<typeof userApi>;

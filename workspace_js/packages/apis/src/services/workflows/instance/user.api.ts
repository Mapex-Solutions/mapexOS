import type {
	InstanceId,
	InstanceCreate,
	InstanceUpdate,
	InstanceQuery,
	InstanceResponse,
	ExecuteRequest,
	ExecuteResponse,
	PaginatedResponse,
} from '@mapexos/schemas';
import type { GetToken } from '@src/common';
import { AxiosInstance } from 'axios';
import { createApiFactory } from '@src/common';
import {
	ZodInstanceIdSchema,
	ZodInstanceCreateSchema,
	ZodInstanceUpdateSchema,
	ZodInstanceQuerySchema,
	ZodExecuteRequestSchema,
} from '@mapexos/schemas';

/**
 * Creates Workflow Instance config user API with CRUD endpoints.
 *
 * @param http - Axios instance for HTTP requests
 * @param getToken - Function to retrieve JWT token
 * @returns Object containing instance config CRUD API methods
 */
export function userApi(http: AxiosInstance, getToken: GetToken | undefined) {
	const factory = createApiFactory(http);

	return factory({
		basePath: '/api/v1/workflow_instances',
		useAuthJWT: true,
		getToken,
		methods: {
			// COUNTER - GET /counter
			counter: {
				method: 'GET',
				path: '/counter',
				responseType: {} as { count: number },
			},

			// LIST - GET /
			list: {
				method: 'GET',
				path: '',
				queryParams: {} as InstanceQuery,
				querySchema: ZodInstanceQuerySchema,
				responseType: {} as PaginatedResponse<InstanceResponse>,
			},

			// CREATE - POST /
			create: {
				method: 'POST',
				path: '',
				bodyParams: {} as InstanceCreate,
				bodySchema: ZodInstanceCreateSchema,
				responseType: {} as InstanceResponse,
			},

			// GET BY ID - GET /:instanceId
			getById: {
				method: 'GET',
				path: '/:instanceId',
				pathParams: {} as InstanceId,
				paramSchema: ZodInstanceIdSchema,
				responseType: {} as InstanceResponse,
			},

			// UPDATE - PUT /:instanceId
			update: {
				method: 'PUT',
				path: '/:instanceId',
				pathParams: {} as InstanceId,
				bodyParams: {} as InstanceUpdate,
				paramSchema: ZodInstanceIdSchema,
				bodySchema: ZodInstanceUpdateSchema,
				responseType: {} as InstanceResponse,
			},

			// EXECUTE - POST /:instanceId/execute
			execute: {
				method: 'POST',
				path: '/:instanceId/execute',
				pathParams: {} as InstanceId,
				bodyParams: {} as ExecuteRequest,
				paramSchema: ZodInstanceIdSchema,
				bodySchema: ZodExecuteRequestSchema,
				responseType: {} as ExecuteResponse,
			},

			// DELETE - DELETE /:instanceId
			delete: {
				method: 'DELETE',
				path: '/:instanceId',
				pathParams: {} as InstanceId,
				paramSchema: ZodInstanceIdSchema,
				responseType: {} as { success: boolean },
			},
		},
	});
}

export type UserApiMethods = ReturnType<typeof userApi>;

import type {
	DefinitionId,
	DefinitionCreate,
	DefinitionUpdate,
	DefinitionQuery,
	DefinitionResponse,
	PaginatedResponse,
} from '@mapexos/schemas';
import type { GetToken } from '@src/common';
import { AxiosInstance } from 'axios';
import { createApiFactory } from '@src/common';
import {
	ZodDefinitionIdSchema,
	ZodDefinitionCreateSchema,
	ZodDefinitionUpdateSchema,
	ZodDefinitionQuerySchema,
} from '@mapexos/schemas';

/**
 * Creates Workflow Definition user API with external CRUD endpoints.
 *
 * @param http - Axios instance for HTTP requests
 * @param getToken - Function to retrieve JWT token
 * @returns Object containing definition external API methods
 */
export function userApi(http: AxiosInstance, getToken: GetToken | undefined) {
	const factory = createApiFactory(http);

	return factory({
		basePath: '/api/v1/workflow_definitions',
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
				queryParams: {} as DefinitionQuery,
				querySchema: ZodDefinitionQuerySchema,
				responseType: {} as PaginatedResponse<DefinitionResponse>,
			},

			// CREATE - POST /
			create: {
				method: 'POST',
				path: '',
				bodyParams: {} as DefinitionCreate,
				bodySchema: ZodDefinitionCreateSchema,
				responseType: {} as DefinitionResponse,
			},

			// GET BY ID - GET /:workflowId
			getById: {
				method: 'GET',
				path: '/:workflowId',
				pathParams: {} as DefinitionId,
				paramSchema: ZodDefinitionIdSchema,
				responseType: {} as DefinitionResponse,
			},

			// UPDATE - PATCH /:workflowId
			update: {
				method: 'PATCH',
				path: '/:workflowId',
				pathParams: {} as DefinitionId,
				bodyParams: {} as DefinitionUpdate,
				paramSchema: ZodDefinitionIdSchema,
				bodySchema: ZodDefinitionUpdateSchema,
				responseType: {} as DefinitionResponse,
			},

			// DELETE - DELETE /:workflowId
			delete: {
				method: 'DELETE',
				path: '/:workflowId',
				pathParams: {} as DefinitionId,
				paramSchema: ZodDefinitionIdSchema,
				responseType: {} as { success: boolean },
			},
		},
	});
}

export type UserApiMethods = ReturnType<typeof userApi>;

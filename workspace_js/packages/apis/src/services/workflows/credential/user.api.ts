import type {
	CredentialId,
	CredentialPluginId,
	CredentialCreate,
	CredentialUpdate,
	CredentialQuery,
	CredentialResponse,
	CredentialTestResult,
	CredentialSchemaResponse,
	LoadOptionsParams,
	LoadOptionsBody,
	LoadOptionsItem,
	PaginatedResponse,
} from '@mapexos/schemas';
import type { GetToken } from '@src/common';
import { AxiosInstance } from 'axios';
import { createApiFactory } from '@src/common';
import {
	ZodCredentialIdSchema,
	ZodCredentialPluginIdSchema,
	ZodCredentialCreateSchema,
	ZodCredentialUpdateSchema,
	ZodCredentialQuerySchema,
	ZodLoadOptionsParamsSchema,
	ZodLoadOptionsBodySchema,
} from '@mapexos/schemas';

/**
 * Creates Workflow Credential user API with external CRUD endpoints.
 *
 * @param http - Axios instance for HTTP requests
 * @param getToken - Function to retrieve JWT token
 * @returns Object containing credential external API methods
 */
export function userApi(http: AxiosInstance, getToken: GetToken | undefined) {
	const factory = createApiFactory(http);

	return factory({
		basePath: '/api/v1/credentials',
		useAuthJWT: true,
		getToken,
		methods: {
			// LIST - GET /
			list: {
				method: 'GET',
				path: '',
				queryParams: {} as CredentialQuery,
				querySchema: ZodCredentialQuerySchema,
				responseType: {} as PaginatedResponse<CredentialResponse>,
			},

			// CREATE - POST /
			create: {
				method: 'POST',
				path: '',
				bodyParams: {} as CredentialCreate,
				bodySchema: ZodCredentialCreateSchema,
				responseType: {} as CredentialResponse,
			},

			// GET SCHEMA - GET /schema/:pluginId
			getSchema: {
				method: 'GET',
				path: '/schema/:pluginId',
				pathParams: {} as CredentialPluginId,
				paramSchema: ZodCredentialPluginIdSchema,
				responseType: {} as CredentialSchemaResponse,
			},

			// GET BY ID - GET /:id
			getById: {
				method: 'GET',
				path: '/:id',
				pathParams: {} as CredentialId,
				paramSchema: ZodCredentialIdSchema,
				responseType: {} as CredentialResponse,
			},

			// UPDATE - PATCH /:id
			update: {
				method: 'PATCH',
				path: '/:id',
				pathParams: {} as CredentialId,
				bodyParams: {} as CredentialUpdate,
				paramSchema: ZodCredentialIdSchema,
				bodySchema: ZodCredentialUpdateSchema,
				responseType: {} as CredentialResponse,
			},

			// DELETE - DELETE /:id
			delete: {
				method: 'DELETE',
				path: '/:id',
				pathParams: {} as CredentialId,
				paramSchema: ZodCredentialIdSchema,
				responseType: {} as { success: boolean },
			},

			// TEST - POST /:id/test
			test: {
				method: 'POST',
				path: '/:id/test',
				pathParams: {} as CredentialId,
				paramSchema: ZodCredentialIdSchema,
				responseType: {} as CredentialTestResult,
			},

			// LOAD OPTIONS - POST /:id/load_options/:resourceKey
			loadOptions: {
				method: 'POST',
				path: '/:id/load_options/:resourceKey',
				pathParams: {} as LoadOptionsParams,
				bodyParams: {} as LoadOptionsBody,
				paramSchema: ZodLoadOptionsParamsSchema,
				bodySchema: ZodLoadOptionsBodySchema,
				responseType: {} as LoadOptionsItem[],
			},
		},
	});
}

export type UserApiMethods = ReturnType<typeof userApi>;

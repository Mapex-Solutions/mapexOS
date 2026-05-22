import type {
	OrganizationId,
	OrganizationCreate,
	OrganizationUpdate,
	OrganizationQuery,
	OrganizationResponse,
	TreeQuery,
	TreeResponse,
	PaginatedResponse,
} from '@mapexos/schemas';
import type { GetToken } from '@src/common';
import { AxiosInstance } from 'axios';
import { createApiFactory } from '@src/common';
import {
	ZodOrganizationIdSchema,
	ZodOrganizationCreateSchema,
	ZodOrganizationUpdateSchema,
	ZodOrganizationQuerySchema,
	ZodTreeQuerySchema,
} from '@mapexos/schemas';

/**
 * Creates Organization API with external CRUD endpoints.
 *
 * @param http - Axios instance for HTTP requests
 * @param getToken - Function to retrieve JWT token
 * @returns Object containing organization API methods
 */
export function userApi(http: AxiosInstance, getToken: GetToken | undefined) {
	const factory = createApiFactory(http);

	return factory({
		basePath: '/api/v1/organizations',
		useAuthJWT: true,
		getToken,
		methods: {
			// LIST - GET /
			list: {
				method: 'GET',
				path: '',
				queryParams: {} as OrganizationQuery,
				querySchema: ZodOrganizationQuerySchema,
				responseType: {} as PaginatedResponse<OrganizationResponse>,
			},

			// TREE - GET /tree (cursor pagination for UI)
			tree: {
				method: 'GET',
				path: '/tree',
				queryParams: {} as TreeQuery,
				querySchema: ZodTreeQuerySchema,
				responseType: {} as TreeResponse,
			},

			// CREATE - POST /
			create: {
				method: 'POST',
				path: '',
				bodyParams: {} as OrganizationCreate,
				bodySchema: ZodOrganizationCreateSchema,
				responseType: {} as OrganizationResponse,
			},

			// GET BY ID - GET /:organizationId
			getById: {
				method: 'GET',
				path: '/:organizationId',
				pathParams: {} as OrganizationId,
				paramSchema: ZodOrganizationIdSchema,
				responseType: {} as OrganizationResponse,
			},

			// UPDATE - PATCH /:organizationId
			update: {
				method: 'PATCH',
				path: '/:organizationId',
				pathParams: {} as OrganizationId,
				bodyParams: {} as OrganizationUpdate,
				paramSchema: ZodOrganizationIdSchema,
				bodySchema: ZodOrganizationUpdateSchema,
				responseType: {} as OrganizationResponse,
			},

			// DELETE - DELETE /:organizationId
			delete: {
				method: 'DELETE',
				path: '/:organizationId',
				pathParams: {} as OrganizationId,
				paramSchema: ZodOrganizationIdSchema,
				responseType: {} as { success: boolean },
			},
		},
	});
}

export type UserApiMethods = ReturnType<typeof userApi>;

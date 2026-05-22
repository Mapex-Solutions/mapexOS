import type {
	MembershipId,
	MembershipCreate,
	MembershipUpdate,
	MembershipQuery,
	MembershipResponse,
	PaginatedResponse,
} from '@mapexos/schemas';
import type { GetToken } from '@src/common';
import { AxiosInstance } from 'axios';
import { createApiFactory } from '@src/common';
import {
	ZodMembershipIdSchema,
	ZodMembershipCreateSchema,
	ZodMembershipUpdateSchema,
	ZodMembershipQuerySchema,
} from '@mapexos/schemas';

/**
 * Creates Membership API with CRUD endpoints.
 *
 * @param http - Axios instance for HTTP requests
 * @param getToken - Function to retrieve JWT token
 * @returns Object containing membership API methods
 */
export function membershipApi(http: AxiosInstance, getToken: GetToken | undefined) {
	const factory = createApiFactory(http);

	return factory({
		basePath: '/api/v1/memberships',
		useAuthJWT: true,
		getToken,
		methods: {
			// LIST - GET /
			list: {
				method: 'GET',
				path: '',
				queryParams: {} as MembershipQuery,
				querySchema: ZodMembershipQuerySchema,
				responseType: {} as PaginatedResponse<MembershipResponse>,
			},

			// CREATE - POST /
			create: {
				method: 'POST',
				path: '',
				bodyParams: {} as MembershipCreate,
				bodySchema: ZodMembershipCreateSchema,
				responseType: {} as MembershipResponse,
			},

			// GET BY ID - GET /:membershipId
			getById: {
				method: 'GET',
				path: '/:membershipId',
				pathParams: {} as MembershipId,
				paramSchema: ZodMembershipIdSchema,
				responseType: {} as MembershipResponse,
			},

			// UPDATE - PATCH /:membershipId
			update: {
				method: 'PATCH',
				path: '/:membershipId',
				pathParams: {} as MembershipId,
				bodyParams: {} as MembershipUpdate,
				paramSchema: ZodMembershipIdSchema,
				bodySchema: ZodMembershipUpdateSchema,
				responseType: {} as MembershipResponse,
			},

			// DELETE - DELETE /:membershipId
			delete: {
				method: 'DELETE',
				path: '/:membershipId',
				pathParams: {} as MembershipId,
				paramSchema: ZodMembershipIdSchema,
				responseType: {} as { success: boolean },
			},
		},
	});
}

export type MembershipApiMethods = ReturnType<typeof membershipApi>;

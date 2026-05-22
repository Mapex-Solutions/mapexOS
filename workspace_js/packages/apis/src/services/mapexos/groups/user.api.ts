import type {
	GroupId,
	GroupCreate,
	GroupUpdate,
	GroupQuery,
	GroupResponse,
	GroupMembersQuery,
	GroupMemberResponse,
	GroupMemberAdd,
	GroupMemberId,
	PaginatedResponse,
} from '@mapexos/schemas';
import type { GetToken } from '@src/common';
import { AxiosInstance } from 'axios';
import { createApiFactory } from '@src/common';
import {
	ZodGroupIdSchema,
	ZodGroupCreateSchema,
	ZodGroupUpdateSchema,
	ZodGroupQuerySchema,
	ZodGroupMembersQuerySchema,
	ZodGroupMemberAddSchema,
	ZodGroupMemberIdSchema,
} from '@mapexos/schemas';

/**
 * Creates Group API with external CRUD endpoints.
 *
 * @param http - Axios instance for HTTP requests
 * @param getToken - Function to retrieve JWT token
 * @returns Object containing group API methods
 */
export function userApi(http: AxiosInstance, getToken: GetToken | undefined) {
	const factory = createApiFactory(http);

	return factory({
		basePath: '/api/v1/groups',
		useAuthJWT: true,
		getToken,
		methods: {
			// LIST - GET /
			list: {
				method: 'GET',
				path: '',
				queryParams: {} as GroupQuery,
				querySchema: ZodGroupQuerySchema,
				responseType: {} as PaginatedResponse<GroupResponse>,
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
				bodyParams: {} as GroupCreate,
				bodySchema: ZodGroupCreateSchema,
				responseType: {} as GroupResponse,
			},

			// GET BY ID - GET /:groupId
			getById: {
				method: 'GET',
				path: '/:groupId',
				pathParams: {} as GroupId,
				paramSchema: ZodGroupIdSchema,
				responseType: {} as GroupResponse,
			},

			// UPDATE - PATCH /:groupId
			update: {
				method: 'PATCH',
				path: '/:groupId',
				pathParams: {} as GroupId,
				bodyParams: {} as GroupUpdate,
				paramSchema: ZodGroupIdSchema,
				bodySchema: ZodGroupUpdateSchema,
				responseType: {} as GroupResponse,
			},

			// DELETE - DELETE /:groupId
			delete: {
				method: 'DELETE',
				path: '/:groupId',
				pathParams: {} as GroupId,
				paramSchema: ZodGroupIdSchema,
				responseType: {} as { success: boolean },
			},

			// GET MEMBERS - GET /:groupId/members (paginated, max 100 per page)
			getMembers: {
				method: 'GET',
				path: '/:groupId/members',
				pathParams: {} as GroupId,
				queryParams: {} as GroupMembersQuery,
				paramSchema: ZodGroupIdSchema,
				querySchema: ZodGroupMembersQuerySchema,
				responseType: {} as PaginatedResponse<GroupMemberResponse>,
			},

			// ADD MEMBER - POST /:groupId/members
			addMember: {
				method: 'POST',
				path: '/:groupId/members',
				pathParams: {} as GroupId,
				bodyParams: {} as GroupMemberAdd,
				paramSchema: ZodGroupIdSchema,
				bodySchema: ZodGroupMemberAddSchema,
				responseType: {} as { success: boolean },
			},

			// REMOVE MEMBER - DELETE /:groupId/members/:userId
			removeMember: {
				method: 'DELETE',
				path: '/:groupId/members/:userId',
				pathParams: {} as GroupMemberId,
				paramSchema: ZodGroupMemberIdSchema,
				responseType: {} as { success: boolean },
			},
		},
	});
}

export type UserApiMethods = ReturnType<typeof userApi>;

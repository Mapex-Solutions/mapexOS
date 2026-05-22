import type {
	RetentionPolicyQuery,
	RetentionPolicyPaginatedResult,
	RetentionPolicyParams,
	RetentionPolicyResponse,
	RetentionPolicyUpsert,
} from '@mapexos/schemas';

import type { GetToken } from '@src/common';
import { AxiosInstance } from 'axios';
import { createApiFactory } from '@src/common';
import {
	ZodRetentionPolicyQuerySchema,
	ZodRetentionPolicyParamsSchema,
	ZodRetentionPolicyUpsertSchema,
} from '@mapexos/schemas';

/**
 * Creates Retention user API with external endpoints.
 *
 * @param http - Axios instance for HTTP requests
 * @param getToken - Function to retrieve JWT token
 * @returns Object containing retention external API methods
 */
export function userApi(http: AxiosInstance, getToken: GetToken | undefined) {
	const factory = createApiFactory(http);

	return factory({
		basePath: '/api/v1/retention',
		useAuthJWT: true,
		getToken,
		methods: {
			/**
			 * LIST RETENTION POLICIES - GET /
			 *
			 * Uses offset-based pagination. Returns items with pagination metadata.
			 *
			 * Query params:
			 *   - page: page number (default: 1)
			 *   - perPage: items per page (default: 20)
			 *   - sort: sort field
			 *   - includeChildren: include child org policies
			 *   - type: filter by type (comma-separated for $in)
			 */
			listRetentionPolicies: {
				method: 'GET',
				path: '/',
				queryParams: {} as RetentionPolicyQuery,
				querySchema: ZodRetentionPolicyQuerySchema,
				responseType: {} as RetentionPolicyPaginatedResult,
			},

			/**
			 * GET RETENTION POLICY BY ID - GET /:retentionPolicyId
			 *
			 * Retrieves a single retention policy by its MongoDB ObjectId.
			 */
			getRetentionPolicyById: {
				method: 'GET',
				path: '/:retentionPolicyId',
				pathParams: {} as RetentionPolicyParams,
				paramSchema: ZodRetentionPolicyParamsSchema,
				responseType: {} as RetentionPolicyResponse,
			},

			/**
			 * UPSERT RETENTION POLICY - PUT /
			 *
			 * Creates or updates a retention policy by orgId + type.
			 * OrgId and pathKey are populated from RequestContext middleware.
			 */
			upsertRetentionPolicy: {
				method: 'PUT',
				path: '/',
				bodyParams: {} as RetentionPolicyUpsert,
				bodySchema: ZodRetentionPolicyUpsertSchema,
				responseType: {} as RetentionPolicyResponse,
			},

			/**
			 * DELETE RETENTION POLICY BY ID - DELETE /:retentionPolicyId
			 *
			 * Removes a retention policy by its MongoDB ObjectId.
			 */
			deleteRetentionPolicyById: {
				method: 'DELETE',
				path: '/:retentionPolicyId',
				pathParams: {} as RetentionPolicyParams,
				paramSchema: ZodRetentionPolicyParamsSchema,
				responseType: {} as { success: boolean },
			},
		},
	});
}

export type UserApiMethods = ReturnType<typeof userApi>;

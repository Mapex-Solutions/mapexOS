import type {
	MigrationPlanQuery,
	MigrationPlanParams,
	MigrationPlanCreateRequest,
	MigrationPlanUpdateRequest,
	MigrationPlanResponse,
	MigrationExecutionQuery,
	MigrationExecutionResponse,
	PaginatedResponse,
} from '@mapexos/schemas';
import type { GetToken } from '@src/common';
import { AxiosInstance } from 'axios';
import { createApiFactory } from '@src/common';
import {
	ZodMigrationPlanQuerySchema,
	ZodMigrationPlanParamsSchema,
	ZodMigrationPlanCreateSchema,
	ZodMigrationPlanUpdateSchema,
	ZodMigrationExecutionQuerySchema,
} from '@mapexos/schemas';

/**
 * Creates the Template Migration user API (external JWT endpoints on the
 * assets service).
 *
 * Covers the operator surface of template migration plans: list, detail,
 * create, edit, cancel, and per-asset execution listing.
 *
 * @param http - Axios instance for HTTP requests
 * @param getToken - Function to retrieve JWT token
 * @returns Object containing template migration API methods
 */
export function userApi(http: AxiosInstance, getToken: GetToken | undefined) {
	const factory = createApiFactory(http);

	return factory({
		basePath: '/api/v1/asset_templates/migrations',
		useAuthJWT: true,
		getToken,
		methods: {
			// LIST - GET /
			list: {
				method: 'GET',
				path: '',
				queryParams: {} as MigrationPlanQuery,
				querySchema: ZodMigrationPlanQuerySchema,
				responseType: {} as PaginatedResponse<MigrationPlanResponse>,
			},

			// GET - GET /:id
			get: {
				method: 'GET',
				path: '/:id',
				pathParams: {} as MigrationPlanParams,
				paramSchema: ZodMigrationPlanParamsSchema,
				responseType: {} as MigrationPlanResponse,
			},

			// EXECUTIONS - GET /:id/executions
			executions: {
				method: 'GET',
				path: '/:id/executions',
				pathParams: {} as MigrationPlanParams,
				paramSchema: ZodMigrationPlanParamsSchema,
				queryParams: {} as MigrationExecutionQuery,
				querySchema: ZodMigrationExecutionQuerySchema,
				responseType: {} as PaginatedResponse<MigrationExecutionResponse>,
			},

			// CREATE - POST /
			create: {
				method: 'POST',
				path: '',
				bodyParams: {} as MigrationPlanCreateRequest,
				bodySchema: ZodMigrationPlanCreateSchema,
				responseType: {} as MigrationPlanResponse,
			},

			// UPDATE - PATCH /:id
			update: {
				method: 'PATCH',
				path: '/:id',
				pathParams: {} as MigrationPlanParams,
				paramSchema: ZodMigrationPlanParamsSchema,
				bodyParams: {} as MigrationPlanUpdateRequest,
				bodySchema: ZodMigrationPlanUpdateSchema,
				responseType: {} as MigrationPlanResponse,
			},

			// CANCEL - DELETE /:id
			cancel: {
				method: 'DELETE',
				path: '/:id',
				pathParams: {} as MigrationPlanParams,
				paramSchema: ZodMigrationPlanParamsSchema,
				responseType: {} as { success: boolean },
			},
		},
	});
}

export type UserApiMethods = ReturnType<typeof userApi>;

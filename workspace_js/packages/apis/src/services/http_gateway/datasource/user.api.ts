import type {
	DataSourceId,
	DataSourceCreate,
	DataSourceUpdate,
	DataSourceQuery,
	DataSourceResponse,
	PaginatedResponse,
} from '@mapexos/schemas';
import type { GetToken } from '@src/common';
import { AxiosInstance } from 'axios';
import { createApiFactory } from '@src/common';
import {
	ZodDataSourceIdSchema,
	ZodDataSourceCreateSchema,
	ZodDataSourceUpdateSchema,
	ZodDataSourceQuerySchema,
} from '@mapexos/schemas';

/**
 * Creates DataSource user API with external CRUD endpoints.
 *
 * @param http - Axios instance for HTTP requests
 * @param getToken - Function to retrieve JWT token
 * @returns Object containing data source external API methods
 */
export function userApi(http: AxiosInstance, getToken: GetToken | undefined) {
	const factory = createApiFactory(http);

	return factory({
		basePath: '/api/v1/data_sources',
		useAuthJWT: true,
		getToken,
		methods: {
			// LIST - GET /
			list: {
				method: 'GET',
				path: '',
				queryParams: {} as DataSourceQuery,
				querySchema: ZodDataSourceQuerySchema,
				responseType: {} as PaginatedResponse<DataSourceResponse>,
			},

			// CREATE - POST /
			create: {
				method: 'POST',
				path: '',
				bodyParams: {} as DataSourceCreate,
				bodySchema: ZodDataSourceCreateSchema,
				responseType: {} as DataSourceResponse,
			},

			// GET BY ID - GET /:dataSourceId
			getById: {
				method: 'GET',
				path: '/:dataSourceId',
				pathParams: {} as DataSourceId,
				paramSchema: ZodDataSourceIdSchema,
				responseType: {} as DataSourceResponse,
			},

			// UPDATE - PATCH /:dataSourceId
			update: {
				method: 'PATCH',
				path: '/:dataSourceId',
				pathParams: {} as DataSourceId,
				bodyParams: {} as DataSourceUpdate,
				paramSchema: ZodDataSourceIdSchema,
				bodySchema: ZodDataSourceUpdateSchema,
				responseType: {} as DataSourceResponse,
			},

			// DELETE - DELETE /:dataSourceId
			delete: {
				method: 'DELETE',
				path: '/:dataSourceId',
				pathParams: {} as DataSourceId,
				paramSchema: ZodDataSourceIdSchema,
				responseType: {} as { success: boolean },
			},
		},
	});
}

export type UserApiMethods = ReturnType<typeof userApi>;

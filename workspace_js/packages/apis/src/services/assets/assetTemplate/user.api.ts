import type {
	AssetTemplateId,
	AssetTemplateCreate,
	AssetTemplateUpdate,
	AssetTemplateQuery,
	AssetTemplateResponse,
	PaginatedResponse,
} from '@mapexos/schemas';
import type { GetToken } from '@src/common';
import { AxiosInstance } from 'axios';
import { createApiFactory } from '@src/common';
import {
	ZodAssetTemplateIdSchema,
	ZodAssetTemplateCreateSchema,
	ZodAssetTemplateUpdateSchema,
	ZodAssetTemplateQuerySchema,
} from '@mapexos/schemas';

/**
 * Creates Asset Template user API with external CRUD endpoints.
 *
 * Provides standard CRUD operations for asset templates including:
 * - List with pagination and filtering
 * - Create new templates
 * - Get template by ID
 * - Update template
 * - Delete template
 *
 * @param http - Axios instance for HTTP requests
 * @param getToken - Function to retrieve JWT token
 * @returns Object containing asset template API methods
 */
export function userApi(http: AxiosInstance, getToken: GetToken | undefined) {
	const factory = createApiFactory(http);

	return factory({
		basePath: '/api/v1/asset_templates',
		useAuthJWT: true,
		getToken,
		methods: {
			// LIST - GET /
			list: {
				method: 'GET',
				path: '',
				queryParams: {} as AssetTemplateQuery,
				querySchema: ZodAssetTemplateQuerySchema,
				responseType: {} as PaginatedResponse<AssetTemplateResponse>,
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
				bodyParams: {} as AssetTemplateCreate,
				bodySchema: ZodAssetTemplateCreateSchema,
				responseType: {} as AssetTemplateResponse,
			},

			// GET BY ID - GET /:assetTemplateId
			getById: {
				method: 'GET',
				path: '/:assetTemplateId',
				pathParams: {} as AssetTemplateId,
				paramSchema: ZodAssetTemplateIdSchema,
				responseType: {} as AssetTemplateResponse,
			},

			// UPDATE - PATCH /:assetTemplateId
			update: {
				method: 'PATCH',
				path: '/:assetTemplateId',
				pathParams: {} as AssetTemplateId,
				bodyParams: {} as AssetTemplateUpdate,
				paramSchema: ZodAssetTemplateIdSchema,
				bodySchema: ZodAssetTemplateUpdateSchema,
				responseType: {} as AssetTemplateResponse,
			},

			// DELETE - DELETE /:assetTemplateId
			delete: {
				method: 'DELETE',
				path: '/:assetTemplateId',
				pathParams: {} as AssetTemplateId,
				paramSchema: ZodAssetTemplateIdSchema,
				responseType: {} as { success: boolean },
			},

			// GET AVAILABLE FIELDS - GET /:assetTemplateId/available_fields
			getAvailableFields: {
				method: 'GET',
				path: '/:assetTemplateId/available_fields',
				pathParams: {} as AssetTemplateId,
				paramSchema: ZodAssetTemplateIdSchema,
				responseType: {} as { availableFields: string[] },
			},
		},
	});
}

export type UserApiMethods = ReturnType<typeof userApi>;

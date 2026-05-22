import type {
	ExecutionId,
	ExecutionQuery,
	ExecutionResponse,
	SignalRequest,
	PaginatedResponse,
} from '@mapexos/schemas';
import type { GetToken } from '@src/common';
import { AxiosInstance } from 'axios';
import { createApiFactory } from '@src/common';
import {
	ZodExecutionIdSchema,
	ZodExecutionQuerySchema,
	ZodSignalRequestSchema,
} from '@mapexos/schemas';

/**
 * Creates Workflow Execution user API with external endpoints.
 *
 * @param http - Axios instance for HTTP requests
 * @param getToken - Function to retrieve JWT token
 * @returns Object containing execution external API methods
 */
export function userApi(http: AxiosInstance, getToken: GetToken | undefined) {
	const factory = createApiFactory(http);

	return factory({
		basePath: '/api/v1/workflow_executions',
		useAuthJWT: true,
		getToken,
		methods: {
			// LIST - GET /
			list: {
				method: 'GET',
				path: '',
				queryParams: {} as ExecutionQuery,
				querySchema: ZodExecutionQuerySchema,
				responseType: {} as PaginatedResponse<ExecutionResponse>,
			},

			// GET BY ID - GET /:executionId
			getById: {
				method: 'GET',
				path: '/:executionId',
				pathParams: {} as ExecutionId,
				paramSchema: ZodExecutionIdSchema,
				responseType: {} as ExecutionResponse,
			},

			// CANCEL - POST /:executionId/cancel
			cancel: {
				method: 'POST',
				path: '/:executionId/cancel',
				pathParams: {} as ExecutionId,
				paramSchema: ZodExecutionIdSchema,
				responseType: {} as { success: boolean },
			},

			// SIGNAL - POST /:executionId/signal
			signal: {
				method: 'POST',
				path: '/:executionId/signal',
				pathParams: {} as ExecutionId,
				bodyParams: {} as SignalRequest,
				paramSchema: ZodExecutionIdSchema,
				bodySchema: ZodSignalRequestSchema,
				responseType: {} as { success: boolean },
			},
		},
	});
}

export type UserApiMethods = ReturnType<typeof userApi>;

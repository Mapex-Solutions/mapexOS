import type { ScriptTest } from '@mapexos/schemas';
import type { GetToken } from '@src/common';
import { AxiosInstance } from 'axios';
import { createApiFactory } from '@src/common';
import { ZodScriptTestSchema } from '@mapexos/schemas';

/** Response type for getSamplePayload endpoint */
export interface SamplePayloadResponse {
	samplePayload: Record<string, any>;
}

/**
 * Creates Scripts API for testing JavaScript code execution.
 *
 * @param http - Axios instance for HTTP requests
 * @param getToken - Function to retrieve JWT token
 * @returns Object containing scripts API methods
 */
export function userApi(http: AxiosInstance, getToken: GetToken | undefined) {
	const factory = createApiFactory(http);

	return factory({
		basePath: '/api/v1/scripts',
		useAuthJWT: true,
		getToken,
		methods: {
			// TEST SCRIPT - POST /test
			test: {
				method: 'POST',
				path: '/test',
				bodyParams: {} as ScriptTest,
				bodySchema: ZodScriptTestSchema,
				responseType: {} as any, // TODO: Define response type based on backend
			},

			// GET SAMPLE PAYLOAD - GET /templates/:templateId/sample_payload
			// Executes template scripts against scriptTest and returns processed payload
			getSamplePayload: {
				method: 'GET',
				path: '/templates/:templateId/sample_payload',
				pathParams: {} as { templateId: string },
				queryParams: {} as { orgId?: string; isSystem?: string },
				responseType: {} as SamplePayloadResponse,
			},
		},
	});
}

export type UserApiMethods = ReturnType<typeof userApi>;

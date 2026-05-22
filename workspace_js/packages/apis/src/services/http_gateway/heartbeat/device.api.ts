import type { AxiosInstance, AxiosRequestConfig } from 'axios';
import type { HeartbeatRequest } from '@mapexos/schemas';
import { ZodHeartbeatRequestSchema } from '@mapexos/schemas';

/**
 * Creates Heartbeat device API for HTTP-protocol assets in explicit mode.
 *
 * The endpoint authenticates per DataSource (apiKey/jwt/oauth2/ip_whitelist),
 * NOT via the user JWT used by other Mapex CRUD APIs. Callers must supply
 * the right auth header in `requestConfig.headers` (e.g. x-api-key).
 *
 * @param http - Axios instance for HTTP requests
 * @returns Object with `post` method to send a heartbeat
 *
 * @example
 * const heartbeatApi = createHeartbeatApi(http);
 * await heartbeatApi.post(
 *     'data-source-id',
 *     { assetUUID: 'asset-uuid' },
 *     { headers: { 'x-api-key': 'device-key' } },
 * );
 */
export function deviceApi(http: AxiosInstance) {
	return {
		/**
		 * POST /api/v1/heartbeat?ds={dataSourceId}
		 *
		 * Body { assetUUID } is validated client-side against ZodHeartbeatRequestSchema
		 * before sending. Server returns 204 on success, 422 on validation, 403 on
		 * disabled DataSource, 401 on auth failure.
		 */
		async post(
			dataSourceId: string,
			body: HeartbeatRequest,
			requestConfig?: AxiosRequestConfig,
		): Promise<void> {
			const validated = ZodHeartbeatRequestSchema.parse(body);
			await http.post(
				`/api/v1/heartbeat?ds=${encodeURIComponent(dataSourceId)}`,
				validated,
				requestConfig,
			);
		},
	};
}

import type { AxiosInstance } from 'axios';

import { userApi } from './user.api';

/**
 * Creates Device module API. Owns the device-side refresh-token endpoint;
 * the bearer used here is the device's NATS user JWT, not the operator's
 * session JWT (see ./user.api.ts for the rationale).
 *
 * @param http - Axios instance for HTTP requests
 * @returns Object with device API methods
 *
 * @example
 * const deviceApi = createDeviceApi(http);
 * await deviceApi.refreshToken(currentJWT);
 */
export function createDeviceApi(http: AxiosInstance) {
	return {
		...userApi(http),
	};
}

export type DeviceApiMethods = ReturnType<typeof createDeviceApi>;

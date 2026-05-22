import type { AxiosInstance } from 'axios';

import { deviceApi } from './device.api';

/**
 * Creates Heartbeat module API for HTTP-protocol assets in explicit mode.
 *
 * The heartbeat endpoint authenticates per DataSource (not via user JWT),
 * so callers supply the right auth header per request.
 *
 * @param http - Axios instance for HTTP requests
 * @returns Object with heartbeat API methods
 */
export function createHeartbeatApi(http: AxiosInstance) {
	return {
		...deviceApi(http),
	};
}

export type HeartbeatApiMethods = ReturnType<typeof createHeartbeatApi>;

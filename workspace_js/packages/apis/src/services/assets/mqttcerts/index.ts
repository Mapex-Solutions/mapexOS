import type { AxiosInstance } from 'axios';
import type { GetToken } from '@src/common';

import { userApi } from './user.api';

/**
 * Creates the MQTT certificates module API. The endpoints wrap the
 * mqttcerts bounded context on the assets service: issue, revoke,
 * and list revoked certs for a given asset. User methods (JWT) are
 * spread to the root level since there is no internal counterpart.
 */
export function createMqttCertsApi(http: AxiosInstance, getToken: GetToken | undefined) {
	return {
		...userApi(http, getToken),
	};
}

export type MqttCertsApiMethods = ReturnType<typeof createMqttCertsApi>;

export * from './user.api';

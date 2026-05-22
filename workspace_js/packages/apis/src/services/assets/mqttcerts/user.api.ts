import type { AxiosInstance } from 'axios';
import type {
	IssueCertRequest,
	IssueCertResponse,
	RevokedCertResponse,
} from '@mapexos/schemas';

import type { GetToken } from '@src/common';

/**
 * mqttcerts user API. Wraps the JWT-gated external endpoints:
 *   POST   /api/v1/mqtt_certs              issue
 *   DELETE /api/v1/mqtt_certs/:serial      revoke
 *   GET    /api/v1/mqtt_certs?assetUUID=   list revoked
 *
 * The factory pattern used by other api modules (asset/user.api.ts) is
 * not strictly needed here — endpoints are simple JSON, no list/cursor
 * pagination. We use the http instance directly.
 */
export function userApi(http: AxiosInstance, getToken: GetToken | undefined) {
	async function authHeader(): Promise<Record<string, string>> {
		const token = (await getToken?.()) ?? '';
		return token ? { Authorization: `Bearer ${token}` } : {};
	}

	return {
		/**
		 * Issue a new device cert. Returns the PEM bundle once; the
		 * frontend MUST zip + download immediately. On 409 (asset has
		 * an existing currentCert), the caller can retry with force=true.
		 */
		async issueCert(assetUUID: string, force = false): Promise<IssueCertResponse> {
			const headers = await authHeader();
			const body: IssueCertRequest = { assetUUID, force };
			const res = await http.post('/api/v1/mqtt_certs', body, { headers });
			return res.data?.data ?? res.data;
		},

		/**
		 * Revoke a cert by serial. Moves the row to mqttRevokedCertificates.
		 */
		async revokeCert(serial: string): Promise<void> {
			const headers = await authHeader();
			await http.delete(`/api/v1/mqtt_certs/${encodeURIComponent(serial)}`, { headers });
		},

		/**
		 * List revoked certs for a single asset (TTL 30d on the server).
		 */
		async listRevoked(assetUUID: string): Promise<RevokedCertResponse[]> {
			const headers = await authHeader();
			const res = await http.get('/api/v1/mqtt_certs', { headers, params: { assetUUID } });
			return res.data?.data ?? res.data;
		},
	};
}

export type MqttCertsUserApiMethods = ReturnType<typeof userApi>;

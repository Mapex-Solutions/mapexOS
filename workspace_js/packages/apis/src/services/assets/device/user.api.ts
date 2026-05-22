import type { AxiosInstance } from 'axios';
import { type RefreshTokenResponse, ZodRefreshTokenResponseSchema } from '@mapexos/schemas';

/**
 * Creates Device user API for the device-side refresh-token endpoint.
 *
 * Differs from the operator-facing `userApi` factories in this package
 * because the auth posture is fundamentally different — the bearer here
 * is the device's CURRENT JWT (NATS user JWT signed by the platform
 * signing key), NOT the operator's session JWT. The shared factory's
 * `getToken` callback is closure-bound at construction time and cannot
 * be overridden per call, so this module bypasses the factory and
 * builds the request directly via the shared HTTP instance.
 *
 * @param http - Axios instance for HTTP requests (shared baseURL + interceptors)
 * @returns Object containing device-side API methods
 */
export function userApi(http: AxiosInstance) {
	return {
		/**
		 * Refreshes the device's MQTT JWT before expiry. The platform
		 * verifies the bearer's signature + trust-anchors the issuer,
		 * asserts the JWT's jti still matches the asset's persisted jti,
		 * and returns a rotated JWT plus the new user pub key + expiry.
		 *
		 * @param currentJWT - The device's CURRENT JWT (passed as Bearer)
		 * @returns The rotated credential payload, Zod-parsed at the boundary
		 */
		async refreshToken(currentJWT: string): Promise<RefreshTokenResponse> {
			const res = await http.post('/api/v1/devices/refresh_token', undefined, {
				headers: { Authorization: `Bearer ${currentJWT}` },
			});
			return ZodRefreshTokenResponseSchema.parse(res.data?.data);
		},
	};
}

export type DeviceUserApiMethods = ReturnType<typeof userApi>;

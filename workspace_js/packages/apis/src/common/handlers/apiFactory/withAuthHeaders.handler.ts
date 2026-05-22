import { mountHeadersWithJWT } from '@src/tools';
import { GetToken } from '@src/common';

/**
 * Conditionally adds authentication JWT headers to the provided headers object.
 *
 * This function enhances the given headers with JWT authentication headers when
 * authentication is required. If authentication is not needed, it returns the
 * original headers unchanged.
 *
 * @param useAuthJWT - Flag indicating whether to include JWT authentication headers.
 * @param getToken - Optional function to retrieve the authentication token.
 * @param headers - Optional base headers object to extend with authentication headers.
 * @returns A promise that resolves to the final headers object, with or without authentication headers.
 */
export async function withAuthHeaders(
	useAuthJWT: boolean,
	getToken?: GetToken,
	headers: Record<string, string> = {},
) {
	return useAuthJWT ? { ...headers, ...(await mountHeadersWithJWT(getToken)) } : headers;
}
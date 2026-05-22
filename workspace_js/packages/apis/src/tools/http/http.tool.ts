import axios, { AxiosInstance, CreateAxiosDefaults } from 'axios';

import { ApiConfig, ApiInterceptors, GetToken } from '@src/common';
import { isEmpty } from 'lodash';

/**
 * Creates an Axios HTTP client instance with the specified configuration.
 *
 * @param {ApiConfig} config - The configuration object for the HTTP client.
 * @param {ApiInterceptors} interceptors - Interceptors for the HTTP client.
 * @returns An Axios instance configured with the provided base URL and headers.
 */
export function createHttp(config: ApiConfig, interceptors?: ApiInterceptors): AxiosInstance {
	const { baseURL, headers = {}, httpsAgent } = config;

	const newInstanceParams: CreateAxiosDefaults = {
		baseURL,
		headers: { 'Content-Type': 'application/json', ...headers },
	};

	if (httpsAgent) newInstanceParams.httpsAgent = httpsAgent;
	const httpInstance = axios.create(newInstanceParams);

	if (!isEmpty(interceptors) && interceptors?.onRequest) {
		httpInstance.interceptors.request.use(interceptors.onRequest);
	}

	if (!isEmpty(interceptors) && interceptors?.onResponse) {
		httpInstance.interceptors.response.use(interceptors.onResponse);
	}

	if (!isEmpty(interceptors) && interceptors?.onError) {
		httpInstance.interceptors.response.use(null, interceptors.onError);
	}

	return httpInstance;
}

/**
 * Asynchronously constructs headers with a JSON Web Token (JWT) for authorization.
 *
 * @param getToken - A function that retrieves the JWT token. It should return a promise that resolves to the token string.
 * @returns A promise that resolves to an object containing the Authorization header with the JWT, or an empty object if the token is not available.
 */
export async function mountHeadersWithJWT(getToken: GetToken | undefined) {
	if (getToken) {
		const token = await getToken();
		return !isEmpty(token) ? { Authorization: `Bearer ${token}` } : {};
	}
	return {};
}
import type { ApiMethod, ApiClientConfig, SessionsConfig, ApiInterceptors } from '@src/common';
import type { AxiosInstance, AxiosResponse, AxiosRequestHeaders } from 'axios';

import axios from 'axios';
import { isEmpty } from 'lodash';
import { SchemaError, zodValidationError } from '@mapexos/validations';
import { replacePathParams } from './replacePathParams.handler';

/**
 * This type needs to stay here because the TS is losssing the reference to the ApiMethod type.
 */
// type ApiClient<T extends Record<string, ApiMethod<any, any, any, any>>> = {
// 	[K in keyof T]: T[K] extends ApiMethod<infer P, infer Q, infer B, infer R>
// 		? P extends {}
// 			? B extends {}
// 				? Q extends {}
// 					? (params: P, body: B, query: Q) => Promise<R>
// 					: (params: P, body: B) => Promise<R>
// 				: Q extends {}
// 					? (params: P, query: Q) => Promise<R>
// 					: (params: P) => Promise<R>
// 			: B extends {}
// 				? Q extends {}
// 					? (body: B, query: Q) => Promise<R>
// 					: (body: B) => Promise<R>
// 				: Q extends {}
// 					? (query: Q) => Promise<R>
// 					: () => Promise<R>
// 		: never
// };

type InferResponseType<R> = unknown extends R ? any : R;
;

type ApiClient<T extends Record<string, ApiMethod<any, any, any, any>>> = {
	[K in keyof T]: T[K] extends ApiMethod<infer P, infer Q, infer B, infer R>
		? P extends {}
			? B extends {}
				? Q extends {}
					? (params: P, body: B, query: Q) => Promise<InferResponseType<R>>
					: (params: P, body: B) => Promise<InferResponseType<R>>
				: Q extends {}
					? (params: P, query: Q) => Promise<InferResponseType<R>>
					: (params: P) => Promise<InferResponseType<R>>
			: B extends {}
				? Q extends {}
					? (body: B, query: Q) => Promise<InferResponseType<R>>
					: (body: B) => Promise<InferResponseType<R>>
				: Q extends {}
					? (query: Q) => Promise<InferResponseType<R>>
					: () => Promise<InferResponseType<R>>
		: never
};

/**
 * Creates an API factory that generates API clients with typed methods.
 *
 * This factory function takes an Axios instance and returns a function that can create
 * API clients based on a configuration object. The created clients have methods that
 * automatically handle path parameters, query parameters, and request bodies according
 * to the API method definitions.
 *
 * @param http - The Axios instance to use for making HTTP requests
 * @returns A function that creates API clients based on configuration
 */
export function createApiFactory(http: AxiosInstance, sessionsConfig?: SessionsConfig) {

	const interceptors = !isEmpty(sessionsConfig?.interceptors) ? sessionsConfig?.interceptors : {} as ApiInterceptors;

	/**
	 * Creates an API client with typed methods based on the provided configuration.
	 *
	 * @param {ApiClientConfig} config - The configuration object that defines the API client's behavior,
	 *                 including base path and method definitions
	 * @returns A typed API client object with methods corresponding to the configuration
	 */
	return function createApiClient<T extends Record<string, ApiMethod<any, any, any>>>(
		config: ApiClientConfig<T>,
	): ApiClient<T> {
		const client = {} as ApiClient<T>;

		for (const [methodName, methodConfig] of Object.entries(config.methods)) {
			const fullPath = config.basePath + methodConfig.path;

			/**
			 * Create the method that makes the API request
			 * @type {ApiClient<T>[keyof T]}
			 */
			client[methodName as keyof T] = (async (...args: any[]) => {
				let url = fullPath;

				let params;
				let body;
				let query;

				let headers: Record<string, any> = { 'Content-Type': 'application/json' };

				/** Extract params body and query */
				if (methodConfig?.pathParams) params = args.shift();
				if (methodConfig?.bodyParams) body = args.shift();
				if (methodConfig?.queryParams) query = args.shift();

				/** Call the before request hook */
				if (methodConfig?.beforeRequest) {

					/**
					 * Method may modify the params, body, or query before making the request
					 * @return { params, body, query } or a Promise that resolves to them
					 */
					const { params: _params, body: _body, query: _query } = await methodConfig.beforeRequest(params, body, query);

					params = _params;
					body = _body;
					query = _query;
				}

				/**
				 * Validate the request using schemas
				 */
				if (methodConfig?.paramSchema) {
					const requestPayload = await methodConfig?.paramSchema.safeParseAsync(params);
					if (!requestPayload.success) throw new SchemaError(zodValidationError(requestPayload));
					else params = requestPayload.data;
				}

				if (methodConfig?.bodySchema) {
					const requestPayload = await methodConfig?.bodySchema.safeParseAsync(body);
					if (!requestPayload.success) throw new SchemaError(zodValidationError(requestPayload));
					else body = requestPayload.data;
				}

				if (methodConfig?.querySchema) {
					const requestPayload = await methodConfig?.querySchema.safeParseAsync(query);
					if (!requestPayload.success) throw new SchemaError(zodValidationError(requestPayload));
					else query = requestPayload.data;
				}


				/** * Replace path parameters in the URL */
				if (params) url = replacePathParams(url, params);

				/**
				 * Add query parameters to the URL.
				 *
				 * Arrays are serialized as repeated keys (e.g. `?kinds=a&kinds=b`)
				 * to match fiber's `QueryParser` slice binding (which expects repeated
				 * keys when `EnableSplittingOnParsers` is `false` — the default).
				 *
				 * `null` / `undefined` values are skipped.
				 */
				if (query) {
					const params = new URLSearchParams();
					for (const [key, value] of Object.entries(query as Record<string, unknown>)) {
						if (value === undefined || value === null) continue;
						if (Array.isArray(value)) {
							for (const item of value) {
								if (item === undefined || item === null) continue;
								params.append(key, String(item));
							}
						} else {
							params.append(key, String(value));
						}
					}
					const queryString = params.toString();
					if (queryString) {
						url += `?${queryString}`;
					}
				}

				/** Check the headers */
				if (methodConfig?.headers) {
					headers = { ...headers, ...methodConfig.headers };
				}

				/**
				 * Set the JWT or APIKey
				 */
				if (config.useAuthJWT && config.getToken) {
					const token = await config.getToken();
					if (token) {
						headers['Authorization'] = `Bearer ${token}`;
					}
				}

				/**
				 * Make the API request using the provided Axios instance
				 */
				let axiosInstance = http;

				if (methodConfig.newAxiosInstance) {
					axiosInstance = axios.create({ ...http.defaults, baseURL: `${http.defaults.baseURL}` });

					if (!isEmpty(interceptors) && interceptors?.onRequest) {
						axiosInstance.interceptors.request.use(interceptors.onRequest);
					}

					if (!isEmpty(interceptors) && interceptors?.onResponse) {
						axiosInstance.interceptors.response.use(interceptors.onResponse);
					}

					if (!isEmpty(interceptors) && interceptors?.onError) {
						axiosInstance.interceptors.response.use(null, interceptors.onError);
					}
				}

				const restData = await axiosInstance.request({
					url,
					method: methodConfig.method,
					data: body,
					headers,
					...methodConfig.axiosConfig,
				});

				/**
				 * Add the after request hook if defined
				 */
				if (methodConfig?.afterRequest) return methodConfig.afterRequest(restData);
				return restData?.data?.data;

			}) as ApiClient<T>[keyof T];
		}

		return client;
	};
}
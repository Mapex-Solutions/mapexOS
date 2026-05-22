import type { AxiosResponse, AxiosRequestConfig, AxiosHeaders } from 'axios';

import { GetToken } from '@src/common';


export type InferResponseType<R> = {} extends R ? any : R;

/**
 * We will define the method type using the AxiosRequestConfig and AxiosResponse types.
 */
export interface ApiMethod<P = {}, Q = {}, B = {}, R = {}> {
	newAxiosInstance?: boolean;

	// Signature of the method
	pathParams?: P;
	queryParams?: Q;
	bodyParams?: B;

	// Request path and HTTP method
	path: string;
	method: 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE';
	headers?: Record<symbol, any>;
	axiosConfig?: AxiosRequestConfig;

	// Response type and any additional settings for Axios.get, Axios.post, etc.
	responseType?: R;

	// Hooks
	beforeRequest?: (params: P, body: B, query: Q) => Promise<{ params: P, body: B, query: Q }> | {
		params: P,
		body: B,
		query: Q
	};

	afterRequest?: (response: AxiosResponse) => Promise<R> | R;

	// Validation schema
	paramSchema?: any;
	bodySchema?: any;
	querySchema?: any;
}

// Definindo a configuração do client
export interface ApiClientConfig<T extends Record<string, ApiMethod<any, any, any>>> {
	basePath: string;
	useAuthJWT?: boolean;
	getToken?: GetToken;
	methods: T;
}
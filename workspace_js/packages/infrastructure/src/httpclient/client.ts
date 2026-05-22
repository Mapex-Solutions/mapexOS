/**
 * HTTPClient
 * Same structure as workspace_go/packages/infrastructure/httpclient/client.go
 *
 * Provides a generic HTTP client for making requests to external services.
 * It supports API Key authentication and automatic JSON serialization/deserialization.
 *
 * Features:
 *   - Configurable base URL and API Key
 *   - Automatic JSON marshaling/unmarshaling
 *   - Timeout support
 *   - Generic response type handling
 *
 * Example usage:
 *   const client = new HTTPClient({
 *     baseURL: 'http://localhost:5003',
 *     apiKey: 'my-secret-key',
 *     timeout: 5000,
 *   });
 *
 *   const result = await client.get<RouteGroupResponse[]>('/api/internal/v1/routegroups?ids=id1,id2');
 */

import type { HTTPClientConfig } from './types';
import { HTTPError } from './types';

const DEFAULT_TIMEOUT = 10000; // 10 seconds

export class HTTPClient {
	private baseURL: string;
	private apiKey?: string;
	private timeout: number;

	constructor(config: HTTPClientConfig) {
		this.baseURL = config.baseURL;
		this.apiKey = config.apiKey;
		this.timeout = config.timeout ?? DEFAULT_TIMEOUT;
	}

	/**
	 * Get performs a GET request to the specified endpoint.
	 *
	 * @param endpoint - API endpoint path (e.g., "/api/internal/v1/routegroups?ids=id1,id2")
	 * @returns The response data parsed as type T
	 */
	async get<T>(endpoint: string): Promise<T> {
		return this.doRequest<T>('GET', endpoint);
	}

	/**
	 * Post performs a POST request to the specified endpoint.
	 *
	 * @param endpoint - API endpoint path
	 * @param body - Request body to be marshaled to JSON
	 * @returns The response data parsed as type T
	 */
	async post<T>(endpoint: string, body?: unknown): Promise<T> {
		return this.doRequest<T>('POST', endpoint, body);
	}

	/**
	 * Put performs a PUT request to the specified endpoint.
	 */
	async put<T>(endpoint: string, body?: unknown): Promise<T> {
		return this.doRequest<T>('PUT', endpoint, body);
	}

	/**
	 * Delete performs a DELETE request to the specified endpoint.
	 */
	async delete<T>(endpoint: string): Promise<T> {
		return this.doRequest<T>('DELETE', endpoint);
	}

	/**
	 * doRequest performs the actual HTTP request with proper error handling.
	 */
	private async doRequest<T>(method: string, endpoint: string, body?: unknown): Promise<T> {
		const url = this.baseURL + endpoint;

		const headers: Record<string, string> = {
			'Content-Type': 'application/json',
		};

		if (this.apiKey) {
			headers['X-API-Key'] = this.apiKey;
		}

		const controller = new AbortController();
		const timeoutId = setTimeout(() => controller.abort(), this.timeout);

		try {
			const response = await fetch(url, {
				method,
				headers,
				body: body ? JSON.stringify(body) : undefined,
				signal: controller.signal,
			});

			clearTimeout(timeoutId);

			const responseText = await response.text();

			// Check for non-2xx status codes
			if (!response.ok) {
				throw new HTTPError(response.status, responseText);
			}

			// Parse JSON response
			if (!responseText) {
				return {} as T;
			}

			return JSON.parse(responseText) as T;
		} catch (error) {
			clearTimeout(timeoutId);

			if (error instanceof HTTPError) {
				throw error;
			}

			if (error instanceof Error && error.name === 'AbortError') {
				throw new Error(`Request timeout after ${this.timeout}ms`);
			}

			throw new Error(`Request failed: ${error}`);
		}
	}
}

/**
 * HTTPClient Types
 * Same structure as workspace_go/packages/infrastructure/httpclient/client.go
 */

/**
 * HTTPClientConfig defines the configuration for creating a new HTTPClient.
 */
export interface HTTPClientConfig {
	/** Base URL of the service (e.g., "http://localhost:5003") */
	baseURL: string;
	/** API Key for authentication (sent as X-API-Key header) */
	apiKey?: string;
	/** Request timeout in milliseconds (default: 10000) */
	timeout?: number;
}

/**
 * HTTPError represents an HTTP error response.
 */
export class HTTPError extends Error {
	constructor(
		public statusCode: number,
		public body: string,
	) {
		super(`Request failed with status ${statusCode}: ${body}`);
		this.name = 'HTTPError';
	}
}

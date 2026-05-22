import { isEmpty } from 'lodash';
import { ApiConfig, ApiInterceptors, SessionsConfig, GetToken } from '@src/common';
import { createHttp } from '@src/tools';

import { createAssetApi } from './asset';
import { createAssetTemplateApi } from './assetTemplate';
import { createDeviceApi } from './device';
import { createMqttCertsApi } from './mqttcerts';

/**
 * Creates complete Assets Service API with all modules.
 *
 * Currently includes:
 * - asset: Asset CRUD operations + internal scripts endpoint
 * - assetTemplate: Asset template CRUD operations
 *
 * @param config - API configuration (baseURL, timeout, headers, etc.)
 * @param sessionsConfig - Optional session configuration (interceptors, getToken)
 * @returns Object with HTTP instance and all module APIs
 *
 * @example
 * const assetsApi = createAssetsApi({
 *   baseURL: 'http://localhost:5002',
 *   timeout: 30000
 * }, {
 *   getToken: () => localStorage.getItem('token'),
 *   interceptors: myInterceptors
 * });
 *
 * // Usage:
 * await assetsApi.asset.list({ page: 1, perPage: 20 });
 * await assetsApi.asset.create({ name: 'Sensor 001', ... });
 * await assetsApi.asset.internal.getScriptsByAssetUUID({ assetUUID: 'dev-eui-12345' });
 * await assetsApi.assetTemplate.list({ page: 1 });
 * await assetsApi.assetTemplate.create({ name: 'LoRaWAN Template', ... });
 */
export function createAssetsApi(config: ApiConfig, sessionsConfig?: SessionsConfig) {
	/** Global/Local interceptors for the assets service */
	const interceptors = !isEmpty(sessionsConfig?.interceptors)
		? sessionsConfig?.interceptors
		: config?.interceptors || {} as ApiInterceptors;

	const getToken: GetToken | undefined = sessionsConfig?.getToken;

	/** Create single HTTP instance for the entire service */
	const http = createHttp(config, interceptors);

	/** Initialize all modules with the same HTTP instance */
	const asset = createAssetApi(http, getToken);
	const assetTemplate = createAssetTemplateApi(http, getToken);
	const device = createDeviceApi(http);
	const mqttcerts = createMqttCertsApi(http, getToken);

	return {
		http,          // Expose HTTP instance
		asset,         // Asset module (includes user + internal methods)
		assetTemplate, // Asset Template module (CRUD operations)
		device,        // Device module (device-side JWT refresh)
		mqttcerts,     // MQTT cert lifecycle (issue, revoke, list revoked)
	};
}

/**
 * Type for the complete Assets Service API
 */
export type AssetsApi = ReturnType<typeof createAssetsApi>;

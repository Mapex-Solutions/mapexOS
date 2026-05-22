/**
 * Test data factories for asset registration E2E tests
 */

/**
 * Asset identification step data
 */
export interface AssetIdentificationData {
  /** Asset name (required) */
  name: string;
  /** Asset ID (required) */
  assetId: string;
  /** Status: true = Active, false = Inactive */
  enabled?: boolean;
  /** Optional description */
  description?: string;
  /** Debug mode toggle */
  debugEnabled?: boolean;
}

/**
 * Asset connectivity step data
 */
export interface AssetConnectivityData {
  /** Communication protocol */
  protocol: 'HTTP' | 'MQTT';
  /** MQTT username (required when protocol is MQTT) */
  mqttUsername?: string;
  /** MQTT client ID (required when protocol is MQTT) */
  mqttClientId?: string;
  /** Device JWT TTL (required when protocol is MQTT). Format: 30d, 1y, 5y, 10y. */
  mqttTokenTTL?: string;
  /** Geographic latitude */
  latitude?: number;
  /** Geographic longitude */
  longitude?: number;
}

/**
 * Full asset registration test data
 */
export interface AssetRegistrationData {
  /** Step 1: Identification */
  identification: AssetIdentificationData;
  /** Step 4: Connectivity */
  connectivity: AssetConnectivityData;
}

/**
 * Generate a valid asset registration dataset
 *
 * @param {Partial<AssetRegistrationData>} overrides - Fields to override
 * @returns {AssetRegistrationData} Complete test data
 */
export function createAssetData(overrides?: Partial<AssetRegistrationData>): AssetRegistrationData {
  const timestamp = Date.now();
  return {
    identification: {
      name: `Test Asset ${timestamp}`,
      assetId: `A${timestamp}`,
      enabled: true,
      description: 'E2E test asset',
      debugEnabled: false,
      ...overrides?.identification,
    },
    connectivity: {
      protocol: 'HTTP',
      ...overrides?.connectivity,
    },
  };
}

/**
 * Generate minimal asset data (only required fields)
 *
 * @returns {AssetRegistrationData} Minimal valid data
 */
export function createMinimalAssetData(): AssetRegistrationData {
  const timestamp = Date.now();
  return {
    identification: {
      name: `Min Asset ${timestamp}`,
      assetId: `A${timestamp}`,
    },
    connectivity: {
      protocol: 'HTTP',
    },
  };
}

/**
 * Generate asset data with MQTT connectivity
 *
 * @param {Partial<AssetRegistrationData>} overrides - Fields to override
 * @returns {AssetRegistrationData} Asset data with MQTT protocol
 */
export function createMqttAssetData(overrides?: Partial<AssetRegistrationData>): AssetRegistrationData {
  const timestamp = Date.now();
  return {
    identification: {
      name: `MQTT Asset ${timestamp}`,
      assetId: `M${timestamp}`,
      enabled: true,
      ...overrides?.identification,
    },
    connectivity: {
      protocol: 'MQTT',
      mqttUsername: 'test-user',
      mqttClientId: 'client-001',
      mqttTokenTTL: '1y',
      ...overrides?.connectivity,
    },
  };
}

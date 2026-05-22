import type { AssetFormData, MqttConfig, HealthMonitorFormConfig } from '../interfaces';
import { MQTT_AUTH_TYPE_CERT } from '../interfaces/createEditAsset.interface';

/**
 * Initial MQTT configuration values. AuthType defaults to `cert` —
 * mTLS is the recommended path because the platform issues the cert
 * end-to-end and the broker enforces the mode at CONNECT. Operators
 * who want password-based devices flip the selector in Step 4.
 */
export const INITIAL_MQTT_CONFIG: MqttConfig = {
  clientId: '',
  username: '',
  authType: MQTT_AUTH_TYPE_CERT,
  password: '',
  certTTL: { value: 1, unit: 'year' },
};

/**
 * Initial health monitoring configuration values
 */
export const INITIAL_HEALTH_MONITOR: HealthMonitorFormConfig = {
  enabled: false,
  thresholdMinutes: 10,
  requiredMisses: 3,
  offlineRouteGroupIds: [],
  onlineRouteGroupIds: [],
  selectedOfflineRouteGroups: [],
  selectedOnlineRouteGroups: [],
};

/**
 * Initial form data values
 */
export const INITIAL_ASSET_FORM_DATA: AssetFormData = {
  name: '',
  assetId: '',
  enabled: true,
  description: '',
  assetTemplateId: null,
  routeGroupIds: [],
  protocol: 'HTTP',
  latitude: null,
  longitude: null,
  mqttConfig: { ...INITIAL_MQTT_CONFIG },
  debugEnabled: false,
  healthMonitor: { ...INITIAL_HEALTH_MONITOR },
};

/**
 * Total number of steps in the form
 */
export const TOTAL_STEPS = 6;

/**
 * Step numbers enum for better readability
 */
export const STEP = {
  IDENTIFICATION: 1,
  ASSET_TEMPLATE: 2,
  ROUTE_GROUPS: 3,
  CONNECTIVITY: 4,
  HEALTH_MONITORING: 5,
  REVIEW: 6,
} as const;

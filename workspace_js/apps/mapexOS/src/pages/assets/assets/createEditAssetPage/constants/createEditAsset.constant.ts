import type { AssetFormData, MqttConfig, HealthMonitorFormConfig, LorawanFormConfig, AssetTypeOption } from '../interfaces';
import { MQTT_AUTH_TYPE_CERT, LORAWAN_KIND_DEVICE } from '../interfaces/createEditAsset.interface';

/**
 * Initial MQTT configuration values. AuthType defaults to `cert` —
 * mTLS is the recommended path because the platform issues the cert
 * end-to-end and the broker enforces the mode at CONNECT. Operators
 * who want password-based devices flip the selector in the connectivity step.
 */
export const INITIAL_MQTT_CONFIG: MqttConfig = {
  clientId: '',
  username: '',
  authType: MQTT_AUTH_TYPE_CERT,
  password: '',
  certTTL: { value: 1, unit: 'year' },
};

/**
 * Initial LoRaWAN configuration values. Defaults to a `device` (the common
 * case) with OTAA activation, Class A, EU868 region. The operator picks the
 * gateway card on the first type step to switch to the gateway field set.
 * Key material starts blank — never pre-filled, never returned from the server.
 */
export const INITIAL_LORAWAN_CONFIG: LorawanFormConfig = {
  kind: LORAWAN_KIND_DEVICE,
  joinEui: '',
  region: '',
  class: 'A',
  macVersion: '1.0.3',
  phyVersion: '',
  activation: 'otaa',
  appKey: '',
  nwkKey: '',
  devAddr: '',
  nwkSKey: '',
  appSKey: '',
  gateway: {
    authMode: 'cert',
    frequencyPlanId: '',
    frequencyPlanIds: [],
    certTTL: { value: 1, unit: 'year' },
    apiKey: '',
  },
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
  attributes: [],
  protocol: 'HTTP',
  latitude: null,
  longitude: null,
  mqttConfig: { ...INITIAL_MQTT_CONFIG },
  lorawanConfig: { ...INITIAL_LORAWAN_CONFIG, gateway: { ...INITIAL_LORAWAN_CONFIG.gateway } },
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
/**
 * Wizard step ids. The step list is dynamic (a LoRaWAN gateway skips
 * assetTemplate + routeGroups), so steps are keyed by id, not a fixed number.
 */
export const STEP = {
  TYPE: 'type',
  IDENTIFICATION: 'identification',
  ATTRIBUTES: 'attributes',
  ASSET_TEMPLATE: 'assetTemplate',
  ROUTE_GROUPS: 'routeGroups',
  CONNECTIVITY: 'connectivity',
  HEALTH_MONITORING: 'health',
  REVIEW: 'review',
} as const;

/**
 * First-step type cards: a flat four-way choice that sets protocol (+ kind for
 * the LoRaWAN cards). Labels/descriptions live in i18n, keyed by `value`.
 */
export const ASSET_TYPE_OPTIONS: AssetTypeOption[] = [
  { value: 'mqtt', icon: 'mdi-transit-connection-variant', protocol: 'MQTT' },
  { value: 'http', icon: 'mdi-web', protocol: 'HTTP' },
  { value: 'lorawanSensor', icon: 'mdi-access-point', protocol: 'LORAWAN', kind: 'device' },
  { value: 'lorawanGateway', icon: 'mdi-router-wireless', protocol: 'LORAWAN', kind: 'gateway' },
];

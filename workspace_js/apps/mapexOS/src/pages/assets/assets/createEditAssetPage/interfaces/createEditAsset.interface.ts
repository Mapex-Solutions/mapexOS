import type { ComputedRef } from 'vue';
import type { AssetTemplateResponse, RouteGroupResponse } from '@mapexos/schemas';

/** Visual group a wizard leaf belongs to in the vertical stepper (UI-only). */
export interface StepMetaGroup {
  id: string;
  label: ComputedRef<string>;
  icon: string;
}

/** One wizard leaf's metadata: id, stepper icon/label/description, optional group. */
export interface StepMetaEntry {
  id: string;
  icon: string;
  label: ComputedRef<string>;
  description: ComputedRef<string>;
  group?: StepMetaGroup;
}

/**
 * Health monitoring configuration for the asset form.
 *
 * heartbeatMode chooses how the platform learns the device is alive:
 *   - 'implicit' (default): js-executor emits a heartbeat for every data event.
 *   - 'explicit': js-executor skips implicit publishes; the path is chosen by
 *     the asset's protocol — MQTT-protocol assets use NATS broker presence
 *     ($SYS.ACCOUNT.* CONNECT/DISCONNECT advisories, no device-side topic);
 *     HTTP-protocol assets POST /api/v1/heartbeat?ds={dataSourceId} with
 *     body { assetUUID }.
 */
export interface HealthMonitorFormConfig {
  enabled: boolean;
  thresholdMinutes: number;
  requiredMisses: number;
  heartbeatMode?: 'implicit' | 'explicit';
  offlineRouteGroupIds: string[];
  onlineRouteGroupIds: string[];
  selectedOfflineRouteGroups: RouteGroupResponse[];
  selectedOnlineRouteGroups: RouteGroupResponse[];
}

/**
 * Auth-mode sentinel values for MQTT assets. The asset declares one
 * mode at a time and the broker enforces mutual exclusion at CONNECT.
 * Mirrors the Go contract enum in
 * `packages/contracts/services/assets/assets/constants.go`.
 */
export const MQTT_AUTH_TYPE_PASSWORD = 'password';
export const MQTT_AUTH_TYPE_CERT = 'cert';
export type MqttAuthType = typeof MQTT_AUTH_TYPE_PASSWORD | typeof MQTT_AUTH_TYPE_CERT;

/**
 * MQTT configuration for asset connectivity.
 *
 * `authType` decides which credential the broker will accept for this
 * asset. `password` is plaintext on the form and only meaningful when
 * `authType` is `password` — the backend bcrypts before persisting and
 * the platform never exposes the plaintext again. Leaving the password
 * blank on edit (password mode) signals "no change" (the existing
 * hash on the asset is kept).
 *
 * `username` and `clientId` are platform-derived from the asset UUID
 * entered in Step 1; the form renders them readonly so operators
 * cannot drift them from the canonical lookup key the broker expects.
 */
export interface MqttConfig {
  /** MQTT client identifier — derived from assetUUID, readonly in the form. */
  clientId: string;

  /** MQTT username — derived as `${orgId}:${assetUUID}`, readonly in the form. */
  username: string;

  /** Auth mode the broker enforces at CONNECT. */
  authType: MqttAuthType;

  /**
   * MQTT password (plaintext, min 8 characters). Only used when
   * `authType` is `password`. Required on create in password mode;
   * optional on edit (blank keeps the existing hash). Cleared from
   * local form state immediately after a successful submit.
   */
  password: string;

  /**
   * Operator-declared validity window applied when the platform
   * signs this asset's MQTT device cert. Only meaningful when
   * `authType` is `cert`. Bounds enforced server-side (1 day .. 10
   * years total); the wizard mirrors the unit enum so the value /
   * unit pair stays in sync with the broker contract.
   */
  certTTL?: CertTTLConfig;
}

/**
 * Cert TTL split into (value, unit) so the wizard can render the
 * pair as (number-input, dropdown). Backend resolves the day count
 * via `value * unit-to-days` (day=1, week=7, month=30, year=365).
 */
export interface CertTTLConfig {
  value: number;
  unit: CertTTLUnit;
}

/** Allowed cert TTL units. Mirrors the Go `oneof=day week month year`. */
export type CertTTLUnit = 'day' | 'week' | 'month' | 'year';

/** Unit options used by the wizard's dropdown. */
export const CERT_TTL_UNITS: readonly CertTTLUnit[] = ['day', 'week', 'month', 'year'] as const;

/**
 * LoRaWAN asset kind. `kind` discriminates the two LoRaWAN asset shapes the
 * wizard collects: an end-device (identity + activation keys) or a gateway
 * (radio infrastructure: frequency plan + connection auth). The two render
 * disjoint field sets. Mirrors the Go `oneof=device gateway`.
 */
export type LorawanKind = 'device' | 'gateway';
export const LORAWAN_KIND_DEVICE = 'device';
export const LORAWAN_KIND_GATEWAY = 'gateway';

/** Device activation. OTAA derives session keys at join (appKey/nwkKey); ABP
 *  ships fixed session keys (devAddr/nwkSKey/appSKey). Mirrors Go `oneof=otaa abp`. */
export type LorawanActivation = 'otaa' | 'abp';

/** Device class. A = uplink-driven; B = scheduled ping slots; C = continuous RX. */
export type LorawanClass = 'A' | 'B' | 'C';

/** LoRaWAN MAC version. Mirrors Go `oneof=1.0.2 1.0.3 1.0.4 1.1`. */
export type LorawanMacVersion = '1.0.2' | '1.0.3' | '1.0.4' | '1.1';

/** Gateway connection auth: `eui` (UDP registered-EUI), `cert` (Basics Station
 *  mTLS), or `key` (Basics Station bearer token). */
export type LorawanGatewayAuthMode = 'eui' | 'cert' | 'key';

/**
 * LoRaWAN gateway block (kind === 'gateway'). Radio infrastructure: a frequency
 * plan (or plans) the Gateway Server tunes the radio to, plus the connection
 * auth mode. No device keys. Mirrors the Go LorawanGatewayConfig.
 */
export interface LorawanGatewayFormConfig {
  authMode: LorawanGatewayAuthMode;
  frequencyPlanId: string;
  frequencyPlanIds: string[];
  /** Only meaningful when authMode is `cert` (Basics Station mTLS). */
  certTTL?: CertTTLConfig;
  /** Bearer token, only meaningful when authMode is `key`. Request-only: the
   *  backend hashes it and never returns it, so it renders blank on edit and is
   *  sent only when filled. */
  apiKey: string;
}

/**
 * LoRaWAN configuration for asset connectivity. `kind` drives which fields the
 * form renders and submits: device identity + activation keys, or the gateway
 * block. Secret key material (appKey/nwkKey or the ABP session keys) is plaintext
 * on the form only — the backend envelope-encrypts it and never returns it, so
 * the fields render blank on edit and are sent only when filled. Mirrors the Go
 * LorawanConfig.
 */
export interface LorawanFormConfig {
  kind: LorawanKind;

  // Device identity + profile (kind === 'device'). The device EUI is the asset
  // UUID (Step 1), so it is not a field here — the handler uses assetId.
  joinEui: string;
  region: string;
  class: LorawanClass;
  macVersion: LorawanMacVersion;
  phyVersion: string;
  activation: LorawanActivation;

  // OTAA root keys (activation === 'otaa').
  appKey: string;
  nwkKey: string;

  // ABP session keys (activation === 'abp').
  devAddr: string;
  nwkSKey: string;
  appSKey: string;

  // Gateway block (kind === 'gateway').
  gateway: LorawanGatewayFormConfig;
}

/**
 * Custom-attribute kinds. Value is typed by kind; mirrors the AssetAttribute contract.
 */
export type AssetAttributeKind = 'integer' | 'string' | 'boolean' | 'date' | 'geo';

/** Geo attribute value (decimal degrees, WGS84). */
export interface AssetAttributeGeoValue {
  lat: number | null;
  lon: number | null;
}

/**
 * One operator-defined custom attribute in the wizard form. `value` is typed by
 * `kind`; `searchable` is captured (default true) but not shown in this iteration.
 */
export interface AssetAttributeForm {
  label: string;
  kind: AssetAttributeKind;
  value: string | number | boolean | AssetAttributeGeoValue | null;
  searchable: boolean;
}

/**
 * Asset form data structure
 */
export interface AssetFormData {
  name: string;
  assetId: string;
  enabled: boolean;
  description: string;
  assetTemplateId: string | null;
  routeGroupIds: string[];
  attributes: AssetAttributeForm[];
  protocol: string;
  latitude: number | null;
  longitude: number | null;

  /** MQTT configuration (required when protocol is MQTT) */
  mqttConfig: MqttConfig;

  /** LoRaWAN configuration (required when protocol is LoRaWAN) */
  lorawanConfig: LorawanFormConfig;

  /** Debug mode enabled for connectivity troubleshooting */
  debugEnabled: boolean;

  /** Health monitoring configuration */
  healthMonitor: HealthMonitorFormConfig;

  // Store complete objects to preserve selection state across steps
  selectedTemplate?: AssetTemplateResponse | null;
  selectedRouteGroups?: RouteGroupResponse[];
}

/**
 * Select option structure for dropdowns
 */
export interface SelectOption {
  label: string;
  value: any;
  disable?: boolean;
}

/**
 * Asset form state (external selections)
 */
export interface AssetFormState {
  selectedTemplate: AssetTemplateResponse | null;
  selectedRouteGroups: RouteGroupResponse[];
  isCreating: boolean;
  currentStep: number;
}

/**
 * Asset type card for the wizard's first step. The flat four-card choice
 * (MQTT/HTTP device, LoRaWAN Sensor, LoRaWAN Gateway) sets the protocol and, for
 * the LoRaWAN cards, the kind — branching the rest of the wizard up front. The
 * label/description are resolved in the component from i18n keyed by `value`.
 */
export interface AssetTypeOption {
  value: string;
  icon: string;
  protocol: 'HTTP' | 'MQTT' | 'LORAWAN';
  kind?: LorawanKind;
}

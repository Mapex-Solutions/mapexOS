import type { AssetTemplateResponse, RouteGroupResponse } from '@mapexos/schemas';

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
 * Asset form data structure
 */
export interface AssetFormData {
  name: string;
  assetId: string;
  enabled: boolean;
  description: string;
  assetTemplateId: string | null;
  routeGroupIds: string[];
  protocol: string;
  latitude: number | null;
  longitude: number | null;

  /** MQTT configuration (required when protocol is MQTT) */
  mqttConfig: MqttConfig;

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

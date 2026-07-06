/** TYPE IMPORTS */
import type { HealthMonitorFormConfig } from '../interfaces';

/** Heartbeat mode an asset can report liveness through. */
export type HeartbeatMode = 'implicit' | 'explicit';

/** Health-monitoring policy a protocol imposes on the form. */
export interface ProtocolHealthPolicy {
  /** Monitoring is locked on — the operator cannot disable it. */
  alwaysEnabled: boolean;
  /** Heartbeat mode the protocol forces, or null when the operator chooses. */
  forcedHeartbeatMode: HeartbeatMode | null;
}

/**
 * Resolves the health-monitoring policy for an asset protocol. The switch keeps
 * each protocol's rule in one place, so adding a protocol is a single new case.
 * @param {string} protocol - Asset protocol (case-insensitive: MQTT/HTTP/LORAWAN).
 * @param {string} [kind] - LoRaWAN kind (device | gateway), when protocol is LoRaWAN.
 * @returns {ProtocolHealthPolicy} The forced monitoring rules for the protocol.
 */
export function resolveProtocolHealthPolicy(protocol: string, kind?: string): ProtocolHealthPolicy {
  switch ((protocol ?? '').toUpperCase()) {
    case 'MQTT':
      // Liveness comes from broker presence (CONNECT/DISCONNECT): always on, explicit.
      return { alwaysEnabled: true, forcedHeartbeatMode: 'explicit' };
    case 'LORAWAN':
      // Liveness is implicit for both LoRaWAN shapes. A gateway's online/offline
      // is known from its LNS connection, so monitoring is always on (like MQTT);
      // a device infers liveness from inbound data and stays operator-toggled.
      return { alwaysEnabled: kind === 'gateway', forcedHeartbeatMode: 'implicit' };
    case 'HTTP':
    default:
      return { alwaysEnabled: false, forcedHeartbeatMode: null };
  }
}

/**
 * Returns a health-monitor config with the protocol's policy applied, leaving
 * operator-chosen fields untouched.
 * @param {HealthMonitorFormConfig} health - Current health-monitor form config.
 * @param {string} protocol - Asset protocol.
 * @param {string} [kind] - LoRaWAN kind, when protocol is LoRaWAN.
 * @returns {HealthMonitorFormConfig} The config with forced fields applied.
 */
export function applyProtocolHealthPolicy(
  health: HealthMonitorFormConfig,
  protocol: string,
  kind?: string
): HealthMonitorFormConfig {
  const policy = resolveProtocolHealthPolicy(protocol, kind);
  const heartbeatMode = policy.forcedHeartbeatMode ?? health.heartbeatMode;
  return {
    ...health,
    enabled: policy.alwaysEnabled ? true : health.enabled,
    ...(heartbeatMode !== undefined ? { heartbeatMode } : {}),
  };
}

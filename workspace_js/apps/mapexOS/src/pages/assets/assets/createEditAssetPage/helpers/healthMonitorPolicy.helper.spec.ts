import { describe, it, expect } from 'vitest';

import type { HealthMonitorFormConfig } from '../interfaces';
import { resolveProtocolHealthPolicy, applyProtocolHealthPolicy } from './healthMonitorPolicy.helper';

/** Minimal health config with monitoring off and no heartbeat mode chosen. */
function baseHealth(overrides: Partial<HealthMonitorFormConfig> = {}): HealthMonitorFormConfig {
  return {
    enabled: false,
    thresholdMinutes: 10,
    requiredMisses: 3,
    offlineRouteGroupIds: [],
    onlineRouteGroupIds: [],
    selectedOfflineRouteGroups: [],
    selectedOnlineRouteGroups: [],
    ...overrides,
  };
}

describe('resolveProtocolHealthPolicy', () => {
  const cases: Array<{
    name: string;
    protocol: string;
    kind?: string;
    alwaysEnabled: boolean;
    forcedHeartbeatMode: 'implicit' | 'explicit' | null;
  }> = [
    { name: 'MQTT is always on and explicit', protocol: 'MQTT', alwaysEnabled: true, forcedHeartbeatMode: 'explicit' },
    { name: 'MQTT is case-insensitive', protocol: 'mqtt', alwaysEnabled: true, forcedHeartbeatMode: 'explicit' },
    { name: 'LoRaWAN device is implicit, operator-toggled', protocol: 'LORAWAN', kind: 'device', alwaysEnabled: false, forcedHeartbeatMode: 'implicit' },
    { name: 'LoRaWAN gateway is always on and implicit', protocol: 'LORAWAN', kind: 'gateway', alwaysEnabled: true, forcedHeartbeatMode: 'implicit' },
    { name: 'HTTP is unconstrained', protocol: 'HTTP', alwaysEnabled: false, forcedHeartbeatMode: null },
    { name: 'unknown protocol falls through to default', protocol: 'COAP', alwaysEnabled: false, forcedHeartbeatMode: null },
  ];

  for (const c of cases) {
    it(c.name, () => {
      const policy = resolveProtocolHealthPolicy(c.protocol, c.kind);
      expect(policy.alwaysEnabled).toBe(c.alwaysEnabled);
      expect(policy.forcedHeartbeatMode).toBe(c.forcedHeartbeatMode);
    });
  }
});

describe('applyProtocolHealthPolicy', () => {
  it('forces MQTT monitoring on and heartbeat explicit even when off', () => {
    const result = applyProtocolHealthPolicy(baseHealth({ enabled: false }), 'MQTT');
    expect(result.enabled).toBe(true);
    expect(result.heartbeatMode).toBe('explicit');
  });

  it('forces a LoRaWAN device to implicit without enabling monitoring', () => {
    const result = applyProtocolHealthPolicy(baseHealth({ enabled: false }), 'LORAWAN', 'device');
    expect(result.enabled).toBe(false);
    expect(result.heartbeatMode).toBe('implicit');
  });

  it('leaves the operator choice intact for HTTP', () => {
    const result = applyProtocolHealthPolicy(baseHealth({ enabled: true, heartbeatMode: 'explicit' }), 'HTTP');
    expect(result.enabled).toBe(true);
    expect(result.heartbeatMode).toBe('explicit');
  });

  it('omits heartbeatMode when unconstrained and unset', () => {
    const result = applyProtocolHealthPolicy(baseHealth(), 'HTTP');
    expect(result.heartbeatMode).toBeUndefined();
  });
});

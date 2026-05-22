import { z } from 'zod';

/**
 * Cross-service NATS payload for the asset heartbeat stream
 * (canonical name: ${ENV}-MAPEXOS-ASSETS-HEARTBEAT;
 * subject pattern: ${env}.mapexos.asset.heartbeat.{orgId}).
 *
 * Mirror of Go: packages/contracts/services/assets/healthmonitor/types.go::HeartbeatEvent
 *
 * Published by: js-executor (TS), http_gateway (Go), mapex-mqtt-broker.
 * Consumed by: assets/healthmonitor (Go).
 */
export const HeartbeatEventSchema = z.object({
  orgId: z.string(),
  assetUUID: z.string(),
  pathKey: z.string(),
  ts: z.number().int(),
});

export type HeartbeatEvent = z.infer<typeof HeartbeatEventSchema>;

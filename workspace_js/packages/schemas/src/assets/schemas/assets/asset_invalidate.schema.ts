import { z } from 'zod';

/**
 * Cross-service fanout wire contract for the canonical subject
 * `${env}.mapexos.fanout.asset.invalidate` (stream `${ENV}-MAPEXOS-FANOUT`).
 * Mirror of Go: packages/contracts/services/assets/assets/types.go::AssetInvalidatePayload
 * Published by: assets service. Consumed by: router (Go), js-executor (TS).
 */
export const AssetInvalidatePayloadSchema = z.object({
  orgId: z.string(),
  assetUUID: z.string(),
});

export type AssetInvalidatePayload = z.infer<typeof AssetInvalidatePayloadSchema>;

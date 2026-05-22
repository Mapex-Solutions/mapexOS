import { z } from 'zod';

/**
 * Cross-service fanout wire contract for the canonical subject
 * `${env}.mapexos.fanout.template.invalidate` (stream `${ENV}-MAPEXOS-FANOUT`).
 * Mirror of Go: packages/contracts/services/assets/assettemplates/types.go::TemplateInvalidatePayload
 * Published by: assets service (assettemplates module). Consumed by: router (Go), events (Go), js-executor (TS).
 */
export const TemplateInvalidatePayloadSchema = z.object({
  orgId: z.string(),
  templateId: z.string(),
});

export type TemplateInvalidatePayload = z.infer<typeof TemplateInvalidatePayloadSchema>;

import { z } from 'zod';

/**
 * Cross-service fanout wire contract for the canonical subject
 * `${env}.mapexos.fanout.workflow.definition.invalidate` (stream `${ENV}-MAPEXOS-FANOUT`).
 * Mirror of Go: packages/contracts/services/workflow/definitions/types.go::DefinitionInvalidatePayload
 * Published by: workflow service (definitions module). Consumed by: js-workflow-executor.
 */
export const DefinitionInvalidatePayloadSchema = z.object({
  orgId: z.string(),
  definitionId: z.string(),
  nodeIds: z.array(z.string()),
});

export type DefinitionInvalidatePayload = z.infer<typeof DefinitionInvalidatePayloadSchema>;

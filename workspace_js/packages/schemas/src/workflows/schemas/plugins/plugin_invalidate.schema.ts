import { z } from 'zod';

/**
 * Cross-service fanout wire contract for the canonical subject
 * `${env}.mapexos.fanout.workflow.plugin.invalidate` (stream `${ENV}-MAPEXOS-FANOUT`).
 * Mirror of Go: packages/contracts/services/workflow/plugins/types.go::PluginInvalidatePayload
 * Published by: workflow service (plugins module). Consumed by: workflow self-fanout.
 */
export const PluginInvalidatePayloadSchema = z.object({
  pluginId: z.string(),
  action: z.enum(['create', 'update', 'delete']),
});

export type PluginInvalidatePayload = z.infer<typeof PluginInvalidatePayloadSchema>;

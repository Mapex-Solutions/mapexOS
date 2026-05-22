import { z } from 'zod';

/**
 * Body shape of POST /api/v1/heartbeat?ds={dataSourceId}.
 *
 * Mirror of Go: packages/contracts/services/http_gateway/events/heartbeat_dto.go::HeartbeatRequestDTO
 *
 * Published by: device firmware (HTTP-protocol assets with HealthMonitorConfig.heartbeatMode='explicit').
 * Consumed by: http_gateway (Go).
 *
 * orgId and pathKey are derived from the resolved DataSource server-side via
 * c.Locals — never from the request body — so a compromised body cannot
 * spoof a different tenant.
 */
export const ZodHeartbeatRequestSchema = z.object({
  assetUUID: z.string().min(1),
});

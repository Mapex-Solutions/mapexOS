import type { z } from 'zod';
import type { ZodHeartbeatRequestSchema } from '../../schemas/events/heartbeat_request.schema';

export type HeartbeatRequest = z.infer<typeof ZodHeartbeatRequestSchema>;

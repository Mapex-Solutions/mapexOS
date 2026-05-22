import { z } from 'zod';
import {
	ZodRouteGroupIdSchema,
	ZodRouteGroupCreateSchema,
	ZodRouteGroupUpdateSchema,
	ZodRouteGroupResponseSchema,
	ZodRouteGroupQuerySchema,
} from '@/router/schemas/routegroups/routegroups.schema';

// Export inferred types
export type RouteGroupId = z.infer<typeof ZodRouteGroupIdSchema>;
export type RouteGroupCreate = z.infer<typeof ZodRouteGroupCreateSchema>;
export type RouteGroupUpdate = z.infer<typeof ZodRouteGroupUpdateSchema>;
export type RouteGroupResponse = z.infer<typeof ZodRouteGroupResponseSchema>;
export type RouteGroupQuery = z.infer<typeof ZodRouteGroupQuerySchema>;

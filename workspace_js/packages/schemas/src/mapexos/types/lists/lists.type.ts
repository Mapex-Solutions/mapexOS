import { z } from 'zod';
import {
	ZodListIdSchema,
	ZodListCreateSchema,
	ZodListUpdateSchema,
	ZodListQuerySchema,
	ZodListResponseSchema,
} from '@/mapexos';

/**
 * List API types
 */
export type ListId = z.infer<typeof ZodListIdSchema>;
export type ListCreate = z.infer<typeof ZodListCreateSchema>;
export type ListUpdate = z.infer<typeof ZodListUpdateSchema>;
export type ListQuery = z.infer<typeof ZodListQuerySchema>;
export type ListResponse = z.infer<typeof ZodListResponseSchema>;

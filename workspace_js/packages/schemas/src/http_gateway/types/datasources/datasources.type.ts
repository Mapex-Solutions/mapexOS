import { z } from 'zod';
import {
	ZodDataSourceIdSchema,
	ZodDataSourceCreateSchema,
	ZodDataSourceUpdateSchema,
	ZodDataSourceResponseSchema,
	ZodDataSourceQuerySchema,
} from '@/http_gateway/schemas/datasources/datasources.schema';

// Export inferred types
export type DataSourceId = z.infer<typeof ZodDataSourceIdSchema>;
export type DataSourceCreate = z.infer<typeof ZodDataSourceCreateSchema>;
export type DataSourceUpdate = z.infer<typeof ZodDataSourceUpdateSchema>;
export type DataSourceResponse = z.infer<typeof ZodDataSourceResponseSchema>;
export type DataSourceQuery = z.infer<typeof ZodDataSourceQuerySchema>;

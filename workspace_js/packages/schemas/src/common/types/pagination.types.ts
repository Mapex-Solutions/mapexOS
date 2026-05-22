import type { z } from 'zod';
import type { ZodPaginationSchema } from '../schemas/pagination.schema';

/**
 * Standard pagination metadata type
 */
export type PaginationType = z.infer<typeof ZodPaginationSchema>;

/**
 * Generic paginated response type
 * @template T - The type of items in the response
 */
export type PaginatedResponse<T> = {
  items: T[];
  pagination: PaginationType;
};

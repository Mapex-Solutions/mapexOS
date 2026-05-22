import { z, NumberIntAndPositive, IsNumber } from '@mapexos/validations';
import type { z as zType } from 'zod';

/**
 * Standard pagination metadata schema
 * Used across all paginated API responses
 */
export const ZodPaginationSchema = z.object({
  page: NumberIntAndPositive,
  perPage: NumberIntAndPositive.max(100),
  totalItems: IsNumber.int().min(0),
  totalPages: IsNumber.int().min(0),
});

/**
 * Generic paginated response schema
 * @template T - The type of items in the response
 */
export const createPaginatedResponseSchema = <T extends zType.ZodTypeAny>(itemSchema: T) =>
  z.object({
    items: z.array(itemSchema),
    pagination: ZodPaginationSchema,
  });

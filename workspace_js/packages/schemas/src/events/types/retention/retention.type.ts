import { z } from 'zod';
import {
	ZodRetentionPolicyUpsertSchema,
	ZodRetentionPolicyResponseSchema,
	ZodRetentionPolicyQuerySchema,
	ZodRetentionPolicyParamsSchema,
	ZodRetentionPolicyPaginatedResultSchema,
} from '@/events/schemas/retention/retention.schema';

/**
 * Retention policy upsert type - PUT body for creating/updating policies
 */
export type RetentionPolicyUpsert = z.infer<typeof ZodRetentionPolicyUpsertSchema>;

/**
 * Retention policy response type - Single policy API response
 */
export type RetentionPolicyResponse = z.infer<typeof ZodRetentionPolicyResponseSchema>;

/**
 * Retention policy query type - GET query parameters for listing
 */
export type RetentionPolicyQuery = z.infer<typeof ZodRetentionPolicyQuerySchema>;

/**
 * Retention policy params type - Route parameters for individual operations
 */
export type RetentionPolicyParams = z.infer<typeof ZodRetentionPolicyParamsSchema>;

/**
 * Retention policy paginated result type - Paginated list response
 */
export type RetentionPolicyPaginatedResult = z.infer<typeof ZodRetentionPolicyPaginatedResultSchema>;

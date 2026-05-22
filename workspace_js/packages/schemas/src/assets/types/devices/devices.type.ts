import { z } from 'zod';
import { ZodRefreshTokenResponseSchema } from '@/assets';

/**
 * External API types (device JWT auth)
 */
export type RefreshTokenResponse = z.infer<typeof ZodRefreshTokenResponseSchema>;

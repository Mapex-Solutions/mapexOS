import { z } from 'zod';

/**
 * Refresh-token response schema - Returned from
 * POST /api/v1/devices/refresh_token. Devices read mqttCredential
 * and replace their stored JWT before the previous one expires.
 */
export const ZodRefreshTokenResponseSchema = z.object({
	mqttCredential: z.string(),
	mqttTokenExpiresAt: z.coerce.date(),
});

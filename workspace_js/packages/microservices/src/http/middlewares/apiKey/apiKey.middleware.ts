import type { Request, Response, NextFunction, RequestHandler } from 'express';

/**
 * API Key authentication middleware for internal service-to-service communication.
 *
 * This middleware validates the X-API-Key header against the expected API key.
 * Used for internal routes that are called by other microservices.
 *
 * @param expectedKey - The expected API key to validate against
 * @returns Express middleware function
 */
export function apiKeyAuthMiddleware(expectedKey: string): RequestHandler {
	return (req: Request, res: Response, next: NextFunction) => {
		const providedKey = req.headers['x-api-key'] as string;

		if (!providedKey || providedKey !== expectedKey) {
			return res.status(401).json({
				success: false,
				message: 'Unauthorized: Invalid API Key',
			});
		}

		return next();
	};
}

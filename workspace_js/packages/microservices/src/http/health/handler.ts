/**
 * Health Handler
 * Same structure as workspace_go/packages/microservices/http/health/handler.go
 */

import type { Request, Response, RequestHandler } from 'express';
import { HttpStatus } from '../status';
import type { HealthService } from './service';

/** Returns an Express handler that checks the health of the service. */
export function healthHandler(service: HealthService): RequestHandler {
	return async (_req: Request, res: Response) => {
		const result = await service.check();

		if (result.status === 'unhealthy') {
			return res.status(HttpStatus.SERVICE_UNAVAILABLE).json({
				httpStatus: HttpStatus.SERVICE_UNAVAILABLE,
				error: null,
				data: result,
			});
		}

		return res.status(HttpStatus.OK).json({
			httpStatus: HttpStatus.OK,
			error: null,
			data: result,
		});
	};
}

import type { Request, Response, RequestHandler } from '@mapexos/microservices';
import { success } from '@mapexos/microservices';

/**
 * Creates a request handler for performing a health check.
 *
 * @returns A RequestHandler function that sends a 200 status with a 'RUNNING' message.
 */
export function healthCheck(): RequestHandler  {
	return function (_: Request, res: Response) {
		return success(res, { status: 'RUNNING' })
	}
}
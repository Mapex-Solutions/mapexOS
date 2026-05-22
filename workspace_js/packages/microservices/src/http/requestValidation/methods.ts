import type { Request, Response, NextFunction, RequestHandler } from '@src/http';
import type { ZodSchema } from 'zod';
import type { Validation } from './types';

import { container } from 'tsyringe';
import { zodValidationError } from '@mapexos/validations';

import { badRequest } from '@src/http';
import { LOGGER_TOKEN, Logger } from '@src/logger';


/**
 * Creates a new validation object containing the provided data transfer object (DTO) schemas.
 *
 * @param bodyDTO - The Zod schema for the request body.
 * @param queryDTO - The Zod schema for the request query parameters.
 * @param paramsDTO - The Zod schema for the request route parameters.
 * @returns A new Validation object containing the provided DTO schemas.
 */
export function newValidation(bodyDTO: ZodSchema, queryDTO: ZodSchema, paramsDTO: ZodSchema): Validation {
	return {
		bodyType: bodyDTO,
		queryType: queryDTO,
		paramsType: paramsDTO,
	} as Validation;
}

/**
 * @function ValidationMiddleware
 * @description Creates a middleware function for validating and binding request data to Zod schemas.
 *
 * @param {Validation} v - The Validation object containing the Zod schemas for request body, query parameters, and route parameters.
 * @returns {RequestHandler} An Express RequestHandler function that validates and binds the request data to the corresponding DTOs.
 *
 * @example
 * ```typescript
 * import { Request, Response, NextFunction } from 'express';
 * import { z } from 'zod';
 * import { Validation, ValidationMiddleware } from './methods';
 *
 * const userSchema = z.object({
 *   name: z.string(),
 *   age: z.number(),
 * });
 *
 * const validation: Validation = NewValidation(userSchema, null, null);
 * const validationMiddleware = ValidationMiddleware(validation);
 *
 * app.use('/users', validationMiddleware, (req: Request, res: Response) => {
 *   const bodyDTO = GetDTO<typeof userSchema.output>(res, 'bodyDTO');
 *   console.log(bodyDTO); // { name: 'John Doe', age: 30 }
 * });
 * ```
 */
export function validationMiddleware(v: Validation): RequestHandler {
	const logger = container.resolve<Logger>(LOGGER_TOKEN);

	return async function (req: Request, res: Response, next: NextFunction): Promise<any> {
		try {

			if (v.bodyType) res.locals['bodyDTO'] = await validateAndBind(v.bodyType, req.body);
			if (v.queryType) res.locals['queryDTO'] = await validateAndBind(v.queryType, req.query);
			if (v.paramsType) res.locals['paramsDTO'] = await validateAndBind(v.paramsType, req.params);
			return next();

		} catch (error: any) {
			logger.error(`Validation error: ${error.message} - Request: ${req.method} ${req.path} - Body: ${JSON.stringify(req.body)}, Query: ${JSON.stringify(req.query)}, Params: ${JSON.stringify(req.params)}`);
      return badRequest(res, error.message);
		}
	};
}

/**
 * Retrieves a data transfer object (DTO) from the response's locals.
 *
 * @template T - The type of the DTO to be retrieved.
 * @param res - The Express response object.
 * @param key - The key under which the DTO is stored in the response's locals.
 * @returns The retrieved DTO.
 * @throws Will throw an error if no data is found for the given key.
 */
export function getDTO<T>(res: Response, key: string): T {
	const data = res.locals[key] as T;
	if (!data) {
		throw new Error(`No data found for key "${key}"`);
	}
	return data;
}

/**
 * Validates and binds the provided request data to a Zod schema.
 *
 * @param schema - The Zod schema to validate and bind the request data against.
 * @param requestData - The request data to be validated and bound.
 * @returns A Promise that resolves to the validated and bound data.
 * @throws Will throw an error if the request data fails validation.
 *
 * @example
 * ```typescript
 * import { z } from 'zod';
 *
 * const userSchema = z.object({
 *   name: z.string(),
 *   age: z.number(),
 * });
 *
 * const userData = { name: 'John Doe', age: 30 };
 *
 * try {
 *   const validatedData = await validateAndBind(userSchema, userData);
 *   console.log(validatedData); // { name: 'John Doe', age: 30 }
 * } catch (error) {
 *   console.error(error); // ZodError: invalid input
 * }
 * ```
 */
async function validateAndBind(schema: ZodSchema, requestData: any): Promise<any> {
	const result = await schema.safeParseAsync(requestData);
	if (result.success) return result.data;
	throw { message: zodValidationError(result) }
}
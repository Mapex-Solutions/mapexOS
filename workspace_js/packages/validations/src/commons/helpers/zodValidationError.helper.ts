import { ZodIssue } from 'zod';
import { ZodErrorResult, SchemaErrorResponse } from '@src/commons/interfaces';

/**
 * Interation about zod errors
 * @param zodResult - Result after call safeParse or safeParseAsync
 * @public
 * @async
 * @example
 *
 * ```ts
 * const validationResult = await zodSchema.safeParseAsync(data);
 * const errorList =
 * ```
 */
export const zodValidationError = (
	zodResult: ZodErrorResult,
): string[] => {
	return zodResult
		.error
		.issues
		.map((issue: ZodIssue) => {
			let { message, path } = issue;
			message = String(message).toLowerCase();
			return `${path.join('.')} ${message}`;
		}) as string[];
};
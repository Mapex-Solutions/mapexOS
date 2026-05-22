import { z, RefinementCtx } from 'zod';

/**
 * Adds a custom issue to the Zod validation context.
 *
 * @param {RefinementCtx} ctx - The Zod refinement context where the issue will be added.
 * @param {string} message - The message describing the issue to be added.
 * @param {(string | number)[]} path - Field to the field where the issue occurred.
 *
 * @returns void - This function does not return a value.
 */
export const addIssue = (ctx: RefinementCtx, message: string, path?: (string | number)[]) => {
	ctx.addIssue({
		code: z.ZodIssueCode.custom,
		message: message,
		path: path,
		input: undefined
	});
};
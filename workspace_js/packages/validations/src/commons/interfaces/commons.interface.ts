import { ZodError } from 'zod';

/**
 * Interface with zod errors
 * @param success - Result validation true or false
 * @param error - List with issues
 */
export interface ZodErrorResult {
  success: boolean;
  error: ZodError;
}


/**
 * Use in all validation with schemas
 * @param message - Error message
 */
export type SchemaErrorResponse = string | string[];
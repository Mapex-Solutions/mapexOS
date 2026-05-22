import type { ZodSchema } from 'zod'

/**
 * Data transfer object (DTO) validation configuration.
 */
export interface Validation {
	bodyType?: ZodSchema;
	queryType?: ZodSchema;
	paramsType?: ZodSchema;
}


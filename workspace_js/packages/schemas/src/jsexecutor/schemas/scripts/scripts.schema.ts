import { z, StringAndNotBeEmpty, StringAndBeEmptyOrOptional, ObjectAndNotBeEmpty, IsBoolean } from '@mapexos/validations';

/**
 * Payload to test a script in the JSExecutor.
 */
export const ZodScriptTestSchema = z.object({
	decode: StringAndBeEmptyOrOptional,
	validation: StringAndBeEmptyOrOptional,
	transform: StringAndNotBeEmpty,
	event: ObjectAndNotBeEmpty,

	/** Enable debug logging to Events MS (default: false) */
	debugEnabled: IsBoolean.optional().default(false),
});
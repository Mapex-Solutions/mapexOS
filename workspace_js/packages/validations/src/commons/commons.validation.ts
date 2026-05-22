import { z as _z } from 'zod';
import { validateRange, stringDateFormatTransform } from './transforms/index';

import {
	IS_REQUIRED,
	MUST_NOT_BE_EMPTY,
	INTEGER_NUMBER,
	POSITIVE_NUMBER,
} from './constants/commons.constant';

export const z = _z;

// ------------------- //
// GENERIC VALIDATIONS //
// ------------------- //

export const IsRecord = z.record(z.string(), z.any(), {
	message: IS_REQUIRED,
});

/**
 * Must be string
 */
export const IsString = z.string({
	message: IS_REQUIRED,
});

/**
 * Must be boolean
 */
export const IsBoolean = z.boolean({
	message: IS_REQUIRED,
});

/**
 * Must be a number
 */
export const IsNumber = z
	.number({
		message: IS_REQUIRED,
	});

/**
 * Field must be string and not be empty
 */
export const StringAndNotBeEmpty = IsString
	.min(1, MUST_NOT_BE_EMPTY);

/**
 * Field must be string and not be empty or optional into payload
 */
export const StringAndNotBeEmptyOrOptional = IsString
	.min(1, MUST_NOT_BE_EMPTY).optional();

/**
 * String must be string and be empty or optional into payload
 */
export const StringAndBeEmptyOrOptional = IsString.optional().or(z.literal(""));

/**
 * Field must be number and not be empty
 */
export const NumberIntAndPositive = IsNumber
	.int(INTEGER_NUMBER)
	.positive(POSITIVE_NUMBER);

/**
 * Transforms a string into a Date object if it matches a valid date format.
 * Supports Unix timestamps in seconds and milliseconds, as well as custom string formats.
 *
 * @param {string} dataString - The string to be transformed into a Date object.
 * @param ctx - The Zod transformation context used to add issues if validation fails.
 * @returns A Date object if the string is a valid date format, otherwise it adds an issue to the context and returns `z.NEVER`.
 */
export const IsStringDateFormat = z
	.string()
	.transform(stringDateFormatTransform);

/**
 * Field must be an object and not be empty
 */
export const ObjectAndNotBeEmpty = z.record(z.string(), z.any())
	.refine((obj) => Object.keys(obj).length > 0, {
		message: "Object cannot be empty",
	});


/**
 * Creates a Zod schema that validates a value as either a positive integer or a non-empty string,
 * and optionally checks if the value falls within a specified numeric range.
 *
 * @param min - The optional minimum value for the numeric range validation. If not provided, no minimum check is applied.
 * @param max - The optional maximum value for the numeric range validation. If not provided, no maximum check is applied.
 * @returns A Zod schema that validates the input as a positive integer or a non-empty string,
 *          and applies the optional range validation if `min` or `max` are specified.
 */
export function NumberPositiveMinAndMax(min = Number.MIN_SAFE_INTEGER, max = Number.MAX_SAFE_INTEGER) {
	return NumberIntAndPositive
		.transform(validateRange(min, max));
}

/**
 * Creates a Zod schema that validates a value as either a positive integer or a non-empty string,
 * and optionally checks if the value falls within a specified numeric range.
 *
 * @param min - The optional minimum value for the numeric range validation. If not provided, no minimum check is applied.
 * @param max - The optional maximum value for the numeric range validation. If not provided, no maximum check is applied.
 * @returns A Zod schema that validates the input as a positive integer or a non-empty string,
 *          and applies the optional range validation if `min` or `max` are specified.
 */
export function NumberOrStringPositive(min = Number.MIN_SAFE_INTEGER, max = Number.MAX_SAFE_INTEGER) {
	return NumberIntAndPositive
		.or(StringAndNotBeEmpty)
		.transform(validateRange(min, max));
}

/**
 * MongoDB ObjectID validation (24 hex characters)
 */
export const IsMongoId = z.string().regex(/^[0-9a-fA-F]{24}$/, 'Invalid MongoDB ObjectID');

/**
 * URL validation
 */
export const IsUrl = z.string().url('Invalid URL format');

/**
 * Duration regex matching the Go-side parser in
 * services/assets/src/modules/assets/application/services/asset_handler_crud.go::parseTTL.
 * Accepts an integer count followed by a single unit suffix:
 *   s = seconds, m = minutes, h = hours, d = days, y = years.
 */
export const DURATION_REGEX = /^\d+(s|m|h|d|y)$/;

/**
 * Returns true when the input matches the platform duration format.
 * Used by Quasar form rules and ad-hoc form-layer checks.
 */
export function isValidDuration(value: string): boolean {
	return DURATION_REGEX.test(value);
}

/**
 * Quasar-style validation rule: returns `true` when the value is valid,
 * or an error string when it is not. Matches the rule signature Quasar's
 * QInput component expects in the `:rules` prop.
 */
export function validateDuration(value: string): string | true {
	return isValidDuration(value) ? true : 'Invalid duration format. Use 30d, 90d, 1y, etc.';
}

/**
 * Zod schema for duration strings. Pairs naturally with the Go contract's
 * `mqtt_token_ttl` field which is parsed by the same regex on the backend.
 */
export const IsDuration = z
	.string()
	.regex(DURATION_REGEX, 'Invalid duration format. Use 30d, 90d, 1y, etc.');
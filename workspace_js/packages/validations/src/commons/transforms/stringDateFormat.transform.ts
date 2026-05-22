import { z, RefinementCtx } from 'zod';

import dayjs from 'dayjs';
import customParseFormat from 'dayjs/plugin/customParseFormat';

dayjs.extend(customParseFormat);

import { addIssue } from '@src/commons/helpers';
import { STRING_DATE_FORMATS } from '@src/commons/constants';

/**
 * Transforms a string representation of a date into a JavaScript Date object.
 * The function supports Unix timestamps in seconds and milliseconds, as well as
 * custom date formats defined in `STRING_DATE_FORMATS`.
 *
 * @param dataString - The string representation of the date to be transformed.
 *                     It can be a Unix timestamp (in seconds or milliseconds) or
 *                     a date string in one of the allowed formats.
 * @param ctx - The Zod refinement context used to add validation issues if the
 *              date string is invalid.
 * @returns A JavaScript Date object if the transformation is successful, or
 *          `z.NEVER` if the date string is invalid.
 */
export function stringDateFormatTransform(dataString: string, ctx: RefinementCtx) {

	// Check if date is milliseconds (new Date().valueOf)
	const isNumeric = /^\d+$/.test(dataString);

	if (isNumeric) {
		if (dataString.length === 10) {
			// Timestamps Unix on seconds
			const unixSeconds = parseInt(dataString, 10);
			const date = dayjs.unix(unixSeconds);
			if (date.isValid()) {
				return dataString;
			}
		} else if (dataString.length === 13) {
			// Timestamps Unix on milliseconds
			const unixMillis = parseInt(dataString, 10);
			const date = dayjs(unixMillis);
			if (date.isValid()) {
				return dataString;
			}
		}
	}

	// Others formats
	for (const format of STRING_DATE_FORMATS) {
		const date = dayjs(dataString, format, true); // `true` para parsing estrito
		if (date.isValid()) {
			return dataString;
		}
	}

	// Fail format not a available
	addIssue(ctx, `Invalid date. Allowed formats: ${[
		...STRING_DATE_FORMATS,
		'Unix timestamp on seconds (10 digits)',
		'Unix timestamp on milliseconds (13 digits)',
	].join(', ')}`);

	// Fail validation
	return z.NEVER;
}
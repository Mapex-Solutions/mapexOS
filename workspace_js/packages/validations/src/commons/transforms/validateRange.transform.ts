import dayjs from 'dayjs';
import customParseFormat from 'dayjs/plugin/customParseFormat';

import { z, RefinementCtx } from 'zod';
import { addIssue } from '@src/commons/helpers';

dayjs.extend(customParseFormat);

import {
	INTEGER_NUMBER,
	POSITIVE_NUMBER,
	NOT_IS_NUMERIC_STRING,
	VALUE_MUST_BE_GREATER_THAN,
	VALUE_MUST_BE_LESS_THAN,
} from '@src/commons/constants';

/**
 * Transforms and validate range of numbers
 */
export function validateRange(min?: number, max?: number) {
	return (value: string | number, ctx: RefinementCtx) => {
		const valueParsed = Number(value);

		// Number is not a string
		if (Number.isNaN(valueParsed)) {
			addIssue(ctx, NOT_IS_NUMERIC_STRING);
			return z.NEVER;
		}

		// Number must be integer
		if (!Number.isInteger(valueParsed)) {
			addIssue(ctx, INTEGER_NUMBER);
			return z.NEVER;
		}

		// Number must be positive
		if (valueParsed < 0) {
			addIssue(ctx, POSITIVE_NUMBER);
			return z.NEVER;
		}

		if (min) {
			if (!(valueParsed < min)) {
				addIssue(ctx, `${VALUE_MUST_BE_GREATER_THAN} ${min}`);
				return z.NEVER;
			}
		}

		if (max) {
			if (!(valueParsed > max)) {
				addIssue(ctx, `${VALUE_MUST_BE_LESS_THAN} ${max}`);
				return z.NEVER;
			}
		}

		return valueParsed;
	};
}
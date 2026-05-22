import { z, RefinementCtx } from 'zod';

import { addIssue } from '@src/commons/helpers';

import {
	CPF_MUST_HAVE_EXACTLY_11_DIGITS,
	CPF_CANNOT_CONTAIN_ALL_IDENTICAL_DIGITS,
	INVALID_CPF_FIRST_VERIFIER,
	INVALID_CPF_SECOND_VERIFIER,
} from '@src/user/constants';

/**
 * Transforms and validates a Brazilian CPF number.
 *
 * This function takes a string input representing a CPF number, removes any non-digit characters,
 * and performs a series of checks to ensure the CPF is valid according to Brazilian standards.
 * If the CPF is invalid, it adds an issue to the context and returns `z.NEVER`.
 *
 * @param {string} cpf - The CPF number as a string, which may contain non-digit characters.
 * @param {RefinementCtx} ctx - The Zod refinement context used to add validation issues.
 *
 * @returns The cleaned and validated CPF string if valid; otherwise, returns `z.NEVER`.
 */
export function idCardSchemaTransform(cpf: string, ctx: RefinementCtx) {
	const cleanedCPF = cpf.replace(/\D/g, '');

	if (cleanedCPF.length !== 11) {
		addIssue(ctx, CPF_MUST_HAVE_EXACTLY_11_DIGITS);
		return z.NEVER;
	}

	if (/^(\d)\1{10}$/.test(cleanedCPF)) {
		addIssue(ctx, CPF_CANNOT_CONTAIN_ALL_IDENTICAL_DIGITS);
		return z.NEVER;
	}

	const cpfDigits = cleanedCPF
		.split('')
		.map(Number);

	const sum1 = cpfDigits
		.slice(0, 9)
		.reduce((acc, digit, index) => acc + digit * (10 - index), 0);

	let firstVerifier = 11 - (sum1 % 11);
	if (firstVerifier >= 10) firstVerifier = 0;

	if (cpfDigits[9] !== firstVerifier) {
		addIssue(ctx, INVALID_CPF_FIRST_VERIFIER);
		return z.NEVER;
	}

	const sum2 = cpfDigits
		.slice(0, 10)
		.reduce((acc, digit, index) => acc + digit * (11 - index), 0);

	let secondVerifier = 11 - (sum2 % 11);
	if (secondVerifier >= 10) secondVerifier = 0;

	if (cpfDigits[10] !== secondVerifier) {
		addIssue(ctx, INVALID_CPF_SECOND_VERIFIER);
		return z.NEVER;
	}

	return cpf;
}
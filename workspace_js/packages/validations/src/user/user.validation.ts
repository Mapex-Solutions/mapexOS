import { z } from 'zod';
import { StringAndNotBeEmpty } from '@src/commons';

import { idCardSchemaTransform } from './transforms';

/**
 * Check if the string is a valid email address
 */
export const Email = StringAndNotBeEmpty
	.email();

/**
 * Schema to validate a password
 * The password must be at least 8 characters long, contain at least one letter, one number, and one special character.
 */
export const Password = z
	.string()
	.min(8, 'Password must have at least 8 characters.')
	.refine((password) => /[a-zA-Z]/.test(password), { message: 'Password must contain at least one letter.' })
	.refine((password) => /\d/.test(password), { message: 'Password must contain at least one number.' })
	.refine((password) => /[!@#$%^&*(),.?":{}|<>]/.test(password), { message: 'Password must contain at least one special character.' });

/**
 * Schema to validate a Brazilian CPF number
 */
export const IdCard = z
	.string()
	.transform(idCardSchemaTransform)
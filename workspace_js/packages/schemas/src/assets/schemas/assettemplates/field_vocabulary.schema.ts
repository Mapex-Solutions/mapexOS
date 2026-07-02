import { z } from 'zod';

/**
 * Field-vocabulary query schema for the canonical-field suggestion endpoint.
 *
 * The `lang` selects which language the backend resolves the human-readable
 * `hint` and group `label` into. The canonical `value` is never translated.
 */
export const ZodFieldVocabularyQuerySchema = z.object({
	lang: z.string(),
});

/**
 * A single suggested canonical field within a vocabulary group.
 *
 * - `value` is the canonical, English-only identifier saved as the field name.
 * - `hint` is already resolved to the requested language by the backend.
 * - `type` maps to a dynamic-field type and pre-selects the field type on pick.
 * - `unit` may be an empty string when the field is dimensionless.
 */
export const ZodFieldVocabularyFieldSchema = z.object({
	value: z.string(),
	hint: z.string(),
	type: z.enum(['number', 'string', 'bool', 'date', 'geo']),
	unit: z.string(),
});

/**
 * A vocabulary group: a category header plus its canonical fields.
 *
 * `label` is already resolved to the requested language by the backend.
 * Groups arrive in a fixed, meaningful order and are rendered as received.
 */
export const ZodFieldVocabularyGroupSchema = z.object({
	category: z.string(),
	label: z.string(),
	fields: z.array(ZodFieldVocabularyFieldSchema),
});

/**
 * Field-vocabulary response payload (the envelope `data` object).
 */
export const ZodFieldVocabularyResponseSchema = z.object({
	groups: z.array(ZodFieldVocabularyGroupSchema),
});

import { z } from 'zod';
import {
	ZodFieldVocabularyQuerySchema,
	ZodFieldVocabularyFieldSchema,
	ZodFieldVocabularyGroupSchema,
	ZodFieldVocabularyResponseSchema,
} from '../../schemas/assettemplates/field_vocabulary.schema';

/**
 * TypeScript types inferred from the field-vocabulary Zod schemas.
 */

export type FieldVocabularyQuery = z.infer<typeof ZodFieldVocabularyQuerySchema>;
export type FieldVocabularyField = z.infer<typeof ZodFieldVocabularyFieldSchema>;
export type FieldVocabularyGroup = z.infer<typeof ZodFieldVocabularyGroupSchema>;
export type FieldVocabularyResponse = z.infer<typeof ZodFieldVocabularyResponseSchema>;

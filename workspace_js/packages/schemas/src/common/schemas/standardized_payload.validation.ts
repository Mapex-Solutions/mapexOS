import { z } from 'zod';
import { StringAndNotBeEmpty, ObjectAndNotBeEmpty } from '@mapexos/validations';

/**
 * Schema for standardized payloads
 */
export const ZodStandardizedPayloadSchema = z.object({
	eventType: StringAndNotBeEmpty,
	eventId: StringAndNotBeEmpty,
	data: ObjectAndNotBeEmpty,
	metadata: ObjectAndNotBeEmpty.optional(),
	created: StringAndNotBeEmpty,
});

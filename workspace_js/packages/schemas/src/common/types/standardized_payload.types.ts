import { z } from 'zod';
import { ZodStandardizedPayloadSchema } from '@/common';

export type StandardizedPayload = z.infer<typeof ZodStandardizedPayloadSchema>

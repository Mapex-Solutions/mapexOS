import type { z } from 'zod';
import type { ZodAuthProjectionSchema } from '@/assets/schemas/auth/auth.schema';

export type AuthProjection = z.infer<typeof ZodAuthProjectionSchema>;

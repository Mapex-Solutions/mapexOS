import { z } from 'zod';
import { ZodOrgContextSchema } from '@/common';

export type OrgContext = z.infer<typeof ZodOrgContextSchema>

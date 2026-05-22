import { z } from 'zod';
import { ZodScriptTestSchema } from '@/jsexecutor';

export type ScriptTest = z.infer<typeof ZodScriptTestSchema>
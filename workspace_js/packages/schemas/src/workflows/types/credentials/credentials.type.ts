import { z } from 'zod';
import {
	ZodCredentialIdSchema,
	ZodCredentialPluginIdSchema,
	ZodCredentialQuerySchema,
	ZodCredentialCreateSchema,
	ZodCredentialUpdateSchema,
	ZodCredentialResponseSchema,
	ZodCredentialTestResultSchema,
	ZodCredentialSchemaItemSchema,
	ZodCredentialSchemaResponseSchema,
	ZodLoadOptionsParamsSchema,
	ZodLoadOptionsBodySchema,
	ZodLoadOptionsItemSchema,
} from '@/workflows/schemas/credentials/credentials.schema';

// DTO types
export type CredentialId = z.infer<typeof ZodCredentialIdSchema>;
export type CredentialPluginId = z.infer<typeof ZodCredentialPluginIdSchema>;
export type CredentialQuery = z.infer<typeof ZodCredentialQuerySchema>;
export type CredentialCreate = z.infer<typeof ZodCredentialCreateSchema>;
export type CredentialUpdate = z.infer<typeof ZodCredentialUpdateSchema>;
export type CredentialResponse = z.infer<typeof ZodCredentialResponseSchema>;
export type CredentialTestResult = z.infer<typeof ZodCredentialTestResultSchema>;
export type CredentialSchemaItem = z.infer<typeof ZodCredentialSchemaItemSchema>;
export type CredentialSchemaResponse = z.infer<typeof ZodCredentialSchemaResponseSchema>;

// LoadOptions proxy types
export type LoadOptionsParams = z.infer<typeof ZodLoadOptionsParamsSchema>;
export type LoadOptionsBody = z.infer<typeof ZodLoadOptionsBodySchema>;
export type LoadOptionsItem = z.infer<typeof ZodLoadOptionsItemSchema>;

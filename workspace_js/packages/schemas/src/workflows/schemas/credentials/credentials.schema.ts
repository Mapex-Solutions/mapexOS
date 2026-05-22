import { z, StringAndNotBeEmpty, StringAndNotBeEmptyOrOptional, IsString, IsBoolean, IsMongoId, NumberIntAndPositive } from '@mapexos/validations';
import { ZodCredentialFieldDefSchema, ZodCredentialTestDefSchema } from '@/workflows/schemas/plugins/plugins.schema';

/**
 * ============================================================================
 * CREDENTIAL INSTANCE DTOs
 * ============================================================================
 */

/**
 * Credential ID parameter schema (for URL params — MongoDB ObjectID)
 */
export const ZodCredentialIdSchema = z.object({
	id: IsMongoId,
});

/**
 * Credential Plugin ID parameter schema (for schema endpoint — pluginId string)
 */
export const ZodCredentialPluginIdSchema = z.object({
	pluginId: StringAndNotBeEmpty,
});

/**
 * Credential Query schema - Used for filtering and pagination
 */
export const ZodCredentialQuerySchema = z.object({
	pluginId: IsString.optional(),
	credentialType: IsString.optional(),
	page: NumberIntAndPositive.optional(),
	perPage: NumberIntAndPositive.max(100).optional(),
});

/**
 * Credential Create schema - Used for creating new credential instances
 */
export const ZodCredentialCreateSchema = z.object({
	name: IsString.min(1).max(255),
	pluginId: StringAndNotBeEmpty,
	credentialType: StringAndNotBeEmpty,
	data: z.record(IsString, z.any()),
});

/**
 * Credential Update schema - Used for updating existing credential instances
 */
export const ZodCredentialUpdateSchema = z.object({
	name: IsString.min(1).max(255).optional(),
	data: z.record(IsString, z.any()).optional(),
});

/**
 * Credential Response schema - What the API returns (NEVER contains secret values)
 */
export const ZodCredentialResponseSchema = z.object({
	id: IsString,
	name: StringAndNotBeEmpty,
	pluginId: StringAndNotBeEmpty,
	credentialType: StringAndNotBeEmpty,
	created: StringAndNotBeEmptyOrOptional,
	updated: StringAndNotBeEmptyOrOptional,
});

/**
 * Credential Test Result schema - Response from POST /api/v1/credentials/:id/test
 */
export const ZodCredentialTestResultSchema = z.object({
	success: z.boolean(),
});

/**
 * ============================================================================
 * LOAD OPTIONS PROXY DTOs
 * ============================================================================
 */

/**
 * LoadOptions path params schema — POST /api/v1/credentials/:id/load_options/:resourceKey
 */
export const ZodLoadOptionsParamsSchema = z.object({
	id: IsMongoId,
	resourceKey: StringAndNotBeEmpty,
});

/**
 * LoadOptions request body schema (optional dependsOn for cascading params)
 */
export const ZodLoadOptionsBodySchema = z.object({
	dependsOn: z.record(IsString, IsString).optional(),
});

/**
 * LoadOptions item schema — single option returned by the proxy
 */
export const ZodLoadOptionsItemSchema = z.object({
	label: IsString,
	value: z.any(),
});

/**
 * ============================================================================
 * CREDENTIAL SCHEMA DTOs (from plugin manifest — UI form rendering)
 * ============================================================================
 */

/**
 * Credential Schema Item - A single credential type within a plugin
 */
export const ZodCredentialSchemaItemSchema = z.object({
	id: StringAndNotBeEmpty,
	name: StringAndNotBeEmpty,
	fields: z.array(ZodCredentialFieldDefSchema),
	test: ZodCredentialTestDefSchema.optional(),
});

/**
 * Credential Schema Response - Returned by GET /api/v1/credentials/schema/:pluginId
 * Returns all available authentication methods for the plugin.
 */
export const ZodCredentialSchemaResponseSchema = z.object({
	credentials: z.array(ZodCredentialSchemaItemSchema),
});

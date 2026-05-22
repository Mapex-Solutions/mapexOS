import { z } from 'zod';
import {
	ZodHandleDefSchema,
	ZodAvailableOutputSchema,
	ZodHttpActionDefSchema,
	ZodActionOutputDefSchema,
	ZodActionDefSchema,
	ZodFetchOptionsPaginationSchema,
	ZodFetchOptionsSearchSchema,
	ZodFetchOptionsDefSchema,
	ZodPropertyRenderingSchema,
	ZodFetchOptionsRuleSchema,
	ZodPropertyFetchOptionsSchema,
	ZodPropertyOptionSchema,
	ZodDisplayOptionsSchema,
	ZodNodePropertyDefSchema,
	ZodNodeHooksSchema,
	ZodNodeTypeManifestSchema,
	ZodCredentialFieldDefSchema,
	ZodCredentialDefSchema,
	ZodPluginMetadataSchema,
	ZodPluginDefaultsSchema,
	ZodPluginIdSchema,
	ZodPluginQuerySchema,
	ZodPluginCreateSchema,
	ZodPluginUpdateSchema,
	ZodPluginResponseSchema,
} from '@/workflows/schemas/plugins/plugins.schema';

// Building block types
export type HandleDef = z.infer<typeof ZodHandleDefSchema>;
export type AvailableOutput = z.infer<typeof ZodAvailableOutputSchema>;

// Action contract types
export type HttpActionDef = z.infer<typeof ZodHttpActionDefSchema>;
export type ActionOutputDef = z.infer<typeof ZodActionOutputDefSchema>;
export type ActionDef = z.infer<typeof ZodActionDefSchema>;

// Fetch options types
export type FetchOptionsPagination = z.infer<typeof ZodFetchOptionsPaginationSchema>;
export type FetchOptionsSearch = z.infer<typeof ZodFetchOptionsSearchSchema>;
export type FetchOptionsDef = z.infer<typeof ZodFetchOptionsDefSchema>;

// Property types
export type PropertyRendering = z.infer<typeof ZodPropertyRenderingSchema>;
export type FetchOptionsRule = z.infer<typeof ZodFetchOptionsRuleSchema>;
export type PropertyFetchOptions = z.infer<typeof ZodPropertyFetchOptionsSchema>;
export type PropertyOption = z.infer<typeof ZodPropertyOptionSchema>;
export type DisplayOptions = z.infer<typeof ZodDisplayOptionsSchema>;
export type NodePropertyDef = z.infer<typeof ZodNodePropertyDefSchema>;

// Node hooks
export type NodeHooks = z.infer<typeof ZodNodeHooksSchema>;

// Node type manifest
export type NodeTypeManifest = z.infer<typeof ZodNodeTypeManifestSchema>;

// Credential types
export type CredentialFieldDef = z.infer<typeof ZodCredentialFieldDefSchema>;
export type CredentialDef = z.infer<typeof ZodCredentialDefSchema>;

// Plugin types
export type PluginMetadata = z.infer<typeof ZodPluginMetadataSchema>;
export type PluginDefaults = z.infer<typeof ZodPluginDefaultsSchema>;

// DTO types
export type PluginId = z.infer<typeof ZodPluginIdSchema>;
export type PluginQuery = z.infer<typeof ZodPluginQuerySchema>;
export type PluginCreate = z.infer<typeof ZodPluginCreateSchema>;
export type PluginUpdate = z.infer<typeof ZodPluginUpdateSchema>;
export type PluginResponse = z.infer<typeof ZodPluginResponseSchema>;

// Legacy compat aliases
/** @deprecated Use AvailableOutput */
export type OutputHint = AvailableOutput;
/** @deprecated Use FetchOptionsDef */
export type LoadOptionsDef = FetchOptionsDef;
/** @deprecated Use ActionDef */
export type OperationDef = ActionDef;

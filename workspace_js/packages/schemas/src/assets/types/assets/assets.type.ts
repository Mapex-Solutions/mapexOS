import { z } from 'zod';
import {
	ZodAssetIdSchema,
	ZodAssetUUIDSchema,
	ZodAssetCreateSchema,
	ZodAssetUpdateSchema,
	ZodAssetQuerySchema,
	ZodAssetResponseSchema,
	ZodAssetInternalIdSchema,
	ZodAssetInternalUpdateSchema,
	ZodAssetUUIDParamSchema,
	ZodAssetScriptsResponseSchema,
} from '@/assets';

/**
 * External API types (JWT auth)
 */
export type AssetId = z.infer<typeof ZodAssetIdSchema>;
export type AssetUUID = z.infer<typeof ZodAssetUUIDSchema>;
export type AssetCreate = z.infer<typeof ZodAssetCreateSchema>;
export type AssetUpdate = z.infer<typeof ZodAssetUpdateSchema>;
export type AssetQuery = z.infer<typeof ZodAssetQuerySchema>;
export type AssetResponse = z.infer<typeof ZodAssetResponseSchema>;

/**
 * Internal API types (API Key auth - MS-to-MS)
 */
export type AssetInternalId = z.infer<typeof ZodAssetInternalIdSchema>;
export type AssetInternalUpdate = z.infer<typeof ZodAssetInternalUpdateSchema>;
export type AssetUUIDParam = z.infer<typeof ZodAssetUUIDParamSchema>;
export type AssetScriptsResponse = z.infer<typeof ZodAssetScriptsResponseSchema>;

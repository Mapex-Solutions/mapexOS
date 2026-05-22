import { z } from 'zod';
import {
	ZodAssetTemplateIdSchema,
	ZodAssetTemplateCreateSchema,
	ZodAssetTemplateUpdateSchema,
	ZodAssetTemplateQuerySchema,
	ZodAssetTemplateResponseSchema,
} from '../../schemas/assettemplates';

/**
 * TypeScript types inferred from Zod schemas for Asset Templates
 */

export type AssetTemplateId = z.infer<typeof ZodAssetTemplateIdSchema>;
export type AssetTemplateCreate = z.infer<typeof ZodAssetTemplateCreateSchema>;
export type AssetTemplateUpdate = z.infer<typeof ZodAssetTemplateUpdateSchema>;
export type AssetTemplateQuery = z.infer<typeof ZodAssetTemplateQuerySchema>;
export type AssetTemplateResponse = z.infer<typeof ZodAssetTemplateResponseSchema>;
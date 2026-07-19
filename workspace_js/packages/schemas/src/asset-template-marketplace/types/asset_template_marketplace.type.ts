import { z } from 'zod';
import {
  ZodAssetTemplateCatalogItemSchema,
  ZodAssetTemplateCatalogQuerySchema,
  ZodAssetTemplateCatalogDataSchema,
  ZodAssetTemplateCatalogListResponseSchema,
  ZodAssetTemplateFacetsQuerySchema,
  ZodAssetTemplateFacetOptionSchema,
  ZodAssetTemplateFacetsSchema,
  ZodAssetTemplateFacetsResponseSchema,
  ZodAssetTemplateBundleDynamicFieldSchema,
  ZodAssetTemplateBundleSchema,
  ZodAssetTemplateBundleResponseSchema,
  ZodMarketplaceInstalledCheckRequestSchema,
  ZodMarketplaceInstalledCheckResponseSchema,
} from '../schemas/asset_template_marketplace.schema';

export type AssetTemplateCatalogItem = z.infer<typeof ZodAssetTemplateCatalogItemSchema>;
export type AssetTemplateCatalogQuery = z.infer<typeof ZodAssetTemplateCatalogQuerySchema>;
export type AssetTemplateCatalogData = z.infer<typeof ZodAssetTemplateCatalogDataSchema>;
export type AssetTemplateCatalogListResponse = z.infer<typeof ZodAssetTemplateCatalogListResponseSchema>;
export type AssetTemplateFacetsQuery = z.infer<typeof ZodAssetTemplateFacetsQuerySchema>;
export type AssetTemplateFacetOption = z.infer<typeof ZodAssetTemplateFacetOptionSchema>;
export type AssetTemplateFacets = z.infer<typeof ZodAssetTemplateFacetsSchema>;
export type AssetTemplateFacetsResponse = z.infer<typeof ZodAssetTemplateFacetsResponseSchema>;
export type AssetTemplateBundleDynamicField = z.infer<typeof ZodAssetTemplateBundleDynamicFieldSchema>;
export type AssetTemplateBundle = z.infer<typeof ZodAssetTemplateBundleSchema>;
export type AssetTemplateBundleResponse = z.infer<typeof ZodAssetTemplateBundleResponseSchema>;
export type MarketplaceInstalledCheckRequest = z.infer<typeof ZodMarketplaceInstalledCheckRequestSchema>;
export type MarketplaceInstalledCheckResponse = z.infer<typeof ZodMarketplaceInstalledCheckResponseSchema>;

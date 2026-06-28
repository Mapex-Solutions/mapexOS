import { z } from 'zod';
import {
  ZodPluginCatalogItemSchema,
  ZodPluginCatalogQuerySchema,
  ZodPluginCatalogDataSchema,
  ZodPluginCatalogListResponseSchema,
  ZodPluginFacetOptionSchema,
  ZodPluginFacetsSchema,
  ZodPluginFacetsResponseSchema,
  ZodPluginManifestResponseSchema,
} from '../schemas/plugin_marketplace.schema';

export type PluginCatalogItem = z.infer<typeof ZodPluginCatalogItemSchema>;
export type PluginCatalogQuery = z.infer<typeof ZodPluginCatalogQuerySchema>;
export type PluginCatalogData = z.infer<typeof ZodPluginCatalogDataSchema>;
export type PluginCatalogListResponse = z.infer<typeof ZodPluginCatalogListResponseSchema>;
export type PluginFacetOption = z.infer<typeof ZodPluginFacetOptionSchema>;
export type PluginFacets = z.infer<typeof ZodPluginFacetsSchema>;
export type PluginFacetsResponse = z.infer<typeof ZodPluginFacetsResponseSchema>;
export type PluginManifestResponse = z.infer<typeof ZodPluginManifestResponseSchema>;

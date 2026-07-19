import { z, IsString, IsBoolean, IsNumber, NumberIntAndPositive, StringAndNotBeEmpty } from '@mapexos/validations';

/**
 * AssetTemplateCatalogItem — a single card in the asset template marketplace
 * listing from `GET /api/v1/asset_templates`. The marketplace service has
 * already resolved `name`/`description` to plain strings via the `?lang` query
 * (they are NOT localized {en-US,pt-BR} objects here). The heavy bundle
 * (`asset_template_information.json`) is fetched separately via the detail
 * endpoint.
 */
export const ZodAssetTemplateCatalogItemSchema = z.object({
  id: StringAndNotBeEmpty,
  marketplaceGuid: IsString,
  slug: StringAndNotBeEmpty,
  name: IsString,
  description: IsString.optional().default(''),
  category: IsString.optional().default(''),
  vendor: StringAndNotBeEmpty,
  vendorName: IsString.optional().default(''),
  model: IsString.optional().default(''),
  version: IsString.optional().default(''),
  icon: IsString.optional().default(''),
  image: IsString.optional().default(''),
  hasImage: IsBoolean.optional().default(false),
  fieldCount: IsNumber.int().optional().default(0),
  hasScripts: IsBoolean.optional().default(false),
  sha256: IsString.optional().default(''),
});

/**
 * AssetTemplateCatalogQuery — query params for the marketplace list endpoint.
 * Filtering and pagination are server-side. `lang` selects which localized
 * strings the service resolves into the flat item fields.
 */
export const ZodAssetTemplateCatalogQuerySchema = z.object({
  search: IsString.optional(),
  category: IsString.optional(),
  vendor: IsString.optional(),
  model: IsString.optional(),
  version: IsString.optional(),
  lang: IsString.optional(),
  page: NumberIntAndPositive.optional(),
  perPage: NumberIntAndPositive.max(100).optional(),
});

/**
 * AssetTemplateCatalogData — the inner data body of the list envelope.
 * `total` is the full server-side count for pagination.
 */
export const ZodAssetTemplateCatalogDataSchema = z.object({
  items: z.array(ZodAssetTemplateCatalogItemSchema),
  total: IsNumber.int().min(0),
  page: IsNumber.int().min(0),
  perPage: IsNumber.int().min(0),
});

/**
 * AssetTemplateCatalogListResponse — the full `{status,errors,data}` envelope
 * returned by `GET /api/v1/asset_templates`.
 */
export const ZodAssetTemplateCatalogListResponseSchema = z.object({
  status: IsNumber.int(),
  errors: z.array(z.any()).optional().nullable(),
  data: ZodAssetTemplateCatalogDataSchema,
});

/**
 * AssetTemplateFacetsQuery — query params for the marketplace facets endpoint.
 * The drill-down is resolved server-side: `vendor` narrows the returned models
 * and `model` narrows the returned versions. `lang` selects the facet labels.
 */
export const ZodAssetTemplateFacetsQuerySchema = z.object({
  vendor: IsString.optional(),
  model: IsString.optional(),
  lang: IsString.optional(),
});

/**
 * AssetTemplateFacetOption — a single selectable facet value with a display
 * label and an optional icon, used to build filter dropdowns.
 */
export const ZodAssetTemplateFacetOptionSchema = z.object({
  value: IsString,
  label: IsString,
  icon: IsString.optional().default(''),
});

/**
 * AssetTemplateFacets — the available filter options the listing UI renders:
 * the three drill-down levels (vendor -> model -> version) plus the flat
 * category facet.
 */
export const ZodAssetTemplateFacetsSchema = z.object({
  categories: z.array(ZodAssetTemplateFacetOptionSchema),
  vendors: z.array(ZodAssetTemplateFacetOptionSchema),
  models: z.array(ZodAssetTemplateFacetOptionSchema),
  versions: z.array(ZodAssetTemplateFacetOptionSchema),
});

/**
 * AssetTemplateFacetsResponse — the `{status,errors,data}` envelope returned by
 * `GET /api/v1/asset_templates/facets`.
 */
export const ZodAssetTemplateFacetsResponseSchema = z.object({
  status: IsNumber.int(),
  errors: z.array(z.any()).optional().nullable(),
  data: ZodAssetTemplateFacetsSchema,
});

/**
 * AssetTemplateBundleDynamicField — one dynamic field entry inside the bundle,
 * matching the entity DynamicField shape the preview table renders.
 */
export const ZodAssetTemplateBundleDynamicFieldSchema = z.object({
  fieldId: IsNumber.int().optional(),
  field: IsString,
  value: IsString.optional(),
  type: z.enum(['string', 'number', 'bool', 'date', 'geo']),
  status: IsNumber.int().optional(),
});

/**
 * AssetTemplateBundle — the `asset_template_information.json` shape the detail
 * modal renders. `name`/`description` are unresolved `{en-US,pt-BR}` locale
 * maps (the detail endpoint returns the full bundle, not a locale-resolved
 * card). `.passthrough()` keeps extra bundle fields (categoryId, etc.) from
 * failing validation.
 */
export const ZodAssetTemplateBundleSchema = z
  .object({
    name: z.record(IsString, IsString),
    description: z.record(IsString, IsString).optional(),
    categoryName: IsString.optional().default(''),
    manufacturerName: IsString.optional().default(''),
    modelName: IsString.optional().default(''),
    version: IsString.optional().default(''),
    assetIdPath: IsString.optional().default(''),
    scriptTest: IsString.optional(),
    scriptProcessor: IsString.optional(),
    scriptValidator: IsString.optional().default(''),
    scriptConversion: IsString.optional().default(''),
    dynamicFields: z.array(ZodAssetTemplateBundleDynamicFieldSchema).optional().default([]),
    availableFields: z.array(IsString).optional().default([]),
    nextFieldId: IsNumber.int().optional(),
  })
  .passthrough();

/**
 * AssetTemplateBundleResponse — the `{status,errors,data}` envelope returned by
 * `GET /api/v1/asset_templates/:vendor/:slug`.
 */
export const ZodAssetTemplateBundleResponseSchema = z.object({
  status: IsNumber.int(),
  errors: z.array(z.any()).optional().nullable(),
  data: ZodAssetTemplateBundleSchema,
});

/**
 * MarketplaceInstalledCheck — POST /api/v1/asset_templates/marketplace/installed.
 * Given a set of marketplace GUIDs, the response lists the subset the caller org
 * has installed, so the listing can show Install vs Uninstall per card in one call.
 */
export const ZodMarketplaceInstalledCheckRequestSchema = z.object({
  marketplaceGuids: z.array(IsString),
});

export const ZodMarketplaceInstalledCheckResponseSchema = z.object({
  installed: z.array(IsString),
});

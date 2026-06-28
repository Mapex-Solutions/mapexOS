import { z, IsString, IsBoolean, IsNumber, NumberIntAndPositive, StringAndNotBeEmpty } from '@mapexos/validations';

/**
 * PluginCatalogItem — a single FLAT list item from
 * `GET /api/v1/workflow_plugins`. The marketplace service has already
 * resolved `name`/`description` to plain strings via the `?lang` query
 * (they are NOT localized {en-US,pt-BR} objects here). The full,
 * unresolved manifest is fetched separately via the detail endpoint.
 */
export const ZodPluginCatalogItemSchema = z.object({
  id: StringAndNotBeEmpty,
  vendor: StringAndNotBeEmpty,
  vendorName: IsString.optional().default(''),
  pluginId: StringAndNotBeEmpty,
  slug: StringAndNotBeEmpty,
  name: IsString,
  description: IsString.optional().default(''),
  category: IsString.optional().default(''),
  capabilities: z.array(IsString).optional().default([]),
  tags: z.array(IsString).optional().default([]),
  icon: IsString.optional().default(''),
  color: IsString.optional().default(''),
  image: IsString.optional().default(''),
  requiresCredentials: IsBoolean.optional().default(false),
  nodeCount: IsNumber.int().optional().default(0),
  triggerCount: IsNumber.int().optional().default(0),
  hasEvents: IsBoolean.optional().default(false),
  hasImage: IsBoolean.optional().default(false),
});

/**
 * PluginCatalogQuery — query params for the marketplace list endpoint.
 * Filtering and pagination are server-side. `lang` selects which
 * localized strings the service resolves into the flat item fields.
 */
export const ZodPluginCatalogQuerySchema = z.object({
  search: IsString.optional(),
  category: IsString.optional(),
  capability: IsString.optional(),
  tag: IsString.optional(),
  page: NumberIntAndPositive.optional(),
  perPage: NumberIntAndPositive.max(100).optional(),
  lang: IsString.optional(),
});

/**
 * PluginCatalogData — the inner data body of the list envelope.
 * `total` is the full server-side count for pagination.
 */
export const ZodPluginCatalogDataSchema = z.object({
  items: z.array(ZodPluginCatalogItemSchema),
  total: IsNumber.int().min(0),
  page: IsNumber.int().min(0),
  perPage: IsNumber.int().min(0),
});

/**
 * PluginCatalogListResponse — the full `{status,errors,data}` envelope
 * returned by `GET /api/v1/workflow_plugins`.
 */
export const ZodPluginCatalogListResponseSchema = z.object({
  status: IsNumber.int(),
  errors: z.array(z.any()).optional().nullable(),
  data: ZodPluginCatalogDataSchema,
});

/**
 * PluginFacetOption — a single selectable facet value with a display
 * label and an optional icon, used to build filter dropdowns.
 */
export const ZodPluginFacetOptionSchema = z.object({
  value: IsString,
  label: IsString,
  icon: IsString.optional().default(''),
});

/**
 * PluginFacets — the available categories and capabilities used to
 * populate the marketplace filter controls.
 */
export const ZodPluginFacetsSchema = z.object({
  categories: z.array(ZodPluginFacetOptionSchema),
  capabilities: z.array(ZodPluginFacetOptionSchema),
});

/**
 * PluginFacetsResponse — the `{status,errors,data}` envelope returned by
 * `GET /api/v1/workflow_plugins/facets`.
 */
export const ZodPluginFacetsResponseSchema = z.object({
  status: IsNumber.int(),
  errors: z.array(z.any()).optional().nullable(),
  data: ZodPluginFacetsSchema,
});

/**
 * PluginManifestResponse — the wrapped `{status,errors,data}` envelope
 * returned by `GET /api/v1/workflow_plugins/:vendor/:slug`. The `data`
 * field is the raw plugin manifest; it is kept intentionally permissive
 * (passthrough) so the consumer can hand it to the existing manifest
 * converter without re-declaring the full manifest shape here.
 */
export const ZodPluginManifestResponseSchema = z.object({
  status: IsNumber.int(),
  errors: z.array(z.any()).optional().nullable(),
  data: z.record(IsString, z.any()),
});

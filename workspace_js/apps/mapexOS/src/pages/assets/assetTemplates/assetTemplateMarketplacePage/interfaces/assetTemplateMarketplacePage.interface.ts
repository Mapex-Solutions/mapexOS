/**
 * Active filter state for the asset template marketplace grid. `search` is a
 * free-text query; the remaining fields are the drill-down facets (a cleared
 * facet is `null`, meaning "no filter").
 */
export interface MarketplacePageFilters {
  /** Free-text search across template name/description. */
  search: string;
  /** Selected category facet value, or null for all. */
  category: string | null;
  /** Selected vendor facet value, or null for all. */
  vendor: string | null;
  /** Selected model facet value, or null for all. */
  model: string | null;
  /** Selected version facet value, or null for all. */
  version: string | null;
}

/**
 * The `{vendor, slug}` reference the grid hands to the detail modal when a card
 * is opened.
 */
export interface MarketplaceSelection {
  /** Template vendor key. */
  vendor: string;
  /** Template slug. */
  slug: string;
}

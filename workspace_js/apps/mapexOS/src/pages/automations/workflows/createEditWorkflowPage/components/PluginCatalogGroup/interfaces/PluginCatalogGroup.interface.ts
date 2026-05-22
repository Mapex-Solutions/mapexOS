import type { CatalogGroup } from '@src/components/workflow/interfaces';

/**
 * Props for PluginCatalogGroup component
 */
export interface PluginCatalogGroupProps {
  /** Catalog group to render */
  group: CatalogGroup;

  /** Whether the parent catalog is collapsed */
  collapsed: boolean;
}

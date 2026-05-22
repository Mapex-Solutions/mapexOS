import type { PluginNodeType } from '@src/components/workflow/interfaces';

/**
 * Props for PluginCatalogItem component
 */
export interface PluginCatalogItemProps {
  /** Node type to render */
  nodeType: PluginNodeType;

  /** Whether the catalog is collapsed */
  collapsed: boolean;
}

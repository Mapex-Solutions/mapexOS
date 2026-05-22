/**
 * Props for PluginCatalog component
 */
export interface PluginCatalogProps {
  /** Whether the catalog is collapsed */
  collapsed: boolean;
}

/**
 * Emits for PluginCatalog component
 */
export interface PluginCatalogEmits {
  (e: 'toggle-collapse'): void;
}
